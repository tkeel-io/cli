<h1 align="center"> tKeel CLI </h1>
<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/tkeel-io/cli)](https://goreportcard.com/report/github.com/tkeel-io/cli)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/tkeel-io/cli)
![GitHub](https://img.shields.io/github/license/tkeel-io/cli?style=plastic)
[![GoDoc](https://godoc.org/github.com/tkeel-io/cli?status.png)](http://godoc.org/github.com/tkeel-io/cli)
</div>

🕹️ tKeel CLI is your main tool for various tasks related to tKeel Platform.

You can use it to **install** the tKeel platform, **manage plugins and users**.

👉 [中文文档](README_zh.md)

### Prerequisites

tKeel CLI can help you install the tKeel platform and help you manage the platform.

> ⚠️ tKeel currently relies on Dapr (Kubernetes mode).

* Install [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* Install [Dapr on Kubernetes](https://docs.dapr.io/operations/hosting/kubernetes/kubernetes-deploy/)

### Install

🔧 We offer a variety of installation options, you choose the one you feel most comfortable installing according to your
preference.

#### Using script to install the latest release

The components required for the entire `tKeel` platform are automatically installed for you by means of scripts we have
written.

##### Linux

Install the latest linux tKeel CLI to `/usr/local/bin`

```bash
$ wget -q https://raw.githubusercontent.com/tkeel-io/cli/master/install/install.sh -O - | /bin/bash
```

##### MacOS

Install the latest darwin tKeel CLI to `/usr/local/bin`

```bash
$ curl -fsSL https://raw.githubusercontent.com/tkeel-io/cli/master/install/install.sh | /bin/bash
```

#### From the Binary Releases

Each release of tKeel CLI includes various OSes and architectures. These binary versions can be manually downloaded and
installed.

1. Download the [tKeel CLI](https://github.com/tkeel-io/cli/releases)
2. Unpack it (e.g. tkeel_linux_amd64.tar.gz, tkeel_windows_amd64.zip)
3. Move it to your desired location.
    * For Linux/MacOS - `/usr/local/bin`
    * For Windows, create a directory and add this to your System PATH. For example create a directory called `c:\tkeel`
      and add this directory to your path, by editing your system environment variable.

### Init tKeel Platform on Kubernetes

Use the init command to initialize tKeel.

```bash
$ tkeel init
```

> For Linux users, if you run your docker cmds with sudo, you need to use "**sudo tkeel init**"

Output should look like so:

```
⌛  Making the jump to hyperspace...
ℹ️  Checking the Dapr runtime status...
↑  Deploying the tKeel Platform to your cluster... 
ℹ️  install plugins...                                                        
ℹ️  install plugins done.                                                                                                        
✅  Deploying the tKeel Platform to your cluster...
↖  Register the plugins ... 
ℹ️  Plugin<plugins>  is registered.                                                                                          
ℹ️  Plugin<keel>  is registered.                                                                                                                        
ℹ️  Plugin<auth>  is registered.                                                                                                                        
✅  Success! tKeel Platform has been installed to namespace keel-system. To verify, run `tkeel plugin list -k' in your terminal. To get started, go here: https://tkeel.io/keel-getting-started
```

### Uninstall tKeel on Kubernetes

To remove tKeel from your Kubernetes cluster, use the `uninstall` command.

```
$ tkeel uninstall
```

### Deploy plugin

You can deploy the plugin app with the Dapr. There
is [deploy-the-plugin-app Doc](https://github.com/dapr/quickstarts/tree/v1.0.0/hello-kubernetes#step-3---deploy-the-nodejs-app-with-the-dapr-sidecar)

### Manage plugins

Use the plugin command to manage plugins.

1. List plugin

```bash
$ tkeel plugin list      
NAME       NAMESPACE    HEALTHY  STATUS    PLUGINSTATUS  REPLICAS  VERSION  AGE  CREATED              
auth       keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00  
plugins    keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00  
keel       keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00
echo-demo  keel-system  False    Running   UNKNOWN       1         0.0.1    1m   2021-10-05 11:25.19  
```

2. Register plugin

```bash
$ tkeel plugin register echo-demo
✅  Success! Plugin<echo-demo> has been Registered to tKeel Platform . To verify, run `tkeel plugin list -k' in your terminal.
```

Check the status

```bash
$ tkeel plugin list              
NAME       NAMESPACE    HEALTHY  STATUS    PLUGINSTATUS  REPLICAS  VERSION  AGE  CREATED              
auth       keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00  
plugins    keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00  
keel       keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00
echo-demo  keel-system  False    Running   ACTIVE        1         0.0.1    2m   2021-10-05 11:25.19  
```

3. Delete plugin

```bash
$ tkeel plugin delete echo-demo
✅  Success! Plugin<echo-demo> has been deleted from tKeel Platform . To verify, run `tkeel plugin list -k' in your terminal.
```

### Quickstart to plugin development

We've designed the [Plugin Development Template](https://github.com/tkeel-io/tkeel-template-go) to organise your code
directly based on this template, and we've helped you to sort out the code hierarchy and provide tools for building APIs
quickly.

You'll love it.

Templates can be quickly downloaded using the `plugin create` command directly.

By default, the template is installed in the current working directory and named `my_plugin`.
> If you are a Windows user, you can use a package manager (e.g. winget, Chocolate) to install `unzip` and then use

```bash
$ tkeel plugin create 
ℹ️  Downloading template...
ℹ️  Unpacking template...
/usr/bin/unzip -o /tmp/template.zip
Archive:  /tmp/template.zip
59700ef3ee2bbe545f9a4c4c84488c8feacaab6b
   creating: tkeel-template-go-main/
  inflating: tkeel-template-go-main/.gitignore
  inflating: tkeel-template-go-main/Dockerfile
  inflating: tkeel-template-go-main/LICENSE
  inflating: tkeel-template-go-main/Makefile
 extracting: tkeel-template-go-main/README.md
   creating: tkeel-template-go-main/api/
   creating: tkeel-template-go-main/api/helloworld/
   creating: tkeel-template-go-main/api/helloworld/v1/
  inflating: tkeel-template-go-main/api/helloworld/v1/greeter.pb.go
  inflating: tkeel-template-go-main/api/helloworld/v1/greeter.proto
  inflating: tkeel-template-go-main/api/helloworld/v1/greeter_grpc.pb.go
  inflating: tkeel-template-go-main/api/helloworld/v1/greeter_http.pb.go
   creating: tkeel-template-go-main/api/openapi/
   creating: tkeel-template-go-main/api/openapi/v1/
  inflating: tkeel-template-go-main/api/openapi/v1/openapi.pb.go
  inflating: tkeel-template-go-main/api/openapi/v1/openapi.proto
  inflating: tkeel-template-go-main/api/openapi/v1/server.pb.go
  inflating: tkeel-template-go-main/api/openapi/v1/server.proto
  inflating: tkeel-template-go-main/api/openapi/v1/server_grpc.pb.go
  inflating: tkeel-template-go-main/api/openapi/v1/server_http.pb.go
   creating: tkeel-template-go-main/cmd/
   creating: tkeel-template-go-main/cmd/hello/
  inflating: tkeel-template-go-main/cmd/hello/main.go
  inflating: tkeel-template-go-main/go.mod
  inflating: tkeel-template-go-main/go.sum
   creating: tkeel-template-go-main/pkg/
   creating: tkeel-template-go-main/pkg/server/
  inflating: tkeel-template-go-main/pkg/server/grpc.go
  inflating: tkeel-template-go-main/pkg/server/http.go
   creating: tkeel-template-go-main/pkg/service/
 extracting: tkeel-template-go-main/pkg/service/README.md
  inflating: tkeel-template-go-main/pkg/service/greeter.go
  inflating: tkeel-template-go-main/pkg/service/openapi.go
   creating: tkeel-template-go-main/pkg/util/
  inflating: tkeel-template-go-main/pkg/util/util.go
   creating: tkeel-template-go-main/third_party/
 extracting: tkeel-template-go-main/third_party/README.md
   creating: tkeel-template-go-main/third_party/google/
   creating: tkeel-template-go-main/third_party/google/api/
  inflating: tkeel-template-go-main/third_party/google/api/annotations.proto
  inflating: tkeel-template-go-main/third_party/google/api/http.proto
  inflating: tkeel-template-go-main/third_party/google/api/httpbody.proto
   creating: tkeel-template-go-main/third_party/google/protobuf/
  inflating: tkeel-template-go-main/third_party/google/protobuf/descriptor.proto
  inflating: tkeel-template-go-main/third_party/google/protobuf/empty.proto
   creating: tkeel-template-go-main/third_party/protoc-gen-openapiv2/
   creating: tkeel-template-go-main/third_party/protoc-gen-openapiv2/options/
  inflating: tkeel-template-go-main/third_party/protoc-gen-openapiv2/options/annotations.proto
  inflating: tkeel-template-go-main/third_party/protoc-gen-openapiv2/options/openapiv2.proto
   creating: tkeel-template-go-main/third_party/validate/
  inflating: tkeel-template-go-main/third_party/validate/README.md
  inflating: tkeel-template-go-main/third_party/validate/validate.proto
✅  Success!! Plugin template created.
```

You can add the name of the directory you want to create after the create command, or you can install the template as a
git with the `-git` flag.
> Note: This usage requires the user to have the `git` command on their system

```bash
$ tkeel plugin create --git my_plugin
Cloning into 'my_plugin'...
remote: Enumerating objects: 95, done.
remote: Counting objects: 100% (95/95), done.
remote: Compressing objects: 100% (56/56), done.
remote: Total 95 (delta 22), reused 87 (delta 20), pack-reused 0
Receiving objects: 100% (95/95), 63.05 KiB | 15.76 MiB/s, done.
Resolving deltas: 100% (22/22), done.
✅  Success!! Plugin template created.
```