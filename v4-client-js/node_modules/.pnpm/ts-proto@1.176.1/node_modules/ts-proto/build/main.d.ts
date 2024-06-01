import { Code } from "ts-poet";
import { ConditionalOutput } from "ts-poet/build/ConditionalOutput";
import { FileDescriptorProto } from "ts-proto-descriptors";
import { Context } from "./context";
import { Options } from "./options";
export declare function generateFile(ctx: Context, fileDesc: FileDescriptorProto): [string, Code];
export type Utils = ReturnType<typeof makeDeepPartial> & ReturnType<typeof makeObjectIdMethods> & ReturnType<typeof makeTimestampMethods> & ReturnType<typeof makeByteUtils> & ReturnType<typeof makeLongUtils> & ReturnType<typeof makeComparisonUtils> & ReturnType<typeof makeNiceGrpcServerStreamingMethodResult> & ReturnType<typeof makeGrpcWebErrorClass> & ReturnType<typeof makeExtensionClass> & ReturnType<typeof makeAssertionUtils>;
/** These are runtime utility methods used by the generated code. */
export declare function makeUtils(options: Options): Utils;
declare function makeLongUtils(options: Options, bytes: ReturnType<typeof makeByteUtils>): {
    numberToLong: ConditionalOutput;
    longToNumber: ConditionalOutput;
    longToString: ConditionalOutput;
    longToBigint: ConditionalOutput;
    Long: ConditionalOutput;
};
declare function makeByteUtils(options: Options): {
    globalThis: ConditionalOutput;
    bytesFromBase64: ConditionalOutput;
    base64FromBytes: ConditionalOutput;
};
declare function makeDeepPartial(options: Options, longs: ReturnType<typeof makeLongUtils>): {
    Builtin: ConditionalOutput;
    DeepPartial: ConditionalOutput;
    Exact: ConditionalOutput;
};
declare function makeObjectIdMethods(): {
    fromJsonObjectId: ConditionalOutput;
    fromProtoObjectId: ConditionalOutput;
    toProtoObjectId: ConditionalOutput;
};
declare function makeTimestampMethods(options: Options, longs: ReturnType<typeof makeLongUtils>, bytes: ReturnType<typeof makeByteUtils>): {
    toTimestamp: ConditionalOutput;
    fromTimestamp: ConditionalOutput;
    fromJsonTimestamp: ConditionalOutput;
};
declare function makeComparisonUtils(): {
    isObject: ConditionalOutput;
    isSet: ConditionalOutput;
};
declare function makeNiceGrpcServerStreamingMethodResult(options: Options): {
    NiceGrpcServerStreamingMethodResult: ConditionalOutput;
};
declare function makeGrpcWebErrorClass(bytes: ReturnType<typeof makeByteUtils>): {
    GrpcWebError: ConditionalOutput;
};
declare function makeExtensionClass(options: Options): {
    Extension: ConditionalOutput;
};
declare function makeAssertionUtils(bytes: ReturnType<typeof makeByteUtils>): {
    fail: ConditionalOutput;
};
export declare const contextTypeVar = "Context extends DataLoaders";
export {};
