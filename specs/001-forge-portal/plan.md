# Implementation Plan: Forge Developer Portal

**Branch**: `001-forge-portal` | **Date**: 2026-07-25 | **Spec**: [spec.md](./spec.md)

## Summary

Replace the vanilla HTML portal with a Vite + React + TypeScript app branded **Forge**. Implement Catalog, Request, Status, and Settings behind a `ProvisioningPort` with Mock default and a GitHub Actions stub/adapter for later wiring.

## Technical Context

**Language/Version**: TypeScript 5.x, React 19

**Primary Dependencies**: Vite, React Router, no UI component library (custom CSS tokens)

**Storage**: localStorage (settings + mock requests)

**Testing**: Manual smoke + `npm run build`

**Target Platform**: Modern browsers (desktop + mobile)

**Project Type**: SPA in `frontend/`

**Performance Goals**: Instant mock UX; no blocking network in mock mode

**Constraints**: Constitution v1.0.0 — no new backend; XRD field fidelity; light teal design

**Scale/Scope**: 4 pages, 1 active catalog service, mock request list

## Constitution Check

- [x] Frontend-only + ProvisioningPort
- [x] Catalog / Request / Status / Settings only
- [x] XRD VM fields preserved
- [x] No purple-on-dark; cards only for interactive surfaces
- [x] Secrets only in localStorage

## Project Structure

### Documentation (this feature)

```text
specs/001-forge-portal/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/provisioning-port.md
└── tasks.md
```

### Source Code

```text
frontend/
  package.json
  vite.config.ts
  tsconfig.json
  index.html
  src/
    main.tsx
    App.tsx
    styles/tokens.css
    styles/global.css
    pages/CatalogPage.tsx
    pages/RequestPage.tsx
    pages/StatusPage.tsx
    pages/SettingsPage.tsx
    components/AppShell.tsx
    components/StatusBadge.tsx
    lib/provisioning/types.ts
    lib/provisioning/port.ts
    lib/provisioning/mock.ts
    lib/provisioning/github.ts
    lib/provisioning/index.ts
  README.md
```

## Implementation Approach

1. Scaffold Vite React-TS in `frontend/` (remove old `index.html` monolith).
2. Tokens + AppShell + routes.
3. ProvisioningPort + mock + settings context.
4. Catalog → Request → Status pages.
5. GitHub adapter stub (create via workflow_dispatch when mode=github).
6. Update root README; verify build.
