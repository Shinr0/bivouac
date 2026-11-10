package vm

import "time"

// RequestStatus reprend le vocabulaire attendu par l'écran Suivi de Forge.
type RequestStatus string

const (
	StatusSubmitted    RequestStatus = "submitted"
	StatusProvisioning RequestStatus = "provisioning"
	StatusReady        RequestStatus = "ready"
	StatusFailed       RequestStatus = "failed"
)

// Parameters correspond au spec.parameters de la XRD VirtualMachine.
type Parameters struct {
	VMName         string   `json:"vmName"`
	MachineType    string   `json:"machineType"`
	Zone           string   `json:"zone"`
	ImageFamily    string   `json:"imageFamily"`
	DiskSizeGb     int      `json:"diskSizeGb"`
	RootVolumeType string   `json:"rootVolumeType"`
	DataDiskSizeGb int      `json:"dataDiskSizeGb"`
	StaticIP       bool     `json:"staticIp"`
	Tags           []string `json:"tags,omitempty"`
}

// Request est la vue exposée au frontend pour une VM demandée.
type Request struct {
	ID         string        `json:"id"`
	ServiceID  string        `json:"serviceId"`
	Parameters Parameters    `json:"parameters"`
	Status     RequestStatus `json:"status"`
	Message    string        `json:"message,omitempty"`
	PublicIP   string        `json:"publicIp,omitempty"`
	ServerID   string        `json:"serverId,omitempty"`
	State      string        `json:"state,omitempty"`
	Requester  string        `json:"requester,omitempty"`
	CreatedAt  time.Time     `json:"createdAt"`
}

// CatalogItem décrit une entrée du catalogue self-service.
type CatalogItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Available   bool   `json:"available"`
}

// FormOptions permet au frontend de construire son formulaire à partir des
// valeurs réellement acceptées par la XRD, sans les dupliquer côté client.
type FormOptions struct {
	MachineTypes    []string `json:"machineTypes"`
	Zones           []string `json:"zones"`
	ImageFamilies   []string `json:"imageFamilies"`
	RootVolumeTypes []string `json:"rootVolumeTypes"`
	MinDiskSizeGb   int      `json:"minDiskSizeGb"`
	MaxDiskSizeGb   int      `json:"maxDiskSizeGb"`
}
