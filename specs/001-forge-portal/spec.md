# Feature Specification: Forge Developer Portal

**Feature Branch**: `001-forge-portal`

**Created**: 2026-07-25

**Status**: Active

**Input**: Rewrite developer portal as Forge IDP — Catalog, Request, Status, Settings with Vite React and mock ProvisioningPort

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Browse catalog and request a VM (Priority: P1)

A platform user opens Forge, sees the service catalog with Virtual Machine as the primary offer, chooses it, fills the XRD-aligned form, and submits a request. In mock mode the request is stored locally and appears as `submitted`.

**Why this priority**: Core IDP self-service loop; without it the portal has no value.

**Independent Test**: From a clean browser profile, open Catalog → Request VM → submit → see success feedback and a new entry on Status.

**Acceptance Scenarios**:

1. **Given** the app is running, **When** the user opens `/`, **Then** they see the Forge brand and a catalog listing Virtual Machine as requestable.
2. **Given** the user selects Virtual Machine, **When** they submit valid `vmName`, `machineType`, `zone`, `diskSizeGb`, `imageFamily`, **Then** a request is created with status `submitted`.
3. **Given** invalid or empty `vmName`, **When** they submit, **Then** the form blocks submit with a clear validation message.

---

### User Story 2 - Track request status (Priority: P2)

A user opens Status and sees their requests with mock lifecycle statuses (`submitted`, `provisioning`, `ready`, `failed`). They can tell which requests progressed.

**Why this priority**: Closes the self-service loop after create.

**Independent Test**: After creating at least one mock request, Status lists it with a visible status badge; advancing mock status (if provided) updates the badge.

**Acceptance Scenarios**:

1. **Given** one or more mock requests exist, **When** the user opens Status, **Then** each request shows name, key parameters, and status.
2. **Given** no requests exist, **When** the user opens Status, **Then** an empty state invites them to use the catalog.

---

### User Story 3 - Configure settings and provider mode (Priority: P3)

A user opens Settings, saves GitHub token + repo (localStorage only), and chooses provider mode `mock` (default) or `github` (adapter present). Mock mode works without credentials.

**Why this priority**: Needed for later real provisioning; not required for demoing Catalog/Request/Status with mocks.

**Acceptance Scenarios**:

1. **Given** Settings open, **When** the user saves token and repo, **Then** values persist across reload (localStorage).
2. **Given** mode is `mock`, **When** the user creates a VM request, **Then** no GitHub API call is made.

### Edge Cases

- Duplicate `vmName` in mock store: reject or show error.
- Missing Settings in `github` mode: block create and prompt Settings.
- Mobile viewport: nav and form remain usable.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Portal brand is **Forge**, visible as primary hero signal on Catalog.
- **FR-002**: Screens: Catalog, Request, Status, Settings (client routing).
- **FR-003**: Catalog includes Virtual Machine (active) and at least one “coming soon” service.
- **FR-004**: Request form fields include XRD: `vmName`, `machineType`, `zone`, `diskSizeGb`, `imageFamily`; optional `spot`, `staticIp`, `dataDiskSizeGb`.
- **FR-005**: All provisioning goes through `ProvisioningPort`.
- **FR-006**: Default port implementation is Mock (localStorage).
- **FR-007**: Status shows requests with statuses `submitted` | `provisioning` | `ready` | `failed`.
- **FR-008**: Settings persist `token`, `repo`, `mode` in localStorage; never commit secrets.
- **FR-009**: Replace legacy `frontend/index.html` monolith with Vite React TS app under `frontend/`.

### Key Entities

- **CatalogItem**: id, title, description, available (bool)
- **VmRequest**: id, parameters (XRD fields), status, createdAt
- **PortalSettings**: token, repo, mode (`mock` | `github`)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: User can complete Catalog → Request → Status in under 2 minutes in mock mode.
- **SC-002**: `npm run build` succeeds with no TypeScript errors.
- **SC-003**: No purple-on-dark theme; teal accent on light atmospheric background.
- **SC-004**: README documents `cd frontend && npm install && npm run dev`.
