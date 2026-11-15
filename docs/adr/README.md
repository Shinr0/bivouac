# Décisions d'architecture (ADR)

Règles : un ADR n'est jamais modifié sur le fond ni supprimé ; une décision révisée est remplacée par un nouvel ADR ([0001](0001-record-architecture-decisions.md)).

| N° | Décision | Statut |
|---|---|---|
| [0001](0001-record-architecture-decisions.md) | Consigner les décisions d'architecture | Accepté |
| [0002](0002-terraform-two-layers.md) | Terraform en deux couches : socle persistant et plateforme éphémère | Accepté |
| [0003](0003-bootstrap-local-state.md) | State local pour la couche bootstrap | Accepté |
| [0004](0004-workload-identity-no-keys.md) | Workload Identity et OIDC partout, aucune clé stockée | Accepté |
| [0005](0005-gke-standard-zonal-spot.md) | GKE Standard zonal, nœuds spot, sans GPU | Accepté |
| [0006](0006-public-nodes-without-nat.md) | Nœuds publics sans Cloud NAT par défaut | Accepté |
| [0007](0007-no-orphan-cloud-resources.md) | Aucune ressource cloud orpheline : règles et exception encadrée | Accepté |
| [0008](0008-argocd-app-of-apps.md) | Argo CD, pattern app-of-apps avec sync waves | Accepté |
| [0009](0009-kro-on-config-connector.md) | API de plateforme : kro au-dessus de Config Connector | Accepté (implémentation en phases 3–4) |
| [0010](0010-kyverno-policies.md) | Kyverno pour les politiques, en CI et dans le cluster | Accepté (implémentation en phase 4) |
| [0011](0011-portal-fastapi-htmx.md) | Portail maison FastAPI + HTMX | Accepté (implémentation en phase 5) |
| [0012](0012-finops-guardrails.md) | Garde-fous FinOps : budget GCP, OpenCost, coûts bornés par construction | Accepté (budget en phase 1, OpenCost en phase 3) |
| [0013](0013-ghcr-images.md) | Images sur GHCR, construites par GitHub Actions | Accepté (implémentation en phase 5) |
| [0014](0014-ai-assisted-workflow.md) | Développer avec Claude Code, sous garde-fous | Accepté |
| [0015](0015-optional-in-platform-ai.md) | IA intégrée à la plateforme : optionnelle, LiteLLM + Vertex AI | Accepté (phases optionnelles 7–8) |

Nouvelle décision : skill `new-adr` ou copier [`0000-template.md`](0000-template.md).
