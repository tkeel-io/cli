// ------------------------------------------------------------
// Copyright 2021 The tKeel Contributors.
// Licensed under the Apache License.
// ------------------------------------------------------------

package kubernetes

// ListOutput represents the application ID, application port and creation time.
type ListOutput struct {
	AppID   string `csv:"APP ID"`
	AppPort string `csv:"APP PORT"`
	Age     string `csv:"AGE"`
	Created string `csv:"CREATED"`
	Status  string `json:"status"`
}

type Plugin struct {
	PluginID     string `json:"plugin_id"`
	Version      string `json:"version"`
	Secret       string `json:"secret"`
	RegisterTime int64  `json:"register_time"`
	Status       string `json:"status"`
}

type PluginResponse struct {
	Ret  int      `json:"ret"`
	Msg  string   `json:"msg"`
	Data []Plugin `json:"data"`
}

func List() ([]Plugin, error) {
	client, err := Client()
	if err != nil {
		return nil, err
	}
	namespace, err := GetTKeelNameSpace(client)
	if client == nil {
		return nil, err
	}

	return ListPlugins(client, namespace)
}

func Register(pluginID string) error {
	clientset, err := Client()
	if err != nil {
		return err
	}

	namespace, err := GetTKeelNameSpace(clientset)
	if err != nil {
		return err
	}

	return RegisterPlugins(clientset, namespace, pluginID)
}

func Delete(pluginID string) error {
	clientset, err := Client()
	if err != nil {
		return err
	}

	namespace, err := GetTKeelNameSpace(clientset)
	if err != nil {
		return err
	}

	return DeletePlugins(clientset, namespace, pluginID)
}
