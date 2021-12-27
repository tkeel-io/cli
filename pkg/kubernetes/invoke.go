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
	"bytes"
	"context"
	"fmt"
	"github.com/dapr/cli/pkg/api"
	"github.com/dapr/cli/pkg/kubernetes"
	"io/ioutil"
	core_v1 "k8s.io/api/core/v1"
	"net/http"
	"os"
	"os/signal"
)

// Invoke is a command to invoke a remote or local dapr instance.
func Invoke(pluginID, method string, data []byte, verb string) (string, error) {
	client, err := Client()
	if err != nil {
		return "", err
	}

	app, err := AppPod(client, pluginID)
	if err != nil {
		return "", err
	}

	res := app.App().Request(client.CoreV1().RESTClient().Verb(verb))
	if data != nil {
		res = res.Body(data)
	}
	res = res.RequestURI(method)

	result := res.Do(context.TODO())
	rawbody, err := result.Raw()
	if err != nil {
		return "", fmt.Errorf("error on Invoke: %w", err)
	}

	if len(rawbody) > 0 {
		return string(rawbody), nil
	}

	return "", nil
}

// Invoke is a command to invoke a remote or local dapr instance.
func InvokeByPortForward(pluginID, method string, data []byte, verb string) (string, error) {
	config, client, err := kubernetes.GetKubeConfigClient()
	if err != nil {
		return "", err
	}

	// manage termination of port forwarding connection on interrupt
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	defer signal.Stop(signals)

	pod, err := AppPod(client, pluginID)
	if err != nil {
		return "", err
	}

	if pod.Status.Phase != core_v1.PodRunning {
		return "", fmt.Errorf("no running pods found for %s", pluginID)
	}

	a := pod.App()
	portForward, err := NewPortForward(
		config,
		a.Namespace, a.PodName,
		"127.0.0.1",
		0,
		a.HTTPPort,
		false,
	)
	if err != nil {
		return "", err
	}

	// initialize port forwarding
	if err = portForward.Init(); err == nil {
		url := makeEndpoint(a, portForward, method)
		fmt.Println(url)
		req, err := http.NewRequest(verb, url, bytes.NewBuffer(data))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")

		var httpc http.Client

		r, err := httpc.Do(req)
		if err != nil {
			return "", err
		}
		defer r.Body.Close()
		return handleResponse(r)
	}

	portForward.Stop()
	return "", nil
}

func makeEndpoint(app App, pf *PortForward, method string) string {
	return fmt.Sprintf("http://127.0.0.1:%s/v%s/invoke/%s/method/%s", fmt.Sprintf("%v", pf.LocalPort), api.RuntimeAPIVersion, app.AppID, method)
}

func handleResponse(response *http.Response) (string, error) {
	rb, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if len(rb) > 0 {
		return string(rb), nil
	}

	return "", nil
}
