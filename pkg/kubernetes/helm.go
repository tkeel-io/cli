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
	"strings"
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

func KeelChart(config InitConfiguration, password string) (*chart.Chart, map[string]map[string]interface{}, error) {
	keelChart, err := tKeelChart(config.Version, tKeelHelmRepo, tkeelKeelHelmChart, helmConf)
	if err != nil {
		return nil, nil, err
	}

	if err = addDaprComponentChartDependency(config, helmConf, keelChart,
		tkeelKeelHelmChart); err != nil {
		return nil, nil, err
	}

	result := make(map[string]map[string]interface{})
	if value, ok := keelChart.Values["middleware"]; ok {
		chartConfig := make(map[string]interface{})
		middlewares := value.(map[string]interface{})
		for service, uri := range middlewares {
			chartConfig[service] = uri
		}
		result[tkeelKeelHelmChart] = chartConfig
	}

	for _, coreComponentName := range coreComponentChartNames {
		coreChart, err2 := tKeelChart(config.Version, tKeelHelmRepo, coreComponentName, helmConf)
		if err2 != nil {
			return nil, nil, err
		}
		if err2 = addDaprComponentChartDependency(config, helmConf, coreChart,
			coreComponentName); err2 != nil {
			return nil, nil, err
		}
		if coreComponentName == tkeelRudderHelmChart {
			coreChart.Values["adminPassword"] = password
		}
		if value, ok := coreChart.Values["middleware"]; ok {
			chartConfig := make(map[string]interface{})
			middlewares := value.(map[string]interface{})
			for service, uri := range middlewares {
				chartConfig[service] = uri
			}
			result[coreComponentName] = chartConfig
		}
		keelChart.AddDependency(coreChart)
	}
	return keelChart, result, err
}

func installTkeel(config InitConfiguration, middlewareChart *chart.Chart, keelChart *chart.Chart) error {
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

	print.InfoStatusEvent(os.Stdout, "install tKeel component<keel, %s> done.", strings.Join(coreComponentChartNames, ", "))

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

func InstallPlugin(config InitConfiguration, repo, chartName, releaseName, version string) error {
	clientset, err := Client()
	if err != nil {
		return err
	}

	namespace, err := GetTKeelNamespace(clientset)
	if err != nil {
		return err
	}

	config.Namespace = namespace
	helmConf, err := helmConfig(config.Namespace, getLog(config.DebugMode))
	if err != nil {
		return err
	}
	pluginChart, err := tKeelChart(version, repo, chartName, helmConf)
	if err != nil {
		return err
	}
	if err = addDaprComponentChartDependency(config, helmConf,
		pluginChart, releaseName); err != nil {
		return err
	}
	installClient := helm.NewInstall(helmConf)
	installClient.ReleaseName = releaseName
	installClient.Namespace = config.Namespace
	installClient.Wait = config.Wait
	installClient.Timeout = time.Duration(config.Timeout) * time.Second

	var values map[string]interface{}
	values, err = chartValues(config)
	if err != nil {
		return err
	}
	if _, err = installClient.Run(pluginChart, values); err != nil {
		return fmt.Errorf("helm install err:%w", err)
	}

	print.InfoStatusEvent(os.Stdout, "install tKeel plugin<%s> done.", chartName)
	return nil
}

func HelmList(namespace string) ([]*release.Release, error) {
	settings := cli.New()
	actionConfig := new(helm.Configuration)
	// You can pass an empty string instead of settings.Namespace() to list
	// all namespaces
	err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), log.Printf)
	if err != nil {
		return nil, fmt.Errorf("error init helm config: %w", err)
	}

	client := helm.NewList(actionConfig)
	// Only list deployed
	client.Deployed = true
	ret, err := client.Run()
	if err != nil {
		return nil, fmt.Errorf("error helm list: %w", err)
	}
	return ret, nil
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
