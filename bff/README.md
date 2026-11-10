# Forge BFF

Backend For Frontend du portail Forge. Il traduit les appels REST du portail en
composites Crossplane `VirtualMachine` dans le cluster.

## Pourquoi un BFF

Sans lui, le navigateur devrait porter un token Kubernetes ou déclencher un
workflow CI pour générer du YAML. Les deux approches déplacent la logique
plateforme dans le client. Ici :

- le portail reste un artefact statique, sans aucun secret ;
- l'authentification cluster passe par le ServiceAccount du pod ;
- la validation est appliquée côté serveur, avant l'API Kubernetes ;
- les valeurs autorisées sont servies par l'API, donc dérivées de la XRD.

Aucun YAML n'est généré ni stocké : l'objet Kubernetes **est** l'état désiré,
et Crossplane le réconcilie en continu.

## API

| Méthode | Route | Description |
|---|---|---|
| `GET` | `/healthz` | Liveness |
| `GET` | `/readyz` | Readiness : vérifie l'API Kubernetes et la présence de la XRD |
| `GET` | `/api/v1/catalog` | Catalogue de services et valeurs acceptées par le formulaire |
| `GET` | `/api/v1/vms` | Liste les demandes, triées par date décroissante |
| `POST` | `/api/v1/vms` | Crée un composite `VirtualMachine` |
| `GET` | `/api/v1/vms/{name}` | Détail d'une demande |
| `DELETE` | `/api/v1/vms/{name}` | Supprime le composite et les ressources Scaleway associées |

### Statuts

Les conditions Crossplane sont traduites en un vocabulaire lisible par le portail :

| Conditions | Statut |
|---|---|
| aucune | `submitted` |
| `Synced=True`, `Ready≠True` | `provisioning` |
| `Ready=True` | `ready` |
| `Synced=False` | `failed` (le message de la condition est remonté) |

## Configuration

| Variable | Défaut | Rôle |
|---|---|---|
| `FORGE_ADDR` | `:8080` | Adresse d'écoute |
| `FORGE_NAMESPACE` | `forge-workloads` | Namespace des composites |
| `FORGE_CORS_ORIGINS` | *(vide)* | Origines autorisées, séparées par des virgules |
| `KUBECONFIG` | *(vide)* | Kubeconfig utilisé hors cluster uniquement |

En production, le frontend et l'API partagent la même origine derrière
l'Ingress : `FORGE_CORS_ORIGINS` reste vide.

## Développement

```bash
go test ./...
go vet ./...

# Hors cluster : utilise le kubeconfig courant
FORGE_CORS_ORIGINS=http://localhost:5173 go run ./cmd/server
```

## Image

```bash
docker build -t rg.fr-par.scw.cloud/forge/forge-bff:0.1.0 .
docker push rg.fr-par.scw.cloud/forge/forge-bff:0.1.0
```

Le binaire est statique et tourne dans une image distroless en `nonroot`, avec
un système de fichiers racine en lecture seule.

## Permissions

Le ServiceAccount `forge-bff` ne dispose que d'un `Role` sur les composites
`virtualmachines` du namespace `forge-workloads`. Il ne peut ni lire de Secret
ni manipuler directement les ressources managées Scaleway.
