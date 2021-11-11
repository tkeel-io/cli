package kubernetes

import (
	"fmt"
	"time"

	helm "helm.sh/helm/v3/pkg/action"
)

// Uninstall removes tKeel from a Kubernetes cluster.
func Uninstall(namespace string, timeout uint, debugMode bool) error {
	config, err := helmConfig(namespace, getLog(debugMode))
	if err != nil {
		return err
	}

	uninstallClient := helm.NewUninstall(config)
	uninstallClient.Timeout = time.Duration(timeout) * time.Second
	_, err = uninstallClient.Run(tKeelReleaseName)
	if err != nil {
		return fmt.Errorf("helm uninstall err:%w", err)
	}
	return nil
}
