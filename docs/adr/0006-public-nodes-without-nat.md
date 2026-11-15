# 0006. Nœuds publics sans Cloud NAT par défaut

- **Statut** : Accepté
- **Date** : 2026-11-15

## Contexte
Des nœuds privés nécessitent Cloud NAT pour sortir sur Internet (images, GitHub), soit environ
33 $/mois s'il restait actif — et un routeur NAT est une ressource orpheline classique.

## Décision
Par défaut, nœuds avec IP publique éphémère et pare-feu strict (aucune entrée hors trafic interne
et plages Google requises). Une variable `private_nodes = true` bascule en nœuds privés + Cloud
NAT pour montrer la configuration « production ».

## Conséquences
- Surface d'exposition des nœuds à maîtriser par le pare-feu.
- `check-orphans` contrôle explicitement les routeurs NAT.
