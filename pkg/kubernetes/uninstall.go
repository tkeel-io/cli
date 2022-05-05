package kubernetes

import (
	"fmt"
	"time"

	"github.com/tkeel-io/cli/fileutil"
	helm "helm.sh/helm/v3/pkg/action"
)

// UninstallPlatform removes tKeel from a Kubernetes cluster.
func UninstallPlatform(namespace string, timeout uint, debugMode bool) error {
	config, err := InitHelmConfig(namespace, getLog(debugMode))
	if err != nil {
		return err
	}

	uninstallClient := helm.NewUninstall(config)
	uninstallClient.Timeout = time.Duration(timeout) * time.Second
	_, err = uninstallClient.Run(tKeelReleaseName)
	if err != nil {
		return fmt.Errorf("helm uninstall err:%w", err)
	}
	_, err = uninstallClient.Run(fmt.Sprintf("tkeel-%s", tkeelCoreHelmChart))
	if err != nil {
		return fmt.Errorf("helm uninstall err:%w", err)
	}
	_, err = uninstallClient.Run(fmt.Sprintf("tkeel-%s", tkeelRudderHelmChart))
	if err != nil {
		return fmt.Errorf("helm uninstall err:%w", err)
	}
	_, err = uninstallClient.Run(tKeelMiddlewareReleaseName)
	if err != nil {
		return fmt.Errorf("helm uninstall err:%w", err)
	}
	return nil
}

// Uninstall removes tKeel's plugin from a Kubernetes cluster.
func Uninstall(pluginID string, debugMode bool) error {
	clientset, err := Client()
	if err != nil {
		return err
	}

	namespace, err := GetTKeelNamespace(clientset)
	if err != nil {
		return err
	}

	_, err = Unregister(pluginID)
	if err != nil {
		return err
	}

	_, err = HelmUninstall(namespace, pluginID)
	return err
}

func UninstallAllPlugin(namespace string, debugMode bool) error {
	list, err := InstalledList()
	if err != nil {
		return err
	}
	for _, plugin := range list {
		err = Uninstall(plugin.Name, debugMode)
		if err != nil {
			return err
		}
	}
	return nil
}

func CleanToken() {
	_, _ = fileutil.LocateAdminToken(fileutil.RewriteFlag())
}
