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
	"k8s.io/client-go/rest"
	"net/url"
	"strings"
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
	res, err = makeEndpoint(res, pluginID, method)
	if err != nil {
		return "", fmt.Errorf("error in make endpoint: %w", err)
	}

	result := res.Do(context.TODO())
	rawbody, err := result.Raw()
	if err != nil {
		return "", fmt.Errorf("error get raw: %w", err)
	}

	if len(rawbody) > 0 {
		return string(rawbody), nil
	}

	return "", nil
}

func makeEndpoint(res *rest.Request, appID, method string) (*rest.Request, error) {
	tempURL, err := url.Parse(fmt.Sprintf("http://127.0.0.1/%s", method))
	if err != nil {
		return nil, err
	}
	res = res.Suffix(tempURL.Path)
	for k, vs := range tempURL.Query() {
		res = res.Param(k, strings.Join(vs, ","))
	}
	return res, nil
}
