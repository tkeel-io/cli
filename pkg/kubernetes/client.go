// ------------------------------------------------------------
// Copyright 2021 The TKeel Contributors.
// Licensed under the Apache License.
// ------------------------------------------------------------

package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dapr/cli/pkg/age"
	dapr "github.com/dapr/cli/pkg/kubernetes"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/net"
	k8s "k8s.io/client-go/kubernetes"
	"strings"
)

var Client = dapr.Client

func ListPlugins(client k8s.Interface, namespace string) ([]Plugin, error) {
	res := client.CoreV1().RESTClient().Get().
		Namespace(namespace).
		Resource("services").
		SubResource("proxy").
		Name(net.JoinSchemeNamePort("", "plugins", "8080")).
		Suffix("list")
	result := res.Do(context.TODO())
	rawbody, err := result.Raw()
	if err != nil {
		return nil, err
	}
	resp := PluginResponse{}
	err = json.Unmarshal(rawbody, &resp)
	if err != nil {
		return nil, err
	} else {
		return resp.Data, nil
	}

}

func RegisterPlugins(client k8s.Interface, namespace, pluginId string) error {
	_ = doRegisterPlugins(client, namespace, pluginId)
	plugins, err := ListPlugins(client, namespace)
	if err != nil {
		return err
	}

	notFound := true
	for _, plugin := range plugins {
		if plugin.PluginId == pluginId && plugin.Status == "ACTIVE" {
			notFound = false
			break
		}
	}
	if notFound {
		return errors.Errorf("plugin<%s> not found.", pluginId)
	}
	
	return nil
}

func doRegisterPlugins(client k8s.Interface, namespace string, pluginId string) error {
	res := client.CoreV1().RESTClient().Post().
		Namespace(namespace).
		Resource("services").
		SubResource("proxy").
		Name(net.JoinSchemeNamePort("", "plugins", "8080")).
		Suffix("register").
		Body([]byte(fmt.Sprintf(`{"id":"%s","secret":"changeme"}`, pluginId)))
	ret := res.Do(context.TODO())
	raw, err := ret.Raw()
	if err != nil {
		return err
	}
	resp := PluginResponse{}
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func DeletePlugins(client k8s.Interface, namespace string, pluginId string) error {
	res := client.CoreV1().RESTClient().Post().
		Namespace(namespace).
		Resource("services").
		SubResource("proxy").
		Name(net.JoinSchemeNamePort("", "plugins", "8080")).
		Suffix("delete").
		Body([]byte(fmt.Sprintf(`{"id":"%s"}`, pluginId)))
	ret := res.Do(context.TODO())
	raw, err := ret.Raw()
	if err != nil {
		return err
	}
	resp := PluginResponse{}
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func GetTKeelNameSpace(client k8s.Interface) (string, error) {
	pod, err := GetPodsInterface(client, "keel")
	if err != nil {
		return "", err
	}
	if pod == nil {
		return "", errors.New("tKeel not initialized")
	}
	return pod.Namespace, nil
}

func GetPodsInterface(client k8s.Interface, label string) (*StatusOutput, error) {
	// Query all namespaces for TKeel pods.
	p, err := ListPodsInterface(client, map[string]string{
		"app": label,
	})
	if err != nil {
		return nil, err
	}

	if len(p.Items) == 0 {
		return nil, err
	}
	pod := p.Items[0]
	replicas := len(p.Items)
	image := pod.Spec.Containers[0].Image
	namespace := pod.GetNamespace()
	age := age.GetAge(pod.CreationTimestamp.Time)
	created := pod.CreationTimestamp.Format("2006-01-02 15:04.05")
	version := image[strings.IndexAny(image, ":")+1:]
	status := ""

	// loop through all replicas and update to Running/Healthy status only if all instances are Running and Healthy
	healthy := "False"
	running := true

	for _, p := range p.Items {
		if len(p.Status.ContainerStatuses) == 0 {
			status = string(p.Status.Phase)
		} else if p.Status.ContainerStatuses[0].State.Waiting != nil {
			status = fmt.Sprintf("Waiting (%s)", p.Status.ContainerStatuses[0].State.Waiting.Reason)
		} else if pod.Status.ContainerStatuses[0].State.Terminated != nil {
			status = "Terminated"
		}

		if len(p.Status.ContainerStatuses) == 0 ||
			p.Status.ContainerStatuses[0].State.Running == nil {
			running = false
			break
		}

		if p.Status.ContainerStatuses[0].Ready {
			healthy = "True"
		}
	}

	if running {
		status = "Running"
	}

	so := StatusOutput{
		Name:      label,
		Namespace: namespace,
		Created:   created,
		Age:       age,
		Status:    status,
		Version:   version,
		Healthy:   healthy,
		Replicas:  replicas,
	}
	return &so, nil
}
