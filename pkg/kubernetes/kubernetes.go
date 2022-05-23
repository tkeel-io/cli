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

	"github.com/tkeel-io/cli/fileutil"
	"github.com/tkeel-io/cli/pkg/utils"
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
	tKeelReleaseName           = "tkeel-platform"
	tKeelMiddlewareReleaseName = "tkeel-middleware"
	tKeelPluginConfigHelmChart = "tkeel-plugin-components"
	tKeelMiddlewareHelmChart   = "tkeel-middleware"
	tkeelKeelHelmChart         = "keel"
	tkeelRudderHelmChart       = "rudder"
	tkeelCoreHelmChart         = "core"
)

const (
	tkeelAdminHost  = "admin.tkeel.io"
	tkeelTenantHost = "tkeel.io"
	tkeelPort       = "30080"
)

var tKeelHelmRepo = "https://tkeel-io.github.io/helm-charts/"
var ErrDaprNotInstall = errors.New("dapr is not installed in your cluster")

var helmConf *helm.Configuration

var defaultPlugins = []string{
	"tkeel/console-portal-admin",
	"tkeel/console-portal-tenant",
	"tkeel/console-plugin-admin-plugins",
}

const (
	Always       = "Always"
	IfNotPresent = "IfNotPresent"
)

type InitConfiguration struct {
	Version           string
	KeelVersion       string
	CoreVersion       string
	RudderVersion     string
	MiddlewareVersion string
	Namespace         string
	Secret            string
	EnableMTLS        bool
	EnableHA          bool
	DaprVersion       string
	Args              []string
	Wait              bool
	Timeout           uint
	DebugMode         bool
	ConfigFile        string
	Password          string
	Repo              *kitconfig.Repo
	DefaultConfig     bool
	ImagePolicy       string
}

type MiddleConfig map[string]map[string]interface{}

func (mc MiddleConfig) Update(data MiddleConfig) {
	if data == nil {
		return
	}
	for key1, value1 := range data {
		value, ok := mc[key1]
		if !ok {
			value = make(map[string]interface{})
		}
		for key2, value2 := range value1 {
			value[key2] = value2
		}
		mc[key1] = value
	}
}

// Init deploys the tKeel operator using the supplied runtime version.
func Init(config InitConfiguration) error {
	installConfig, err := loadInstallConfig(config)
	if err != nil {
		return err
	}
	tKeelHelmRepo = config.Repo.Url

	helmConf, err = InitHelmConfig(config.Namespace, getLog(config.DebugMode))
	if err != nil {
		return err
	}

	tKeelMiddle := MiddleConfig{}

	keelChart, middleConfig, err := ChartMiddleware(tkeelKeelHelmChart, config.KeelVersion, config)
	if err != nil {
		return err
	}
	tKeelMiddle.Update(middleConfig)

	coreChart, middleConfig, err := ChartMiddleware(tkeelCoreHelmChart, config.CoreVersion, config)
	if err != nil {
		return err
	}
	tKeelMiddle.Update(middleConfig)

	rudderChart, middleConfig, err := ChartMiddleware(tkeelRudderHelmChart, config.RudderVersion, config)
	if err != nil {
		return err
	}
	tKeelMiddle.Update(middleConfig)

	middlewareChart, err := tKeelChart(config.MiddlewareVersion, tKeelHelmRepo, tKeelMiddlewareHelmChart, helmConf)
	if err != nil {
		return err
	}

	updateMiddlewareDependencies(middlewareChart, installConfig)

	charts := map[string]*chart.Chart{
		tkeelKeelHelmChart:   keelChart,
		tkeelCoreHelmChart:   coreChart,
		tkeelRudderHelmChart: rudderChart,
	}
	updateImagePolicy(charts, installConfig.ImagePolicy)

	middleware := updateChartsValue(charts, tKeelMiddle, installConfig, config)
	installConfig.SetMiddleware(middleware)

	updateComponentsValues(middlewareChart, installConfig)
	updateIngresValues(middlewareChart, installConfig)
	err = updateConfigMap(middlewareChart, installConfig)
	if err != nil {
		return err
	}

	err = deploy(config, middlewareChart, keelChart, coreChart, rudderChart)
	if err != nil {
		return err
	}

	err = dumpInstallConfig(config.ConfigFile, installConfig)
	if err != nil {
		return err
	}

	err = afterDeploy(config)
	if err != nil {
		return err
	}

	installPlugins(config, installConfig.Plugins, "")

	return nil
}

