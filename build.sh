#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

clear

detect_cpu_count() {
    case "$(uname -s)" in
        Darwin)        sysctl -n hw.logicalcpu ;;
        Linux)         nproc ;;
        MINGW*|MSYS*)  nproc ;;
        *)             echo 4 ;;
    esac
}

MODE="${1:-}"
BUILD_TYPE="Release"

if [ "$MODE" = "debug" ]; then
    BUILD_TYPE="Debug"
    MODE=""
fi

BUILD_DIR="Builds/Ninja"

if [ "$MODE" = "clean" ]; then
    echo "Cleaning..."
    rm -rf "$BUILD_DIR"
fi

if [ "$MODE" = "clean" ] || [ ! -f "$BUILD_DIR/build.ninja" ]; then
    echo "Configuring ($BUILD_TYPE)..."
    cmake -S . -B "$BUILD_DIR" -G Ninja -DCMAKE_BUILD_TYPE="$BUILD_TYPE"
fi

if [ "$MODE" != "configure" ]; then
    echo "Building ($BUILD_TYPE)..."
    cmake --build "$BUILD_DIR" -- -j"$(detect_cpu_count)"
    echo "Build succeeded."
    echo "Binary: $BUILD_DIR/cakec_artefacts/cakec"
fi
