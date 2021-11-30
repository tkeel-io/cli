/*
Copyright 2021 The tKeel Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/tkeel-io/cli/pkg/kubernetes"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/print"
)

const defaultHTTPVerb = http.MethodPost

var (
	invokeAppID     string
	invokeAppMethod string
	invokeData      string
	invokeVerb      string
	invokeDataFile  string
)

var InvokeCmd = &cobra.Command{
	Use:   "invoke",
	Short: "Invoke a method on a given tKeel plugin(application). Supported platforms: kubernetes",
	Example: `
# Invoke a sample method on target app with POST Verb
tkeel invoke --plugin-id target --method v1/sample --dao '{"key":"value"}

# Invoke a sample method on target app with GET Verb
tkeel invoke --plugin-id target --method v1/sample --verb GET
`,
	Run: func(cmd *cobra.Command, args []string) {
		bytePayload := []byte{}
		var err error
		if invokeDataFile != "" && invokeData != "" {
			print.FailureStatusEvent(os.Stdout, "Only one of --dao and --dao-file allowed in the same invoke command")
			os.Exit(1)
		}

		if invokeDataFile != "" {
			bytePayload, err = ioutil.ReadFile(invokeDataFile)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, "Error reading payload from '%s'. Error: %s", invokeDataFile, err)
				os.Exit(1)
			}
		} else if invokeData != "" {
			bytePayload = []byte(invokeData)
		}

		response, err := kubernetes.Invoke(invokeAppID, invokeAppMethod, bytePayload, invokeVerb)
		if err != nil {
			err = fmt.Errorf("error invoking plugin %s: %s", invokeAppID, err)
			print.FailureStatusEvent(os.Stdout, err.Error())
			return
		}

		if response != "" {
			fmt.Println(response)
		}

		print.SuccessStatusEvent(os.Stdout, "Plugin invoked successfully")
	},
}

func init() {
	InvokeCmd.Flags().StringVarP(&invokeAppID, "plugin-id", "p", "", "The application id to invoke")
	InvokeCmd.Flags().StringVarP(&invokeAppMethod, "method", "m", "", "The method to invoke")
	InvokeCmd.Flags().StringVarP(&invokeData, "dao", "d", "", "The JSON serialized dao string (optional)")
	InvokeCmd.Flags().StringVarP(&invokeVerb, "verb", "v", defaultHTTPVerb, "The HTTP verb to use")
	InvokeCmd.Flags().StringVarP(&invokeDataFile, "dao-file", "f", "", "A file containing the JSON serialized dao (optional)")
	InvokeCmd.Flags().BoolP("help", "h", false, "Print this help message")
	InvokeCmd.MarkFlagRequired("plugin-id")
	InvokeCmd.MarkFlagRequired("method")
	RootCmd.AddCommand(InvokeCmd)
}
