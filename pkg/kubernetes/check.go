package kubernetes

import dapr "github.com/dapr/cli/pkg/kubernetes"

type DaprStatus struct {
	Installed   bool   `json:"installed"`
	Version     string `json:"version"`
	Namespace   string `json:"namespace"`
	MTLSEnabled bool   `json:"mtls_enabled"`
	Error       error
}

func Check() *DaprStatus {
	result := &DaprStatus{
		Installed:   false,
		Version:     "",
		Namespace:   "",
		MTLSEnabled: false,
	}

	statusClient, err := dapr.NewStatusClient()
	if err != nil {
		result.Error = err
		return result
	}
	status, err := statusClient.Status()
	if err != nil {
		result.Error = err
		return result
	}
	if len(status) == 0 {
		result.Error = ErrDaprNotInstall
		result.Installed = false
		return result
	} else {
		result.Installed = true
		result.Namespace = status[0].Namespace
	}
	client, err := dapr.Client()
	if err != nil {
		result.Error = err
		return result
	}
	info, err := client.ServerVersion()
	if err != nil {
		result.Error = err
		return result
	}
	result.Version = info.String()
	enabled, err := dapr.IsMTLSEnabled()
	if err != nil {
		result.Error = err
		return result
	}
	result.MTLSEnabled = enabled
	return result
}
