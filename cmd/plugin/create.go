// ------------------------------------------------------------
// Copyright 2021 The tKeel Contributors.
// Licensed under the Apache License.
// ------------------------------------------------------------

package plugin

import (
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/downloadutil"
	"github.com/tkeel-io/cli/pkg/print"
)

const (
	zipDownloadURL = "https://github.com/tkeel-io/tkeel-template-go/archive/refs/heads/main.zip"
	githubRepoURL  = "https://github.com/tkeel-io/tkeel-template-go.git"

	downloadedZipFilename   = "template.zip"
	defaultUnzippedFilename = "tkeel-template-go-main"
	gitConfigDir            = ".git"
)

var (
	_gitMode = false
	_tempDir = "/tmp"
)

var Create = &cobra.Command{
	Use:   "create [dir]",
	Short: "Create a plugin in quickstart template.",
	Example: `
# Create a plugin in quickstart template.
tkeel plugin create 
tkeel plugin create plugin_name
`,
	Run: func(cmd *cobra.Command, args []string) {
		name := "my_plugin"
		if len(args) != 0 {
			name = args[0]
		}

		wd, err := os.Getwd()
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
		}
		targetDir := path.Join(wd, name)

		if _gitMode {
			gitCloneCMD := exec.Command("git", "clone", githubRepoURL, name)
			gitCloneCMD.Stdout = os.Stdout
			gitCloneCMD.Stderr = os.Stdout
			if err = gitCloneCMD.Run(); err != nil {
				print.FailureStatusEvent(os.Stdout, "Git clone err:"+err.Error())
				return
			}
			if err = os.RemoveAll(path.Join(targetDir, gitConfigDir)); err != nil {
				print.FailureStatusEvent(os.Stdout, "Remove .git form template err:"+err.Error())
			}
			print.SuccessStatusEvent(os.Stdout, "Success!! Plugin template is created.")
			return
		}

		print.InfoStatusEvent(os.Stdout, "Downloading template...")
		if runtime.GOOS == "windows" {
			_tempDir = "C:\\WINDOWS\\TEMP"
		}
		tmpDest := path.Join(_tempDir, downloadedZipFilename)
		err = downloadutil.Download(tmpDest, zipDownloadURL)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, "Template download err:"+err.Error())
			return
		}

		print.InfoStatusEvent(os.Stdout, "Unpacking template...")
		unzipcmd := exec.Command("unzip", "-o", tmpDest)
		unzipcmd.Stderr = os.Stdout
		unzipcmd.Stdout = os.Stdout
		if err := unzipcmd.Run(); err != nil {
			print.FailureStatusEvent(os.Stdout, "Unzip err:"+err.Error())
			return
		}

		if err := os.Rename(defaultUnzippedFilename, targetDir); err != nil {
			print.FailureStatusEvent(os.Stdout, "Move err:"+err.Error())
			return
		}

		print.SuccessStatusEvent(os.Stdout, "Success!! Plugin template is created.")
	},
}

func init() {
	Create.Flags().BoolP("help", "h", false, "Print this help message")
	Create.Flags().BoolVarP(&_gitMode, "git", "", false, "use git to download this template")
	PluginCmd.AddCommand(Create)
}
