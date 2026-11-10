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
