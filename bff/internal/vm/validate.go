package vm

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// Listes alignées sur les enum et bornes de la XRD VirtualMachine.
var (
	MachineTypes = []string{
		"DEV1-S", "DEV1-M", "DEV1-L",
		"PLAY2-PICO", "PLAY2-NANO", "PLAY2-MICRO",
		"PRO2-XXS", "PRO2-XS",
	}
	Zones           = []string{"fr-par-1", "fr-par-2", "nl-ams-1", "nl-ams-2", "pl-waw-1"}
	ImageFamilies   = []string{"ubuntu_jammy", "ubuntu_noble", "debian_bookworm", "rockylinux_9"}
	RootVolumeTypes = []string{"sbs_volume", "l_ssd"}
)

const (
	MinDiskSizeGb = 10
	MaxDiskSizeGb = 1000
	MaxDataDiskGb = 10000
)

var nameRE = regexp.MustCompile(`^[a-z]([a-z0-9-]{1,61}[a-z0-9])?$`)

// ValidationError agrège les erreurs de saisie pour un retour 400 exploitable.
type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	parts := make([]string, 0, len(e.Fields))
	for field, msg := range e.Fields {
		parts = append(parts, fmt.Sprintf("%s: %s", field, msg))
	}
	slices.Sort(parts)
	return strings.Join(parts, "; ")
}

// Normalize applique les valeurs par défaut de la XRD puis valide l'entrée.
func Normalize(p Parameters) (Parameters, error) {
	out := p
	out.VMName = strings.TrimSpace(strings.ToLower(p.VMName))

	if out.MachineType == "" {
		out.MachineType = "DEV1-S"
	}
	if out.Zone == "" {
		out.Zone = "fr-par-1"
	}
	if out.ImageFamily == "" {
		out.ImageFamily = "ubuntu_jammy"
	}
	if out.RootVolumeType == "" {
		out.RootVolumeType = "sbs_volume"
	}
	if out.DiskSizeGb == 0 {
		out.DiskSizeGb = 20
	}

	fields := map[string]string{}

	if !nameRE.MatchString(out.VMName) {
		fields["vmName"] = "minuscules, chiffres et tirets ; doit commencer par une lettre (2 à 63 caractères)"
	}
	if !slices.Contains(MachineTypes, out.MachineType) {
		fields["machineType"] = "type non supporté par le catalogue"
	}
	if !slices.Contains(Zones, out.Zone) {
		fields["zone"] = "zone non supportée"
	}
	if !slices.Contains(ImageFamilies, out.ImageFamily) {
		fields["imageFamily"] = "image non supportée"
	}
	if !slices.Contains(RootVolumeTypes, out.RootVolumeType) {
		fields["rootVolumeType"] = "type de volume non supporté"
	}
	if out.DiskSizeGb < MinDiskSizeGb || out.DiskSizeGb > MaxDiskSizeGb {
		fields["diskSizeGb"] = fmt.Sprintf("doit être compris entre %d et %d", MinDiskSizeGb, MaxDiskSizeGb)
	}
	if out.DataDiskSizeGb < 0 || out.DataDiskSizeGb > MaxDataDiskGb {
		fields["dataDiskSizeGb"] = fmt.Sprintf("doit être compris entre 0 et %d", MaxDataDiskGb)
	}

	if len(fields) > 0 {
		return Parameters{}, &ValidationError{Fields: fields}
	}
	return out, nil
}

func Options() FormOptions {
	return FormOptions{
		MachineTypes:    MachineTypes,
		Zones:           Zones,
		ImageFamilies:   ImageFamilies,
		RootVolumeTypes: RootVolumeTypes,
		MinDiskSizeGb:   MinDiskSizeGb,
		MaxDiskSizeGb:   MaxDiskSizeGb,
	}
}
