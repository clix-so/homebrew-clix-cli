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

### 2. GitHub App Token (Required to Bypass Repository Rules)

Since the repository has branch protection rules that prevent direct pushes to `main`, we use a GitHub App to bypass these rules for the release workflow.

#### Steps to create and configure GitHub App:

1. **Create a GitHub App**
   - Go to https://github.com/organizations/clix-so/settings/apps
   - Click "New GitHub App"
   - Fill in the details:
     - **Name**: `github-actions-clix-cli`
     - **Homepage URL**: `https://github.com/clix-so/homebrew-clix-cli`
     - **Webhook**: Uncheck "Active"
     - **Permissions**:
       - Repository permissions:
         - Contents: **Read and write**
         - Metadata: **Read-only**
     - **Where can this GitHub App be installed?**: Only on this account
   - Click "Create GitHub App"

2. **Generate Private Key**
   - After creating the app, scroll to "Private keys" section
   - Click "Generate a private key"
   - Save the downloaded `.pem` file securely

3. **Install the App**
   - Go to the app settings page
   - Click "Install App" in the left sidebar
   - Select your organization (`clix-so`)
   - Choose "Only select repositories"
   - Select `homebrew-clix-cli` repository
   - Click "Install"

4. **Add Secrets to Repository**
   - Go to https://github.com/clix-so/homebrew-clix-cli/settings/secrets/actions
   - Add two secrets:
     - **CLIX_APP_ID**: Copy the "App ID" from the app settings page
     - **CLIX_APP_PRIVATE_KEY**: Copy the entire contents of the `.pem` file (including BEGIN/END lines)

5. **Configure Repository Rules to Bypass the App**
   - Go to https://github.com/clix-so/homebrew-clix-cli/settings/rules
   - Edit the rule that blocks direct pushes to main
   - Scroll to "Bypass list"
   - Click "Add bypass" → "Apps" → Select `github-actions-clix-cli`
   - Save the rule

### 3. GITHUB_TOKEN (automatically provided)

The `GITHUB_TOKEN` is automatically provided by GitHub Actions and is used for:
- Creating GitHub releases

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

### Repository Rule Violations

Error: `push declined due to repository rule violations`

**Solution**:
1. Ensure you've completed the GitHub App setup (see section 2 above)
2. Verify the app is added to the bypass list in repository rules
3. Check that `CLIX_APP_ID` and `CLIX_APP_PRIVATE_KEY` secrets are set correctly

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
