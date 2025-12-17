# Release Checklist

Use this checklist when preparing a new release.

## Pre-Release

- [ ] All tests pass (`npm run typecheck`, `npm run build`)
- [ ] Local testing complete (`npm link`, test all commands)
- [ ] CHANGELOG updated (if applicable)
- [ ] README updated (if needed)
- [ ] Version bump decision made (patch/minor/major)

## Release Steps

### Option 1: Automated Release (Recommended)

1. **Update version and push to main**
   ```bash
   # Update version in package.json
   npm version patch  # or minor, or major

   # Push to main (this automatically triggers the release)
   git push origin main
   ```

2. **Monitor GitHub Actions**
   - Go to [Actions tab](https://github.com/clix-so/homebrew-clix-cli/actions)
   - Wait for Release workflow to complete
   - The workflow will automatically:
     - ✅ Detect version change
     - ✅ Create git tag
     - ✅ Build the project
     - ✅ Publish to npm
     - ✅ Calculate SHA256
     - ✅ Update Homebrew formula
     - ✅ Create GitHub release

3. **Verify the release**
   ```bash
   # Check npm
   npm view @clix-so/clix-cli version

   # Test Homebrew installation (in a clean environment)
   brew tap clix-so/clix-cli
   brew install clix-so/clix-cli/clix
   clix --version
   ```

### Option 2: Manual Release

See [DEPLOYMENT.md](../DEPLOYMENT.md#manual-deployment) for detailed manual steps.

## Post-Release

- [ ] Verify npm package: https://www.npmjs.com/package/@clix-so/clix-cli
- [ ] Verify GitHub release: https://github.com/clix-so/homebrew-clix-cli/releases
- [ ] Test Homebrew installation from clean environment
- [ ] Announce release (if applicable)

## Troubleshooting

### GitHub Actions fails at "Publish to npm"

**Check:**
- Is `NPM_TOKEN` secret still valid?
- Is the version number unique (not already published)?
- Are there any build errors?

**Fix:**
- Regenerate NPM_TOKEN if expired
- Update version number if duplicate
- Fix build errors and re-tag

### Formula update fails

**Check:**
- Did npm publish succeed?
- Is the version number correct?
- Are there sed syntax errors in the workflow?

**Fix:**
- Verify npm package exists
- Update workflow if needed
- Manually update clix.rb and commit

### Users report installation issues

**Check:**
- Does `brew install` work in a clean environment?
- Is the SHA256 correct?
- Are all dependencies bundled?

**Fix:**
- Recalculate SHA256 from the actual npm tarball
- Update clix.rb with correct hash
- Verify dist/ contains all necessary files

## Emergency Rollback

If a release is broken and needs to be rolled back:

```bash
# Unpublish from npm (within 72 hours only)
npm unpublish @clix-so/clix-cli@VERSION

# Delete git tag
git tag -d vVERSION
git push origin :refs/tags/vVERSION

# Delete GitHub release (via GitHub UI)
# Revert clix.rb to previous version
git revert <commit-hash>
git push origin main
```

## Version Guidelines

Follow [Semantic Versioning](https://semver.org/):

- **PATCH** (1.0.x): Bug fixes, minor changes, no breaking changes
- **MINOR** (1.x.0): New features, no breaking changes
- **MAJOR** (x.0.0): Breaking changes, major refactors

## Release Frequency

- **Patch releases**: As needed for bug fixes
- **Minor releases**: Every few weeks for new features
- **Major releases**: Only when necessary for breaking changes
