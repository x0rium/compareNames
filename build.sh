#!/bin/bash

# Check if Go is installed
if ! command -v go &> /dev/null
then
    echo "ERROR: Go is not installed or not in PATH"
    exit 1
fi

# Clean previous build
rm -f compareNames

# Build the project
echo "Building compareNames..."
go build -o compareNames main.go

# Check build status
if [ $? -eq 0 ]; then
    echo "Build successful. Binary created: compareNames"
else
    echo "Build failed"
    exit 1
fi
