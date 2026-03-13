#!/bin/bash

set -euo pipefail

# Default values
MODE=""
APPLY=false

# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --major)
      MODE="major"
      shift
      ;;
    --minor)
      MODE="minor"
      shift
      ;;
    --patch)
      MODE="patch"
      shift
      ;;
    --apply)
      APPLY=true
      shift
      ;;
    *)
      echo "Error: Unknown option '$1'"
      echo ""
      echo "Usage: $0 (--major | --minor | --patch) [--apply]"
      exit 1
      ;;
  esac
done

# Validate that a mode was selected
if [[ -z "$MODE" ]]; then
  echo "Error: Must specify one of --major, --minor, or --patch"
  echo ""
  echo "Usage: $0 (--major | --minor | --patch) [--apply]"
  exit 1
fi

# Get the latest tag
LATEST_TAG=$(git tag -l | sort -V | tail -n 1)

if [[ -z "$LATEST_TAG" ]]; then
  echo "Error: No existing tags found"
  exit 1
fi

# Parse the current version (remove 'v' prefix if present)
CURRENT_VERSION="${LATEST_TAG#v}"
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

# Calculate new version based on mode
case $MODE in
  major)
    MAJOR=$((MAJOR + 1))
    MINOR=0
    PATCH=0
    ;;
  minor)
    MINOR=$((MINOR + 1))
    PATCH=0
    ;;
  patch)
    PATCH=$((PATCH + 1))
    ;;
esac

NEW_VERSION="$MAJOR.$MINOR.$PATCH"
NEW_TAG="v$NEW_VERSION"

echo "Current version: $LATEST_TAG"
echo "New version:    $NEW_TAG"

if [[ "$APPLY" == true ]]; then
  echo ""
  echo "Pushing tag $NEW_TAG..."
  git tag "$NEW_TAG"
  git push origin "$NEW_TAG"
  echo "✓ Release $NEW_TAG pushed successfully"
else
  echo ""
  echo "Run with --apply to push this release"
fi
