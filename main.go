// ------------------------------------------------------------
// Copyright 2021 The tKeel Contributors.
// Licensed under the Apache License.
// ------------------------------------------------------------

package main

import (
	"github.com/tkeel-io/cli/cmd"
)

// Values for version and _apiVersion are injected by the build.
var (
	version     = ""
	_apiVersion = "1.0"
)

func main() {
	cmd.Execute(version, _apiVersion)
}
