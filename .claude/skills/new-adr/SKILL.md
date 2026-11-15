---
name: new-adr
description: Record an architecture decision in docs/adr/ for this repo. Use whenever a tool, resource, version policy or scope choice is made, changed or reverted.
---

# Nouvel ADR

1. Numéro : le plus grand `docs/adr/NNNN-*.md` + 1. Nom de fichier en anglais, kebab-case.
2. Partir de `docs/adr/0000-template.md`. Contenu en français, court :
   Statut, Date (aujourd'hui), Contexte, Décision, Alternatives écartées, Conséquences.
3. Si la décision en remplace une autre : ne pas réécrire l'ancien ADR, seulement sa ligne
   Statut → `Remplacé par [NNNN](NNNN-....md) (AAAA-MM-JJ)`.
4. Régénérer l'index `docs/adr/README.md` (une ligne par ADR : numéro, titre, statut).
5. Répercuter la décision là où elle se voit : `PLAN.md` (section concernée), README si besoin.
6. `make test`, puis commit `docs(adr): NNNN <titre>` si l'utilisateur l'a demandé.
