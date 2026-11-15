# Décisions d'architecture (ADR)

| N° | Décision | Statut |
|---|---|---|
| [0001](0001-record-architecture-decisions.md) | Consigner les décisions d'architecture | Accepté |
| [0002](0002-terraform-two-layers.md) | Terraform en deux couches : socle persistant et plateforme éphémère | Accepté |
| [0003](0003-bootstrap-local-state.md) | State local pour la couche bootstrap | Accepté |
| [0004](0004-workload-identity-no-keys.md) | Workload Identity partout, aucune clé de compte de service | Accepté |
| [0005](0005-gke-standard-zonal.md) | GKE Standard zonal, nœuds spot | Accepté (implémentation en phase 2) |
| [0006](0006-public-nodes-without-nat.md) | Nœuds publics sans Cloud NAT par défaut | Accepté (implémentation en phase 2) |
| [0007](0007-no-persistent-cloud-resources-from-k8s.md) | Ni PVC, ni LoadBalancer, ni IP statique dans la démo | Accepté |
| [0008](0008-argocd-app-of-apps.md) | Argo CD, pattern app-of-apps avec sync waves | Accepté (implémentation en phase 3) |
| [0009](0009-kro-developer-api.md) | kro comme API développeur | Accepté (implémentation en phases 4–5) |
| [0010](0010-kyverno-policies.md) | Kyverno pour les politiques, en CI et dans le cluster | Accepté (implémentation en phase 4) |
| [0011](0011-portal-fastapi-htmx.md) | Portail maison FastAPI + HTMX | Accepté (implémentation en phases 6–7) |
| [0012](0012-litellm-gateway.md) | LiteLLM comme gateway LLM | Accepté — reporté en phase optionnelle 7 (ADR 0018) |
| [0013](0013-vertex-ai-claude.md) | Vertex AI (Claude via Model Garden) comme backend LLM | Accepté — reporté en phase optionnelle 7 (ADR 0018) |
| [0014](0014-kserve-vllm.md) | KServe + vLLM pour les modèles auto-hébergés | Remplacé par [0017](0017-vertex-only-no-gpu.md) (2026-11-15) |
| [0015](0015-finops-guardrails.md) | Garde-fous FinOps : budget GCP, OpenCost, budgets LiteLLM | Accepté (budget GCP en phase 1, le reste en phases 3–5) |
| [0016](0016-ghcr-images.md) | Images sur GHCR, construites par GitHub Actions | Accepté (implémentation en phase 6) |
| [0017](0017-vertex-only-no-gpu.md) | Vertex AI uniquement : ni GPU ni modèle auto-hébergé | Accepté — remplace [0014](0014-kserve-vllm.md) |
| [0018](0018-defer-in-platform-ai.md) | Reporter l'IA intégrée à la plateforme | Accepté |
| [0019](0019-ai-assisted-workflow.md) | Développer la plateforme avec Claude Code, sous garde-fous | Accepté |

Nouvelle décision : copier [`0000-template.md`](0000-template.md).
