package vm

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

// deriveStatus traduit les conditions Crossplane en statut lisible par Forge.
//
// Crossplane expose deux conditions sur un composite :
//   - Synced : la réconciliation avec l'API cloud a réussi
//   - Ready  : toutes les ressources composées sont prêtes
func deriveStatus(obj *unstructured.Unstructured) (RequestStatus, string) {
	conditions, found, err := unstructured.NestedSlice(obj.Object, "status", "conditions")
	if err != nil || !found || len(conditions) == 0 {
		return StatusSubmitted, ""
	}

	var (
		syncedStatus, syncedReason, syncedMessage string
		readyStatus, readyReason, readyMessage    string
	)

	for _, raw := range conditions {
		cond, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		condType, _ := cond["type"].(string)
		status, _ := cond["status"].(string)
		reason, _ := cond["reason"].(string)
		message, _ := cond["message"].(string)

		switch condType {
		case "Synced":
			syncedStatus, syncedReason, syncedMessage = status, reason, message
		case "Ready":
			readyStatus, readyReason, readyMessage = status, reason, message
		}
	}

	// Une erreur de réconciliation est définitive tant qu'elle n'est pas corrigée.
	if syncedStatus == "False" {
		return StatusFailed, firstNonEmpty(syncedMessage, syncedReason, "réconciliation en échec")
	}
	if readyStatus == "True" {
		return StatusReady, ""
	}
	if syncedStatus == "True" {
		return StatusProvisioning, firstNonEmpty(readyMessage, readyReason)
	}
	return StatusSubmitted, ""
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
