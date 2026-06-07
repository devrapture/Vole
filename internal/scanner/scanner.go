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
	AssetsDirs  []string
	IgnoreDirs  []string
	Verbose     bool
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

	var assetAbsPaths []string
	var imageAssets []*ImageAsset

	for _, assetsDir := range s.opts.AssetsDirs {
		assetAbsPath := filepath.Join(projectAbsPath, assetsDir)

		if _, err := os.Stat(assetAbsPath); err != nil {
			return nil, fmt.Errorf("assets directory not found: %s\n set --assets or vole.yml correctly", assetAbsPath)
		}

		assets, err := s.collectAssets(projectAbsPath, assetAbsPath)
		if err != nil {
			return nil, fmt.Errorf("collecting assets from %s: %w", assetAbsPath, err)
		}

		assetAbsPaths = append(assetAbsPaths, assetAbsPath)
		imageAssets = append(imageAssets, assets...)
	}

	refs, err := s.collectReferences(projectAbsPath, assetAbsPaths)

	if err != nil {
		return nil, fmt.Errorf("collecting references: %w", err)
	}

	usedCount := 0
	for _, a := range imageAssets {
		if refs[strings.ToLower(a.Basename)] {
			a.Used = true
			usedCount++
		}
	}

	return &ScanResult{
		ProjectPath:  projectAbsPath,
		AssetsDirs:   assetAbsPaths,
		TotalAssets:  len(imageAssets),
		UsedAssets:   usedCount,
		UnusedAssets: (len(imageAssets) - usedCount),
		Assets:       imageAssets,
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
		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("reading file info for %s: %w", absPath, err)
		}
		relPath, err := filepath.Rel(projectAbsPath, absPath)
		if err != nil {
			return fmt.Errorf("computing relative path for %s: %w", absPath, err)
		}

		imageAssets = append(imageAssets, &ImageAsset{
			AbsPath:   absPath,
			RelPath:   filepath.ToSlash(relPath),
			Basename:  strings.ToLower(filepath.Base(path)),
			SizeBytes: info.Size(),
		})

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("collecting assets: %w", err)
	}
	return imageAssets, nil
}

func (s *Scanner) collectReferences(projectAbsPath string, assetAbsPaths []string) (map[string]bool, error) {
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
			for _, assetAbsPath := range assetAbsPaths {
				if absDir == assetAbsPath {
					return fs.SkipDir
				}
			}
			return nil
		}

		if !IsSourceFile(path) {
			return nil
		}

		absPath := filepath.Join(projectAbsPath, path)
		if s.opts.Verbose {
			fmt.Printf("vole reading: %s\n", path)
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
