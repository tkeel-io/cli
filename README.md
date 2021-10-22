# Keel CLI

[中文](README_zh.md)

TKeel CLI is your main tool for various tasks related to TKeel.

You can use it to install the TKeel platform, manage plugins and users. 
TKeel CLI temporarily works in Kubernetes mode.

### Prerequisites

TKeel CLI can help you install the tKeel platform and help you manage the platform.

TKeel currently relies on Dapr (Kubernetes mode).

* Install [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* Install [Dapr on Kubernetes](https://docs.dapr.io/operations/hosting/kubernetes/kubernetes-deploy/)

### Install TKeel CLI

#### Using script to install the latest release

**Linux**

Install the latest linux TKeel CLI to `/usr/local/bin`

```bash
wget -q https://raw.githubusercontent.com/tkeel-io/cli/master/install/install.sh -O - | /bin/bash
```

**MacOS**

Install the latest darwin TKeel CLI to `/usr/local/bin`

```bash
curl -fsSL https://raw.githubusercontent.com/tkeel-io/cli/master/install/install.sh | /bin/bash
```

#### From the Binary Releases

Each release of TKeel CLI includes various OSes and architectures. These binary versions can be manually downloaded and installed.

1. Download the [TKeel CLI](https://github.com/tkeel-io/cli/releases)
2. Unpack it (e.g. tkeel_linux_amd64.tar.gz, tkeel_windows_amd64.zip)
3. Move it to your desired location.
   * For Linux/MacOS - `/usr/local/bin`
   * For Windows, create a directory and add this to your System PATH. For example create a directory called `c:\tkeel` and add this directory to your path, by editing your system environment variable.

### Init TKeel Platform on Kubernetes

([Prerequisite](#Prerequisites): Dapr and Kubectl is available in the environment)

Use the init command to initialize TKeel. 

```bash
tkeel init
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
✅  Success! TKeel Platform has been installed to namespace keel-system. To verify, run `tkeel plugin list -k' in your terminal. To get started, go here: https://tkeel.io/keel-getting-started
```

### Uninstall TKeel on Kubernetes

To remove TKeel from your Kubernetes cluster, use the `uninstall` command.

```
$ tkeel uninstall
```

### Deploy plugin

([Prerequisite](#Prerequisites): Dapr and Kubectl is available in the environment)

You can deploy the plugin app with the Dapr

[deploy-the-plugin-app](https://github.com/dapr/quickstarts/tree/v1.0.0/hello-kubernetes#step-3---deploy-the-nodejs-app-with-the-dapr-sidecar)


### Manage plugins

Use the plugin command to manage plugins.

1. List plugin

```bash
tkeel plugin list
```

Output should look like so:

```
plugin list              
NAME       NAMESPACE    HEALTHY  STATUS    PLUGINSTATUS  REPLICAS  VERSION  AGE  CREATED              
auth       keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00  
plugins    keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00  
keel       keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00
echo-demo  keel-system  False    Running   UNKNOWN       1         0.0.1    1m   2021-10-05 11:25.19  
```


2. Register plugin

```bash
tkeel plugin register echo-demo
```

Output should look like so:

```
✅  Success! Plugin<echo-demo> has been Registered to TKeel Platform . To verify, run `tkeel plugin list -k' in your terminal.
```

Check the status

```bash
tkeel plugin list
```

Output should look like so:

```
plugin list              
NAME       NAMESPACE    HEALTHY  STATUS    PLUGINSTATUS  REPLICAS  VERSION  AGE  CREATED              
auth       keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00  
plugins    keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00  
keel       keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00
echo-demo  keel-system  False    Running   ACTIVE        1         0.0.1    2m   2021-10-05 11:25.19  
```


3. Delete plugin

```bash
tkeel plugin delete echo-demo
```

Output should look like so:

```
✅  Success! Plugin<echo-demo> has been deleted from TKeel Platform . To verify, run `tkeel plugin list -k' in your terminal.
```
