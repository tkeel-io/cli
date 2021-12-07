import path from "path";
var fs = require("fs");
var util = require("util");
var { execSync, exec } = require("child_process");

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
