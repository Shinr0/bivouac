# 0001. Consigner les décisions d'architecture

- **Statut** : Accepté
- **Date** : 2026-11-15

## Contexte
Bivouac est un laboratoire public : un lecteur (prospect, recruteur, contributeur) doit
comprendre *pourquoi* chaque outil est là, et pas seulement *comment* il est câblé. Les essais
successifs d'un labo font évoluer les choix : il faut en garder la trace.

## Décision
- Chaque choix structurant fait l'objet d'un ADR court (format de Michael Nygard) dans
  `docs/adr/`, numéroté.
- **Un ADR n'est jamais modifié sur le fond ni supprimé.** Une décision révisée donne lieu à un
  nouvel ADR qui la remplace ; l'ancien ne change que sa ligne de statut :
  `Remplacé par [NNNN](NNNN-….md) (AAAA-MM-JJ)`.
- Un ADR peut être accepté avant son implémentation ; la phase concernée est indiquée dans le
  statut, et les détails (syntaxe, versions) sont confirmés à l'implémentation.
- Le skill de projet `new-adr` applique ces règles.

## Conséquences
- L'historique des décisions se lit dans `docs/adr/` : statut, remplacements, dates.
- Toute PR qui introduit un outil, une ressource facturable ou un changement de périmètre ajoute
  un ADR. Modèle : `0000-template.md`.
