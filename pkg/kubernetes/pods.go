package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dapr/cli/pkg/age"
	core_v1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/net"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type DaprPod core_v1.Pod

func (p *DaprPod) App() App {
	ap := App{
		PodName:   p.Name,
		NameSpace: p.Namespace,
	}
	for _, c := range p.Spec.Containers {
		if c.Name == "daprd" {
			for i, a := range c.Args {
				if a == "--app-port" {
					port := c.Args[i+1]
					ap.AppPort = port
				} else if a == "--app-id" {
					id := c.Args[i+1]
					ap.AppID = id
				}
			}
			ap.Created = p.CreationTimestamp.Format("2006-01-02 15:04.05")
			ap.Age = age.GetAge(p.CreationTimestamp.Time)

			image := p.Spec.Containers[0].Image
			ap.Version = image[strings.IndexAny(image, ":")+1:]
		}
	}

	return ap
}

type DaprPodList []DaprPod

func (l DaprPodList) GroupByAppID() map[string]DaprPodList {
	ret := make(map[string]DaprPodList)
	for _, c := range l {
		id := c.App().AppID
		g, ok := ret[id]
		if !ok {
			g = DaprPodList{}
		}
		g = append(g, c)
		ret[id] = g
	}
	return ret
}

type App struct {
	AppID     string `csv:"APP ID"   json:"app_id"   yaml:"app_id"`
	AppPort   string `csv:"APP PORT" json:"app_port" yaml:"app_port"`
	PodName   string `csv:"POD NAME" json:"pod_name" yaml:"pod_name"`
	NameSpace string `csv:"NAMESPACE" json:"name_space" yaml:"name_space"`
	Age       string `csv:"AGE"      json:"age"     yaml:"age"`
	Created   string `csv:"CREATED"  json:"created" yaml:"created"`
	Version   string `csv:"VERSION"  json:"version" yaml:"version"`
}

func (a App) Request(r *rest.Request) *rest.Request {
	r.Namespace(a.NameSpace).
		Resource("pods").
		SubResource("proxy").
		SetHeader("Content-Type", "application/json").
		Name(net.JoinSchemeNamePort("", a.PodName, a.AppPort))
	return r
}

// List outputs plugins.
func ListPluginPods(client k8s.Interface, appIDs ...string) (DaprPodList, error) {
	opts := v1.ListOptions{}
	podList, err := client.CoreV1().Pods(v1.NamespaceAll).List(context.TODO(), opts)
	if err != nil {
		return nil, fmt.Errorf("err get pods list:%w", err)
	}

	fn := func(dp DaprPod) bool {
		return true
	}
	if len(appIDs) > 0 {
		fn = func(dp DaprPod) bool {
			for _, id := range appIDs {
				if dp.App().AppID == id {
					return true
				}
			}
			return false
		}
	}

	l := []DaprPod{}
	for _, p := range podList.Items {
		p := DaprPod(p)
	FindLoop:
		for _, c := range p.Spec.Containers {
			if c.Name == "daprd" {
				if fn(p) {
					l = append(l, p)
				}
				break FindLoop
			}
		}
	}

	return l, nil
}

func AppPod(client k8s.Interface, appID string) (*DaprPod, error) {
	pods, err := ListPluginPods(client, appID)
	if err != nil {
		return nil, err
	}
	if len(pods) == 0 {
		return nil, fmt.Errorf("%s not found", appID)
	}
	appPod := pods[0]
	return &appPod, nil
}

func ListPlugins(client k8s.Interface) ([]*Plugin, error) {
	rudder, err := AppPod(client, "rudder")
	if err != nil {
		return nil, err
	}

	res := rudder.App().Request(client.CoreV1().RESTClient().Get()).
		Suffix("v1/plugins")
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
	rudder, err := AppPod(client, "rudder")
	if err != nil {
		return err
	}

	res := rudder.App().Request(client.CoreV1().RESTClient().Post()).
		Suffix("v1/plugins").
		Body([]byte(fmt.Sprintf(`{"id":"%s","secret":"changeme"}`, pluginID)))

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

// UnregisterPlugins
// curl -XDELETE "http://192.168.123.11:30777/v1/plugins/keel-echo"
func UnregisterPlugins(client k8s.Interface, pluginID string) (*Plugin, error) {
	rudder, err := AppPod(client, "rudder")
	if err != nil {
		return nil, err
	}

	res := rudder.App().Request(client.CoreV1().RESTClient().Delete()).
		Suffix(fmt.Sprintf(`v1/plugins/%s`, pluginID))

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

func GetTKeelNameSpace(client k8s.Interface) (string, error) {
	pods, err := ListPluginPods(client, "keel")
	if err != nil {
		return "", err
	}
	if len(pods) == 0 {
		return "", fmt.Errorf("tKeel not initialized")
	}
	return pods[0].Namespace, nil
}

// GetStatusAndHealthyInPodList loop through all replicas and update to Running/Healthy status only if all instances are Running and Healthy.
func GetStatusAndHealthyInPodList(podList DaprPodList) (status string, healthy string) {
	healthy = "False"
	running := true
	if len(podList) == 0 {
		return status, healthy
	}
	firstPod := podList[0]
	for _, pod := range podList {
		if len(pod.Status.ContainerStatuses) == 0 {
			status = string(pod.Status.Phase)
		} else if pod.Status.ContainerStatuses[0].State.Waiting != nil {
			status = fmt.Sprintf("Waiting (%s)", pod.Status.ContainerStatuses[0].State.Waiting.Reason)
		} else if firstPod.Status.ContainerStatuses[0].State.Terminated != nil {
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
