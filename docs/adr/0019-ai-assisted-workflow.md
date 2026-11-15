# 0019. Développer la plateforme avec Claude Code, sous garde-fous

- **Statut** : Accepté
- **Date** : 2026-11-15

## Contexte
La plateforme est construite par itérations avec un agent de code (Claude Code). Un agent sans
cadre dérive : il anticipe des phases, invente des versions, exécute des commandes facturées ou
destructrices, et relit tout le dépôt à chaque session, ce qui coûte des tokens. Un incident
réel l'a montré dès le début : la configuration `gcloud` globale de la machine pointait vers le
compte d'un client.

## Décision
1. **Cadre écrit** : `CLAUDE.md` (règles), `PLAN.md` (phases), un ADR par décision.
2. **Garde-fous exécutables** (`.claude/settings.json`, versionné) :
   - confirmation obligatoire pour `terraform apply/destroy/import/state`, `make up/down/bootstrap`,
     `git push`, `gh pr merge` ;
   - hook `PreToolUse` qui bloque tout `gcloud` hors de la config isolée du projet.
3. **Contexte économe** :
   - `CLAUDE.md` court et stable (relu à chaque requête, mis en cache) : règles et pointeurs, pas
     de contenu déductible du code ; seule la section utile de `PLAN.md` est lue ;
   - consignes chargées à la demande : skills de projet (`new-adr`, `close-phase`) et, à partir de
     la phase 2, un `CLAUDE.md` par dossier (`infra/`, `gitops/`, `portal/`) ;
   - lectures interdites (`permissions.deny`) : `.terraform/`, states, lock files, `*.tfvars`,
     `.env` — moins de tokens, et aucun state ni secret local dans le contexte ;
   - une session par tâche, `/clear` entre deux, mode plan en début de phase, modèle adapté.
4. **Revues** :
   - locale, avant push : `/code-review` de Claude Code, sans coût de CI ;
   - en CI : `claude-code-action`, uniquement à l'ouverture d'une PR prête ou sur le label
     `ai-review`, limitée au diff, 12 tours, modèle Sonnet, timeout 10 min, PR du dépôt seulement.
     Authentification par **Workload Identity Federation** (OIDC GitHub → jeton Anthropic court) :
     aucune clé API stockée, cohérent avec l'ADR 0004.
5. **Supply chain** : actions épinglées par SHA, paquet PyPI épinglé, et aucune version de moins
   de 7 jours adoptée.

## Alternatives écartées
- Revue IA à chaque push : coût proportionnel au nombre de commits, bruit.
- Clé `ANTHROPIC_API_KEY` en secret GitHub : secret longue durée à faire tourner.
- Graphe de code (Graphify) : essayé ; sur ce dépôt une requête renvoyait ~9 Ko contre ~13 Ko pour
  lire tout le code d'infra, et ses hooks ajoutaient un rappel à presque chaque commande. Retiré.
- Mémoire d'agent externe (pages écrites par l'agent) : doublon des ADR, outils encore jeunes.

## Conséquences
- Prérequis pour la revue CI : règle de fédération dans la console Claude et trois variables de
  dépôt (`ANTHROPIC_FEDERATION_RULE_ID`, `ANTHROPIC_ORGANIZATION_ID`,
  `ANTHROPIC_SERVICE_ACCOUNT_ID`) ; sans elles, le job est ignoré. Fixer une limite de dépense sur
  le workspace associé.
- À réévaluer quand le code grossit (phases 3 à 5) : mesurer le contexte consommé par session
  (`/context`, `/cost`) avant d'ajouter un outil.
