# platform/ — configuration GitOps

Tout ce dossier est synchronisé par ArgoCD. L'Application racine
`platform-bootstrap`, créée par OpenTofu, pointe sur
`platform/bootstrap/applications` et déploie les Applications enfants.

## Ordre de synchronisation

Les `sync-wave` garantissent que les CRD existent avant les ressources qui les
utilisent.

| Wave | Application | Contenu |
|---|---|---|
| 0 | `ingress-nginx` | Contrôleur d'ingress, provisionne le Load Balancer Scaleway |
| 0 | `crossplane` | Chart Helm Crossplane (controllers + CRD de base) |
| 1 | `crossplane-runtime` | `provider-scaleway`, functions, `ClusterProviderConfig` |
| 2 | `platform-apis` | Namespace `forge-workloads`, XRD `VirtualMachine`, Composition Scaleway |
| 3 | `forge` | BFF, frontend, RBAC, Ingress |

## Arborescence

```
platform/
  bootstrap/applications/   Applications ArgoCD (app-of-apps)
  crossplane/runtime/       Provider Scaleway, functions, ProviderConfig
  apis/                     XRD + Composition + namespace des workloads
  forge/                    Manifests du portail et de son BFF
```

## Contrat exposé

La XRD `virtualmachines.compute.platform.local` est en `scope: Namespaced`
(Crossplane v2). Les claims n'existent plus : l'objet créé par le BFF dans
`forge-workloads` est directement le composite.

```yaml
apiVersion: compute.platform.local/v1alpha1
kind: VirtualMachine
metadata:
  name: my-dev-server
  namespace: forge-workloads
spec:
  parameters:
    vmName: my-dev-server
    machineType: DEV1-S
    zone: fr-par-1
    imageFamily: ubuntu_jammy
    diskSizeGb: 20
    rootVolumeType: sbs_volume
    dataDiskSizeGb: 0
    staticIp: false
```

La Composition rend jusqu'à trois ressources managées, dans le namespace du
composite :

| Ressource | Condition | Groupe |
|---|---|---|
| `Server` | toujours | `instance.scaleway.m.upbound.io` |
| `IP` | `staticIp: true` | `instance.scaleway.m.upbound.io` |
| `Volume` | `dataDiskSizeGb > 0` | `block.scaleway.m.upbound.io` |

Le volume de données est rattaché au serveur au tour de réconciliation suivant
sa création, le champ `additionalVolumeIds` du provider n'exposant pas de
sélecteur par référence.

## Avant le premier déploiement

Remplacer `https://github.com/Shinr0/platform-project.git` par l'URL réelle du
dépôt dans les Applications enfants (`20-`, `30-`, `40-`) et dans la variable
`platform_repo_url` du module OpenTofu.

Les images `rg.fr-par.scw.cloud/forge/forge-bff` et `forge-frontend` doivent
être construites et poussées avant la synchronisation de la wave 3.

## Dépendance externe

Le `ClusterProviderConfig` référence le Secret
`crossplane-system/crossplane-scaleway-credentials`, créé par OpenTofu car il
contient une clé API IAM. C'est la seule ressource du cluster qui n'est pas
dans Git.
