// ------------------------------------------------------------
// Copyright 2021 The TKeel Contributors.
// Licensed under the Apache License.
// ------------------------------------------------------------

package kubernetes

import (
	"context"
	"fmt"
	dapr "github.com/dapr/cli/pkg/kubernetes"
	"github.com/dapr/cli/pkg/print"
	"github.com/pkg/errors"
	helm "helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/helm/pkg/strvals"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	tKeelReleaseName             = "tkeel-platform"
	tKeelHelmRepo                = "https://tkeel-io.github.io/helm-charts/"
	tKeelPluginConfigHelmRepo    = "tkeel-plugin-components"
	tKeelPluginComponentHelmRepo = "tkeel-middleware"
	latestVersion                = "latest"
)

var (
	controlPlanePlugins = []string{
		"plugins",
		"keel",
		"auth",
		"iothub",
		"core",
	}
	DaprNotInstall = errors.New("dapr is not installed in your cluster")
)

type InitConfiguration struct {
	Version    string
	Namespace  string
	EnableMTLS bool
	EnableHA   bool
	Args       []string
	Wait       bool
	Timeout    uint
	DebugMode  bool
}

// Init deploys the TKeel operator using the supplied runtime version.
func Init(config InitConfiguration) (err error) {
	print.InfoStatusEvent(os.Stdout, "Checking the Dapr runtime status...")
	err = check(config)
	if err != nil {
		return err
	}

	err = deploy(config)
	if err != nil {
		return err
	}

	err = registerPlugins(config)
	if err != nil {
		return err
	}

	return nil
}

func deploy(config InitConfiguration) (err error) {
	msg := "Deploying the tKeel Platform to your cluster..."

	stopSpinning := print.Spinner(os.Stdout, msg)
	defer stopSpinning(print.Failure)

	err = install(config)
	if err != nil {
		return err
	}
	stopSpinning(print.Success)
	return err
}

func registerPlugins(config InitConfiguration) error {
	msg := "Register the plugins ..."

	stopSpinning := print.Spinner(os.Stdout, msg)
	defer stopSpinning(print.Failure)

	clientset, err := Client()
	if err != nil {
		return err
	}

	for _, pluginId := range controlPlanePlugins {
		err = RegisterPlugins(clientset, config.Namespace, pluginId)
		if err != nil {
			return err
		}
		print.InfoStatusEvent(os.Stdout, "Plugin<%s>  is registered.", pluginId)
	}

	stopSpinning(print.Success)
	return err
}

func check(config InitConfiguration) error {
	client, err := dapr.NewStatusClient()
	if err != nil {
		return fmt.Errorf("can't connect to a Kubernetes cluster: %v", err)
	}
	status, err := client.Status()
	if err != nil {
		return fmt.Errorf("can't connect to a Kubernetes cluster: %v", err)
	}
	if len(status) == 0 {
		return DaprNotInstall
	}
	if status[0].Namespace != config.Namespace {
		return fmt.Errorf("dapr is installed in namespace: `%v`, not in `%v`\nUse `-n %v` flag", status[0].Namespace, config.Namespace, status[0].Namespace)
	}
	return nil
}

func createNamespace(namespace string) error {
	_, client, err := dapr.GetKubeConfigClient()
	if err != nil {
		return fmt.Errorf("can't connect to a Kubernetes cluster: %v", err)
	}

	ns := &v1.Namespace{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: namespace,
		},
	}
	// try to create the namespace if it doesn't exist. ok to ignore error.
	_, _ = client.CoreV1().Namespaces().Create(context.TODO(), ns, meta_v1.CreateOptions{})
	return nil
}

func helmConfig(namespace string, log helm.DebugLog) (*helm.Configuration, error) {
	ac := helm.Configuration{}
	flags := &genericclioptions.ConfigFlags{
		Namespace: &namespace,
	}
	err := ac.Init(flags, namespace, "secret", log)
	return &ac, err
}

func createTempDir() (string, error) {
	dir, err := ioutil.TempDir("", "tkeel")
	if err != nil {
		return "", fmt.Errorf("error creating temp dir: %s", err)
	}
	return dir, nil
}

func locateChartFile(dirPath string) (string, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return "", err
	}
	return filepath.Join(dirPath, files[0].Name()), nil
}

func tKeelChart(version, chartName string, config *helm.Configuration) (*chart.Chart, error) {
	pull := helm.NewPull()
	pull.RepoURL = tKeelHelmRepo
	pull.Settings = &cli.EnvSettings{}

	if version != latestVersion {
		pull.Version = chartVersion(version)
	}

	dir, err := createTempDir()
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	pull.DestDir = dir

	_, err = pull.Run(chartName)
	if err != nil {
		return nil, err
	}

	chartPath, err := locateChartFile(dir)
	if err != nil {
		return nil, err
	}
	return loader.Load(chartPath)
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
			return nil, err
		}
	}
	return chartVals, nil
}

func install(config InitConfiguration) error {
	err := createNamespace(config.Namespace)
	if err != nil {
		return err
	}

	helmConf, err := helmConfig(config.Namespace, getLog(config.DebugMode))
	if err != nil {
		return err
	}

	platformChart, err := tKeelChart(config.Version, tKeelPluginComponentHelmRepo, helmConf)
	if err != nil {
		return err
	}

	for _, corePluginName := range controlPlanePlugins {
		tKeelPluginChart, err := tKeelChart(config.Version, corePluginName, helmConf)
		if err != nil {
			return err
		}

		configChart, err := tKeelChart(config.Version, tKeelPluginConfigHelmRepo, helmConf)
		if err != nil {
			return err
		}
		configChart.Values["pluginID"] = corePluginName
		tKeelPluginChart.AddDependency(configChart)

		if tKeelPluginChart.Values["daprConfig"] != corePluginName {
			print.InfoStatusEvent(os.Stdout, "Update Plugin<%s>'s Values[daprConfig] = %s.", corePluginName, corePluginName)
			tKeelPluginChart.Values["daprConfig"] = corePluginName
		}
		platformChart.AddDependency(tKeelPluginChart)
	}

	installClient := helm.NewInstall(helmConf)
	installClient.ReleaseName = tKeelReleaseName
	installClient.Namespace = config.Namespace
	installClient.Wait = config.Wait
	installClient.Timeout = time.Duration(config.Timeout) * time.Second

	values, err := chartValues(config)
	if err != nil {
		return err
	}

	if _, err = installClient.Run(platformChart, values); err != nil {
		return err
	}

	print.InfoStatusEvent(os.Stdout, "install plugin<%s> done.", strings.Join(controlPlanePlugins, ", "))

	return nil
}

func getLog(DebugMode bool) helm.DebugLog {
	if DebugMode {
		return func(format string, v ...interface{}) {
			print.InfoStatusEvent(os.Stdout, format, v...)
		}
	} else {
		return func(format string, v ...interface{}) {

		}
	}
}