func deploy(config InitConfiguration, middlewareChart *chart.Chart, keelChart *chart.Chart, coreChart *chart.Chart, rudderChart *chart.Chart) (err error) {
	msg := "Deploying the tKeel Platform to your cluster..."

	stopSpinning := print.Spinner(os.Stdout, msg)
	defer stopSpinning(print.Failure)

	err = installTkeel(config, middlewareChart, keelChart, coreChart, rudderChart)
	if err != nil {
		return err
	}
	stopSpinning(print.Success)
	return err
}

func afterDeploy(config InitConfiguration) error {
	_, err := AdminLogin(config.Password)
	if err != nil {
		return err
	}

	err = AddRepo(config.Repo.Name, config.Repo.Url)
	if err != nil {
		if err.Error() == "response error: REPO已存在" {
			return nil
		}
		return err
	}
	return nil
}

func updateMiddlewareDependencies(middlewareChart *chart.Chart, config *kitconfig.InstallConfig) {
	if middlewareChart == nil {
		return
	}
	customizedMiddleware := make(map[string]string)
	middleware := config.GetMiddleware()
	for _, value := range middleware {
		if !value.Customized {
			continue
		}
		var uri *url.URL
		uri, err := url.Parse(value.Url)
		if err == nil {
			customizedMiddleware[uri.Scheme] = value.Url
		}
	}

	dependencies := make([]*chart.Chart, 0)
	for _, k := range middlewareChart.Dependencies() {
		if _, ok := customizedMiddleware[k.Name()]; !ok {
			dependencies = append(dependencies, k)
		}
	}
	middlewareChart.SetDependencies(dependencies...)
}

func updateComponentsValues(middlewareChart *chart.Chart, config *kitconfig.InstallConfig) {
	if middlewareChart == nil {
		return
	}
	CacheURL, err := url.Parse(config.Middleware.Cache.Url)
	if err != nil {
		return
	}
	QueueURL, err := url.Parse(config.Middleware.Queue.Url)
	if err != nil {
		return
	}

	cachePassword, _ := CacheURL.User.Password()
	queuePassword, _ := CacheURL.User.Password()

	components := map[string]interface{}{
		"state": map[string]interface{}{
			CacheURL.Scheme: map[string]interface{}{
				"host":     CacheURL.Host,
				"password": cachePassword,
			},
		},
		"pubsub": map[string]interface{}{
			QueueURL.Scheme: map[string]interface{}{
				"host":     QueueURL.Host,
				"password": queuePassword,
			},
		},
	}

	middlewareChart.Values["components"] = components
}

func updateImagePolicy(charts map[string]*chart.Chart, policy string) {
	for _, v := range charts {
		if v == nil {
			continue
		}
		if image, ok := v.Values["image"]; ok {
			if temp, ok := image.(map[string]interface{}); ok {
				if policy == Always {
					temp["pullPolicy"] = Always
				} else {
					temp["pullPolicy"] = IfNotPresent
				}
			}
		}
	}
}

func updateChartsValue(charts map[string]*chart.Chart, tKeelMiddle MiddleConfig, installConfig *kitconfig.InstallConfig, config InitConfiguration) map[string]*kitconfig.Value {
	middleware := installConfig.GetMiddleware()

	for component, middlewareConfig := range tKeelMiddle {
		for service, iuri := range middlewareConfig {
			if value, exist := middleware[service]; exist {
				middlewareConfig[service] = value.Url
			} else if uri, ok := iuri.(string); ok {
				middleware[service] = &kitconfig.Value{
					Customized: false,
					Url:        uri,
				}
			}
		}
		c := charts[component]
		if c == nil {
			continue
		}
		c.Values["middleware"] = middlewareConfig
		if component == tkeelRudderHelmChart {
			c.Values["tkeelRepo"] = tKeelHelmRepo
			c.Values["adminPassword"] = config.Password
		}
	}
	for _, v := range charts {
		if v == nil {
			continue
		}
		v.Values["daprVersion"] = config.DaprVersion
	}
	return middleware
}

func updateIngresValues(middlewareChart *chart.Chart, config *kitconfig.InstallConfig) {
	if middlewareChart == nil {
		return
	}
	ingress := make(map[string]interface{})
	ingress["host"] = map[string]interface{}{
		"admin":  config.Host.Admin,
		"tenant": config.Host.Tenant,
	}
	ingress["port"] = config.Port
	middlewareChart.Values["ingress"] = ingress
}

