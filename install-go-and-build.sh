#!/bin/bash
# Check if ~/go directory exists
if [ ! -d "$HOME/go" ]; then
  echo "Go installation not found. Installing Go..."
  curl -L https://go.dev/dl/go1.24.1.linux-amd64.tar.gz -O
  tar -xzf go1.24.1.linux-amd64.tar.gz -C $HOME
  rm go1.24.1.linux-amd64.tar.gz
  export PATH=$PATH:$HOME/go/bin
  echo "Go $(go version) installed successfully"
else
  echo "Go installation found, skipping download"
  export PATH=$PATH:$HOME/go/bin
fi

# Run build.sh
echo "Running build script..."
./build.sh
