package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	core_v1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/net"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type AppPod struct {
	AppInfo
	pod *DaprPod
}

type AppInfo struct {
	AppID     string `csv:"APP ID"      json:"app_id"        yaml:"app_id"`
	HTTPPort  int    `csv:"HTTP PORT"   json:"http_port"     yaml:"http_port"`
	GRPCPort  int    `csv:"GRPC PORT"   json:"grpc_port"     yaml:"grpc_port"`
	AppPort   int    `csv:"APP PORT"    json:"app_port"      yaml:"app_port"`
	PodName   string `csv:"POD NAME"    json:"pod_name"      yaml:"pod_name"`
	Namespace string `csv:"NAMESPACE"   json:"namespace"    yaml:"namespace"`
	Age       string `csv:"AGE"      json:"age"     yaml:"age"`
	Created   string `csv:"CREATED"  json:"created" yaml:"created"`
	Version   string `csv:"VERSION"  json:"version" yaml:"version"`
}

type (
	DaprPod     core_v1.Pod
	DaprAppList []*AppPod
)

func GetAppPod(client k8s.Interface, appID string) (*AppPod, error) {
	list, err := ListAppInfos(client, appID)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("%s not found", appID)
	}
	app := list[0]
	return app, nil
}

// ListAppInfos ListPluginPods outputs plugins list.
func ListAppInfos(client k8s.Interface, appIDs ...string) (DaprAppList, error) {
	opts := v1.ListOptions{}
	podList, err := client.CoreV1().Pods(v1.NamespaceAll).List(context.TODO(), opts)
	if err != nil {
		return nil, fmt.Errorf("err get pods list:%w", err)
	}

	fn := func(*AppPod) bool {
		return true
	}
	if len(appIDs) > 0 {
		fn = func(a *AppPod) bool {
			if a == nil {
				return false
			}
			for _, id := range appIDs {
				if id != "" && a.AppID == id {
					return true
				}
			}
			return false
		}
	}

	var l DaprAppList
	for _, p := range podList.Items {
		p := DaprPod(p)
		for _, c := range p.Spec.Containers {
			if c.Name == "daprd" {
				app := getAppInfoFromPod(&p)
				if fn(app) {
					l = append(l, app)
				}
			}
		}
	}

	return l, nil
}

func getAppInfoFromPod(p *DaprPod) (a *AppPod) {
	for _, c := range p.Spec.Containers {
		if c.Name == "daprd" {
			a = &AppPod{
				AppInfo: AppInfo{
					PodName:   p.Name,
					Namespace: p.Namespace,
				},
				pod: p,
			}
			for i, arg := range c.Args {
				if arg == "--app-port" {
					if port, err := strconv.Atoi(c.Args[i+1]); err != nil {
						continue
					} else {
						a.AppPort = port
					}
				} else if arg == "--dapr-http-port" {
					if port, err := strconv.Atoi(c.Args[i+1]); err != nil {
						continue
					} else {
						a.HTTPPort = port
					}
				} else if arg == "--dapr-grpc-port" {
					if port, err := strconv.Atoi(c.Args[i+1]); err != nil {
						continue
					} else {
						a.GRPCPort = port
					}
				} else if arg == "--app-id" {
					id := c.Args[i+1]
					a.AppID = id
				}
			}
		}
	}
	return
}

func (a *AppInfo) Request(r *rest.Request, method string, data []byte) (*rest.Request, error) {
	r = r.Namespace(a.Namespace).
		Resource("pods").
		SubResource("proxy").
		SetHeader("Content-Type", "application/json").
		Name(net.JoinSchemeNamePort("", a.PodName, fmt.Sprintf("%d", a.AppPort)))
	if data != nil {
		r = r.Body(data)
	}

	u, err := url.Parse(method)
	if err != nil {
		return nil, fmt.Errorf("error parse method %s: %w", method, err)
	}

	r = r.Suffix(u.Path)

	for k, vs := range u.Query() {
		r = r.Param(k, strings.Join(vs, ","))
	}

	return r, nil
}

