import { CodeGeneratorRequest, DescriptorProto, EnumDescriptorProto, FieldDescriptorProto, FieldDescriptorProto_Type, FileDescriptorProto, MessageOptions, MethodDescriptorProto, ServiceDescriptorProto } from "ts-proto-descriptors";
import { Code, Import } from "ts-poet";
import { Options } from "./options";
import { Context } from "./context";
/** Based on https://github.com/dcodeIO/protobuf.js/blob/master/src/types.js#L37. */
export declare function basicWireType(type: FieldDescriptorProto_Type): number;
export declare function basicLongWireType(type: FieldDescriptorProto_Type): number | undefined;
/** Returns the type name without any repeated/required/etc. labels. */
export declare function basicTypeName(ctx: Context, field: FieldDescriptorProto, typeOptions?: {
    keepValueType?: boolean;
}): Code;
/** Returns the Reader method for the primitive's read/write call. */
export declare function toReaderCall(field: FieldDescriptorProto): string;
export declare function packedType(type: FieldDescriptorProto_Type): number | undefined;
export declare function getFieldOptionsJsType(field: FieldDescriptorProto, options: Options): FieldDescriptorProto_Type | undefined;
export declare function defaultValue(ctx: Context, field: FieldDescriptorProto): any;
/** Creates code that checks that the field is not the default value. Supports scalars and enums. */
export declare function notDefaultCheck(ctx: Context, field: FieldDescriptorProto, messageOptions: MessageOptions | undefined, place: string): Code;
/** A map of proto type name, e.g. `foo.Message.Inner`, to module/class name, e.g. `foo`, `Message_Inner`. */
export type TypeMap = Map<string, [string, string, DescriptorProto | EnumDescriptorProto]>;
/** Scans all of the proto files in `request` and builds a map of proto typeName -> TS module/name. */
export declare function createTypeMap(request: CodeGeneratorRequest, options: Options): TypeMap;
/** A "Scalar Value Type" as defined in https://developers.google.com/protocol-buffers/docs/proto3#scalar */
export declare function isScalar(field: FieldDescriptorProto): boolean;
export declare function isOptionalProperty(field: FieldDescriptorProto, messageOptions: MessageOptions | undefined, options: Options, isProto3Syntax: boolean): boolean;
/** This includes all scalars, enums and the [groups type](https://developers.google.com/protocol-buffers/docs/reference/java/com/google/protobuf/DescriptorProtos.FieldDescriptorProto.Type.html#TYPE_GROUP) */
export declare function isPrimitive(field: FieldDescriptorProto): boolean;
export declare function isBytes(field: FieldDescriptorProto): boolean;
export declare function isMessage(field: FieldDescriptorProto): boolean;
export declare function isEnum(field: FieldDescriptorProto): boolean;
export declare function isWithinOneOf(field: FieldDescriptorProto): boolean;
export declare function isWithinOneOfThatShouldBeUnion(options: Options, field: FieldDescriptorProto): boolean;
export declare function isRepeated(field: FieldDescriptorProto): boolean;
export declare function isLong(field: FieldDescriptorProto): boolean;
export declare function isWholeNumber(field: FieldDescriptorProto): boolean;
export declare function isMapType(ctx: Context, messageDesc: DescriptorProto, field: FieldDescriptorProto): boolean;
export declare function isObjectId(field: FieldDescriptorProto): boolean;
export declare function isTimestamp(field: FieldDescriptorProto): boolean;
export declare function isValueType(ctx: Context, field: FieldDescriptorProto): boolean;
export declare function isAnyValueType(field: FieldDescriptorProto): boolean;
export declare function isAnyValueTypeName(typeName: string): boolean;
export declare function isBytesValueType(field: FieldDescriptorProto): boolean;
export declare function isFieldMaskType(field: FieldDescriptorProto): boolean;
export declare function isFieldMaskTypeName(typeName: string): boolean;
export declare function isListValueType(field: FieldDescriptorProto): boolean;
export declare function isListValueTypeName(typeName: string): boolean;
export declare function isStructType(field: FieldDescriptorProto): boolean;
export declare function isStructTypeName(typeName: string): boolean;
export declare function isLongValueType(field: FieldDescriptorProto): boolean;
export declare function isEmptyType(typeName: string): boolean;
export declare function valueTypeName(ctx: Context, typeName: string): Code | undefined;
export declare function wrapperTypeName(typeName: string): string | undefined;
/** Maps `.some_proto_namespace.Message` to a TypeName. */
export declare function messageToTypeName(ctx: Context, protoType: string, typeOptions?: {
    keepValueType?: boolean;
    repeated?: boolean;
}): Code;
export declare function getEnumMethod(ctx: Context, enumProtoType: string, methodSuffix: string): Import;
/** Return the TypeName for any field (primitive/message/etc.) as exposed in the interface. */
export declare function toTypeName(ctx: Context, messageDesc: DescriptorProto | undefined, field: FieldDescriptorProto, ensureOptional?: boolean): Code;
/**
 * For a protobuf map field, if the generated code should use the javascript Map type.
 *
 * If the type of a protobuf map key corresponds to the Long type, we always use the Map type. This avoids generating
 * invalid code such as below (using Long as key of a javascript object):
 *
 * export interface Foo {
 *  bar: { [key: Long]: Long }
 * }
 *
 * See https://github.com/stephenh/ts-proto/issues/708 for more details.
 */
export declare function shouldGenerateJSMapType(ctx: Context, message: DescriptorProto, field: FieldDescriptorProto): boolean;
export declare function detectMapType(ctx: Context, messageDesc: DescriptorProto, fieldDesc: FieldDescriptorProto): {
    messageDesc: DescriptorProto;
    keyField: FieldDescriptorProto;
    keyType: Code;
    valueField: FieldDescriptorProto;
    valueType: Code;
} | undefined;
export declare function rawRequestType(ctx: Context, methodDesc: MethodDescriptorProto, typeOptions?: {
    keepValueType?: boolean;
    repeated?: boolean;
}): Code;
export declare function observableType(ctx: Context, asType?: boolean): Code;
export declare function requestType(ctx: Context, methodDesc: MethodDescriptorProto, partial?: boolean): Code;
export declare function responseType(ctx: Context, methodDesc: MethodDescriptorProto, typeOptions?: {
    keepValueType?: boolean;
    repeated?: boolean;
}): Code;
export declare function responsePromise(ctx: Context, methodDesc: MethodDescriptorProto): Code;
export declare function responseObservable(ctx: Context, methodDesc: MethodDescriptorProto): Code;
export declare function responsePromiseOrObservable(ctx: Context, methodDesc: MethodDescriptorProto): Code;
export interface BatchMethod {
    methodDesc: MethodDescriptorProto;
    uniqueIdentifier: string;
    singleMethodName: string;
    inputFieldName: string;
    inputType: Code;
    outputFieldName: string;
    outputType: Code;
    mapType: boolean;
}
export declare function detectBatchMethod(ctx: Context, fileDesc: FileDescriptorProto, serviceDesc: ServiceDescriptorProto, methodDesc: MethodDescriptorProto): BatchMethod | undefined;
export declare function isJsTypeFieldOption(options: Options, field: FieldDescriptorProto): boolean;
