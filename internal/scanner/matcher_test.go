package scanner

import (
	"testing"
)

func TestIsImageFile(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"logo.png", true},
		{"photo.jpg", true},
		{"photo.jpeg", true},
		{"animation.gif", true},
		{"icon.svg", true},
		{"image.webp", true},
		{"favicon.ico", true},
		{"image.avif", true},
		{"logo.PNG", true},
		{"logo.JPG", true},
		{"document.txt", false},
		{"main.go", false},
		{"noext", false},
		{"", false},
		{"data.json", false},
	}
	for _, tt := range tests {
		got := IsImageFile(tt.path)
		if got != tt.want {
			t.Errorf("IsImageFile(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestIsSourceFile(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"App.js", true},
		{"App.jsx", true},
		{"App.ts", true},
		{"App.tsx", true},
		{"styles.css", true},
		{"styles.scss", true},
		{"styles.less", true},
		{"index.html", true},
		{"data.json", true},
		{"readme.md", true},
		{"logo.png", false},
		{"main.go", false},
		{"Makefile", false},
	}
	for _, tt := range tests {
		got := IsSourceFile(tt.path)
		if got != tt.want {
			t.Errorf("IsSourceFile(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestIsIgnoredDirs(t *testing.T) {
	tests := []struct {
		name  string
		extra []string
		want  bool
	}{
		{"node_modules", nil, true},
		{"dist", nil, true},
		{"build", nil, true},
		{".git", nil, true},
		{"coverage", nil, true},
		{".next", nil, true},
		{"out", nil, true},
		{"src", nil, false},
		{"public", nil, false},
		{"custom_cache", []string{"custom_cache"}, true},
		{"custom_cache", []string{}, false},
		{"node_modules", []string{"custom_cache"}, true},
		{"", nil, false},
	}
	for _, tt := range tests {
		got := IsIgnoredDirs(tt.name, tt.extra)
		if got != tt.want {
			t.Errorf("IsIgnoredDirs(%q, %v) = %v, want %v", tt.name, tt.extra, got, tt.want)
		}
	}
}

func TestExtractReferences(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    map[string]bool
	}{
		{
			name:    "empty content",
			content: "",
			want:    map[string]bool{},
		},
		{
			name:    "no references",
			content: "const x = 1;",
			want:    map[string]bool{},
		},
		{
			name: "import statement",
			content: `import logo from "./assets/logo.png"`,
			want:    map[string]bool{"logo.png": true},
		},
		{
			name: "require call",
			content: `const img = require("./assets/photo.jpg")`,
			want:    map[string]bool{"photo.jpg": true},
		},
		{
			name: "img src attribute",
			content: `<img src="/images/avatar.svg" />`,
			want:    map[string]bool{"avatar.svg": true},
		},
		{
			name: "url() in CSS",
			content: `background: url("../assets/bg.webp")`,
			want:    map[string]bool{"bg.webp": true},
		},
		{
			name: "bare quoted string",
			content: `"icon.gif"`,
			want:    map[string]bool{"icon.gif": true},
		},
		{
			name: "mixed case extension",
			content: `import logo from "./assets/logo.PNG"`,
			want:    map[string]bool{"logo.png": true},
		},
		{
			name: "http URLs skipped",
			content: `const url = "https://example.com/image.png"`,
			want:    map[string]bool{},
		},
		{
			name: "multiple references",
			content: `import a from "./a.png"
const b = require("./b.jpg")
<img src="/c.svg" />`,
			want:    map[string]bool{"a.png": true, "b.jpg": true, "c.svg": true},
		},
		{
			name: "duplicate references deduplicated",
			content: `import logo from "./logo.png"
<img src="/logo.png" />`,
			want:    map[string]bool{"logo.png": true},
		},
		{
			name: "require with single quotes",
			content: `const img = require('./photo.jpg')`,
			want:    map[string]bool{"photo.jpg": true},
		},
		{
			name: "url with no quotes",
			content: `background: url(../icon.ico)`,
			want:    map[string]bool{"icon.ico": true},
		},
		{
			name: "path with directory prefix",
			content: `import logo from "./assets/subdir/logo.png"`,
			want:    map[string]bool{"logo.png": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractReferences(tt.content)
			if len(got) != len(tt.want) {
				t.Errorf("ExtractReferences() returned %d refs, want %d: %v", len(got), len(tt.want), got)
				return
			}
			for k := range tt.want {
				if !got[k] {
					t.Errorf("ExtractReferences() missing key %q", k)
				}
			}
		})
	}
}

func TestExtractReferencesExcludesHTTP(t *testing.T) {
	content := `const a = "http://example.com/img.png"
const b = "https://cdn.example.com/photo.jpg"
const c = "/local/image.svg"`
	refs := ExtractReferences(content)
	if refs["img.png"] {
		t.Error("http:// URL should not be extracted")
	}
	if refs["photo.jpg"] {
		t.Error("https:// URL should not be extracted")
	}
	if !refs["image.svg"] {
		t.Error("local path should be extracted")
	}
}
