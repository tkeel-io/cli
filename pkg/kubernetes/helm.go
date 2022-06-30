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
	"os"
	"time"

	"github.com/tkeel-io/cli/pkg/print"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/helm/pkg/strvals"
)

func InitHelmConfig(namespace string, log action.DebugLog) (*action.Configuration, error) {
	helmConf = &action.Configuration{}
	flags := &genericclioptions.ConfigFlags{
		Namespace: &namespace,
	}
	err := helmConf.Init(flags, namespace, "secret", log)
	if err != nil {
		err = fmt.Errorf("helm configuration init err:%w", err)
	}
	return helmConf, err
}

func tKeelChart(version, repo, chartName string, config *action.Configuration) (*chart.Chart, error) {
	if version == "" {
		return nil, nil // nolint
	}
	pull := action.NewPull()
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

func ChartMiddleware(name, version string, config InitConfiguration) (*chart.Chart, MiddleConfig, error) {
	if version == "" {
		return nil, nil, nil
	}
	result := MiddleConfig{}
	c, err := tKeelChart(version, tKeelHelmRepo, name, helmConf)
	if err != nil {
		return nil, nil, err
	}

	if err = addDaprComponentChartDependency(config, helmConf, c, name); err != nil {
		return nil, nil, err
	}

	if value, exists := c.Values["middleware"]; exists {
		chartConfig := make(map[string]interface{})
		if middlewares, ok := value.(map[string]interface{}); ok {
			for service, uri := range middlewares {
				chartConfig[service] = uri
			}
			result[name] = chartConfig
		}
	}
	return c, result, nil
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

	installClient := action.NewInstall(helmConf)
	installClient.Namespace = config.Namespace
	installClient.Timeout = time.Duration(config.Timeout) * time.Second
	installClient.Wait = config.Wait

	if middlewareChart != nil {
		print.InfoStatusEvent(os.Stdout, "install tKeel middleware begin.")
		installClient.ReleaseName = tKeelMiddlewareReleaseName
		if _, err = installClient.Run(middlewareChart, values); err != nil {
			return fmt.Errorf("helm install err:%w", err)
		}
		print.InfoStatusEvent(os.Stdout, "install tKeel middleware done.")
	}

	if keelChart != nil {
		print.InfoStatusEvent(os.Stdout, "install tKeel keelChart begin.")
		installClient.ReleaseName = tKeelReleaseName
		if _, err = installClient.Run(keelChart, values); err != nil {
			return fmt.Errorf("helm install err:%w", err)
		}
		print.InfoStatusEvent(os.Stdout, "install tKeel platform  done.")
	}

	if coreChart != nil {
		print.InfoStatusEvent(os.Stdout, "install tKeel coreChart begin.")
		installClient.ReleaseName = fmt.Sprintf("tkeel-%s", tkeelCoreHelmChart)
		if _, err = installClient.Run(coreChart, values); err != nil {
			return fmt.Errorf("helm install err:%w", err)
		}
		print.InfoStatusEvent(os.Stdout, "install tKeel component <core> done.")
	}

	if rudderChart != nil {
		print.InfoStatusEvent(os.Stdout, "install tKeel rudderChart begin.")
		installClient.ReleaseName = fmt.Sprintf("tkeel-%s", tkeelRudderHelmChart)
		if _, err = installClient.Run(rudderChart, values); err != nil {
			return fmt.Errorf("helm install err:%w", err)
		}
		print.InfoStatusEvent(os.Stdout, "install tKeel component <rudder> done.")
	}
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

	upgradeClient := action.NewUpgrade(helmConf)
	upgradeClient.Namespace = config.Namespace
	upgradeClient.Timeout = time.Duration(config.Timeout) * time.Second
	upgradeClient.Wait = config.Wait

	if middlewareChart != nil {
		if _, err = upgradeClient.Run(tKeelMiddlewareReleaseName, middlewareChart, values); err != nil {
			return fmt.Errorf("helm install err:%w", err)
		}
		print.InfoStatusEvent(os.Stdout, "upgrade tKeel middleware done.")
	}

	if keelChart != nil {
		if _, err = upgradeClient.Run(tKeelReleaseName, keelChart, values); err != nil {
			return fmt.Errorf("helm upgrade err:%w", err)
		}
		print.InfoStatusEvent(os.Stdout, "upgrade tKeel platform done.")
	}

	if coreChart != nil {
		if _, err = upgradeClient.Run(fmt.Sprintf("tkeel-%s", tkeelCoreHelmChart), coreChart, values); err != nil {
			return fmt.Errorf("helm install err:%w", err)
		}
		print.InfoStatusEvent(os.Stdout, "upgrade tKeel component <core> done.")
	}

	if rudderChart != nil {
		if _, err = upgradeClient.Run(fmt.Sprintf("tkeel-%s", tkeelRudderHelmChart), rudderChart, values); err != nil {
			return fmt.Errorf("helm install err:%w", err)
		}
		print.InfoStatusEvent(os.Stdout, "upgrade tKeel component <rudder> done.")
	}

	return nil
}

func addDaprComponentChartDependency(config InitConfiguration, helmConf *action.Configuration,
	root *chart.Chart, daprAppID string) error {
	version := config.KeelVersion
	if version == "" {
		version = "latest"
	}
	componentChart, err := tKeelChart(version, tKeelHelmRepo, tKeelPluginConfigHelmChart, helmConf)
	if err != nil {
		return err
	}

	root.Values["daprConfig"] = daprAppID
	componentChart.Values["pluginID"] = daprAppID
	componentChart.Values["secret"] = config.Secret
	root.AddDependency(componentChart)
	return nil
}
