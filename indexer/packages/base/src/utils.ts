/**
 * Checks if any object has any defined properties.
 * @param object 
 * @returns True if there is at least 1 defined property, false otherwise.
 */
export function hasDefinedProperties(object: Object): boolean {
  for (const entry of Object.entries(object)) {
    if (entry[1] !== undefined) {
      return true;
    }
  }

  return false;
}
