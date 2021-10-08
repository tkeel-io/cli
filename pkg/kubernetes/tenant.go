package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/util/net"
	k8s "k8s.io/client-go/kubernetes"
)

func TenantCreate(tenantTitle string) (*TenantCreateResp, error) {
	clientset, err := Client()
	if err != nil {
		return nil, err
	}

	namespace, err := GetTKeelNameSpace(clientset)
	if err != nil {
		return nil, err
	}
	return CreateTenant(clientset, namespace, tenantTitle)
}

func CreateTenant(client k8s.Interface, namespace, tenantTitle string) (*TenantCreateResp, error) {
	res := client.CoreV1().RESTClient().Post().
		Namespace(namespace).
		Resource("services").
		SubResource("proxy").
		Name(net.JoinSchemeNamePort("", "keel", "8080")).
		Suffix("auth/tenant/create").
		Body([]byte(fmt.Sprintf(`{"title":"%s"}`, tenantTitle)))
	ret := res.Do(context.TODO())
	raw, err := ret.Raw()
	resp := TenantCreateResponse{}
	err = json.Unmarshal(raw, &resp)
	return &resp.Data, err
}

func TenantList() (*TenantListData, error) {
	clientset, err := Client()
	if err != nil {
		return nil, err
	}

	namespace, err := GetTKeelNameSpace(clientset)
	if err != nil {
		return nil, err
	}
	return ListTenant(clientset, namespace)
}

func ListTenant(client k8s.Interface, namespace string) (*TenantListData, error) {
	res := client.CoreV1().RESTClient().Post().
		Namespace(namespace).
		Resource("services").
		SubResource("proxy").
		Name(net.JoinSchemeNamePort("", "keel", "8080")).
		Suffix("auth/tenant/list").
		Body([]byte(`{}`))
	ret := res.Do(context.TODO())
	raw, err := ret.Raw()
	resp := TenantListResponse{}
	err = json.Unmarshal(raw, &resp)
	return &resp.Data, err
}

type TenantCreateResponse struct {
	Ret  int              `json:"ret"`
	Msg  string           `json:"msg"`
	Data TenantCreateResp `json:"data"`
}
type TenantListResponse struct {
	Ret  int            `json:"ret"`
	Msg  string         `json:"msg"`
	Data TenantListData `json:"data"`
}

type TenantListData struct {
	TenantList []TenantCreateResp `json:"tenant_list"`
}
type TenantCreateResp struct {
	TenantID    string `json:"tenant_id"`
	Title       string `json:"title"`
	CreatedTime int64  `json:"created_time"`
	TenantAdmin User   `json:"tenant_admin"`
}

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Password   string `json:"password"`
	TenantID   string `json:"tenant_id"`
	Email      string `json:"email"`
	CreateTime int64  `json:"create_time"`
}
