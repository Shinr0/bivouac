package vm

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func composite(conditions ...map[string]any) *unstructured.Unstructured {
	raw := make([]any, 0, len(conditions))
	for _, c := range conditions {
		raw = append(raw, c)
	}
	return &unstructured.Unstructured{
		Object: map[string]any{
			"status": map[string]any{"conditions": raw},
		},
	}
}

func TestDeriveStatus(t *testing.T) {
	cases := map[string]struct {
		obj  *unstructured.Unstructured
		want RequestStatus
	}{
		"sans condition, la demande vient d'être créée": {
			obj:  &unstructured.Unstructured{Object: map[string]any{}},
			want: StatusSubmitted,
		},
		"synced sans ready, provisioning en cours": {
			obj:  composite(map[string]any{"type": "Synced", "status": "True"}),
			want: StatusProvisioning,
		},
		"ready, la VM est disponible": {
			obj: composite(
				map[string]any{"type": "Synced", "status": "True"},
				map[string]any{"type": "Ready", "status": "True"},
			),
			want: StatusReady,
		},
		"erreur de reconciliation": {
			obj: composite(map[string]any{
				"type": "Synced", "status": "False",
				"reason": "ReconcileError", "message": "quota exceeded",
			}),
			want: StatusFailed,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, _ := deriveStatus(tc.obj)
			if got != tc.want {
				t.Errorf("deriveStatus() = %q, want %q", got, tc.want)
			}
		})
	}
}
