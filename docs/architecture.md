# Architecture

Vue d'ensemble de la plateforme cible. Le détail de chaque choix est dans [`adr/`](adr/README.md),
la spécification complète dans [`PLAN.md`](../PLAN.md).

## Couches

| Couche | Cycle de vie | Contenu | Outil |
|---|---|---|---|
| Socle | Persistant (`make bootstrap`) | APIs, bucket de state, secret `github-token`, budget | Terraform, state local |
| Plateforme | Éphémère (`make up` / `make down`) | VPC, GKE zonal, node pools spot, IAM / Workload Identity, Argo CD | Terraform, state GCS |
| Services | Éphémère, réconcilié | Kyverno, ESO, Config Connector, kro, OpenCost, portail (LiteLLM en option) | Argo CD (GitOps) |
| Demandes | Éphémère, par PR | `TeamSpace` (puis `LLMAccess` en option) | Git → Argo CD → kro → Config Connector (GCP) + K8s |

## Vue d'ensemble

```
                        ┌──────────────────────── GitHub (dépôt public) ────────────────────────┐
                        │  gitops/  (apps Argo CD)     claims/  (demandes des équipes)          │
                        └───────────▲──────────────────────────────▲────────────────────────────┘
                                    │ sync                         │ PR ouverte par le portail/agent
┌───────────── GKE Standard zonal (europe-west4) ─────────────────┼────────────────────────────┐
│                                                                  │                            │
│  Argo CD (app-of-apps) ──► Kyverno ──► kro ──► Config Connector + K8s ──► ressources d'équipe │
│                                                                                               │
│  Portail (FastAPI + HTMX) ── catalogue, formulaires depuis schémas, chat assistant ───────────┘
│        │                                                                                      │
│        ▼                                                                                      │
│  Gateway LLM (LiteLLM, API OpenAI-compatible, clés virtuelles + budgets par équipe)           │
│        └──► Vertex AI (Claude via Model Garden)        [seul backend, facturé à l'usage]      │
│                                                                                               │
│  OpenCost (coûts par namespace/équipe)   Prometheus/Grafana (léger)                           │
│                                                                                               │
│  Pool : cpu-spot (e2-standard-4, spot) — aucun GPU                                            │
└───────────────────────────────────────────────────────────────────────────────────────────────┘
        ▲ Workload Identity (KSA ↔ GSA roles/aiplatform.user)        ▲ accès : port-forward / IAP
```

## Flux d'une demande

1. Une équipe (ou l'assistant IA, après confirmation) ouvre une **pull request** ajoutant un claim sous `claims/`.
2. La CI exécute les politiques **Kyverno** sur le claim ; un humain relit et merge.
3. **Argo CD** synchronise le claim ; Kyverno le revalide à l'admission.
4. **kro** crée, dans l'ordre, les ressources de l'équipe — namespace, quotas, RBAC — et les objets
   **Config Connector** qui deviennent un bucket GCS et un compte de service GCP lié par Workload
   Identity ; il expose leur statut.
5. Le portail affiche le statut et les coûts (OpenCost, LiteLLM).

L'IA n'a **aucun** droit d'écriture sur le cluster : sa seule action est d'ouvrir une PR.

## Invariants

- Aucun secret dans Git, aucune clé de compte de service ([ADR 0004](adr/0004-workload-identity-no-keys.md)).
- Ni PVC, ni LoadBalancer, ni IP statique ([ADR 0007](adr/0007-no-orphan-cloud-resources.md)).
- Seule exception : les ressources GCP des équipes, créées par Config Connector, rapides, quasi gratuites, supprimées avec la demande et avant le cluster au `make down` ([ADR 0009](adr/0009-kro-on-config-connector.md)).
- Label `env=bivouac` sur toute ressource GCP ; `make check-orphans` échoue s'il en reste après `down`.
