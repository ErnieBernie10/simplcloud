#!/bin/bash

# This script copies assets from the source directory to the destination directory.
mkdir -p ./build/static
mkdir -p ./build/templates

npx @tailwindcss/cli -i ./src/internal/web/input.css -o ./build/static/output.css

cp -r ./src/internal/web/static/* ./build/static/
cp -r ./src/internal/web/templates/* ./build/templates/