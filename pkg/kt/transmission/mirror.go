package transmission

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gitlayzer/kt-connect/pkg/kt/util"
	"github.com/rs/zerolog/log"
	"io"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

const mirrorMaxPayloadBytes = 1024 * 1024

type MirrorConfig struct {
	Target       string
	SampleRate   int
	RedactRules  string
	LogPath      string
	LocalAddress string
}

type MirrorLogEntry struct {
	Timestamp  string `json:"timestamp"`
	RemoteAddr string `json:"remoteAddr"`
	LocalAddr  string `json:"localAddr"`
	Payload    string `json:"payload"`
	Truncated  bool   `json:"truncated"`
	Redacted   bool   `json:"redacted"`
}

type mirrorRedactRule struct {
	pattern     *regexp.Regexp
	replacement string
}

func (m MirrorConfig) Enabled() bool {
	return m.Target != "" || m.LogPath != ""
}

func (m MirrorConfig) normalizedSampleRate() int {
	if m.SampleRate <= 0 {
		return 0
	}
	if m.SampleRate > 100 {
		return 100
	}
	return m.SampleRate
}

func (m MirrorConfig) shouldSample() bool {
	rate := m.normalizedSampleRate()
	if rate == 0 {
		return false
	}
	if rate == 100 {
		return true
	}
	return rand.Intn(100) < rate
}

func (m MirrorConfig) parseRedactRules() []mirrorRedactRule {
	if strings.TrimSpace(m.RedactRules) == "" {
		return nil
	}
	rules := strings.Split(m.RedactRules, ";")
	parsed := make([]mirrorRedactRule, 0, len(rules))
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		parts := strings.SplitN(rule, "=", 2)
		if len(parts) != 2 {
			log.Warn().Msgf("Invalid mirror redact rule: %s", rule)
			continue
		}
		re, err := regexp.Compile(parts[0])
		if err != nil {
			log.Warn().Err(err).Msgf("Invalid mirror redact regex: %s", parts[0])
			continue
		}
		parsed = append(parsed, mirrorRedactRule{pattern: re, replacement: parts[1]})
	}
	return parsed
}

func (m MirrorConfig) applyRedaction(payload []byte, rules []mirrorRedactRule) ([]byte, bool) {
	if len(rules) == 0 || len(payload) == 0 {
		return payload, false
	}
	redacted := false
	content := string(payload)
	for _, rule := range rules {
		newContent := rule.pattern.ReplaceAllString(content, rule.replacement)
		if newContent != content {
			redacted = true
		}
		content = newContent
	}
	return []byte(content), redacted
}

type mirrorRecorder struct {
	buf       bytes.Buffer
	remaining int
	truncated bool
	mu        sync.Mutex
}

func newMirrorRecorder(limit int) *mirrorRecorder {
	return &mirrorRecorder{remaining: limit}
}

func (r *mirrorRecorder) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.remaining <= 0 {
		r.truncated = true
		return len(p), nil
	}
	if len(p) > r.remaining {
		r.buf.Write(p[:r.remaining])
		r.remaining = 0
		r.truncated = true
		return len(p), nil
	}
	r.buf.Write(p)
	r.remaining -= len(p)
	return len(p), nil
}

func (r *mirrorRecorder) Bytes() []byte {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]byte(nil), r.buf.Bytes()...)
}

func (r *mirrorRecorder) Truncated() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.truncated
}

func StartMirrorProxy(localPort int, mirror MirrorConfig) (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return -1, err
	}
	addr := listener.Addr().(*net.TCPAddr)
	proxyPort := addr.Port
	log.Info().Msgf("Mirror proxy listening on 127.0.0.1:%d for local port %d", proxyPort, localPort)

	rand.Seed(time.Now().UnixNano())
	go func() {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Warn().Err(err).Msgf("Mirror proxy accept failed")
				return
			}
			go handleMirrorConnection(conn, localPort, mirror)
		}
	}()
	return proxyPort, nil
}

func handleMirrorConnection(client net.Conn, localPort int, mirror MirrorConfig) {
	defer client.Close()
	localConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", localPort))
	if err != nil {
		log.Error().Err(err).Msgf("Mirror proxy failed to connect to local service")
		return
	}
	defer localConn.Close()

	shouldSample := mirror.shouldSample()
	recorder := newMirrorRecorder(mirrorMaxPayloadBytes)
	rules := mirror.parseRedactRules()

	done := make(chan struct{}, 2)
	go func() {
		reader := io.TeeReader(client, recorder)
		if _, err := io.Copy(localConn, reader); err != nil {
			log.Debug().Err(err).Msgf("Mirror proxy copy client->local interrupted")
		}
		done <- struct{}{}
	}()
	go func() {
		if _, err := io.Copy(client, localConn); err != nil {
			log.Debug().Err(err).Msgf("Mirror proxy copy local->client interrupted")
		}
		done <- struct{}{}
	}()

	<-done
	payload := recorder.Bytes()
	truncated := recorder.Truncated()

	if shouldSample && len(payload) > 0 {
		redactedPayload, redacted := mirror.applyRedaction(payload, rules)
		if mirror.Target != "" {
			if err := mirrorToTarget(mirror.Target, redactedPayload); err != nil {
				log.Warn().Err(err).Msgf("Mirror to target failed")
			}
		}
		if mirror.LogPath != "" {
			if err := writeMirrorLog(mirror.LogPath, MirrorLogEntry{
				Timestamp:  util.GetTimestamp(),
				RemoteAddr: client.RemoteAddr().String(),
				LocalAddr:  mirror.LocalAddress,
				Payload:    base64.StdEncoding.EncodeToString(redactedPayload),
				Truncated:  truncated,
				Redacted:   redacted,
			}); err != nil {
				log.Warn().Err(err).Msgf("Mirror log write failed")
			}
		}
	}
}

func mirrorToTarget(target string, payload []byte) error {
	conn, err := net.Dial("tcp", target)
	if err != nil {
		return err
	}
	defer conn.Close()
	if len(payload) == 0 {
		return nil
	}
	_, err = conn.Write(payload)
	return err
}

func writeMirrorLog(path string, entry MirrorLogEntry) error {
	if err := util.CreateDirIfNotExist(path); err != nil {
		return err
	}
	fileName := fmt.Sprintf("mirror-%s-%s.json", util.GetTimestamp(), strings.ToLower(util.RandomString(6)))
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	fullPath := filepath.Join(path, fileName)
	return os.WriteFile(fullPath, data, 0644)
}
