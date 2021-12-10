var { execSync, exec, spawnSync, spawn } = require("child_process");
var util = require("util");
import { Command } from "./core.d";

/**
 * 异步执行命令，不支持交互式输入
 * @param {string}  cmd - 命令
 * @return {object} 命令的执行结果
 */
export async function asyncCmder(cmd: string): Promise<any> {
  const result = await util.promisify(exec)(cmd);
  return result;
}

/**
 * 同步执行命令，不支持交互式输入
 * @param {string}  cmd - 命令
 * @return {string}} 命令的执行结果
 */
export function cmder(cmd: string) {
  const result = execSync(cmd).toString();
  return result;
}

/**
 * 异步执行命令，支持交互式输入
 * @param {string}  cmd - 命令
 * @param {Array} arg - 命令的参数
 * @param {string} [content] - 可选参数，用户输入
 * @return {string} 命令的执行结果
 */
export async function asyncSpawner(cmd: Command) {
  const sp = spawn(cmd.cmd, cmd.args);

  if (cmd.options) {
    sp.stdin.write(cmd.options);
    sp.stdin.end();
  }

  let result;
  for await (const chunk of sp.stdout) {
    result = chunk.toString();
  }

  return result;
}

/**
 * 同步执行命令，支持交互式输入
 * @param {string}  cmd - 命令
 * @param {Array} arg - 命令的参数
 * @param {object} [options] - 可选参数，object 类型
 * @return {string} 命令的执行结果
 */
export function spawner(cmd: Command) {
  const result = spawnSync(cmd.cmd, cmd.args, cmd.options).stdout.toString();
  return result;
}
