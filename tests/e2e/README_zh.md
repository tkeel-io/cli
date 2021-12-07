# E2E 测试

Tkeel CLI E2E 测试。

## 测试环境

使用 kind 创建运行在 docker 容器中的 k8s 集群环境 。

### 安装

linux

```bash
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.11.1/kind-linux-amd64
chmod +x ./kind
mv ./kind  /usr/local/bin/kind
```

macOS

```bash
brew Install kind
```

### 创建默认集群

```bash
kind creat cluster
```

### 镜像构建

kind 创建集群默认使用的 node image 仅包含 k8s 相关的组件。 本 E2E 测试项目依赖 dapr 、nodejs 等相关组件，因此需要打包自定义的 node image。

构造自定义镜像构建依赖 [kubernetes](https://github.com/kubernetes/kubernetes) 源码和 [kind](https://github.com/kubernetes-sigs/kind) 源码。

#### base image

一个小型 Docker 镜像，用于运行 containers、systemd 和 kubernetes 组件。修改 [Dockerfile](https://github.com/kubernetes-sigs/kind/blob/main/images/base/Dockerfile) ，加入需要的其他组件。在 `/image/base` 目录执行 make quick 构建 base image。

#### node image

构建依赖 [building-kubernetes-with-docker](https://github.com/kubernetes/community/blob/master/contributors/devel/development.md#building-kubernetes-with-docker) 和 base image

> NOTE：
>
> 1. 确保 kubernetes 源码在 `$GOPATH/src` 目录
> 2. 确保 kubernetes 源码包含 gitversion 信息
> 3. 建议在 liunx 平台构建

构建 node image

```bash
kind build node-image --base-iamge diy_base_image:1.0 --image new_node_image:1.1
```

使用自定义镜像和配置文件创建集群

```yaml
# kind.yaml 配置文件
# a cluster with 1 control-plane node and 1 worker
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  # 自定义的镜像
  image: mx2542/node:1.0@sha256:131c887156ffb257854c1d08ed664e635ca31b1b73146bec27024554fde72670
  extraMounts:
  # 挂载卷，挂载项目代码
   - hostPath: .
     containerPath: /cli
 - role: worker
```

```bash
kind creat cluster --config kind.yaml
```

## 测试设计

我们设计了如下测试用例类 ，主要包含三个部分：

- 测试用例属性
- 测试执行方法
- 执行结果存储

Case 类实例化的过程即是测试执行过程。实例化完成后会自动将执行的结果绑定实例，供其他用例调用。如果需要对测试执行数据进行处理，实例化时在外部将处理逻辑以函数的形式传入即可。以 Case 实例组织测试场景，优化测试用例之间的测试数据依赖问题。

```typescript
export class Case {
  public id!: string;
  // 用例 id
  public name!: string;
  // 用例名称
  public describe!: string;
  // 用例描述
  public command!: string;
  // 用例输入: tkeel 命令
  public expectation: any;
  // 预期结果
  public actuality: any;
  // 实际结果
  public store: any;
  // 用例输出,通常是供其他用例调用的测试数据
  public asyncStore: any;
  // 同上，异步方式

  public filePath = "src/datas/";
  public txtMap = {
    tkeel: "tkeel 帮助文案",
    plugin_help: "plugin 帮助文案",
  };

  getter(fileName: string) {
    // 当前预期结果以 txt 形式存储，因此需要一个读取的方法
    if (fileName in this.txtMap) {
      const content = fs.readFileSync(
        path.resolve(`${this.filePath}${fileName}.txt`),
        "utf-8"
      );
      return content;
    } else {
      return fileName;
    }
  }

  async asyncCmder(args: string) {
    // 异步执行 cmd
    const result = await util.promisify(exec)(args);
    return result;
  }

  cmder(args: string) {
    // 同步执行 cmd
    const result = execSync(args).toString();
    return result;
  }

  static async asyncInit(
    // 异步方法，初始化 CASE 实例
    id: string,
    name: string,
    describe: string,
    command: string,
    expectation: string,
    callback: Function
  ) {
    const c = new Case();
    c.id = id;
    c.name = name;
    c.describe = describe;
    c.expectation = c.getter(expectation);
    c.actuality = await c.asyncCmder(command);
    c.asyncStore = callback(c.actuality);
    return c;
  }

  static init(
    // 同步方法，初始化 CASE 实例
    id: string,
    name: string,
    describe: string,
    command: string,
    expectation: string,
    callback: Function
  ) {
    const c = new Case();
    c.id = id;
    c.name = name;
    c.describe = describe;
    c.expectation = c.getter(expectation);
    c.actuality = c.cmder(command);
    c.store = callback(c.actuality);
    return c;
  }
}
```

### 用例编写

同步用例

```typescript
import { Case } from "../core";

export const c1 = Case.init(
  "c001", // 用例id
  "tkeel是否安装", // 用例名称
  "输入 tkeel 验证是否安装成功", // 用例描述
  "tkeel", // 用例输入：待执行的命令
  "tkeel", // 预期结果
  (arg: string) => {
    // 将执行结果数据进行处理
    const result = arg;
    return "test date for  c002";
  }
);

export const c2 = Case.init(
  c1.store,
  "tkeel是否安装",
  "这是测试用例 c002",
  "tkeel",
  "tkeel",
  (arg: string) => {
    // 处理逻辑
    const result = arg;
    return "this is a test date for c003";
  }
);
```

异步用例

```typescript
export async function cases() {
  // 实例化 case 并打包
  const c1 = await Case.asyncInit(
    "c001",
    "tkeel是否安装",
    "输入 tkeel 验证是否安装成功",
    "tkeel",
    "tkeel",
    (arg: string) => {
      // 处理逻辑
      const result = arg;
      return "c002 from c001 store";
    }
  );

  const c2 = await Case.asyncInit(
    c1.asyncStore,
    "tkeel是否安装",
    "这是测试用例 c002",
    "tkeel",
    "tkeel",
    (arg: string) => {
      // 处理逻辑

      const result = arg;
      return "this is a test date for c003";
    }
  );

  const testCases = {
    c1,
    c2,
  };

  return testCases;
}
```

### 校验预期

```typescript
import { c1, c2, cases } from "./cli.case";

describe("tkeel cli", () => {
  var testCases: any;

  beforeAll(async () => {
    testCases = await cases();
  });

  test("001", () => {
    // 同步用例
    expect(c1.actuality).toBe(c1.expectation);
    console.log(c2);
  });

  test("002", async () => {
    // 异步用例
    expect(testCases.c1.actuality).toBe(testCases.c1.expectation);
    console.log(testCases.c2);
  });
});
```

### 本地执行用例

依赖

- node.js
- npm
- tkeel CLI

```javascript
cd e2e
// 安装node依赖
npm install
// 运行测试
npm run test
```

### 容器执行用例

运行挂载到容器内的脚本

```bash
docker exec kind_dokcer_container_name /bin/bash -c ". /cli/.github/scripts/run_e2e.sh"
```
