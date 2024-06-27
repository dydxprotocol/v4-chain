/**
 * Removes all '/', '\', and spaces from the beginning and end of the path
 * Replaces all '/', '\', and spaces between sections of the path
 * Adds prefix and suffix if supplied
 * Replace ':pathParam' with '{pathParam}'
 */
export declare function normalisePath(path: string, withPrefix?: string, withSuffix?: string, skipPrefixAndSuffixIfEmpty?: boolean): string;
export declare function convertColonPathParams(path: string): string;
export declare function convertBracesPathParams(path: string): string;
