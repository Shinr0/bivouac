# Plan éditorial — « Une IDP éphémère, plan par plan »

Fil rouge LinkedIn + blog adossé à la construction de ce dépôt (PLAN.md).
Révisé le 2026-11-15 : l'IA *intégrée* à la plateforme (LiteLLM, Vertex, assistant) passe en
épisode bonus ; l'IA entre dans la série par **la façon de construire** (Claude Code, épisode 1).

Articles : dépôt `Shinr0.github.io`, branche `feature/serie-idp-ephemere-ep0`,
`content/posts/idp-ephemere/`.

---

## 1. Le landscape et ce que la plateforme en couvre

Source : [platformengineering.org/platform-tooling](https://platformengineering.org/platform-tooling)
(landscape du 2026-03-17, ~210 fiches, 38 catégories), consulté le 2026-11-15.

| Plan | Catégories du landscape | Brique | Phase |
|---|---|---|---|
| Developer Control | VCS · Developer Portal · Workload Specification | GitHub · portail FastAPI + HTMX · API kro (`TeamSpace`) | 4, 5 |
| Integration & Delivery | CI · Registry · CD · Platform Orchestrator · IaC | GitHub Actions · GHCR · Argo CD · kro · Terraform | 1–4 |
| Resource | Cloud · Compute · Cluster Management · Networking | GCP · GKE zonal spot · VPC · Config Connector | 2, 3 |
| Observability | Monitoring & Logging · FinOps | Prometheus/Grafana · OpenCost · budget GCP · `check-orphans` | 2, 3 |
| Security | Policy Control · Secrets · ID Management · Network Security | Kyverno · Secret Manager + ESO · Workload Identity · NetworkPolicy | 1, 3, 4 |
| *(Copilot/Agent/LLM — vide dans le landscape)* | | Claude Code pour construire ; LiteLLM + Vertex en bonus | transverse, 7–8 |

Absents du landscape mais assumés (ADR) : kro, Config Connector, OpenCost, External Secrets, Vertex/LiteLLM.

---

## 2. Fil rouge

**Promesse** : « Je construis une Internal Developer Platform complète sur GCP, qui n'existe que le
temps d'une démo, avec un agent IA comme co-développeur sous garde-fous. Un plan à la fois. »

À chaque épisode : la carte des 5 plans (brique du jour surlignée) · le compteur réel (`up`,
`down`, coût) · une décision (ADR) · « ce que je dirais à un client ».

**Règle d'or** : publier seulement une fois la phase terminée et mesurée. Chiffres réels uniquement.

Rythme : 1 épisode / 1–2 semaines ; post LinkedIn (1 200–1 800 caractères) + lien de l'article en
premier commentaire. Hashtags (3 max) : `#PlatformEngineering` `#GoogleCloud` `#Kubernetes`,
ou `#ClaudeCode` `#DevOps` `#IA` selon l'épisode.

---

## 3. Épisodes

| # | Sujet | Prérequis | État |
|---|---|---|---|
| 0 | Le plan de bataille : une IDP qui s'autodétruit | Phase 1 | article rédigé |
| 1 | Construire sa plateforme avec Claude Code, sans dérive ni facture surprise | Outillage (fait) | article rédigé |
| 2 | Resource : `make up`, `make down`, zéro orphelin (+ coût réel d'une revue IA en CI) | Phase 2 | — |
| 3 | Integration & Delivery : GitOps, Terraform n'installe qu'une chose | Phase 3 | — |
| 4 | Security : zéro secret, zéro clé (GCP et CI) | Phase 3 | — |
| 5 | Developer Control : kro + Config Connector — objets cloud et K8s en self-service, validés par Kyverno (et ce qui reste imparfait : débogage, versionnage, pas de plan) | Phase 4 | — |
| 6 | Observability : combien coûte une équipe ? | Phases 3–4 | — |
| 7 | Developer Control : un portail sans Backstage | Phase 5 | — |
| 8 | Bilan, chiffres, version locale kind | Phase 6 | — |
| bonus | L'IA dans la plateforme : LLM par équipe, assistant qui ouvre des PR | Phases 7–8 | optionnel |

À capturer pendant le développement : chronos `up`/`down`, sortie `check-orphans`, facture du
jour, UI Argo CD, claim refusé (CI + admission), tableau OpenCost, GIF du portail, et à chaque
épisode : contexte consommé par session (`/context`, `/cost`) et coût des revues IA.

---

## 4. Post LinkedIn — Épisode 0

> Une plateforme interne de développement, ça coûte souvent plus cher à entretenir qu'à héberger.
>
> Alors j'en construis une qui n'existe que le temps d'une démo.
>
> make up → un projet GCP vide devient une plateforme en self-service.
> make down → tout disparaît. Et un script échoue s'il reste un seul disque facturé.
>
> Point de départ : le landscape de platformengineering.org. 5 plans, 38 catégories, ~210 outils.
> Et une case vide : la catégorie « Copilot/Agent/LLM » ne contient aucun outil.
>
> Ma règle : le plus petit chemin qui traverse les 5 plans.
> 🧭 Developer Control → une API kro, un portail léger (pas de Backstage)
> 🚚 Integration & Delivery → Terraform, Argo CD
> ☁️ Resource → GKE zonal spot + Config Connector : un bucket et un compte de service GCP par équipe, à la demande
> 🔎 Observability → OpenCost, budget
> 🔐 Security → Kyverno, Workload Identity, zéro secret dans Git
>
> Et l'IA ? Avant de la mettre dans la plateforme, je m'en sers pour la construire : Claude Code
> écrit le code, sous des garde-fous que je publie avec lui. Ce sera l'épisode 1.
>
> Un plan par épisode, des chiffres réels à chaque fois. Tout est open source.
>
> Ce que je dirais à un client : ne commencez pas par le portail. Commencez par le plus petit
> chemin qui traverse les cinq plans.
>
> #PlatformEngineering #GoogleCloud #Kubernetes

Premier commentaire : lien de l'article + https://github.com/Shinr0/platform-project

---

## 5. Post LinkedIn — Épisode 1

> Mon agent IA a bloqué son propre commit.
>
> Je construis ma plateforme avec Claude Code. Le sujet n'est plus « sait-il coder ? » mais
> « comment le garder sur les rails, et combien il coûte ? ».
>
> Premier incident, avant tout déploiement : ma config gcloud pointait vers le compte d'un client.
> Une règle écrite (« utilise mon compte perso ») ne suffit pas. J'ai ajouté un garde-fou
> exécutable : un hook qui bloque tout gcloud hors de la config isolée du projet.
>
> Quelques minutes plus tard, il bloquait le git commit de l'agent… parce que le message de commit
> contenait le mot « gcloud ». Leçon : un garde-fou est du code. Il se teste. (12 cas dans
> make test, exécutés en CI.)
>
> Le cadre complet :
> 📜 Règles et phases écrites (CLAUDE.md, PLAN.md, ADR)
> 🛑 Confirmation obligatoire pour apply, destroy, push
> ✂️ Un CLAUDE.md de 24 lignes au lieu de « relis le plan de plus de 20 Ko avant chaque tâche »
> 🧩 Des consignes chargées à la demande (skills) et des lectures interdites (states, .env)
> 🔍 Revue IA en CI à coût borné : à l'ouverture de PR seulement, diff seul, tours plafonnés, Sonnet
> 🔐 Aucune clé API : authentification OIDC, comme Workload Identity côté GCP
> 📦 Aucune dépendance de moins de 7 jours
>
> Ce que je dirais à un client : un agent de code, c'est un nouveau membre de l'équipe, rapide et
> sans mémoire. Onboarding écrit, moindre privilège, contrôles testés, revue, budget.
>
> #ClaudeCode #PlatformEngineering #DevOps

Premier commentaire : lien de l'article + lien vers `.claude/` dans le dépôt.
