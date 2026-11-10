# Plateforme IA éphémère sur GCP — Spécification de déploiement

> Document d'entrée pour Claude Code. Il décrit **quoi** construire et **dans quel ordre**.
> Les choix d'outils sont arrêtés ; les versions exactes, la syntaxe des APIs et des charts
> doivent être **vérifiées dans la documentation officielle au moment de l'implémentation**
> (voir § 12 « Points à vérifier »). Ne rien inventer : en cas de doute, lire la doc ou demander.

---

## 1. Contexte et objectif

Vitrine technique d'un consultant freelance Cloud DevOps (GCP, Kubernetes, Terraform,
Platform Engineering) : une **plateforme interne de développement (IDP) en self-service, avec de
l'IA par-dessus**, publiée en dépôt GitHub public et démontrable en live devant un prospect.

Deux volets IA :
1. **LLM-as-a-Service** : une équipe demande un accès LLM (ou un modèle auto-hébergé) par un
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
5. **Coûts bornés par construction** : spot partout, GPU à zéro par défaut, budgets LLM plafonnés.
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
│        ├──► Vertex AI (Claude via Model Garden)        [par défaut, facturé à l'usage]        │
│        └──► KServe + vLLM sur nœud GPU L4 spot         [optionnel, à la demande]              │
│                                                                                               │
│  OpenCost (coûts par namespace/équipe)   Prometheus/Grafana (léger)                           │
│                                                                                               │
│  Pools : cpu-spot (2× e2-standard-4)    gpu-spot (g2-standard-4, L4, min=0)                   │
└───────────────────────────────────────────────────────────────────────────────────────────────┘
        ▲ Workload Identity (KSA ↔ GSA roles/aiplatform.user)        ▲ accès : port-forward / IAP
```

### Choix arrêtés (justifier chacun dans un ADR `docs/adr/`)

| Couche | Choix | Pourquoi |
|---|---|---|
| IaC socle GCP | Terraform (ou OpenTofu, compatible) | Standard du marché, certification de l'auteur |
| Cluster | GKE **Standard zonal** | Le crédit gratuit GKE couvre les frais de gestion d'un cluster zonal ; contrôle des node pools GPU |
| GitOps | Argo CD, pattern app-of-apps | Maîtrisé par l'auteur, visuel en démo |
| API développeur | **kro** (ResourceGraphDefinitions) ; Crossplane seulement si une ressource GCP doit être exposée en self-service (optionnel, phase 5) | kro est léger et natif K8s ; Crossplane reste pertinent pour le cloud |
| Politiques | Kyverno | Plus simple qu'OPA pour du 100 % Kubernetes ; utilisable en CI (CLI) |
| Portail | **Maison** : FastAPI + HTMX (pas de SPA) | Backstage jugé trop lourd ; le portail n'est qu'un client de l'API (les CRD) — remplaçable |
| Gateway LLM | LiteLLM proxy | API OpenAI-compatible, clés virtuelles, budgets, multi-backends |
| Modèles managés | Vertex AI (Claude via Model Garden) | Pas de GPU à maintenir, facturé au token |
| Modèles auto-hébergés | KServe + vLLM | KServe = control plane de serving CNCF ; vLLM = moteur |
| FinOps | OpenCost + budget GCP + budgets LiteLLM | Argument FinOps de l'offre |
| Images | GHCR, buildées par GitHub Actions | Gratuit, hors périmètre détruit |

**Explicitement hors périmètre** : Backstage, llm-d (trop lourd pour la démo), multi-cluster,
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
- Node pool `gpu-spot` (activable par variable `enable_gpu`) : g2-standard-4 (1× L4), spot,
  autoscaling **0**–1, taint `nvidia.com/gpu=present:NoSchedule`, installation du driver gérée par GKE.
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
| 3 | LiteLLM (+ Postgres pour clés/budgets) | Postgres **sans PVC** en démo (emptyDir), assumé et documenté |
| 3 | KServe (+ cert-manager si requis) | uniquement si `enable_gpu` ; documenter le mode de déploiement choisi |
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

### `LLMAccess`
Donne à une équipe un accès au gateway :
- `backend` : `vertex` (défaut) | `selfhosted`
- `model` : enum restreint (alias définis dans LiteLLM)
- `monthlyBudgetUsd` : entier, borné (ex. max 50)
- Si `selfhosted` : crée un `InferenceService` KServe (vLLM, petit modèle open-weights) sur le pool GPU,
  avec scale-to-zero si le mode KServe choisi le permet.
- La création de la clé virtuelle LiteLLM se fait par un petit **Job/contrôleur** (ou un hook
  du portail) qui appelle l'API admin LiteLLM et stocke la clé dans un Secret du namespace de l'équipe.
  Choisir l'option la plus simple et la documenter dans un ADR.

Exemples de claims dans `claims/examples/`.

### Politiques Kyverno (`platform/policies/`)
- Interdiction des PVC (démo).
- Labels obligatoires `team`, `owner` sur les claims.
- `monthlyBudgetUsd` ≤ plafond global.
- Interdiction des Services `type: LoadBalancer` hors namespace autorisé.
- Pods GPU uniquement depuis les ressources générées par `LLMAccess`.
- Les mêmes politiques tournent en CI avec la CLI Kyverno sur les PR de `claims/`.

---

## 7. Portail et assistant IA (`portal/`)

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
make gpu-on      # active enable_gpu (terraform apply ciblé) ou scale le pool
make gpu-off
make down        # destruction complète + contrôle d'orphelins
make check-orphans
make local-up    # chemin kind, sans GCP (§ 9)
make test        # pytest + lint YAML + kyverno CLI + terraform validate/fmt
```

### `make up` (scripts/up.sh)
1. Vérifier les prérequis (gcloud authentifié, projet, outils, quota GPU si `enable_gpu`).
2. `terraform -chdir=infra/platform init` (backend GCS) puis `apply -auto-approve`.
3. Récupérer les credentials du cluster.
4. Attendre que l'Application racine Argo CD et toutes ses enfants soient `Synced/Healthy`
   (timeout explicite, message clair en cas d'échec).
5. Smoke tests : le portail répond, le gateway répond, un appel Vertex de test aboutit.
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
- pas de GPU, pas de KServe ; backend LLM = mode « mock » ou clé API fournie par l'utilisateur
  via variable d'environnement (jamais committée) ;
- External Secrets remplacé par un Secret créé localement depuis `.env`.
Objectif : démontrer kro, Kyverno, le portail et le flux PR sans compte cloud.

---

## 10. Prérequis manuels (une seule fois, à documenter dans le README)

1. Projet GCP avec facturation (compte de la société).
2. **Activer les modèles Claude dans Vertex AI Model Garden** et accepter les conditions
   (par projet, non automatisable proprement).
3. **Demander le quota GPU** (L4, région europe-west4) si le volet auto-hébergé est utilisé —
   sur un projet neuf il est souvent à zéro ; le délai de réponse est variable.
4. Créer un token GitHub à portée minimale (fine-grained, accès au seul dépôt, PR + contents)
   et le stocker dans le secret `github-token` de Secret Manager.
5. `make bootstrap`.

---

## 11. Garde-fous de coûts

| Mécanisme | Où |
|---|---|
| Budget GCP + alertes | `infra/bootstrap` |
| Nœuds spot, cluster zonal, pas de NAT par défaut | `infra/platform` |
| Pool GPU min=0, `enable_gpu=false` par défaut | `infra/platform` |
| Budgets par clé LiteLLM + plafond global | config LiteLLM / RGD `LLMAccess` |
| Interdiction PVC / LoadBalancer non autorisés | Kyverno |
| `check-orphans` obligatoire en fin de `down` | scripts |

Ordres de grandeur relevés en novembre 2026 (prix tiers, **à revérifier** dans le calculateur GCP) :
- e2-standard-4 spot à europe-west4 ≈ 0,056 $/h ; g2-standard-4 (L4) spot ≈ 0,37 $/h.
- Frais de gestion GKE 0,10 $/h, couverts par le crédit gratuit (74,40 $/mois) pour **un** cluster zonal.
- Une démo de 2 h avec GPU ≈ quelques dollars.

---

## 12. Points à vérifier avant d'implémenter (ne pas deviner)

- **Versions** à épingler (dernières stables) : providers Terraform google/helm/kubernetes, charts
  Argo CD, Kyverno, kro, External Secrets, OpenCost, KServe, LiteLLM. Épingler partout, pas de `latest`.
- **kro** : `apiVersion` exacte de `ResourceGraphDefinition`, syntaxe du simple schema
  (marqueurs `required`, `default`, `enum`, `minimum`/`maximum`), `includeWhen` — doc kro.run.
- **LiteLLM** : syntaxe `model_list` pour Claude sur Vertex (`vertex_ai/<model>`, `vertex_project`,
  `vertex_location`), endpoint `global`/`eu`, auth par ADC/Workload Identity, API admin des clés virtuelles et budgets.
- **Vertex / Claude** : modèles disponibles pour le projet et **région** à utiliser. Les modèles récents
  peuvent n'être servis que via les endpoints `global`, `us` ou `eu` et pas via une région précise
  comme europe-west1 → le cluster reste en europe-west4, seul l'appel Vertex cible `eu`/`global`.
- **KServe** : mode de déploiement le plus léger compatible scale-to-zero (Knative ou non), prérequis
  (cert-manager, Gateway API), et ressource adaptée (`InferenceService` vs `LLMInferenceService`, ce dernier plus récent).
- **GKE GPU** : installation automatique du driver sur node pool L4, taints/tolérances.
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
| 5 | LiteLLM + Vertex (WI, sans clé) + RGD `LLMAccess` (backend vertex) | Une équipe obtient une clé et appelle Claude via le gateway ; budget appliqué |
| 6 | Portail : catalogue, formulaires, ouverture de PR | Demande via UI → PR GitHub → merge → ressource créée |
| 7 | Assistant IA (outils, confirmation avant PR) | Scénario de démo bout en bout scripté dans `docs/demo.md` |
| 8 | Volet GPU : pool L4 spot, KServe + vLLM, `backend: selfhosted` | `make gpu-on`, appel au modèle auto-hébergé via le même gateway, `make gpu-off` |
| 9 | Chemin local kind, polish README (schéma, GIF de démo), article de blog | Un visiteur lance `make local-up` sans compte cloud |

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
