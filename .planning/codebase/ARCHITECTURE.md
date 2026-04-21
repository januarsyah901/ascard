# Architecture

## Pattern: Single Page Application (SPA)
The project is built as a client-side SPA using React and Vite. Routing is handled on the client side, providing a seamless transition between pages without full reloads.

## Design System: Bento Grid & Editorial Layout
The UI follows a modern, premium design aesthetic characterized by:
- **Bento Grids**: Used in KnowledgeBase and LandingPage for structured content.
- **Editorial Typography**: Large headline fonts and high contrast.
- **Glassmorphism/Minimalism**: Subtle borders and clean surfaces.

## Global State
- **LanguageContext**: A React Context provider that wraps the entire application (`App.jsx`) to provide multi-language support (ID/EN) across all components.

## Routing Strategy
- **React Router (v7)**: Uses a declarative `Routes` structure in `App.jsx` to map URL paths to Page components.
