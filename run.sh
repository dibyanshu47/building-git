#!/bin/bash

go build -o /tmp/my-build-git-go ./cmd/mygit

# Check if the build was successful
if [ $? -ne 0 ]; then
  echo "Go build failed. Exiting..."
  exit 1
fi

exec /tmp/my-build-git-go "$@"
