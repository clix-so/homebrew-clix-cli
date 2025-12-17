# GitHub Actions Setup Guide

This guide explains how to configure GitHub secrets for automated releases.

## Required Secrets

### 1. NPM_TOKEN

This token is required to publish packages to npm.

#### Steps to create NPM_TOKEN:

1. **Login to npm**
   ```bash
   npm login
   ```

2. **Generate an access token**
   - Go to https://www.npmjs.com/settings/YOUR_USERNAME/tokens
   - Click "Generate New Token" → "Classic Token"
   - Select "Automation" type
   - Give it a name like "github-actions-clix-cli"
   - Copy the generated token (it will only be shown once)

3. **Add to GitHub Secrets**
   - Go to your GitHub repository
   - Navigate to Settings → Secrets and variables → Actions
   - Click "New repository secret"
   - Name: `NPM_TOKEN`
   - Value: Paste your npm token
   - Click "Add secret"

### 2. GITHUB_TOKEN (automatically provided)

The `GITHUB_TOKEN` is automatically provided by GitHub Actions and doesn't need manual setup. It's used to:
- Create GitHub releases
- Commit and push the updated Homebrew formula

## Workflow Permissions

The workflow requires the following permissions (already configured in [release.yml](../workflows/release.yml)):

```yaml
permissions:
  contents: write    # To create releases and push commits
  id-token: write    # For npm provenance
```

## Testing the Workflow

### 1. Manual Trigger (Recommended for first test)

You can manually trigger the workflow:

1. Go to GitHub → Actions → Release workflow
2. Click "Run workflow"
3. Select the branch (usually `main`)
4. Click "Run workflow"

**Note**: The workflow will only proceed if the version in package.json has changed from the previous commit.

### 2. Automatic Release (Standard workflow)

The workflow automatically triggers when you push a version change:

```bash
# Update version in package.json
npm version patch  # or minor, or major

# Push to main (this automatically triggers the release)
git push origin main
```

The workflow will:
- Detect the version change in package.json
- Automatically create a git tag (e.g., v1.0.1)
- Proceed with build, publish, and release steps

### 3. Monitor the Workflow

1. Go to GitHub → Actions
2. Click on the running workflow
3. Watch each step complete:
   - ✅ Build
   - ✅ Publish to npm
   - ✅ Calculate SHA256
   - ✅ Update Homebrew formula
   - ✅ Create GitHub release

## What the Workflow Does

1. **Builds the project**
   - Runs `npm ci` to install dependencies
   - Runs `npm run build` to create the distribution files

2. **Publishes to npm**
   - Publishes `@clix-so/clix-cli` to npm registry
   - Uses provenance for supply chain security

3. **Updates Homebrew Formula**
   - Downloads the published tarball from npm
   - Calculates SHA256 hash
   - Updates [clix.rb](../../clix.rb) with new version and hash
   - Commits and pushes the changes

4. **Creates GitHub Release**
   - Creates a release on GitHub with release notes
   - Attaches built files and the updated formula

## Troubleshooting

### NPM_TOKEN Invalid or Expired

Error: `npm ERR! code E401`

**Solution**:
1. Generate a new npm token (see steps above)
2. Update the `NPM_TOKEN` secret in GitHub

### Permission Denied When Pushing

Error: `remote: Permission to clix-so/homebrew-clix-cli.git denied`

**Solution**:
1. Go to Settings → Actions → General
2. Under "Workflow permissions", select:
   - ✅ Read and write permissions
3. Click Save

### Formula Update Failed

Error: `sed: no such file or directory`

**Solution**:
The workflow expects `clix.rb` in the repository root. Ensure it exists and is committed.

### Package Not Found on npm

Error: `npm ERR! code E404`

**Solution**:
1. Wait a few seconds (npm propagation delay)
2. Verify the package was published: https://www.npmjs.com/package/@clix-so/clix-cli
3. Check that the version in package.json matches

## Manual Rollback

If a release fails and you need to rollback:

```bash
# Unpublish from npm (within 72 hours)
npm unpublish @clix-so/clix-cli@VERSION

# Delete the git tag
git tag -d vVERSION
git push origin :refs/tags/vVERSION

# Delete the GitHub release (via GitHub UI)
```

## Security Best Practices

1. **Never commit tokens to the repository**
2. **Use "Automation" type tokens** for npm (not "Publish" tokens)
3. **Regularly rotate npm tokens** (every 3-6 months)
4. **Enable 2FA on your npm account**
5. **Review workflow logs** for any suspicious activity
