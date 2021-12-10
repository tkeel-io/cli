# E2E test

Tkeel CLI E2E test.

[中文文档](https://github.com/lunz1207/cli/blob/test/tests/e2e/README_zh.md)

## Test design

We designed the following test case classes, which mainly consist of three parts:

- Test case attributes
- Test execution method
- Execution result storage

The process of instantiating the test class is the test execution process。After the instantiation is completed, the execution result will be automatically bound to the instance for other use cases to call. If you need to process the test execution data, you can pass in the processing logic externally in the form of a function during instantiation. Use Case instances to organize test scenarios to optimize test data dependency between test cases.

```typescript
class Case {
  public id!: string;
  public name!: string;
  public describe!: string;
  public command!: string;
  public expectation: any;
  public actuality: any;
  public store: any;
  public asyncStore: any;

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

### Writing case

Sync case

```typescript
import { Case } from "../core";

const c1 = Case.init(
  "001",
  "tkeel is install",
  "input tkeel and check log",
  { cmd: "tkeel" },
  "Things Keel Platform",
  (args: string) => {
    return "test date for c2";
  }
);

const c2 = Case.init(
  c1.store,
  "tkeel version",
  "input tkeel -v",
  { cmd: "tkeel", args: ["-v"] },
  "Keel CLI version: edge"
);
```

### Run test

Required

- node.js
- npm
- tkeel CLI

```javascript
cd test/e2e
// install node package
npm install
// run test
npm run test
```
