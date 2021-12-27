// ------------------------------------------------------------
// Copyright 2021 The tKeel Contributors.
// Licensed under the Apache License.
// ------------------------------------------------------------

package kubernetes

import (
	"encoding/json"
	"fmt"

	v1 "github.com/tkeel-io/tkeel-interface/openapi/v1"
)

// ListOutput represents the application ID, application port and creation time.
type ListOutput struct {
	AppID   string `csv:"APP ID"`
	AppPort string `csv:"APP PORT"`
	Age     string `csv:"AGE"`
	Created string `csv:"CREATED"`
	Status  string `json:"status"`
}

type RegisterAddons struct {
	Addons   string `json:"addons,omitempty"`   // addons name.
	Upstream string `json:"upstream,omitempty"` // upstream path.
}

type Plugin struct {
	ID                string                  `json:"id,omitempty"`                 // plugin id.
	PluginVersion     string                  `json:"plugin_version,omitempty"`     // plugin version.
	TkeelVersion      string                  `json:"tkeel_version,omitempty"`      // plugin depend tkeel version.
	AddonsPoint       []*v1.AddonsPoint       `json:"addons_point,omitempty"`       // plugin declares addons.
	ImplementedPlugin []*v1.ImplementedPlugin `json:"implemented_plugin,omitempty"` // plugin implemented plugin list.
	Secret            string                  `json:"secret,omitempty"`             // plugin registered secret.
	RegisterTimestamp int64                   `json:"register_timestamp,omitempty"` // register timestamp.
	ActiveTenantes    []string                `json:"active_tenantes,omitempty"`    // active tenant's id list.
	RegisterAddons    []*RegisterAddons       `json:"register_addons,omitempty"`    // register addons router.
	Status            v1.PluginStatus         `json:"status,omitempty"`             // register plugin status.
}

func (p *Plugin) String() string {
	b, err := json.Marshal(p)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

type ListResponse struct {
	PluginList []*Plugin `json:"plugin_list"`
}

type DeleteResponse struct {
	Plugin *Plugin `json:"plugin"`
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

	pluginsMap := make(map[string]*Plugin)
	for _, plugin := range tKeelPlugins {
		pluginsMap[plugin.ID] = plugin
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
		pluginStatus := "NOT_REGISTER"
		if plugin, ok := pluginsMap[appID]; ok {
			pluginStatus = plugin.Status.String()
		}
		if appID == "keel" || appID == "rudder" || appID == "core" {
			continue
		}
		statuses = append(statuses, StatusOutput{
			Name:         appID,
			Namespace:    info.Namespace,
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

	// check
	plugins, err := ListPlugins(clientset)
	if err != nil {
		return err
	}

	for _, plugin := range plugins {
		if plugin.ID == pluginID && plugin.Status == v1.PluginStatus_RUNNING {
			return nil
		}
	}
	return fmt.Errorf("plugin<%s> not found", pluginID)
}

func Unregister(pluginID string) (*Plugin, error) {
	clientset, err := Client()
	if err != nil {
		return nil, err
	}

	return UnregisterPlugins(clientset, pluginID)
}
