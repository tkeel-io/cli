<h1 align="center"> tKeel CLI </h1>
<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/tkeel-io/cli)](https://goreportcard.com/report/github.com/tkeel-io/cli)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/tkeel-io/cli)
![GitHub](https://img.shields.io/github/license/tkeel-io/cli?style=plastic)
[![GoDoc](https://godoc.org/github.com/tkeel-io/cli?status.png)](http://godoc.org/github.com/tkeel-io/cli)
</div>

ð¹ï¸ tKeel CLI æ¯æ¨ç¨äºåç§ tKeel ç¸å³ä»»å¡æä½çç®æä½¿ç¨å·¥å·ã

æ¨å¯ä»¥ä½¿ç¨å®æ¥ **å®è£ tKeel å¹³å°**ã**ç®¡çæä»¶** ä»¥å **ç¨æ·æ¨¡å**ã

### å®è£é¡»ç¥

tKeel CLI å¯ä»¥å¸®å©æ¨å®è£ tKeel å¹³å°å¹¶ä¸å¸®å©æ¨ç®¡çå¹³å°ã

> â ï¸ tKeel ç°é¶æ®µä¾èµäº Daprï¼Kubernetes modeï¼ã

- å®è£ [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- å®è£ [Dapr on Kubernetes](https://docs.dapr.io/operations/hosting/kubernetes/kubernetes-deploy/)

### å®è£

ð§ æä»¬æä¾äºå¤ç§å®è£æ¹å¼ï¼æ¨æ ¹æ®æ¨çåå¥½éæ©æ¨è§å¾æåéçå®è£æ¹å¼è¿è¡å®è£ã

#### ä½¿ç¨èæ¬å®è£ææ°çæ¬

éè¿æä»¬ç¼åå¥½çèæ¬èªå¨ä¸ºæ¨å®è£ `tKeel Cli`ã

##### Linux

éè¿ Bash èæ¬å°ææ°ç tKeel Cli å®è£è³ Linux ç³»ç»ç `/usr/local/bin`

```bash
$ wget -q https://raw.githubusercontent.com/tkeel-io/cli/master/install/install.sh -O - | /bin/bash
```

##### MacOS

éè¿ Bash èæ¬å°ææ°ç tKeel Cli å®è£è³ MacOS(darwin) ç³»ç»ç `/usr/local/bin`

```bash
$ curl -fsSL https://raw.githubusercontent.com/tkeel-io/cli/master/install/install.sh | /bin/bash
```

#### éè¿åè¡çäºè¿å¶ç¨åº

æ¯ä¸ªåè¡çæ¬ç tKeel CLI åæ¬åç§æä½ç³»ç»åæ¶æãè¿äºäºè¿å¶çæ¬å¯ä»¥æå¨ä¸è½½åå®è£ã

1. ä¸è½½ [tKeel CLI](https://github.com/tkeel-io/cli/releases)
2. å°ä¸è½½çæä»¶è§£å (e.g. tkeel_linux_amd64.tar.gz, tkeel_windows_amd64.zip)
3. æå®ç§»å°ä½ æ³è¦çä½ç½®
    * å¦æä½ æ¯ Linux/MacOS ç¨æ· - `/usr/local/bin`
    * å¦æä½ æ¯ Windows ç¨æ· - åå»ºä¸ä¸ªç®å½å¹¶å°å¶æ·»å å°ä½ ç `ç³»ç» PATH `ä¸­ãä¾å¦ï¼éè¿ç¼è¾ç³»ç»ç¯å¢åéï¼åå»ºä¸ä¸ªåä¸º`c:\tkeel`çç®å½ï¼å¹¶å°è¿ä¸ªç®å½æ·»å å°ä½ ç `ç³»ç» PATH` ä¸­ã

### å¨ Kubernetes åå§ tKeel å¹³å°

> è¯·æ³¨æ [å®è£é¡»ç¥](#å®è£é¡»ç¥) ç¡®ä¿ä½ çç³»ç»ä¸­æææç¯å¢ã

ä½¿ç¨å½ä»¤è¡åå§ `tKeel`

```bash
$ tkeel init
```

> æ³¨æï¼Linux ç¨æ·è¯·æ³¨æï¼å¦æä½ ç docker éè¦ä½¿ç¨ sudo æéæè½ä½¿ç¨ï¼é£ä¹è¯·ä½ ä½¿ç¨ `sudo tkeel init`

Output should look like so:

```
â  Making the jump to hyperspace...
â¹ï¸  Checking the Dapr runtime status...
â  Deploying the tKeel Platform to your cluster... 
â¹ï¸  install plugins...                                                        
â¹ï¸  install plugins done.                                                                                                        
â  Deploying the tKeel Platform to your cluster...
â  Register the plugins ... 
â¹ï¸  Plugin<plugins>  is registered.                                                                                          
â¹ï¸  Plugin<keel>  is registered.                                                                                                                        
â¹ï¸  Plugin<auth>  is registered.                                                                                                                        
â  Success! tKeel Platform has been installed to namespace keel-system. To verify, run `tkeel plugin list' in your terminal. To get started, go here: https://tkeel.io/keel-getting-started
```

### å¸è½½ tKeel å¹³å°

è¦ä»ä½ ç Kubernetes éç¾¤ä¸­ç§»é¤ `tKeel`ï¼è¯·ä½¿ç¨ `uninstall`å½ä»¤ã

```bash
$ tkeel uninstall
```

### é¨ç½²æä»¶

ä½ è½éè¿ Dapr é¨ç½² tKeel çæä»¶ï¼
è¯¦ç»è¯·è§ [deploy-the-plugin-app ææ¡£](https://github.com/dapr/quickstarts/tree/v1.0.0/hello-kubernetes#step-3---deploy-the-nodejs-app-with-the-dapr-sidecar)

### ç®¡çæä»¶

ä½¿ç¨æä»¶å½ä»¤å»ç®¡çå¹³å°ä¸çæä»¶ã

#### å±ç¤ºæææä»¶

```bash
$ tkeel plugin list
```

æ¨ä¼å¾å°åæ¯è¿æ ·çä¸ä¸²è¾åº:

```bash
$ plugin list              
NAME       NAMESPACE    HEALTHY  STATUS    PLUGINSTATUS  REPLICAS  VERSION  AGE  CREATED              
auth       keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00  
plugins    keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00  
keel       keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00
echo-demo  keel-system  False    Running   UNKNOWN       1         0.0.1    1m   2021-10-05 11:25.19  
```

#### æ³¨åæä»¶

```bash
$ tkeel plugin register echo-demo
â  Success! Plugin<echo-demo> has been Registered to tKeel Platform . To verify, run `tkeel plugin list' in your terminal.
```

ä½¿ç¨` plugin list ` å¯ä»¥æ¥çæä»¶ç¶æ

```bash
$ tkeel plugin list              
NAME       NAMESPACE    HEALTHY  STATUS    PLUGINSTATUS  REPLICAS  VERSION  AGE  CREATED              
auth       keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00  
plugins    keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00  
keel       keel-system  True     Running   ACTIVE        1         0.0.1    37m  2021-10-07 16:07.00
echo-demo  keel-system  False    Running   ACTIVE        1         0.0.1    2m   2021-10-05 11:25.19  
```

#### å é¤æä»¶

```bash
$ tkeel plugin uninstall echo-demo
â  Remove "echo-demo" success!
```

### ç®¡çåç»å½

ä½¿ç¨æ¥ä¸æ¥çå½ä»¤å¯ä»¥ç´æ¥ç»å½ï¼è·åç®¡çå token
> éç¨ä¸å¯è§æ¹å¼è¾å¥å¯ç 

```shell
tkeel admin login
```
