# E2E test

Tkeel CLI E2E test.

[中文文档](https://github.com/lunz1207/cli/blob/test/tests/e2e/README_zh.md)

## Test environment

Use kind to create a k8s cluster environment running in a docker container.

### Install

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

### Creat default cluster

```bash
kind creat cluster
```

### Build image

The node image used by kind to create the cluster by default contains only k8s-related components. The E2E test project need dapr, nodejs and other components, you need to package a custom node image.

Build custom image need [kubernetes](https://github.com/kubernetes/kubernetes) source and [kind](https://github.com/kubernetes-sigs/kind) source.

#### base image

A Docker image for running nested containers, systemd, and Kubernetes components. edite [Dockerfile](https://github.com/kubernetes-sigs/kind/blob/main/images/base/Dockerfile), add other components needed.run `make quick` on `/image/base` to build base image.

#### node image

Build need [building-kubernetes-with-docker](https://github.com/kubernetes/community/blob/master/contributors/devel/development.md#building-kubernetes-with-docker) and base image

> NOTE：
>
> 1. ensure kubernetes source on `$GOPATH/src`
> 2. ensure kubernetes source contain gitversion
> 3. use liunx to build

Build node image

```bash
kind build node-image --base-iamge diy_base_image:1.0 --image new_node_image:1.1
```

Create cluster by custom image and config file

```yaml
# kind.yaml config file
# a cluster with 1 control-plane node and 1 worker
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  # custom image
  image: mx2542/node:1.0@sha256:131c887156ffb257854c1d08ed664e635ca31b1b73146bec27024554fde72670
  extraMounts:
  # volumeMounts
   - hostPath: .
     containerPath: /cli
 - role: worker
```

```bash
kind creat cluster --config kind.yaml
```

## Test design

We designed the following test case classes, which mainly consist of three parts:

- Test case attributes
- Test execution method
- Execution result storage

The process of instantiating the test class is the test execution process。After the instantiation is completed, the execution result will be automatically bound to the instance for other use cases to call. If you need to process the test execution data, you can pass in the processing logic externally in the form of a function during instantiation. Use Case instances to organize test scenarios to optimize test data dependency between test cases.

```typescript
export class Case {
  public id!: string;
  public name!: string;
  public describe!: string;
  public command!: string;
  public expectation: any;
  public actuality: any;
  public store: any;
  public asyncStore: any;

  public filePath = "src/datas/";
  public txtMap = {
    tkeel: "tkeel 帮助文案",
    plugin_help: "plugin 帮助文案",
  };

  getter(fileName: string) {
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
    // run cmd
    const result = await util.promisify(exec)(args);
    return result;
  }

  cmder(args: string) {
    // run  cmd by async
    const result = execSync(args).toString();
    return result;
  }

  static async asyncInit(
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

### Writing case

Sync case

```typescript
import { Case } from "../core";

export const c1 = Case.init(
  "c001",
  "tkeel是否安装",
  "输入 tkeel 验证是否安装成功",
  "tkeel",
  "tkeel",
  (arg: string) => {
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
    const result = arg;
    return "this is a test date for c003";
  }
);
```

async case

```typescript
export async function cases() {
  const c1 = await Case.asyncInit(
    "c001",
    "tkeel是否安装",
    "输入 tkeel 验证是否安装成功",
    "tkeel",
    "tkeel",
    (arg: string) => {
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

### Check expectations

```typescript
import { c1, c2, cases } from "./cli.case";

describe("tkeel cli", () => {
  var testCases: any;

  beforeAll(async () => {
    testCases = await cases();
  });

  test("001", () => {
    expect(c1.actuality).toBe(c1.expectation);
    console.log(c2);
  });

  test("002", async () => {
    expect(testCases.c1.actuality).toBe(testCases.c1.expectation);
    console.log(testCases.c2);
  });
});
```

## Run test

### Local run

Required

- node.js
- npm
- tkeel CLI

```javascript
cd e2e
// install node package
npm install
// run test
npm run test
```

### Contains run

Run the script mounted in the container

```bash
docker exec kind_dokcer_container_name /bin/bash -c ". /cli/.github/scripts/run_e2e.sh"
```
