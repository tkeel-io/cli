/*
Copyright 2021 The tKeel Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kubernetes

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tkeel-io/cli/pkg/print"
	helm "helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/helm/pkg/strvals"
)

func helmConfig(namespace string, log helm.DebugLog) (*helm.Configuration, error) {
	helmConf := helm.Configuration{}
	flags := &genericclioptions.ConfigFlags{
		Namespace: &namespace,
	}
	err := helmConf.Init(flags, namespace, "secret", log)
	if err != nil {
		err = fmt.Errorf("helm configuration init err:%w", err)
	}
	return &helmConf, err
}

func tKeelChart(version, repo, chartName string, config *helm.Configuration) (*chart.Chart, error) {
	pull := helm.NewPull()
	pull.RepoURL = repo
	pull.Settings = &cli.EnvSettings{}
	if version != "latest" {
		pull.Version = version
	}
	dir, err := createTempDir()
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	pull.DestDir = dir

	_, err = pull.Run(chartName)
	if err != nil {
		return nil, fmt.Errorf("err helm run pull:%w", err)
	}

	chartPath, err := locateChartFile(dir)
	if err != nil {
		return nil, fmt.Errorf("err locate chart file:%w", err)
	}
	c, err := loader.Load(chartPath)
	if err != nil {
		return nil, fmt.Errorf("err load chart file and parse:%w", err)
	}

	return c, nil
}

func chartValues(config InitConfiguration) (map[string]interface{}, error) {
	chartVals := map[string]interface{}{}
	globalVals := []string{
		fmt.Sprintf("global.ha.enabled=%t", config.EnableHA),
		fmt.Sprintf("global.mtls.enabled=%t", config.EnableMTLS),
	}
	globalVals = append(globalVals, config.Args...)

	for _, v := range globalVals {
		if err := strvals.ParseInto(v, chartVals); err != nil {
			return nil, fmt.Errorf("parse value err:%w", err)
		}
	}
	return chartVals, nil
}

func KeelChart(config InitConfiguration) (*chart.Chart, *chart.Chart, *chart.Chart, map[string]map[string]interface{}, error) {
	if config.Version != "latest" {
		if config.CoreVersion == "latest" {
			config.CoreVersion = config.Version
		}
		if config.RudderVersion == "latest" {
			config.RudderVersion = config.Version
		}
	}
	keelChart, err := tKeelChart(config.Version, tKeelHelmRepo, tkeelKeelHelmChart, helmConf)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if err = addDaprComponentChartDependency(config, helmConf, keelChart,
		tkeelKeelHelmChart); err != nil {
		return nil, nil, nil, nil, err
	}

	result := make(map[string]map[string]interface{})
	if value, exists := keelChart.Values["middleware"]; exists {
		chartConfig := make(map[string]interface{})
		if middlewares, ok := value.(map[string]interface{}); ok {
			for service, uri := range middlewares {
				chartConfig[service] = uri
			}
			result[tkeelKeelHelmChart] = chartConfig
		}
	}

	coreChart, err := tKeelChart(config.CoreVersion, tKeelHelmRepo, tkeelCoreHelmChart, helmConf)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if err = addDaprComponentChartDependency(config, helmConf, coreChart,
		tkeelCoreHelmChart); err != nil {
		return nil, nil, nil, nil, err
	}
	if value, ok := coreChart.Values["middleware"]; ok {
		chartConfig := make(map[string]interface{})
		if middlewares, ok := value.(map[string]interface{}); ok {
			for service, uri := range middlewares {
				chartConfig[service] = uri
			}
			result[tkeelCoreHelmChart] = chartConfig
		}
	}

	rudderChart, err := tKeelChart(config.RudderVersion, tKeelHelmRepo, tkeelRudderHelmChart, helmConf)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if err = addDaprComponentChartDependency(config, helmConf, rudderChart,
		tkeelRudderHelmChart); err != nil {
		return nil, nil, nil, nil, err
	}
	if value, ok := rudderChart.Values["middleware"]; ok {
		chartConfig := make(map[string]interface{})
		if middlewares, ok := value.(map[string]interface{}); ok {
			for service, uri := range middlewares {
				chartConfig[service] = uri
			}
			result[tkeelRudderHelmChart] = chartConfig
		}
	}
	return keelChart, coreChart, rudderChart, result, err
}

func installTkeel(config InitConfiguration, middlewareChart *chart.Chart, keelChart *chart.Chart, coreChart *chart.Chart, rudderChart *chart.Chart) error {
	err := createNamespace(config.Namespace)
	if err != nil {
		return err
	}
	var values map[string]interface{}
	values, err = chartValues(config)
	if err != nil {
		return err
	}

	installClient := helm.NewInstall(helmConf)
	installClient.Namespace = config.Namespace
	installClient.Timeout = time.Duration(config.Timeout) * time.Second

	installClient.ReleaseName = tKeelMiddlewareReleaseName
	installClient.Wait = true
	if _, err = installClient.Run(middlewareChart, values); err != nil {
		return fmt.Errorf("helm install err:%w", err)
	}
	print.InfoStatusEvent(os.Stdout, "install tKeel middleware done.")

	installClient.ReleaseName = tKeelReleaseName
	installClient.Wait = config.Wait
	if _, err = installClient.Run(keelChart, values); err != nil {
		return fmt.Errorf("helm install err:%w", err)
	}
	print.InfoStatusEvent(os.Stdout, "install tKeel platform  done.")

	installClient.ReleaseName = fmt.Sprintf("tkeel-%s", tkeelCoreHelmChart)
	installClient.Wait = config.Wait
	if _, err = installClient.Run(coreChart, values); err != nil {
		return fmt.Errorf("helm install err:%w", err)
	}
	print.InfoStatusEvent(os.Stdout, "install tKeel component <core> done.")

	installClient.ReleaseName = fmt.Sprintf("tkeel-%s", tkeelRudderHelmChart)
	installClient.Wait = config.Wait
	if _, err = installClient.Run(rudderChart, values); err != nil {
		return fmt.Errorf("helm install err:%w", err)
	}
	print.InfoStatusEvent(os.Stdout, "install tKeel component <rudder> done.")

	return nil
}

func upgradeTkeel(config InitConfiguration, middlewareChart *chart.Chart, keelChart *chart.Chart, coreChart *chart.Chart, rudderChart *chart.Chart) error {
	err := createNamespace(config.Namespace)
	if err != nil {
		return err
	}
	var values map[string]interface{}
	values, err = chartValues(config)
	if err != nil {
		return err
	}

	installClient := helm.NewUpgrade(helmConf)
	installClient.Namespace = config.Namespace
	installClient.Timeout = time.Duration(config.Timeout) * time.Second
	installClient.Wait = config.Wait

	if _, err = installClient.Run(tKeelMiddlewareReleaseName, middlewareChart, values); err != nil {
		return fmt.Errorf("helm install err:%w", err)
	}
	print.InfoStatusEvent(os.Stdout, "install tKeel middleware done.")

	if _, err = installClient.Run(tKeelReleaseName, keelChart, values); err != nil {
		return fmt.Errorf("helm upgrade err:%w", err)
	}
	print.InfoStatusEvent(os.Stdout, "upgrade tKeel platform  done.")

	if coreChart != nil {
		if _, err = installClient.Run(fmt.Sprintf("tkeel-%s", tkeelCoreHelmChart), coreChart, values); err != nil {
			return fmt.Errorf("helm install err:%w", err)
		}
		print.InfoStatusEvent(os.Stdout, "upgrade tKeel component <core> done.")
	}

	if rudderChart != nil {
		if _, err = installClient.Run(fmt.Sprintf("tkeel-%s", tkeelRudderHelmChart), rudderChart, values); err != nil {
			return fmt.Errorf("helm install err:%w", err)
		}
		print.InfoStatusEvent(os.Stdout, "upgrade tKeel component <rudder> done.")
	}

	return nil
}

func addDaprComponentChartDependency(config InitConfiguration, helmConf *helm.Configuration,
	root *chart.Chart, daprAppID string) error {
	componentChart, err := tKeelChart(config.Version, tKeelHelmRepo,
		tKeelPluginConfigHelmChart, helmConf)
	if err != nil {
		return err
	}

	root.Values["daprConfig"] = daprAppID
	componentChart.Values["pluginID"] = daprAppID
	componentChart.Values["secret"] = config.Secret
	root.AddDependency(componentChart)
	return nil
}

func HelmUninstall(namespace, pluginName string) (*release.UninstallReleaseResponse, error) {
	settings := cli.New()

	actionConfig := new(helm.Configuration)
	// You can pass an empty string instead of settings.Namespace() to list
	// all namespaces
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, fmt.Errorf("error init helm config: %w", err)
	}

	client := helm.NewUninstall(actionConfig)
	ret, err := client.Run(pluginName)
	if err != nil {
		return nil, fmt.Errorf("error helm uninstall: %w", err)
	}
	return ret, nil
}
