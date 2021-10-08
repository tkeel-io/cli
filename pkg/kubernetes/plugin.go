// ------------------------------------------------------------
// Copyright 2021 The TKeel Contributors.
// Licensed under the Apache License.
// ------------------------------------------------------------

package kubernetes

func Register(pluginId string) error {
	clientset, err := Client()
	if err != nil {
		return err
	}

	namespace, err := GetTKeelNameSpace(clientset)
	if err != nil {
		return err
	}

	return RegisterPlugins(clientset, namespace, pluginId)
}

func Delete(pluginId string) error {
	clientset, err := Client()
	if err != nil {
		return err
	}

	namespace, err := GetTKeelNameSpace(clientset)
	if err != nil {
		return err
	}

	return DeletePlugins(clientset, namespace, pluginId)
}
