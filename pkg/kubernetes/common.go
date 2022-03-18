package kubernetes

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tkeel-io/cli/fileutil"
)

func getAdminToken() (string, error) {
	f, err := fileutil.LocateAdminToken(fileutil.RWFlag())
	if err != nil {
		return "", errors.Wrap(err, "open admin token failed, run `tkeel admin login` to update token.")
	}
	tokenb := make([]byte, 512)
	n, err := f.Read(tokenb)
	if err != nil {
		return "", errors.Wrap(err, "read token failed, run `tkeel admin login` to update token.")
	}
	if n == 0 {
		return "", errors.New("invalid token, run `tkeel admin login` to update token.")
	}

	return fmt.Sprintf("Bearer %s", tokenb[:n]), nil
}

func setAuthenticate(token string) HTTPRequestOption {
	return InvokeSetHTTPHeader("Authorization", token)
}
