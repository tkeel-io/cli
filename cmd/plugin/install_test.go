package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInstallArg(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  struct {
			repo, plugin, version string
		}
	}{
		{"url test", "http://github.com/install/pluginname", struct{ repo, plugin, version string }{"http://github.com/install", "pluginname", "latest"}},
		{"short test", "tkeel/pluginname@v0.2.0", struct{ repo, plugin, version string }{officialRepo, "pluginname", "0.2.0"}},
		{"short test no repo", "pluginname", struct{ repo, plugin, version string }{repo: officialRepo, plugin: "pluginname", version: "latest"}},
		{"invalid test", "test/plugin", struct{ repo, plugin, version string }{repo: "test", plugin: "plugin", version: "latest"}},
		{"- test", "test/hello-plugin", struct{ repo, plugin, version string }{repo: "test", plugin: "hello-plugin", version: "latest"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo, plugin, version := parseInstallArg(test.input)
			assert.Equal(t, repo, test.want.repo)
			assert.Equal(t, plugin, test.want.plugin)
			assert.Equal(t, version, test.want.version)
		})
	}
}
