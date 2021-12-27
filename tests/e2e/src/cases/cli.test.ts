import { cases } from "./cli.case";

describe("tkeel cli", () => {
  test(cases.c1.name, async () => {
    expect(cases.c1.actuality).toContain(cases.c1.expectation);
  });

  test(cases.c2.name, async () => {
    expect(cases.c2.actuality).toContain(cases.c2.expectation);
  });

  test(cases.c3.name, async () => {
    expect(cases.c3.actuality).toContain(cases.c3.expectation);
  });

  test(cases.c3_1.name, async () => {
    expect(cases.c3_1.actuality).toContain(cases.c3_1.expectation);
  });

  test(cases.c3_2.name, async () => {
    expect(cases.c3_2.actuality).toContain(cases.c3_2.expectation);
  });

  test(cases.c3_3.name, async () => {
    expect(cases.c3_3.actuality).toContain(cases.c3_3.expectation);
  });

  test(cases.c4.name, async () => {
    expect(cases.c4.actuality).toContain(cases.c4.expectation);
  });

  test(cases.c5.name, async () => {
    expect(cases.c5.actuality).toContain(cases.c5.expectation);
  });

  test(cases.c6.name, async () => {
    expect(cases.c6.actuality).toContain(cases.c6.expectation);
  });

  test(cases.c6_1.name, async () => {
    expect(cases.c6_1.actuality).toContain(cases.c6_1.expectation);
  });

  test(cases.c6_2.name, async () => {
    expect(cases.c6_2.actuality).toContain(cases.c6_2.expectation);
  });

  test(cases.c6_3.name, async () => {
    expect(cases.c6_3.asyncStore.Name).toBe(cases.c6_3.expectation);
  });

  test(cases.c6_4.name, async () => {
    expect(cases.c6_4.actuality).toContain(cases.c6_4.expectation);
  });

  test(cases.c6_5.name, async () => {
    expect(cases.c6_5.actuality).toContain(cases.c6_5.expectation);
  });

  test(cases.c7.name, async () => {
    expect(cases.c7.actuality).toContain(cases.c7.expectation);
  });

  test(cases.c7_1.name, async () => {
    expect(cases.c7_1.actuality).toContain(cases.c7_1.expectation);
  });
});
