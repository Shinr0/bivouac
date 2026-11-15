# 0002. Terraform en deux couches : socle persistant et plateforme éphémère

- **Statut** : Accepté
- **Date** : 2026-11-15

## Contexte
La plateforme ne doit exister que le temps d'un essai ou d'une démo (`make up` / `make down`),
mais certains éléments doivent survivre : le state Terraform lui-même, le secret GitHub, le budget,
l'activation des APIs (lente et gratuite).

## Décision
- `infra/bootstrap/` : socle **persistant**, appliqué une fois (`make bootstrap`), jamais détruit
  par `make down`. Coût quasi nul (un bucket de quelques Ko, un secret, un budget).
- `infra/platform/` : couche **éphémère**, state dans le bucket GCS du socle.
- Terraform (compatible OpenTofu ; aucune fonctionnalité propre à l'un ou l'autre n'est utilisée).
- Toutes les ressources portent `env=bivouac` (via `default_labels` du provider) pour permettre
  le contrôle d'orphelins.

## Alternatives écartées
- Une seule couche avec `-target` pour épargner le socle : fragile, facile à détruire par erreur.
- Pulumi / Config Connector pour les fondations : moins répandus pour ce rôle chez les clients visés.

## Conséquences
- Le `destroy` de la couche éphémère ne peut pas toucher au socle.
- Les APIs ne sont jamais désactivées (`disable_on_destroy = false`).
