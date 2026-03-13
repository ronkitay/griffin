#!/bin/bash

set -euo pipefail

# Default values
MODE=""
APPLY=false

# Convert long options to short options
for arg in "$@"; do
  shift
  case "$arg" in
    '--major')  set -- "$@" '-M' ;;
    '--minor')  set -- "$@" '-m' ;;
    '--patch')  set -- "$@" '-p' ;;
    '--apply')  set -- "$@" '-a' ;;
    '--help')   set -- "$@" '-h' ;;
    *)          set -- "$@" "$arg" ;;
  esac
done

# Parse options using getopts
while getopts ":Mmpha" option; do
  case "${option}" in
    M)
      if [[ -n "$MODE" ]]; then
        echo "Error: Cannot specify multiple modes (-M, -m, -p / --major, --minor, --patch)"
        exit 1
      fi
      MODE="major"
      ;;
    m)
      if [[ -n "$MODE" ]]; then
        echo "Error: Cannot specify multiple modes (-M, -m, -p / --major, --minor, --patch)"
        exit 1
      fi
      MODE="minor"
      ;;
    p)
      if [[ -n "$MODE" ]]; then
        echo "Error: Cannot specify multiple modes (-M, -m, -p / --major, --minor, --patch)"
        exit 1
      fi
      MODE="patch"
      ;;
    a)
      APPLY=true
      ;;
    h)
      echo "Usage: $0 (-M | -m | -p | --major | --minor | --patch) [-a | --apply]"
      echo ""
      echo "Options:"
      echo "  -M, --major   Bump major version"
      echo "  -m, --minor   Bump minor version"
      echo "  -p, --patch   Bump patch version"
      echo "  -a, --apply   Actually push the tag (default: dry-run)"
      echo "  -h, --help    Show this help message"
      exit 0
      ;;
    *)
      echo "Error: Invalid option '$OPTARG'"
      echo "Usage: $0 (-M | -m | -p | --major | --minor | --patch) [-a | --apply]"
      exit 1
      ;;
  esac
done
shift $((OPTIND - 1))

# Validate that a mode was selected
if [[ -z "$MODE" ]]; then
  echo "Error: Must specify one of -M, -m, -p (or --major, --minor, --patch)"
  echo "Usage: $0 (-M | -m | -p | --major | --minor | --patch) [-a | --apply]"
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
