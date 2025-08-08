#!/bin/bash

# Check if MJML is already installed
if command -v mjml >/dev/null 2>&1; then
    echo "MJML CLI is already installed."
    mjml --version
else
    echo "MJML CLI not found. Installing..."
    
    # Check if npm is available
    if command -v npm >/dev/null 2>&1; then
        echo "Installing MJML CLI globally..."
        npm install -g mjml
        echo "MJML CLI installed successfully!"
        mjml --version
    else
        echo "Error: npm is not installed. Please install Node.js and npm first."
        echo "Visit: https://nodejs.org/"
        exit 1
    fi
fi
