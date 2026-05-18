#!/usr/bin/env bash
set -euo pipefail

BINARY="ollama-oneshot"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info()    { echo -e "${GREEN}[+]${NC} $*"; }
warn()    { echo -e "${YELLOW}[!]${NC} $*"; }
error()   { echo -e "${RED}[x]${NC} $*" >&2; exit 1; }

check_go() {
    if ! command -v go &>/dev/null; then
        error "Go is not installed. Install it from https://go.dev/dl/"
    fi
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    info "Go $GO_VERSION detected"
}

build() {
    info "Building $BINARY..."
    cd "$REPO_ROOT"
    VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0")
    go build -ldflags "-s -w -X main.version=$VERSION" -o "$BINARY" .
    info "Build complete"
}

install_binary() {
    mkdir -p "$INSTALL_DIR"

    if [ -f "$INSTALL_DIR/$BINARY" ]; then
        warn "Existing installation found at $INSTALL_DIR/$BINARY — overwriting"
    fi

    mv "$REPO_ROOT/$BINARY" "$INSTALL_DIR/$BINARY"
    chmod +x "$INSTALL_DIR/$BINARY"
    info "Installed to $INSTALL_DIR/$BINARY"
}

check_path() {
    if ! echo "$PATH" | tr ':' '\n' | grep -qx "$INSTALL_DIR"; then
        warn "$INSTALL_DIR is not in your PATH"
        echo ""
        echo "  Add this to your shell config (~/.bashrc, ~/.zshrc, etc.):"
        echo ""
        echo "    export PATH=\"\$HOME/.local/bin:\$PATH\""
        echo ""
    else
        info "$INSTALL_DIR is in PATH"
    fi
}

setup_env() {
    if [ ! -f "$REPO_ROOT/.env" ]; then
        warn ".env not found — creating default"
        cat > "$REPO_ROOT/.env" <<EOF
OLLAMA_HOST=127.0.0.1:11434

DEFAULT_MODEL=kimi-k2.6:cloud
DEFAULT_TOOL=claude

PROMPT_ENHANCEMENT=true
PROMPT_ENHANCEMENT_MODEL=deepseek-v4-flash
EOF
        info ".env created at $REPO_ROOT/.env"
    fi
}

main() {
    echo ""
    echo "  ollama-oneshot installer"
    echo "  ─────────────────────────"
    echo ""

    check_go
    setup_env
    build
    install_binary
    check_path

    echo ""
    info "Done. Run: $BINARY --help"
    echo ""
}

main "$@"
