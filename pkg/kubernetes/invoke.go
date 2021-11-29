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

	"github.com/pkg/errors"
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
	res = res.Suffix(makeEndpoint(pluginID, method))
	if len(data) > 0 {
		res = res.Body(data)
	}

	result := res.Do(context.TODO())
	rawbody, err := result.Raw()
	if err != nil {
		return "", errors.Wrap(err, "get raw body err")
	}

	if len(rawbody) > 0 {
		return string(rawbody), nil
	}

	return "", nil
}

func makeEndpoint(appID, method string) string {
	return method
}
