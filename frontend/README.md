# Forge — frontend

Portail IDP **Forge** (Catalog · Demande · Suivi · Settings).

## Prérequis

Node.js 20+ (ce poste utilise fnm / Node 22).

```bash
cd frontend
npm install
npm run dev
```

Le serveur de développement proxyfie `/api` vers `http://localhost:8080`, où
tourne le BFF. Surcharger avec `FORGE_BFF_URL` si nécessaire.

Build production :

```bash
npm run build
npm run lint
npm run preview
```

## Modes de provisionnement

| Mode | Comportement |
|---|---|
| `crossplane` *(défaut)* | Appelle le BFF, qui crée un composite `VirtualMachine` dans le cluster |
| `mock` | Tout en localStorage, sans cluster ni backend |

Le portail ne stocke aucun secret : l'authentification auprès de Kubernetes est
portée par le ServiceAccount du BFF. Settings ne contient que le mode et l'URL
de base de l'API, relative (`/api/v1`) par défaut puisque l'Ingress sert le
frontend et l'API sur la même origine.

## Image

```bash
docker build -t rg.fr-par.scw.cloud/forge/forge-frontend:0.1.0 .
```

Les assets buildés sont servis par nginx non privilégié sur le port 8080, avec
un fallback `index.html` pour le routage React Router.
