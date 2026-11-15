# 0021. kro au-dessus de Config Connector comme API de plateforme

- **Statut** : Accepté — remplace [0020](0020-crossplane-v2-platform-api.md)
- **Date** : 2026-11-15

## Contexte
La plateforme doit exposer en self-service des objets Kubernetes **et** des objets cloud GCP.
L'ADR 0020 retenait Crossplane v2. Un retour d'expérience public (r/kubernetes, « kro in 2026 »)
décrit un modèle en production : **opérateur de première partie comme couche primitive, kro
par-dessus comme catalogue** (« ils s'empilent plutôt que de rivaliser »), avec un rayon d'impact
plus petit que Crossplane. Il signale aussi qu'un rendu Helm de ressources Config Connector n'a
pas d'ordre et converge mal, ce que l'ordonnancement de kro résout.

## Décision
- **Config Connector** crée les ressources GCP des équipes. Installé par son **opérateur**
  (l'add-on GKE accuse jusqu'à 12 mois de retard selon la documentation Google), d'abord en
  **mode cluster** avec un compte de service GCP aux rôles minimaux (Workload Identity), puis en
  **mode namespacé** (un compte de service par équipe) comme évolution.
- **kro** compose le catalogue : une ResourceGraphDefinition `TeamSpace` (namespace, quotas,
  LimitRange, NetworkPolicy, RBAC, ServiceAccount + `StorageBucket`, `IAMServiceAccount` et
  liaisons IAM Config Connector). Version épinglée, adoptée après 7 jours.
- **Kyverno** fait le contrôle d'admission que kro n'a pas, en CI et dans le cluster, dont
  l'interdiction de l'annotation `cnrm.cloud.google.com/deletion-policy: abandon`.
- Chaque RGD est couverte par un test de bout en bout sur kind en CI avant d'en écrire une autre.

## Alternatives écartées
- Crossplane v2 (ADR 0020) : mature et multi-cloud, mais écosystème plus large de contrôleurs à
  maintenir et provider tiers pour GCP ; non natif GCP.
- Helm + ApplicationSet au-dessus de Config Connector : technologie éprouvée, mais rendu sans
  ordre ni statut agrégé par demande.
- Pipeline Terraform par équipe : vrai `plan`, pas de problème d'ordre au `down`, mais pas de
  réconciliation continue ni d'API dans le cluster.

## Conséquences
- kro est pré-stable (`v1alpha1`) : version épinglée, une seule RGD validée avant d'élargir le
  catalogue. Le versionnage des RGD reste un problème ouvert de l'écosystème ; sans impact sur
  une démo éphémère, documenté dans la série.
- Débogage : une instance kro bloquée indique rarement la ressource fautive ; l'erreur GCP est
  sur l'objet Config Connector. Procédure documentée dans `docs/`.
- Pas d'équivalent de `terraform plan` : la PR montre la demande, Kyverno la valide.
- Exception encadrée à l'ADR 0007, comme pour l'ADR 0020 : `make down` supprime les demandes,
  **attend que Config Connector ait supprimé les ressources GCP**, puis détruit le cluster ;
  `check-orphans` couvre buckets et comptes de service.
- Chemin local (kind) : les CRD Config Connector n'y existent pas ; une variante de la RGD sans
  ressource GCP, à préciser en phase 6.
