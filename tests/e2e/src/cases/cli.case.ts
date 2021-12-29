import { Case } from "../core";

export const cases = {
  c1: Case.init(
    "001",
    "tkeel 安装",
    "输入 tkeel 验证是否安装成功",
    { cmd: "tk" },
    "Things Keel Platform"
  ),
  c2: Case.init(
    "002",
    "tkeel 版本",
    "输入 tkeel -v",
    { cmd: "tk", args: ["-v"] },
    "Keel CLI version: edge"
  ),
  c3: Case.init(
    "003",
    "tkeel 自动补全",
    "输入 tkeel completion",
    { cmd: "tk", args: ["completion"] },
    "Generates shell completion scripts"
  ),
  c3_1: Case.init(
    "003_1",
    "tkeel bash 自动补全",
    "输入 tkeel completion bash",
    { cmd: "tk", args: ["completion", "bash"] },
    "bash completion for tkeel"
  ),
  c3_2: Case.init(
    "003_2",
    "tkeel powershell 自动补全",
    "输入 tkeel completion powershell",
    { cmd: "tk", args: ["completion", "powershell"] },
    "powershell completion for tkeel"
  ),
  c3_3: Case.init(
    "003_3",
    "tkeel zsh 自动补全",
    "输入 tkeel completion zsh",
    { cmd: "tk", args: ["completion", "zsh"] },
    "zsh completion for tkeel"
  ),
  c4: Case.init(
    "004",
    "tkeel 帮助",
    "输入 tkeel help",
    { cmd: "tk", args: ["help"] },
    `Use "tkeel [command] --help" for more information about a command`
  ),
  c5: Case.init(
    "005",
    "tkeel 初始化",
    "输入 tkeel init -h",
    { cmd: "tk", args: ["init", "-h"] },
    "Initialize Keel in Kubernete"
  ),
  c6: Case.init(
    "006",
    "tkeel 插件",
    "输入 tkeel plugin",
    { cmd: "tk", args: ["plugin"] },
    "Get status of tKeel plugins from Kubernetes"
  ),
  c6_1: Case.init(
    "006_1",
    "tkeel 无插件",
    "输入 tkeel plugin list",
    { cmd: "tk", args: ["plugin", "list"] },
    "No status returned. Is tKeel plugins not install in your cluster?"
  ),
  c6_2: Case.init(
    "006_2",
    "tkeel 安装插件",
    "安装 keel-echo 插件",
    {
      cmd: "tk",
      args: [
        "plugin",
        "install",
        "https://tkeel-io.github.io/helm-charts/keel-echo@v0.2.0",
        "tkeel-echo",
        "-n",
        "testing",
      ],
    },
    `Install "keel-echo" success! It's named "tkeel-echo" in k8s`
  ),
  c6_3: Case.init(
    "006_3",
    "tkeel 查看插件列表",
    "安装 keel-echo 插件",
    {
      cmd: "tk",
      args: ["plugin", "list", "-o", "json"],
    },
    "keel-echo",
    (args: string) => {
      const result = JSON.parse(args);
      return result;
    }
  ),
  c6_4: Case.init(
    "006_4",
    "tkeel 注册插件",
    "注册 keel-echo 插件",
    {
      cmd: "tk",
      args: ["plugin", "register", "keel-echo"],
    },
    "Success! Plugin<keel-echo> has been Registered to tKeel Platform"
  ),
  c6_5: Case.init(
    "006_5",
    "tkeel 反注册插件",
    "反注册 keel-echo 插件",
    {
      cmd: "tk",
      args: ["plugin", "unregister", "keel-echo"],
    },
    `Unregister plugin: {"id":"keel-echo"`
  ),
  c7: Case.init(
    "007",
    "tkeel 创建租户",
    "创建租户abcdefg",
    { cmd: "tk", args: ["tenant", "create", "abcdefg", "-k"] },
    "Success!"
  ),
  c7_1: Case.init(
    "007",
    "tkeel 查看租户",
    "查看租户列表",
    { cmd: "tk", args: ["tenant", "list", "-k"] },
    "abcdefg"
  ),
};
