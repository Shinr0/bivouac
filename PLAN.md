# Plateforme IA éphémère sur GCP — Spécification de déploiement

> Document d'entrée pour Claude Code. Il décrit **quoi** construire et **dans quel ordre**.
> Les choix d'outils sont arrêtés ; les versions exactes, la syntaxe des APIs et des charts
> doivent être **vérifiées dans la documentation officielle au moment de l'implémentation**
> (voir § 12 « Points à vérifier »). Ne rien inventer : en cas de doute, lire la doc ou demander.
>
> **Révision du 2026-11-15** : pas de GPU ni de modèle auto-hébergé ; les LLM passent uniquement
> par Vertex AI ([ADR 0017](docs/adr/0017-vertex-only-no-gpu.md)). L'ancien volet KServe/vLLM est retiré.
>
> **Révision du 2026-11-15 (bis)** : l'IA **intégrée à la plateforme** (LiteLLM, Vertex AI, RGD `LLMAccess`,
> assistant du portail) est **reportée** en phases optionnelles 7 et 8
> ([ADR 0018](docs/adr/0018-defer-in-platform-ai.md)). L'IA intervient d'abord **dans la façon de construire**
> la plateforme : Claude Code, contexte économe, garde-fous et revues
> ([ADR 0019](docs/adr/0019-ai-assisted-workflow.md), § 16). Les sections qui décrivent l'IA intégrée
> restent valables pour ces phases optionnelles.

---

## 1. Contexte et objectif

Vitrine technique d'un consultant freelance Cloud DevOps (GCP, Kubernetes, Terraform,
Platform Engineering) : une **plateforme interne de développement (IDP) en self-service, avec de
l'IA par-dessus**, publiée en dépôt GitHub public et démontrable en live devant un prospect.

Deux volets IA, **reportés en phases optionnelles 7–8** (ADR 0018) :
1. **LLM-as-a-Service** : une équipe demande un accès LLM (Claude sur Vertex AI) par un
   simple claim Git ; la plateforme provisionne, applique quotas et budgets, mesure les coûts.
2. **Assistant de plateforme** : un agent en langage naturel, intégré au portail, qui aide à
   formuler une demande et **ouvre une pull request** — il n'applique jamais rien directement.

### Exigence principale : éphémère de A à Z

- `make up` crée **toute** l'infrastructure depuis un projet GCP vide (hors socle persistant, § 4).
- `make down` détruit **tout**, sans ressource orpheline facturée.
- L'infra ne vit que le temps d'une démo → coût proche de zéro hors démo.
- Cible : `make up` en moins de ~20 min (à mesurer), `make down` en moins de ~15 min.
- Idempotence : relancer `make up` ou `make down` ne casse rien.

---

## 2. Principes non négociables

1. **Tout est code** : Terraform pour GCP, GitOps (Argo CD) pour tout ce qui tourne dans le cluster.
   Aucune action manuelle hors les prérequis listés au § 10.
2. **Zéro secret dans Git, zéro clé JSON de compte de service.** Workload Identity partout.
   Les seuls secrets (token GitHub) vivent dans Secret Manager (socle persistant).
3. **Rien ne survit à `make down` en dehors du socle persistant.** Toutes les ressources GCP
   créées portent le label `env=platform-ai-demo` pour permettre un contrôle d'orphelins.
4. **L'IA ne touche jamais le cluster directement** : elle produit des PR, la validation passe
   par Kyverno (en CI et dans le cluster) + revue humaine.
5. **Coûts bornés par construction** : spot partout, aucun GPU, budgets LLM plafonnés.
6. **Lisible pour un recruteur** : README clair, schéma d'architecture, ADR pour chaque choix.

---

