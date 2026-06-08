package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDedupStrings(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  int
	}{
		{"empty", nil, 0},
		{"no duplicates", []string{"/a", "/b", "/c"}, 3},
		{"all duplicates", []string{"/a", "/a", "/a"}, 1},
		{"some duplicates", []string{"/a", "/b", "/a", "/c", "/b"}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dedupStrings(tt.input)
			if len(got) != tt.want {
				t.Errorf("dedupStrings() returned %d items, want %d: %v", len(got), tt.want, got)
			}
		})
	}
}

func TestDedupAssets(t *testing.T) {
	tests := []struct {
		name  string
		input []*ImageAsset
		want  int
	}{
		{"empty", nil, 0},
		{"no duplicates", []*ImageAsset{
			{AbsPath: "/a.png"},
			{AbsPath: "/b.png"},
		}, 2},
		{"all duplicates", []*ImageAsset{
			{AbsPath: "/a.png"},
			{AbsPath: "/a.png"},
		}, 1},
		{"some duplicates", []*ImageAsset{
			{AbsPath: "/a.png"},
			{AbsPath: "/b.png"},
			{AbsPath: "/a.png"},
		}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dedupAssets(tt.input)
			if len(got) != tt.want {
				t.Errorf("dedupAssets() returned %d items, want %d", len(got), tt.want)
			}
		})
	}
}

func TestCollectAssets(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "assets"), 0755)
	os.WriteFile(filepath.Join(dir, "assets", "logo.png"), []byte("png"), 0644)
	os.WriteFile(filepath.Join(dir, "assets", "banner.svg"), []byte("svg"), 0644)
	os.WriteFile(filepath.Join(dir, "assets", "notes.txt"), []byte("text"), 0644)

	s := NewScanner(Options{})
	assets, err := s.collectAssets(dir, filepath.Join(dir, "assets"))
	if err != nil {
		t.Fatal(err)
	}

	if len(assets) != 2 {
		t.Fatalf("expected 2 image files, got %d", len(assets))
	}

	found := map[string]bool{}
	for _, a := range assets {
		found[a.Basename] = true
	}
	if !found["logo.png"] {
		t.Error("missing logo.png")
	}
	if !found["banner.svg"] {
		t.Error("missing banner.svg")
	}
}

func TestCollectAssetsWithIgnoredDir(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "assets", "node_modules"), 0755)
	os.MkdirAll(filepath.Join(dir, "assets", "icons"), 0755)
	os.WriteFile(filepath.Join(dir, "assets", "logo.png"), []byte("png"), 0644)
	os.WriteFile(filepath.Join(dir, "assets", "node_modules", "ignore.png"), []byte("png"), 0644)
	os.WriteFile(filepath.Join(dir, "assets", "icons", "icon.svg"), []byte("svg"), 0644)

	s := NewScanner(Options{})
	assets, err := s.collectAssets(dir, filepath.Join(dir, "assets"))
	if err != nil {
		t.Fatal(err)
	}

	for _, a := range assets {
		if a.Basename == "ignore.png" {
			t.Error("should not have collected ignore.png from node_modules")
		}
	}

	if len(assets) != 2 {
		t.Fatalf("expected 2 assets (logo.png, icon.svg), got %d", len(assets))
	}
}

func TestCollectReferences(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "src"), 0755)
	os.WriteFile(filepath.Join(dir, "src", "App.tsx"), []byte(`
import logo from "./assets/logo.png"
const img = require("./assets/banner.svg")
<img src="/icons/icon.svg" />
`), 0644)

	s := NewScanner(Options{})
	refs, err := s.collectReferences(dir, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"logo.png", "banner.svg", "icon.svg"}
	for _, e := range expected {
		if !refs[e] {
			t.Errorf("missing reference %q", e)
		}
	}
}

func TestCollectReferencesSkipsAssetDirs(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "public"), 0755)
	os.WriteFile(filepath.Join(dir, "public", "ignored.tsx"), []byte(`import "./logo.png"`), 0644)

	os.MkdirAll(filepath.Join(dir, "src"), 0755)
	os.WriteFile(filepath.Join(dir, "src", "App.tsx"), []byte(`import "./logo.png"`), 0644)

	s := NewScanner(Options{})
	refs, err := s.collectReferences(dir, []string{filepath.Join(dir, "public")})
	if err != nil {
		t.Fatal(err)
	}

	// Only the reference from src/ should be counted
	if len(refs) != 1 {
		t.Errorf("expected 1 reference from src/, got %d", len(refs))
	}
}

func TestCollectReferencesSkipsIgnoredDirs(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "node_modules"), 0755)
	os.WriteFile(filepath.Join(dir, "node_modules", "pkg.tsx"), []byte(`import "./logo.png"`), 0644)

	os.MkdirAll(filepath.Join(dir, "src"), 0755)
	os.WriteFile(filepath.Join(dir, "src", "App.tsx"), []byte(`import "./logo.png"`), 0644)

	s := NewScanner(Options{IgnoreDirs: []string{"extra"}})
	refs, err := s.collectReferences(dir, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(refs) != 1 {
		t.Errorf("expected 1 reference (src/ only), got %d", len(refs))
	}
}

