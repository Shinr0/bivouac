# 0016. Référentiel d'équipes et d'applications dans Git, au format Backstage

- **Statut** : Accepté (implémentation en phase 4, affichage dans le portail en phase 5)
- **Date** : 2026-11-15

## Contexte
Une plateforme doit savoir **à qui appartient quoi** : quelles équipes existent, qui les manage,
comment les contacter, quelles applications elles opèrent et sur quel socle technique, et à quelle
équipe rattacher une ressource ou un coût. Sans référentiel, pas d'attribution des coûts, pas de
propriétaire pour une ressource orpheline, pas d'escalade en cas d'incident.

Contrainte propre à Bivouac : la plateforme est éphémère, le référentiel ne l'est pas. Il doit
survivre à `make down`.

## Décision
- **Source de vérité : Git**, dossier `catalog/`, un fichier YAML par entité, au **format
  descripteur du catalogue Backstage** (`apiVersion: backstage.io/v1alpha1`), standard de fait :
  - `Group` (équipe) : `spec.type: team`, `spec.profile` (nom affiché, **adresse de groupe**),
    `spec.parent` (rattachement hiérarchique, donc le management), `spec.members` ;
  - `User` : uniquement si nécessaire, pour les membres ;
  - `System` et `Component` (application) : `spec.owner` (équipe), `spec.lifecycle`, `spec.type`,
    `spec.system`, `spec.dependsOn` ; le socle technique en `metadata.tags` ;
  - `Resource` : les ressources GCP issues d'un `TeamSpace` (bucket, compte de service).
- Ce que le format ne prévoit pas passe par des **annotations préfixées** `bivouac/` (manager,
  canal de discussion, astreinte, centre de coût) et par `metadata.links` (runbook, tableaux de bord).
- **Intégration à la plateforme** :
  - une demande `TeamSpace` doit référencer un `Group` existant ; la CI le vérifie ;
  - le label `team=<groupe>` est posé sur les ressources Kubernetes et GCP : OpenCost attribue
    les coûts, `check-orphans` sait à qui appartient un reste ;
  - le portail (ADR 0011) affiche une page par équipe : manager, contacts, applications,
    ressources, coûts.
- **Données personnelles** : le dépôt est public, les équipes et personnes du labo sont fictives
  et les contacts sont des adresses de groupe. En contexte client, `Group` et `User` sont
  synchronisés depuis l'annuaire (Google Workspace / Cloud Identity) plutôt qu'écrits à la main.

## Alternatives écartées
- Backstage et sa base de catalogue : lourd à opérer pour un labo éphémère ; le format est
  conservé, le moteur non. Un client qui adopte Backstage ingère `catalog/` tel quel.
- Une CRD `Team` dans le cluster : disparaît à chaque `make down`.
- Tableur ou CMDB : ni versionné, ni validé, ni relu en pull request.
- Portail SaaS (Port, Cortex, OpsLevel…) : dépendance externe pour une donnée centrale.

## Conséquences
- Toute évolution d'équipe ou d'application passe par une pull request relue, comme le reste.
- Le format ne publiant pas de schéma de validation sur sa page de référence, la CI valide le
  catalogue avec un schéma propre au dépôt (champs obligatoires, références résolues, annotations
  `bivouac/` attendues).
