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
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/tkeel-io/cli/fileutil"
	kitconfig "github.com/tkeel-io/kit/config"
	"helm.sh/helm/v3/pkg/chart"
	"sigs.k8s.io/yaml"

	dapr "github.com/dapr/cli/pkg/kubernetes"
	"github.com/dapr/cli/pkg/print"
	"github.com/pkg/errors"
	helm "helm.sh/helm/v3/pkg/action"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	tKeelReleaseName               = "tkeel-platform"
	tKeelMiddlewareReleaseName     = "tkeel-middleware"
	tKeelPluginConfigHelmChart     = "tkeel-plugin-components"
	tKeelPluginMiddlewareHelmChart = "tkeel-middleware"
	tkeelKeelHelmChart             = "keel"
	tkeelRudderHelmChart           = "rudder"
	tkeelCoreHelmChart             = "core"
)

const (
	tkeelAdminHost  = "admin.tkeel.io"
	tkeelTenantHost = "tkeel.io"
	tkeelPort       = "30080"
)

var tKeelHelmRepo = "https://tkeel-io.github.io/helm-charts/"
var ErrDaprNotInstall = errors.New("dapr is not installed in your cluster")

var helmConf *helm.Configuration
var coreComponentChartNames = []string{tkeelRudderHelmChart, tkeelCoreHelmChart}

type InitConfiguration struct {
	Version       string
	Namespace     string
	Secret        string
	EnableMTLS    bool
	EnableHA      bool
	Args          []string
	Wait          bool
	Timeout       uint
	DebugMode     bool
	ConfigFile    string
	Password      string
	Repo          *kitconfig.Repo
	DefaultConfig bool
}

// middleware.yaml
type MiddlewareConfig struct {
	Queue           string `json:"queue" yaml:"queue" mapstructure:"queue,omitempty"`
	Database        string `json:"database" yaml:"database" mapstructure:"database,omitempty"`
	Cache           string `json:"cache" yaml:"cache" mapstructure:"cache,omitempty"`
	Search          string `json:"search" yaml:"search" mapstructure:"search,omitempty"`
	TSDB            string `json:"tsdb" yaml:"tsdb" mapstructure:"tsdb,omitempty"`
	ServiceRegistry string `json:"service_registry" yaml:"service_registry" mapstructure:"service_registry,omitempty"`
}

func (c *MiddlewareConfig) Empty() bool {
	return c.Queue == "" && c.Database == "" && c.Cache == "" && c.Search == "" && c.TSDB == "" && c.ServiceRegistry == ""
}

// Init deploys the tKeel operator using the supplied runtime version.
func Init(config InitConfiguration) (err error) {
	installConfig := &kitconfig.InstallConfig{}
	if config.ConfigFile != "" {
		installConfig, err = loadInstallConfig(config)
		if err != nil {
			return err
		}
	}
	installConfig.Namespace = config.Namespace
	if installConfig.Repo != nil {
		if installConfig.Repo.Url != "" {
			config.Repo.Url = installConfig.Repo.Url
		}
		if installConfig.Repo.Name != "" {
			config.Repo.Name = installConfig.Repo.Name
		}
	} else {
		installConfig.Repo = &kitconfig.Repo{
			Url:  config.Repo.Url,
			Name: config.Repo.Name,
		}
	}
	if installConfig.Host == nil {
		installConfig.Host = &kitconfig.Host{
			Admin:  tkeelAdminHost,
			Tenant: tkeelTenantHost,
		}
	}
	if installConfig.Port == "" {
		installConfig.Port = tkeelPort
	}

	tKeelHelmRepo = config.Repo.Url

	helmConf, err = helmConfig(config.Namespace, getLog(config.DebugMode))
	if err != nil {
		return err
	}

	keelChart, componentMiddleware, err := KeelChart(config, config.Password)
	if err != nil {
		return err
	}

	customizedMiddleware := make(map[string]string)
	middleware := installConfig.GetMiddleware()
	for _, value := range middleware {
		if !value.Customized {
			continue
		}
		items := strings.Split(value.Url, ",")
		if len(items) > 0 {
			uri, err := url.Parse(items[0])
			if err == nil {
				customizedMiddleware[uri.Scheme] = value.Url
			}
		}
	}

	middlewareChart, err := tKeelChart(config.Version, tKeelHelmRepo, tKeelPluginMiddlewareHelmChart, helmConf)
	if err != nil {
		return err
	}

	dependencies := make([]*chart.Chart, 0)
	for _, k := range middlewareChart.Dependencies() {
		if _, ok := customizedMiddleware[k.Name()]; !ok {
			dependencies = append(dependencies, k)
		}
	}
	middlewareChart.SetDependencies(dependencies...)

	// cover middleware config with custom middleware
	for component, middlewareConfig := range componentMiddleware {
		for service, iuri := range middlewareConfig {
			value, exist := middleware[service]
			if exist {
				middlewareConfig[service] = value.Url
			} else {
				if uri, ok := iuri.(string); ok {
					middleware[service] = &kitconfig.Value{
						Customized: false,
						Url:        uri,
					}
				}
			}
		}
		if component == tkeelKeelHelmChart {
			keelChart.Values["middleware"] = middlewareConfig
		} else if component == tkeelRudderHelmChart {
			keelChart.Values[component] = map[string]interface{}{
				"middleware": middlewareConfig,
				"tkeelRepo":  tKeelHelmRepo,
			}
		} else {
			keelChart.Values[component] = map[string]interface{}{
				"middleware": middlewareConfig,
			}
		}
	}

	installConfig.SetMiddleware(middleware)

	updateComponentsValues(middlewareChart, installConfig)
	updateIngresValues(middlewareChart, installConfig)
	err = updateConfigMap(middlewareChart, installConfig)
	if err != nil {
		return err
	}

	err = deploy(config, middlewareChart, keelChart)
	if err != nil {
		return err
	}

	err = dumpInstallConfig(config.ConfigFile, installConfig)
	if err != nil {
		return err
	}

	_, err = AdminLogin(config.Password)
	if err != nil {
		return err
	}

	err = AddRepo(config.Repo.Name, config.Repo.Url)
	if err != nil {
		return err
	}

	return nil
}

