package kubernetes

import (
	"fmt"
	"github.com/dapr/cli/pkg/age"
	dapr "github.com/dapr/cli/pkg/kubernetes"
	"os"
	"strings"
	"sync"

	"errors"

	"github.com/dapr/cli/pkg/print"
	k8s "k8s.io/client-go/kubernetes"
)

type StatusClient struct {
	client k8s.Interface
}

// StatusOutput represents the status of a named TKeel resource.
type StatusOutput struct {
	Name         string `csv:"NAME"`
	Namespace    string `csv:"NAMESPACE"`
	Healthy      string `csv:"HEALTHY"`
	Status       string `csv:"STATUS"`
	PluginStatus string `csv:"PLUGINSTATUS"`
	Replicas     int    `csv:"REPLICAS"`
	Version      string `csv:"VERSION"`
	Age          string `csv:"AGE"`
	Created      string `csv:"CREATED"`
}

// Create a new k8s client for status commands.
func NewStatusClient() (*StatusClient, error) {
	clientset, err := Client()
	if err != nil {
		return nil, err
	}
	return &StatusClient{
		client: clientset,
	}, nil
}

// List status for TKeel resources.
func (s *StatusClient) Status() ([]StatusOutput, error) {
	client := s.client
	if client == nil {
		return nil, errors.New("kubernetes client not initialized")
	}

	namespace, err := GetTKeelNameSpace(client)
	if client == nil {
		return nil, err
	}

	tKeelPlugins, err := ListPlugins(client, namespace)
	if err != nil {
		return nil, err
	}

	pluginsMap := make(map[string]Plugin)
	for _, plugin := range tKeelPlugins {
		pluginsMap[plugin.PluginId] = plugin
	}

	daprAppStatus, err := s.daprAppsStatus()
	if err != nil {
		return nil, err
	}

	list := make([]StatusOutput, 0, len(daprAppStatus))
	for _, app := range daprAppStatus {
		status := "UNKNOWN"
		if plugin, ok := pluginsMap[app.Name]; ok {
			status = plugin.Status
		}
		list = append(list, StatusOutput{
			app.Name,
			app.Namespace,
			app.Healthy,
			app.Status,
			status,
			app.Replicas,
			app.Version,
			app.Age,
			app.Created,
		})
	}
	return list, nil
}

func (s *StatusClient) daprAppsStatus() ([]StatusOutput, error) {
	client := s.client
	if client == nil {
		return nil, errors.New("kubernetes client not initialized")
	}

	daprApps, err := dapr.List()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(len(daprApps))

	m := sync.Mutex{}
	statuses := []StatusOutput{}

	for _, app := range daprApps {
		go func(label string) {
			defer wg.Done()
			// Query all namespaces for TKeel pods.
			p, err := ListPodsInterface(client, map[string]string{
				"app": label,
			})
			if err != nil {
				print.WarningStatusEvent(os.Stdout, "Failed to get status for %s: %s", label, err.Error())
				return
			}

			if len(p.Items) == 0 {
				return
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

			s := StatusOutput{
				Name:      label,
				Namespace: namespace,
				Created:   created,
				Age:       age,
				Status:    status,
				Version:   version,
				Healthy:   healthy,
				Replicas:  replicas,
			}

			m.Lock()
			statuses = append(statuses, s)
			m.Unlock()
		}(app.AppID)
	}

	wg.Wait()
	return statuses, nil
}
