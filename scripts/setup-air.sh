#!/bin/bash

# Setup script for Air development environment

echo "🚀 Setting up Air for automatic reload..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

# Install Air
echo "📦 Installing Air..."
go install github.com/air-verse/air@latest

# Check if Air was installed successfully
if ! command -v air &> /dev/null; then
    echo "❌ Failed to install Air. Please check your Go installation."
    exit 1
fi

echo "✅ Air installed successfully!"

# Create tmp directory if it doesn't exist
mkdir -p tmp

echo "🎉 Setup complete!"
echo ""
echo "To start development with auto-reload:"
echo "  make dev"
echo ""
echo "Or run directly:"
echo "  air"
echo ""
echo "Air will automatically reload your application when you make changes to:"
echo "  - Go files (*.go)"
echo "  - Template files (*.html, *.tpl, *.tmpl)"
echo "  - Static files" 