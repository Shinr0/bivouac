# Règles du projet

IDP éphémère sur GCP construite par phases. **Phase en cours** : tableau « Avancement » du README.
Ne lire de `PLAN.md` que la section utile à la tâche (phases : § 13 ; workflow IA : § 16).

- Une phase à la fois ; ne pas anticiper la suivante.
- Syntaxe de CRD, chart ou API externe et versions : vérifier dans la doc ou le registre officiel,
  ne jamais deviner. Tout épingler (jamais `latest`), actions GitHub par SHA ; aucun outil tiers
  publié depuis moins de 7 jours.
- Aucun secret, clé JSON ou token dans le dépôt. Workload Identity / OIDC uniquement.
- GCP : uniquement via `make` (config gcloud isolée, compte perso). La config gcloud globale de la
  machine est celle d'un **client** ; un hook bloque tout `gcloud` direct.
- Aucun PVC, LoadBalancer ou IP statique sans ADR ; toute ressource GCP porte `env=platform-ai-demo`.
- Jamais `make up/down/bootstrap`, `terraform apply/destroy` ni `git push` sans demande explicite.
- Code et identifiants en anglais, documentation en français.

## Commandes
- `make test` : contrôles statiques, à lancer après chaque modification.
- `make whoami` : compte et projet GCP effectifs.
- Skills : `new-adr` (toute nouvelle décision), `close-phase` (fin de phase).
- Avant de proposer un push : `make test` vert, puis suggérer `/code-review`.

## Repères
- Garde-fous dans `.claude/` (permissions, hook gcloud) : ne pas les affaiblir sans ADR.
- Ancienne version : branche `archive/v1-crossplane-forge`, ne pas réutiliser sans demande.
