package kubernetes

import (
	"fmt"
	"os"
	"time"

	"github.com/tkeel-io/cli/fileutil"
	"github.com/tkeel-io/cli/pkg/print"
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

func UninstallAllPlugin() error {
	tenantList, err := TenantList()
	if err != nil {
		return err
	}
	pluginList, err := InstalledPlugin()
	if err != nil {
		return err
	}
	for _, plugin := range pluginList {
		// TODO 为所有租户禁用插件
		print.InfoStatusEvent(os.Stdout, "Removing plugin %s ...", plugin.Name)
		for _, tenant := range tenantList {
			err = DisablePlugin(plugin.Name, tenant.ID)
			if err != nil {
				return err
			}
		}
		err = UninstallPlugin(plugin.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func CleanToken() {
	_, _ = fileutil.LocateAdminToken(fileutil.RewriteFlag())
}
