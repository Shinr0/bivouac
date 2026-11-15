# 0020. Crossplane v2 comme API de plateforme, avec des objets cloud

- **Statut** : Accepté — remplace [0009](0009-kro-developer-api.md)
- **Date** : 2026-11-15

## Contexte
L'ADR 0009 retenait kro, qui ne compose que des objets Kubernetes. La plateforme doit exposer en
self-service des **objets cloud GCP** en plus des objets Kubernetes : c'est ce qu'attend un client
d'une IDP. kro ne peut créer de ressources GCP qu'à travers un second contrôleur (Config Connector),
et son API est encore en `v1alpha1` (v0.10.0 au 2026-11-15). Helm, lui, ne fournit ni API typée
dans le cluster, ni statut par demande, ni réconciliation de l'ensemble.

## Décision
- **Crossplane v2** et le provider GCP (familles de ressources strictement nécessaires) exposent
  les capacités de la plateforme : une définition (XRD) et une composition par capacité, sous
  `platform/apis/`. Une même composition crée les objets GCP et les objets Kubernetes.
- Première capacité, `TeamSpace` : namespace, quotas, LimitRange, NetworkPolicy, RoleBinding,
  **bucket GCS** et **compte de service GCP** lié au compte Kubernetes de l'équipe par Workload
  Identity. Des ressources **rapides à créer et quasi gratuites** ; pas de Cloud SQL, Memorystore
  ou autre ressource lente ou facturée au temps dans la démo.
- Ressources managées labellisées `env=platform-ai-demo` et `team=<nom>`.
- La syntaxe v2 (XR namespacés, groupes d'API des ressources managées, fonctions de composition,
  politiques de gestion) est vérifiée dans la documentation officielle au moment de l'implémentation.

## Alternatives écartées
- kro + Config Connector : natif GCP et léger, mais deux outils et kro en alpha.
- Config Connector seul : pas d'abstraction, les équipes manipuleraient des ressources GCP brutes.
- Helm (un fichier de valeurs par équipe) : du templating, pas une API ; ne crée pas d'objets cloud.

## Conséquences
- Exception encadrée à l'[ADR 0007](0007-no-persistent-cloud-resources-from-k8s.md) : des
  ressources cloud sont créées par un contrôleur. Leur nettoyage devient une étape de `make down` :
  supprimer les demandes, **attendre que Crossplane ait détruit les ressources GCP**, puis
  seulement détruire le cluster.
- Kyverno refuse toute ressource managée configurée pour survivre à sa suppression.
- `check-orphans` couvre aussi les buckets et comptes de service portant le label du projet.
- Le compte de service GCP de Crossplane (Workload Identity) ne reçoit que les rôles requis par
  les types de ressources exposés, jamais `owner` ni `editor`.
- Crossplane est plus lourd que kro (CRD du provider, mémoire) : n'installer que les familles de
  ressources utilisées, et mesurer l'impact sur le pool de nœuds.
- Le chemin local (kind) utilise une composition sans ressources cloud pour la même définition.
