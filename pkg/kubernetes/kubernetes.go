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
	"github.com/mitchellh/mapstructure"
	"github.com/tkeel-io/cli/fileutil"
	"helm.sh/helm/v3/pkg/chart"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"strings"

	dapr "github.com/dapr/cli/pkg/kubernetes"
	"github.com/dapr/cli/pkg/print"
	"github.com/pkg/errors"
	helm "helm.sh/helm/v3/pkg/action"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	tKeelReleaseName               = "tkeel-platform"
	tKeelMiddlewareReleaseName     = "tkeel-platform-middleware"
	tKeelPluginConfigHelmChart     = "tkeel-plugin-components"
	tKeelPluginMiddlewareHelmChart = "tkeel-middleware"
	tkeelKeelHelmChart             = "keel"
	tkeelRudderHelmChart           = "rudder"
	tkeelCoreHelmChart             = "core"
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
	Repo          *Repo
	DefaultConfig bool
}

type InstallConfig struct {
	Middleware map[string]interface{} `yaml:"middleware"`
	Repo       map[string]interface{} `yaml:"repo"`
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

type Repo struct {
	Url  string
	Name string
}

// Init deploys the tKeel operator using the supplied runtime version.
func Init(config InitConfiguration) (err error) {
	print.InfoStatusEvent(os.Stdout, "Checking the Dapr runtime status...")
	err = check(config)
	if err != nil {
		return err
	}

	helmConf, err = helmConfig(config.Namespace, getLog(config.DebugMode))
	if err != nil {
		return err
	}

	var middlewareMap = make(map[string]string)
	middlewareConfig, err := loadMiddlewareConfig(config.ConfigFile)
	if err != nil {
		return err
	}
	err = mapstructure.Decode(middlewareConfig, &middlewareMap)
	if err != nil {
		return err
	}

	tKeelHelmRepo = config.Repo.Url

	keelChart, componentMiddleware, err := KeelChart(config, config.Password)
	if err != nil {
		return err
	}

	middlewareChart, err := MiddlewareChart(config)
	if err != nil {
		return err
	}

	dependencies := make(map[string]*chart.Chart)
	for _, k := range middlewareChart.Dependencies() {
		dependencies[k.Name()] = k
	}
	for _, v := range middlewareMap {
		items := strings.Split(v, ",")
		res, err := url.Parse(items[0])
		if err != nil {
			return err
		}
		if _, ok := dependencies[res.Scheme]; ok {
			delete(dependencies, res.Scheme)
		}
	}
	newDependency := make([]*chart.Chart, 0)
	for _, v := range dependencies {
		newDependency = append(newDependency, v)
	}
	middlewareChart.SetDependencies(newDependency...)
	// cover middleware config with custom middleware
	for component, middleware := range componentMiddleware {
		for k := range middleware {
			if v, ok := middlewareMap[k]; ok {
				middleware[k] = v
			}
		}
		if component == tkeelKeelHelmChart {
			keelChart.Values["middleware"] = middleware
		} else {
			keelChart.Values[component] = map[string]interface{}{"middleware": middleware}
		}
	}

	err = deploy(config, middlewareChart, keelChart)
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

	// err = registerPlugins(config)
	// if err != nil {
	// 	return err
	// }

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

// func registerPlugins(config InitConfiguration) error {
// 	msg := "Register the plugins ..."

// 	stopSpinning := print.Spinner(os.Stdout, msg)
// 	defer stopSpinning(print.Failure)

// 	clientset, err := Client()
// 	if err != nil {
// 		return err
// 	}

// 	for _, pluginID := range controlPlanePlugins {
// 		err = RegisterPlugins(clientset, pluginID)
// 		if err != nil {
// 			return err
// 		}
// 		print.InfoStatusEvent(os.Stdout, "Plugin<%s>  is registered.", pluginID)
// 	}

// 	stopSpinning(print.Success)
// 	return err
// }

// load middleware config form file
func loadMiddlewareConfig(config string) (*MiddlewareConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	config = strings.Replace(config, "~", home, 1)
	middlewareConfig := &MiddlewareConfig{}
	file, err := fileutil.LocateFile(os.O_RDONLY, config)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &middlewareConfig)
	if err != nil {
		return nil, err
	}
	return middlewareConfig, nil
}

func check(config InitConfiguration) error {
	client, err := dapr.NewStatusClient()
	if err != nil {
		return fmt.Errorf("can't connect to a Kubernetes cluster: %w", err)
	}
	status, err := client.Status()
	if err != nil {
		return fmt.Errorf("can't connect to a Kubernetes cluster: %w", err)
	}
	if len(status) == 0 {
		return ErrDaprNotInstall
	}
	enabled, err := dapr.IsMTLSEnabled()
	if err != nil {
		return fmt.Errorf("can't connect to a Kubernetes cluster: %w", err)
	}
	if !enabled {
		return errors.New("dapr mtls is disabled")
	}
	if status[0].Namespace != config.Namespace {
		return fmt.Errorf("dapr is installed in namespace: `%v`, not in `%v`\nUse `-n %v` flag", status[0].Namespace, config.Namespace, status[0].Namespace)
	}
	return nil
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
