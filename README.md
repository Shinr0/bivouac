# Platform Engineering — Scaleway via Crossplane

Internal Developer Platform permettant de provisionner des instances Scaleway en
self-service depuis un portail web, sans passer par Git ni par la CI.

## Architecture

```
Forge (React, statique)
  └─ BFF Go — ServiceAccount Kubernetes
       └─ API Kubernetes : composite VirtualMachine dans forge-workloads
            └─ Crossplane v2 + provider-scaleway
                 └─ Instance, Block Volume, IP publique Scaleway
```

En parallèle, ArgoCD synchronise la configuration de la plateforme (Crossplane,
XRD, Compositions, manifests Forge) depuis `platform/`.

### Trois couches, trois propriétaires

| Couche | Outil | Contenu | Fréquence de changement |
|---|---|---|---|
| Infrastructure jour 0 | OpenTofu | VPC, Kapsule, IAM, ArgoCD | rare |
| Configuration plateforme | ArgoCD / Git | Crossplane, XRD, Compositions, Forge | à chaque évolution du contrat |
| Ressources utilisateurs | BFF → API Kubernetes | Composites `VirtualMachine` | en continu |

La séparation est délibérée : les demandes utilisateurs ne transitent pas par
Git. Un commit par VM créerait un couplage entre le rythme des utilisateurs et
celui de la plateforme, avec des conflits de merge et une latence inutile. Les
composites vivent dans etcd, où Crossplane les réconcilie ; ArgoCD ne gère que
ce qui doit être revu et versionné.

## Composants

| Composant | Rôle |
|---|---|
| **Scaleway Kapsule** | Cluster Kubernetes managé, control plane mutualisé gratuit |
| **Crossplane v2** | Réconcilie les composites vers l'API Scaleway |
| **provider-scaleway** | Ressources managées Instance, Block, VPC, RDB… |
| **ArgoCD** | GitOps sur la configuration de la plateforme |
| **BFF** (`bff/`) | API REST Go, seul détenteur des droits cluster |
| **Forge** (`frontend/`) | Portail Catalog / Demande / Suivi / Settings |

## Structure

```
bff/                      API Go (client dynamique Kubernetes)
frontend/                 Portail Forge (Vite + React + TypeScript)
platform/                 Configuration GitOps synchronisée par ArgoCD
  bootstrap/applications/   Applications ArgoCD (app-of-apps)
  crossplane/runtime/       Provider Scaleway, functions, ProviderConfig
  apis/                     XRD VirtualMachine + Composition Scaleway
  forge/                    Déploiements BFF et frontend, RBAC, Ingress
specs/                    Spec Kit — spec / plan / tasks
.specify/                 Constitution et templates Spec Kit
```

Le socle OpenTofu vit dans le dépôt voisin `infra-rcscale-scaleway`.

## Démarrage

### 1. Infrastructure

```bash
cd ../infra-rcscale-scaleway
export SCW_ACCESS_KEY=... SCW_SECRET_KEY=... SCW_DEFAULT_PROJECT_ID=...
cp terraform.tfvars.example terraform.tfvars   # renseigner project_id et platform_repo_url
tofu init && tofu apply

tofu output -raw kubeconfig > ~/.kube/forge-platform.yaml
export KUBECONFIG=~/.kube/forge-platform.yaml
```

OpenTofu crée le cluster puis installe ArgoCD et l'Application racine. ArgoCD
déroule ensuite les waves décrites dans `platform/README.md`.

### 2. Images

```bash
cd bff       && docker build -t rg.fr-par.scw.cloud/forge/forge-bff:0.1.0 . && docker push ...
cd frontend  && docker build -t rg.fr-par.scw.cloud/forge/forge-frontend:0.1.0 . && docker push ...
```

### 3. Développement local

```bash
# API, contre le cluster courant
cd bff && FORGE_CORS_ORIGINS=http://localhost:5173 go run ./cmd/server

# Portail, /api est proxyfié vers le BFF
cd frontend && npm install && npm run dev
```

Sans cluster, basculer **Settings → mode `mock`** : le portail fonctionne
entièrement en localStorage.

## Contrat self-service

`VirtualMachine` (`compute.platform.local/v1alpha1`), namespacé :

| Paramètre | Type | Défaut | Description |
|---|---|---|---|
| `vmName` | string | *requis* | Nom de l'instance |
| `machineType` | string | `DEV1-S` | Type commercial Scaleway |
| `zone` | string | `fr-par-1` | Zone de disponibilité |
| `imageFamily` | string | `ubuntu_jammy` | Label d'image |
| `diskSizeGb` | integer | `20` | Volume racine |
| `rootVolumeType` | string | `sbs_volume` | `sbs_volume` ou `l_ssd` |
| `dataDiskSizeGb` | integer | `0` | Volume de données additionnel |
| `staticIp` | boolean | `false` | IP publique réservée plutôt que dynamique |

Les valeurs acceptées sont servies par `GET /api/v1/catalog` : le formulaire du
portail se construit à partir de l'API, jamais à partir d'une liste codée en dur.

## Commandes utiles

```bash
kubectl get providers
kubectl get xrd,compositions
kubectl -n forge-workloads get virtualmachines
kubectl -n forge-workloads get servers.instance.scaleway.m.upbound.io
kubectl -n argocd get applications
kubectl -n crossplane-system logs -l app=crossplane
```

## Dépannage

- **Provider pas `Healthy`** : `kubectl describe provider provider-scaleway`, puis
  vérifier que le Secret `crossplane-scaleway-credentials` existe bien dans
  `crossplane-system`.
- **Composite bloqué en `provisioning`** :
  `kubectl -n forge-workloads describe virtualmachine <nom>` expose les
  conditions des ressources composées.
- **Application ArgoCD `OutOfSync` sur les XRD** : les CRD Crossplane n'étaient
  pas encore installées au moment du premier passage ; la wave suivante rejoue
  automatiquement grâce à la politique de retry.
- **`readyz` du BFF en échec** : la XRD `virtualmachines` n'est pas installée ou
  le RoleBinding n'a pas été appliqué.
