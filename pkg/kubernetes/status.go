package kubernetes

import (
	k8s "k8s.io/client-go/kubernetes"
)

type StatusClient struct {
	client k8s.Interface
}

// StatusOutput represents the status of a named tKeel resource.
type StatusOutput struct {
	Name         string `csv:"NAME"`
	Namespace    string `csv:"NAMESPACE"`
	Healthy      string `csv:"HEALTHY"`
	Status       string `csv:"STATUS"`
	PluginStatus string `csv:"PLUGINSTATUS"`
	Replicas     int    `csv:"REPLICAS"`
	Version      string `csv:"VERSION"`
	Age          string `csv:"AGE"`
	Created      string `csv:"CREATED"`
}

// NewStatusClient Create a new k8s client for status commands.
func NewStatusClient() (*StatusClient, error) {
	clientset, err := Client()
	if err != nil {
		return nil, err
	}
	return &StatusClient{
		client: clientset,
	}, nil
}
