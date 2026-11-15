# 0003. State local pour la couche bootstrap

- **Statut** : Accepté
- **Date** : 2026-11-15

## Contexte
Le socle crée le bucket qui hébergera le state de la couche éphémère : il ne peut pas y stocker
son propre state avant d'exister.

## Décision
Le state de `infra/bootstrap/` reste **local** (`infra/bootstrap/terraform.tfstate`, ignoré par
git). Le socle contient peu de ressources, toutes recréables ou importables, et **aucune valeur
secrète** (le secret est créé vide, sa valeur est ajoutée hors Terraform). Un seul opérateur
l'applique, rarement.

## Alternatives écartées
- Bucket dédié créé à la main : une étape manuelle de plus pour un gain faible.
- Migration ultérieure dans le bucket (`terraform init -migrate-state`) : possible sans rien casser
  si plusieurs personnes doivent appliquer le socle.

## Conséquences
- En cas de perte du state local : `terraform import` des ressources (bucket, secret, budget).
- Le fichier `terraform.tfstate` du socle est à sauvegarder par l'opérateur.
