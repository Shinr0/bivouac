# 0014. Développer avec Claude Code, sous garde-fous

- **Statut** : Accepté
- **Date** : 2026-11-15

## Contexte
Bivouac est construit par itérations avec un agent de code (Claude Code). Sans cadre, un agent
dérive : il anticipe des phases, invente des versions, exécute des commandes facturées ou
destructrices, relit tout le dépôt à chaque session. La configuration `gcloud` globale du poste de
travail pointe par ailleurs vers le compte d'un client.

## Décision
1. **Cadre écrit** : `CLAUDE.md` court et stable (règles, commandes, pointeurs), `PLAN.md`
   (phases), un ADR par décision. Seule la section utile de `PLAN.md` est lue.
2. **Garde-fous exécutables** (`.claude/settings.json`, versionné) :
   - confirmation obligatoire pour `terraform apply/destroy/import/state`, `make up/down/bootstrap`,
     `git push`, `gh pr merge` ;
   - hook `PreToolUse` qui bloque tout `gcloud` hors de la configuration isolée du projet, testé
     par `make test` ;
   - lectures interdites : `.terraform/`, states, lock files, `*.tfvars`, `.env`.
3. **Contexte économe** : skills de projet chargés à la demande (`new-adr`, `close-phase`), un
   `CLAUDE.md` par dossier à partir de la phase 2, une session par tâche.
4. **Revues** : `/code-review` en local avant push ; en CI, `claude-code-action` à l'ouverture d'une
   PR prête ou sur le label `ai-review`, limitée au diff, 12 tours, modèle Sonnet, 10 min, PR du
   dépôt uniquement, authentifiée par fédération d'identité OIDC (aucune clé, ADR 0004).
5. **Supply chain** : actions épinglées par SHA, paquets par version exacte, aucune version publiée
   depuis moins de 7 jours.

## Alternatives écartées
- Revue IA à chaque push : coût proportionnel au nombre de commits, bruit.
- Graphe de code ou mémoire d'agent externes : gain négligeable mesuré sur un dépôt de cette
  taille, doublon des ADR ; à réévaluer quand le code grossit, mesures (`/context`, `/cost`) à l'appui.

## Conséquences
- Prérequis de la revue CI : règle de fédération dans la console Claude, trois variables de dépôt
  (`ANTHROPIC_FEDERATION_RULE_ID`, `ANTHROPIC_ORGANIZATION_ID`, `ANTHROPIC_SERVICE_ACCOUNT_ID`) et
  une limite de dépense sur le workspace ; sans elles, le job est ignoré.
