package kubernetes

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/dapr/cli/pkg/age"
	dapr "github.com/dapr/cli/pkg/kubernetes"
	"github.com/dapr/cli/pkg/print"
	k8score "k8s.io/api/core/v1"
	k8s "k8s.io/client-go/kubernetes"
)

type StatusClient struct {
	client k8s.Interface
}

// StatusOutput represents the status of a named tKeel resource.
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

// List status for tKeel resources.
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
		pluginsMap[plugin.PluginID] = plugin
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
	if s.client == nil {
		return nil, errors.New("kubernetes client not initialized")
	}

	daprApps, err := dapr.List()
	if err != nil {
		return nil, fmt.Errorf("err dapr do list: %w", err)
	}

	var (
		wg sync.WaitGroup
		m  sync.Mutex
	)
	statuses := make([]StatusOutput, 0, len(daprApps))
	wg.Add(len(daprApps))
	for _, app := range daprApps {
		go func(label string) {
			defer wg.Done()
			// Query all namespaces for tKeel pods.
			podList, err := ListPodsInterface(s.client, map[string]string{
				"app": label,
			})
			if err != nil {
				print.WarningStatusEvent(os.Stdout, "Failed to get status for %s: %s", label, err.Error())
				return
			}

			if len(podList.Items) == 0 {
				return
			}
			firstPod := podList.Items[0]
			replicas := len(podList.Items)
			namespace, podAge, created, version := extractPod(firstPod)
			status, healthy := GetStatusAndHealthyInPodList(podList)

			m.Lock()
			statuses = append(statuses, StatusOutput{
				Name:      label,
				Namespace: namespace,
				Created:   created,
				Age:       podAge,
				Status:    status,
				Version:   version,
				Healthy:   healthy,
				Replicas:  replicas,
			})
			m.Unlock()
		}(app.AppID)
	}

	wg.Wait()
	return statuses, nil
}

// GetStatusAndHealthyInPodList loop through all replicas and update to Running/Healthy status only if all instances are Running and Healthy.
func GetStatusAndHealthyInPodList(podList *k8score.PodList) (status string, healthy string) {
	healthy = "False"
	running := true
	if len(podList.Items) == 0 {
		return status, healthy
	}
	firstPod := podList.Items[0]
	for _, pod := range podList.Items {
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

func extractPod(firstPod k8score.Pod) (string, string, string, string) {
	image := firstPod.Spec.Containers[0].Image
	namespace := firstPod.GetNamespace()
	podAge := age.GetAge(firstPod.CreationTimestamp.Time)
	created := firstPod.CreationTimestamp.Format("2006-01-02 15:04.05")
	version := image[strings.IndexAny(image, ":")+1:]
	return namespace, podAge, created, version
}
