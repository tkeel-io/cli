// ------------------------------------------------------------
// Copyright 2021 The tKeel Contributors.
// Licensed under the Apache License.
// ------------------------------------------------------------

package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dapr/cli/pkg/age"
	dapr "github.com/dapr/cli/pkg/kubernetes"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/net"
	k8s "k8s.io/client-go/kubernetes"
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
		return nil, fmt.Errorf("k8s query err:%w", err)
	}
	resp := PluginResponse{}
	err = json.Unmarshal(rawbody, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshak json to struct err: %w", err)
	}
	return resp.Data, nil
}

func RegisterPlugins(client k8s.Interface, namespace, pluginID string) error {
	_ = doRegisterPlugins(client, namespace, pluginID)
	plugins, err := ListPlugins(client, namespace)
	if err != nil {
		return err
	}

	notFound := true
	for _, plugin := range plugins {
		if plugin.PluginID == pluginID && plugin.Status == "ACTIVE" {
			notFound = false
			break
		}
	}
	if notFound {
		return errors.Errorf("plugin<%s> not found.", pluginID)
	}

	return nil
}

func doRegisterPlugins(client k8s.Interface, namespace string, pluginID string) error {
	res := client.CoreV1().RESTClient().Post().
		Namespace(namespace).
		Resource("services").
		SubResource("proxy").
		Name(net.JoinSchemeNamePort("", "plugins", "8080")).
		Suffix("register").
		Body([]byte(fmt.Sprintf(`{"id":"%s","secret":"changeme"}`, pluginID)))
	ret := res.Do(context.TODO())
	raw, err := ret.Raw()
	if err != nil {
		return fmt.Errorf("k8s query err:%w", err)
	}
	resp := PluginResponse{}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return fmt.Errorf("unmarshall json to struct err:%w", err)
	}
	return nil
}

func DeletePlugins(client k8s.Interface, namespace string, pluginID string) error {
	res := client.CoreV1().RESTClient().Post().
		Namespace(namespace).
		Resource("services").
		SubResource("proxy").
		Name(net.JoinSchemeNamePort("", "plugins", "8080")).
		Suffix("delete").
		Body([]byte(fmt.Sprintf(`{"id":"%s"}`, pluginID)))
	ret := res.Do(context.TODO())
	raw, err := ret.Raw()
	if err != nil {
		return fmt.Errorf("k8s qeury err:%w", err)
	}
	resp := PluginResponse{}
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return fmt.Errorf("unmarshal json to struct err:%w", err)
	}
	return nil
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
	// Query all namespaces for tKeel pods.
	podList, err := ListPodsInterface(client, map[string]string{
		"app": label,
	})
	if err != nil {
		return nil, err
	}

	if len(podList.Items) == 0 {
		return nil, err
	}
	pod := podList.Items[0]
	replicas := len(podList.Items)
	image := pod.Spec.Containers[0].Image
	namespace := pod.GetNamespace()
	podAge := age.GetAge(pod.CreationTimestamp.Time)
	created := pod.CreationTimestamp.Format("2006-01-02 15:04.05")
	version := image[strings.IndexAny(image, ":")+1:]
	status, healthy := GetStatusAndHealthyInPodList(podList)
	return &StatusOutput{
		Name:      label,
		Namespace: namespace,
		Created:   created,
		Age:       podAge,
		Status:    status,
		Version:   version,
		Healthy:   healthy,
		Replicas:  replicas,
	}, nil
}
