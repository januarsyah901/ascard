# Project Structure

## Root
- `index.html`: Entry point for the browser.
- `package.json`: Dependency and script management.
- `vite.config.js`: Vite build configuration.

## src/
- `main.jsx`: React application bootstrap.
- `App.jsx`: Root component with routing and context providers.
- `index.css`: Global styles and Tailwind imports.

### pages/
- `LandingPage.jsx`: Main entry point and platform overview.
- `KnowledgeBase.jsx`: Educational content about gambling mechanics.
- `SimulationLab.jsx`: Interactive gambling simulator (Slot machine).
- `Intervention.jsx`: AI-driven cognitive intervention and myth-busting.

### components/
- `Navigation.jsx`: Global navbar with language switcher.
- `Footer.jsx`: Global footer.

### context/
- `LanguageContext.jsx`: Multi-language strings and state logic.
- `CurrencyContext.jsx`: Currency conversion (USD/IDR) and formatting.

### assets/
- Static images and local resources.

## public/
- `icon.png`: Main platform icon.
- Static assets served directly.
