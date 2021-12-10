import { c1, c2, c3, c4 } from "./cli.case";

describe("tkeel cli", () => {
  test(c1.name, async () => {
    expect(c1.actuality).toContain(c1.expectation);
  });

  test(c2.name, async () => {
    expect(c2.actuality).toContain(c2.expectation);
  });

  test(c3.name, async () => {
    expect(c3.actuality).toContain(c3.expectation);
  });

  test(c4.name, async () => {
    expect(c4.actuality).toContain(c4.expectation);
  });
});
