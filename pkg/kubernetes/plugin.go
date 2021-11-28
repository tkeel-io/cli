// ------------------------------------------------------------
// Copyright 2021 The tKeel Contributors.
// Licensed under the Apache License.
// ------------------------------------------------------------



package kubernetes

import (
	"fmt"
)

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

func List() ([]StatusOutput, error) {
	client, err := Client()
	if err != nil {
		return nil, err
	}

	tKeelPlugins, err := ListPlugins(client)
	if err != nil {
		return nil, err
	}

	pluginsMap := make(map[string]Plugin)
	for _, plugin := range tKeelPlugins {
		pluginsMap[plugin.PluginID] = plugin
	}

	apps, err := ListPluginPods(client)
	if err != nil {
		return nil, fmt.Errorf("err dapr do list: %w", err)
	}

	appGroups := apps.GroupByAppID()
	statuses := make([]StatusOutput, 0, len(appGroups))

	for appID, lp := range appGroups {
		if len(lp) == 0 {
			continue
		}
		firstPod := lp[0]
		replicas := len(lp)
		info := firstPod.App()
		status, healthy := GetStatusAndHealthyInPodList(lp)
		pluginStatus := "UNKNOWN"
		if plugin, ok := pluginsMap[appID]; ok {
			pluginStatus = plugin.Status
		}
		statuses = append(statuses, StatusOutput{
			Name:         appID,
			Namespace:    info.NameSpace,
			Created:      info.Created,
			Age:          info.Age,
			Status:       status,
			Version:      info.Version,
			Healthy:      healthy,
			Replicas:     replicas,
			PluginStatus: pluginStatus,
		})
	}
	return statuses, nil
}

func Register(pluginID string) error {
	clientset, err := Client()
	if err != nil {
		return err
	}

	err = RegisterPlugins(clientset, pluginID)
	if err != nil {
		return err
	}

	//check
	plugins, err := ListPlugins(clientset)
	if err != nil {
		return err
	}

	for _, plugin := range plugins {
		if plugin.PluginID == pluginID && plugin.Status == "ACTIVE" {
			return nil
		}
	}
	return fmt.Errorf("plugin<%s> not found", pluginID)

}

func Remove(pluginID string) error {
	clientset, err := Client()
	if err != nil {
		return err
	}

	return DeletePlugins(clientset, pluginID)
}
