# Development Conventions

## Code Style
- **Functional Components**: All React components are defined as functional components using `export default function`.
- **Hooks**: Extensive use of `useState`, `useEffect`, and custom hooks like `useLanguage`.
- **JSX**: standard React 19 JSX syntax.

## Styling
- **Utility-First**: Tailwind CSS is used for almost all styling.
- **Custom Design Tokens**: Colors like `primary`, `on-primary`, `surface`, etc., are likely mapped to a design system (Material 3 style).
- **Responsive Design**: Mobile-first approach using Tailwind's `md:`, `lg:` prefixes.

## Localization & Currency
- All user-facing strings MUST be moved to `LanguageContext.jsx`.
- Use the `t` object from `useLanguage()` to access localized strings.
- All monetary values MUST use `formatValue()` from `useCurrency()` to ensure proper conversion between USD and IDR.

## File Naming
- **Components/Pages**: PascalCase (`LandingPage.jsx`).
- **Context/Hooks**: PascalCase or camelCase depending on usage (`LanguageContext.jsx`).
- **Assets**: kebab-case or underscore_case.
