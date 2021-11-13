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
	windowsOS  = "windows"
	windowsTmp = "C:\\WINDOWS\\TEMP"
	unixTmp    = "/tmp"

	templateZipDownloadURL = "https://github.com/tkeel-io/tkeel-template-go/archive/refs/heads/main.zip"
	githubRepoURL          = "https://github.com/tkeel-io/tkeel-template-go.git"

	downloadedZipFilename   = "template.zip"
	defaultUnzippedFilename = "tkeel-template-go-main"
	gitConfigDir            = ".git"
)

var (
	_gitMode = false
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
		dirName := "my_plugin"
		if len(args) != 0 {
			dirName = args[0]
		}

		wd, err := os.Getwd()
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
		}
		targetDest := path.Join(wd, dirName)
		if i, err := os.Stat(targetDest); err == nil {
			print.FailureStatusEvent(os.Stdout, "The '%s' is exited!", i.Name())
			return
		}
		if _gitMode {
			createByGit(targetDest, dirName)
			return
		}

		createByUnzip(targetDest, wd)
	},
}

func init() {
	Create.Flags().BoolP("help", "h", false, "Print this help message")
	Create.Flags().BoolVarP(&_gitMode, "git", "", false, "use git to download this template")
	PluginCmd.AddCommand(Create)
}

func createByGit(targetDest, dirname string) {
	gitCloneCMD := exec.Command("git", "clone", githubRepoURL, dirname)
	gitCloneCMD.Stdout = os.Stdout
	gitCloneCMD.Stderr = os.Stdout
	if err := gitCloneCMD.Run(); err != nil {
		print.FailureStatusEvent(os.Stdout, "Git clone err:"+err.Error())
		return
	}
	if err := os.RemoveAll(path.Join(targetDest, gitConfigDir)); err != nil {
		print.FailureStatusEvent(os.Stdout, "Remove .git form template err:"+err.Error())
	}
	print.SuccessStatusEvent(os.Stdout, "Success!! Plugin template is created.")
}

func createByUnzip(targetDest, workingDir string) {
	tempDir := unixTmp
	if runtime.GOOS == windowsOS {
		tempDir = windowsTmp
	}

	print.InfoStatusEvent(os.Stdout, "Downloading template...")
	tmpDest := path.Join(tempDir, downloadedZipFilename)
	if err := downloadutil.Download(tmpDest, templateZipDownloadURL); err != nil {
		print.FailureStatusEvent(os.Stdout, "Template download err:"+err.Error())
		return
	}

	var (
		unzip     = "unzip"
		unzipArgs = []string{"-o", tmpDest}
	)

	if runtime.GOOS == windowsOS {
		unzip = "powershell"
		unzipArgs = []string{"-Command", "Expand-Archive", "-Path", "'" + tmpDest + "'", "-DestinationPath", "'" + workingDir + "'"}
	}

	print.InfoStatusEvent(os.Stdout, "Unpacking template...")
	unzipcmd := exec.Command(unzip, unzipArgs...)
	unzipcmd.Stderr = os.Stdout
	unzipcmd.Stdout = os.Stdout
	if err := unzipcmd.Run(); err != nil {
		print.FailureStatusEvent(os.Stdout, "Unzip err:"+err.Error())
		return
	}

	if err := os.Rename(defaultUnzippedFilename, targetDest); err != nil {
		print.FailureStatusEvent(os.Stdout, "Move err:"+err.Error())
		return
	}

	print.SuccessStatusEvent(os.Stdout, "Success!! Plugin template is created.")
}
