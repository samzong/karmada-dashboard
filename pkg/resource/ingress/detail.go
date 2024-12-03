package ingress

import (
	"context"
	"log"

	v1 "k8s.io/api/networking/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "k8s.io/client-go/kubernetes"
)

// IngressDetail API resource provides mechanisms to inject containers with configuration data while keeping
// containers agnostic of Kubernetes
type IngressDetail struct {
	// Extends list item structure.
	Ingress `json:",inline"`

	// Spec is the desired state of the Ingress.
	Spec v1.IngressSpec `json:"spec"`

	// Status is the current state of the Ingress.
	Status v1.IngressStatus `json:"status"`

	// List of non-critical errors, that occurred during resource retrieval.
	Errors []error `json:"errors"`
}

// GetIngressDetail returns detailed information about an ingress
func GetIngressDetail(client client.Interface, namespace, name string) (*IngressDetail, error) {
	log.Printf("Getting details of %s ingress in %s namespace", name, namespace)

	rawIngress, err := client.NetworkingV1().Ingresses(namespace).Get(context.TODO(), name, metaV1.GetOptions{})

	if err != nil {
		return nil, err
	}

	return getIngressDetail(rawIngress), nil
}

func getIngressDetail(i *v1.Ingress) *IngressDetail {
	return &IngressDetail{
		Ingress: toIngress(i),
		Spec:    i.Spec,
		Status:  i.Status,
	}
}