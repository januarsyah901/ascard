# AsCard

> A modern React application for interactive card-based experiences with simulations, knowledge base, and interventions.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Node.js](https://img.shields.io/badge/Node.js->=18.0.0-green.svg)](https://nodejs.org/)
[![React](https://img.shields.io/badge/React-^19.2.5-blue.svg)](https://react.dev/)
[![Vite](https://img.shields.io/badge/Vite-^8.0.0-purple.svg)](https://vitejs.dev/)

## 📋 Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [Development](#development)
- [Building](#building)
- [Version Management](#version-management)
- [Release Process](#release-process)
- [Contributing](#contributing)
- [License](#license)

## ✨ Features

- **Interactive Card Interface**: Engaging card-based UI for presenting content
- **Simulation Engine**: Run interactive simulations with real-time feedback
- **Knowledge Base**: Comprehensive documentation and knowledge management
- **Intervention System**: Smart intervention mechanisms for user engagement
- **Responsive Design**: Mobile-first responsive layout with Tailwind CSS
- **Modern Build System**: Vite for ultra-fast development and optimized builds
- **Code Quality**: ESLint configuration for consistent code standards

## 🛠️ Tech Stack

- **Frontend Framework**: React 19.2.5
- **Build Tool**: Vite 8.0.9
- **Styling**: Tailwind CSS 4.2.3
- **Routing**: React Router DOM 7.14.1
- **Icons**: Lucide React 1.8.0
- **Linting**: ESLint 9.39.4
- **Package Manager**: npm

## 📦 Installation

### Prerequisites

- Node.js >= 18.0.0
- npm >= 9.0.0

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/ascard.git
   cd ascard
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

## 🚀 Quick Start

### Development Server

Start the development server with hot module replacement:

```bash
npm run dev
```

The application will be available at `http://localhost:5173`

### Build for Production

Create an optimized production build:

```bash
npm run build
```

Output will be in the `dist/` directory.

### Preview Production Build

Preview the production build locally:

```bash
npm run preview
```

### Linting

Check code quality:

```bash
npm run lint
```

## 📁 Project Structure

```
ascard/
├── src/
│   ├── components/        # Reusable React components
│   ├── pages/            # Page components
│   ├── App.jsx           # Main App component
│   └── main.jsx          # Application entry point
├── public/               # Static assets
├── dist/                 # Production build output
├── index.html           # HTML template
├── intervention.html    # Intervention page template
├── knowledge.html       # Knowledge page template
├── landing.html         # Landing page template
├── simulation.html      # Simulation page template
├── vite.config.js       # Vite configuration
├── eslint.config.js     # ESLint configuration
├── package.json         # Project dependencies
├── CHANGELOG.md         # Version history and changes
├── CONTRIBUTING.md      # Contribution guidelines
├── VERSION              # Current version
└── README.md            # This file
```

## 🔄 Development

### Code Quality Standards

- Follow ESLint rules configured in `eslint.config.js`
- Use functional components and React hooks
- Maintain consistent code formatting
- Write meaningful commit messages

### Hot Module Replacement (HMR)

Development server includes HMR for instant updates without page reload.

## 🏗️ Building

The build process:
1. Bundles React components and dependencies
2. Optimizes assets and images
3. Generates production-ready output in `dist/`

For detailed build configuration, see `vite.config.js`

## 📦 Version Management

This project follows **Semantic Versioning (SemVer)** format: `MAJOR.MINOR.PATCH`

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes

Current version: See [VERSION](./VERSION) file

### Version Tracking

- `VERSION` file contains the current version
- `package.json` version field is automatically kept in sync
- All releases are tagged in git as `v{version}`

## 🎯 Release Process

### Automated Release Workflow

Releases follow a structured process:

1. **Update Version**
   ```bash
   npm version [major|minor|patch]
   ```

2. **Update CHANGELOG**
   - Document changes in [CHANGELOG.md](./CHANGELOG.md)
   - Include date, version, and categorized changes

3. **Create Git Tag**
   ```bash
   git tag -a v{version} -m "Release version {version}"
   git push origin v{version}
   ```

4. **Build Release**
   ```bash
   npm run build
   ```

5. **Create Release Notes**
   - Use git tags and CHANGELOG as source
   - Include installation instructions
   - Link to CHANGELOG entry

### Hotfix Releases

For critical fixes in production:

1. Create branch from latest release tag
   ```bash
   git checkout -b hotfix/issue-description v{version}
   ```

2. Make necessary fixes
3. Bump PATCH version
4. Follow standard release process

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for detailed guidelines.

### Quick Contribution Steps

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

## 📧 Support

For issues, questions, or suggestions:
- Open an [issue](https://github.com/yourusername/ascard/issues)
- Check [discussions](https://github.com/yourusername/ascard/discussions)
- Review [documentation](./docs)

---

**Last Updated**: April 2026 | **Version**: 1.0.0
