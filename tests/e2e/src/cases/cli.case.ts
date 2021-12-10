import { Case } from "../core";

export const c1 = Case.init(
  "001",
  "tkeel 是否安装成功",
  "输入 tkeel 验证是否安装成功",
  { cmd: "tkeel" },
  "Things Keel Platform"
);

export const c2 = Case.init(
  "002",
  "tkeel 版本",
  "输入 tkeel -v",
  { cmd: "tkeel", args: ["-v"] },
  "Keel CLI version: edge"
);

export const c3 = Case.init(
  "003",
  "tkeel 自动补全",
  "输入 tkeel completion",
  { cmd: "tkeel", args: ["completion"] },
  "Generates shell completion scripts"
);

export const c4 = Case.init(
  "004",
  "tkeel 帮助",
  "输入 tkeel help",
  { cmd: "tkeel", args: ["help"] },
  `Use "tkeel [command] --help" for more information about a command`
);
