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
✅  Success! tKeel Platform has been installed to namespace keel-system. To verify, run `tkeel plugin list' in your terminal. To get started, go here: https://tkeel.io/keel-getting-started
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
 NAME                                        PLUGIN                                      PLUGIN VERSION  REPO         REGISTER AT          STATE    DESCRIPTION
  console-plugin-admin-custom-config          console-plugin-admin-custom-config          2.0.2           tkeel        2022-08-17 20:35:37  RUNNING  前端管理平台的「平台配置」插件
  console-plugin-admin-notification-configs   console-plugin-admin-notification-configs   2.0.2           tkeel        2022-08-17 20:35:44  RUNNING  前端管理平台的「通知方式配置」插件
  console-plugin-admin-plugins                console-plugin-admin-plugins                2.0.2           tkeel        2022-08-17 20:35:49  RUNNING  前端管理平台的「插件管理」插件
  console-plugin-admin-service-monitoring     console-plugin-admin-service-monitoring     2.0.2           tkeel        2022-08-17 20:35:53  RUNNING  前端管理平台的「服务监控」插件。
  console-plugin-admin-tenants                console-plugin-admin-tenants                2.0.2           tkeel        2022-08-17 20:35:57  RUNNING  前端管理平台的「租户管理」插件
  console-plugin-admin-usage-statistics       console-plugin-admin-usage-statistics       2.0.2           tkeel        2022-08-17 20:36:00  RUNNING  前端管理平台的「用量统计」插件。
  console-plugin-tenant-alarm-policy          console-plugin-tenant-alarm-policy          2.0.2           tkeel        2022-08-17 20:36:04  RUNNING  前端租户平台的「告警策略」插件
  console-plugin-tenant-alarms                console-plugin-tenant-alarms                2.0.2           tkeel        2022-08-17 20:36:08  RUNNING  前端租户平台的「告警记录」插件
  console-plugin-tenant-data-query            console-plugin-tenant-data-query            2.0.2           tkeel        2022-08-17 20:36:11  RUNNING  前端租户平台的「数据查询」插件
  console-plugin-tenant-data-subscription     console-plugin-tenant-data-subscription     2.0.2           tkeel        2022-08-17 20:36:14  RUNNING  前端租户平台的「数据订阅」插件
  console-plugin-tenant-device-templates      console-plugin-tenant-device-templates      2.0.2           tkeel        2022-08-17 20:36:18  RUNNING  前端租户平台的「设备模板」插件
  console-plugin-tenant-devices               console-plugin-tenant-devices               2.0.2           tkeel        2022-08-17 20:36:21  RUNNING  前端租户平台的「设备列表」插件
  console-plugin-tenant-networks              console-plugin-tenant-networks              2.0.2           tkeel        2022-08-17 20:36:26  RUNNING  前端租户平台的「网络服务」插件
  console-plugin-tenant-notification-objects  console-plugin-tenant-notification-objects  2.0.2           tkeel        2022-08-17 20:36:30  RUNNING  前端租户平台的「通知对象」插件
  console-plugin-tenant-plugins               console-plugin-tenant-plugins               2.0.2           tkeel        2022-08-17 20:36:34  RUNNING  前端租户平台的「插件管理」插件
  console-plugin-tenant-roles                 console-plugin-tenant-roles                 2.0.2           tkeel        2022-08-17 20:36:37  RUNNING  前端租户平台的「角色管理」插件
  console-plugin-tenant-routing-rules         console-plugin-tenant-routing-rules         2.0.2           tkeel        2022-08-17 20:36:41  RUNNING  前端租户平台的「数据路由」插件
  console-plugin-tenant-users                 console-plugin-tenant-users                 2.0.2           tkeel        2022-08-17 20:36:44  RUNNING  前端租户平台的「用户管理」插件
  console-portal-admin                        console-portal-admin                        2.0.2           tkeel        2022-08-17 20:36:48  RUNNING  前端管理平台
  console-portal-tenant                       console-portal-tenant                       2.0.2           tkeel        2022-08-17 20:36:52  RUNNING  前端租户平台
  core-broker                                 core-broker                                 2.0.2           tkeel        2022-08-17 21:18:12  RUNNING  后端租户平台的「数据订阅」插件
  fluxswitch                                  fluxswitch                                  2.0.2           tkeel        2022-08-17 21:18:10  RUNNING  为物联网设备和物联网平台之间建立一个安全的双向TCP通道
  iothub                                      iothub                                      2.0.2           tkeel        2022-08-17 21:18:23  RUNNING  设备接入插件
  rule-manager                                rule-manager                                2.0.2           tkeel        2022-08-17 21:19:58  RUNNING  后端租户平台的「数据路由」插件
  tkeel-alarm                                 tkeel-alarm                                 2.0.2           tkeel        2022-08-17 21:19:29  RUNNING  监控告警插件
  tkeel-calc                                  tkeel-calc                                  0.0.1           helm-charts  2022-12-01 21:04:58  RUNNING  A Helm chart for Kubernetes
  tkeel-calc-console                          tkeel-calc-console                          0.0.1           helm-charts  2022-12-01 21:40:58  RUNNING  A Helm chart for Kubernetes
  tkeel-calc-mul                              tkeel-calc-mul                              0.0.1           helm-charts  2022-12-01 21:08:58  RUNNING  A Helm chart for Kubernetes
  tkeel-device                                tkeel-device                                2.0.2           tkeel        2022-08-17 21:19:16  RUNNING  设备管理插件
  tkeel-docs                                  tkeel-docs                                  2.0.2           tkeel        2022-08-17 21:23:46  RUNNING  帮助文档
  tkeel-monitor                               tkeel-monitor                               2.0.2           tkeel        2022-08-17 21:19:37  RUNNING  服务监控与用量统计
```

#### 注册插件

```bash
$ tkeel plugin register echo-demo
✅  Success! Plugin<echo-demo> has been Registered to tKeel Platform . To verify, run `tkeel plugin list' in your terminal.
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
$ tkeel plugin uninstall echo-demo
✅  Remove "echo-demo" success!
```

### 管理员登录

使用接下来的命令可以直接登录，获取管理员 token
> 采用不可见方式输入密码

```shell
tkeel admin login
```
