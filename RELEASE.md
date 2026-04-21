# Release Management Guide

This document provides detailed instructions for managing releases in the AsCard project.

## Table of Contents

- [Version Strategy](#version-strategy)
- [Pre-Release Checklist](#pre-release-checklist)
- [Release Types](#release-types)
- [Release Steps](#release-steps)
- [Post-Release](#post-release)
- [Hotfix Releases](#hotfix-releases)
- [Git Tagging](#git-tagging)
- [Troubleshooting](#troubleshooting)

## Version Strategy

### Semantic Versioning

AsCard follows [Semantic Versioning 2.0.0](https://semver.org/):

```
MAJOR.MINOR.PATCH
  |      |     |
  |      |     └─ Bug fixes (backward compatible)
  |      └────── New features (backward compatible)
  └──────────── Breaking changes (incompatible API)
```

### Version Examples

- `1.0.0` - Initial release
- `1.1.0` - New features added
- `1.1.1` - Bug fix
- `2.0.0` - Major breaking changes
- `1.0.0-beta.1` - Pre-release version

### Current Version

Check current version:

```bash
cat VERSION                    # Standalone version file
cat package.json | grep version  # Package version
node -e "console.log(require('./package.json').version)"
```

## Pre-Release Checklist

Before creating any release, ensure:

- [ ] All changes are committed
- [ ] Branch is clean (`git status` shows no uncommitted changes)
- [ ] Code passes linting: `npm run lint`
- [ ] Build succeeds: `npm run build`
- [ ] CHANGELOG.md is updated
- [ ] No security vulnerabilities: `npm audit`
- [ ] All tests pass (if applicable)
- [ ] Branch is up-to-date with main: `git pull origin main`

## Release Types

### PATCH Release

**Use for**: Bug fixes and minor improvements

Increment: `1.0.0` → `1.0.1`

```bash
npm run release:patch
```

**What it does**:
- Updates version in package.json
- Updates VERSION file
- Creates git commit with version bump
- Creates git tag `v1.0.1`
- Pushes to origin with tags

### MINOR Release

**Use for**: New features (backward compatible)

Increment: `1.0.0` → `1.1.0`

```bash
npm run release:minor
```

**What it does**:
- Updates version in package.json
- Updates VERSION file
- Creates git commit with version bump
- Creates git tag `v1.1.0`
- Pushes to origin with tags

### MAJOR Release

**Use for**: Breaking changes

Increment: `1.0.0` → `2.0.0`

```bash
npm run release:major
```

**What it does**:
- Updates version in package.json
- Updates VERSION file
- Creates git commit with version bump
- Creates git tag `v2.0.0`
- Pushes to origin with tags

## Release Steps

### Step 1: Prepare Changes

1. Ensure all changes are committed:
   ```bash
   git status
   ```

2. Build and test:
   ```bash
   npm run build
   npm run lint
   ```

3. Verify no uncommitted changes:
   ```bash
   git diff
   git diff --staged
   ```

### Step 2: Update CHANGELOG

1. Open `CHANGELOG.md`
2. Replace `[Unreleased]` section with new version:
   ```markdown
   ## [1.1.0] - 2026-04-22

   ### Added
   - Feature description
   - Another feature

   ### Fixed
   - Bug fix description

   ### Changed
   - Change description
   ```

3. Add new Unreleased section at top:
   ```markdown
   ## [Unreleased]

   ### Added
   - (Placeholder for unreleased features)

   ### Changed
   - (Placeholder for unreleased changes)
   ```

4. Commit CHANGELOG changes:
   ```bash
   git add CHANGELOG.md
   git commit -m "docs: update CHANGELOG for v1.1.0"
   ```

### Step 3: Create Release

Choose appropriate release command:

```bash
# For bug fixes
npm run release:patch

# For new features
npm run release:minor

# For breaking changes
npm run release:major
```

This will:
- Update package.json version
- Update VERSION file
- Create version commit
- Create version tag
- Push to remote with tags

### Step 4: Verify Release

Confirm release was created:

```bash
# Check git tag
git tag -l | tail -5

# Check package.json
cat package.json | grep version

# Check VERSION file
cat VERSION

# View release commit
git log -1
```

### Step 5: Create Release Notes (Manual)

On GitHub:

1. Go to [Releases](https://github.com/yourusername/ascard/releases)
2. Click "Draft a new release"
3. Select the tag just created (e.g., `v1.1.0`)
4. Title: `Release v1.1.0`
5. Description:
   ```markdown
   ## What's New

   [Copy relevant sections from CHANGELOG]

   ### Installation

   ```bash
   npm install
   ```

   ### Notable Changes

   - Feature 1
   - Feature 2

   ### Bug Fixes

   - Fix 1
   - Fix 2

   ### Contributors

   Thanks to all contributors!
   ```

6. Click "Publish release"

## Post-Release

### After Successful Release

1. **Verify on npm** (if published):
   ```bash
   npm info ascard
   ```

2. **Announce release** (if applicable):
   - GitHub Releases page
   - Project documentation
   - Community channels

3. **Monitor issues**:
   - Watch for bug reports
   - Address critical issues quickly
   - Plan next release

### Updating Dependent Projects

If other projects depend on this one:

1. Update dependencies in dependent projects
2. Test compatibility
3. Create releases as needed

## Hotfix Releases

### When to Use

Use hotfixes for critical bugs in production that can't wait for the next regular release.

### Hotfix Process

1. **Create hotfix branch from latest tag**:
   ```bash
   git fetch origin
   git checkout -b hotfix/critical-bug origin/v1.0.0
   ```

2. **Fix the issue**:
   - Make minimal changes
   - Test thoroughly
   - Commit with proper message

3. **Update CHANGELOG**:
   ```markdown
   ## [1.0.1] - 2026-04-22

   ### Fixed
   - Critical bug description
   ```

4. **Create patch release**:
   ```bash
   npm run release:patch
   ```

5. **Merge back to main** (if using main branch):
   ```bash
   git checkout main
   git pull origin main
   git merge hotfix/critical-bug
   ```

## Git Tagging

### Manual Tag Creation

If automated process fails:

```bash
# Create annotated tag
git tag -a v1.1.0 -m "Release version 1.1.0"

# Push tag
git push origin v1.1.0

# Push all tags
git push origin --tags
```

### Tag Naming Convention

- Format: `v{MAJOR}.{MINOR}.{PATCH}`
- Examples: `v1.0.0`, `v1.2.3`, `v2.0.0-beta.1`

### List All Tags

```bash
# All tags
git tag -l

# Latest tag
git describe --tags --abbrev=0

# Tags with descriptions
git tag -l -n10
```

### View Tag Details

```bash
git show v1.0.0
git log v1.0.0..v1.1.0    # Changes between versions
```

### Delete Tag

```bash
# Delete local tag
git tag -d v1.0.0

# Delete remote tag
git push origin --delete v1.0.0
```

## Troubleshooting

### Issue: "Cannot find main branch"

```bash
# Check current branch
git branch -a

# Push to correct branch
git push origin HEAD:main --tags
```

### Issue: "Tag already exists"

```bash
# Delete existing tag
git tag -d v1.0.0
git push origin --delete v1.0.0

# Create new tag
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

### Issue: "Changes not committed"

```bash
# Check status
git status

# Stage changes
git add .

# Commit
git commit -m "feat: description"

# Try release again
npm run release:patch
```

### Issue: "VERSION file out of sync"

```bash
# Check current version
cat VERSION
cat package.json | grep version

# Manually update VERSION
echo "1.1.0" > VERSION

# Commit and create tag
git add VERSION package.json
git commit -m "bump: version 1.1.0"
git tag -a v1.1.0 -m "Release version 1.1.0"
git push origin main --tags
```

### Issue: Build fails before release

```bash
# Check errors
npm run build
npm run lint

# Fix errors
# ... make corrections ...

# Retry build and release
npm run build
npm run release:patch
```

## Release Calendar

Plan major releases:

- Plan features and breaking changes
- Set release dates
- Announce timeline to community
- Prepare detailed release notes
- Test extensively before release

## Rollback Procedure

If critical issue is found in release:

1. **Identify issue** and create fix
2. **Revert tag** (if possible):
   ```bash
   git revert v1.1.0
   git push origin main --tags
   ```

3. **Create hotfix**:
   ```bash
   npm run release:patch
   ```

4. **Communicate with users**:
   - Announce issue on GitHub Issues
   - Release notice of rollback
   - Provide migration path if needed

## Best Practices

- Keep releases focused and coherent
- Write clear CHANGELOG entries
- Use semantic versioning consistently
- Tag all releases in git
- Create annotated tags with messages
- Document breaking changes prominently
- Test releases before announcing
- Plan major versions well in advance
- Respond quickly to critical issues
- Keep version history clean

---

**Last Updated**: April 2026 | **Maintained By**: AsCard Team
