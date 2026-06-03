package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	ProjectPath string
	AssetsDir   string
	IgnoreDirs  []string
	verbose     bool
}

type Scanner struct {
	opts Options
}

func NewScanner(opts Options) *Scanner {
	return &Scanner{
		opts: opts,
	}
}

func (s *Scanner) Scan() (*ScanResult, error) {
	projectAbsPath, err := filepath.Abs(s.opts.ProjectPath)
	if err != nil {
		return nil, fmt.Errorf("resolving project path: %w", err)
	}
	if _, err := os.Stat(projectAbsPath); err != nil {
		return nil, fmt.Errorf("project path does not exist: %s", projectAbsPath)
	}

	assetAbsPath := filepath.Join(projectAbsPath, s.opts.AssetsDir)
	if _, err := os.Stat(assetAbsPath); err != nil {
		return nil, fmt.Errorf("assets directory not found: %s\n set --assets to the correct sub-path", assetAbsPath)
	}

	imageAsset, err := s.collectAssets(projectAbsPath, assetAbsPath)
	if err != nil {
		return nil, fmt.Errorf("collecting assets: %w", err)
	}

	refs, err := s.collectReferences(projectAbsPath, assetAbsPath)
	if err != nil {
		return nil, fmt.Errorf("collecting references: %w", err)
	}

	usedCount := 0
	for _, a := range imageAsset {
		if refs[strings.ToLower(a.Basename)] {
			a.Used = true
			usedCount++
		}
	}

	return &ScanResult{
		ProjectPath:  projectAbsPath,
		AssetsDir:    assetAbsPath,
		TotalAssets:  len(imageAsset),
		UsedAssets:   usedCount,
		UnusedAssets: (len(imageAsset) - usedCount),
		Assets:       imageAsset,
	}, nil
}

// collectAssets walks the assets directory and returns every image file found.
func (s *Scanner) collectAssets(projectAbsPath, assetAbsPath string) ([]*ImageAsset, error) {
	var imageAssets []*ImageAsset

	err := fs.WalkDir(os.DirFS(assetAbsPath), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != "." && IsIgnoredDirs(d.Name(), s.opts.IgnoreDirs) {
				return fs.SkipDir
			}
		}

		if !IsImageFile(path) {
			return nil
		}

		absPath := filepath.Join(assetAbsPath, path)

		relPath, err := filepath.Rel(projectAbsPath, absPath)
		if err != nil {
			return fmt.Errorf("computing relative path for %s: %w", absPath, err)
		}

		imageAssets = append(imageAssets, &ImageAsset{
			AbsPath:  absPath,
			RelPath:  filepath.ToSlash(relPath),
			Basename: strings.ToLower(filepath.Base(path)),
		})

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("collecting assets: %w", err)
	}
	return imageAssets, nil
}

func (s *Scanner) collectReferences(projectAbsPath, assetAbsPath string) (map[string]bool, error) {
	refs := make(map[string]bool)
	err := fs.WalkDir(os.DirFS(projectAbsPath), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path == "." {
				return nil
			}

			if IsIgnoredDirs(d.Name(), s.opts.IgnoreDirs) {
				return fs.SkipDir
			}

			absDir := filepath.Join(projectAbsPath, path)
			if absDir == assetAbsPath {
				return fs.SkipDir
			}
			return nil
		}

		if !IsSourceFile(path) {
			return nil
		}

		absPath := filepath.Join(projectAbsPath, path)
		if s.opts.verbose {
			fmt.Printf("vole reading: %/s\n", path)
		}

		content, err := os.ReadFile(absPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "vole warning: could not read %s: %v\n", path, err)
			return nil
		}

		for basename := range ExtractReferences(string(content)) {
			refs[basename] = true
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return refs, nil
}
