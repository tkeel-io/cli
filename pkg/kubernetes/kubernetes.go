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
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"path/filepath"
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

	err = deploy(config, controlPlanePlugins)
	if err != nil {
		return err
	}

	err = registerPlugins(config)
	if err != nil {
		return err
	}

	return nil
}

func deploy(config InitConfiguration, pluginNames []string) (err error) {
	msg := "Deploying the tKeel Platform to your cluster..."

	stopSpinning := print.Spinner(os.Stdout, msg)
	defer stopSpinning(print.Failure)

	err = install(config, pluginNames)
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
