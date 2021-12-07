package kubernetes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	k8s "k8s.io/client-go/kubernetes"
)

func TenantCreate(title, adminName, adminPW string) error {
	if len(title) == 0 {
		return errors.New("title param nil")
	}
	client, err := Client()
	if err != nil {
		return err
	}

	tenant := &TenantCreateIn{Title: title}
	if len(adminName) != 0 && len(adminPW) != 0 {
		admin := TenantAdmin{UserName: adminName, Password: adminPW}
		tenant.Admin = admin
	}

	return CreateTenant(client, tenant)
}

func CreateTenant(client k8s.Interface, tenant *TenantCreateIn) error {
	rudder, err := AppPod(client, "rudder")
	if err != nil {
		return err
	}
	body, _ := json.Marshal(tenant)
	res := rudder.App().Request(client.CoreV1().RESTClient().Post()).
		Suffix("v1/tenants").Body(body)

	ret := res.Do(context.TODO())
	if ret.Error() != nil {
		return fmt.Errorf("k8s query result err: %w", ret.Error())
	}
	_, err = ret.Raw()
	if err != nil {
		return fmt.Errorf("do k8s query err: %w", err)
	}

	return nil
}

func TenantList() ([]Tenant, error) {
	client, err := Client()
	if err != nil {
		return nil, err
	}
	rudder, err := AppPod(client, "rudder")
	if err != nil {
		return nil, err
	}
	res := rudder.App().Request(client.CoreV1().RESTClient().Get()).
		Suffix("v1/tenants")

	ret := res.Do(context.TODO())
	if ret.Error() != nil {
		return nil, fmt.Errorf("k8s query result err: %w", err)
	}
	raw, err := ret.Raw()
	if err != nil {
		return nil, fmt.Errorf("do k8s query err: %w", err)
	}
	resp := TenantListResponse{}
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		err = fmt.Errorf("unmarshal err :%w", err)
	}
	return resp.Data, err
}

type TenantCreateIn struct {
	Title string      `json:"title"`
	Admin TenantAdmin `json:"admin"`
}

type TenantAdmin struct {
	UserName string `json:"username"` //nolint
	Password string `json:"password"` //nolint
}

type TenantCreateResponse struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data TenantCreateResp `json:"data"`
}

type TenantListResponse struct {
	Ret  int      `json:"ret"`
	Msg  string   `json:"msg"`
	Data []Tenant `json:"data"`
}

type TenantListData struct {
	TenantList []TenantCreateResp `json:"tenant_list"`
}

type TenantCreateResp struct {
	TenantID int         `json:"tenant_id"`
	Title    string      `json:"title"`
	Admin    TenantAdmin `json:"admin"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Tenant struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Remark string `json:"remark"`
}
