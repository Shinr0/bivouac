# 0009. API de plateforme : kro au-dessus de Config Connector

- **Statut** : Accepté (implémentation en phases 3–4)
- **Date** : 2026-11-15

## Contexte
Les équipes doivent demander en self-service, par un objet simple versionné dans Git, à la fois
des objets Kubernetes (namespace, quotas, droits) et des objets cloud GCP (bucket, compte de
service), sans connaître les dizaines de ressources générées derrière.

## Décision
- **Config Connector**, l'opérateur Kubernetes de Google pour les ressources GCP, est la couche
  primitive. Installé par son **opérateur** (l'add-on GKE accuse jusqu'à 12 mois de retard selon
  la documentation Google), d'abord en **mode cluster** avec un compte de service GCP aux rôles
  minimaux (Workload Identity), puis en **mode namespacé** (un compte de service par équipe).
- **kro** compose le catalogue par-dessus : une ResourceGraphDefinition par capacité, qui crée les
  ressources dans l'ordre de leurs dépendances et expose un statut. Première capacité,
  `TeamSpace` : namespace, quotas, LimitRange, NetworkPolicy, RBAC, ServiceAccount, plus
  `StorageBucket`, `IAMServiceAccount` et liaisons IAM Config Connector.
- **Kyverno** assure le contrôle d'admission que kro n'a pas (ADR 0010).
- Versions épinglées ; une seule RGD, couverte par un test de bout en bout sur kind en CI, avant
  d'élargir le catalogue.

## Alternatives écartées
- Crossplane v2 + provider GCP : mature et multi-cloud, mais un écosystème entier de contrôleurs
  à maintenir et un provider GCP tiers.
- Helm (+ ApplicationSet) au-dessus de Config Connector : technologie éprouvée, mais pas d'API
  typée dans le cluster, pas de statut par demande, et un rendu **sans ordre** de création.
- Pipeline Terraform par équipe : vrai `plan`, mais pas de réconciliation continue ni d'API
  dans le cluster.
- kro seul : ne crée que des objets Kubernetes.

## Conséquences
- kro est pré-stable (`v1alpha1`) ; le versionnage des RGD reste un problème ouvert de
  l'écosystème, sans impact sur un labo éphémère.
- Débogage : une instance kro bloquée nomme rarement la ressource fautive ; l'erreur GCP figure
  sur l'objet Config Connector. Procédure à documenter dans `docs/`.
- Pas d'équivalent de `terraform plan` : la PR montre la demande, Kyverno la valide.
- Ressources GCP créées par un contrôleur : exception encadrée par l'ADR 0007.
- Chemin local (kind) : variante de chaque RGD sans ressource GCP.
