package kubernetes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tkeel-io/kit/result"
	repoAPI "github.com/tkeel-io/tkeel/api/repo/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	_listReposMethodFormat  = "v1/repos"
	_addRepoMethodFormat    = "v1/repos/%s"
	_deleteRepoMethodFormat = "v1/repos/%s"
)

type RepoListOutput struct {
	Name   string `csv:"REPO NAME"`
	Remote string `csv:"REMOTE"`
}

func ListRepo() ([]RepoListOutput, error) {
	token, err := getAdminToken()
	if err != nil {
		return nil, errors.Wrap(err, "error getting admin token")
	}

	resp, err := InvokeByPortForward(_pluginRudder, _listReposMethodFormat, nil, http.MethodGet, InvokeSetHTTPHeader("Authorization", token))
	if err != nil {
		return nil, errors.Wrap(err, "error invoke")
	}

	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return nil, errors.Wrap(err, "error unmarsh")
	}

	if r.Code != http.StatusOK {
		return nil, errors.New("response error: unexpected status code")
	}

	listResponse := repoAPI.ListRepoResponse{}
	err = r.Data.UnmarshalTo(&listResponse)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshal")
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
	data, err := json.Marshal(repoAPI.CreateRepoRequest{Url: url})
	if err != nil {
		return errors.Wrap(err, "json marshal error")
	}

	resp, err := InvokeByPortForward(_pluginRudder, method, data, http.MethodPost, InvokeSetHTTPHeader("Authorization", token))
	if err != nil {
		return errors.Wrap(err, "invoke by port forward error")
	}
	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return errors.Wrap(err, "error unmarshal")
	}

	if r.Code != http.StatusOK {
		return errors.New("response error: unexpected status code")
	}

	return nil
}

func DeleteRepo(name string) error {
	method := fmt.Sprintf(_deleteRepoMethodFormat, name)
	token, err := getAdminToken()
	if err != nil {
		return errors.Wrap(err, "get admin token error")
	}
	resp, err := InvokeByPortForward(_pluginRudder, method, nil, http.MethodDelete, InvokeSetHTTPHeader("Authorization", token))
	if err != nil {
		return errors.Wrap(err, "invoke error")
	}

	var r = &result.Http{}
	if err = protojson.Unmarshal([]byte(resp), r); err != nil {
		return errors.Wrap(err, "error unmarshal")
	}

	if r.Code != http.StatusOK {
		return errors.New("response error: unexpected status code")
	}

	return nil
}
