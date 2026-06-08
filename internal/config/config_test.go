package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadNoConfig(t *testing.T) {
	dir := t.TempDir()

	cfg, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if len(cfg.Assets) != 0 {
		t.Errorf("expected empty assets, got %v", cfg.Assets)
	}
	if len(cfg.Ignore) != 0 {
		t.Errorf("expected empty ignore, got %v", cfg.Ignore)
	}
}

func TestLoadWithYml(t *testing.T) {
	dir := t.TempDir()

	content := []byte("assets:\n  - public\n  - assets\nignore:\n  - node_modules\n  - dist\n")
	os.WriteFile(filepath.Join(dir, "vole.yml"), content, 0644)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}

	expectedAssets := []string{"public", "assets"}
	expectedIgnore := []string{"node_modules", "dist"}

	if len(cfg.Assets) != len(expectedAssets) {
		t.Errorf("expected assets %v, got %v", expectedAssets, cfg.Assets)
	}
	for i, a := range expectedAssets {
		if cfg.Assets[i] != a {
			t.Errorf("asset[%d] = %q, want %q", i, cfg.Assets[i], a)
		}
	}
	if len(cfg.Ignore) != len(expectedIgnore) {
		t.Errorf("expected ignore %v, got %v", expectedIgnore, cfg.Ignore)
	}
	for i, ig := range expectedIgnore {
		if cfg.Ignore[i] != ig {
			t.Errorf("ignore[%d] = %q, want %q", i, cfg.Ignore[i], ig)
		}
	}
}

func TestLoadWithYaml(t *testing.T) {
	dir := t.TempDir()

	content := []byte("assets:\n  - public\nignore:\n  - custom\n")
	os.WriteFile(filepath.Join(dir, "vole.yaml"), content, 0644)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(cfg.Assets) != 1 || cfg.Assets[0] != "public" {
		t.Errorf("expected assets [public], got %v", cfg.Assets)
	}
	if len(cfg.Ignore) != 1 || cfg.Ignore[0] != "custom" {
		t.Errorf("expected ignore [custom], got %v", cfg.Ignore)
	}
}

func TestLoadPrefersYmlOverYaml(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "vole.yml"), []byte("assets:\n  - from_yml\n"), 0644)
	os.WriteFile(filepath.Join(dir, "vole.yaml"), []byte("assets:\n  - from_yaml\n"), 0644)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(cfg.Assets) != 1 || cfg.Assets[0] != "from_yml" {
		t.Errorf("expected assets from vole.yml (preferred), got %v", cfg.Assets)
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "vole.yml"), []byte("invalid: yaml: [\n"), 0644)

	_, err := Load(dir)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestLoadConfigFileError(t *testing.T) {
	dir := t.TempDir()

	// Create a directory with the same name as the config to cause a stat error
	os.MkdirAll(filepath.Join(dir, "vole.yml"), 0755)

	_, err := Load(dir)
	if err == nil {
		t.Fatal("expected error when vole.yml is a directory")
	}
}

func TestFindConfigPathNotFound(t *testing.T) {
	dir := t.TempDir()

	path, err := findConfigPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	if path != "" {
		t.Errorf("expected empty path, got %q", path)
	}
}

func TestFindConfigPathFound(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "vole.yml"), []byte{}, 0644)

	path, err := findConfigPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	if path == "" {
		t.Fatal("expected non-empty path")
	}
}

func TestFindConfigPathYaml(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "vole.yaml"), []byte{}, 0644)

	path, err := findConfigPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	if path == "" {
		t.Fatal("expected non-empty path")
	}
}
