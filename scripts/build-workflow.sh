#!/usr/bin/env bash
# Build the .alfredworkflow package.
#
# Steps:
#   1. Build the two binaries the packaged workflow invokes
#      (cmd/note-md-template-alfred, cmd/note-md-template-paste-alfred) as
#      universal (amd64+arm64) binaries via lipo, so the bundle runs
#      natively on both Intel and Apple Silicon.
#   2. If CODESIGN_IDENTITY is set, codesign each binary (and notarize it if
#      NOTARY_KEY_ID is also set). Unset by default, so an ordinary local
#      `make build-workflow` stays unsigned; .github/workflows/release.yml
#      sets these for tagged releases.
#   3. Copy workflow/ (info.plist, icon.png) into the build dir.
#   4. Zip into dist/<name>-<version>.alfredworkflow.
set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
WORKFLOW_DIR="${REPO_ROOT}/workflow"
DIST_DIR="${REPO_ROOT}/dist"
BUILD_DIR="${REPO_ROOT}/.build"

BINARIES=(note-md-template-alfred note-md-template-paste-alfred)

echo "→ Preparing build directory"
rm -rf "$BUILD_DIR"
cp -r "$WORKFLOW_DIR/" "$BUILD_DIR/"

echo "→ Building universal binaries (amd64 + arm64)"
for bin in "${BINARIES[@]}"; do
  GOOS=darwin GOARCH=amd64 go build -o "${BUILD_DIR}/${bin}-amd64" "./cmd/${bin}"
  GOOS=darwin GOARCH=arm64 go build -o "${BUILD_DIR}/${bin}-arm64" "./cmd/${bin}"
  lipo -create -output "${BUILD_DIR}/${bin}" "${BUILD_DIR}/${bin}-amd64" "${BUILD_DIR}/${bin}-arm64"
  rm "${BUILD_DIR}/${bin}-amd64" "${BUILD_DIR}/${bin}-arm64"
  chmod +x "${BUILD_DIR}/${bin}"
  lipo -info "${BUILD_DIR}/${bin}"

  if [ -n "${CODESIGN_IDENTITY:-}" ]; then
    echo "→ Signing ${bin} (${CODESIGN_IDENTITY})"
    codesign --force --options runtime --timestamp --sign "$CODESIGN_IDENTITY" \
      "${BUILD_DIR}/${bin}"
    codesign --verify --strict --verbose=2 "${BUILD_DIR}/${bin}"

    if [ -n "${NOTARY_KEY_ID:-}" ]; then
      "${REPO_ROOT}/scripts/notarize-binary.sh" "${BUILD_DIR}/${bin}"
    fi
  else
    echo "→ Skipping signing of ${bin} (CODESIGN_IDENTITY not set) — unsigned local/dev build"
  fi
done

VERSION=$(/usr/libexec/PlistBuddy -c "Print :version" "${BUILD_DIR}/info.plist")
WORKFLOW_NAME=$(/usr/libexec/PlistBuddy -c "Print :name" "${BUILD_DIR}/info.plist" | tr '[:upper:] ' '[:lower:]-')

mkdir -p "$DIST_DIR"
OUTPUT="${DIST_DIR}/${WORKFLOW_NAME}-${VERSION}.alfredworkflow"
rm -f "$OUTPUT" # ensure a clean zip (zip -r updates rather than replaces)

echo "→ Packaging: ${OUTPUT}"
(cd "$BUILD_DIR" && zip -r "$OUTPUT" . -x "*.DS_Store" --quiet)

echo "✓ Build complete: ${OUTPUT}"
