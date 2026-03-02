#!/bin/sh

set -e

# Determine OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalize architecture names
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
esac

# Create tools directory
mkdir -p "${HOME}/tools"

# Download and extract the binary
RELEASE_URL="https://github.com/ronkitay/griffin/releases/latest/download/griffin-${OS}-${ARCH}"

echo "Downloading griffin for ${OS}/${ARCH}..."
curl -sL -o "${HOME}/tools/griffin" "$RELEASE_URL"
chmod +x "${HOME}/tools/griffin"

# Check what needs to be added to .zshrc
ZSHRC="${HOME}/.zshrc"
NEED_PATH=false
NEED_INTEGRATION=false

if [ ! -f "$ZSHRC" ] || ! grep -q '${HOME}/tools' "$ZSHRC"; then
  NEED_PATH=true
fi

if [ ! -f "$ZSHRC" ] || ! grep -q 'griffin shell-integration' "$ZSHRC"; then
  NEED_INTEGRATION=true
fi

# If any changes are needed, back up the file first
if [ "$NEED_PATH" = true ] || [ "$NEED_INTEGRATION" = true ]; then
  if [ -f "$ZSHRC" ]; then
    cp "$ZSHRC" "$ZSHRC.griffin-backup"
    echo "Backed up .zshrc to .zshrc.griffin-backup"
  fi
  
  if [ "$NEED_PATH" = true ]; then
    echo 'export PATH="${HOME}/tools:$PATH"' >> "$ZSHRC"
    echo "Added ${HOME}/tools to PATH in .zshrc"
  fi
  
  if [ "$NEED_INTEGRATION" = true ]; then
    echo 'source <(griffin shell-integration)' >> "$ZSHRC"
    echo "Added griffin shell integration to .zshrc"
  fi
fi

# Remove quarantine attribute (macOS)
xattr -d com.apple.quarantine "${HOME}/tools/griffin" 2>/dev/null || true

echo "Installation complete! Run 'source ~/.zshrc' or restart your terminal."
