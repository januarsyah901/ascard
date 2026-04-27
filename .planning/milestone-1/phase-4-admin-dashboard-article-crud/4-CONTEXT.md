# Phase 4 Context: Admin Dashboard & Article CRUD

## Decisions
- **Database**: PostgreSQL (Local, Port 5432)
- **ORM**: Prisma
- **Backend**: Node.js/Express (to be added as a separate server or integrated with Vite dev server)
- **Authentication**: Simple Admin Login (JWT-based or Session)
- **UI/UX**: Premium Dark Mode dashboard (matching the Clinical Authority aesthetic)
- **Article Schema**:
  - id (UUID)
  - title (String)
  - content (Markdown/Text)
  - category (String)
  - createdAt/updatedAt (DateTime)

## Architecture Changes
- Add `server/` directory for the backend logic.
- Add `prisma/` directory for schema and migrations.
- Update `package.json` with backend dependencies.
- Update `Vite` config if needed to proxy API requests.
