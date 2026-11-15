# Règles du projet

Bivouac : laboratoire de platform engineering sur GCP, monté et démonté à la demande, construit par phases. **Phase en cours** : tableau « Avancement » du README.
Ne lire de `PLAN.md` que la section utile à la tâche (phases : § 13 ; workflow IA : § 16).

- Une phase à la fois ; ne pas anticiper la suivante.
- Syntaxe de CRD, chart ou API externe et versions : vérifier dans la doc ou le registre officiel,
  ne jamais deviner. Tout épingler (jamais `latest`), actions GitHub par SHA ; aucun outil tiers
  publié depuis moins de 7 jours.
- Aucun secret, clé JSON ou token dans le dépôt. Workload Identity / OIDC uniquement.
- GCP : uniquement via `make` (config gcloud isolée, compte perso). La config gcloud globale de la
  machine est celle d'un **client** ; un hook bloque tout `gcloud` direct.
- Aucun PVC, LoadBalancer ou IP statique sans ADR ; toute ressource GCP porte `env=bivouac`.
- Jamais `make up/down/bootstrap`, `terraform apply/destroy` ni `git push` sans demande explicite.
- Code et identifiants en anglais, documentation en français.
- Toute décision est historisée : un ADR n'est jamais modifié sur le fond ni supprimé ; une
  décision révisée = nouvel ADR qui remplace l'ancien (skill `new-adr`). Les ADR 0001–0015 sont la base.

## Commandes
- `make test` : contrôles statiques, à lancer après chaque modification.
- `make whoami` : compte et projet GCP effectifs.
- Skills : `new-adr` (toute nouvelle décision), `close-phase` (fin de phase), `write-episode`
  (article de blog + post LinkedIn d'une phase close ; site dans `~/perso/Shinr0.github.io`).
- Avant de proposer un push : `make test` vert, puis suggérer `/code-review`.

## Repères
- Garde-fous dans `.claude/` (permissions, hook gcloud) : ne pas les affaiblir sans ADR.
- Ancienne version : branche `archive/v1-crossplane-forge`, ne pas réutiliser sans demande.
