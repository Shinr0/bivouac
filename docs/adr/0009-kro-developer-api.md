# 0009. kro comme API développeur

- **Statut** : Accepté (implémentation en phases 4–5)
- **Date** : 2026-11-10

## Contexte
Les équipes doivent demander des ressources par un objet simple (« claim »), sans connaître les
dizaines d'objets Kubernetes générés derrière.

## Décision
Les capacités de la plateforme sont exposées par des **ResourceGraphDefinitions kro**
(`TeamSpace`, `LLMAccess`), qui génèrent des CRD et composent des ressources Kubernetes natives.
Crossplane n'est envisagé que si une ressource GCP doit être exposée en self-service.

## Alternatives écartées
- Crossplane seul : puissant pour le cloud, mais lourd (providers, functions) pour composer des
  objets 100 % Kubernetes. La v1 de ce dépôt l'utilisait (branche `archive/v1-crossplane-forge`).
- Helm charts par équipe : pas d'API, pas de statut, pas de validation de schéma.
- Opérateur maison : coût de maintenance disproportionné.

## Conséquences
- Syntaxe du schéma et `apiVersion` à vérifier sur kro.run au moment de l'implémentation.
- Le portail lit le schéma OpenAPI des CRD générées pour construire ses formulaires.
