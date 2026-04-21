# Codebase Concerns

## Technical Debt
- **Tailwind CDN**: The project currently relies on the Tailwind CDN in `index.html`. This should be removed in favor of the Tailwind CLI/PostCSS pipeline for production.
- **State Persistence**: User balance and logs in `SimulationLab` are reset on page refresh. Implementing `localStorage` persistence is a high-priority enhancement.
- **Static HTML Files**: There are legacy/placeholder HTML files in the root (`landing.html`, `knowledge.html`, etc.) that should be cleaned up as the React app is now the primary entry point.

## Scalability
- **LanguageContext Size**: As the platform grows, the `LanguageContext.jsx` file will become very large. Consider splitting translations into separate JSON files (e.g., `id.json`, `en.json`).

## Performance
- **Image Optimization**: Some images are loaded from external URLs. These should be localized or optimized for faster load times.
