package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
)

func GetConfig(client k8s.Interface, name, namespace string) (*corev1.ConfigMap, error) {
	return client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, v1.GetOptions{})
}
