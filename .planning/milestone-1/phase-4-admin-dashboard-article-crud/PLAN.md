# Plan: Phase 4 - Admin Dashboard & Article CRUD

## Objective
Implement a secure admin dashboard with article management powered by Prisma and PostgreSQL.

## Waves

### Wave 1: Infrastructure & Database
- [ ] Initialize Prisma with PostgreSQL provider.
- [ ] Define `Article` and `Admin` models in `schema.prisma`.
- [ ] Setup a Node.js/Express server in `server/` directory.
- [ ] Implement database connection and initial migration.

### Wave 2: Backend API & Auth
- [ ] Implement Admin Login endpoint with JWT.
- [ ] Create CRUD endpoints for Articles (GET, POST, PUT, DELETE).
- [ ] Setup authentication middleware for protected routes.

### Wave 3: Frontend Dashboard UI
- [ ] Create `/admin/login` page with premium dark mode styling.
- [ ] Create `/admin/dashboard` layout and navigation.
- [ ] Implement Article List view with delete/edit actions.
- [ ] Implement Article Editor (Create/Update) with markdown support.

### Wave 4: Integration & Polish
- [ ] Connect Frontend to Backend API via Axios/Fetch.
- [ ] Add loading states and toast notifications for CRUD actions.
- [ ] Final UI/UX audit for consistency with the "Bongkar Ilusi" brand.

## Verification
- Admin can login with valid credentials.
- Unauthorized users cannot access the dashboard.
- Articles can be created, edited, and deleted.
- Data persists in the PostgreSQL database.
