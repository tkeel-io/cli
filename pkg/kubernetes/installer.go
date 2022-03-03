package kubernetes

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	terrors "github.com/tkeel-io/kit/errors"
	"github.com/tkeel-io/kit/result"
	repoApi "github.com/tkeel-io/tkeel/api/repo/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	_installerListFormat    = "apis/rudder/v1/repos/%s/installers"
	_installerListAllFormat = "apis/rudder/v1/repos/installers"
	_installerInfoFormat    = "apis/rudder/v1/repos/%s/installers/%s/%s"
)

type InstallerListOutPut struct {
	Name    string `csv:"NAME"`
	Version string `csv:"VERSION"`
	Repo    string `csv:"REPO"`
	Status  string `csv:"STATUS"`
}

func InstallerList(repo string) ([]InstallerListOutPut, error) {
	token, err := getAdminToken()
	if err != nil {
		return nil, errors.Wrap(err, "error getting admin token")
	}
	method := fmt.Sprintf(_installerListFormat, repo)

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

	response := repoApi.ListRepoInstallerResponse{}
	err = r.Data.UnmarshalTo(&response)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshal")
	}

	var list = make([]InstallerListOutPut, 0, len(response.BriefInstallers))
	for _, installer := range response.BriefInstallers {
		list = append(list, InstallerListOutPut{installer.Name, installer.Version, installer.Repo, repoApi.InstallerState_name[int32(installer.State)]})
	}
	return list, nil
}

func InstallerListAll() ([]InstallerListOutPut, error) {

	token, err := getAdminToken()
	if err != nil {
		return nil, errors.Wrap(err, "error getting admin token")
	}
	method := _installerListAllFormat

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

	response := repoApi.ListAllRepoInstallerResponse{}
	err = r.Data.UnmarshalTo(&response)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshal")
	}

	var list = make([]InstallerListOutPut, 0, len(response.BriefInstallers))
	for _, installer := range response.BriefInstallers {
		list = append(list, InstallerListOutPut{installer.Name, installer.Version, installer.Repo, repoApi.InstallerState_name[int32(installer.State)]})
	}
	return list, nil
}

func InstallerInfo(repo, installer, version string) ([]InstallerListOutPut, error) {
	token, err := getAdminToken()
	if err != nil {
		return nil, errors.Wrap(err, "error getting admin token")
	}
	method := fmt.Sprintf(_installerInfoFormat, repo, installer, version)

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

	response := repoApi.GetRepoInstallerResponse{}
	err = r.Data.UnmarshalTo(&response)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshal")
	}

	var list = make([]InstallerListOutPut, 0, 1)
	list = append(list, InstallerListOutPut{response.Installer.Name, response.Installer.Version, response.Installer.Repo, repoApi.InstallerState_name[int32(response.Installer.State)]})
	return list, nil
}