func deploy(config InitConfiguration, middlewareChart *chart.Chart, keelChart *chart.Chart) (err error) {
	msg := "Deploying the tKeel Platform to your cluster..."

	stopSpinning := print.Spinner(os.Stdout, msg)
	defer stopSpinning(print.Failure)

	err = installTkeel(config, middlewareChart, keelChart)
	if err != nil {
		return err
	}
	stopSpinning(print.Success)
	return err
}

func updateComponentsValues(middlewareChart *chart.Chart, config *kitconfig.InstallConfig) {
	redisUrl, err := url.Parse(config.Middleware.Cache.Url)
	if err != nil {
		return
	}
	kafkaUrl, err := url.Parse(config.Middleware.Queue.Url)
	if err != nil {
		return
	}
	components := make(map[string]interface{})
	state := make(map[string]interface{})
	pubsub := make(map[string]interface{})
	kafka := map[string]interface{}{
		"host": kafkaUrl.Host,
	}
	password, ok := redisUrl.User.Password()
	if !ok {
		password = ""
	}
	redis := map[string]interface{}{
		"host":     redisUrl.Host,
		"password": password,
	}
	state["redis"] = redis
	pubsub["kafka"] = kafka
	components["state"] = state
	components["pubsub"] = pubsub
	middlewareChart.Values["components"] = components
}

func updateIngresValues(middlewareChart *chart.Chart, config *kitconfig.InstallConfig) {
	ingress := make(map[string]interface{})
	ingress["host"] = map[string]interface{}{
		"admin":  config.Host.Admin,
		"tenant": config.Host.Tenant,
	}
	ingress["port"] = config.Port
	middlewareChart.Values["ingress"] = ingress
}

func updateConfigMap(middlewareChart *chart.Chart, config *kitconfig.InstallConfig) error {
	for _, f := range middlewareChart.Templates {
		if f.Name == "templates/configmap.yaml" {
			configmap := make(map[string]interface{})
			err := yaml.Unmarshal(f.Data, &configmap)
			if err != nil {
				return err
			}
			bConfig, err := yaml.Marshal(config)
			if err != nil {
				return err
			}
			if value, ok := configmap["metadata"].(map[string]interface{}); ok {
				value["namespace"] = config.Namespace
			}
			configmap["data"] = map[string]string{"config": string(bConfig)}
			data, err := yaml.Marshal(configmap)
			if err != nil {
				return err
			}
			f.Data = data
			break
		}
	}
	return nil
}

func dumpInstallConfig(configFile string, config *kitconfig.InstallConfig) error {
	if configFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configFile = filepath.Join(home, ".tkeel/config.yaml")
	}
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configFile, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

// load middleware config form file
func loadInstallConfig(config InitConfiguration) (*kitconfig.InstallConfig, error) {
	installConfig := &kitconfig.InstallConfig{}
	file, err := fileutil.LocateFile(fileutil.RewriteFlag(), config.ConfigFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &installConfig)
	if err != nil {
		return nil, err
	}
	return installConfig, nil
}

func createNamespace(namespace string) error {
	_, client, err := dapr.GetKubeConfigClient()
	if err != nil {
		return fmt.Errorf("can't connect to a Kubernetes cluster: %w", err)
	}

	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	// try to create the namespace if it doesn't exist. ok to ignore error.
	_, _ = client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	return nil
}

func createTempDir() (string, error) {
	dir, err := ioutil.TempDir("", "tkeel")
	if err != nil {
		return "", fmt.Errorf("error creating temp dir: %w", err)
	}
	return dir, nil
}

func locateChartFile(dirPath string) (string, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return "", fmt.Errorf("read dir err:%w", err)
	}
	return filepath.Join(dirPath, files[0].Name()), nil
}

func getLog(debugMode bool) helm.DebugLog {
	if debugMode {
		return func(format string, v ...interface{}) {
			print.InfoStatusEvent(os.Stdout, format, v...)
		}
	}
	return func(format string, v ...interface{}) {}
}
