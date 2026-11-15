---
name: close-phase
description: Close a PLAN.md phase in this repo - check the acceptance criterion, run tests, update README progress and docs, prepare the commit. Use when the user says a phase is done or asks to wrap one up.
---

# Clôturer une phase

1. Relire le critère d'acceptation de la phase dans `PLAN.md` § 13. S'il exige une action sur GCP
   (`make up`, `make down`…), demander à l'utilisateur de la lancer ou de l'autoriser ; ne jamais
   la lancer d'office.
2. `make test` doit être vert. Sinon, corriger d'abord.
3. Vérifier qu'aucune règle de `CLAUDE.md` n'est enfreinte dans le diff de la phase : versions
   épinglées, aucun secret, labels GCP, pas de PVC/LoadBalancer/IP statique, pas de travail
   d'une phase suivante.
4. Mettre à jour : tableau « Avancement » du README (✅ / ⏳ suivante), sections du README
   touchées, ADR créés pendant la phase (skill `new-adr`).
5. Noter pour l'article de la série les mesures réelles de la phase (durées, coûts, captures)
   dans `content/plan-editorial.md`.
6. Proposer le message de commit (`feat(phase-N): …`) et suggérer `/code-review` avant push.
