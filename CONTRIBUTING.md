# Contributing to AsCard

Thank you for your interest in contributing to AsCard! We appreciate your effort to improve this project. This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Release Policy](#release-policy)
- [Reporting Issues](#reporting-issues)

## Code of Conduct

This project adheres to a Code of Conduct that all contributors are expected to follow:

- Be respectful and inclusive
- Provide constructive feedback
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Started

### Prerequisites

- Node.js >= 18.0.0
- npm >= 9.0.0
- Git

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/your-username/ascard.git
   cd ascard
   ```

3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/original-owner/ascard.git
   ```

## Development Setup

1. Install dependencies:
   ```bash
   npm install
   ```

2. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

4. Make your changes and ensure code quality:
   ```bash
   npm run lint
   ```

## Making Changes

### Branch Naming Convention

Use descriptive branch names with a prefix:

- `feature/` - New features
  - Example: `feature/add-dark-mode`
- `fix/` - Bug fixes
  - Example: `fix/card-overflow-issue`
- `docs/` - Documentation updates
  - Example: `docs/update-api-docs`
- `refactor/` - Code refactoring
  - Example: `refactor/simplify-state-management`
- `test/` - Adding or updating tests
  - Example: `test/add-simulation-tests`
- `chore/` - Maintenance tasks
  - Example: `chore/update-dependencies`

### Code Style

- Follow the existing code style
- Use functional components with React hooks
- Use descriptive variable and function names
- Keep components small and focused
- Add comments for complex logic

### ESLint

Before committing, ensure your code passes linting:

```bash
npm run lint
```

Fix linting issues automatically where possible:

```bash
npm run lint -- --fix
```

## Commit Guidelines

We follow the Conventional Commits specification for clear commit history.

### Commit Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that don't affect code meaning (formatting, semi-colons, etc.)
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `perf`: Code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to build process, dependencies, or other tools

### Scope

Optional, but recommended. Describes the section of the codebase affected:
- `components`
- `pages`
- `router`
- `styles`
- `build`

### Subject

- Use imperative mood ("add", "remove", not "adds", "removed")
- Capitalize the first letter
- Do not include period at the end
- Keep it under 50 characters

### Body

- Provide detailed explanation of changes
- Explain what and why, not how
- Wrap at 72 characters
- Reference related issues

### Footer

Include references to related issues:
```
Fixes #123
Closes #456
Related to #789
```

### Example Commits

```
feat(components): add interactive card component

Implement new interactive card component with hover effects
and animation transitions. This replaces the static card
component and provides better user engagement.

Fixes #42
```

```
fix(router): prevent page flicker on navigation

Add proper loading state handling in Router to eliminate
visual flicker when transitioning between pages.

Closes #88
```

## Pull Request Process

### Before Submitting

1. Ensure your branch is up to date with main:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. Test your changes:
   ```bash
   npm run build
   npm run lint
   ```

3. Update documentation if needed
4. Add or update tests if applicable

### Creating a Pull Request

1. Push your changes to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. Go to the original repository on GitHub
3. Click "New Pull Request"
4. Describe your changes:
   - What problem does this solve?
   - What changes were made?
   - Are there any breaking changes?
   - Screenshots or demos if applicable

### PR Title Format

Use the same convention as commits:

```
feat(components): add interactive card component
fix(router): prevent page flicker on navigation
docs: update installation instructions
```

### PR Description Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## How Has This Been Tested?
Describe the testing you've done

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex logic
- [ ] Documentation updated
- [ ] No new warnings generated
- [ ] Tests pass locally
- [ ] Build passes locally
```

### Review Process

- At least one maintainer review required
- Automated checks must pass
- All discussions must be resolved
- Maintainer approval required to merge

## Release Policy

### Versioning

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR.MINOR.PATCH** (e.g., 1.2.3)
- **MAJOR**: Incompatible API changes
- **MINOR**: Backward-compatible features
- **PATCH**: Backward-compatible bug fixes

### Release Workflow

#### Patch Release (Bug Fixes)

```bash
npm run release:patch
```

This will:
1. Bump version in package.json and VERSION file
2. Create git tag `v{version}`
3. Push to main branch with tags

#### Minor Release (New Features)

```bash
npm run release:minor
```

This will:
1. Bump version in package.json and VERSION file
2. Create git tag `v{version}`
3. Push to main branch with tags

#### Major Release (Breaking Changes)

```bash
npm run release:major
```

This will:
1. Bump version in package.json and VERSION file
2. Create git tag `v{version}`
3. Push to main branch with tags

### CHANGELOG Updates

Before each release:

1. Review changes since last release
2. Update [CHANGELOG.md](./CHANGELOG.md)
3. Categorize changes: Added, Changed, Deprecated, Removed, Fixed, Security
4. Include version number and date

Example CHANGELOG entry:

```markdown
## [1.1.0] - 2026-04-22

### Added
- New interactive card animations
- Dark mode support

### Fixed
- Page flicker issue on navigation
- Memory leak in simulation engine

### Changed
- Improved card styling and layout
```

## Reporting Issues

### Bug Reports

When reporting a bug, include:

1. **Title**: Clear, descriptive summary
2. **Environment**: Node version, OS, browser
3. **Steps to Reproduce**: Clear steps to reproduce the issue
4. **Expected Behavior**: What should happen
5. **Actual Behavior**: What actually happens
6. **Screenshots/Videos**: If applicable
7. **Additional Context**: Any other relevant information

### Feature Requests

When suggesting a feature:

1. **Title**: Clear description of feature
2. **Use Case**: Why is this feature needed?
3. **Proposed Solution**: How should it work?
4. **Alternatives**: Other solutions you've considered
5. **Additional Context**: Any other relevant information

### Issue Template

Use the provided issue templates on GitHub to ensure all necessary information is included.

## Questions?

- Check existing [issues](https://github.com/yourusername/ascard/issues)
- Review [discussions](https://github.com/yourusername/ascard/discussions)
- Open a new discussion if your question isn't covered

## Recognition

Contributors will be recognized:
- In commit history
- In pull request discussions
- In release notes for significant contributions

Thank you for contributing to AsCard!

---

**Last Updated**: April 2026 | **Maintained By**: AsCard Team
