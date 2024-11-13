#!/bin/bash

# Variables
REPO="github.com/oliverst8/gottis/cmd/gottis" # replace with your actual module path (repository URL)
INSTALL_DIR="/usr/local/bin"

# Ensure the Go binary is available
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go and try again."
    exit 1
fi

# Install the Go program
echo "Installing the Go program from $REPO..."
GOBIN="$INSTALL_DIR" go install "$REPO@latest" || {
    echo "Failed to install the Go program. Check the module path and try again."
    exit 1
}

echo "Installation complete! You can now run gottis globally."
