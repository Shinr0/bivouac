---
name: write-episode
description: Draft the Bivouac blog episode (Hugo site repo) and its LinkedIn post once a PLAN.md phase is closed, from real measurements only, then walk the publication checklist. Use when the user asks to write, prepare or publish the next episode or article of the series.
---

# Rédiger un épisode de la série Bivouac

Dépôts : code dans `~/perso/bivouac`, site dans `~/perso/Shinr0.github.io` (lire son `CLAUDE.md`).
Plan éditorial et brouillons LinkedIn : `~/perso/bivouac/content/plan-editorial.md`.

## 1. Préconditions
- La phase du `PLAN.md` correspondant à l'épisode est close (skill `close-phase`).
- Les mesures réelles de la phase sont disponibles (durées `make up`/`make down`, coûts, sorties de
  commandes, captures). S'il en manque, **s'arrêter et demander** : aucun chiffre n'est estimé.

## 2. Matière
Lire, et seulement cela : l'entrée de l'épisode dans le plan éditorial, la section de `PLAN.md`
de la phase, les ADR créés ou concernés, `git log` de la phase.

## 3. Article (site, branche `episode-NN`)
- Dossier `content/posts/bivouac/NN-slug/`, front matter selon le `CLAUDE.md` du site.
- 1 200 à 2 000 mots, en français, ton factuel. Trame : le problème → ce qui a été construit →
  les décisions (liens vers les ADR sur GitHub) → **les chiffres mesurés** (tableau) → ce qui a
  raté ou surpris → les limites assumées → « Ce que je dirais à un client » → l'épisode suivant.
- Visuel : `scripts/bivouac-cover.py content/posts/bivouac/NN-slug --label "ÉPISODE N · …"
  --subtitle "…" --caption "…"` (sous-titre ~50 caractères, légende ~80).

## 4. Vérifications avant de proposer la PR
- Chaque chiffre renvoie à une mesure ou une sortie de commande de la phase.
- Chaque affirmation sur un outil (fonctionnalité, version, statut, adoption) est vérifiée à la
  source : documentation officielle, registre, dépôt. Sinon la reformuler ou la retirer.
- Liens testés (`curl -sL -o /dev/null -w '%{http_code}'`), y compris vers le dépôt `bivouac`.
- Cohérence avec les épisodes précédents (nombre d'ADR, tableau des épisodes, vocabulaire).
- Build Hugo complet sans `ERROR`.

## 5. Post LinkedIn (dans le plan éditorial)
- 1 200 à 1 800 caractères : accroche (un fait ou un incident réel), le récit en quelques lignes,
  une liste à émojis, « Ce que je dirais à un client : … », puis
  « L'article complet et le dépôt sont en commentaire 👇 », 3 hashtags maximum.
- Premier commentaire : URL de l'article + `https://github.com/Shinr0/bivouac`.
- Image à joindre : le `hero.png` de l'épisode. Aucun lien dans le texte du post.

## 6. Publication (l'utilisateur décide de chaque étape)
1. Pousser la branche et ouvrir la PR du site ; le job `build` doit être vert.
2. L'utilisateur relit et merge ; le workflow déploie sur `gh-pages`.
3. Vérifier l'URL en ligne (HTTP 200, balises `og:` de l'article, visuel servi).
   Article daté dans le futur : après la date, « Run workflow » sur le workflow de déploiement.
4. L'utilisateur publie le post avec l'image, puis le lien en premier commentaire.
5. Mettre à jour l'état de l'épisode dans le plan éditorial.
