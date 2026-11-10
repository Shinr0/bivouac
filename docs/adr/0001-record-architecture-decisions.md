# 0001. Consigner les décisions d'architecture

- **Statut** : Accepté
- **Date** : 2026-11-10

## Contexte
Ce dépôt est une vitrine : un lecteur (prospect, recruteur) doit comprendre *pourquoi* chaque
outil est là, pas seulement *comment* il est câblé.

## Décision
Chaque choix structurant fait l'objet d'un ADR court (format de Michael Nygard) dans `docs/adr/`,
numéroté, jamais réécrit : une décision révisée est remplacée par un nouvel ADR qui la référence.
Les ADR des phases à venir sont rédigés dès maintenant avec le statut « Accepté » pour le choix
d'outil ; les détails d'implémentation sont confirmés dans la phase concernée.

## Conséquences
- Modèle : `0000-template.md`.
- Toute PR introduisant un nouvel outil ou une nouvelle ressource facturable ajoute un ADR.