func ListPlugins(client k8s.Interface) ([]*Plugin, error) {
	rudder, err := GetAppPod(client, "rudder")
	if err != nil {
		return nil, err
	}

	res, err := rudder.Request(client.CoreV1().RESTClient().Get(), "v1/plugins", nil)
	if err != nil {
		return nil, err
	}

	result := res.Do(context.TODO())
	if result.Error() != nil {
		return nil, fmt.Errorf("k8s query resutl err: %w", err)
	}
	rawbody, err := result.Raw()
	if err != nil {
		return nil, fmt.Errorf("k8s query err: %w", err)
	}
	resp := &ListResponse{}
	err = json.Unmarshal(rawbody, resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshak json to struct err: %w", err)
	}
	return resp.PluginList, nil
}

func RegisterPlugins(client k8s.Interface, pluginID string) error {
	rudder, err := GetAppPod(client, "rudder")
	if err != nil {
		return err
	}

	res, err := rudder.Request(client.CoreV1().RESTClient().Post(), "v1/plugins",
		[]byte(fmt.Sprintf(`{"id":"%s","secret":"changeme"}`, pluginID)))
	if err != nil {
		return err
	}

	ret := res.Do(context.TODO())
	if ret.Error() != nil {
		return fmt.Errorf("k8s query result err: %w", err)
	}
	statusCode := http.StatusOK
	ret.StatusCode(&statusCode)
	if statusCode != http.StatusNoContent {
		return fmt.Errorf("k8s query status_code not invaild: %d", statusCode)
	}
	return nil
}

func UnregisterPlugins(client k8s.Interface, pluginID string) (*Plugin, error) {
	rudder, err := GetAppPod(client, "rudder")
	if err != nil {
		return nil, err
	}

	res, err := rudder.Request(client.CoreV1().RESTClient().Delete(), fmt.Sprintf(`v1/plugins/%s`, pluginID),
		[]byte(fmt.Sprintf(`{"id":"%s","secret":"changeme"}`, pluginID)))
	if err != nil {
		return nil, err
	}

	ret := res.Do(context.TODO())
	if ret.Error() != nil {
		return nil, fmt.Errorf("k8s query ret err: %w", ret.Error())
	}
	raw, err := ret.Raw()
	if err != nil {
		return nil, fmt.Errorf("k8s qeury err: %w", err)
	}
	resp := DeleteResponse{}
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal json to struct err:%w", err)
	}
	return resp.Plugin, nil
}

func GetTKeelNamespace(client k8s.Interface) (string, error) {
	keel, err := GetAppPod(client, "keel")
	if err != nil {
		return "", err
	}
	return keel.Namespace, nil
}

// GetStatusAndHealthyInPodList loop through all replicas and update to Running/Healthy status only if all instances are Running and Healthy.
func GetStatusAndHealthyInPodList(appList DaprAppList) (status string, healthy string) {
	healthy = "False"
	running := true
	if len(appList) == 0 {
		return status, healthy
	}
	firstApp := appList[0]
	for _, app := range appList {
		pod := app.pod
		if len(pod.Status.ContainerStatuses) == 0 {
			status = string(pod.Status.Phase)
		} else if pod.Status.ContainerStatuses[0].State.Waiting != nil {
			status = fmt.Sprintf("Waiting (%s)", pod.Status.ContainerStatuses[0].State.Waiting.Reason)
		} else if firstApp.pod.Status.ContainerStatuses[0].State.Terminated != nil {
			status = "Terminated"
		}

		if len(pod.Status.ContainerStatuses) == 0 ||
			pod.Status.ContainerStatuses[0].State.Running == nil {
			running = false
			break
		}

		if pod.Status.ContainerStatuses[0].Ready {
			healthy = "True"
		}
	}
	if running {
		status = "Running"
	}
	return status, healthy
}

func GroupByAppID(l DaprAppList) map[string]DaprAppList {
	ret := make(map[string]DaprAppList)
	for _, c := range l {
		id := c.AppID
		g, ok := ret[id]
		if !ok {
			g = DaprAppList{}
		}
		g = append(g, c)
		ret[id] = g
	}
	return ret
}
