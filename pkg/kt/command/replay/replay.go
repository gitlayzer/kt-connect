package replay

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gitlayzer/kt-connect/pkg/kt/transmission"
	"github.com/rs/zerolog/log"
	"net"
	"os"
	"path/filepath"
	"sort"
)

func Replay(logPath, target string) error {
	if logPath == "" {
		return fmt.Errorf("mirror log path is required")
	}
	if target == "" {
		return fmt.Errorf("target address is required")
	}
	files, err := collectLogFiles(logPath)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no mirror logs found in %s", logPath)
	}
	for _, file := range files {
		entry, err := readMirrorLog(file)
		if err != nil {
			return err
		}
		payload, err := base64.StdEncoding.DecodeString(entry.Payload)
		if err != nil {
			return fmt.Errorf("invalid mirror log payload in %s: %w", file, err)
		}
		if len(payload) == 0 {
			continue
		}
		if err := sendPayload(target, payload); err != nil {
			return err
		}
		log.Info().Msgf("Replayed mirror log %s to %s", filepath.Base(file), target)
	}
	return nil
}

func collectLogFiles(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return []string{path}, nil
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, filepath.Join(path, entry.Name()))
	}
	sort.Strings(files)
	return files, nil
}

func readMirrorLog(path string) (*transmission.MirrorLogEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entry transmission.MirrorLogEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("invalid mirror log %s: %w", path, err)
	}
	return &entry, nil
}

func sendPayload(target string, payload []byte) error {
	conn, err := net.Dial("tcp", target)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(payload)
	return err
}
