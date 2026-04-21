# Testing Strategy

## Automated Testing
- **Status**: Not yet implemented.
- **Future Recommendation**: 
    - **Vitest**: For unit testing logic (especially RNG and i18n logic).
    - **Playwright/Cypress**: For E2E testing of the simulation and intervention flows.

## Manual Verification
- **Cross-browser checks**: Ensuring Tailwind 4 features render correctly.
- **Localization Audit**: verifying all text changes correctly when the language toggle is clicked.
- **Game Logic**: Verifying spin mechanics, balance updates, and log accuracy in `SimulationLab`.
