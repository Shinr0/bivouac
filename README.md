# Bivouac

**Le camp de base pour explorer le platform engineering sur GCP : monté en 20 minutes, démonté
sans laisser de trace.**

Bivouac est un **laboratoire** pour rechercher et tester des concepts de platform engineering sur
Google Cloud — API de plateforme, GitOps, politiques, FinOps, self-service — sur une vraie
plateforme interne de développement (IDP), sans payer d'infrastructure entre deux essais.

- `make up` construit toute la plateforme depuis un projet GCP vide ;
- on expérimente : une API de plateforme (kro + Config Connector), une politique Kyverno, un
  nouveau composant ;
- `make down` détruit tout, et un contrôle échoue s'il reste la moindre ressource facturée.

- **Construite avec un agent IA, sous contrôle** : Claude Code itère sur le dépôt, encadré par
  des règles (`CLAUDE.md`), des garde-fous (permissions, hooks), un contexte économe
  et des revues IA à coût borné — voir PLAN.md § 16 et [ADR 0014](docs/adr/0014-ai-assisted-workflow.md).
- **IA dans la plateforme (phases optionnelles)** : accès LLM par équipe via LiteLLM + Vertex AI,
  assistant qui ouvre des pull requests — en option ([ADR 0015](docs/adr/0015-optional-in-platform-ai.md)).

> Projet en construction, phase par phase. Spécification complète : [`PLAN.md`](PLAN.md).

## Architecture

```
GitHub (gitops/, claims/) ──sync──► Argo CD ──► Kyverno ──► kro ──► Config Connector + K8s ──► ressources par équipe
                                        │
GKE Standard zonal (spot)               ├─ Portail FastAPI + HTMX (catalogue, demandes)
                                        ├─ (optionnel) LiteLLM ──► Vertex AI (Claude)
                                        └─ OpenCost, Prometheus/Grafana
```

Détails : [`docs/architecture.md`](docs/architecture.md) · Décisions : [`docs/adr/`](docs/adr/README.md)

| Couche | Outil |
|---|---|
| IaC | Terraform (compatible OpenTofu), deux couches : socle persistant / plateforme éphémère |
| Cluster | GKE Standard zonal, nœuds spot, aucun GPU |
| GitOps | Argo CD, app-of-apps |
| API développeur | kro au-dessus de Config Connector (objets GCP et Kubernetes en self-service) |
| Politiques | Kyverno (CI + admission) |
| Portail | FastAPI + HTMX |
| Gateway LLM *(optionnel)* | LiteLLM → Vertex AI (Claude), pas de modèle auto-hébergé |
| Outillage IA | Claude Code, skills de projet, revue `claude-code-action` (OIDC, sans clé) |
| FinOps | Budget GCP, OpenCost (+ budgets LiteLLM si phase 7) |

## Avancement

| Phase | Contenu | État |
|---|---|---|
| 1 | Squelette, ADR, Makefile, `infra/bootstrap` | ✅ |
| 2 | `infra/platform` (VPC, GKE, IAM, WI) + `up` / `down` / `check-orphans` | ⏳ |
| 3 | Argo CD app-of-apps, Kyverno, ESO, Config Connector, kro, OpenCost | — |
| 4 | API `TeamSpace` (namespace, quotas, bucket GCS, compte de service WI), politiques, CI Kyverno | — |
| 5 | Portail : catalogue, formulaires, PR | — |
| 6 | Chemin local kind, polish | — |
| 7 *(opt.)* | LiteLLM + Vertex, API `LLMAccess` | optionnel |
| 8 *(opt.)* | Assistant IA dans le portail | optionnel |
| transverse | Workflow Claude Code : contexte économe, garde-fous, revues IA | ✅ |

## Prérequis

### Outils

| Outil | Version testée |
|---|---|
| Terraform (ou OpenTofu) | 1.17 (≥ 1.9 requis) |
| Google Cloud SDK (`gcloud`) | 590 |
| GNU Make, Bash | 4.x / 5.x |

`kubectl`, `helm`, la CLI Kyverno et Python seront requis à partir des phases suivantes.

### Actions manuelles (une seule fois)

Ce sont les **seules** opérations hors code (PLAN.md § 10).

1. **Projet GCP avec facturation et identifiants isolés.** Le projet utilise sa propre
   configuration gcloud (`~/.config/gcloud-bivouac`, ADC comprises) : la configuration
   gcloud globale de la machine n'est jamais lue ni modifiée. Renseigner `GCP_ACCOUNT` et
   `GCP_PROJECT_ID` dans `.env` (étape 3), puis :
   ```bash
   make gcloud-login   # ouvre le navigateur : connexion CLI + ADC sur GCP_ACCOUNT
   make whoami         # vérifie compte (CLI et ADC) et projet
   ```
   Tous les scripts refusent de s'exécuter si le compte actif n'est pas `GCP_ACCOUNT`.
   Le compte doit pouvoir activer des APIs sur le projet et créer un budget sur le compte de
   facturation (`roles/billing.costsManager` ou administrateur de facturation).
