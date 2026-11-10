package vm

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

const (
	apiVersion  = "compute.platform.local/v1alpha1"
	kind        = "VirtualMachine"
	serviceID   = "virtual-machine"
	managedBy   = "forge-bff"
	requesterAn = "forge.platform.local/requester"
)

// GVR cible le composite namespacé défini par la XRD.
var GVR = schema.GroupVersionResource{
	Group:    "compute.platform.local",
	Version:  "v1alpha1",
	Resource: "virtualmachines",
}

var (
	// ErrAlreadyExists est renvoyée quand une VM du même nom existe déjà.
	ErrAlreadyExists = errors.New("virtual machine already exists")
	// ErrNotFound est renvoyée quand la VM demandée n'existe pas.
	ErrNotFound = errors.New("virtual machine not found")
)

// Service manipule les composites VirtualMachine dans un namespace donné.
type Service struct {
	client    dynamic.Interface
	namespace string
}

func NewService(client dynamic.Interface, namespace string) *Service {
	return &Service{client: client, namespace: namespace}
}

func (s *Service) resource() dynamic.ResourceInterface {
	return s.client.Resource(GVR).Namespace(s.namespace)
}

// Create écrit le composite. Crossplane prend ensuite le relais : aucun YAML
// n'est généré ni stocké, l'objet Kubernetes est la source de vérité.
func (s *Service) Create(ctx context.Context, params Parameters, requester string) (Request, error) {
	normalized, err := Normalize(params)
	if err != nil {
		return Request{}, err
	}

	obj := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": apiVersion,
			"kind":       kind,
			"metadata": map[string]any{
				"name":      normalized.VMName,
				"namespace": s.namespace,
				"labels": map[string]any{
					"app.kubernetes.io/managed-by": managedBy,
					"forge.platform.local/service": serviceID,
				},
				"annotations": map[string]any{
					requesterAn: requester,
				},
			},
			"spec": map[string]any{
				"parameters": specParameters(normalized),
			},
		},
	}

	created, err := s.resource().Create(ctx, obj, metav1.CreateOptions{})
	switch {
	case apierrors.IsAlreadyExists(err):
		return Request{}, fmt.Errorf("%w: %s", ErrAlreadyExists, normalized.VMName)
	case err != nil:
		return Request{}, fmt.Errorf("create composite: %w", err)
	}

	return toRequest(created), nil
}

func (s *Service) List(ctx context.Context) ([]Request, error) {
	list, err := s.resource().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list composites: %w", err)
	}

	requests := make([]Request, 0, len(list.Items))
	for i := range list.Items {
		requests = append(requests, toRequest(&list.Items[i]))
	}

	// Plus récentes en tête, comme l'attend l'écran Suivi.
	sort.Slice(requests, func(i, j int) bool {
		return requests[i].CreatedAt.After(requests[j].CreatedAt)
	})
	return requests, nil
}

func (s *Service) Get(ctx context.Context, name string) (Request, error) {
	obj, err := s.resource().Get(ctx, name, metav1.GetOptions{})
	switch {
	case apierrors.IsNotFound(err):
		return Request{}, fmt.Errorf("%w: %s", ErrNotFound, name)
	case err != nil:
		return Request{}, fmt.Errorf("get composite: %w", err)
	}
	return toRequest(obj), nil
}

// Delete supprime le composite ; Crossplane détruit les ressources Scaleway.
func (s *Service) Delete(ctx context.Context, name string) error {
	err := s.resource().Delete(ctx, name, metav1.DeleteOptions{})
	switch {
	case apierrors.IsNotFound(err):
		return fmt.Errorf("%w: %s", ErrNotFound, name)
	case err != nil:
		return fmt.Errorf("delete composite: %w", err)
	}
	return nil
}

func specParameters(p Parameters) map[string]any {
	params := map[string]any{
		"vmName":         p.VMName,
		"machineType":    p.MachineType,
		"zone":           p.Zone,
		"imageFamily":    p.ImageFamily,
		"diskSizeGb":     int64(p.DiskSizeGb),
		"rootVolumeType": p.RootVolumeType,
		"dataDiskSizeGb": int64(p.DataDiskSizeGb),
		"staticIp":       p.StaticIP,
	}
	if len(p.Tags) > 0 {
		tags := make([]any, 0, len(p.Tags))
		for _, tag := range p.Tags {
			tags = append(tags, tag)
		}
		params["tags"] = tags
	}
	return params
}

func toRequest(obj *unstructured.Unstructured) Request {
	status, message := deriveStatus(obj)

	req := Request{
		ID:         string(obj.GetUID()),
		ServiceID:  serviceID,
		Parameters: readParameters(obj),
		Status:     status,
		Message:    message,
		PublicIP:   nestedString(obj, "status", "publicIp"),
		ServerID:   nestedString(obj, "status", "serverId"),
		State:      nestedString(obj, "status", "state"),
		Requester:  obj.GetAnnotations()[requesterAn],
		CreatedAt:  obj.GetCreationTimestamp().Time,
	}

	if req.ID == "" {
		req.ID = obj.GetName()
	}
	if req.CreatedAt.IsZero() {
		req.CreatedAt = time.Now().UTC()
	}
	if req.Parameters.VMName == "" {
		req.Parameters.VMName = obj.GetName()
	}
	return req
}

func readParameters(obj *unstructured.Unstructured) Parameters {
	raw, found, err := unstructured.NestedMap(obj.Object, "spec", "parameters")
	if err != nil || !found {
		return Parameters{}
	}

	params := Parameters{
		VMName:         asString(raw["vmName"]),
		MachineType:    asString(raw["machineType"]),
		Zone:           asString(raw["zone"]),
		ImageFamily:    asString(raw["imageFamily"]),
		RootVolumeType: asString(raw["rootVolumeType"]),
		DiskSizeGb:     asInt(raw["diskSizeGb"]),
		DataDiskSizeGb: asInt(raw["dataDiskSizeGb"]),
	}
	if v, ok := raw["staticIp"].(bool); ok {
		params.StaticIP = v
	}
	if items, ok := raw["tags"].([]any); ok {
		for _, item := range items {
			if tag := asString(item); tag != "" {
				params.Tags = append(params.Tags, tag)
			}
		}
	}
	return params
}

func nestedString(obj *unstructured.Unstructured, fields ...string) string {
	value, found, err := unstructured.NestedString(obj.Object, fields...)
	if err != nil || !found {
		return ""
	}
	return value
}

func asString(value any) string {
	s, _ := value.(string)
	return s
}

func asInt(value any) int {
	switch v := value.(type) {
	case int64:
		return int(v)
	case int32:
		return int(v)
	case int:
		return v
	case float64:
		return int(v)
	default:
		return 0
	}
}
