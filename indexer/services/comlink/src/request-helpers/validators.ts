import _ from 'lodash';

/**
 * @function validateArray
 * @param inputArray array sanitized from query
 * @param setArray array of enum values to check against
 * @description Checks whether all the values  of the `inputArray` are present in the `setArray`
 */
export function validateArray(
  inputArray: string[],
  setArray: string[],
): boolean {
  return inputArray.length > 0 && _.difference(inputArray, setArray).length === 0;
}
