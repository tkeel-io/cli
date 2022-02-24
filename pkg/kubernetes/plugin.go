// ------------------------------------------------------------
// Copyright 2021 The tKeel Contributors.
// Licensed under the Apache License.
// ------------------------------------------------------------

package kubernetes

import (
	"encoding/json"
	"fmt"
	terrors "github.com/tkeel-io/kit/errors"
	tenantApi "github.com/tkeel-io/tkeel/api/tenant/v1"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/tkeel-io/cli/fileutil"
	"github.com/tkeel-io/kit/result"
	v1 "github.com/tkeel-io/tkeel-interface/openapi/v1"
	pluginAPI "github.com/tkeel-io/tkeel/api/plugin/v1"
	repoAPI "github.com/tkeel-io/tkeel/api/repo/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	_getInstallerListFromRepoFormat = "apis/rudder/v1/repos/%s/installers"
	_showInstallerFormat            = "apis/rudder/v1/repos/%s/installers/%s/%s"
	_installPluginFormat            = "apis/rudder/v1/plugins/%s"
	_showPluginFormat               = "apis/rudder/v1/plugins/%s"
	_uninstallPluginFormat          = "apis/rudder/v1/plugins/%s"
	_getInstalledPluginListFormat   = "apis/rudder/v1/plugins"

	_enablePluginFormat  = "apis/rudder/v1/tenants/%s/plugins"
	_disablePluginFormat = "apis/rudder/v1/tenants/%s/plugins/%s"
	_enabledPluginFormat = "apis/security/v1/tenants/%s/plugins"
)

var ErrInvalidToken = errors.New("invalid token")

// ListOutput represents the application ID, application port and creation time.
type ListOutput struct {
	AppID   string `csv:"APP ID"`
	AppPort string `csv:"APP PORT"`
	Age     string `csv:"AGE"`
	Created string `csv:"CREATED"`
	Status  string `json:"status"`
}

type InstalledListOutput struct {
	Name          string `csv:"NAME"`
	Plugin        string `csv:"PLUGIN"`
	PluginVersion string `csv:"PLUGIN VERSION"`
	Repo          string `csv:"REPO"`
	RegisterAt    string `csv:"REGISTER_AT"`
	Status        string `csv:"STATE"`
}

type RepoPluginListOutput struct {
	Name    string `csv:"NAME"`
	Version string `csv:"VERSION"`
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
	PluginId string `json:"plugin_id"`
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

func InstalledList() ([]InstalledListOutput, error) {
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
		})
	}
	return list, nil
}

func PluginInfo(pluginId string) ([]InstalledListOutput, error) {
	token, err := getAdminToken()
	if err != nil {
		return nil, errors.Wrap(err, "get token error")
	}
	method := fmt.Sprintf(_showPluginFormat, pluginId)

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

func EnablePlugin(pluginId, tenantId string) error {
	token, err := getAdminToken()
	if err != nil {
		return errors.Wrap(err, "get token error")
	}
	method := fmt.Sprintf(_enablePluginFormat, tenantId)
	req := EnablePluginRequest{PluginId: pluginId}
	data, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "can't marshal'")
	}

	resp, err := InvokeByPortForward(_pluginKeel, method, data, http.MethodPost, setAuthenticate(token))
	if err != nil {
		return errors.Wrap(err, "invoke "+method+" error")
	}
	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return errors.Wrap(err, "can't unmarshal'")
	}

	if r.Code == terrors.Success.Reason {
		return nil
	}

	return errors.New("enable failed")
}

func DisablePlugin(pluginId, tenantId string) error {
	token, err := getAdminToken()
	if err != nil {
		return errors.Wrap(err, "get token error")
	}
	method := fmt.Sprintf(_disablePluginFormat, tenantId, pluginId)
	resp, err := InvokeByPortForward(_pluginKeel, method, nil, http.MethodDelete, setAuthenticate(token))
	if err != nil {
		return errors.Wrap(err, "invoke "+method+" error")
	}
	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return errors.Wrap(err, "can't unmarshal'")
	}

	if r.Code == terrors.Success.Reason {
		return nil
	}

	return errors.New("register failed")
}

func Unregister(pluginID string) (*Plugin, error) {
	clientset, err := Client()
	if err != nil {
		return nil, err
	}

	return UnregisterPlugins(clientset, pluginID)
}

