package cleaner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devrapture/vole/internal/scanner"
)

func TestCleanNoUnused(t *testing.T) {
	result, err := scanner.NewScanner(scanner.Options{
		ProjectPath: t.TempDir(),
	}).Scan()
	if err != nil {
		t.Fatal(err)
	}

	cleanResult, err := Clean(result, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if len(cleanResult.Deleted) != 0 {
		t.Errorf("expected 0 deleted, got %d", len(cleanResult.Deleted))
	}
	if cleanResult.SpaceSavedBytes != 0 {
		t.Errorf("expected 0 bytes saved, got %d", cleanResult.SpaceSavedBytes)
	}
}

func TestCleanDryRun(t *testing.T) {
	dir := t.TempDir()
	assetDir := filepath.Join(dir, "assets")
	os.MkdirAll(assetDir, 0755)
	filePath := filepath.Join(assetDir, "unused.png")
	os.WriteFile(filePath, []byte("fake png"), 0644)

	result := &scanner.ScanResult{
		ProjectPath: dir,
		Assets: []*scanner.ImageAsset{
			{AbsPath: filePath, RelPath: "assets/unused.png", Basename: "unused.png", SizeBytes: 8},
		},
		TotalAssets:  1,
		UnusedAssets: 1,
	}

	cleanResult, err := Clean(result, Options{DryRun: true})
	if err != nil {
		t.Fatal(err)
	}

	if len(cleanResult.Deleted) != 0 {
		t.Errorf("dry-run should not delete files, got %d deleted", len(cleanResult.Deleted))
	}
	if len(cleanResult.Skipped) != 1 {
		t.Errorf("dry-run should skip 1 file, got %d", len(cleanResult.Skipped))
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("dry-run should not remove the file")
	}
}

func TestCleanDeletesFiles(t *testing.T) {
	dir := t.TempDir()
	assetDir := filepath.Join(dir, "assets")
	os.MkdirAll(assetDir, 0755)
	filePath := filepath.Join(assetDir, "unused.png")
	os.WriteFile(filePath, []byte("fake png"), 0644)

	result := &scanner.ScanResult{
		ProjectPath: dir,
		Assets: []*scanner.ImageAsset{
			{AbsPath: filePath, RelPath: "assets/unused.png", Basename: "unused.png", SizeBytes: 8},
		},
		TotalAssets:  1,
		UnusedAssets: 1,
	}

	cleanResult, err := Clean(result, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if len(cleanResult.Deleted) != 1 {
		t.Fatalf("expected 1 deleted, got %d", len(cleanResult.Deleted))
	}
	if cleanResult.Deleted[0] != filePath {
		t.Errorf("expected deleted path %q, got %q", filePath, cleanResult.Deleted[0])
	}
	if cleanResult.SpaceSavedBytes != 8 {
		t.Errorf("expected 8 bytes saved, got %d", cleanResult.SpaceSavedBytes)
	}
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("file should have been removed")
	}
}

func TestCleanDeleteError(t *testing.T) {
	dir := t.TempDir()
	nonexistent := filepath.Join(dir, "nonexistent.png")

	result := &scanner.ScanResult{
		ProjectPath: dir,
		Assets: []*scanner.ImageAsset{
			{AbsPath: nonexistent, RelPath: "nonexistent.png", Basename: "nonexistent.png"},
		},
		TotalAssets:  1,
		UnusedAssets: 1,
	}

	cleanResult, err := Clean(result, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if len(cleanResult.Deleted) != 0 {
		t.Errorf("expected 0 deleted for nonexistent file, got %d", len(cleanResult.Deleted))
	}
	if len(cleanResult.Errors) != 1 {
		t.Errorf("expected 1 error for nonexistent file, got %d", len(cleanResult.Errors))
	}
}

func TestCleanMultipleFiles(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "assets"), 0755)

	files := []struct {
		name string
		size int64
	}{
		{"a.png", 100},
		{"b.png", 200},
		{"c.png", 300},
	}
	var assets []*scanner.ImageAsset
	for _, f := range files {
		p := filepath.Join(dir, "assets", f.name)
		os.WriteFile(p, make([]byte, f.size), 0644)
		assets = append(assets, &scanner.ImageAsset{
			AbsPath: p, RelPath: "assets/" + f.name, Basename: f.name, SizeBytes: f.size,
		})
	}

	result := &scanner.ScanResult{
		ProjectPath:  dir,
		Assets:       assets,
		TotalAssets:  3,
		UnusedAssets: 3,
	}

	cleanResult, err := Clean(result, Options{})
	if err != nil {
		t.Fatal(err)
	}

	if len(cleanResult.Deleted) != 3 {
		t.Fatalf("expected 3 deleted, got %d", len(cleanResult.Deleted))
	}
	if cleanResult.SpaceSavedBytes != 600 {
		t.Errorf("expected 600 bytes saved, got %d", cleanResult.SpaceSavedBytes)
	}
}

func TestCleanVerbose(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "unused.png")
	os.WriteFile(p, []byte("data"), 0644)

	result := &scanner.ScanResult{
		ProjectPath: dir,
		Assets: []*scanner.ImageAsset{
			{AbsPath: p, RelPath: "unused.png", Basename: "unused.png", SizeBytes: 4},
		},
		UnusedAssets: 1,
	}

	cleanResult, err := Clean(result, Options{Verbose: true})
	if err != nil {
		t.Fatal(err)
	}

	if len(cleanResult.Deleted) != 1 {
		t.Errorf("expected 1 deleted, got %d", len(cleanResult.Deleted))
	}
}
