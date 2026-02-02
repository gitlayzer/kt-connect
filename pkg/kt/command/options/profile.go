package options

import (
	"fmt"
	"github.com/gitlayzer/kt-connect/pkg/kt/util"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Profile struct {
	Name                     string `yaml:"name"`
	KubeContext              string `yaml:"kubeContext,omitempty"`
	Namespace                string `yaml:"namespace,omitempty"`
	ProxyMode                string `yaml:"proxyMode,omitempty"`
	ExchangeMode             string `yaml:"exchangeMode,omitempty"`
	ExchangeExpose           string `yaml:"exchangeExpose,omitempty"`
	ExchangeSkipPortChecking *bool  `yaml:"exchangeSkipPortChecking,omitempty"`
	MeshMode                 string `yaml:"meshMode,omitempty"`
	MeshExpose               string `yaml:"meshExpose,omitempty"`
	MeshSkipPortChecking     *bool  `yaml:"meshSkipPortChecking,omitempty"`
	PreviewExpose            string `yaml:"previewExpose,omitempty"`
	PreviewExternal          *bool  `yaml:"previewExternal,omitempty"`
	PreviewSkipPortChecking  *bool  `yaml:"previewSkipPortChecking,omitempty"`
}

func profileDir() string {
	return filepath.Join(util.KtProfileDir, "profiles")
}

func profilePath(name string) string {
	return filepath.Join(profileDir(), fmt.Sprintf("%s.yaml", name))
}

func currentProfilePath() string {
	return filepath.Join(profileDir(), "current")
}

func ensureProfileDir() error {
	return os.MkdirAll(profileDir(), 0755)
}

func SaveProfile(profile Profile) error {
	if profile.Name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}
	if err := ensureProfileDir(); err != nil {
		return err
	}
	data, err := yaml.Marshal(profile)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(profilePath(profile.Name), data, 0644)
}

func LoadProfile(name string) (*Profile, error) {
	if name == "" {
		return nil, fmt.Errorf("profile name cannot be empty")
	}
	data, err := ioutil.ReadFile(profilePath(name))
	if err != nil {
		return nil, err
	}
	var profile Profile
	if err := yaml.Unmarshal(data, &profile); err != nil {
		return nil, err
	}
	if profile.Name == "" {
		profile.Name = name
	}
	return &profile, nil
}

func DeleteProfile(name string) error {
	if name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}
	return os.Remove(profilePath(name))
}

func ListProfileNames() ([]string, error) {
	dirEntries, err := ioutil.ReadDir(profileDir())
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	names := make([]string, 0, len(dirEntries))
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}
		if entry.Name() == "current" {
			continue
		}
		if filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}
		names = append(names, strings.TrimSuffix(entry.Name(), ".yaml"))
	}
	return names, nil
}

func ProfileExists(name string) (bool, error) {
	if name == "" {
		return false, fmt.Errorf("profile name cannot be empty")
	}
	_, err := os.Stat(profilePath(name))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func SetCurrentProfileName(name string) error {
	if name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}
	if err := ensureProfileDir(); err != nil {
		return err
	}
	return ioutil.WriteFile(currentProfilePath(), []byte(name), 0644)
}

func ClearCurrentProfile() error {
	err := os.Remove(currentProfilePath())
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func CurrentProfileName() (string, bool, error) {
	data, err := ioutil.ReadFile(currentProfilePath())
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, err
	}
	name := strings.TrimSpace(string(data))
	if name == "" {
		return "", false, nil
	}
	return name, true, nil
}

func LoadCurrentProfile() (*Profile, error) {
	name, ok, err := CurrentProfileName()
	if err != nil || !ok {
		return nil, err
	}
	return LoadProfile(name)
}

func ApplyProfileDefaults(opt *DaemonOptions, profile *Profile) {
	if profile == nil {
		return
	}
	if profile.KubeContext != "" {
		opt.Global.Context = profile.KubeContext
	}
	if profile.Namespace != "" {
		opt.Global.Namespace = profile.Namespace
	}
	if profile.ProxyMode != "" {
		opt.Connect.Mode = profile.ProxyMode
	}
	if profile.ExchangeMode != "" {
		opt.Exchange.Mode = profile.ExchangeMode
	}
	if profile.ExchangeExpose != "" {
		opt.Exchange.Expose = profile.ExchangeExpose
	}
	if profile.ExchangeSkipPortChecking != nil {
		opt.Exchange.SkipPortChecking = *profile.ExchangeSkipPortChecking
	}
	if profile.MeshMode != "" {
		opt.Mesh.Mode = profile.MeshMode
	}
	if profile.MeshExpose != "" {
		opt.Mesh.Expose = profile.MeshExpose
	}
	if profile.MeshSkipPortChecking != nil {
		opt.Mesh.SkipPortChecking = *profile.MeshSkipPortChecking
	}
	if profile.PreviewExpose != "" {
		opt.Preview.Expose = profile.PreviewExpose
	}
	if profile.PreviewExternal != nil {
		opt.Preview.External = *profile.PreviewExternal
	}
	if profile.PreviewSkipPortChecking != nil {
		opt.Preview.SkipPortChecking = *profile.PreviewSkipPortChecking
	}
}
