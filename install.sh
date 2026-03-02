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
REPO="your-username/griffin"  # Update with actual repo
RELEASE_URL="https://github.com/${REPO}/releases/latest/download/griffin-${OS}-${ARCH}"

echo "Downloading griffin for ${OS}/${ARCH}..."
curl -sL -o "${HOME}/tools/griffin" "$RELEASE_URL"
chmod +x "${HOME}/tools/griffin"

# Add to PATH in .zshrc if not already present
ZSHRC="${HOME}/.zshrc"
if [ -f "$ZSHRC" ]; then
  if ! grep -q '${HOME}/tools' "$ZSHRC"; then
    echo 'export PATH="${HOME}/tools:$PATH"' >> "$ZSHRC"
    echo "Added ${HOME}/tools to PATH in .zshrc"
  fi
else
  echo 'export PATH="${HOME}/tools:$PATH"' > "$ZSHRC"
  echo "Created .zshrc with ${HOME}/tools in PATH"
fi

# Remove quarantine attribute (macOS)
xattr -d com.apple.quarantine "${HOME}/tools/griffin" 2>/dev/null || true

echo "Installation complete! Run 'source ~/.zshrc' or restart your terminal."
