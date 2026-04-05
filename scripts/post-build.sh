#!/bin/bash
# Called by goreleaser after each build. Signs and notarizes macOS binaries.
# Env vars set by goreleaser: HOOK_TARGET (e.g. darwin_arm64), HOOK_PATH (binary path)

case "$HOOK_TARGET" in
  darwin_*)
    echo "Signing $HOOK_PATH..."
    codesign --force --options runtime \
      --entitlements ./entitlements.plist \
      --sign "Developer ID Application: Bayu Ardianto (9BDSN9TDX3)" \
      "$HOOK_PATH"
    codesign --verify --verbose "$HOOK_PATH"

    echo "Notarizing $HOOK_PATH..."
    ZIP=$(mktemp /tmp/cake-notarize-XXXXXX.zip)
    ditto -c -k --keepParent "$HOOK_PATH" "$ZIP"
    xcrun notarytool submit "$ZIP" \
      --keychain-profile notary \
      --wait
    rm "$ZIP"
    ;;
esac
