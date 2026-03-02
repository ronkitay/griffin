#!/bin/sh

set -e

# Determine OS and architecture
OS=$(uname -s)
ARCH=$(uname -m)

# Normalize OS and architecture names to match release artifacts
case "$OS" in
  Darwin) OS="Darwin" ;;
  Linux) OS="Linux" ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
  x86_64) ARCH="x86_64" ;;
  amd64) ARCH="x86_64" ;;
  aarch64) ARCH="arm64" ;;
  arm64) ARCH="arm64" ;;
  i386|i686) ARCH="i386" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Create tools directory
mkdir -p "${HOME}/tools"

# Determine file extension based on OS
case "$OS" in
  Darwin|Linux) EXT="tar.gz" ;;
  Windows) EXT="zip" ;;
esac

# Construct the download URL
RELEASE_URL="https://github.com/ronkitay/griffin/releases/latest/download/griffin_${OS}_${ARCH}.${EXT}"

echo "Downloading griffin for ${OS}/${ARCH}..."
TEMP_FILE="${HOME}/tools/griffin.${EXT}"
curl -sL -o "$TEMP_FILE" "$RELEASE_URL" || { echo "Failed to download from $RELEASE_URL"; exit 1; }

# Extract the binary
case "$EXT" in
  tar.gz)
    tar -xzf "$TEMP_FILE" -C "${HOME}/tools" griffin || { echo "Failed to extract $TEMP_FILE"; exit 1; }
    ;;
  zip)
    unzip -o "$TEMP_FILE" -d "${HOME}/tools" griffin || { echo "Failed to extract $TEMP_FILE"; exit 1; }
    ;;
esac

rm "$TEMP_FILE"
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
