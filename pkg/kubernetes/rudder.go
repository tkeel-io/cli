package kubernetes

import (
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/tkeel-io/cli/pkg/print"
	"github.com/tkeel-io/kit/result"
	oauth2 "github.com/tkeel-io/tkeel/api/oauth2/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	_pluginRudder     = "rudder"
	_adminLoginMethod = "apis/rudder/v1/oauth2/admin"
	_TokenFile        = ".token"
	_tkeelRudderDir   = ".tkeel/rudder"
)

func AdminLogin(password string) (token string, err error) {
	u, err := url.Parse(_adminLoginMethod)
	if err != nil {
		return "", errors.Wrap(err, "parse admin login method error")
	}
	val := u.Query()
	val.Set("password", password)
	u.RawQuery = val.Encode()

	resp, err := Invoke(_pluginRudder, u.String(), nil, http.MethodGet)
	if err != nil {
		return "", errors.Wrap(err, "invoking admin login err")
	}
	token, err = parseToken(resp)
	if err != nil {
		return "", errors.Wrap(err, "parse token err")
	}

	f, err := openRudderTokenFile()
	if err != nil {
		return "", errors.Wrap(err, "open rudder token failed")
	}
	if _, err = f.WriteString(token); err != nil {
		return "", errors.Wrap(err, "write token to file failed")
	}

	return token, nil
}

func openRudderTokenFile() (*os.File, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		print.FailureStatusEvent(os.Stdout, "create tkeel rudder token file, get user homedir failed")
		return nil, errors.Wrap(err, "get user home dir failed")
	}
	rudderTokenFile := path.Join(homedir, _tkeelRudderDir, _TokenFile)
	_, err = os.Stat(rudderTokenFile)
	if err != nil {
		err = os.MkdirAll(path.Join(homedir, _tkeelRudderDir), os.ModeDir|os.ModePerm)
		if err != nil {
			return nil, errors.Wrap(err, "create tkeel rudder dir failed")
		}
	}

	f, err := os.OpenFile(rudderTokenFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "open tkeel rudder token file error")
	}

	return f, nil
}

func parseToken(body string) (string, error) {
	tokenResponse := oauth2.IssueTokenResponse{}

	var r = &result.Http{}
	if err := protojson.Unmarshal([]byte(body), r); err != nil {
		return "", errors.Wrap(err, "unmarshal response context error")
	}

	if r.Code != http.StatusOK {
		print.FailureStatusEvent(os.Stdout, "invalid response code")
		print.FailureStatusEvent(os.Stdout, "response context: %s", body)
		return "", errors.New("invalid response")
	}

	if err := r.Data.UnmarshalTo(&tokenResponse); err != nil {
		return "", errors.Wrap(err, "unmarshal response data to token error")
	}

	return tokenResponse.AccessToken, nil
}