## 3. Architecture cible

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
│        └──► Vertex AI (Claude via Model Garden)        [seul backend, facturé à l'usage]      │
│                                                                                               │
│  OpenCost (coûts par namespace/équipe)   Prometheus/Grafana (léger)                           │
│                                                                                               │
│  Pool : cpu-spot (e2-standard-4, spot) — aucun GPU                                            │
└───────────────────────────────────────────────────────────────────────────────────────────────┘
        ▲ Workload Identity (KSA ↔ GSA roles/aiplatform.user)        ▲ accès : port-forward / IAP
```

### Choix arrêtés (justifier chacun dans un ADR `docs/adr/`)

| Couche | Choix | Pourquoi |
|---|---|---|
| IaC socle GCP | Terraform (ou OpenTofu, compatible) | Standard du marché, certification de l'auteur |
| Cluster | GKE **Standard zonal** | Le crédit gratuit GKE couvre les frais de gestion d'un cluster zonal ; contrôle des node pools (spot) |
| GitOps | Argo CD, pattern app-of-apps | Maîtrisé par l'auteur, visuel en démo |
| API développeur | **kro** (ResourceGraphDefinitions) ; Crossplane seulement si une ressource GCP doit être exposée en self-service (optionnel, phase 5) | kro est léger et natif K8s ; Crossplane reste pertinent pour le cloud |
| Politiques | Kyverno | Plus simple qu'OPA pour du 100 % Kubernetes ; utilisable en CI (CLI) |
| Portail | **Maison** : FastAPI + HTMX (pas de SPA) | Backstage jugé trop lourd ; le portail n'est qu'un client de l'API (les CRD) — remplaçable |
| Gateway LLM | LiteLLM proxy | API OpenAI-compatible, clés virtuelles, budgets, multi-backends |
| Modèles | Vertex AI (Claude via Model Garden), **seul backend** | Pas de GPU à maintenir, facturé au token |
| FinOps | OpenCost + budget GCP + budgets LiteLLM | Argument FinOps de l'offre |
| Images | GHCR, buildées par GitHub Actions | Gratuit, hors périmètre détruit |

**Explicitement hors périmètre** : GPU et modèles auto-hébergés (KServe, vLLM, llm-d), Backstage, multi-cluster,
Kratix, service mesh.

---

## 4. Deux couches Terraform

### 4.1 `infra/bootstrap/` — socle **persistant** (créé une fois, quasi gratuit)

Ne jamais détruit par `make down`. Contient :
- Bucket GCS du state Terraform de la couche éphémère (versioning activé).
- Activation des APIs : container, compute, aiplatform, secretmanager, iam, billingbudgets…
- Secret Manager : secret `github-token` (la **valeur** est ajoutée à la main, § 10).
- Budget GCP sur le projet avec alertes (50 %, 90 %, 100 %) — montant en variable (défaut 50 €).
- Optionnel : un compte de service « deployer » utilisé par la CI.

State de cette couche : local ou bucket dédié créé manuellement — documenter le choix.

### 4.2 `infra/platform/` — couche **éphémère** (`make up` / `make down`)

- VPC dédié + sous-réseau (ranges secondaires pods/services), règles de pare-feu minimales.
- Cluster GKE Standard **zonal**, `deletion_protection = false`, Workload Identity activée,
  release channel `REGULAR`, labels `env=platform-ai-demo`.
- Node pool `cpu-spot` : e2-standard-4, spot, autoscaling 1–3 (à ajuster après mesure).
- Aucun node pool GPU.
- Compte de service Google `litellm-vertex` avec `roles/aiplatform.user` + binding Workload Identity
  vers le KSA `llm-gateway/litellm`.
- Compte de service Google pour External Secrets (lecture Secret Manager) + binding WI.
- Bootstrap Argo CD : installation du chart Argo CD via le provider Helm **uniquement**, puis
  création de l'Application racine (app-of-apps) pointant vers `gitops/`. Tout le reste est synchronisé par Argo CD.
- **Pas** de load balancer, d'IP statique ni de disque persistant créés par Terraform (sources d'orphelins).

Choix réseau par défaut : nœuds avec IP publique + pare-feu strict, **sans Cloud NAT**
(économie ~33 $/mois). Variable `private_nodes = true` pour basculer en nœuds privés + Cloud NAT.

---

## 5. Contenu GitOps (`gitops/`)

App-of-apps avec **sync waves** pour l'ordre :

| Wave | Application | Notes |
|---|---|---|
| 0 | namespaces, CRDs de base | |
| 1 | Kyverno + politiques | Politiques en mode `Enforce` |
| 1 | External Secrets Operator + `ClusterSecretStore` GCP | via Workload Identity |
| 2 | kro + ResourceGraphDefinitions | |
| 2 | OpenCost, kube-prometheus-stack (allégé) | rétention courte, pas de PVC (emptyDir) |
| 3 | *(phase 7, optionnelle)* LiteLLM (+ Postgres pour clés/budgets) | Postgres **sans PVC** en démo (emptyDir), assumé et documenté |
| 4 | Portail | image GHCR |
| 5 | `claims/` | instances kro créées par les équipes (via PR) |

Règle : **aucun `PersistentVolumeClaim`** dans la démo (un PVC = un disque PD qui peut survivre au
`destroy`). Une politique Kyverno l'interdit, sauf dérogation explicite.

---

## 6. API de la plateforme (kro)

À définir sous `platform/rgd/`. Deux RGD minimum :

### `TeamSpace`
Crée pour une équipe : Namespace, ResourceQuota, LimitRange, NetworkPolicy par défaut,
RoleBinding, labels de coût (`team=<nom>`).
Champs : `team` (requis), `owner` (email), `cpuQuota`, `memoryQuota`.

### `LLMAccess` *(phase 7, optionnelle)*
Donne à une équipe un accès au gateway :
- `model` : enum restreint (alias LiteLLM pointant vers des modèles Claude sur Vertex AI)
- `monthlyBudgetUsd` : entier, borné (ex. max 50)
- La création de la clé virtuelle LiteLLM se fait par un petit **Job/contrôleur** (ou un hook
  du portail) qui appelle l'API admin LiteLLM et stocke la clé dans un Secret du namespace de l'équipe.
  Choisir l'option la plus simple et la documenter dans un ADR.

Exemples de claims dans `claims/examples/`.

### Politiques Kyverno (`platform/policies/`)
- Interdiction des PVC (démo).
- Labels obligatoires `team`, `owner` sur les claims.
- `monthlyBudgetUsd` ≤ plafond global.
- Interdiction des Services `type: LoadBalancer` hors namespace autorisé.
- Interdiction de toute demande de GPU (`nvidia.com/gpu`) dans les pods.
- Les mêmes politiques tournent en CI avec la CLI Kyverno sur les PR de `claims/`.

---

## 7. Portail et assistant IA (`portal/`) — assistant et coûts LLM en phase 8, optionnelle

Python 3.12+, FastAPI, Jinja2 + HTMX, aucun framework JS lourd.

Pages :
1. **Catalogue** : liste des RGD disponibles (lues depuis l'API Kubernetes) et des claims existants
   par équipe, avec leur statut kro.
2. **Nouvelle demande** : formulaire **généré depuis le schéma OpenAPI** de la CRD → produit le YAML du claim
   → ouvre une **PR GitHub** (token depuis Secret Manager via External Secrets).
3. **Assistant** : chat qui appelle le gateway LiteLLM (backend Vertex). L'agent dispose d'outils
   (function calling) **en lecture seule + une écriture contrôlée** :
   - `list_capabilities()` — RGD et leurs schémas
   - `get_team_status(team)` — claims, état, coûts OpenCost
   - `validate_claim(yaml)` — validation schéma + `kyverno apply` local
   - `open_pull_request(yaml, title, description)` — **seule action d'écriture**, après confirmation
     explicite de l'utilisateur dans l'UI
4. **Coûts** : coût par équipe (OpenCost) + consommation LLM (API spend LiteLLM).

Contraintes :
- Le portail n'a **aucun** droit d'écriture Kubernetes (RBAC lecture seule sur ses ressources).
- Authentification : en mode privé (défaut), accès via `kubectl port-forward` ; en mode public,
  IAP ou OAuth2-proxy — à documenter, désactivé par défaut.
- Tests `pytest` : génération de formulaire depuis un schéma, rendu YAML, outils de l'agent
  (avec client LLM mocké), refus d'écriture sans confirmation.

---

## 8. Séquences `make up` / `make down`

### Makefile (cibles minimales)

```
make bootstrap   # une seule fois : couche persistante
make up          # couche éphémère + attente que tout soit Healthy
make demo        # port-forwards portail / Argo CD / Grafana, affiche les URLs
make down        # destruction complète + contrôle d'orphelins
make check-orphans
make local-up    # chemin kind, sans GCP (§ 9)
make test        # pytest + lint YAML + kyverno CLI + terraform validate/fmt
```

### `make up` (scripts/up.sh)
1. Vérifier les prérequis (gcloud authentifié sur le bon compte, projet, outils).
2. `terraform -chdir=infra/platform init` (backend GCS) puis `apply -auto-approve`.
3. Récupérer les credentials du cluster.
4. Attendre que l'Application racine Argo CD et toutes ses enfants soient `Synced/Healthy`
   (timeout explicite, message clair en cas d'échec).
5. Smoke tests : le portail répond ; si les phases IA sont activées, le gateway répond et un appel Vertex de test aboutit.
6. Afficher la durée totale.

### `make down` (scripts/down.sh) — l'ordre compte
1. Supprimer les claims puis l'Application racine Argo CD **en cascade** et attendre la
   suppression effective (éviter les finalizers bloquants : documenter le contournement).
2. Supprimer explicitement tout `Service type=LoadBalancer`, `Gateway`/`Ingress` et PVC restants,
   puis attendre que GCP ait libéré les ressources associées (forwarding rules, NEG, disques).
3. `terraform -chdir=infra/platform destroy -auto-approve`.
4. `make check-orphans` : lister via `gcloud` les ressources résiduelles facturables
   (instances, disques, adresses, forwarding-rules, NEG, target-pools, backend-services, routeurs NAT)
   filtrées sur le projet / label. **Code de sortie ≠ 0 s'il en reste.**
5. Afficher la durée totale.

---

## 9. Chemin local (kind) — pour un visiteur sans GCP

`make local-up` : cluster kind, même app-of-apps avec un overlay `local` :
- backend LLM = mode « mock » ou clé API fournie par l'utilisateur
  via variable d'environnement (jamais committée) ;
- External Secrets remplacé par un Secret créé localement depuis `.env`.
Objectif : démontrer kro, Kyverno, le portail et le flux PR sans compte cloud.

---

## 10. Prérequis manuels (une seule fois, à documenter dans le README)

1. Projet GCP avec facturation (compte de la société).
2. *(Phases 7–8 uniquement)* **Activer les modèles Claude dans Vertex AI Model Garden** et accepter
   les conditions (par projet, non automatisable proprement).
3. Créer un token GitHub à portée minimale (fine-grained, accès au seul dépôt, PR + contents)
   et le stocker dans le secret `github-token` de Secret Manager.
4. `make bootstrap`.

---

## 11. Garde-fous de coûts

| Mécanisme | Où |
|---|---|
| Budget GCP + alertes | `infra/bootstrap` |
| Nœuds spot, cluster zonal, pas de NAT par défaut | `infra/platform` |
| Aucun GPU (ni pool, ni pod) | `infra/platform` + Kyverno |
| Budgets par clé LiteLLM + plafond global | config LiteLLM / RGD `LLMAccess` |
| Interdiction PVC / LoadBalancer non autorisés | Kyverno |
| `check-orphans` obligatoire en fin de `down` | scripts |

Ordres de grandeur relevés en novembre 2026 (prix tiers, **à revérifier** dans le calculateur GCP) :
- e2-standard-4 spot à europe-west4 ≈ 0,056 $/h.
- Frais de gestion GKE 0,10 $/h, couverts par le crédit gratuit (74,40 $/mois) pour **un** cluster zonal.
- Une démo de 2 h ≈ moins d'un dollar d'infra + la consommation Vertex AI (au token, plafonnée par LiteLLM).

---

## 12. Points à vérifier avant d'implémenter (ne pas deviner)

- **Versions** à épingler (dernières stables) : providers Terraform google/helm/kubernetes, charts
  Argo CD, Kyverno, kro, External Secrets, OpenCost, LiteLLM. Épingler partout, pas de `latest`.
- **kro** : `apiVersion` exacte de `ResourceGraphDefinition`, syntaxe du simple schema
  (marqueurs `required`, `default`, `enum`, `minimum`/`maximum`), `includeWhen` — doc kro.run.
- **LiteLLM** : syntaxe `model_list` pour Claude sur Vertex (`vertex_ai/<model>`, `vertex_project`,
  `vertex_location`), endpoint `global`/`eu`, auth par ADC/Workload Identity, API admin des clés virtuelles et budgets.
- **Vertex / Claude** : modèles disponibles pour le projet et **région** à utiliser. Les modèles récents
  peuvent n'être servis que via les endpoints `global`, `us` ou `eu` et pas via une région précise
  comme europe-west1 → le cluster reste en europe-west4, seul l'appel Vertex cible `eu`/`global`.
- **Destroy GKE** : ressources réseau créées par des contrôleurs (NEG, firewall) et leur nettoyage.

---

## 13. Plan de réalisation par phases

Chaque phase se termine par un commit propre, des tests verts et une mise à jour du README.

| Phase | Contenu | Critère d'acceptation |
|---|---|---|
| 1 | Squelette dépôt, README, ADR, Makefile, `infra/bootstrap` | `terraform validate` OK ; README explique les prérequis |
| 2 | `infra/platform` (VPC, GKE, pools, IAM, WI) + `up.sh`/`down.sh`/`check-orphans` | Cycle `up` → `down` → `check-orphans` sans orphelin, deux fois de suite |
| 3 | Argo CD app-of-apps, Kyverno, ESO, kro, OpenCost | Tout `Synced/Healthy` après `make up` |
| 4 | RGD `TeamSpace` + politiques + exemples de claims + CI Kyverno | Un claim d'exemple crée l'espace d'équipe ; un claim invalide est refusé en CI et dans le cluster |
| 5 | Portail : catalogue, formulaires, ouverture de PR | Demande via UI → PR GitHub → merge → ressource créée |
| 6 | Chemin local kind, polish README (schéma, GIF de démo), `docs/demo.md` | Un visiteur lance `make local-up` sans compte cloud |
| 7 *(optionnelle)* | LiteLLM + Vertex (WI, sans clé) + RGD `LLMAccess` | Une équipe obtient une clé et appelle Claude via le gateway ; budget appliqué |
| 8 *(optionnelle)* | Assistant IA dans le portail (outils, confirmation avant PR) | Scénario de démo bout en bout scripté dans `docs/demo.md` |
| transverse | Workflow Claude Code : CLAUDE.md, skills, garde-fous, revues IA (§ 16) | En place dès la phase 1, revu à chaque phase |

---

## 14. Structure du dépôt cible

```
.
├── README.md                # pitch, schéma, démarrage rapide, coûts
├── PLAN.md                  # ce document
├── CLAUDE.md                # conventions pour Claude Code (§ 15)
├── Makefile
├── .env.example
├── docs/
│   ├── architecture.md
│   ├── demo.md              # scénario de démo pas à pas
│   └── adr/
├── infra/
│   ├── bootstrap/
│   └── platform/
├── gitops/
│   ├── root-app.yaml
│   ├── apps/                # une Application par composant
│   └── overlays/{gcp,local}/
├── platform/
│   ├── rgd/
│   └── policies/
├── claims/
│   └── examples/
├── portal/
│   ├── app/
│   ├── tests/
│   └── Dockerfile
├── scripts/                 # up.sh, down.sh, check-orphans.sh, local-up.sh
└── .github/workflows/       # lint, tests, kyverno sur claims, build image GHCR
```

---

## 15. À copier dans `CLAUDE.md`

```markdown
# Règles du projet
- Lire PLAN.md avant toute tâche ; travailler phase par phase (§ 13), ne pas anticiper la suivante.
- Vérifier la doc officielle pour toute syntaxe de CRD, chart ou API externe ; ne jamais deviner une version.
- Épingler toutes les versions (providers, charts, images). Jamais de `latest`.
- Aucun secret, clé JSON ou token dans le dépôt. Workload Identity uniquement.
- Aucun PVC, LoadBalancer ou IP statique sans ADR justifiant le nettoyage au destroy.
- Toute ressource GCP porte le label env=platform-ai-demo.
- Ne jamais exécuter `make up`, `terraform apply` ou `destroy` sans demande explicite de l'utilisateur.
- Après chaque phase : `make test` vert, README à jour, commit.
- Langue : code et identifiants en anglais, documentation en français.
```

---

## 16. Construire la plateforme avec Claude Code (transverse)

La plateforme est développée avec Claude Code. Le dépôt encadre l'agent pour qu'il reste dans
le plan et que le coût en tokens reste maîtrisé ([ADR 0019](docs/adr/0019-ai-assisted-workflow.md)).

| Risque | Garde-fou | Où |
|---|---|---|
| Dérive de périmètre (anticiper une phase, inventer une version) | Règles `CLAUDE.md`, phases de PLAN.md, ADR par décision | `CLAUDE.md`, `PLAN.md`, `docs/adr/` |
| Action destructrice ou facturée sans accord | Confirmation obligatoire : `terraform apply/destroy`, `make up/down/bootstrap`, `git push` | `.claude/settings.json` (`permissions.ask`) |
| Mauvais compte GCP (compte client) | Hook qui bloque `gcloud` hors config isolée + vérification d'identité des scripts | `.claude/hooks/guard-gcloud.sh`, `scripts/lib.sh` |
| Contexte trop lourd (coût de chaque requête) | `CLAUDE.md` court ; consignes à la demande (skills, `CLAUDE.md` par dossier) ; lectures interdites ; une session par tâche | `CLAUDE.md`, `.claude/skills/`, `permissions.deny` |
| Régression non vue | `make test` + CI ; revue Claude **locale** (`/code-review`) avant push | `Makefile`, `.github/workflows/ci.yml` |
| Revue IA trop chère | Revue CI déclenchée à l'ouverture / au label `ai-review`, diff seul, tours plafonnés, Sonnet, auth OIDC sans clé | `.github/workflows/claude-review.yml` |
| Supply chain des outils | Versions épinglées (SHA d'actions, version PyPI) et délai de 7 jours avant adoption | workflows |
