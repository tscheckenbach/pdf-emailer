#!/bin/bash

# Script to use Go 1.24.5 for development
# This script temporarily sets the PATH to use the Homebrew Go version

echo "Switching to Go 1.24.5..."
export PATH="/opt/homebrew/bin:$PATH"

# Verify the Go version
go version

echo ""
echo "Go 1.24.5 is now active in this shell session."
echo "To make this permanent, add the following to your shell profile:"
echo "export PATH=\"/opt/homebrew/bin:\$PATH\""
echo ""
echo "You can now run: go test ./..." 