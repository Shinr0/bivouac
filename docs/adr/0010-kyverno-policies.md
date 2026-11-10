# 0010. Kyverno pour les politiques, en CI et dans le cluster

- **Statut** : Accepté (implémentation en phase 4)
- **Date** : 2026-11-10

## Contexte
Les claims arrivent par PR, potentiellement rédigées par un agent IA. Ils doivent être validés
avant le merge et à l'admission.

## Décision
Kyverno, politiques en `Enforce`, versionnées sous `platform/policies/`. Les **mêmes fichiers**
sont exécutés en CI par la CLI Kyverno sur les PR qui touchent `claims/`.

## Alternatives écartées
- OPA Gatekeeper : Rego plus expressif mais plus coûteux à lire pour un public non spécialiste.
- ValidatingAdmissionPolicy (CEL) seule : pas d'outillage CI équivalent aussi simple, pas de mutation/génération.

## Conséquences
- Un claim invalide est rejeté deux fois : en CI (feedback rapide) et à l'admission (garantie).
