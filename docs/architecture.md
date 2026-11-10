# Architecture

Vue d'ensemble de la plateforme cible. Le détail de chaque choix est dans [`adr/`](adr/README.md),
la spécification complète dans [`PLAN.md`](../PLAN.md).

## Couches

| Couche | Cycle de vie | Contenu | Outil |
|---|---|---|---|
| Socle | Persistant (`make bootstrap`) | APIs, bucket de state, secret `github-token`, budget | Terraform, state local |
| Plateforme | Éphémère (`make up` / `make down`) | VPC, GKE zonal, node pools spot, IAM / Workload Identity, Argo CD | Terraform, state GCS |
| Services | Éphémère, réconcilié | Kyverno, ESO, kro, OpenCost, LiteLLM, KServe, portail | Argo CD (GitOps) |
| Demandes | Éphémère, par PR | Claims `TeamSpace`, `LLMAccess` | Git → Argo CD → kro |

## Vue d'ensemble

```
                        ┌──────────────────────── GitHub (dépôt public) ────────────────────────┐
                        │  gitops/  (apps Argo CD)     claims/  (demandes des équipes)          │
                        └───────────▲──────────────────────────────▲────────────────────────────┘
                                    │ sync                         │ PR ouverte par le portail/agent
┌───────────── GKE Standard zonal (europe-west4) ─────────────────┼────────────────────────────┐
│                                                                  │                            │
│  Argo CD (app-of-apps) ──► Kyverno ──► kro (RGD) ──► ressources par équipe                    │
│                                                                                               │
│  Portail (FastAPI + HTMX) ── catalogue, formulaires depuis schémas, chat assistant ───────────┘
│        │                                                                                      │
│        ▼                                                                                      │
│  Gateway LLM (LiteLLM, API OpenAI-compatible, clés virtuelles + budgets par équipe)           │
│        ├──► Vertex AI (Claude via Model Garden)        [par défaut, facturé à l'usage]        │
│        └──► KServe + vLLM sur nœud GPU L4 spot         [optionnel, à la demande]              │
│                                                                                               │
│  OpenCost (coûts par namespace/équipe)   Prometheus/Grafana (léger)                           │
│                                                                                               │
│  Pools : cpu-spot (2× e2-standard-4)    gpu-spot (g2-standard-4, L4, min=0)                   │
└───────────────────────────────────────────────────────────────────────────────────────────────┘
        ▲ Workload Identity (KSA ↔ GSA roles/aiplatform.user)        ▲ accès : port-forward / IAP
```

## Flux d'une demande

1. Une équipe (ou l'assistant IA, après confirmation) ouvre une **pull request** ajoutant un claim sous `claims/`.
2. La CI exécute les politiques **Kyverno** sur le claim ; un humain relit et merge.
3. **Argo CD** synchronise le claim ; Kyverno le revalide à l'admission.
4. **kro** génère les ressources (namespace, quotas, accès LLM…) et expose un statut.
5. Le portail affiche le statut et les coûts (OpenCost, LiteLLM).

L'IA n'a **aucun** droit d'écriture sur le cluster : sa seule action est d'ouvrir une PR.

## Invariants

- Aucun secret dans Git, aucune clé de compte de service ([ADR 0004](adr/0004-workload-identity-no-keys.md)).
- Aucune ressource cloud créée par un contrôleur Kubernetes : ni PVC, ni LoadBalancer, ni IP statique ([ADR 0007](adr/0007-no-persistent-cloud-resources-from-k8s.md)).
- Label `env=platform-ai-demo` sur toute ressource GCP ; `make check-orphans` échoue s'il en reste après `down`.
