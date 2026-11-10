# Tasks: Forge Developer Portal

**Input**: Design documents from `/specs/001-forge-portal/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

## Phase 1: Setup

- [x] T001 Replace `frontend/index.html` monolith with Vite React TS scaffold
- [x] T002 [P] Add CSS tokens and global styles
- [x] T003 [P] Add `frontend/README.md`

## Phase 2: Foundational

- [x] T004 Implement types + `ProvisioningPort`
- [x] T005 Implement Mock port + localStorage
- [x] T006 [P] Implement GitHub adapter
- [x] T007 Wire `getPort()` factory
- [x] T008 Create `AppShell` + React Router routes

## Phase 3: User Story 1 — Catalog + Request (P1)

- [x] T009 [US1] Implement `CatalogPage.tsx`
- [x] T010 [US1] Implement `RequestPage.tsx`
- [x] T011 [US1] Connect create flow to `ProvisioningPort.createRequest`

## Phase 4: User Story 2 — Status (P2)

- [x] T012 [US2] Implement `StatusPage.tsx` + `StatusBadge`
- [x] T013 [US2] Add mock status advance control

## Phase 5: User Story 3 — Settings (P3)

- [x] T014 [US3] Implement `SettingsPage.tsx`
- [x] T015 [US3] Persist settings; github mode requires creds

## Phase 6: Polish & Docs

- [x] T016 Update root `README.md`
- [x] T017 Adjust `.gitignore`
- [x] T018 Run `npm run build`
