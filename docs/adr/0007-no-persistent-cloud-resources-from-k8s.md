# 0007. Ni PVC, ni LoadBalancer, ni IP statique dans la démo

- **Statut** : Accepté
- **Date** : 2026-11-10

## Contexte
Les ressources créées par des contrôleurs Kubernetes (disques PD derrière un PVC, forwarding rules
et NEG derrière un Service `LoadBalancer` ou un Ingress) ne sont pas connues de Terraform et
survivent souvent à `terraform destroy`. Ce sont les principales sources d'orphelins facturés.

## Décision
- Aucun `PersistentVolumeClaim` : stockage `emptyDir` (Prometheus, Postgres de LiteLLM), perte de
  données assumée à chaque `down`.
- Aucun Service `LoadBalancer`, Ingress ou Gateway exposé : accès par `kubectl port-forward`.
- Aucune IP statique.
- Ces règles sont **appliquées** par Kyverno (phase 4) et **contrôlées** par `make check-orphans`.
- Toute exception exige un ADR décrivant son nettoyage dans `down.sh`.

## Conséquences
- La démo repart d'un état vierge à chaque `up` (données de coût et clés LiteLLM incluses).
- Le mode public (IAP / OAuth2-proxy) sera une exception documentée.
