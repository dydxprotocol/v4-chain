import { Context } from "./context";
import { Code } from "ts-poet";
export type StructFieldNames = {
    nullValue: string;
    numberValue: string;
    stringValue: string;
    boolValue: string;
    structValue: string;
    listValue: string;
};
/** Whether we need to generate `.wrap` and `.unwrap` methods for the given type. */
export declare function isWrapperType(fullProtoTypeName: string): boolean;
/**
 * Converts ts-proto's idiomatic Struct/Value/ListValue representation to the proto messages.
 *
 * We do this deeply b/c NestJS does not invoke wrappers recursively.
 */
export declare function generateWrapDeep(ctx: Context, fullProtoTypeName: string, fieldNames: StructFieldNames): Code[];
/**
 * Converts proto's Struct/Value?listValue messages to ts-proto's idiomatic representation.
 *
 * We do this deeply b/c NestJS does not invoke wrappers recursively.
 */
export declare function generateUnwrapDeep(ctx: Context, fullProtoTypeName: string, fieldNames: StructFieldNames): Code[];
/**
 * Converts ts-proto's idiomatic Struct/Value/ListValue representation to the proto messages.
 *
 * We do this shallow's b/c ts-proto's encode methods handle the recursion.
 */
export declare function generateWrapShallow(ctx: Context, fullProtoTypeName: string, fieldNames: StructFieldNames): Code[];
/**
 * Converts proto's Struct/Value?listValue messages to ts-proto's idiomatic representation.
 *
 * We do this shallowly b/c ts-proto's decode methods handle recursion.
 */
export declare function generateUnwrapShallow(ctx: Context, fullProtoTypeName: string, fieldNames: StructFieldNames): Code[];
