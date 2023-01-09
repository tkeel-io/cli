// ------------------------------------------------------------
// Copyright 2021 The tKeel Contributors.
// Licensed under the Apache License.
// ------------------------------------------------------------

package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dapr/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/client/redis"
	terrors "github.com/tkeel-io/kit/errors"
	tenantApi "github.com/tkeel-io/tkeel/api/tenant/v1"

	"github.com/pkg/errors"
	"github.com/tkeel-io/kit/result"
	v1 "github.com/tkeel-io/tkeel-interface/openapi/v1"
	pluginAPI "github.com/tkeel-io/tkeel/api/plugin/v1"
	"google.golang.org/protobuf/encoding/protojson"
	mate_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	_installPluginFormat          = "apis/rudder/v1/plugins/%s"
	_showPluginFormat             = "apis/rudder/v1/plugins/%s"
	_uninstallPluginFormat        = "apis/rudder/v1/plugins/%s"
	_getInstalledPluginListFormat = "apis/rudder/v1/plugins"
	_registerPluginFormat         = "apis/rudder/v1/tm/plugins/register?id=%s"
	_enablePluginFormat           = "apis/rudder/v1/tm/plugins/%s/tenants/%s"
	_disablePluginFormat          = "apis/rudder/v1/tm/plugins/%s/tenants/%s"
	_enabledPluginFormat          = "apis/rudder/v1/tenants/%s/plugins"
)

type InstalledListOutput struct {
	Name          string `csv:"NAME"`
	Plugin        string `csv:"PLUGIN"`
	PluginVersion string `csv:"PLUGIN VERSION"`
	Repo          string `csv:"REPO"`
	RegisterAt    string `csv:"REGISTER_AT"`
	Status        string `csv:"STATE"`
	Description   string `csv:"DESCRIPTION"`
}