func ListPluginsFromTenant(tenant string) ([]RepoPluginListOutput, error) {
	token, err := getAdminToken()
	if err != nil {
		return nil, errors.Wrap(err, "get token error")
	}
	method := fmt.Sprintf(_enabledPluginFormat, tenant)
	body, err := InvokeByPortForward(_pluginKeel, method, nil, http.MethodGet, setAuthenticate(token))
	if err != nil {
		return nil, errors.Wrap(err, "invoke "+method+" error")
	}

	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(body), r); err != nil {
		return nil, errors.Wrap(err, "unmarshal response context error")
	}

	if r.Code != terrors.Success.Reason {
		return nil, fmt.Errorf("invalid response: %s", r.Msg)
	}

	fmt.Println(body)
	response := tenantApi.ListTenantPluginResponse{}
	if err = r.Data.UnmarshalTo(&response); err != nil {
		return nil, errors.Wrap(err, "cant handle response data")
	}

	l := make([]RepoPluginListOutput, 0, len(response.Plugins))
	for _, i := range response.Plugins {
		l = append(l, RepoPluginListOutput{i, ""})
	}

	return l, nil
}

func ListPluginsFromRepo(repo string) ([]RepoPluginListOutput, error) {
	token, err := getAdminToken()
	if err != nil {
		return nil, errors.Wrap(err, "get token error")
	}
	// if auth token is not bearer type ?
	method := fmt.Sprintf(_getInstallerListFromRepoFormat, repo)
	body, err := InvokeByPortForward(_pluginKeel, method, nil, http.MethodGet, setAuthenticate(token))
	if err != nil {
		return nil, errors.Wrap(err, "invoke "+method+" error")
	}

	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(body), r); err != nil {
		return nil, errors.Wrap(err, "unmarshal response context error")
	}

	if r.Code != terrors.Success.Reason {
		return nil, fmt.Errorf("invalid response: %s", r.Msg)
	}

	listResponse := repoAPI.ListRepoInstallerResponse{}
	if err = r.Data.UnmarshalTo(&listResponse); err != nil {
		return nil, errors.Wrap(err, "cant handle response data")
	}

	l := make([]RepoPluginListOutput, 0, len(listResponse.BriefInstallers))
	for _, i := range listResponse.BriefInstallers {
		l = append(l, RepoPluginListOutput{i.Name, i.Version})
	}

	return l, nil
}

func Install(repo, plugin, version, name string, config []byte) error {
	token, err := getAdminToken()
	if err != nil {
		return err
	}
	method := fmt.Sprintf(_installPluginFormat, name)
	req := pluginAPI.Installer{
		Name:          plugin,
		Version:       version,
		Repo:          repo,
		Configuration: config,
		Type:          1,
	}
	data, err := json.Marshal(req) //nolint
	if err != nil {
		return errors.Wrap(err, "marshal plugin request failed")
	}
	resp, err := InvokeByPortForward(_pluginKeel, method, data, http.MethodPost, setAuthenticate(token))
	if err != nil {
		return errors.Wrap(err, "invoke "+method+" error")
	}

	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return errors.Wrap(err, "can't unmarshal'")
	}

	if r.Code == terrors.Success.Reason {
		return nil
	}

	return errors.New("can't handle this")
}

func UninstallPlugin(pluginID string) error {
	token, err := getAdminToken()
	if err != nil {
		return err
	}

	method := fmt.Sprintf(_uninstallPluginFormat, pluginID)
	resp, err := InvokeByPortForward(_pluginKeel, method, nil, http.MethodDelete, setAuthenticate(token))
	if err != nil {
		return errors.Wrap(err, "invoke "+method+" error")
	}
	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return errors.Wrap(err, "can't unmarshal'")
	}

	if r.Code == terrors.Success.Reason {
		return nil
	}
	return errors.New("uninstall plugin failed")
}

func getAdminToken() (string, error) {
	f, err := fileutil.LocateAdminToken()
	if err != nil {
		return "", errors.Wrap(err, "open admin token failed")
	}
	tokenb := make([]byte, 512)
	n, err := f.Read(tokenb)
	if err != nil {
		return "", errors.Wrap(err, "read token failed")
	}
	if n == 0 {
		return "", ErrInvalidToken
	}

	return fmt.Sprintf("Bearer %s", tokenb[:n]), nil
}

func setAuthenticate(token string) HTTPRequestOption {
	return InvokeSetHTTPHeader("Authorization", token)
}
