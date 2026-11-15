# 0007. Aucune ressource cloud orpheline : règles et exception encadrée

- **Statut** : Accepté
- **Date** : 2026-11-15

## Contexte
Les ressources créées par des contrôleurs Kubernetes ne sont pas connues de Terraform et survivent
souvent à `terraform destroy` : disques derrière un PVC, forwarding rules et NEG derrière un
Service `LoadBalancer` ou un Ingress, ressources GCP créées depuis le cluster. Ce sont les
principales sources de facture fantôme.

## Décision
- Aucun `PersistentVolumeClaim` : stockage `emptyDir`, perte de données assumée à chaque `down`.
- Aucun Service `LoadBalancer`, Ingress ou Gateway exposé : accès par `kubectl port-forward`.
- Aucune IP statique.
- **Exception encadrée** : les ressources GCP des équipes créées par Config Connector (ADR 0009).
  Elles sont rapides à créer, quasi gratuites à l'arrêt, supprimées avec leur demande, et
  `make down` **attend leur suppression effective avant de détruire le cluster**.
- Règles appliquées par Kyverno (dont l'interdiction de l'annotation
  `cnrm.cloud.google.com/deletion-policy: abandon`) et contrôlées par `make check-orphans`, qui
  échoue s'il reste une ressource facturable portant le label du projet.
- Toute autre exception exige un ADR décrivant son nettoyage dans `down.sh`.

## Conséquences
- Le labo repart d'un état vierge à chaque `up` (données de coût comprises).
- Un mode d'exposition public (IAP / OAuth2-proxy) sera une exception documentée.
