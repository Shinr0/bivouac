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
3. **Cartographie** : [Graphify](https://github.com/Graphify-Labs/graphify) construit un graphe du
   dépôt (AST tree-sitter, local, sans LLM ni télémétrie) ; son hook oriente l'agent vers
   `graphify query` avant les lectures brutes. Version épinglée (`make graph-install`), graphe
   régénéré par `make graph`, non commité.
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
- Graphe commité : diff volumineux à chaque build, sans valeur pour un dépôt de cette taille.

## Conséquences
- Prérequis pour la revue CI : règle de fédération dans la console Claude et trois variables de
  dépôt (`ANTHROPIC_FEDERATION_RULE_ID`, `ANTHROPIC_ORGANIZATION_ID`,
  `ANTHROPIC_SERVICE_ACCOUNT_ID`) ; sans elles, le job est ignoré. Fixer une limite de dépense sur
  le workspace associé.
- Sur un petit dépôt, le gain de Graphify est faible ; il croît avec le code (phases 3 à 6). À mesurer.
