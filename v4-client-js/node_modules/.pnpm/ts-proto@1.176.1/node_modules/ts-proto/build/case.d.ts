import { Options } from "./options";
/** Converts `key` to TS/JS camel-case idiom, unless overridden not to. */
export declare function maybeSnakeToCamel(key: string, options: Pick<Options, "snakeToCamel">): string;
export declare function snakeToCamel(s: string): string;
export declare function camelToSnake(s: string): string;
export declare function capitalize(s: string): string;
export declare function uncapitalize(s: string): string;
export declare function camelCaseGrpc(s: string): string;
