package scanner

import (
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var ImageExtensions = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".webp": true, ".svg": true, ".gif": true, ".ico": true, ".avif": true,
}

var IgnoredDirs = map[string]bool{
	"node_modules": true, "dist": true, "build": true, ".git": true, "coverage": true, ".next": true, "out": true,
}

var SourceExtensions = map[string]bool{
	".js": true, ".jsx": true, ".ts": true, ".tsx": true, ".css": true, ".scss": true, ".less": true, ".html": true, ".json": true, ".md": true,
}

// referencePatterns catches every common way a React project can reference
// an image. Capture group 1 always holds the raw path or filename.
var referencePatterns = []*regexp.Regexp{
	// import logo from "./assets/logo.png"
	// import "./assets/banner.svg"
	regexp.MustCompile(
		`(?i)import\s+(?:[^"']*\s+from\s+)?["']([^"']+\.(?:png|jpg|jpeg|webp|svg|gif|ico|avif))["']`,
	),
	// require("./assets/photo.jpg")
	regexp.MustCompile(
		`(?i)require\s*\(\s*["']([^"']+\.(?:png|jpg|jpeg|webp|svg|gif|ico|avif))["']\s*\)`,
	),
	// <img src="/images/avatar.png" />
	regexp.MustCompile(
		`(?i)src\s*=\s*["']([^"']+\.(?:png|jpg|jpeg|webp|svg|gif|ico|avif))["']`,
	),
	// url("../assets/bg.webp")  url('../bg.png')  url(../icon.svg)
	regexp.MustCompile(
		`(?i)url\s*\(\s*["']?([^"')]+\.(?:png|jpg|jpeg|webp|svg|gif|ico|avif))["']?\s*\)`,
	),
	// Any bare quoted string ending in an image extension.
	// Catches JSON configs, styled-components, and public-folder refs.
	regexp.MustCompile(
		`(?i)["']([^"']*\.(?:png|jpg|jpeg|webp|svg|gif|ico|avif))["']`,
	),
}

func IsImageFile(path string) bool {
	return ImageExtensions[strings.ToLower(filepath.Ext(path))]
}

func IsSourceFile(path string) bool {
	return SourceExtensions[strings.ToLower(filepath.Ext(path))]
}

func IsIgnoredDirs(name string, extra []string) bool {
	if IgnoredDirs[name] {
		return true
	}
	return slices.Contains(extra, name)
}

func ExtractReferences(content string) map[string]bool {
	refs := make(map[string]bool)
	for _, re := range referencePatterns {
		for _, m := range re.FindAllStringSubmatch(content, -1) {
			if len(m) < 2 {
				continue
			}
			raw := m[1]
			// Skip CDN / HTTP URLs — vole only tracks local assets.
			if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
				continue
			}
			base := strings.ToLower(filepath.Base(raw))
			if base != "" && base != "." {
				refs[base] = true
			}
		}
	}
	return refs
}
