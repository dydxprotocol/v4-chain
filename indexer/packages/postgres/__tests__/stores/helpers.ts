export function checkLengthAndContains(
  // eslint-disable-next-line  @typescript-eslint/no-explicit-any
  inputArray: any[],
  expectedLength: number,
  // eslint-disable-next-line  @typescript-eslint/no-explicit-any
  expectedValue?: any,
): void {
  expect(inputArray.length).toEqual(expectedLength);
  if (expectedValue !== undefined) {
    expect(inputArray).toEqual(
      expect.arrayContaining([
        expect.objectContaining(expectedValue),
      ]),
    );
  }
}
