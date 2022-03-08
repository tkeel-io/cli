package kubernetes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	terrors "github.com/tkeel-io/kit/errors"
	"github.com/tkeel-io/kit/result"
	tenantApi "github.com/tkeel-io/tkeel/api/tenant/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	_listTenantUserMethodFormat   = "apis/rudder/v1/tenants/%s/users"
	_createTenantUserMethodFormat = "apis/rudder/v1/tenants/%s/users"
	_deleteTenantUserMethodFormat = "apis/rudder/v1/tenants/%s/users/%s"
	_infoTenantUserMethodFormat   = "apis/rudder/v1/tenants/%s/users/%s"
)

type UserListOutPut struct {
	ID       string `csv:"ID"`
	Username string `csv:"USERNAME"`
	TenantID string `csv:"TENANT ID"`
}

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TenantUserCreate(tenantID, username, password string) error {
	token, err := getAdminToken()
	if err != nil {
		return err
	}
	method := fmt.Sprintf(_createTenantUserMethodFormat, tenantID)
	userinfo := UserInfo{Username: username, Password: password}
	data, err := json.Marshal(userinfo)
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

	if r.Code != terrors.Success.Reason {
		return errors.New("response error: " + r.Msg)
	}

	return nil
}

// tenant user manage.
func TenantUserDelete(tenantID, userID string) error {
	token, err := getAdminToken()
	if err != nil {
		return err
	}
	method := fmt.Sprintf(_deleteTenantUserMethodFormat, tenantID, userID)
	resp, err := InvokeByPortForward(_pluginKeel, method, nil, http.MethodDelete, setAuthenticate(token))
	if err != nil {
		return errors.Wrap(err, "invoke "+method+" error")
	}

	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return errors.Wrap(err, "can't unmarshal'")
	}

	if r.Code != terrors.Success.Reason {
		return errors.New("response error: " + r.Msg)
	}

	return nil
}

func TenantUserInfo(tenantID, userID string) ([]UserListOutPut, error) {
	token, err := getAdminToken()
	if err != nil {
		return nil, errors.Wrap(err, "error getting admin token")
	}
	method := fmt.Sprintf(_infoTenantUserMethodFormat, tenantID, userID)

	resp, err := InvokeByPortForward(_pluginKeel, method, nil, http.MethodGet, setAuthenticate(token))
	if err != nil {
		return nil, errors.Wrap(err, "error invoke")
	}

	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return nil, errors.Wrap(err, "error unmarshal")
	}

	if r.Code != terrors.Success.Reason {
		return nil, errors.New("response error: " + r.Msg)
	}

	response := tenantApi.GetUserResponse{}
	err = r.Data.UnmarshalTo(&response)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshal")
	}

	var list = make([]UserListOutPut, 0, 1)
	list = append(list, UserListOutPut{response.UserId, response.Username, response.TenantId})
	return list, nil
}

func TenantUserList(tenantID string) ([]UserListOutPut, error) {
	token, err := getAdminToken()
	if err != nil {
		return nil, errors.Wrap(err, "error getting admin token")
	}
	method := fmt.Sprintf(_listTenantUserMethodFormat, tenantID)

	resp, err := InvokeByPortForward(_pluginKeel, method, nil, http.MethodGet, setAuthenticate(token))
	if err != nil {
		return nil, errors.Wrap(err, "error invoke")
	}

	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return nil, errors.Wrap(err, "error unmarshal")
	}

	if r.Code != terrors.Success.Reason {
		return nil, errors.New("response error: " + r.Msg)
	}

	response := tenantApi.ListUserResponse{}
	err = r.Data.UnmarshalTo(&response)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshal")
	}

	var list = make([]UserListOutPut, 0, len(response.Users))
	for _, user := range response.Users {
		list = append(list, UserListOutPut{user.UserId, user.Username, user.TenantId})
	}
	return list, nil
}
