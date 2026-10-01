#!/bin/sh
# anon installer for macOS and Linux.
# Usage: curl -fsSL __SERVER__/install.sh | sh
set -e

SERVER="__SERVER__"
BIN_DIR="$HOME/.local/bin"

case "$(uname -s)" in
  Darwin) os=darwin ;;
  Linux)  os=linux ;;
  *) echo "Unsupported OS: $(uname -s)"; exit 1 ;;
esac
case "$(uname -m)" in
  x86_64|amd64)  arch=amd64 ;;
  arm64|aarch64) arch=arm64 ;;
  *) echo "Unsupported CPU: $(uname -m)"; exit 1 ;;
esac

echo "Downloading anon ($os/$arch)..."
mkdir -p "$BIN_DIR"
curl -fsSL "$SERVER/dl/anon-$os-$arch" -o "$BIN_DIR/anon.tmp"
chmod +x "$BIN_DIR/anon.tmp"
mv "$BIN_DIR/anon.tmp" "$BIN_DIR/anon"

# Add ~/.local/bin to PATH if it is not there yet.
case ":$PATH:" in
  *":$BIN_DIR:"*) ;;
  *)
    line='export PATH="$HOME/.local/bin:$PATH"'
    for rc in "$HOME/.zshrc" "$HOME/.bashrc"; do
      if [ -f "$rc" ] || [ "$rc" = "$HOME/.zshrc" -a "$os" = darwin ]; then
        grep -qF "$line" "$rc" 2>/dev/null || printf '\n%s\n' "$line" >> "$rc"
      fi
    done
    ;;
esac

if [ -z "$ANON_KEY" ]; then
  printf "Team key: "
  read -r ANON_KEY < /dev/tty
fi
"$BIN_DIR/anon" config "$SERVER" "$ANON_KEY" > /dev/null

echo ""
echo "✓ anon installed!"
echo "  1. Telegram e bot ke /start pathan (message receive korar jonno)"
echo "  2. Notun terminal khule try korun: anon members"
