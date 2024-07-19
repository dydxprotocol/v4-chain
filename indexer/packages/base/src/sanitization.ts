import traverse from 'traverse';

// Common request headers which should be redacted. Normalized to all-lowercase.
export const DEFAULT_SECRET_KEYS = [
  'Authorization',
  'X-Routing-Key', // Used by PagerDuty.
];

const DEFAULT_REDACTED_PLACEHOLDER = '[REDACTED]';
const JSON_CIRCULAR_PLACEHOLDER = '[CIRCULAR]';

/**
 * Creates a deep copy of an object with values redacted where the key matches `secretKeys`.
 */
export function redact<T>(
  obj: T,
  secretKeys: string[] = DEFAULT_SECRET_KEYS,
  placeholder: string = DEFAULT_REDACTED_PLACEHOLDER,
): T {
  const normalizedSecretKeys = secretKeys.map((s) => s.toLowerCase());

  // eslint-disable-next-line array-callback-return
  return traverse(obj).map(function traverseFunction(this: traverse.TraverseContext, value: {}) {
    if (normalizedSecretKeys.includes(this.key?.toLowerCase() as string) && value !== null) {
      this.update(placeholder);
    }
  });
}

/**
 * Creates a deep copy of an object with circular references removed or replaced.
 */
export function removeCircularReferences<T>(
  obj: T,
  placeholder: string | null = null,
): T {
  // eslint-disable-next-line array-callback-return
  return traverse(obj).map(function traverseFunction(this: traverse.TraverseContext, _value: {}) {
    if (this.circular) {
      if (placeholder !== null) {
        this.update(placeholder);
      } else {
        this.remove();
      }
    }
  });
}

/**
 * A modified JSON.stringify() that can be used with unknown input when it's important not to throw.
 *
 * Examples known to break JSON.stringify():
 *   - circular reference
 *   - BigInt
 *   - had an issue with a certain kind of buffer object in an Axios error when using HTTPS
 */
export function safeJsonStringify(value: unknown, spaces?: number): string {
  try {
    const noCircular = removeCircularReferences(value, JSON_CIRCULAR_PLACEHOLDER);
    return JSON.stringify(noCircular, (_key, val) => {
      // If val is an object, assume it's safe to recurse (also handle case where val === null).
      if (typeof val === 'object') {
        return val;
      }

      // Handle BigInt.
      if (typeof val === 'bigint') {
        return val.toString();
      }

      // TODO: Figure out if the buffer issue mentioned above can be handled here.
      return val;
    }, spaces);
  } catch (error) {
    return `${value}`;
  }
}