func TestScanEndToEnd(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "src"), 0755)
	os.WriteFile(filepath.Join(dir, "src", "App.tsx"), []byte(`
import logo from "./assets/logo.png"
<img src="/assets/banner.svg" />
`), 0644)

	os.MkdirAll(filepath.Join(dir, "assets"), 0755)
	os.WriteFile(filepath.Join(dir, "assets", "logo.png"), []byte("png"), 0644)
	os.WriteFile(filepath.Join(dir, "assets", "banner.svg"), []byte("svg"), 0644)
	os.WriteFile(filepath.Join(dir, "assets", "unused.gif"), []byte("gif"), 0644)

	s := NewScanner(Options{
		ProjectPath: dir,
		AssetsDirs:  []string{"assets"},
	})
	result, err := s.Scan()
	if err != nil {
		t.Fatal(err)
	}

	if result.TotalAssets != 3 {
		t.Errorf("expected 3 total assets, got %d", result.TotalAssets)
	}
	if result.UsedAssets != 2 {
		t.Errorf("expected 2 used assets, got %d", result.UsedAssets)
	}
	if result.UnusedAssets != 1 {
		t.Errorf("expected 1 unused asset, got %d", result.UnusedAssets)
	}

	used := map[string]bool{}
	for _, a := range result.Assets {
		if a.Used {
			used[a.Basename] = true
		}
	}
	if !used["logo.png"] {
		t.Error("logo.png should be used")
	}
	if !used["banner.svg"] {
		t.Error("banner.svg should be used")
	}

	unused := result.UnusedList()
	if len(unused) != 1 || unused[0].Basename != "unused.gif" {
		t.Errorf("unused should be [unused.gif], got %v", unused)
	}
}

func TestScanDuplicateAssetDirs(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "src"), 0755)
	os.WriteFile(filepath.Join(dir, "src", "App.tsx"), []byte(`import "./assets/logo.png"`), 0644)

	os.MkdirAll(filepath.Join(dir, "assets"), 0755)
	os.WriteFile(filepath.Join(dir, "assets", "logo.png"), []byte("png"), 0644)

	s := NewScanner(Options{
		ProjectPath: dir,
		AssetsDirs:  []string{"assets", "assets"},
	})
	result, err := s.Scan()
	if err != nil {
		t.Fatal(err)
	}

	if result.TotalAssets != 1 {
		t.Errorf("expected 1 asset after dedup, got %d", result.TotalAssets)
	}
	if len(result.AssetsDirs) != 1 {
		t.Errorf("expected 1 asset dir after dedup, got %d", len(result.AssetsDirs))
	}
}

func TestScanProjectPathNotExist(t *testing.T) {
	s := NewScanner(Options{
		ProjectPath: "/nonexistent/path",
	})
	_, err := s.Scan()
	if err == nil {
		t.Fatal("expected error for nonexistent project path")
	}
}

func TestScanAssetsDirNotExist(t *testing.T) {
	dir := t.TempDir()

	s := NewScanner(Options{
		ProjectPath: dir,
		AssetsDirs:  []string{"nonexistent"},
	})
	_, err := s.Scan()
	if err == nil {
		t.Fatal("expected error for nonexistent assets dir")
	}
}

func TestScanWithIgnoreDirs(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "src"), 0755)
	os.WriteFile(filepath.Join(dir, "src", "App.tsx"), []byte(`import "./assets/logo.png"`), 0644)

	os.MkdirAll(filepath.Join(dir, "assets"), 0755)
	os.WriteFile(filepath.Join(dir, "assets", "logo.png"), []byte("png"), 0644)

	os.MkdirAll(filepath.Join(dir, "node_modules"), 0755)
	os.WriteFile(filepath.Join(dir, "node_modules", "pkg.js"), []byte(`import "./assets/logo.png"`), 0644)

	s := NewScanner(Options{
		ProjectPath: dir,
		AssetsDirs:  []string{"assets"},
	})
	result, err := s.Scan()
	if err != nil {
		t.Fatal(err)
	}

	if result.TotalAssets != 1 {
		t.Errorf("expected 1 asset, got %d", result.TotalAssets)
	}
}

func TestScanEmptyAssetsDir(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "assets"), 0755)
	os.MkdirAll(filepath.Join(dir, "src"), 0755)
	os.WriteFile(filepath.Join(dir, "src", "App.tsx"), []byte(`const x = 1`), 0644)

	s := NewScanner(Options{
		ProjectPath: dir,
		AssetsDirs:  []string{"assets"},
	})
	result, err := s.Scan()
	if err != nil {
		t.Fatal(err)
	}

	if result.TotalAssets != 0 {
		t.Errorf("expected 0 assets, got %d", result.TotalAssets)
	}
	if result.UsedAssets != 0 {
		t.Errorf("expected 0 used, got %d", result.UsedAssets)
	}
}

func TestScanMultipleAssetDirs(t *testing.T) {
	dir := t.TempDir()

	os.MkdirAll(filepath.Join(dir, "src"), 0755)
	os.WriteFile(filepath.Join(dir, "src", "App.tsx"), []byte(`
import logo from "./public/logo.png"
import icon from "./icons/icon.svg"
`), 0644)

	os.MkdirAll(filepath.Join(dir, "public"), 0755)
	os.WriteFile(filepath.Join(dir, "public", "logo.png"), []byte("png"), 0644)

	os.MkdirAll(filepath.Join(dir, "icons"), 0755)
	os.WriteFile(filepath.Join(dir, "icons", "icon.svg"), []byte("svg"), 0644)
	os.WriteFile(filepath.Join(dir, "icons", "unused.gif"), []byte("gif"), 0644)

	s := NewScanner(Options{
		ProjectPath: dir,
		AssetsDirs:  []string{"public", "icons"},
	})
	result, err := s.Scan()
	if err != nil {
		t.Fatal(err)
	}

	if result.TotalAssets != 3 {
		t.Errorf("expected 3 total assets, got %d", result.TotalAssets)
	}
	if result.UsedAssets != 2 {
		t.Errorf("expected 2 used, got %d", result.UsedAssets)
	}
	if result.UnusedAssets != 1 {
		t.Errorf("expected 1 unused, got %d", result.UnusedAssets)
	}
}
