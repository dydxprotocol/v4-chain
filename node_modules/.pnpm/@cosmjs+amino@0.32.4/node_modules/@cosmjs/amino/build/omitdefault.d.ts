/**
 * Returns the given input. If the input is the default value
 * of protobuf, undefined is retunred. Use this when creating Amino JSON converters.
 */
export declare function omitDefault<T extends string | number | bigint | boolean>(input: T): T | undefined;
