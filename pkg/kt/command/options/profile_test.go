package options

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gitlayzer/kt-connect/pkg/kt/util"
)

func TestProfileLifecycle(t *testing.T) {
	baseDir := t.TempDir()
	originalProfileDir := util.KtProfileDir
	util.KtProfileDir = filepath.Join(baseDir, "profile")
	t.Cleanup(func() {
		util.KtProfileDir = originalProfileDir
	})

	profile := Profile{
		Name:         "dev",
		KubeContext:  "kc-dev",
		Namespace:    "default",
		ProxyMode:    "tun2socks",
		ExchangeMode: "selector",
	}
	if err := SaveProfile(profile); err != nil {
		t.Fatalf("save profile: %v", err)
	}

	if exists, err := ProfileExists("dev"); err != nil || !exists {
		t.Fatalf("profile exists: %v, %v", exists, err)
	}

	loaded, err := LoadProfile("dev")
	if err != nil {
		t.Fatalf("load profile: %v", err)
	}
	if loaded.Name != "dev" || loaded.KubeContext != "kc-dev" || loaded.Namespace != "default" {
		t.Fatalf("unexpected profile: %#v", loaded)
	}

	if err := SetCurrentProfileName("dev"); err != nil {
		t.Fatalf("set current profile: %v", err)
	}
	current, ok, err := CurrentProfileName()
	if err != nil || !ok || current != "dev" {
		t.Fatalf("current profile: %v, %v, %v", current, ok, err)
	}

	names, err := ListProfileNames()
	if err != nil {
		t.Fatalf("list profiles: %v", err)
	}
	if !reflect.DeepEqual(names, []string{"dev"}) {
		t.Fatalf("unexpected profile list: %#v", names)
	}

	if err := DeleteProfile("dev"); err != nil {
		t.Fatalf("delete profile: %v", err)
	}
}

func TestApplyProfileDefaults(t *testing.T) {
	opt := &DaemonOptions{
		Global:   &GlobalOptions{},
		Connect:  &ConnectOptions{},
		Exchange: &ExchangeOptions{},
		Mesh:     &MeshOptions{},
		Preview:  &PreviewOptions{},
		Forward:  &ForwardOptions{},
		Recover:  &RecoverOptions{},
		Clean:    &CleanOptions{},
		Birdseye: &BirdseyeOptions{},
		Config:   &ConfigOptions{},
	}
	previewExternal := true
	profile := &Profile{
		KubeContext:     "kc-dev",
		Namespace:       "default",
		ProxyMode:       "tun2socks",
		ExchangeMode:    "selector",
		ExchangeExpose:  "8080",
		MeshMode:        "auto",
		MeshExpose:      "8081",
		PreviewExpose:   "8082",
		PreviewExternal: &previewExternal,
	}

	ApplyProfileDefaults(opt, profile)

	if opt.Global.Context != "kc-dev" || opt.Global.Namespace != "default" {
		t.Fatalf("global defaults not applied")
	}
	if opt.Connect.Mode != "tun2socks" {
		t.Fatalf("connect defaults not applied")
	}
	if opt.Exchange.Mode != "selector" || opt.Exchange.Expose != "8080" {
		t.Fatalf("exchange defaults not applied")
	}
	if opt.Mesh.Mode != "auto" || opt.Mesh.Expose != "8081" {
		t.Fatalf("mesh defaults not applied")
	}
	if opt.Preview.Expose != "8082" || !opt.Preview.External {
		t.Fatalf("preview defaults not applied")
	}
}
