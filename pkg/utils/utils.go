package utils

import (
	"os"
	"strings"
)

// ParseInstallArg parse the first arg, get repo, plugin and version information.
// More efficient and concise support for both formats：
// url style install target plugin: https://tkeel-io.github.io/helm-charts/A@version
// short style install official plugin： tkeel/B@version or C@version.
func ParseInstallArg(arg string, defaultRepo string) (repo, name, version string) {
	version = ""
	name = arg

	if sp := strings.Split(arg, "@"); len(sp) == 2 {
		name, version = sp[0], sp[1]
	}

	if version != "" && version[0] == 'v' {
		version = version[1:]
	}

	repo = defaultRepo
	if spi := strings.LastIndex(name, "/"); spi != -1 {
		repo, name = name[:spi], name[spi+1:]
		if repo == "" || strings.EqualFold(repo, "tkeel") {
			repo = defaultRepo
			return
		}
	}
	return
}

func GetRealPath(path string) (string, error) {
	if path != "" && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = strings.Replace(path, "~", home, 1)
	}
	return path, nil
}
