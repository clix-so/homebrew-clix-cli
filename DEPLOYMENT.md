# Deployment Guide

This guide explains how to publish the Clix CLI to npm and update the Homebrew formula.

## Deployment Methods

There are two ways to deploy:

1. **Automated (Recommended)**: Using GitHub Actions - see [Automated Deployment](#automated-deployment-recommended)
2. **Manual**: Step-by-step manual process - see [Manual Deployment](#manual-deployment)

---

## Automated Deployment (Recommended)

The easiest way to release is using GitHub Actions, which automatically:
- Builds the project
- Publishes to npm
- Calculates SHA256
- Updates Homebrew formula
- Creates GitHub release

### Prerequisites

- npm account with publish access to `@clix-so/clix-cli`
- GitHub access to `clix-so/homebrew-clix-cli` repository
- `NPM_TOKEN` configured in GitHub Secrets (see [.github/SETUP.md](.github/SETUP.md))

### Steps

1. **Update version in package.json**
   ```bash
   # Update version in package.json
   npm version patch  # or minor, or major
   ```

2. **Commit and push to main**
   ```bash
   # Commit the version change
   git add package.json package-lock.json
   git commit -m "chore: bump version to vX.X.X"

   # Push to main branch
   git push origin main
   ```

3. **Wait for GitHub Actions to complete**
   - GitHub Actions automatically detects the version change in package.json
   - Go to GitHub → Actions → [Release workflow](../../actions)
   - Watch the workflow complete all steps
   - The workflow will automatically:
     - ✅ Detect version change
     - ✅ Create git tag (e.g., v1.0.1)
     - ✅ Build the project
     - ✅ Publish to npm
     - ✅ Calculate SHA256 for the tarball
     - ✅ Update `clix.rb` with new version and hash
     - ✅ Commit and push the updated formula
     - ✅ Create a GitHub release

4. **Verify the release**
   ```bash
   # Check npm
   npm view @clix-so/clix-cli version

   # Check Homebrew (after a few minutes for cache)
   brew update
   brew install clix-so/clix-cli/clix
   clix --version
   ```

That's it! The entire release process is automated when you push a version change to main.

### How It Works

The workflow:
1. **Triggers on main branch push** when `package.json` is modified
2. **Compares versions** between the current and previous commit
3. **Skips release** if version hasn't changed (no unnecessary releases)
4. **Checks for existing tags** to prevent duplicate releases
5. **Creates git tag** automatically (e.g., v1.0.1)
6. **Publishes and updates** everything automatically

### Manual Trigger (Optional)

You can also manually trigger the workflow:
1. Go to GitHub → Actions → Release workflow
2. Click "Run workflow"
3. Select the branch
4. Click "Run workflow"

Note: Manual triggers still require a version change to proceed with the release.

For setup details, see [.github/SETUP.md](.github/SETUP.md).

---

## Manual Deployment

If you need to deploy manually without GitHub Actions:

### Prerequisites

- npm account with publish access to `@clix-so/clix-cli`
- GitHub access to `clix-so/homebrew-clix-cli` repository

## Steps

### 1. Build and Test Locally

```bash
# Build the project
npm run build

# Test locally with npm link
npm link
clix --help
clix --version

# Unlink when done testing
npm unlink -g @clix-so/clix-cli
```

### 2. Update Version

Update the version in [package.json](package.json):

```bash
npm version patch  # for bug fixes
npm version minor  # for new features
npm version major  # for breaking changes
```

### 3. Publish to npm

```bash
# Login to npm (if not already logged in)
npm login

# Publish the package
npm publish --access public
```

### 4. Calculate SHA256 for Homebrew Formula

After publishing, download the tarball and calculate its SHA256:

```bash
# Download the published tarball
npm pack @clix-so/clix-cli@<version>

# Calculate SHA256
sha256sum clix-so-clix-cli-<version>.tgz
# or on macOS:
shasum -a 256 clix-so-clix-cli-<version>.tgz
```

### 5. Update Homebrew Formula

Update [clix.rb](clix.rb):

1. Update the `url` with the new version
2. Update the `sha256` with the calculated hash from step 4

```ruby
class Clix < Formula
  desc "A CLI tool for integrating and managing the Clix SDK in your mobile projects"
  homepage "https://github.com/clix-so/homebrew-clix-cli"
  url "https://registry.npmjs.org/@clix-so/clix-cli/-/clix-cli-1.0.0.tgz"
  sha256 "YOUR_CALCULATED_SHA256_HERE"
  license "MIT"

  depends_on "node@18"

  def install
    system "npm", "install", *std_npm_args(prefix: libexec)
    bin.install_symlink libexec/"bin/clix"
  end

  test do
    assert_match "A CLI tool for integrating and managing the Clix SDK", shell_output("#{bin}/clix --help")
  end
end
```

### 6. Test Homebrew Formula Locally

```bash
# Install from local formula
brew install --build-from-source ./clix.rb

# Test the installation
clix --help
clix --version

# Uninstall
brew uninstall clix
```

### 7. Commit and Push

```bash
git add clix.rb package.json package-lock.json
git commit -m "chore: release v<version>"
git push origin main

# Create a git tag
git tag v<version>
git push origin v<version>
```

### 8. Users Can Install

Once pushed, users can install with:

```bash
brew tap clix-so/clix-cli
brew install clix-so/clix-cli/clix
```

Or update to the latest version:

```bash
brew update
brew upgrade clix
```

## Troubleshooting

### Formula SHA256 Mismatch

If users get a SHA256 mismatch error, verify:
1. The URL in the formula points to the correct version
2. The SHA256 was calculated from the exact tarball at that URL
3. Re-download and recalculate if needed

### Installation Fails

Check:
1. Node.js 18+ is available (`brew install node@18`)
2. All dependencies are listed correctly in package.json
3. The build output in `dist/` is complete and includes all chunks

### CLI Doesn't Work After Install

1. Verify the shebang in [dist/cli.js](dist/cli.js) is correct (`#!/usr/bin/env node`)
2. Check that all dependencies are bundled (run `npm run build` and inspect dist/)
3. Test with `node dist/cli.js --help` directly