type RepoPluginListOutput struct {
	Name        string `csv:"NAME"`
	Version     string `csv:"VERSION"`
	Status      string `csv:"STATUS"`
	Description string `csv:"DESCRIPTION"`
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

type EnablePluginRequest struct {
	PluginID string `json:"plugin_id"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

	apps, err := ListAppInfos(client)
	if err != nil {
		return nil, fmt.Errorf("err dapr do list: %w", err)
	}

	appGroups := GroupByAppID(apps)
	statuses := make([]StatusOutput, 0, len(appGroups))

	for appID, apps := range appGroups {
		if len(apps) == 0 {
			continue
		}
		firstApp := apps[0]
		replicas := len(apps)
		status, healthy := GetStatusAndHealthyInPodList(apps)
		pluginStatus := "NOT_REGISTER"
		if plugin, ok := pluginsMap[appID]; ok {
			pluginStatus = plugin.Status.String()
		}
		if appID == "keel" || appID == "rudder" || appID == "core" {
			continue
		}
		statuses = append(statuses, StatusOutput{
			Name:         appID,
			Namespace:    firstApp.Namespace,
			Created:      firstApp.Created,
			Age:          firstApp.Age,
			Status:       status,
			Version:      firstApp.Version,
			Healthy:      healthy,
			Replicas:     replicas,
			PluginStatus: pluginStatus,
		})
	}
	return statuses, nil
}

func InstalledPlugin() ([]InstalledListOutput, error) {
	token, err := getAdminToken()
	if err != nil {
		return nil, errors.Wrap(err, "get token error")
	}

	resp, err := InvokeByPortForward(_pluginKeel, _getInstalledPluginListFormat, nil, http.MethodGet, setAuthenticate(token))
	if err != nil {
		return nil, errors.Wrap(err, "InvokeByPortForward error")
	}
	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return nil, errors.Wrap(err, "unmarshal response context error")
	}

	if r.Code != terrors.Success.Reason {
		return nil, fmt.Errorf("invalid response: %s", r.Msg)
	}

	listResponse := pluginAPI.ListPluginResponse{}
	err = r.Data.UnmarshalTo(&listResponse)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal error")
	}
	list := make([]InstalledListOutput, 0, len(listResponse.PluginList))
	for _, o := range listResponse.GetPluginList() {
		registeredAt := time.Unix(o.RegisterTimestamp, 0).Format("2006-01-02 15:04:05")
		list = append(list, InstalledListOutput{
			Name:          o.Id,
			Plugin:        o.InstallerBrief.Name,
			Repo:          o.InstallerBrief.Repo,
			PluginVersion: o.InstallerBrief.Version,
			RegisterAt:    registeredAt,
			Status:        v1.PluginStatus_name[int32(o.Status)],
			Description:   o.InstallerBrief.Desc,
		})
	}
	return list, nil
}

func PluginInfo(pluginID string) ([]InstalledListOutput, error) {
	token, err := getAdminToken()
	if err != nil {
		return nil, errors.Wrap(err, "get token error")
	}
	method := fmt.Sprintf(_showPluginFormat, pluginID)

	resp, err := InvokeByPortForward(_pluginKeel, method, nil, http.MethodGet, setAuthenticate(token))
	if err != nil {
		return nil, errors.Wrap(err, "InvokeByPortForward error")
	}
	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return nil, errors.Wrap(err, "unmarshal response context error")
	}

	if r.Code != terrors.Success.Reason {
		return nil, fmt.Errorf("invalid response: %s", r.Msg)
	}

	response := pluginAPI.GetPluginResponse{}
	err = r.Data.UnmarshalTo(&response)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal error")
	}
	list := make([]InstalledListOutput, 0, 1)
	registeredAt := time.Unix(response.Plugin.GetPlugin().RegisterTimestamp, 0).Format("2006-01-02 15:04:05")
	list = append(list, InstalledListOutput{
		Name:          response.Plugin.GetPlugin().Id,
		Plugin:        response.Plugin.GetPlugin().InstallerBrief.Name,
		Repo:          response.Plugin.GetPlugin().InstallerBrief.Repo,
		PluginVersion: response.Plugin.GetPlugin().InstallerBrief.Version,
		RegisterAt:    registeredAt,
		Status:        v1.PluginStatus_name[int32(response.Plugin.GetPlugin().Status)],
	})
	return list, nil
}

func RegisterPlugin(plugin string) error {
	token, err := getAdminToken()
	if err != nil {
		return errors.Wrap(err, "get token error")
	}
	method := fmt.Sprintf(_registerPluginFormat, plugin)

	resp, err := InvokeByPortForward(_pluginKeel, method, nil, http.MethodGet, setAuthenticate(token))
	if err != nil {
		return errors.Wrap(err, "error invoke")
	}
	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return errors.Wrap(err, "error unmarshal")
	}

	if r.Code == terrors.SUCCESS_CODE {
		return nil
	}
	return errors.Wrap(errors.New(r.Msg), "register failed")
}

func EnablePlugin(pluginID, tenantID string) error {
	token, err := getAdminToken()
	if err != nil {
		return errors.Wrap(err, "get token error")
	}
	method := fmt.Sprintf(_enablePluginFormat, pluginID, tenantID)

	resp, err := InvokeByPortForward(_pluginKeel, method, nil, http.MethodPost, setAuthenticate(token))
	if err != nil {
		return errors.Wrap(err, "error invoke")
	}
	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return errors.Wrap(err, "error unmarshal")
	}

	if r.Code == terrors.Success.Reason {
		return nil
	}

	return errors.Wrap(errors.New(r.Msg), "enable failed")
}

func DisablePlugin(pluginID, tenantID string) error {
	token, err := getAdminToken()
	if err != nil {
		return errors.Wrap(err, "get token error")
	}
	method := fmt.Sprintf(_disablePluginFormat, pluginID, tenantID)
	resp, err := InvokeByPortForward(_pluginKeel, method, nil, http.MethodDelete, setAuthenticate(token))
	if err != nil {
		return errors.Wrap(err, "error invoke")
	}
	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return errors.Wrap(err, "error unmarshal'")
	}

	if r.Code != terrors.Success.Reason {
		return errors.Wrap(errors.New(r.Msg), "error response code")
	}

	return nil
}

func Unregister(pluginID string) (*Plugin, error) {
	clientset, err := Client()
	if err != nil {
		return nil, err
	}

	return UnregisterPlugins(clientset, pluginID)
}

func ListPluginsOfTenant(tenant string) ([]RepoPluginListOutput, error) {
	token, err := getAdminToken()
	if err != nil {
		return nil, errors.Wrap(err, "get token error")
	}
	method := fmt.Sprintf(_enabledPluginFormat, tenant)
	body, err := InvokeByPortForward(_pluginKeel, method, nil, http.MethodGet, setAuthenticate(token))
	if err != nil {
		return nil, errors.Wrap(err, "error invoke")
	}

	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(body), r); err != nil {
		return nil, errors.Wrap(err, "error unmarshal")
	}

	if r.Code != terrors.Success.Reason {
		return nil, errors.Wrap(errors.New(r.Msg), "error response code")
	}

	response := tenantApi.ListTenantPluginResponse{}
	if err = r.Data.UnmarshalTo(&response); err != nil {
		return nil, errors.Wrap(err, "error unmarshal response")
	}

	pluginList, err := InstalledPlugin()
	if err != nil {
		return nil, err
	}
	l := make([]RepoPluginListOutput, 0, len(response.Plugins))
	for _, i := range response.Plugins {
		for _, plugin := range pluginList {
			if plugin.Name == i {
				temp := RepoPluginListOutput{
					i,
					plugin.PluginVersion,
					plugin.Status,
					plugin.Description,
				}
				l = append(l, temp)
			}
		}
	}

	return l, nil
}

func Install(repo, plugin, version, name string, config []byte) error {
	token, err := getAdminToken()
	if err != nil {
		return err
	}
	if version == "" {
		version = "latest"
	}
	method := fmt.Sprintf(_installPluginFormat, name)
	req := pluginAPI.Installer{
		Name:          plugin,
		Version:       version,
		Repo:          repo,
		Configuration: config,
		Type:          1,
	}
	data, err := json.Marshal(&req)
	if err != nil {
		return errors.Wrap(err, "error marshal")
	}
	resp, err := InvokeByPortForward(_pluginKeel, method, data, http.MethodPost, setAuthenticate(token))
	if err != nil {
		return errors.Wrap(err, "error invoke")
	}

	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return errors.Wrap(err, "error unmarshal")
	}

	if r.Code != terrors.Success.Reason {
		return errors.Wrap(errors.New(r.Msg), "error response code")
	}

	return nil
}

func PluginUpgrade(repo, plugin, version, name string, config []byte) error {
	token, err := getAdminToken()
	if err != nil {
		return errors.Wrap(err, "get token error")
	}
	method := fmt.Sprintf(_installPluginFormat, name)
	req := pluginAPI.Installer{
		Name:          plugin,
		Version:       version,
		Repo:          repo,
		Configuration: config,
		Type:          1,
	}
	data, err := json.Marshal(&req)
	if err != nil {
		return errors.Wrap(err, "error marshal")
	}
	resp, err := InvokeByPortForward(_pluginKeel, method, data, http.MethodPut, setAuthenticate(token))
	if err != nil {
		return errors.Wrap(err, "error invoke")
	}

	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return errors.Wrap(err, "error unmarshal")
	}

	if r.Code != terrors.Success.Reason {
		return errors.Wrap(errors.New(r.Msg), "error response code")
	}

	return nil
}

func UninstallPlugin(pluginID string) error {
	token, err := getAdminToken()
	if err != nil {
		return errors.Wrap(err, "get token error")
	}

	method := fmt.Sprintf(_uninstallPluginFormat, pluginID)
	resp, err := InvokeByPortForward(_pluginKeel, method, nil, http.MethodDelete, setAuthenticate(token))
	if err != nil {
		return errors.Wrap(err, "error invoke")
	}
	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return errors.Wrap(err, "error unmarshal")
	}

	if r.Code != terrors.Success.Reason {
		return errors.Wrap(errors.New(r.Msg), "error response code")
	}
	return nil
}

type PluginEnableInfo struct {
	ID        string `json:"id"`
	Installer struct {
		Repo       string `json:"repo"`
		Name       string `json:"name"`
		Version    string `json:"version"`
		Desc       string `json:"desc"`
		Maintainer []struct {
			Name  string `json:"name"`
			URL   string `json:"url"`
			Email string `json:"email"`
		} `json:"maintainer"`
	} `json:"installer"`
	TkeelVersion      string         `json:"tkeel_version"`
	Secret            string         `json:"secret"`
	RegisterTimestamp int            `json:"register_timestamp"`
	Version           string         `json:"version"`
	Status            int            `json:"status"`
	EnableTenantes    []EnableTenant `json:"enable_tenantes"`
}
type EnableTenant struct {
	TenantID        string `json:"tenant_id"`
	OperatorID      string `json:"operator_id"`
	EnableTimestamp int    `json:"enable_timestamp"`
}

func CleanInvalidTenants(pluginID string, tenants []string, namespace string) error {
	tenantMap := make(map[string]struct{})
	for _, tenant := range tenants {
		tenantMap[tenant] = struct{}{}
	}
	password, err := GetRedisPassword(namespace)
	if err != nil {
		return err
	}
	pf, err := GetPodPortForward("tkeel-middleware-redis-master-0", namespace, 6379)
	if err != nil {
		return err
	}
	defer pf.Stop()

	err = pf.Init()
	if err != nil {
		return err
	}

	ctx := context.Background()
	rdb := redis.NewClient("127.0.0.1", pf.LocalPort, password, 0)
	res := rdb.HGet(ctx, fmt.Sprintf("rudder||p_%s", pluginID), "data")
	if res.Err() != nil {
		return res.Err()
	}
	data, err := res.Bytes()
	info := &PluginEnableInfo{}
	err = json.Unmarshal(data, info)
	newEnableTenants := make([]EnableTenant, 0)
	for _, tenant := range info.EnableTenantes {
		if tenant.TenantID == "_tKeel_system" {
			newEnableTenants = append(newEnableTenants, tenant)
		} else if _, ok := tenantMap[tenant.TenantID]; ok {
			newEnableTenants = append(newEnableTenants, tenant)
		}
	}
	info.EnableTenantes = newEnableTenants
	data, err = json.Marshal(info)
	rdb.HSet(ctx, fmt.Sprintf("rudder||p_%s", pluginID), "data", string(data))
	return nil
}

func GetRedisPassword(namespace string) (string, error) {
	_, client, err := kubernetes.GetKubeConfigClient()
	if err != nil {
		return "", fmt.Errorf("get kube config error: %w", err)
	}
	opts := mate_v1.GetOptions{}
	secret, err := client.CoreV1().Secrets(namespace).Get(context.TODO(), "tkeel-middleware-redis", opts)
	if err != nil {
		return "", fmt.Errorf("get secret error: %w", err)
	}

	if value := secret.Data["redis-password"]; value != nil {
		return string(value), err
	}
	return "", fmt.Errorf("get redis password error")
}
