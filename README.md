# Bongkar Ilusi (Breaking Illusions)

> An interactive educational web application that reveals cognitive biases and the mathematical reality behind probability-based systems through hands-on simulations and data-driven insights.

**Tagline**: *Data Driven • Clinically Validated • Non-Profit Mission*

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Node.js](https://img.shields.io/badge/Node.js->=18.0.0-green.svg)](https://nodejs.org/)
[![React](https://img.shields.io/badge/React-^19.2.5-blue.svg)](https://react.dev/)
[![Vite](https://img.shields.io/badge/Vite-^8.0.0-purple.svg)](https://vitejs.dev/)

## 📋 Table of Contents

- [Project Overview](#project-overview)
- [What is Bongkar Ilusi?](#what-is-bongkar-ilusi)
- [How It Works](#how-it-works)
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

## 🎯 Project Overview

**Bongkar Ilusi** (Breaking Illusions) is an interactive educational platform designed to teach users about cognitive biases, probability misconceptions, and the mathematical reality behind probability-based systems. Through immersive simulations and data-driven interventions, the application helps people understand:

- How cognitive biases affect decision-making
- The statistical reality of random events
- Why the "house edge" is mathematically unbeatable
- Common gambling fallacies (Gambler's Fallacy, Chasing Losses, Hot Hand Fallacy)
- The psychological mechanisms that lead to poor financial decisions

### Mission

To provide a non-profit, educational tool that raises awareness about cognitive biases and helps individuals make better-informed decisions through evidence-based learning.

### Target Audience

- Students & Educators
- Psychology & Finance professionals
- Individuals interested in behavioral economics
- Anyone seeking to understand probability and decision-making biases
- General public with curiosity about financial literacy

## 🧠 What is Bongkar Ilusi?

Bongkar Ilusi is divided into three main interactive modules:

### 1. **Simulation Lab** 🎲

A risk-free virtual environment where users can:
- Simulate betting on a system with a 96.4% RTP (Return to Player)
- Experience the mathematics of probability firsthand
- Observe their own behavior and biases in real-time
- Track their "losses" and emotional responses
- Receive real-time AI analyst commentary on their cognitive state

**Key Insights:**
- Only 0.05% win rate (1 in 2000) emphasizes how rarely winning occurs
- Virtual balance prevents real financial loss while providing authentic emotional experience
- RNG (Random Number Generator) transparency shows the algorithm is fair
- "Tilt Level" monitoring reveals emotional escalation patterns

### 2. **Intervention System** 🛡️

Evidence-based interventions designed to interrupt negative patterns:
- Cognitive reframing techniques
- Statistical reality checks
- Real-time psychological support messages
- Decision-making frameworks
- Resources for problematic behavior

### 3. **Knowledge Base** 📚

Comprehensive educational content including:
- Behavioral economics fundamentals
- Common cognitive biases explained
- Mathematical probability basics
- Historical case studies
- Scientific research summaries
- Practical financial literacy guides

---

## ✨ Features

### Core Educational Modules

- **Interactive Simulation Lab**: Risk-free virtual environment to experience probability systems
- **Real-Time AI Analyst**: Psychological monitoring and behavioral feedback during simulations
- **Intervention System**: Evidence-based cognitive interventions and support
- **Comprehensive Knowledge Base**: Educational content on behavioral economics and cognitive biases
- **Balance Analytics**: Visual representation of financial trends and volatility patterns
- **RNG Transparency**: Real-time view of random number generation for algorithmic transparency
- **Tilt Level Monitoring**: Emotional state tracking and cognitive bias indicators
- **Reality Check Prompts**: Periodic intervention messages based on behavior patterns

### Technical Features

- **Multi-Language Support**: LanguageContext for internationalization
- **Multi-Currency Support**: CurrencyContext for global accessibility
- **Responsive Design**: Mobile-first responsive layout with Tailwind CSS
- **Modern Build System**: Vite for ultra-fast development and optimized builds
- **Code Quality**: ESLint configuration for consistent code standards
- **Client-Side Processing**: All calculations run locally for privacy

---

## 🛠️ How It Works

### Simulation Mechanism

1. **User Setup**: Player starts with virtual balance (Rp 5,000,000 / ~$333 USD equivalent)
2. **Betting**: User selects or inputs custom bet amount
3. **Simulation Spin**: System generates random outcome with configurable RTP
4. **RNG Algorithm**: 
   - Displays real-time probability calculations
   - Shows win/loss outcome
   - Updates balance accordingly
5. **Behavior Monitoring**: System tracks:
   - Consecutive losses (triggers Tilt Level)
   - Betting patterns
   - Psychological state indicators
6. **Interventions**: 
   - Every 15 spins: Reality Check modal
   - Escalating messages based on Tilt Level (Low → Medium → High → Critical)
   - AI Analyst commentary on behavioral patterns

### Cognitive Bias Framework

The application helps users identify common cognitive biases:

| Bias | Description | How Detected |
|------|-------------|--------------|
| **Gambler's Fallacy** | Belief that past outcomes affect future probability | Consecutive losses > 8 |
| **Chasing Losses** | Increasing bets to recover losses | Bet escalation pattern |
| **Hot Hand Fallacy** | Belief in "winning streaks" | Positive balance reinforcement |
| **Loss Aversion** | Emotional response to losses > gains | Tilt Level increase |
| **Sunk Cost Fallacy** | Continuing due to past investment | Session duration |

---

## 📊 Tech Stack

- **Frontend Framework**: React 19.2.5
- **Build Tool**: Vite 8.0.9
- **Styling**: Tailwind CSS 4.2.3
- **Routing**: React Router DOM 7.14.1
- **Icons**: Lucide React 1.8.0
- **Linting**: ESLint 9.39.4
- **Package Manager**: npm
- **Language**: JavaScript/JSX

---

## 📦 Installation

### Prerequisites

- Node.js >= 18.0.0
- npm >= 9.0.0

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/januarsyah901/ascard.git
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

---

## 📁 Project Structure

```
ascard/
├── src/
│   ├── pages/
│   │   ├── LandingPage.jsx          # Hero & project overview
│   │   ├── SimulationLab.jsx        # Main betting simulator
│   │   ├── KnowledgeBase.jsx        # Educational content
│   │   └── Intervention.jsx         # Cognitive interventions
│   │
│   ├── components/
│   │   ├── Navigation.jsx           # Top navigation
│   │   └── Footer.jsx               # Footer
│   │
│   ├── context/
│   │   ├── LanguageContext.jsx      # i18n support
│   │   └── CurrencyContext.jsx      # Currency conversion
│   │
│   ├── assets/                      # Images & media
│   ├── App.jsx                      # Main app component
│   ├── main.jsx                     # React entry point
│   ├── App.css                      # App styles
│   └── index.css                    # Global styles
│
├── public/                          # Static files
├── dist/                            # Production build
├── vite.config.js                   # Build config
├── eslint.config.js                 # Linting rules
├── package.json                     # Dependencies
├── README.md                        # This file
├── CHANGELOG.md                     # Version history
├── CONTRIBUTING.md                  # Contributing guide
├── RELEASE.md                       # Release guide
├── VERSION                          # Current version
└── .version-config.json             # Version config
```

---

## 🔄 Development

### Code Quality Standards

- Follow ESLint rules configured in `eslint.config.js`
- Use functional components and React hooks
- Maintain consistent code formatting
- Write meaningful commit messages (follow Conventional Commits)

### Hot Module Replacement (HMR)

Development server includes HMR for instant updates without page reload.

### Adding New Features

1. Create components in `src/components/`
2. Create pages in `src/pages/` if needed
3. Use context for global state (Language, Currency)
4. Maintain responsive design with Tailwind CSS
5. Update tests and documentation

---

## 🏗️ Building

The build process:
1. Bundles React components and dependencies
2. Optimizes assets and images
3. Generates production-ready output in `dist/`

For detailed build configuration, see `vite.config.js`

---

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

---

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

### Release Commands

```bash
# Patch release (1.0.0 → 1.0.1)
npm run release:patch

# Minor release (1.0.0 → 1.1.0)
npm run release:minor

# Major release (1.0.0 → 2.0.0)
npm run release:major
```

For detailed instructions, see [RELEASE.md](./RELEASE.md)

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for detailed guidelines.

### Quick Contribution Steps

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Contribution Areas

- **Bug fixes**: Report issues and submit fixes
- **New features**: Suggest and implement improvements
- **Documentation**: Improve guides and tutorials
- **Translations**: Add language support
- **Content**: Contribute educational materials
- **UX/UI**: Enhance user experience

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

---

## 📧 Support

For issues, questions, or suggestions:
- Open an [issue](https://github.com/januarsyah901/ascard/issues)
- Check [discussions](https://github.com/januarsyah901/ascard/discussions)
- Review [documentation](./docs)

### Resources

- [CHANGELOG](./CHANGELOG.md) - Version history
- [CONTRIBUTING](./CONTRIBUTING.md) - How to contribute
- [RELEASE](./RELEASE.md) - Release management
- [PROJECT](https://github.com/januarsyah901/ascard) - Main repository

---

## 🎓 Educational Value

Bongkar Ilusi is designed with evidence-based learning principles:

- **Interactive Learning**: Hands-on experience with probability
- **Real-Time Feedback**: Immediate psychological insights
- **Behavior Monitoring**: Identify your own cognitive patterns
- **Evidence-Based**: Grounded in behavioral economics research
- **Risk-Free**: Virtual money prevents real financial harm
- **Accessible**: Free, non-profit educational tool

---

## 🔐 Privacy & Security

- **No Data Collection**: All calculations run locally in your browser
- **No Accounts Required**: Complete anonymity
- **No Tracking**: No cookies or analytics
- **Open Source**: Fully transparent code
- **HTTPS Only**: Secure connections

---

## 🚀 Roadmap

Future features in development:

- [ ] Multiplayer simulation sessions
- [ ] Advanced analytics dashboard
- [ ] AI-powered personalized interventions
- [ ] Integration with financial literacy programs
- [ ] Mobile-native app (iOS/Android)
- [ ] Offline mode support
- [ ] Advanced cognitive assessment tools
- [ ] Research collaboration features

---

**Last Updated**: April 2026 | **Version**: 1.0.0 | **Repository**: [github.com/januarsyah901/ascard](https://github.com/januarsyah901/ascard)
