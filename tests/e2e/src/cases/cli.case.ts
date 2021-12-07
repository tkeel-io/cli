import { Case } from "../core";

export const c1 = Case.init(
  "001",
  "tkeel 安装成功",
  "输入 tkeel 验证是否安装成功",
  "tkeel",
  "Things Keel Platform",
  (arg: string) => {
    // 处理逻辑
    const result = arg;
    return "test date for case 002";
  }
);

export const c2 = Case.init(
  c1.store,
  "tkeel 版本",
  "输入 tkeel -v",
  "tkeel -v",
  "Keel CLI version: edge",
  (arg: string) => {
    // 处理逻辑
    const result = arg;
    return "this is a test date for case 003";
  }
);
