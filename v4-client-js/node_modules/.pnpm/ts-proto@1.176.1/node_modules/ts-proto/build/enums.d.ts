import { Code } from "ts-poet";
import { EnumDescriptorProto, EnumValueDescriptorProto } from "ts-proto-descriptors";
import SourceInfo from "./sourceInfo";
import { Context } from "./context";
type UnrecognizedEnum = {
    present: false;
} | {
    present: true;
    name: string;
};
export declare function generateEnum(ctx: Context, fullName: string, enumDesc: EnumDescriptorProto, sourceInfo: SourceInfo): Code;
/** Generates a function with a big switch statement to decode JSON -> our enum. */
export declare function generateEnumFromJson(ctx: Context, fullName: string, enumDesc: EnumDescriptorProto, unrecognizedEnum: UnrecognizedEnum): Code;
/** Generates a function with a big switch statement to encode our enum -> JSON. */
export declare function generateEnumToJson(ctx: Context, fullName: string, enumDesc: EnumDescriptorProto, unrecognizedEnum: UnrecognizedEnum): Code;
/** Generates a function with a big switch statement to encode our string enum -> int value. */
export declare function generateEnumToNumber(ctx: Context, fullName: string, enumDesc: EnumDescriptorProto, unrecognizedEnum: UnrecognizedEnum): Code;
export declare function getMemberName(ctx: Context, enumDesc: EnumDescriptorProto, valueDesc: EnumValueDescriptorProto): string;
export {};
