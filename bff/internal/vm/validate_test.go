package vm

import (
	"errors"
	"testing"
)

func TestNormalizeAppliesDefaults(t *testing.T) {
	got, err := Normalize(Parameters{VMName: "My-Dev-Server"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.VMName != "my-dev-server" {
		t.Errorf("vmName = %q, want lowercase", got.VMName)
	}
	if got.MachineType != "DEV1-S" {
		t.Errorf("machineType = %q, want DEV1-S", got.MachineType)
	}
	if got.Zone != "fr-par-1" {
		t.Errorf("zone = %q, want fr-par-1", got.Zone)
	}
	if got.DiskSizeGb != 20 {
		t.Errorf("diskSizeGb = %d, want 20", got.DiskSizeGb)
	}
	if got.RootVolumeType != "sbs_volume" {
		t.Errorf("rootVolumeType = %q, want sbs_volume", got.RootVolumeType)
	}
}

func TestNormalizeRejectsInvalidInput(t *testing.T) {
	cases := map[string]struct {
		params Parameters
		field  string
	}{
		"nom en majuscules interdit après trim": {
			params: Parameters{VMName: "1invalid"},
			field:  "vmName",
		},
		"type machine hors catalogue": {
			params: Parameters{VMName: "web", MachineType: "e2-micro"},
			field:  "machineType",
		},
		"zone GCP refusée": {
			params: Parameters{VMName: "web", Zone: "europe-west1-b"},
			field:  "zone",
		},
		"disque trop petit": {
			params: Parameters{VMName: "web", DiskSizeGb: 5},
			field:  "diskSizeGb",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Normalize(tc.params)
			if err == nil {
				t.Fatal("expected validation error")
			}

			var invalid *ValidationError
			if !errors.As(err, &invalid) {
				t.Fatalf("expected *ValidationError, got %T", err)
			}
			if _, ok := invalid.Fields[tc.field]; !ok {
				t.Errorf("expected error on field %q, got %v", tc.field, invalid.Fields)
			}
		})
	}
}
