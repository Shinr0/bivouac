# 0003. State local pour la couche bootstrap

- **Statut** : Accepté
- **Date** : 2026-11-10

## Contexte
Le socle crée le bucket qui hébergera le state de la couche éphémère : il ne peut pas y stocker
son propre state avant d'exister (problème de l'œuf et de la poule).

## Décision
Le state de `infra/bootstrap/` reste **local** (`infra/bootstrap/terraform.tfstate`, ignoré par git).

Justification : le socle contient peu de ressources, toutes recréables ou importables ; il ne
contient **aucune valeur secrète** (le secret est créé vide, sa valeur est ajoutée hors Terraform) ;
un seul opérateur l'applique, rarement.

## Alternatives écartées
- Bucket dédié créé à la main : une étape manuelle de plus pour un gain faible.
- Migration du state dans le bucket après création (`terraform init -migrate-state`) :
  possible plus tard sans casser quoi que ce soit, si plusieurs personnes doivent appliquer le socle.

## Conséquences
- En cas de perte du state local : `terraform import` des ressources (bucket, secret, budget).
- Le fichier `terraform.tfstate` du socle doit être sauvegardé par l'opérateur.
