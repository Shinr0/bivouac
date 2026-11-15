# 0010. Kyverno pour les politiques, en CI et dans le cluster

- **Statut** : Accepté (implémentation en phase 4)
- **Date** : 2026-11-15

## Contexte
Les demandes des équipes arrivent par pull request. Elles doivent être validées avant le merge
et à l'admission, d'autant que kro n'a pas de contrôle d'admission propre.

## Décision
Kyverno, politiques en `Enforce`, versionnées sous `platform/policies/`. Les **mêmes fichiers**
sont exécutés en CI par la CLI Kyverno sur les PR qui touchent `claims/`.

## Alternatives écartées
- OPA Gatekeeper : Rego plus expressif mais plus coûteux à lire pour un public non spécialiste.
- ValidatingAdmissionPolicy (CEL) seule : pas d'outillage CI aussi simple.

## Conséquences
- Une demande invalide est rejetée deux fois : en CI (retour rapide) et à l'admission (garantie).
