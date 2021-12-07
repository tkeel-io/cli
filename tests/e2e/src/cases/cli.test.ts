import { c1, c2 } from "./cli.case";

describe("tkeel cli", () => {
  test("001", async () => {
    expect(c1.actuality).toContain(c1.expectation);
  });

  test("002", async () => {
    expect(c2.actuality).toContain(c2.expectation);
  });
});
