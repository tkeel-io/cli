package kubernetes

import (
	"encoding/json"
	"fmt"
	"net/http"

	terrors "github.com/tkeel-io/kit/errors"

	"github.com/pkg/errors"
	"github.com/tkeel-io/kit/result"
	repoAPI "github.com/tkeel-io/tkeel/api/repo/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	_listReposMethodFormat  = "apis/rudder/v1/repos"
	_addRepoMethodFormat    = "apis/rudder/v1/repos/%s"
	_deleteRepoMethodFormat = "apis/rudder/v1/repos/%s"
)

type RepoListOutput struct {
	Name   string `csv:"REPO NAME"`
	Remote string `csv:"REMOTE"`
}

type AddRepoRequest struct {
	URL string `json:"url"`
}

func ListRepo() ([]RepoListOutput, error) {
	token, err := getAdminToken()
	if err != nil {
		return nil, errors.Wrap(err, "error get token")
	}

	resp, err := InvokeByPortForward(_pluginKeel, _listReposMethodFormat, nil, http.MethodGet, setAuthenticate(token))
	if err != nil {
		return nil, errors.Wrap(err, "error invoke")
	}

	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return nil, errors.Wrap(err, "error unmarshal")
	}

	if r.Code != terrors.Success.Reason {
		return nil, errors.Wrap(errors.New(r.Code), "error response code")
	}

	listResponse := repoAPI.ListRepoResponse{}
	err = r.Data.UnmarshalTo(&listResponse)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshal response")
	}

	var list = make([]RepoListOutput, 0, len(listResponse.Repos))
	for _, repo := range listResponse.Repos {
		list = append(list, RepoListOutput{repo.Name, repo.Url})
	}
	return list, nil
}

func AddRepo(name, url string) error {
	token, err := getAdminToken()
	if err != nil {
		return errors.Wrap(err, "get token error")
	}
	method := fmt.Sprintf(_addRepoMethodFormat, name)
	req := AddRepoRequest{URL: url}
	data, err := json.Marshal(req)
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

func DeleteRepo(name string) error {
	method := fmt.Sprintf(_deleteRepoMethodFormat, name)
	token, err := getAdminToken()
	if err != nil {
		return errors.Wrap(err, "get admin token error")
	}
	resp, err := InvokeByPortForward(_pluginKeel, method, nil, http.MethodDelete, setAuthenticate(token))
	if err != nil {
		return errors.Wrap(err, "invoke error")
	}

	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return errors.Wrap(err, "error unmarshal")
	}

	if r.Code != terrors.Success.Reason {
		return errors.New("response error: " + r.Msg)
	}

	return nil
}
