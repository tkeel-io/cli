/**
 * @todo: 重构 Cli 数据结构
 */
export type Command = {
  cmd: string;
  args?: Array<any>;
  options?: object;
};

/**
 * @todo: 定义 Http 数据结构，支持 http 借口测试 
 */
export type Http = {
  headers: object;
  body: object;
  cookie: object;
};

/**
 * @todo: 定义 testCase 数据结构
 */
export type TestCase = {
  testCaseId: string;
  testCaseDescription: string;
  testSteps: string;
  testData: Command | Http;
  expectedResults: string | object;
  actualResults;
  preRequisites?: { callback: Function };
  isAutomation: true | false
};