func updateConfigMap(middlewareChart *chart.Chart, config *kitconfig.InstallConfig) error {
	if middlewareChart == nil {
		return nil
	}
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
			configmap["data"] = map[string]string{
				"config":      string(bConfig),
				"TENANT_HOST": config.Host.Tenant,
				"ADMIN_HOST":  config.Host.Admin,
			}
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
			return errors.Wrap(err, "dump install config error")
		}
		configFile = filepath.Join(home, ".tkeel/config.yaml")
	}
	data, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "marshal install config error")
	}
	file, err := fileutil.LocateFile(fileutil.RewriteFlag(), configFile)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return errors.Wrap(err, "write install config error")
	}
	return nil
}

// load middleware config form file.
func loadInstallConfig(config InitConfiguration) (*kitconfig.InstallConfig, error) {
	installConfig := &kitconfig.InstallConfig{}
	if config.ConfigFile != "" {
		file, err := fileutil.LocateFile(fileutil.RWFlag(), config.ConfigFile)
		if err != nil {
			return nil, errors.Wrap(err, "load install config error")
		}
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, errors.Wrap(err, "read install config error")
		}
		err = yaml.Unmarshal(data, &installConfig)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal install config error")
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
	if installConfig.Plugins == nil {
		installConfig.Plugins = defaultPlugins
	}
	if installConfig.ImagePolicy == "" {
		installConfig.ImagePolicy = config.ImagePolicy
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

func Upgrade(config InitConfiguration) error {
	installConfig, err := loadInstallConfig(config)
	if err != nil {
		return err
	}
	tKeelHelmRepo = config.Repo.Url

	helmConf, err = InitHelmConfig(config.Namespace, getLog(config.DebugMode))
	if err != nil {
		return err
	}

	tKeelMiddle := MiddleConfig{}

	keelChart, middleConfig, err := ChartMiddleware(tkeelKeelHelmChart, config.KeelVersion, config)
	if err != nil {
		return err
	}
	tKeelMiddle.Update(middleConfig)

	coreChart, middleConfig, err := ChartMiddleware(tkeelCoreHelmChart, config.CoreVersion, config)
	if err != nil {
		return err
	}
	tKeelMiddle.Update(middleConfig)

	rudderChart, middleConfig, err := ChartMiddleware(tkeelRudderHelmChart, config.RudderVersion, config)
	if err != nil {
		return err
	}
	tKeelMiddle.Update(middleConfig)

	middlewareChart, err := tKeelChart(config.MiddlewareVersion, tKeelHelmRepo, tKeelMiddlewareHelmChart, helmConf)
	if err != nil {
		return err
	}

	charts := map[string]*chart.Chart{
		tkeelKeelHelmChart:   keelChart,
		tkeelCoreHelmChart:   coreChart,
		tkeelRudderHelmChart: rudderChart,
	}
	updateImagePolicy(charts, installConfig.ImagePolicy)

	middleware := updateChartsValue(charts, tKeelMiddle, installConfig, config)
	installConfig.SetMiddleware(middleware)

	updateComponentsValues(middlewareChart, installConfig)
	updateIngresValues(middlewareChart, installConfig)
	err = updateConfigMap(middlewareChart, installConfig)
	if err != nil {
		return err
	}

	msg := "Deploying the tKeel Platform to your cluster..."

	stopSpinning := print.Spinner(os.Stdout, msg)
	defer stopSpinning(print.Failure)

	err = upgradeTkeel(config, middlewareChart, keelChart, coreChart, rudderChart)
	if err != nil {
		return err
	}
	stopSpinning(print.Success)

	return nil
}

func installPlugins(config InitConfiguration, plugins []string, defaultVersion string) {
	for _, plugin := range plugins {
		repo, name, version := utils.ParseInstallArg(plugin, config.Repo.Name)
		if version == "" {
			version = defaultVersion
		}
		if err := Install(repo, name, version, name, nil); err != nil {
			print.FailureStatusEvent(os.Stdout, "Install %q failed, Because: %s", name, err.Error())
			continue
		}
		print.SuccessStatusEvent(os.Stdout, "Install %q success! It's named %q in k8s", plugin, name)
	}
}
