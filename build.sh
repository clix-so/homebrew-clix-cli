#!/bin/bash
set -e

echo "Building Clix CLI with Bun..."

# Install dependencies
bun install

# Build the CLI
bun build src/cli.tsx --target=bun --outdir=dist --minify

# Make the output executable
chmod +x dist/cli.js

# Add shebang if not present
if ! head -n 1 dist/cli.js | grep -q "^#!"; then
  echo "#!/usr/bin/env bun" | cat - dist/cli.js > dist/cli.tmp
  mv dist/cli.tmp dist/cli.js
  chmod +x dist/cli.js
fi

echo "✅ Build complete! Output: dist/cli.js"
