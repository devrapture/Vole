#!/usr/bin/env bash
set -euo pipefail

REPO="devrapture/Vole"
BIN_NAME="vole-clean"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
VERSION="${VERSION:-latest}"

detect_os() {
  case "$(uname -s)" in
    Darwin) echo "darwin" ;;
    Linux) echo "linux" ;;
    *) echo "unsupported" ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64) echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    *) echo "unsupported" ;;
  esac
}

OS="$(detect_os)"
ARCH="$(detect_arch)"

if [ "$OS" = "unsupported" ] || [ "$ARCH" = "unsupported" ]; then
  echo "Unsupported platform: $(uname -s) $(uname -m)"
  exit 1
fi

if [ "$VERSION" = "latest" ]; then
  URL="https://github.com/${REPO}/releases/latest/download/${BIN_NAME}-${OS}-${ARCH}"
else
  URL="https://github.com/${REPO}/releases/download/${VERSION}/${BIN_NAME}-${OS}-${ARCH}"
fi

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Downloading ${BIN_NAME} for ${OS}/${ARCH}..."
curl -fsSL "$URL" -o "$TMP_DIR/$BIN_NAME"

chmod +x "$TMP_DIR/$BIN_NAME"

if [ ! -w "$INSTALL_DIR" ]; then
  echo "Installing to $INSTALL_DIR requires sudo..."
  sudo mkdir -p "$INSTALL_DIR"
  sudo mv "$TMP_DIR/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
else
  mkdir -p "$INSTALL_DIR"
  mv "$TMP_DIR/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
fi

echo "Installed:"
"$INSTALL_DIR/$BIN_NAME" --help | head -n 2
