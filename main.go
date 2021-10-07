// ------------------------------------------------------------
// Copyright 2021 The TKeel Contributors.
// Licensed under the Apache License.
// ------------------------------------------------------------

package main

import (
	"github.com/tkeel-io/cli/cmd"
)

// Values for version and apiVersion are injected by the build.
var (
	version    = ""
	apiVersion = "1.0"
)

func main() {
	cmd.Execute(version, apiVersion)
}
