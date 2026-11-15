# Règles du projet
- Lire PLAN.md avant toute tâche ; travailler phase par phase (§ 13), ne pas anticiper la suivante.
- Vérifier la doc officielle pour toute syntaxe de CRD, chart ou API externe ; ne jamais deviner une version.
- Épingler toutes les versions (providers, charts, images). Jamais de `latest`.
- Aucun secret, clé JSON ou token dans le dépôt. Workload Identity uniquement.
- Aucun PVC, LoadBalancer ou IP statique sans ADR justifiant le nettoyage au destroy.
- GCP : uniquement le compte perso `GCP_ACCOUNT` via la config gcloud isolée du Makefile
  (`CLOUDSDK_CONFIG`). Le compte gcloud global de la machine est celui d'un client : ne jamais
  lancer `gcloud`/`terraform` hors `make` ou sans exporter `CLOUDSDK_CONFIG` et `GOOGLE_APPLICATION_CREDENTIALS`.
- Toute ressource GCP porte le label env=platform-ai-demo.
- Ne jamais exécuter `make up`, `terraform apply` ou `destroy` sans demande explicite de l'utilisateur.
- Après chaque phase : `make test` vert, README à jour, commit.
- Langue : code et identifiants en anglais, documentation en français.

# Repères
- État d'avancement : tableau « Avancement » du README.
- Ancienne version (Crossplane VM + portail Forge) : branche `archive/v1-crossplane-forge`, ne pas réutiliser sans demande.
- Actions GitHub épinglées par SHA de commit (commentaire avec le tag).
- Outils tiers (actions, paquets) : n'adopter qu'une version publiée depuis au moins 7 jours.
- Avant de proposer un push : `make test` vert, puis suggérer `/code-review` (revue locale gratuite).
- Garde-fous Claude Code dans `.claude/` : ne pas les affaiblir sans ADR. Carte du dépôt : `make graph`.

## graphify

This project has a knowledge graph at graphify-out/ with god nodes, community structure, and cross-file relationships.

Rules:
- When graphify-out/graph.json exists, run `graphify query "<question>"` before your first grep, search or source-file read for a codebase question, and put this rule in every subagent prompt that explores code. Read raw files only for lines the graph or an earlier read already located. Use `graphify path "<A>" "<B>"` for relationships and `graphify explain "<concept>"` for focused concepts. These return a scoped subgraph, usually much smaller than GRAPH_REPORT.md or raw grep output.
- If graphify-out/wiki/index.md exists, use it for broad navigation instead of raw source browsing.
- Read graphify-out/GRAPH_REPORT.md only for broad architecture review or when query/path/explain do not surface enough context.
- After modifying code, run `graphify update .` to keep the graph current (AST-only, no API cost).
- A committed graphify-out/graph.json is rewritten whole on every build, so a large diff to it in a pull request is expected. Do not flag it as churn or suggest skipping the graph rebuild or its hooks.
