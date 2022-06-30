package kubernetes

import (
	dapr "github.com/dapr/cli/pkg/kubernetes"
	helm "helm.sh/helm/v3/pkg/action"
)

type DaprStatus struct {
	Installed   bool   `json:"installed"`
	Version     string `json:"version"`
	Namespace   string `json:"namespace"`
	MTLSEnabled bool   `json:"mtls_enabled"`
}

type TKeelStatus struct {
	Installed bool   `json:"installed"`
	Version   string `json:"version"`
	Namespace string `json:"namespace"`
}

func CheckDapr() (*DaprStatus, error) {
	result := &DaprStatus{
		Installed:   false,
		Version:     "",
		Namespace:   "",
		MTLSEnabled: false,
	}

	statusClient, err := dapr.NewStatusClient()
	if err != nil {
		return nil, err
	}
	status, err := statusClient.Status()
	if err != nil {
		return nil, err
	}
	if len(status) == 0 {
		return nil, ErrDaprNotInstall
	}
	result.Installed = true
	result.Namespace = status[0].Namespace
	for _, item := range status {
		if item.Name == "dapr-sentry" {
			result.Version = item.Version
			break
		}
	}
	enabled, err := dapr.IsMTLSEnabled()
	if err != nil {
		return nil, err
	}
	result.MTLSEnabled = enabled
	return result, nil
}

func CheckTKeel() (*TKeelStatus, error) {
	result := &TKeelStatus{
		Installed: false,
		Version:   "",
		Namespace: "",
	}
	helmConf, err := InitHelmConfig("", getLog(false))
	if err != nil {
		return nil, err
	}
	list := helm.NewList(helmConf)
	list.Filter = "tkeel-platform"
	releases, err := list.Run()
	if err != nil {
		return nil, err
	}
	if len(releases) == 0 {
		return nil, ErrTKeelNotInstall
	}
	result.Installed = true
	result.Namespace = releases[0].Namespace
	result.Version = releases[0].Chart.Metadata.Version
	return result, nil
}
