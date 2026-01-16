#!/bin/bash

echo "Starting Development Environment Setup..."

# 1. Check Go Version
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.22+ first."
    exit 1
fi

echo "Go is installed."

# 2. Install System Dependencies (Ubuntu/Debian)
# Checks if apt is available (for Linux)
if command -v apt &> /dev/null; then
    echo "Installing Protobuf Compiler..."
    sudo apt update && sudo apt install -y protobuf-compiler
else
    echo "Warning: Not on Ubuntu/Debian. Please ensure 'protoc' is installed manually."
fi

# 3. Install Go Tools (The Kratos Toolchain)
echo "Installing Go tools (Kratos, Wire, Protoc Plugins)..."

# Kratos CLI
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest

# Wire (Dependency Injection)
go install github.com/google/wire/cmd/wire@latest

# Protobuf Plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest

# 4. Install Project Dependencies
echo "Downloading project dependencies (go mod)..."
go mod tidy

echo "Setup Complete! You can now run 'kratos run' or use 'wire' to generate code."
