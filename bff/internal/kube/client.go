// Package kube construit le client dynamique utilisé pour manipuler les
// composites Crossplane, dont le schéma n'est pas connu à la compilation.
package kube

import (
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewDynamicClient privilégie le ServiceAccount monté dans le pod et retombe
// sur un kubeconfig local pour le développement hors cluster.
func NewDynamicClient(kubeconfigPath string) (dynamic.Interface, error) {
	cfg, err := restConfig(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	client, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("build dynamic client: %w", err)
	}
	return client, nil
}

func restConfig(kubeconfigPath string) (*rest.Config, error) {
	if cfg, err := rest.InClusterConfig(); err == nil {
		return cfg, nil
	}

	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	if kubeconfigPath != "" {
		rules.ExplicitPath = kubeconfigPath
	}

	cfg, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		rules, &clientcmd.ConfigOverrides{},
	).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("load kubeconfig: %w", err)
	}
	return cfg, nil
}
