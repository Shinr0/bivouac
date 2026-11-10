# Platform Forge Constitution

Product brand for the developer portal: **Forge**.

## Core Principles

### I. BFF-Mediated Provisioning (NON-NEGOTIABLE)
The portal is a Vite + React + TypeScript UI that never talks to a cloud API or
to Kubernetes directly. All provisioning goes through the `ProvisioningPort`
interface, whose production implementation calls the **BFF** (`bff/`, Go). The
BFF is the only component holding cluster credentials, via its ServiceAccount.
A **Mock** adapter (localStorage) stays available for offline UI work. No YAML
is generated in the browser and no provisioning path goes through Git commits
or CI workflows.

### II. IDP Scope Lock
Screens: **Catalog**, **Request** (self-service form), **Status** (request
tracking), **Settings**. No SSO, no team directory, no Backstage-style docs hub,
no platform health dashboard until a later constitution amendment.

### III. XRD Contract Fidelity
The XRD is the single source of truth for the self-service contract. The BFF
serves the accepted values (`GET /api/v1/catalog`) and the frontend builds its
form from that response rather than hardcoding enums. Virtual Machine requests
expose `vmName`, `machineType`, `zone`, `imageFamily`, `diskSizeGb`,
`rootVolumeType`, `dataDiskSizeGb`, `staticIp`. Field names and defaults stay
aligned with `platform/apis/xrd-virtualmachine.yaml`.

### IV. Ownership Boundaries
Three layers, each with one owner, and no overlap:

| Layer | Owner | Contents |
|---|---|---|
| Day 0 infrastructure | OpenTofu (`infra-rcscale-scaleway`) | VPC, Kapsule, IAM, credentials Secret, ArgoCD |
| Platform configuration | ArgoCD / Git (`platform/`) | Crossplane, providers, XRD, Compositions, Forge manifests |
| User resources | BFF → Kubernetes API | Namespaced Composite Resources |

Never manage Crossplane CRs with `kubernetes_manifest` in OpenTofu: the plan
performs a dry-run against CRDs that do not exist yet.

### V. Small Diffs, No Mid-Flight Redesign
Prefer the smallest change that completes the current Spec Kit task. Visual
direction is fixed: light cool gray-blue atmosphere, teal/cyan accent,
expressive fonts (not Inter/Roboto/Arial), no purple-on-dark theme, cards only
for interactive surfaces.

### VI. Verify Before Done
Run `npm run build` and `npm run lint` in `frontend/`, `go vet ./...` and
`go test ./...` in `bff/`, and `tofu validate` in the infrastructure module
after meaningful work. Never commit credentials: the Scaleway API key lives only
in the Secret created by OpenTofu, and the portal stores no secret at all.

## Stack Constraints

- Cloud: **Scaleway** (Kapsule, Instances, Block Storage).
- Control plane: **Crossplane v2**, namespaced Composite Resources. Claims are
  deprecated and MUST NOT be reintroduced.
- Composition style: `mode: Pipeline` with `function-go-templating` and
  `function-auto-ready`.
- Frontend root: `frontend/`. Backend root: `bff/`. GitOps root: `platform/`.
- Language: French UI copy and comments; code/identifiers in English.

## Development Workflow

1. Spec Kit governs changes: specify → plan → tasks → implement.
2. Human reviews each Markdown artifact before the next phase.
3. A change to the XRD is a contract change: update the Composition, the BFF
   validation, and the frontend defaults in the same task.

## Governance

This constitution supersedes ad-hoc preferences. Amendments require a version
bump and dated note. Scope expansions (docs hub, teams, additional catalog
services) need an explicit amendment before implementation.

**Version**: 2.0.0 | **Ratified**: 2026-07-25 | **Last Amended**: 2026-08-02

### Amendment log

- **2.0.0** (2026-08-02) — Migration GCP → Scaleway. Principle I replaced:
  the GitHub Actions adapter and the browser-side YAML generation are removed
  in favour of a Go BFF calling the Kubernetes API. Crossplane v1 claims
  replaced by v2 namespaced composites. Added Principle IV (ownership
  boundaries) and extended verification to the backend and infrastructure.
- **1.0.0** (2026-07-25) — Initial ratification.
