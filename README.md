# vole-clean

**Find and remove unused image assets in React/TypeScript projects.**

vole-clean scans your project, identifies image files inside your asset directories
that are never referenced in source code, and lets you safely delete them.

## Features

- **Scans** JS, JSX, TS, TSX, CSS, SCSS, Less, HTML, JSON, and Markdown files for image references
- **Detects** `import`, `require()`, `<img src>`, CSS `url()`, and bare quoted strings
- **Ignores** `node_modules`, `dist`, `build`, `.git`, `coverage`, `.next`, `out` by default
- **Multiple** asset directories — scan `public`, `assets`, or any set of directories
- **Configurable** via YAML file (`vole.yml` / `vole.yaml`) or CLI flags
- **Interactive** delete prompt with a `--yes` flag to skip it
- **Verbose** mode — watch every file being read and every deletion
- **Fast** — written in Go, walks your filesystem concurrently with zero dependencies for core logic
- **Cross-platform** — macOS, Linux, Windows binaries distributed via GitHub Releases and npm

## Installation

### Go install

```bash
go install github.com/devrapture/vole@latest
```

### npm (global)

```bash
npm install -g vole-clean
```

### Binary download

Download the latest binary for your platform from the
[GitHub Releases](https://github.com/devrapture/vole/releases) page.

### Build from source

```bash
git clone https://github.com/devrapture/vole.git
cd vole
make build
./bin/vole-clean
```

## Usage

```bash
# Scan the current project (default: assets in ./public)
vole-clean

# Scan a specific project and asset directory
vole-clean --project /path/to/project --assets src/assets

# Scan multiple asset directories
vole-clean --assets public --assets assets

# Skip the delete prompt
vole-clean --yes

# Verbose — log every file being read
vole-clean --verbose

# Ignore extra directories
vole-clean --ignore custom_cache --ignore .cache
```

## CLI Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--project` | `string` | `.` | Path to the React project root |
| `--assets` | `string` | `public` | Image assets sub-directory (repeatable) |
| `--ignore` | `stringSlice` | | Extra directory names to skip (repeatable) |
| `--verbose` | `bool` | `false` | Log every file being read and every deletion |
| `--yes` | `bool` | `false` | Delete without confirmation prompt |
## Configuration

vole-clean looks for a `vole.yml` (preferred) or `vole.yaml` in the project root.
CLI flags override values from the config file.

```yaml
# vole.yml — place in your project root

# Directories to scan for image assets.
# Overrides the --assets flag when set.
assets:
  - public
  - assets
  - src/images

# Extra directory names to skip during scanning.
# These are appended to the built-in ignore list
# (node_modules, dist, build, .git, coverage, .next, out).
ignore:
  - cypress
  - .cache
  - storybook-static
```

### Field reference

| Field | Type | CLI equivalent | Description |
|-------|------|----------------|-------------|
| `assets` | `[]string` | `--assets` | Asset directories to scan for images |
| `ignore` | `[]string` | `--ignore` | Extra directory names to skip |

## How it works

### 1. Scan

1. **Collect assets** — walks each configured asset directory and collects every
   image file (PNG, JPG, JPEG, GIF, SVG, WebP, ICO, AVIF).
2. **Collect references** — walks the project (skipping asset dirs and ignored
   dirs) and reads every source file, extracting image references using the
   patterns below.

### 2. Report

Displays total assets, used count, unused count, and lists every unused file.

### 3. Clean

Deletes the unused files and reports the space saved.

### Supported reference patterns

vole-clean detects image references in all of these forms:

```ts
// ES module imports
import logo from "./assets/logo.png"
import "./assets/banner.svg"

// CommonJS require
const img = require("./assets/photo.jpg")

// JSX img tags
<img src="/images/avatar.svg" />

// CSS url() — with and without quotes
background: url("../assets/bg.webp")
background: url(../icon.ico)

// Any quoted string — catches JSON configs, styled-components, etc.
"icon.gif"
```

HTTP and HTTPS URLs are automatically excluded:

```ts
// These will NOT be treated as local asset references
const url = "https://cdn.example.com/image.png"
```

## Development

```bash
make build        # Build to ./bin/vole-clean
make run          # Run directly
make test         # Run all tests
make test-verbose # Run all tests with verbose output
```

### Project structure

```
├── cmd/             # CLI entry point (cobra command)
├── internal/
│   ├── cleaner/     # File deletion logic
│   ├── config/      # vole.yml / vole.yaml parsing
│   └── scanner/     # File walking, reference extraction, result types
├── npm/             # npm distribution wrapper script
├── scripts/         # Build helper scripts
└── .github/         # CI / release workflows
```

## npm distribution

vole-clean is also published as an npm package. The package includes
platform-specific binaries for macOS (arm64/amd64), Linux (arm64/amd64), and
Windows (amd64). When installed globally via npm, the `vole-clean` command
automatically selects the correct binary for your platform.

```bash
npm install -g vole-clean
```

### npm packaging

The package.json `files` array ships only the `npm/` wrapper, `dist/` binaries,
the README, and a LICENSE. The source code and Go toolchain are **not** included
in the npm tarball — it's purely a binary distribution.

## License

MIT
