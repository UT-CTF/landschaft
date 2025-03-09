#!/bin/bash
set -e

VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS="-X github.com/UT-CTF/landschaft/cmd.Version=${VERSION} -X github.com/UT-CTF/landschaft/cmd.BuildTime=${BUILD_TIME}"

echo "Building Landschaft ${VERSION} (${BUILD_TIME})"
echo "----------------------------------------"

# Create output directory
mkdir -p build

# Build for Linux (amd64)
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o build/landschaft-linux-amd64

# Build for Windows (amd64)
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o build/landschaft-windows-amd64.exe

echo "----------------------------------------"
echo "Build complete! Binaries available in the build/ directory:"
ls -lh build/

echo "
Linux:   build/landschaft-linux-amd64
Windows: build/landschaft-windows-amd64.exe"
