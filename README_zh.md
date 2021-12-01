<h1 align="center"> tKeel CLI </h1>
<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/tkeel-io/cli)](https://goreportcard.com/report/github.com/tkeel-io/cli)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/tkeel-io/cli)
![GitHub](https://img.shields.io/github/license/tkeel-io/cli?style=plastic)
[![GoDoc](https://godoc.org/github.com/tkeel-io/cli?status.png)](http://godoc.org/github.com/tkeel-io/cli)
</div>

🕹️ tKeel CLI 是您用于各种 tKeel 相关任务操作的简易使用工具。

您可以使用它来 **安装 tKeel 平台**、**管理插件** 以及 **用户模块**。

### 安装须知

tKeel CLI 可以帮助您安装 tKeel 平台并且帮助您管理平台。

> ⚠️ tKeel 现阶段依赖于 Dapr（Kubernetes mode）。

- 安装 [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- 安装 [Dapr on Kubernetes](https://docs.dapr.io/operations/hosting/kubernetes/kubernetes-deploy/)

### 安装

🔧 我们提供了多种安装方式，您根据您的偏好选择您觉得最合适的安装方式进行安装。

#### 使用脚本安装最新版本

通过我们编写好的脚本自动为您安装 `tKeel Cli`。

##### Linux

通过 Bash 脚本将最新版 tKeel Cli 安装至 Linux 系统的 `/usr/local/bin`

```bash
$ wget -q https://raw.githubusercontent.com/tkeel-io/cli/master/install/install.sh -O - | /bin/bash
```

##### MacOS

通过 Bash 脚本将最新版 tKeel Cli 安装至 MacOS(darwin) 系统的 `/usr/local/bin`

```bash
$ curl -fsSL https://raw.githubusercontent.com/tkeel-io/cli/master/install/install.sh | /bin/bash
```

#### 通过发行的二进制程序

每个发行版本的 tKeel CLI 包括各种操作系统和架构。这些二进制版本可以手动下载和安装。

1. 下载 [tKeel CLI](https://github.com/tkeel-io/cli/releases)
2. 将下载的文件解压 (e.g. tkeel_linux_amd64.tar.gz, tkeel_windows_amd64.zip)
3. 把它移到你想要的位置
    * 如果你是 Linux/MacOS 用户 - `/usr/local/bin`
    * 如果你是 Windows 用户 - 创建一个目录并将其添加到你的 `系统 PATH `中。例如，通过编辑系统环境变量，创建一个名为`c:\tkeel`的目录，并将这个目录添加到你的 `系统 PATH` 中。

### 在 Kubernetes 初始 tKeel 平台

> 请注意 [安装须知](#安装须知) 确保你的系统中有所有环境。

使用命令行初始 `tKeel`

```bash
$ tkeel init
```

> 注意：Linux 用户请注意，如果你的 docker 需要使用 sudo 权限才能使用，那么请你使用 `sudo tkeel init`

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

### 卸载 tKeel 平台

要从你的 Kubernetes 集群中移除 `tKeel`，请使用 `uninstall`命令。

```bash
$ tkeel uninstall
```

### 部署插件

你能通过 Dapr 部署 tKeel 的插件，
详细请见 [deploy-the-plugin-app 文档](https://github.com/dapr/quickstarts/tree/v1.0.0/hello-kubernetes#step-3---deploy-the-nodejs-app-with-the-dapr-sidecar)

### 管理插件

使用插件命令去管理平台上的插件。

#### 展示所有插件

```bash
$ tkeel plugin list
```

您会得到像是这样的一串输出:

```bash
$ plugin list              
NAME       NAMESPACE    HEALTHY  STATUS    PLUGINSTATUS  REPLICAS  VERSION  AGE  CREATED              
auth       keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00  
plugins    keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00  
keel       keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00
echo-demo  keel-system  False    Running   UNKNOWN       1         0.0.1    1m   2021-10-05 11:25.19  
```

#### 注册插件

```bash
$ tkeel plugin register echo-demo
✅  Success! Plugin<echo-demo> has been Registered to tKeel Platform . To verify, run `tkeel plugin list -k' in your terminal.
```

使用` plugin list ` 可以查看插件状态

```bash
$ tkeel plugin list              
NAME       NAMESPACE    HEALTHY  STATUS    PLUGINSTATUS  REPLICAS  VERSION  AGE  CREATED              
auth       keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00  
plugins    keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00  
keel       keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00
echo-demo  keel-system  False    Running   ACTIVE        1         0.0.1    2m   2021-10-05 11:25.19  
```

#### 删除插件

```bash
$ tkeel plugin delete echo-demo
✅  Success! Plugin<echo-demo> has been deleted from tKeel Platform . To verify, run `tkeel plugin list -k' in your terminal.
```
