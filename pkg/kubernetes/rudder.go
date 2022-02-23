package kubernetes

import (
	"encoding/base64"
	"fmt"
	terrors "github.com/tkeel-io/kit/errors"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/tkeel-io/cli/fileutil"
	"github.com/tkeel-io/kit/result"
	oauth2 "github.com/tkeel-io/tkeel/api/oauth2/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	_pluginRudder     = "rudder"
	_pluginKeel       = "keel"
	_adminLoginMethod = "v1/oauth2/admin"
)

func AdminLogin(password string) (token string, err error) {
	password = base64.StdEncoding.EncodeToString([]byte(password))
	u, err := url.Parse(_adminLoginMethod)
	if err != nil {
		return "", errors.Wrap(err, "parse admin login method error")
	}
	val := u.Query()
	val.Set("password", password)
	u.RawQuery = val.Encode()

	resp, err := InvokeByPortForward(_pluginRudder, u.String(), nil, http.MethodGet)
	if err != nil {
		return "", errors.Wrap(err, "invoking admin login err")
	}
	token, err = getToken(resp)
	if err != nil {
		return "", errors.Wrap(err, "get token err")
	}

	f, err := fileutil.LocateAdminToken()
	if err != nil {
		return "", errors.Wrap(err, "open rudder token failed")
	}
	if _, err = f.WriteString(token); err != nil {
		return "", errors.Wrap(err, "write token to file failed")
	}

	return token, nil
}

func getToken(body string) (string, error) {
	tokenResponse := oauth2.IssueTokenResponse{}

	var r = &result.Http{}
	if err := protojson.Unmarshal([]byte(body), r); err != nil {
		return "", errors.Wrap(err, "unmarshal response context error")
	}

	if r.Code == "io.tkeel.rudder.api.oauth2.v1.OAUTH2_ERR_PASSWORD_NOT_MATCH" {
		return "", errors.New("invalid password")
	}

	if r.Code != terrors.Success.Reason {
		return "", fmt.Errorf("invalid response: %s", r.Msg)
	}

	if err := r.Data.UnmarshalTo(&tokenResponse); err != nil {
		return "", errors.Wrap(err, "unmarshal response data to token error")
	}

	return tokenResponse.AccessToken, nil
}
