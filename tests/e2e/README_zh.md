# E2E 测试

Tkeel CLI E2E 测试。

## 测试设计

我们设计了如下测试用例类 ，主要包含三个部分：

- 测试用例属性
- 测试执行方法
- 执行结果存储

Case 实例化的过程即是测试执行过程。实例化完成后会自动将执行的结果绑定实例，供其他用例调用。如果需要对测试结果进行处理，实例化时在外部将处理逻辑以函数的形式传入。

```typescript
class Case {
  public id!: string;
  // 用例 id
  public name!: string;
  // 用例名称
  public describe!: string;
  // 用例描述
  public command!: Command;
  // 用例输入: tkeel 命令
  public expectation: any;
  // 预期结果
  public actuality: any;
  // 实际结果
  public store: any;
  // 用例输出,通常是供其他用例调用的测试数据
  public asyncStore: any;
  // 同上，异步方式

  static async asyncInit(
    id: string,
    name: string,
    describe: string,
    command: Command,
    expectation: string,
    callback?: Function
  ) {
    const c = new Case();
    c.id = id;
    c.name = name;
    c.describe = describe;
    c.expectation = expectation;
    c.actuality = await asyncSpawner(command);
    if (callback) {
      c.asyncStore = callback(c.actuality);
    }
    return c;
  }

  static init(
    id: string,
    name: string,
    describe: string,
    command: Command,
    expectation: string,
    callback?: Function
  ) {
    const c = new Case();
    c.id = id;
    c.name = name;
    c.describe = describe;
    c.expectation = expectation;
    c.actuality = spawner(command);
    if (callback) {
      c.asyncStore = callback(c.actuality);
    }
    return c;
  }
}
```

### 用例编写

```typescript
import { Case } from "../core";

const c1 = Case.init(
  "001",
  "tkeel 是否安装成功",
  "输入 tkeel 验证是否安装成功",
  { cmd: "tkeel" },
  "Things Keel Platform",
  (args:string)=>{
    return 'test date for c2'
  }
);

const c2 = Case.init(
  c1.store,
  "tkeel 版本",
  "输入 tkeel -v",
  { cmd: "tkeel", args: ["-v"] },
  "Keel CLI version: edge"
);
```

### 编写测试

```typescript
import { c1 } from "./cli.case";

describe("tkeel cli", () => {
  test(c1.name, async () => {
    expect(c1.actuality).toContain(c1.expectation);
  });
});
```

### 执行测试

依赖

- node.js
- npm
- tkeel CLI

```javascript
cd tests/e2e
// 安装node依赖
npm install
// 运行测试
npm run test
```
