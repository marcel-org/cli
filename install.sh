#!/bin/bash

set -e

SKILL_MARKER_START="<!-- marcel:start -->"
SKILL_MARKER_END="<!-- marcel:end -->"

info()  { printf '\033[1;36m%s\033[0m\n' "$*"; }

inject_block() {
  local file="$1"
  local content="$2"
  local block
  block="$(printf '%s\n%s\n%s' "$SKILL_MARKER_START" "$content" "$SKILL_MARKER_END")"

  if [ ! -f "$file" ]; then
    printf '%s\n' "$block" > "$file"
    return
  fi

  if grep -qF "$SKILL_MARKER_START" "$file"; then
    local tmp
    tmp="$(mktemp)"
    awk -v start="$SKILL_MARKER_START" -v end="$SKILL_MARKER_END" '
      $0 == start { skip=1; next }
      $0 == end   { skip=0; next }
      !skip       { print }
    ' "$file" > "$tmp"
    mv "$tmp" "$file"
    printf '\n%s\n' "$block" >> "$file"
  else
    printf '\n%s\n' "$block" >> "$file"
  fi
}

register_skill() {
  local repo_dir="$1"
  local skill_file="$repo_dir/integrations/SKILL.md"

  [ -f "$skill_file" ] || return 0

  local skill_content
  skill_content="$(cat "$skill_file")"

  if command -v claude &>/dev/null; then
    local skill_dir="$HOME/.claude/skills/marcel"
    mkdir -p "$skill_dir"
    cp "$skill_file" "$skill_dir/SKILL.md"
    inject_block "$HOME/.claude/CLAUDE.md" "$skill_content"
    info "  ✓ Claude Code skill registered"
  fi

  if command -v codex &>/dev/null; then
    mkdir -p "$HOME/.codex"
    inject_block "$HOME/.codex/AGENTS.md" "$skill_content"
    info "  ✓ Codex skill registered"
  fi
}

REPO="marcel-org/cli"
INSTALL_DIR="$HOME/.local/bin"

echo "Installing Marcel CLI..."

mkdir -p "$INSTALL_DIR"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
    linux|darwin) ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

BINARY_NAME="marcel-${OS}-${ARCH}"

if command -v go >/dev/null 2>&1; then
    echo "Building from source..."
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"

    if command -v git >/dev/null 2>&1; then
        git clone "https://github.com/${REPO}.git" .
        cd src
        go build -o ../marcel .
        cd ..
    else
        echo "Error: git is required to build from source"
        exit 1
    fi
else
    echo "Go not found. Please install Go or use a pre-built binary."
    exit 1
fi

mv marcel "$INSTALL_DIR/marcel"
chmod +x "$INSTALL_DIR/marcel"

register_skill "$TEMP_DIR"

rm -rf "$TEMP_DIR" 2>/dev/null || true

if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo ""
    echo "Add $INSTALL_DIR to your PATH by adding this to your shell config:"
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
fi

echo ""
echo "Marcel CLI installed successfully!"
echo ""
echo "Configuration:"
echo "  1. Get your token from Marcel web app settings"
echo "  2. Set environment variable: export MARCEL_TOKEN=\"your_token_here\""
echo "  3. Run: marcel"
echo ""
