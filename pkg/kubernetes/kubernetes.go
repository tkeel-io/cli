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
	"os"
	"path/filepath"

	dapr "github.com/dapr/cli/pkg/kubernetes"
	"github.com/dapr/cli/pkg/print"
	"github.com/pkg/errors"
	helm "helm.sh/helm/v3/pkg/action"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	tKeelReleaseName             = "tkeel-platform"
	tKeelHelmRepo                = "https://tkeel-io.github.io/helm-charts/"
	tKeelPluginConfigHelmRepo    = "tkeel-plugin-components"
	tKeelPluginComponentHelmRepo = "tkeel-middleware"
	latestVersion                = "latest"
)

var controlPlanePlugins = []string{
	"plugins",
	"keel",
	"auth",
	"iothub",
	"core",
}

var ErrDaprNotInstall = errors.New("dapr is not installed in your cluster")

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

// Init deploys the tKeel operator using the supplied runtime version.
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

	for _, pluginID := range controlPlanePlugins {
		err = RegisterPlugins(clientset, pluginID)
		if err != nil {
			return err
		}
		print.InfoStatusEvent(os.Stdout, "Plugin<%s>  is registered.", pluginID)
	}

	stopSpinning(print.Success)
	return err
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