2. *(Phases optionnelles 7–8 uniquement)* **Activer les modèles Claude dans Vertex AI Model
   Garden** et accepter les conditions (console GCP › Vertex AI › Model Garden).
3. **Configuration locale** :
   ```bash
   cp .env.example .env   # GCP_ACCOUNT, GCP_PROJECT_ID, GCP_BILLING_ACCOUNT_ID
   ```
4. **Socle persistant** :
   ```bash
   make bootstrap
   ```
   Crée : activation des APIs, bucket de state `<PROJECT_ID>-tfstate` (versionné), secret vide
   `github-token`, budget mensuel (50 € par défaut, alertes à 50 / 90 / 100 %). Le state de cette
   couche reste local — voir [ADR 0003](docs/adr/0003-bootstrap-local-state.md).
5. **Token GitHub** : créer un token *fine-grained* limité à ce dépôt (permissions
   *Contents* et *Pull requests* en lecture/écriture), puis l'ajouter au secret sans qu'il
   passe par l'historique du shell :
   ```bash
   export CLOUDSDK_CONFIG=~/.config/gcloud-bivouac   # config isolée
   read -rs TOKEN && printf '%s' "$TOKEN" | \
     gcloud secrets versions add github-token --data-file=- && unset TOKEN
   ```

## Utilisation

```bash
make help        # liste des cibles
make whoami      # compte et projet GCP utilisés
make test        # contrôles statiques (terraform fmt/validate, puis lint, kyverno, pytest)
make bootstrap   # socle persistant, idempotent
make up          # plateforme éphémère          (phase 2)
make down        # destruction + check-orphans  (phase 2)
```

## Travailler sur le dépôt avec Claude Code

Le dépôt est fait pour être modifié avec [Claude Code](https://claude.com/claude-code), sous
garde-fous versionnés (PLAN.md § 16, [ADR 0014](docs/adr/0014-ai-assisted-workflow.md)) :

| Élément | Rôle |
|---|---|
| `CLAUDE.md`, `PLAN.md` | Règles et phases : l'agent ne travaille que sur la phase en cours |
| `.claude/settings.json` | Confirmation exigée pour `terraform apply/destroy`, `make up/down/bootstrap`, `git push` |
| `.claude/hooks/guard-gcloud.sh` | Bloque tout `gcloud` hors de la config isolée du projet (testé par `make test`) |
| `.claude/skills/` | Routines chargées à la demande : `new-adr`, `close-phase` |
| `permissions.deny` | Interdit la lecture de `.terraform/`, states, lock files, `*.tfvars`, `.env` |
| `/code-review` | Revue locale avant push, sans coût de CI |
| `.github/workflows/claude-review.yml` | Revue IA des PR, bornée en coût, authentifiée par OIDC (aucune clé) |

Activer la revue IA en CI (une fois) : dans la console Claude, *Settings › Workload identity*,
enregistrer l'émetteur `https://token.actions.githubusercontent.com`, créer un compte de service et
une règle de fédération limitée à `repo:<owner>/platform-project:`, fixer une limite de dépense sur
le workspace, puis créer les variables de dépôt `ANTHROPIC_FEDERATION_RULE_ID`,
`ANTHROPIC_ORGANIZATION_ID` et `ANTHROPIC_SERVICE_ACCOUNT_ID`. Sans elles, le job est ignoré.

## Coûts

| Élément | Ordre de grandeur |
|---|---|
| Socle (bucket, secret, budget) | ≈ 0 € / mois |
| Nœud e2-standard-4 spot (europe-west4) | ≈ 0,056 $/h |
| Frais de gestion GKE (cluster zonal) | 0,10 $/h, couverts par le crédit gratuit GKE |
| Vertex AI (Claude), phases optionnelles | au token, plafonné par les budgets LiteLLM |
| Démo de 2 h (infra) | < 1 $ |

Prix relevés en novembre 2026, à revérifier dans le calculateur GCP. Garde-fous : budget GCP,
nœuds spot, aucun GPU, interdiction des PVC/LoadBalancer, contrôle d'orphelins
après chaque `down` ([ADR 0012](docs/adr/0012-finops-guardrails.md)).

## Structure du dépôt

```
├── PLAN.md               spécification et phases
├── CLAUDE.md             conventions pour Claude Code
├── Makefile
├── docs/                 architecture, ADR
├── infra/bootstrap/      socle persistant (Terraform, state local)
├── scripts/              scripts appelés par le Makefile
├── .claude/              garde-fous et skills pour Claude Code
└── .github/workflows/    CI
```

La structure cible complète est décrite dans PLAN.md § 14. L'ancienne version de ce dépôt
(Crossplane + portail React) est conservée dans la branche `archive/v1-crossplane-forge`.
