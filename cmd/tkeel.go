// ------------------------------------------------------------
// Copyright 2021 The tKeel Contributors.
// Licensed under the Apache License.
// ------------------------------------------------------------

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tkeel-io/cli/cmd/plugin"
	"github.com/tkeel-io/cli/cmd/tenant"
	"github.com/tkeel-io/cli/pkg/api"
	"github.com/tkeel-io/cli/pkg/print"
)

var RootCmd = &cobra.Command{
	Use:   "tkeel",
	Short: "Keel CLI",
	Long: `
      __             __
     / /_____  ___  / /
    / //_/ _ \/ _ \/ /
   / ,< /  __/  __/ /
  /_/|_|\___/\___/_/
									   
===============================
Things Keel Platform`,
}

var logAsJSON bool

// Execute adds all child commands to the root command.
func Execute(version, apiVersion string) {
	RootCmd.Version = version
	api.PlatformAPIVersion = apiVersion

	cobra.OnInitialize(initConfig)

	setVersion()

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func setVersion() {
	template := fmt.Sprintf("Keel CLI version: %s \n", RootCmd.Version)
	RootCmd.SetVersionTemplate(template)
}

func initConfig() {
	if logAsJSON {
		print.EnableJSONFormat()
	}

	viper.SetEnvPrefix("keel")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&logAsJSON, "log-as-json", "", false, "Log output in JSON format")

	RootCmd.AddCommand(plugin.PluginCmd)
	RootCmd.AddCommand(tenant.TenantCmd)
}
