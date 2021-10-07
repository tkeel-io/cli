package kubernetes

import (
	"time"

	helm "helm.sh/helm/v3/pkg/action"
)

// Uninstall removes TKeel from a Kubernetes cluster.
func Uninstall(namespace string, timeout uint) error {
	config, err := helmConfig(namespace)
	if err != nil {
		return err
	}

	uninstallClient := helm.NewUninstall(config)
	uninstallClient.Timeout = time.Duration(timeout) * time.Second
	_, err = uninstallClient.Run(tKeelReleaseName)
	return err
}
