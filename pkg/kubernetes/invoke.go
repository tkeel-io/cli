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
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"

	"github.com/dapr/cli/pkg/api"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"k8s.io/client-go/rest"
)

// Invoke is a command to invoke a remote or local dapr instance.
func Invoke(pluginID, method string, data []byte, verb string, reqOpts ...RestRequestOption) (string, error) {
	client, err := Client()
	if err != nil {
		return "", err
	}

	app, err := GetAppPod(client, pluginID)
	if err != nil {
		return "", err
	}

	return invoke(client.CoreV1().RESTClient(), &app.AppInfo, method, data, verb, reqOpts...)
}

func invoke(client rest.Interface, app *AppInfo, method string, data []byte, verb string, reqOpts ...RestRequestOption) (string, error) {
	req, err := app.Request(client.Verb(verb), method, data)
	if err != nil {
		return "", fmt.Errorf("error get request: %w", err)
	}

	for i := 0; i < len(reqOpts); i++ {
		if err = reqOpts[i](req); err != nil {
			return "", err
		}
	}

	result := req.Do(context.TODO())
	rawbody, err := result.Raw()
	if err != nil {
		return "", fmt.Errorf("error get raw: %w", err)
	}

	if len(rawbody) > 0 {
		return string(rawbody), nil
	}

	return "", nil
}

// InvokeByPortForward is a command to invoke a remote or local dapr instance.
func InvokeByPortForward(pluginID, method string, data []byte, verb string, reqOpts ...HTTPRequestOption) (string, error) {
	portForward, err := GetPortforward(pluginID, WithHTTPPort, WithAppPod)
	if err != nil {
		return "", err
	}
	// initialize port forwarding.
	if err = portForward.Init(); err == nil {
		url := makeEndpoint(portForward.App, portForward, method)
		// initialize port forwarding.
		// fmt.Println(verb, url)
		var req *http.Request
		req, err = http.NewRequest(verb, url, bytes.NewBuffer(data))
		if err != nil {
			return "", fmt.Errorf("error creat http request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		for i := 0; i < len(reqOpts); i++ {
			if err = reqOpts[i](req); err != nil {
				return "", err
			}
		}

		var httpc http.Client

		var r *http.Response
		r, err = httpc.Do(req)
		if err != nil {
			return "", fmt.Errorf("error do http request: %w", err)
		}
		defer r.Body.Close()
		return readResponse(r)
	}

	portForward.Stop()
	return "", err
}

func makeEndpoint(app *AppPod, pf *PortForward, method string) string {
	return fmt.Sprintf("http://127.0.0.1:%s/v%s/invoke/%s/method/%s", fmt.Sprintf("%v", pf.LocalPort), api.RuntimeAPIVersion, app.AppID, method)
}

// not use dapr api.
func makeWsEndpoint(pf *PortForward, method string) string {
	return fmt.Sprintf("ws://127.0.0.1:%s/%s", fmt.Sprintf("%v", pf.LocalPort), method)
}

func readResponse(response *http.Response) (string, error) {
	rb, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error read http response: %w", err)
	}

	if len(rb) > 0 {
		return string(rb), nil
	}

	return "", nil
}

// WebsocketByPortForward websocket request to the k8s pod.
func WebsocketByPortForward(pluginID, method string, data []byte) (string, error) {
	portForward, err := GetPortforward(pluginID, WithAppPort)
	if err != nil {
		return "", err
	}

	// manage termination of port forwarding connection on interrupt
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	defer signal.Stop(signals)
	go func() {
		<-signals
		os.Exit(0)
	}()

	// initialize port forwarding
	if err = portForward.Init(); err == nil {
		defer portForward.Stop()
		url := makeWsEndpoint(portForward, method)
		fmt.Println(url)

		dialer := websocket.Dialer{}
		connect, resp, err := dialer.Dial(url, nil)
		if nil != err {
			fmt.Println(err)
			return "", errors.Wrap(err, "connect error")
		}
		defer resp.Body.Close()
		defer connect.Close()

		err = connect.WriteMessage(websocket.TextMessage, data)
		if nil != err {
			fmt.Println(err)
			return "", errors.Wrap(err, "websocket write error")
		}

		for {
			messageType, messageData, err := connect.ReadMessage()
			if nil != err {
				return "", errors.Wrap(err, "websocket read error")
			}
			switch messageType {
			case websocket.TextMessage:
				fmt.Println(string(messageData))
			case websocket.BinaryMessage:
				fmt.Println(messageData)
			case websocket.CloseMessage:
			case websocket.PingMessage:
			case websocket.PongMessage:
			default:
			}
		}
	}
	return "", nil
}

type HTTPRequestOption func(*http.Request) error

func InvokeSetHTTPHeader(header, val string) HTTPRequestOption {
	return func(r *http.Request) error {
		r.Header.Set(header, val)
		return nil
	}
}

type RestRequestOption func(*rest.Request) error

func InvokeSetRestRequestHeader(header string, val string) RestRequestOption {
	return func(r *rest.Request) error {
		r.SetHeader(header, val)
		return nil
	}
}
