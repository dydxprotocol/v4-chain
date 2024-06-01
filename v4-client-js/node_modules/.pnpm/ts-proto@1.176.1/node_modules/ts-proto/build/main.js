"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.contextTypeVar = exports.makeUtils = exports.generateFile = void 0;
const ts_poet_1 = require("ts-poet");
const ConditionalOutput_1 = require("ts-poet/build/ConditionalOutput");
const ts_proto_descriptors_1 = require("ts-proto-descriptors");
const case_1 = require("./case");
const enums_1 = require("./enums");
const generate_async_iterable_1 = require("./generate-async-iterable");
const generate_generic_service_definition_1 = require("./generate-generic-service-definition");
const generate_grpc_js_1 = require("./generate-grpc-js");
const generate_grpc_web_1 = require("./generate-grpc-web");
const generate_nestjs_1 = require("./generate-nestjs");
const generate_nice_grpc_1 = require("./generate-nice-grpc");
const generate_services_1 = require("./generate-services");
const generate_struct_wrappers_1 = require("./generate-struct-wrappers");
const options_1 = require("./options");
const schema_1 = require("./schema");
const sourceInfo_1 = require("./sourceInfo");
const types_1 = require("./types");
const utils_1 = require("./utils");
const visit_1 = require("./visit");
function generateFile(ctx, fileDesc) {
    var _a;
    const { options, utils } = ctx;
    if (options.useOptionals === false) {
        console.warn("ts-proto: Passing useOptionals as a boolean option is deprecated and will be removed in a future version. Please pass the string 'none' instead of false.");
        options.useOptionals = "none";
    }
    else if (options.useOptionals === true) {
        console.warn("ts-proto: Passing useOptionals as a boolean option is deprecated and will be removed in a future version. Please pass the string 'messages' instead of true.");
        options.useOptionals = "messages";
    }
    // Google's protofiles are organized like Java, where package == the folder the file
    // is in, and file == a specific service within the package. I.e. you can have multiple
    // company/foo.proto and company/bar.proto files, where package would be 'company'.
    //
    // We'll match that structure by setting up the module path as:
    //
    // company/foo.proto --> company/foo.ts
    // company/bar.proto --> company/bar.ts
    //
    // We'll also assume that the fileDesc.name is already the `company/foo.proto` path, with
    // the package already implicitly in it, so we won't re-append/strip/etc. it out/back in.
    const suffix = `${options.fileSuffix}.ts`;
    const moduleName = fileDesc.name.replace(".proto", suffix);
    const chunks = [];
    // Indicate this file's source protobuf package for reflective use with google.protobuf.Any
    if (options.exportCommonSymbols) {
        chunks.push((0, ts_poet_1.code) `export const protobufPackage = '${fileDesc.package}';`);
    }
    // Syntax, unlike most fields, is not repeated and thus does not use an index
    const sourceInfo = sourceInfo_1.default.fromDescriptor(fileDesc);
    const headerComment = sourceInfo.lookup(sourceInfo_1.Fields.file.syntax, undefined);
    (0, utils_1.maybeAddComment)(options, headerComment, chunks, (_a = fileDesc.options) === null || _a === void 0 ? void 0 : _a.deprecated);
    // Apply formatting to methods here, so they propagate globally
    for (let svc of fileDesc.service) {
        for (let i = 0; i < svc.method.length; i++) {
            svc.method[i] = new utils_1.FormattedMethodDescriptor(svc.method[i], options);
        }
    }
    // first make all the type declarations
    (0, visit_1.visit)(fileDesc, sourceInfo, (fullName, message, sInfo, fullProtoTypeName) => {
        chunks.push(generateInterfaceDeclaration(ctx, fullName, message, sInfo, (0, utils_1.maybePrefixPackage)(fileDesc, fullProtoTypeName)));
    }, options, (fullName, enumDesc, sInfo) => {
        chunks.push((0, enums_1.generateEnum)(ctx, fullName, enumDesc, sInfo));
    });
    // If nestJs=true export [package]_PACKAGE_NAME and [service]_SERVICE_NAME const
    if (options.nestJs) {
        if (options.exportCommonSymbols) {
            const prefix = (0, case_1.camelToSnake)(fileDesc.package.replace(/\./g, "_"));
            chunks.push((0, ts_poet_1.code) `export const ${prefix}_PACKAGE_NAME = '${fileDesc.package}';`);
        }
        if (options.useDate === options_1.DateOption.DATE &&
            fileDesc.messageType.find((message) => message.field.find((field) => field.typeName === ".google.protobuf.Timestamp"))) {
            chunks.push(makeProtobufTimestampWrapper());
        }
    }
    // We add `nestJs` here because enough though it doesn't use our encode/decode methods
    // for most/vanilla messages, we do generate static wrap/unwrap methods for the special
    // Struct/Value/wrapper types and use the `wrappers[...]` to have NestJS know about them.
    if (options.outputEncodeMethods ||
        options.outputJsonMethods ||
        options.outputTypeAnnotations ||
        options.outputTypeRegistry ||
        options.nestJs) {
        // then add the encoder/decoder/base instance
        (0, visit_1.visit)(fileDesc, sourceInfo, (fullName, message, _sInfo, fullProtoTypeName) => {
            const fullTypeName = (0, utils_1.maybePrefixPackage)(fileDesc, fullProtoTypeName);
            const outputWrapAndUnwrap = (0, generate_struct_wrappers_1.isWrapperType)(fullTypeName);
            // Only decode, fromPartial, and wrap use the createBase method
            if ((options.outputEncodeMethods && options.outputEncodeMethods !== "encode-no-creation") ||
                options.outputPartialMethods ||
                outputWrapAndUnwrap) {
                chunks.push(generateBaseInstanceFactory(ctx, fullName, message, fullTypeName));
            }
            const staticMembers = [];
            if (options.outputTypeAnnotations || options.outputTypeRegistry) {
                staticMembers.push((0, ts_poet_1.code) `$type: '${fullTypeName}' as const`);
            }
            if (options.outputExtensions) {
                for (const extension of message.extension) {
                    const { name, type, extensionInfo } = generateExtension(ctx, message, extension);
                    staticMembers.push((0, ts_poet_1.code) `${name}: <${ctx.utils.Extension}<${type}>> ${extensionInfo}`);
                }
            }
            if (options.outputEncodeMethods) {
                if (options.outputEncodeMethods === true ||
                    options.outputEncodeMethods === "encode-only" ||
                    options.outputEncodeMethods === "encode-no-creation") {
                    staticMembers.push(generateEncode(ctx, fullName, message));
                    if (options.outputExtensions && options.unknownFields && message.extensionRange.length) {
                        staticMembers.push(generateSetExtension(ctx, fullName));
                    }
                }
                if (options.outputEncodeMethods === true || options.outputEncodeMethods === "decode-only") {
                    staticMembers.push(generateDecode(ctx, fullName, message));
                    if (options.outputExtensions && options.unknownFields && message.extensionRange.length) {
                        staticMembers.push(generateGetExtension(ctx, fullName));
                    }
                }
            }
            if (options.useAsyncIterable) {
                staticMembers.push((0, generate_async_iterable_1.generateEncodeTransform)(ctx.utils, fullName));
                staticMembers.push((0, generate_async_iterable_1.generateDecodeTransform)(ctx.utils, fullName));
            }
            if (options.outputJsonMethods) {
                if (options.outputJsonMethods === true || options.outputJsonMethods === "from-only") {
                    staticMembers.push(generateFromJson(ctx, fullName, fullTypeName, message));
                }
                if (options.outputJsonMethods === true || options.outputJsonMethods === "to-only") {
                    staticMembers.push(generateToJson(ctx, fullName, fullTypeName, message));
                }
            }
            if (options.outputPartialMethods) {
                staticMembers.push(generateFromPartial(ctx, fullName, message));
            }
            const structFieldNames = {
                nullValue: (0, case_1.maybeSnakeToCamel)("null_value", ctx.options),
                numberValue: (0, case_1.maybeSnakeToCamel)("number_value", ctx.options),
                stringValue: (0, case_1.maybeSnakeToCamel)("string_value", ctx.options),
                boolValue: (0, case_1.maybeSnakeToCamel)("bool_value", ctx.options),
                structValue: (0, case_1.maybeSnakeToCamel)("struct_value", ctx.options),
                listValue: (0, case_1.maybeSnakeToCamel)("list_value", ctx.options),
            };
            if (options.nestJs) {
                staticMembers.push(...(0, generate_struct_wrappers_1.generateWrapDeep)(ctx, fullTypeName, structFieldNames));
                staticMembers.push(...(0, generate_struct_wrappers_1.generateUnwrapDeep)(ctx, fullTypeName, structFieldNames));
            }
            else {
                staticMembers.push(...(0, generate_struct_wrappers_1.generateWrapShallow)(ctx, fullTypeName, structFieldNames));
                staticMembers.push(...(0, generate_struct_wrappers_1.generateUnwrapShallow)(ctx, fullTypeName, structFieldNames));
            }
            if (staticMembers.length > 0) {
                chunks.push((0, ts_poet_1.code) `
            export const ${(0, ts_poet_1.def)(fullName)} = {
              ${(0, ts_poet_1.joinCode)(staticMembers, { on: ",\n\n" })}
            };
          `);
            }
            if (options.outputTypeRegistry) {
                const messageTypeRegistry = (0, utils_1.impFile)(options, "messageTypeRegistry@./typeRegistry");
                chunks.push((0, ts_poet_1.code) `
            ${messageTypeRegistry}.set(${fullName}.$type, ${fullName});
          `);
            }
        }, options);
    }
    if (options.outputExtensions) {
        for (const extension of fileDesc.extension) {
            const { name, type, extensionInfo } = generateExtension(ctx, undefined, extension);
            chunks.push((0, ts_poet_1.code) `export const ${name}: ${ctx.utils.Extension}<${type}> = ${extensionInfo};`);
        }
    }
    if (options.nestJs) {
        if (fileDesc.messageType.find((message) => message.field.find(types_1.isStructType))) {
            chunks.push(makeProtobufStructWrapper(options));
        }
    }
    let hasServerStreamingMethods = false;
    let hasStreamingMethods = false;
    (0, visit_1.visitServices)(fileDesc, sourceInfo, (serviceDesc, sInfo) => {
        if (options.nestJs) {
            // NestJS is sufficiently different that we special case the client/server interfaces
            // generate nestjs grpc client interface
            chunks.push((0, generate_nestjs_1.generateNestjsServiceClient)(ctx, fileDesc, sInfo, serviceDesc));
            // and the service controller interface
            chunks.push((0, generate_nestjs_1.generateNestjsServiceController)(ctx, fileDesc, sInfo, serviceDesc));
            // generate nestjs grpc service controller decorator
            chunks.push((0, generate_nestjs_1.generateNestjsGrpcServiceMethodsDecorator)(ctx, serviceDesc));
            let serviceConstName = `${(0, case_1.camelToSnake)(serviceDesc.name)}_NAME`;
            if (!serviceDesc.name.toLowerCase().endsWith("service")) {
                serviceConstName = `${(0, case_1.camelToSnake)(serviceDesc.name)}_SERVICE_NAME`;
            }
            chunks.push((0, ts_poet_1.code) `export const ${serviceConstName} = "${serviceDesc.name}";`);
        }
        const uniqueServices = [...new Set(options.outputServices)].sort();
        uniqueServices.forEach((outputService) => {
            if (outputService === options_1.ServiceOption.GRPC) {
                chunks.push((0, generate_grpc_js_1.generateGrpcJsService)(ctx, fileDesc, sInfo, serviceDesc));
            }
            else if (outputService === options_1.ServiceOption.NICE_GRPC) {
                chunks.push((0, generate_nice_grpc_1.generateNiceGrpcService)(ctx, fileDesc, sInfo, serviceDesc));
            }
            else if (outputService === options_1.ServiceOption.GENERIC) {
                chunks.push((0, generate_generic_service_definition_1.generateGenericServiceDefinition)(ctx, fileDesc, sInfo, serviceDesc));
            }
            else if (outputService === options_1.ServiceOption.DEFAULT) {
                // This service could be Twirp or grpc-web or JSON (maybe). So far all of their
                // interfaces are fairly similar so we share the same service interface.
                chunks.push((0, generate_services_1.generateService)(ctx, fileDesc, sInfo, serviceDesc));
                if (options.outputClientImpl === true) {
                    chunks.push((0, generate_services_1.generateServiceClientImpl)(ctx, fileDesc, serviceDesc));
                }
                else if (options.outputClientImpl === "grpc-web") {
                    chunks.push((0, generate_grpc_web_1.generateGrpcClientImpl)(ctx, fileDesc, serviceDesc));
                    chunks.push((0, generate_grpc_web_1.generateGrpcServiceDesc)(fileDesc, serviceDesc));
                    serviceDesc.method.forEach((method) => {
                        if (!method.clientStreaming) {
                            chunks.push((0, generate_grpc_web_1.generateGrpcMethodDesc)(ctx, serviceDesc, method));
                        }
                        if (method.serverStreaming) {
                            hasServerStreamingMethods = true;
                        }
                    });
                }
            }
        });
        serviceDesc.method.forEach((methodDesc, _index) => {
            if (methodDesc.serverStreaming || methodDesc.clientStreaming) {
                hasStreamingMethods = true;
            }
        });
    });
    if (options.outputServices.includes(options_1.ServiceOption.DEFAULT) &&
        options.outputClientImpl &&
        fileDesc.service.length > 0) {
        if (options.outputClientImpl === true) {
            chunks.push((0, generate_services_1.generateRpcType)(ctx, hasStreamingMethods));
        }
        else if (options.outputClientImpl === "grpc-web") {
            chunks.push((0, generate_grpc_web_1.addGrpcWebMisc)(ctx, hasServerStreamingMethods));
        }
    }
    if (options.context) {
        chunks.push((0, generate_services_1.generateDataLoaderOptionsType)());
        chunks.push((0, generate_services_1.generateDataLoadersType)());
    }
    if (options.outputSchema) {
        chunks.push(...(0, schema_1.generateSchema)(ctx, fileDesc, sourceInfo));
    }
    // https://www.typescriptlang.org/docs/handbook/2/modules.html:
    // > In TypeScript, just as in ECMAScript 2015, any file containing a top-level import or export is considered a module.
    // > Conversely, a file without any top-level import or export declarations is treated as a script whose contents are available in the global scope (and therefore to modules as well).
    //
    // Thus, to mark an empty file a module, we need to add `export {}` to it.
    if (options.esModuleInterop && chunks.length === 0) {
        chunks.push((0, ts_poet_1.code) `export {};`);
    }
    chunks.push(...Object.values(utils).map((v) => {
        if (v instanceof ConditionalOutput_1.ConditionalOutput) {
            return (0, ts_poet_1.code) `${v.ifUsed}`;
        }
        else {
            return (0, ts_poet_1.code) ``;
        }
    }));
    // Finally, reset method definitions to their original state (unformatted)
    // This is mainly so that the `meta-typings` tests pass
    for (let svc of fileDesc.service) {
        for (let i = 0; i < svc.method.length; i++) {
            const methodInfo = svc.method[i];
            (0, utils_1.assertInstanceOf)(methodInfo, utils_1.FormattedMethodDescriptor);
            svc.method[i] = methodInfo.getSource();
        }
    }
    return [moduleName, (0, ts_poet_1.joinCode)(chunks, { on: "\n\n" })];
}
exports.generateFile = generateFile;
/** These are runtime utility methods used by the generated code. */
function makeUtils(options) {
    const bytes = makeByteUtils(options);
    const longs = makeLongUtils(options, bytes);
    return {
        ...bytes,
        ...makeDeepPartial(options, longs),
        ...makeObjectIdMethods(),
        ...makeTimestampMethods(options, longs, bytes),
        ...longs,
        ...makeComparisonUtils(),
        ...makeNiceGrpcServerStreamingMethodResult(options),
        ...makeGrpcWebErrorClass(bytes),
        ...makeExtensionClass(options),
        ...makeAssertionUtils(bytes),
    };
}
exports.makeUtils = makeUtils;
function makeProtobufTimestampWrapper() {
    const wrappers = (0, ts_poet_1.imp)("wrappers@protobufjs");
    return (0, ts_poet_1.code) `
      ${wrappers}['.google.protobuf.Timestamp'] = {
        fromObject(value: Date) {
          return {
            seconds: value.getTime() / 1000,
            nanos: (value.getTime() % 1000) * 1e6,
          };
        },
        toObject(message: { seconds: number; nanos: number }) {
          return new Date(message.seconds * 1000 + message.nanos / 1e6);
        },
      } as any;`;
}
function makeProtobufStructWrapper(options) {
    const wrappers = (0, ts_poet_1.imp)("wrappers@protobufjs");
    const Struct = (0, utils_1.impProto)(options, "google/protobuf/struct", "Struct");
    return (0, ts_poet_1.code) `
    ${wrappers}['.google.protobuf.Struct'] = {
      fromObject: ${Struct}.wrap,
      toObject: ${Struct}.unwrap,
    } as any;`;
}
function makeLongUtils(options, bytes) {
    // Regardless of which `forceLong` config option we're using, we always use
    // the `long` library to either represent or at least sanity-check 64-bit values
    const util = (0, utils_1.impFile)(options, `util@protobufjs/minimal`);
    const configure = (0, utils_1.impFile)(options, `configure@protobufjs/minimal`);
    const LongImp = (0, ts_poet_1.imp)("Long=long");
    // Instead of exposing `LongImp` directly, let callers think that they are getting the
    // `imp(Long)` but really it is that + our long initialization snippet. This means the
    // initialization code will only be emitted in files that actually use the Long import.
    const Long = (0, ts_poet_1.conditionalOutput)("Long", (0, ts_poet_1.code) `
      if (${util}.Long !== ${LongImp}) {
        ${util}.Long = ${LongImp} as any;
        ${configure}();
      }
    `);
    // TODO This is unused?
    const numberToLong = (0, ts_poet_1.conditionalOutput)("numberToLong", (0, ts_poet_1.code) `
      function numberToLong(number: number) {
        return ${Long}.fromNumber(number);
      }
    `);
    const longToString = (0, ts_poet_1.conditionalOutput)("longToString", (0, ts_poet_1.code) `
      function longToString(long: ${Long}) {
        return long.toString();
      }
    `);
    const longToBigint = (0, ts_poet_1.conditionalOutput)("longToBigint", (0, ts_poet_1.code) `
      function longToBigint(long: ${Long}) {
        return BigInt(long.toString());
      }
    `);
    const longToNumber = (0, ts_poet_1.conditionalOutput)("longToNumber", (0, ts_poet_1.code) `
      function longToNumber(long: ${Long}): number {
        if (long.gt(${bytes.globalThis}.Number.MAX_SAFE_INTEGER)) {
          throw new ${bytes.globalThis}.Error("Value is larger than Number.MAX_SAFE_INTEGER")
        }
        return long.toNumber();
      }
    `);
    return { numberToLong, longToNumber, longToString, longToBigint, Long };
}
function makeByteUtils(options) {
    const globalThisPolyfill = (0, ts_poet_1.conditionalOutput)("gt", (0, ts_poet_1.code) `
      declare const self: any | undefined;
      declare const window: any | undefined;
      declare const global: any | undefined;
      const gt: any = (() => {
        if (typeof globalThis !== "undefined") return globalThis;
        if (typeof self !== "undefined") return self;
        if (typeof window !== "undefined") return window;
        if (typeof global !== "undefined") return global;
        throw "Unable to locate global object";
      })();
    `);
    const globalThis = options.globalThisPolyfill ? globalThisPolyfill : (0, ts_poet_1.conditionalOutput)("globalThis", (0, ts_poet_1.code) ``);
    function getBytesFromBase64Snippet() {
        const bytesFromBase64NodeSnippet = (0, ts_poet_1.code) `
      return Uint8Array.from(${globalThis}.Buffer.from(b64, 'base64'));
    `;
        const bytesFromBase64BrowserSnippet = (0, ts_poet_1.code) `
      const bin = ${globalThis}.atob(b64);
      const arr = new Uint8Array(bin.length);
      for (let i = 0; i < bin.length; ++i) {
          arr[i] = bin.charCodeAt(i);
      }
      return arr;
    `;
        switch (options.env) {
            case options_1.EnvOption.NODE:
                return bytesFromBase64NodeSnippet;
            case options_1.EnvOption.BROWSER:
                return bytesFromBase64BrowserSnippet;
            default:
                return (0, ts_poet_1.code) `
        if ((${globalThis} as any).Buffer) {
          ${bytesFromBase64NodeSnippet}
          } else {
            ${bytesFromBase64BrowserSnippet}
          }
        `;
        }
    }
    const bytesFromBase64 = (0, ts_poet_1.conditionalOutput)("bytesFromBase64", (0, ts_poet_1.code) `
      function bytesFromBase64(b64: string): Uint8Array {
        ${getBytesFromBase64Snippet()}
      }
    `);
    function getBase64FromBytesSnippet() {
        const base64FromBytesNodeSnippet = (0, ts_poet_1.code) `
      return ${globalThis}.Buffer.from(arr).toString('base64');
    `;
        const base64FromBytesBrowserSnippet = (0, ts_poet_1.code) `
      const bin: string[] = [];
      arr.forEach((byte) => {
        bin.push(${globalThis}.String.fromCharCode(byte));
      });
      return ${globalThis}.btoa(bin.join(''));
    `;
        switch (options.env) {
            case options_1.EnvOption.NODE:
                return base64FromBytesNodeSnippet;
            case options_1.EnvOption.BROWSER:
                return base64FromBytesBrowserSnippet;
            default:
                return (0, ts_poet_1.code) `
          if ((${globalThis} as any).Buffer) {
            ${base64FromBytesNodeSnippet}
          } else {
            ${base64FromBytesBrowserSnippet}
          }
        `;
        }
    }
    const base64FromBytes = (0, ts_poet_1.conditionalOutput)("base64FromBytes", (0, ts_poet_1.code) `
      function base64FromBytes(arr: Uint8Array): string {
        ${getBase64FromBytesSnippet()}
      }
    `);
    return { globalThis, bytesFromBase64, base64FromBytes };
}
function makeDeepPartial(options, longs) {
    let oneofCase = "";
    if (options.oneof === options_1.OneofOption.UNIONS) {
        oneofCase = `
      : T extends { ${maybeReadonly(options)}$case: string }
      ? { [K in keyof Omit<T, '$case'>]?: DeepPartial<T[K]> } & { ${maybeReadonly(options)}$case: T['$case'] }
    `;
    }
    const maybeExport = options.exportCommonSymbols ? "export" : "";
    // Allow passing longs as numbers or strings, nad we'll convert them
    const maybeLong = options.forceLong === options_1.LongOption.LONG ? (0, ts_poet_1.code) ` : T extends ${longs.Long} ? string | number | Long ` : "";
    const Builtin = (0, ts_poet_1.conditionalOutput)("Builtin", (0, ts_poet_1.code) `type Builtin = Date | Function | Uint8Array | string | number | boolean |${options.forceLong === options_1.LongOption.BIGINT ? " bigint |" : ""} undefined;`);
    // Based on https://github.com/sindresorhus/type-fest/pull/259
    const maybeExcludeType = (0, options_1.addTypeToMessages)(options) ? `| '$type'` : "";
    const Exact = (0, ts_poet_1.conditionalOutput)("Exact", (0, ts_poet_1.code) `
      type KeysOfUnion<T> = T extends T ? keyof T : never;
      ${maybeExport} type Exact<P, I extends P> = P extends ${Builtin}
        ? P
        : P &
        { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P> ${maybeExcludeType}>]: never };
    `);
    // Based on the type from ts-essentials
    const keys = (0, options_1.addTypeToMessages)(options) ? (0, ts_poet_1.code) `Exclude<keyof T, '$type'>` : (0, ts_poet_1.code) `keyof T`;
    const DeepPartial = (0, ts_poet_1.conditionalOutput)("DeepPartial", (0, ts_poet_1.code) `
      ${maybeExport} type DeepPartial<T> =  T extends ${Builtin}
        ? T
        ${maybeLong}
        : T extends globalThis.Array<infer U>
        ? globalThis.Array<DeepPartial<U>>
        : T extends ReadonlyArray<infer U>
        ? ReadonlyArray<DeepPartial<U>>${oneofCase}
        : T extends {}
        ? { [K in ${keys}]?: DeepPartial<T[K]> }
        : Partial<T>;
    `);
    return { Builtin, DeepPartial, Exact };
}
function makeObjectIdMethods() {
    const mongodb = (0, ts_poet_1.imp)("mongodb*mongodb");
    const fromProtoObjectId = (0, ts_poet_1.conditionalOutput)("fromProtoObjectId", (0, ts_poet_1.code) `
      function fromProtoObjectId(oid: ObjectId): ${mongodb}.ObjectId {
        return new ${mongodb}.ObjectId(oid.value);
      }
    `);
    const fromJsonObjectId = (0, ts_poet_1.conditionalOutput)("fromJsonObjectId", (0, ts_poet_1.code) `
      function fromJsonObjectId(o: any): ${mongodb}.ObjectId {
        if (o instanceof ${mongodb}.ObjectId) {
          return o;
        } else if (typeof o === "string") {
          return new ${mongodb}.ObjectId(o);
        } else {
          return ${fromProtoObjectId}(ObjectId.fromJSON(o));
        }
      }
    `);
    const toProtoObjectId = (0, ts_poet_1.conditionalOutput)("toProtoObjectId", (0, ts_poet_1.code) `
      function toProtoObjectId(oid: ${mongodb}.ObjectId): ObjectId {
        const value = oid.toString();
        return { value };
      }
    `);
    return { fromJsonObjectId, fromProtoObjectId, toProtoObjectId };
}
function makeTimestampMethods(options, longs, bytes) {
    const Timestamp = (0, utils_1.impProto)(options, "google/protobuf/timestamp", "Timestamp");
    const NanoDate = (0, ts_poet_1.imp)("NanoDate=nano-date");
    let seconds = "Math.trunc(date.getTime() / 1_000)";
    let toNumberCode = "t.seconds";
    const makeToNumberCode = (methodCall) => `t.seconds${options.useOptionals === "all" ? "?" : ""}.${methodCall}`;
    if (options.forceLong === options_1.LongOption.LONG) {
        toNumberCode = makeToNumberCode("toNumber()");
        seconds = (0, ts_poet_1.code) `${longs.numberToLong}(${seconds})`;
    }
    else if (options.forceLong === options_1.LongOption.BIGINT) {
        toNumberCode = (0, ts_poet_1.code) `${bytes.globalThis}.Number(${makeToNumberCode("toString()")})`;
        seconds = (0, ts_poet_1.code) `BigInt(${seconds})`;
    }
    else if (options.forceLong === options_1.LongOption.STRING) {
        toNumberCode = (0, ts_poet_1.code) `${bytes.globalThis}.Number(t.seconds)`;
        seconds = (0, ts_poet_1.code) `${seconds}.toString()`;
    }
    const maybeTypeField = (0, options_1.addTypeToMessages)(options) ? `$type: 'google.protobuf.Timestamp',` : "";
    const toTimestamp = (0, ts_poet_1.conditionalOutput)("toTimestamp", options.useDate === options_1.DateOption.STRING
        ? (0, ts_poet_1.code) `
          function toTimestamp(dateStr: string): ${Timestamp} {
            const date = new ${bytes.globalThis}.Date(dateStr);
            const seconds = ${seconds};
            const nanos = (date.getTime() % 1_000) * 1_000_000;
            return { ${maybeTypeField} seconds, nanos };
          }
        `
        : options.useDate === options_1.DateOption.STRING_NANO
            ? (0, ts_poet_1.code) `
          function toTimestamp(dateStr: string): ${Timestamp} {
            const nanoDate = new ${NanoDate}(dateStr);

            const date = {
              getTime: (): number => nanoDate.valueOf(),
            } as const;
            const seconds = ${seconds};

            let nanos = nanoDate.getMilliseconds() * 1_000_000;
            nanos += nanoDate.getMicroseconds() * 1_000;
            nanos += nanoDate.getNanoseconds();

            return { ${maybeTypeField} seconds, nanos };
          }
        `
            : (0, ts_poet_1.code) `
          function toTimestamp(date: Date): ${Timestamp} {
            const seconds = ${seconds};
            const nanos = (date.getTime() % 1_000) * 1_000_000;
            return { ${maybeTypeField} seconds, nanos };
          }
        `);
    const fromTimestamp = (0, ts_poet_1.conditionalOutput)("fromTimestamp", options.useDate === options_1.DateOption.STRING
        ? (0, ts_poet_1.code) `
          function fromTimestamp(t: ${Timestamp}): string {
            let millis = (${toNumberCode} || 0) * 1_000;
            millis += (t.nanos || 0) / 1_000_000;
            return new ${bytes.globalThis}.Date(millis).toISOString();
          }
        `
        : options.useDate === options_1.DateOption.STRING_NANO
            ? (0, ts_poet_1.code) `
          function fromTimestamp(t: ${Timestamp}): string {
            const seconds = ${toNumberCode} || 0;
            const nanos = (t.nanos || 0) % 1_000;
            const micros = Math.trunc(((t.nanos || 0) % 1_000_000) / 1_000)
            let millis = seconds * 1_000;
            millis += Math.trunc((t.nanos || 0) / 1_000_000);

            const nanoDate = new ${NanoDate}(millis);
            nanoDate.setMicroseconds(micros);
            nanoDate.setNanoseconds(nanos);

            return nanoDate.toISOStringFull();
          }
        `
            : (0, ts_poet_1.code) `
          function fromTimestamp(t: ${Timestamp}): Date {
            let millis = (${toNumberCode} || 0) * 1_000;
            millis += (t.nanos || 0) / 1_000_000;
            return new ${bytes.globalThis}.Date(millis);
          }
        `);
    const fromJsonTimestamp = (0, ts_poet_1.conditionalOutput)("fromJsonTimestamp", options.useDate === options_1.DateOption.DATE
        ? (0, ts_poet_1.code) `
        function fromJsonTimestamp(o: any): Date {
          if (o instanceof ${bytes.globalThis}.Date) {
            return o;
          } else if (typeof o === "string") {
            return new ${bytes.globalThis}.Date(o);
          } else {
            return ${fromTimestamp}(Timestamp.fromJSON(o));
          }
        }
      `
        : (0, ts_poet_1.code) `
        function fromJsonTimestamp(o: any): Timestamp {
          if (o instanceof ${bytes.globalThis}.Date) {
            return ${toTimestamp}(o);
          } else if (typeof o === "string") {
            return ${toTimestamp}(new ${bytes.globalThis}.Date(o));
          } else {
            return Timestamp.fromJSON(o);
          }
        }
      `);
    return { toTimestamp, fromTimestamp, fromJsonTimestamp };
}
function makeComparisonUtils() {
    const isObject = (0, ts_poet_1.conditionalOutput)("isObject", (0, ts_poet_1.code) `
    function isObject(value: any): boolean {
      return typeof value === 'object' && value !== null;
    }`);
    const isSet = (0, ts_poet_1.conditionalOutput)("isSet", (0, ts_poet_1.code) `
    function isSet(value: any): boolean {
      return value !== null && value !== undefined;
    }`);
    return { isObject, isSet };
}
function makeNiceGrpcServerStreamingMethodResult(options) {
    const NiceGrpcServerStreamingMethodResult = (0, ts_poet_1.conditionalOutput)("ServerStreamingMethodResult", options.outputIndex
        ? (0, ts_poet_1.code) `
        type ServerStreamingMethodResult<Response> = {
          [Symbol.asyncIterator](): AsyncIterator<Response, void>;
        };
      `
        : (0, ts_poet_1.code) `
        export type ServerStreamingMethodResult<Response> = {
          [Symbol.asyncIterator](): AsyncIterator<Response, void>;
        };
      `);
    return { NiceGrpcServerStreamingMethodResult };
}
function makeGrpcWebErrorClass(bytes) {
    const GrpcWebError = (0, ts_poet_1.conditionalOutput)("GrpcWebError", (0, ts_poet_1.code) `
      export class GrpcWebError extends ${bytes.globalThis}.Error {
        constructor(message: string, public code: grpc.Code, public metadata: grpc.Metadata) {
          super(message);
        }
      }
    `);
    return { GrpcWebError };
}
function makeExtensionClass(options) {
    const Reader = (0, utils_1.impFile)(options, "Reader@protobufjs/minimal");
    const Writer = (0, utils_1.impFile)(options, "Writer@protobufjs/minimal");
    const Extension = (0, ts_poet_1.conditionalOutput)("Extension", (0, ts_poet_1.code) `
      export interface Extension <T> {
        number: number;
        tag: number;
        singularTag?: number;
        encode?: (message: T) => Uint8Array[];
        decode?: (tag: number, input: Uint8Array[]) => T;
        repeated: boolean;
        packed: boolean;
      }
    `);
    return { Extension };
}
function makeAssertionUtils(bytes) {
    const fail = (0, ts_poet_1.conditionalOutput)("fail", (0, ts_poet_1.code) `
      function fail(message?: string): never {
        throw new ${bytes.globalThis}.Error(message ?? "Failed");
      }
    `);
    return { fail };
}
// Create the interface with properties
function generateInterfaceDeclaration(ctx, fullName, messageDesc, sourceInfo, fullTypeName) {
    var _a;
    const { options, currentFile } = ctx;
    const chunks = [];
    (0, utils_1.maybeAddComment)(options, sourceInfo, chunks, (_a = messageDesc.options) === null || _a === void 0 ? void 0 : _a.deprecated);
    // interface name should be defined to avoid import collisions
    chunks.push((0, ts_poet_1.code) `export interface ${(0, ts_poet_1.def)(fullName)} {`);
    if ((0, options_1.addTypeToMessages)(options)) {
        chunks.push((0, ts_poet_1.code) `$type${options.outputTypeAnnotations === "optional" ? "?" : ""}: '${fullTypeName}',`);
    }
    // When oneof=unions, we generate a single property with an ADT per `oneof` clause.
    const processedOneofs = new Set();
    messageDesc.field.forEach((fieldDesc, index) => {
        var _a;
        if ((0, types_1.isWithinOneOfThatShouldBeUnion)(options, fieldDesc)) {
            const { oneofIndex } = fieldDesc;
            if (!processedOneofs.has(oneofIndex)) {
                processedOneofs.add(oneofIndex);
                chunks.push(generateOneofProperty(ctx, messageDesc, oneofIndex, sourceInfo));
            }
            return;
        }
        const info = sourceInfo.lookup(sourceInfo_1.Fields.message.field, index);
        (0, utils_1.maybeAddComment)(options, info, chunks, (_a = fieldDesc.options) === null || _a === void 0 ? void 0 : _a.deprecated);
        const fieldKey = (0, utils_1.safeAccessor)((0, utils_1.getFieldName)(fieldDesc, options));
        const isOptional = (0, types_1.isOptionalProperty)(fieldDesc, messageDesc.options, options, currentFile.isProto3Syntax);
        const type = (0, types_1.toTypeName)(ctx, messageDesc, fieldDesc, isOptional);
        chunks.push((0, ts_poet_1.code) `${maybeReadonly(options)}${fieldKey}${isOptional ? "?" : ""}: ${type}, `);
    });
    if (ctx.options.unknownFields) {
        chunks.push((0, ts_poet_1.code) `_unknownFields?: {[key: number]: Uint8Array[]} | undefined,`);
    }
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
function generateOneofProperty(ctx, messageDesc, oneofIndex, sourceInfo) {
    const { options } = ctx;
    const fields = messageDesc.field.filter((field) => (0, types_1.isWithinOneOf)(field) && field.oneofIndex === oneofIndex);
    const mbReadonly = maybeReadonly(options);
    const unionType = (0, ts_poet_1.joinCode)(fields.map((f) => {
        let fieldName = (0, case_1.maybeSnakeToCamel)(f.name, options);
        let typeName = (0, types_1.toTypeName)(ctx, messageDesc, f);
        return (0, ts_poet_1.code) `{ ${mbReadonly}$case: '${fieldName}', ${mbReadonly}${fieldName}: ${typeName} }`;
    }), { on: " | " });
    const name = (0, case_1.maybeSnakeToCamel)(messageDesc.oneofDecl[oneofIndex].name, options);
    return (0, ts_poet_1.code) `${mbReadonly}${name}?: ${unionType} | ${(0, utils_1.nullOrUndefined)(options)},`;
    /*
    // Ideally we'd put the comments for each oneof field next to the anonymous
    // type we've created in the type union above, but ts-poet currently lacks
    // that ability. For now just concatenate all comments into one big one.
    let comments: Array<string> = [];
    const info = sourceInfo.lookup(Fields.message.oneof_decl, oneofIndex);
    maybeAddComment(options, info, (text) => comments.push(text));
    messageDesc.field.forEach((field, index) => {
      if (!isWithinOneOf(field) || field.oneofIndex !== oneofIndex) {
        return;
      }
      const info = sourceInfo.lookup(Fields.message.field, index);
      const name = maybeSnakeToCamel(field.name, options);
      maybeAddComment(options, info, (text) => comments.push(name + '\n' + text));
    });
    if (comments.length) {
      prop = prop.addJavadoc(comments.join('\n'));
    }
    return prop;
    */
}
// Create a function that constructs 'base' instance with default values for decode to use as a prototype
function generateBaseInstanceFactory(ctx, fullName, messageDesc, fullTypeName) {
    const { options, currentFile } = ctx;
    const fields = [];
    // When oneof=unions, we generate a single property with an ADT per `oneof` clause.
    const processedOneofs = new Set();
    for (const field of messageDesc.field) {
        if ((0, types_1.isWithinOneOfThatShouldBeUnion)(ctx.options, field)) {
            const { oneofIndex } = field;
            if (!processedOneofs.has(oneofIndex)) {
                processedOneofs.add(oneofIndex);
                const name = options.useJsonName
                    ? (0, utils_1.getFieldName)(field, options)
                    : (0, case_1.maybeSnakeToCamel)(messageDesc.oneofDecl[oneofIndex].name, ctx.options);
                fields.push((0, ts_poet_1.code) `${(0, utils_1.safeAccessor)(name)}: ${(0, utils_1.nullOrUndefined)(options)}`);
            }
            continue;
        }
        if (!options.initializeFieldsAsUndefined &&
            (0, types_1.isOptionalProperty)(field, messageDesc.options, options, currentFile.isProto3Syntax)) {
            continue;
        }
        const fieldKey = (0, utils_1.safeAccessor)((0, utils_1.getFieldName)(field, options));
        const val = (0, types_1.isWithinOneOf)(field)
            ? (0, utils_1.nullOrUndefined)(options)
            : (0, types_1.isMapType)(ctx, messageDesc, field)
                ? (0, types_1.shouldGenerateJSMapType)(ctx, messageDesc, field)
                    ? "new Map()"
                    : "{}"
                : (0, types_1.isRepeated)(field)
                    ? "[]"
                    : (0, types_1.defaultValue)(ctx, field);
        fields.push((0, ts_poet_1.code) `${fieldKey}: ${val}`);
    }
    if ((0, options_1.addTypeToMessages)(options)) {
        fields.unshift((0, ts_poet_1.code) `$type: '${fullTypeName}'`);
    }
    if (ctx.options.unknownFields && ctx.options.initializeFieldsAsUndefined) {
        fields.push((0, ts_poet_1.code) `_unknownFields: {}`);
    }
    return (0, ts_poet_1.code) `
    function createBase${fullName}(): ${fullName} {
      return { ${(0, ts_poet_1.joinCode)(fields, { on: "," })} };
    }
  `;
}
function getDecodeReadSnippet(ctx, field) {
    const { options, utils } = ctx;
    let readSnippet;
    if ((0, types_1.isPrimitive)(field)) {
        readSnippet = (0, ts_poet_1.code) `reader.${(0, types_1.toReaderCall)(field)}()`;
        if ((0, types_1.isBytes)(field)) {
            if (options.env === options_1.EnvOption.NODE) {
                readSnippet = (0, ts_poet_1.code) `${readSnippet} as Buffer`;
            }
        }
        else if ((0, types_1.basicLongWireType)(field.type) !== undefined) {
            if ((0, types_1.isJsTypeFieldOption)(options, field)) {
                switch (field.options.jstype) {
                    case ts_proto_descriptors_1.FieldOptions_JSType.JS_NUMBER:
                        readSnippet = (0, ts_poet_1.code) `${utils.longToNumber}(${readSnippet} as Long)`;
                        break;
                    case ts_proto_descriptors_1.FieldOptions_JSType.JS_STRING:
                        readSnippet = (0, ts_poet_1.code) `${utils.longToString}(${readSnippet} as Long)`;
                        break;
                }
            }
            else if (options.forceLong === options_1.LongOption.LONG) {
                readSnippet = (0, ts_poet_1.code) `${readSnippet} as Long`;
            }
            else if (options.forceLong === options_1.LongOption.STRING) {
                readSnippet = (0, ts_poet_1.code) `${utils.longToString}(${readSnippet} as Long)`;
            }
            else if (options.forceLong === options_1.LongOption.BIGINT) {
                readSnippet = (0, ts_poet_1.code) `${utils.longToBigint}(${readSnippet} as Long)`;
            }
            else {
                readSnippet = (0, ts_poet_1.code) `${utils.longToNumber}(${readSnippet} as Long)`;
            }
        }
        else if ((0, types_1.isEnum)(field)) {
            if (options.stringEnums) {
                const fromJson = (0, types_1.getEnumMethod)(ctx, field.typeName, "FromJSON");
                readSnippet = (0, ts_poet_1.code) `${fromJson}(${readSnippet})`;
            }
            else {
                readSnippet = (0, ts_poet_1.code) `${readSnippet} as any`;
            }
        }
    }
    else if ((0, types_1.isValueType)(ctx, field)) {
        const type = (0, types_1.basicTypeName)(ctx, field, { keepValueType: true });
        const unwrap = (decodedValue) => {
            if ((0, types_1.isListValueType)(field) || (0, types_1.isStructType)(field) || (0, types_1.isAnyValueType)(field) || (0, types_1.isFieldMaskType)(field)) {
                return (0, ts_poet_1.code) `${type}.unwrap(${decodedValue})`;
            }
            return (0, ts_poet_1.code) `${decodedValue}.value`;
        };
        const decoder = (0, ts_poet_1.code) `${type}.decode(reader, reader.uint32())`;
        readSnippet = (0, ts_poet_1.code) `${unwrap(decoder)}`;
    }
    else if ((0, types_1.isTimestamp)(field) &&
        (options.useDate === options_1.DateOption.DATE ||
            options.useDate === options_1.DateOption.STRING ||
            options.useDate === options_1.DateOption.STRING_NANO)) {
        const type = (0, types_1.basicTypeName)(ctx, field, { keepValueType: true });
        readSnippet = (0, ts_poet_1.code) `${utils.fromTimestamp}(${type}.decode(reader, reader.uint32()))`;
    }
    else if ((0, types_1.isObjectId)(field) && options.useMongoObjectId) {
        const type = (0, types_1.basicTypeName)(ctx, field, { keepValueType: true });
        readSnippet = (0, ts_poet_1.code) `${utils.fromProtoObjectId}(${type}.decode(reader, reader.uint32()))`;
    }
    else if ((0, types_1.isMessage)(field)) {
        const type = (0, types_1.basicTypeName)(ctx, field);
        if (field.type == ts_proto_descriptors_1.FieldDescriptorProto_Type.TYPE_GROUP) {
            readSnippet = (0, ts_poet_1.code) `${type}.decode(reader)`;
        }
        else {
            readSnippet = (0, ts_poet_1.code) `${type}.decode(reader, reader.uint32())`;
        }
    }
    else {
        throw new Error(`Unhandled field ${field}`);
    }
    return readSnippet;
}
/** Creates a function to decode a message by loop overing the tags. */
function generateDecode(ctx, fullName, messageDesc) {
    const { options, currentFile } = ctx;
    const chunks = [];
    let createBase = (0, ts_poet_1.code) `createBase${fullName}()`;
    if (options.usePrototypeForDefaults) {
        createBase = (0, ts_poet_1.code) `Object.create(${createBase}) as ${fullName}`;
    }
    const Reader = (0, utils_1.impFile)(ctx.options, "Reader@protobufjs/minimal");
    // create the basic function declaration
    chunks.push((0, ts_poet_1.code) `
    decode(
      input: ${Reader} | Uint8Array,
      length?: number,
    ): ${fullName} {
      const reader = input instanceof ${Reader} ? input : ${Reader}.create(input);
      let end = length === undefined ? reader.len : reader.pos + length;
  `);
    chunks.push((0, ts_poet_1.code) `const message = ${createBase}${maybeAsAny(options)};`);
    // start the tag loop
    chunks.push((0, ts_poet_1.code) `
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
  `);
    // add a case for each incoming field
    messageDesc.field.forEach((field) => {
        const fieldName = (0, utils_1.getFieldName)(field, options);
        const messageProperty = (0, utils_1.getPropertyAccessor)("message", fieldName);
        chunks.push((0, ts_poet_1.code) `case ${field.number}:`);
        const tag = ((field.number << 3) | (0, types_1.basicWireType)(field.type)) >>> 0;
        const tagCheck = (0, ts_poet_1.code) `
      if (tag !== ${tag}) {
        break;
      }
    `;
        // get a generic 'reader.doSomething' bit that is specific to the basic type
        const readSnippet = getDecodeReadSnippet(ctx, field);
        // and then use the snippet to handle repeated fields if necessary
        const initializerNecessary = !options.initializeFieldsAsUndefined &&
            (0, types_1.isOptionalProperty)(field, messageDesc.options, options, currentFile.isProto3Syntax);
        if ((0, types_1.isRepeated)(field)) {
            const maybeNonNullAssertion = ctx.options.useOptionals === "all" || ctx.options.useOptionals === "deprecatedOnly" ? "!" : "";
            const mapType = (0, types_1.detectMapType)(ctx, messageDesc, field);
            if (mapType) {
                // We need a unique const within the `cast` statement
                const varName = `entry${field.number}`;
                const generateMapType = (0, types_1.shouldGenerateJSMapType)(ctx, messageDesc, field);
                let valueSetterSnippet;
                if (generateMapType) {
                    valueSetterSnippet = `${messageProperty}${maybeNonNullAssertion}.set(${varName}.key, ${varName}.value)`;
                }
                else {
                    valueSetterSnippet = `${messageProperty}${maybeNonNullAssertion}[${varName}.key] = ${varName}.value`;
                }
                const initializerSnippet = initializerNecessary
                    ? `
            if (${messageProperty} === undefined ${(0, utils_1.withOrMaybeCheckIsNull)(options, messageProperty)}) {
              ${messageProperty} = ${generateMapType ? "new Map()" : "{}"};
            }`
                    : "";
                chunks.push((0, ts_poet_1.code) `
          ${tagCheck}
          const ${varName} = ${readSnippet};
          if (${varName}.value !== undefined ${(0, utils_1.withAndMaybeCheckIsNotNull)(options, `${varName}.value`)}) {
            ${initializerSnippet}
            ${valueSetterSnippet};
          }
        `);
            }
            else {
                const initializerSnippet = initializerNecessary
                    ? `
            if (${messageProperty} === undefined ${(0, utils_1.withOrMaybeCheckIsNull)(options, messageProperty)}) {
              ${messageProperty} = [];
            }`
                    : "";
                if ((0, types_1.packedType)(field.type) === undefined) {
                    chunks.push((0, ts_poet_1.code) `
            ${tagCheck}
            ${initializerSnippet}
            ${messageProperty}${maybeNonNullAssertion}.push(${readSnippet});
          `);
                }
                else {
                    const packedTag = ((field.number << 3) | 2) >>> 0;
                    chunks.push((0, ts_poet_1.code) `
            if (tag === ${tag}) {
              ${initializerSnippet}
              ${messageProperty}${maybeNonNullAssertion}.push(${readSnippet});

              continue;
            }

            if (tag === ${packedTag}) {
              ${initializerSnippet}
              const end2 = reader.uint32() + reader.pos;
              while (reader.pos < end2) {
                ${messageProperty}${maybeNonNullAssertion}.push(${readSnippet});
              }

              continue;
            }

            break;
          `);
                }
            }
        }
        else if ((0, types_1.isWithinOneOfThatShouldBeUnion)(options, field)) {
            const oneofNameWithMessage = options.useJsonName
                ? messageProperty
                : (0, utils_1.getPropertyAccessor)("message", (0, case_1.maybeSnakeToCamel)(messageDesc.oneofDecl[field.oneofIndex].name, options));
            chunks.push((0, ts_poet_1.code) `
        ${tagCheck}
        ${oneofNameWithMessage} = { $case: '${fieldName}', ${fieldName}: ${readSnippet} };
      `);
        }
        else {
            chunks.push((0, ts_poet_1.code) `
        ${tagCheck}
        ${messageProperty} = ${readSnippet};
      `);
        }
        if (!(0, types_1.isRepeated)(field) || (0, types_1.packedType)(field.type) === undefined) {
            chunks.push((0, ts_poet_1.code) `continue;`);
        }
    });
    chunks.push((0, ts_poet_1.code) `}`);
    chunks.push((0, ts_poet_1.code) `
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
  `);
    if (options.unknownFields) {
        let unknownFieldsInitializerSnippet = "";
        let maybeNonNullAssertion = options.initializeFieldsAsUndefined ? "!" : "";
        if (!options.initializeFieldsAsUndefined) {
            unknownFieldsInitializerSnippet = `
        if (message._unknownFields === undefined ${(0, utils_1.withOrMaybeCheckIsNull)(options, `message._unknownFields`)}) {
          message._unknownFields = {};
        }
      `;
        }
        chunks.push((0, ts_poet_1.code) `
      const startPos = reader.pos;
      reader.skipType(tag & 7);
      const buf = reader.buf.slice(startPos, reader.pos);

      ${unknownFieldsInitializerSnippet}
      const list = message._unknownFields${maybeNonNullAssertion}[tag];

      if (list === undefined ${(0, utils_1.withOrMaybeCheckIsNull)(options, `message._unknownFields`)}) {
        message._unknownFields${maybeNonNullAssertion}[tag] = [buf];
      } else {
        list.push(buf);
      }
    `);
    }
    else {
        chunks.push((0, ts_poet_1.code) `
        reader.skipType(tag & 7);
    `);
    }
    // and then wrap up the while/return
    chunks.push((0, ts_poet_1.code) `}`);
    chunks.push((0, ts_poet_1.code) `return message;`);
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
/** Returns a generic writer.doSomething based on the basic type */
function getEncodeWriteSnippet(ctx, field) {
    const { options, utils } = ctx;
    if ((0, types_1.isEnum)(field) && options.stringEnums) {
        const tag = ((field.number << 3) | (0, types_1.basicWireType)(field.type)) >>> 0;
        const toNumber = (0, types_1.getEnumMethod)(ctx, field.typeName, "ToNumber");
        return (place) => (0, ts_poet_1.code) `writer.uint32(${tag}).${(0, types_1.toReaderCall)(field)}(${toNumber}(${place}))`;
    }
    else if ((0, types_1.isLong)(field) && options.forceLong === options_1.LongOption.BIGINT) {
        const tag = ((field.number << 3) | (0, types_1.basicWireType)(field.type)) >>> 0;
        const fieldType = (0, types_1.toReaderCall)(field);
        switch (fieldType) {
            case "int64":
            case "sint64":
            case "sfixed64":
                return (place, placeAlt) => (0, ts_poet_1.code) `if (BigInt.asIntN(64, ${place}) !== ${placeAlt !== null && placeAlt !== void 0 ? placeAlt : place}) {
          throw new ${utils.globalThis}.Error('value provided for field ${place} of type ${fieldType} too large');
        }
        writer.uint32(${tag}).${(0, types_1.toReaderCall)(field)}(${place}.toString())`;
            case "uint64":
            case "fixed64":
                return (place, placeAlt) => (0, ts_poet_1.code) `if (BigInt.asUintN(64, ${place}) !== ${placeAlt !== null && placeAlt !== void 0 ? placeAlt : place}) {
          throw new ${utils.globalThis}.Error('value provided for field ${place} of type ${fieldType} too large');
        }
        writer.uint32(${tag}).${(0, types_1.toReaderCall)(field)}(${place}.toString())`;
            default:
                throw new Error(`unexpected BigInt type: ${fieldType}`);
        }
    }
    else if ((0, types_1.isScalar)(field) || (0, types_1.isEnum)(field)) {
        const tag = ((field.number << 3) | (0, types_1.basicWireType)(field.type)) >>> 0;
        return (place) => (0, ts_poet_1.code) `writer.uint32(${tag}).${(0, types_1.toReaderCall)(field)}(${place})`;
    }
    else if ((0, types_1.isObjectId)(field) && options.useMongoObjectId) {
        const tag = ((field.number << 3) | 2) >>> 0;
        const type = (0, types_1.basicTypeName)(ctx, field, { keepValueType: true });
        return (place) => (0, ts_poet_1.code) `${type}.encode(${utils.toProtoObjectId}(${place}), writer.uint32(${tag}).fork()).ldelim()`;
    }
    else if ((0, types_1.isTimestamp)(field) &&
        (options.useDate === options_1.DateOption.DATE ||
            options.useDate === options_1.DateOption.STRING ||
            options.useDate === options_1.DateOption.STRING_NANO)) {
        const tag = ((field.number << 3) | 2) >>> 0;
        const type = (0, types_1.basicTypeName)(ctx, field, { keepValueType: true });
        return (place) => (0, ts_poet_1.code) `${type}.encode(${utils.toTimestamp}(${place}), writer.uint32(${tag}).fork()).ldelim()`;
    }
    else if ((0, types_1.isValueType)(ctx, field)) {
        const maybeTypeField = (0, options_1.addTypeToMessages)(options) ? `$type: '${field.typeName.slice(1)}',` : "";
        const type = (0, types_1.basicTypeName)(ctx, field, { keepValueType: true });
        const wrappedValue = (place) => {
            if ((0, types_1.isAnyValueType)(field) || (0, types_1.isListValueType)(field) || (0, types_1.isStructType)(field) || (0, types_1.isFieldMaskType)(field)) {
                return (0, ts_poet_1.code) `${type}.wrap(${place})`;
            }
            return (0, ts_poet_1.code) `{${maybeTypeField} value: ${place}!}`;
        };
        const tag = ((field.number << 3) | 2) >>> 0;
        return (place) => (0, ts_poet_1.code) `${type}.encode(${wrappedValue(place)}, writer.uint32(${tag}).fork()).ldelim()`;
    }
    else if ((0, types_1.isMessage)(field)) {
        const type = (0, types_1.basicTypeName)(ctx, field);
        if (field.type == ts_proto_descriptors_1.FieldDescriptorProto_Type.TYPE_GROUP) {
            const startTag = ((field.number << 3) | 3) >>> 0, endTag = ((field.number << 3) | 4) >>> 0;
            return (place) => (0, ts_poet_1.code) `${type}.encode(${place}, writer.uint32(${startTag})).uint32(${endTag})`;
        }
        const tag = ((field.number << 3) | 2) >>> 0;
        return (place) => (0, ts_poet_1.code) `${type}.encode(${place}, writer.uint32(${tag}).fork()).ldelim()`;
    }
    else {
        throw new Error(`Unhandled field ${field}`);
    }
}
/** Creates a function to encode a message by loop overing the tags. */
function generateEncode(ctx, fullName, messageDesc) {
    const { options, utils, typeMap, currentFile } = ctx;
    const chunks = [];
    const Writer = (0, utils_1.impFile)(ctx.options, "Writer@protobufjs/minimal");
    // create the basic function declaration
    chunks.push((0, ts_poet_1.code) `
    encode(
      ${messageDesc.field.length > 0 || options.unknownFields ? "message" : "_"}: ${fullName},
      writer: ${Writer} = ${Writer}.create(),
    ): ${Writer} {
  `);
    const processedOneofs = new Set();
    const oneOfFieldsDict = messageDesc.field
        .filter((field) => (0, types_1.isWithinOneOfThatShouldBeUnion)(options, field))
        .reduce((result, field) => ((result[field.oneofIndex] || (result[field.oneofIndex] = [])).push(field), result), {});
    // then add a case for each field
    messageDesc.field.forEach((field) => {
        const fieldName = (0, utils_1.getFieldName)(field, options);
        const messageProperty = (0, utils_1.getPropertyAccessor)("message", fieldName);
        // get a generic writer.doSomething based on the basic type
        const writeSnippet = getEncodeWriteSnippet(ctx, field);
        const isOptional = (0, types_1.isOptionalProperty)(field, messageDesc.options, options, currentFile.isProto3Syntax);
        if ((0, types_1.isRepeated)(field)) {
            if ((0, types_1.isMapType)(ctx, messageDesc, field)) {
                const valueType = typeMap.get(field.typeName)[2].field[1];
                const maybeTypeField = (0, options_1.addTypeToMessages)(options) ? `$type: '${field.typeName.slice(1)}',` : "";
                const entryWriteSnippet = (0, types_1.isValueType)(ctx, valueType)
                    ? (0, ts_poet_1.code) `
              if (value !== undefined ${(0, utils_1.withOrMaybeCheckIsNotNull)(options, `value`)}) {
                ${writeSnippet(`{ ${maybeTypeField} key: key as any, value }`)};
              }
            `
                    : writeSnippet(`{ ${maybeTypeField} key: key as any, value }`);
                const useMapType = (0, types_1.shouldGenerateJSMapType)(ctx, messageDesc, field);
                const optionalAlternative = isOptional ? (useMapType ? " || new Map()" : " || {}") : "";
                if (useMapType) {
                    chunks.push((0, ts_poet_1.code) `
            (${messageProperty}${optionalAlternative}).forEach((value, key) => {
              ${entryWriteSnippet}
            });
          `);
                }
                else {
                    chunks.push((0, ts_poet_1.code) `
            Object.entries(${messageProperty}${optionalAlternative}).forEach(([key, value]) => {
              ${entryWriteSnippet}
            });
          `);
                }
            }
            else if ((0, types_1.packedType)(field.type) === undefined) {
                const listWriteSnippet = (0, ts_poet_1.code) `
          for (const v of ${messageProperty}) {
            ${writeSnippet("v!")};
          }
        `;
                if (isOptional) {
                    chunks.push((0, ts_poet_1.code) `
            if (${messageProperty} !== undefined && ${messageProperty}.length !== 0) {
              ${listWriteSnippet}
            }
          `);
                }
                else {
                    chunks.push(listWriteSnippet);
                }
            }
            else if ((0, types_1.isEnum)(field) && options.stringEnums) {
                // This is a lot like the `else` clause, but we wrap `fooToNumber` around it.
                // Ideally we'd reuse `writeSnippet` here, but `writeSnippet` has the `writer.uint32(tag)`
                // embedded inside of it, and we want to drop that so that we can encode it packed
                // (i.e. just one tag and multiple values).
                const tag = ((field.number << 3) | 2) >>> 0;
                const toNumber = (0, types_1.getEnumMethod)(ctx, field.typeName, "ToNumber");
                const listWriteSnippet = (0, ts_poet_1.code) `
          writer.uint32(${tag}).fork();
          for (const v of ${messageProperty}) {
            writer.${(0, types_1.toReaderCall)(field)}(${toNumber}(v));
          }
          writer.ldelim();
        `;
                if (isOptional) {
                    chunks.push((0, ts_poet_1.code) `
            if (${messageProperty} !== undefined && ${messageProperty}.length !== 0) {
              ${listWriteSnippet}
            }
          `);
                }
                else {
                    chunks.push(listWriteSnippet);
                }
            }
            else {
                // Ideally we'd reuse `writeSnippet` but it has tagging embedded inside of it.
                const tag = ((field.number << 3) | 2) >>> 0;
                const rhs = (x) => ((0, types_1.isLong)(field) && options.forceLong === options_1.LongOption.BIGINT ? `${x}.toString()` : x);
                let listWriteSnippet = (0, ts_poet_1.code) `
          writer.uint32(${tag}).fork();
          for (const v of ${messageProperty}) {
            writer.${(0, types_1.toReaderCall)(field)}(${rhs("v")});
          }
          writer.ldelim();
        `;
                if ((0, types_1.isLong)(field) && options.forceLong === options_1.LongOption.BIGINT) {
                    const fieldType = (0, types_1.toReaderCall)(field);
                    switch (fieldType) {
                        case "int64":
                        case "sint64":
                        case "sfixed64":
                            listWriteSnippet = (0, ts_poet_1.code) `
                writer.uint32(${tag}).fork();
                for (const v of ${messageProperty}) {
                  if (BigInt.asIntN(64, v) !== v) {
                    throw new ${utils.globalThis}.Error('a value provided in array field ${fieldName} of type ${fieldType} is too large');
                  }
                  writer.${(0, types_1.toReaderCall)(field)}(${rhs("v")});
                }
                writer.ldelim();
              `;
                            break;
                        case "uint64":
                        case "fixed64":
                            listWriteSnippet = (0, ts_poet_1.code) `
                writer.uint32(${tag}).fork();
                for (const v of ${messageProperty}) {
                  if (BigInt.asUintN(64, v) !== v) {
                    throw new ${utils.globalThis}.Error('a value provided in array field ${fieldName} of type ${fieldType} is too large');
                  }
                  writer.${(0, types_1.toReaderCall)(field)}(${rhs("v")});
                }
                writer.ldelim();
              `;
                            break;
                        default:
                            throw new Error(`unexpected BigInt type: ${fieldType}`);
                    }
                }
                if (isOptional) {
                    chunks.push((0, ts_poet_1.code) `
            if (${messageProperty} !== undefined ${(0, utils_1.withAndMaybeCheckIsNotNull)(options, messageProperty)} && ${messageProperty}.length !== 0) {
              ${listWriteSnippet}
            }
          `);
                }
                else {
                    chunks.push(listWriteSnippet);
                }
            }
        }
        else if ((0, types_1.isWithinOneOfThatShouldBeUnion)(options, field)) {
            if (!processedOneofs.has(field.oneofIndex)) {
                processedOneofs.add(field.oneofIndex);
                const oneofNameWithMessage = options.useJsonName
                    ? messageProperty
                    : (0, utils_1.getPropertyAccessor)("message", (0, case_1.maybeSnakeToCamel)(messageDesc.oneofDecl[field.oneofIndex].name, options));
                chunks.push((0, ts_poet_1.code) `switch (${oneofNameWithMessage}?.$case) {`);
                for (const oneOfField of oneOfFieldsDict[field.oneofIndex]) {
                    const writeSnippet = getEncodeWriteSnippet(ctx, oneOfField);
                    const oneOfFieldName = (0, case_1.maybeSnakeToCamel)(oneOfField.name, ctx.options);
                    chunks.push((0, ts_poet_1.code) `case "${oneOfFieldName}":
            ${writeSnippet(`${oneofNameWithMessage}.${oneOfFieldName}`)};
            break;`);
                }
                chunks.push((0, ts_poet_1.code) `}`);
            }
        }
        else if ((0, types_1.isWithinOneOf)(field)) {
            // Oneofs don't have a default value check b/c they need to denote which-oneof presence
            chunks.push((0, ts_poet_1.code) `
        if (${messageProperty} !== undefined ${(0, utils_1.withAndMaybeCheckIsNotNull)(options, messageProperty)}) {
          ${writeSnippet(`${messageProperty}`)};
        }
      `);
        }
        else if ((0, types_1.isMessage)(field)) {
            chunks.push((0, ts_poet_1.code) `
        if (${messageProperty} !== undefined ${(0, utils_1.withAndMaybeCheckIsNotNull)(options, messageProperty)}) {
          ${writeSnippet(`${messageProperty}`)};
        }
      `);
        }
        else if ((0, types_1.isScalar)(field) || (0, types_1.isEnum)(field)) {
            const isJsType = (0, types_1.isScalar)(field) && (0, types_1.isJsTypeFieldOption)(options, field);
            const body = isJsType && options.forceLong === options_1.LongOption.BIGINT
                ? writeSnippet(`BigInt(${messageProperty})`)
                : writeSnippet(`${messageProperty}`);
            chunks.push((0, ts_poet_1.code) `
        if (${(0, types_1.notDefaultCheck)(ctx, field, messageDesc.options, `${messageProperty}`)}) {
          ${body};
        }
      `);
        }
        else {
            chunks.push((0, ts_poet_1.code) `${writeSnippet(`${messageProperty}`)};`);
        }
    });
    if (options.unknownFields) {
        chunks.push((0, ts_poet_1.code) `if (message._unknownFields !== undefined) {
      for (const [key, values] of Object.entries(message._unknownFields)) {
        const tag = parseInt(key, 10);
        for (const value of values) {
          writer.uint32(tag);
          (writer as any)['_push'](
            (val: Uint8Array, buf: Buffer, pos: number) => buf.set(val, pos),
            value.length,
            value
          );
        }
      }
    }`);
    }
    chunks.push((0, ts_poet_1.code) `return writer;`);
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
function generateSetExtension(ctx, fullName) {
    return (0, ts_poet_1.code) `
    setExtension <T> (message: ${fullName}, extension: ${ctx.utils.Extension}<T>, value: T): void {
      const encoded = extension.encode!(value);

      if (message._unknownFields !== undefined) {
        delete message._unknownFields[extension.tag];

        if (extension.singularTag !== undefined) {
          delete message._unknownFields[extension.singularTag];
        }
      }

      if (encoded.length !== 0) {
        if (message._unknownFields === undefined) {
          message._unknownFields = {};
        }

        message._unknownFields[extension.tag] = encoded;
      }
    }
  `;
}
function generateGetExtension(ctx, fullName) {
    return (0, ts_poet_1.code) `
    getExtension <T> (message: ${fullName}, extension: ${ctx.utils.Extension}<T>): T | undefined {
      let results: T | undefined = undefined;

      if (message._unknownFields === undefined) {
        return undefined;
      }

      let list = message._unknownFields[extension.tag];

      if (list !== undefined) {
        results = extension.decode!(extension.tag, list);
      }

      if (extension.singularTag === undefined) {
        return results;
      }

      list = message._unknownFields[extension.singularTag];

      if (list !== undefined) {
        const results2 = extension.decode!(extension.singularTag, list);

        if (results !== undefined && (results as any).length !== 0) {
          results = (results as any).concat(results2);
        } else {
          results = results2;
        }
      }

      return results;
    }
  `;
}
function generateExtension(ctx, message, extension) {
    var _a;
    const type = (0, types_1.toTypeName)(ctx, message, extension);
    const packedTag = (0, types_1.isRepeated)(extension) && (0, types_1.packedType)(extension.type) !== undefined ? ((extension.number << 3) | 2) >>> 0 : undefined;
    const singularTag = ((extension.number << 3) | (0, types_1.basicWireType)(extension.type)) >>> 0;
    const tag = packedTag !== null && packedTag !== void 0 ? packedTag : singularTag;
    const chunks = [];
    chunks.push((0, ts_poet_1.code) `{`);
    chunks.push((0, ts_poet_1.code) `number: ${extension.number},`);
    chunks.push((0, ts_poet_1.code) `tag: ${tag},`);
    if (packedTag !== undefined)
        chunks.push((0, ts_poet_1.code) `singularTag: ${singularTag},`);
    chunks.push((0, ts_poet_1.code) `repeated: ${extension.label == ts_proto_descriptors_1.FieldDescriptorProto_Label.LABEL_REPEATED},`);
    chunks.push((0, ts_poet_1.code) `packed: ${((_a = extension.options) === null || _a === void 0 ? void 0 : _a.packed) ? true : false},`);
    const Reader = (0, utils_1.impFile)(ctx.options, "Reader@protobufjs/minimal");
    const Writer = (0, utils_1.impFile)(ctx.options, "Writer@protobufjs/minimal");
    if (ctx.options.outputEncodeMethods === true ||
        ctx.options.outputEncodeMethods === "encode-only" ||
        ctx.options.outputEncodeMethods === "encode-no-creation") {
        chunks.push((0, ts_poet_1.code) `
      encode: (value: ${type}): Uint8Array[] => {
        const encoded: Uint8Array[] = [];
    `);
        function getEncodeSnippet(ctx, field) {
            const { options, utils } = ctx;
            if ((0, types_1.isEnum)(field) && options.stringEnums) {
                const toNumber = (0, types_1.getEnumMethod)(ctx, field.typeName, "ToNumber");
                return (place) => (0, ts_poet_1.code) `writer.${(0, types_1.toReaderCall)(field)}(${toNumber}(${place}))`;
            }
            else if ((0, types_1.isLong)(field) && options.forceLong === options_1.LongOption.BIGINT) {
                return (place) => (0, ts_poet_1.code) `writer.${(0, types_1.toReaderCall)(field)}(${place}.toString())`;
            }
            else if ((0, types_1.isScalar)(field) || (0, types_1.isEnum)(field)) {
                return (place) => (0, ts_poet_1.code) `writer.${(0, types_1.toReaderCall)(field)}(${place})`;
            }
            else if ((0, types_1.isObjectId)(field) && options.useMongoObjectId) {
                const type = (0, types_1.basicTypeName)(ctx, field, { keepValueType: true });
                return (place) => (0, ts_poet_1.code) `${type}.encode(${utils.toProtoObjectId}(${place}), writer.fork()).ldelim()`;
            }
            else if ((0, types_1.isTimestamp)(field) &&
                (options.useDate === options_1.DateOption.DATE ||
                    options.useDate === options_1.DateOption.STRING ||
                    options.useDate === options_1.DateOption.STRING_NANO)) {
                const type = (0, types_1.basicTypeName)(ctx, field, { keepValueType: true });
                return (place) => (0, ts_poet_1.code) `${type}.encode(${utils.toTimestamp}(${place}), writer.fork()).ldelim()`;
            }
            else if ((0, types_1.isValueType)(ctx, field)) {
                const maybeTypeField = (0, options_1.addTypeToMessages)(options) ? `$type: '${field.typeName.slice(1)}',` : "";
                const type = (0, types_1.basicTypeName)(ctx, field, { keepValueType: true });
                const wrappedValue = (place) => {
                    if ((0, types_1.isAnyValueType)(field) || (0, types_1.isListValueType)(field) || (0, types_1.isStructType)(field) || (0, types_1.isFieldMaskType)(field)) {
                        return (0, ts_poet_1.code) `${type}.wrap(${place})`;
                    }
                    return (0, ts_poet_1.code) `{${maybeTypeField} value: ${place}!}`;
                };
                return (place) => (0, ts_poet_1.code) `${type}.encode(${wrappedValue(place)}, writer.fork()).ldelim()`;
            }
            else if ((0, types_1.isMessage)(field)) {
                const type = (0, types_1.basicTypeName)(ctx, field);
                if (field.type == ts_proto_descriptors_1.FieldDescriptorProto_Type.TYPE_GROUP) {
                    const endTag = ((field.number << 3) | 4) >>> 0;
                    return (place) => (0, ts_poet_1.code) `${type}.encode(${place}, writer).uint32(${endTag})`;
                }
                return (place) => (0, ts_poet_1.code) `${type}.encode(${place}, writer.fork()).ldelim()`;
            }
            else {
                throw new Error(`Unhandled field ${field}`);
            }
        }
        const writeSnippet = getEncodeSnippet(ctx, extension);
        if ((0, types_1.isRepeated)(extension)) {
            if (packedTag === undefined) {
                chunks.push((0, ts_poet_1.code) `
          for (const v of value) {
            const writer = ${Writer}.create();
            ${writeSnippet("v")};
            encoded.push(writer.finish());
          }
        `);
            }
            else {
                const rhs = (x) => (0, types_1.isLong)(extension) && ctx.options.forceLong === options_1.LongOption.BIGINT ? `${x}.toString()` : x;
                chunks.push((0, ts_poet_1.code) `
          const writer = ${Writer}.create();
          writer.fork();
          for (const v of value) {
            ${writeSnippet(rhs("v"))};
          }
          writer.ldelim();
          encoded.push(writer.finish());
        `);
            }
        }
        else if ((0, types_1.isScalar)(extension) || (0, types_1.isEnum)(extension)) {
            chunks.push((0, ts_poet_1.code) `
        if (${(0, types_1.notDefaultCheck)(ctx, extension, message === null || message === void 0 ? void 0 : message.options, "value")}) {
          const writer = ${Writer}.create();
          ${writeSnippet("value")};
          encoded.push(writer.finish());
        }
      `);
        }
        else {
            chunks.push((0, ts_poet_1.code) `
        const writer = ${Writer}.create();
        ${writeSnippet("value")};
        encoded.push(writer.finish());
      `);
        }
        chunks.push((0, ts_poet_1.code) `
        return encoded;
      },
    `);
    }
    if (ctx.options.outputEncodeMethods === true || ctx.options.outputEncodeMethods === "decode-only") {
        chunks.push((0, ts_poet_1.code) `decode: (tag: number, input: Uint8Array[]): ${type} => {`);
        // get a generic 'reader.doSomething' bit that is specific to the basic type
        const readSnippet = getDecodeReadSnippet(ctx, extension);
        if ((0, types_1.isRepeated)(extension)) {
            chunks.push((0, ts_poet_1.code) `const values: ${type} = [];`);
            // start loop over all buffers
            chunks.push((0, ts_poet_1.code) `
        for (const buffer of input) {
          const reader = ${Reader}.create(buffer);
      `);
            if ((0, types_1.packedType)(extension.type) === undefined) {
                chunks.push((0, ts_poet_1.code) `
          values.push(${readSnippet});
        `);
            }
            else {
                chunks.push((0, ts_poet_1.code) `
          if (tag == ${packedTag}) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              values.push(${readSnippet});
            }
          } else {
            values.push(${readSnippet});
          }
        `);
            }
            chunks.push((0, ts_poet_1.code) `
          }

          return values;
        },
      `);
        }
        else {
            // pick the last entry, since it overrides all previous entries if not repeated
            chunks.push((0, ts_poet_1.code) `
          const reader = ${Reader}.create(input[input.length -1] ?? ${ctx.utils.fail}());
          return ${readSnippet};
        },
      `);
        }
    }
    chunks.push((0, ts_poet_1.code) `}`);
    return {
        name: (0, case_1.maybeSnakeToCamel)(extension.name, ctx.options),
        type,
        extensionInfo: (0, ts_poet_1.joinCode)(chunks, { on: "\n" }),
    };
}
/**
 * Creates a function to decode a message from JSON.
 *
 * This is very similar to decode, we loop through looking for properties, with
 * a few special cases for https://developers.google.com/protocol-buffers/docs/proto3#json.
 * */
function generateFromJson(ctx, fullName, fullTypeName, messageDesc) {
    const { options, utils, currentFile } = ctx;
    const chunks = [];
    // create the basic function declaration
    chunks.push((0, ts_poet_1.code) `
    fromJSON(${messageDesc.field.length > 0 ? "object" : "_"}: any): ${fullName} {
      return {
  `);
    if ((0, options_1.addTypeToMessages)(options)) {
        chunks.push((0, ts_poet_1.code) `$type: ${fullName}.$type,`);
    }
    const oneofFieldsCases = messageDesc.oneofDecl.map((oneof, oneofIndex) => messageDesc.field.filter(types_1.isWithinOneOf).filter((field) => field.oneofIndex === oneofIndex));
    const canonicalFromJson = {
        ["google.protobuf.FieldMask"]: {
            paths: (from) => (0, ts_poet_1.code) `typeof(${from}) === 'string'
        ? ${from}.split(",").filter(${ctx.utils.globalThis}.Boolean)
        : ${ctx.utils.globalThis}.Array.isArray(${from}?.paths)
        ? ${from}.paths.map(${ctx.utils.globalThis}.String)
        : []`,
        },
    };
    // add a check for each incoming field
    messageDesc.field.forEach((field) => {
        var _a;
        const fieldName = (0, utils_1.getFieldName)(field, options);
        const fieldKey = (0, utils_1.safeAccessor)(fieldName);
        const jsonName = (0, utils_1.getFieldJsonName)(field, options);
        const jsonProperty = (0, utils_1.getPropertyAccessor)("object", jsonName);
        const jsonPropertyOptional = (0, utils_1.getPropertyAccessor)("object", jsonName, true);
        // get code that extracts value from incoming object
        const readSnippet = (from) => {
            var _a;
            if ((0, types_1.isEnum)(field)) {
                const fromJson = (0, types_1.getEnumMethod)(ctx, field.typeName, "FromJSON");
                return (0, ts_poet_1.code) `${fromJson}(${from})`;
            }
            else if ((0, types_1.isPrimitive)(field)) {
                // Convert primitives using the String(value)/Number(value)/bytesFromBase64(value)
                if ((0, types_1.isBytes)(field)) {
                    if (options.env === options_1.EnvOption.NODE) {
                        return (0, ts_poet_1.code) `Buffer.from(${utils.bytesFromBase64}(${from}))`;
                    }
                    else {
                        return (0, ts_poet_1.code) `${utils.bytesFromBase64}(${from})`;
                    }
                }
                else if ((0, types_1.isLong)(field) && (0, types_1.isJsTypeFieldOption)(options, field)) {
                    const fieldType = (_a = (0, types_1.getFieldOptionsJsType)(field, ctx.options)) !== null && _a !== void 0 ? _a : field.type;
                    const cstr = (0, case_1.capitalize)((0, types_1.basicTypeName)(ctx, { ...field, type: fieldType }, { keepValueType: true }).toCodeString([]));
                    return (0, ts_poet_1.code) `${utils.globalThis}.${cstr}(${from})`;
                }
                else if ((0, types_1.isLong)(field) && options.forceLong === options_1.LongOption.LONG) {
                    const cstr = (0, case_1.capitalize)((0, types_1.basicTypeName)(ctx, field, { keepValueType: true }).toCodeString([]));
                    return (0, ts_poet_1.code) `${cstr}.fromValue(${from})`;
                }
                else if ((0, types_1.isLong)(field) && options.forceLong === options_1.LongOption.BIGINT) {
                    return (0, ts_poet_1.code) `BigInt(${from})`;
                }
                else {
                    const cstr = (0, case_1.capitalize)((0, types_1.basicTypeName)(ctx, field, { keepValueType: true }).toCodeString([]));
                    return (0, ts_poet_1.code) `${utils.globalThis}.${cstr}(${from})`;
                }
            }
            else if ((0, types_1.isObjectId)(field) && options.useMongoObjectId) {
                return (0, ts_poet_1.code) `${utils.fromJsonObjectId}(${from})`;
            }
            else if ((0, types_1.isTimestamp)(field) &&
                (options.useDate === options_1.DateOption.STRING || options.useDate === options_1.DateOption.STRING_NANO)) {
                return (0, ts_poet_1.code) `${utils.globalThis}.String(${from})`;
            }
            else if ((0, types_1.isTimestamp)(field) &&
                (options.useDate === options_1.DateOption.DATE || options.useDate === options_1.DateOption.TIMESTAMP)) {
                return (0, ts_poet_1.code) `${utils.fromJsonTimestamp}(${from})`;
            }
            else if ((0, types_1.isAnyValueType)(field) || (0, types_1.isStructType)(field)) {
                return (0, ts_poet_1.code) `${from}`;
            }
            else if ((0, types_1.isFieldMaskType)(field)) {
                const type = (0, types_1.basicTypeName)(ctx, field, { keepValueType: true });
                return (0, ts_poet_1.code) `${type}.unwrap(${type}.fromJSON(${from}))`;
            }
            else if ((0, types_1.isListValueType)(field)) {
                return (0, ts_poet_1.code) `[...${from}]`;
            }
            else if ((0, types_1.isValueType)(ctx, field)) {
                const valueType = (0, types_1.valueTypeName)(ctx, field.typeName);
                if ((0, types_1.isLongValueType)(field) && options.forceLong === options_1.LongOption.LONG) {
                    return (0, ts_poet_1.code) `${(0, case_1.capitalize)(valueType.toCodeString([]))}.fromValue(${from})`;
                }
                else if ((0, types_1.isLongValueType)(field) && options.forceLong === options_1.LongOption.BIGINT) {
                    return (0, ts_poet_1.code) `BigInt(${from})`;
                }
                else if ((0, types_1.isBytesValueType)(field)) {
                    return (0, ts_poet_1.code) `new ${(0, case_1.capitalize)(valueType.toCodeString([]))}(${from})`;
                }
                else {
                    return (0, ts_poet_1.code) `${(0, case_1.capitalize)(valueType.toCodeString([]))}(${from})`;
                }
            }
            else if ((0, types_1.isMessage)(field)) {
                if ((0, types_1.isRepeated)(field) && (0, types_1.isMapType)(ctx, messageDesc, field)) {
                    const { valueField, valueType } = (0, types_1.detectMapType)(ctx, messageDesc, field);
                    if ((0, types_1.isPrimitive)(valueField)) {
                        // TODO Can we not copy/paste this from ^?
                        if ((0, types_1.isBytes)(valueField)) {
                            if (options.env === options_1.EnvOption.NODE) {
                                return (0, ts_poet_1.code) `Buffer.from(${utils.bytesFromBase64}(${from} as string))`;
                            }
                            else {
                                return (0, ts_poet_1.code) `${utils.bytesFromBase64}(${from} as string)`;
                            }
                        }
                        else if ((0, types_1.isLong)(valueField) && options.forceLong === options_1.LongOption.LONG) {
                            return (0, ts_poet_1.code) `Long.fromValue(${from} as Long | string)`;
                        }
                        else if ((0, types_1.isLong)(valueField) && options.forceLong === options_1.LongOption.BIGINT) {
                            return (0, ts_poet_1.code) `BigInt(${from} as string | number | bigint | boolean)`;
                        }
                        else if ((0, types_1.isEnum)(valueField)) {
                            const fromJson = (0, types_1.getEnumMethod)(ctx, valueField.typeName, "FromJSON");
                            return (0, ts_poet_1.code) `${fromJson}(${from})`;
                        }
                        else {
                            const cstr = (0, case_1.capitalize)(valueType.toCodeString([]));
                            return (0, ts_poet_1.code) `${cstr}(${from})`;
                        }
                    }
                    else if ((0, types_1.isObjectId)(valueField) && options.useMongoObjectId) {
                        return (0, ts_poet_1.code) `${utils.fromJsonObjectId}(${from})`;
                    }
                    else if ((0, types_1.isTimestamp)(valueField) &&
                        (options.useDate === options_1.DateOption.STRING || options.useDate === options_1.DateOption.STRING_NANO)) {
                        return (0, ts_poet_1.code) `${utils.globalThis}.String(${from})`;
                    }
                    else if ((0, types_1.isTimestamp)(valueField) &&
                        (options.useDate === options_1.DateOption.DATE || options.useDate === options_1.DateOption.TIMESTAMP)) {
                        return (0, ts_poet_1.code) `${utils.fromJsonTimestamp}(${from})`;
                    }
                    else if ((0, types_1.isValueType)(ctx, valueField)) {
                        return (0, ts_poet_1.code) `${from} as ${valueType}`;
                    }
                    else if ((0, types_1.isAnyValueType)(valueField)) {
                        return (0, ts_poet_1.code) `${from}`;
                    }
                    else {
                        return (0, ts_poet_1.code) `${valueType}.fromJSON(${from})`;
                    }
                }
                else {
                    const type = (0, types_1.basicTypeName)(ctx, field);
                    return (0, ts_poet_1.code) `${type}.fromJSON(${from})`;
                }
            }
            else {
                throw new Error(`Unhandled field ${field}`);
            }
        };
        const noDefaultValue = !options.initializeFieldsAsUndefined &&
            (0, types_1.isOptionalProperty)(field, messageDesc.options, options, currentFile.isProto3Syntax);
        // and then use the snippet to handle repeated fields if necessary
        if ((_a = canonicalFromJson[fullTypeName]) === null || _a === void 0 ? void 0 : _a[fieldName]) {
            chunks.push((0, ts_poet_1.code) `${fieldName}: ${canonicalFromJson[fullTypeName][fieldName]("object")},`);
        }
        else if ((0, types_1.isRepeated)(field)) {
            if ((0, types_1.isMapType)(ctx, messageDesc, field)) {
                const fieldType = (0, types_1.toTypeName)(ctx, messageDesc, field);
                const i = convertFromObjectKey(ctx, messageDesc, field, "key");
                if ((0, types_1.shouldGenerateJSMapType)(ctx, messageDesc, field)) {
                    const fallback = noDefaultValue ? (0, utils_1.nullOrUndefined)(options) : "new Map()";
                    chunks.push((0, ts_poet_1.code) `
            ${fieldKey}: ${ctx.utils.isObject}(${jsonProperty})
              ? Object.entries(${jsonProperty}).reduce<${fieldType}>((acc, [key, value]) => {
                  acc.set(${i}, ${readSnippet("value")});
                  return acc;
                }, new Map())
              : ${fallback},
          `);
                }
                else {
                    const fallback = noDefaultValue ? (0, utils_1.nullOrUndefined)(options) : "{}";
                    chunks.push((0, ts_poet_1.code) `
            ${fieldKey}: ${ctx.utils.isObject}(${jsonProperty})
              ? Object.entries(${jsonProperty}).reduce<${fieldType}>((acc, [key, value]) => {
                  acc[${i}] = ${readSnippet("value")};
                  return acc;
                }, {})
              : ${fallback},
          `);
                }
            }
            else {
                const fallback = noDefaultValue ? (0, utils_1.nullOrUndefined)(options) : "[]";
                const readValueSnippet = readSnippet("e");
                if (readValueSnippet.toString() === (0, ts_poet_1.code) `e`.toString()) {
                    chunks.push((0, ts_poet_1.code) `${fieldKey}: ${ctx.utils.globalThis}.Array.isArray(${jsonPropertyOptional}) ? [...${jsonProperty}] : [],`);
                }
                else {
                    // Explicit `any` type required to make TS with noImplicitAny happy. `object` is also `any` here.
                    chunks.push((0, ts_poet_1.code) `
            ${fieldKey}: ${ctx.utils.globalThis}.Array.isArray(${jsonPropertyOptional}) ? ${jsonProperty}.map((e: any) => ${readValueSnippet}): ${fallback},
          `);
                }
            }
        }
        else if ((0, types_1.isWithinOneOfThatShouldBeUnion)(options, field)) {
            const cases = oneofFieldsCases[field.oneofIndex];
            const firstCase = cases[0];
            const lastCase = cases[cases.length - 1];
            if (field === firstCase) {
                const fieldName = (0, case_1.maybeSnakeToCamel)(messageDesc.oneofDecl[field.oneofIndex].name, options);
                chunks.push((0, ts_poet_1.code) `${fieldName}: `);
            }
            const ternaryIf = (0, ts_poet_1.code) `${ctx.utils.isSet}(${jsonProperty})`;
            const ternaryThen = (0, ts_poet_1.code) `{ $case: '${fieldName}', ${fieldKey}: ${readSnippet(`${jsonProperty}`)}`;
            chunks.push((0, ts_poet_1.code) `${ternaryIf} ? ${ternaryThen}} : `);
            if (field === lastCase) {
                chunks.push((0, ts_poet_1.code) `${(0, utils_1.nullOrUndefined)(options)},`);
            }
        }
        else if ((0, types_1.isAnyValueType)(field)) {
            chunks.push((0, ts_poet_1.code) `${fieldKey}: ${ctx.utils.isSet}(${jsonPropertyOptional})
        ? ${readSnippet(`${jsonProperty}`)}
        : ${(0, utils_1.nullOrUndefined)(options)},
      `);
        }
        else if ((0, types_1.isStructType)(field)) {
            chunks.push((0, ts_poet_1.code) `${fieldKey}: ${ctx.utils.isObject}(${jsonProperty})
          ? ${readSnippet(`${jsonProperty}`)}
          : ${(0, utils_1.nullOrUndefined)(options)},`);
        }
        else if ((0, types_1.isListValueType)(field)) {
            chunks.push((0, ts_poet_1.code) `
        ${fieldKey}: ${ctx.utils.globalThis}.Array.isArray(${jsonProperty})
          ? ${readSnippet(`${jsonProperty}`)}
          : ${(0, utils_1.nullOrUndefined)(options)},
      `);
        }
        else {
            const fallback = (0, types_1.isWithinOneOf)(field) || noDefaultValue ? (0, utils_1.nullOrUndefined)(options) : (0, types_1.defaultValue)(ctx, field);
            chunks.push((0, ts_poet_1.code) `
        ${fieldKey}: ${ctx.utils.isSet}(${jsonProperty})
          ? ${readSnippet(`${jsonProperty}`)}
          : ${fallback},
      `);
        }
    });
    // and then wrap up the switch/while/return
    chunks.push((0, ts_poet_1.code) `};`);
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
function generateCanonicalToJson(fullName, fullProtobufTypeName, { useOptionals, useNullAsOptional }) {
    if ((0, types_1.isFieldMaskTypeName)(fullProtobufTypeName)) {
        const returnType = useOptionals === "all" ? `string | ${(0, utils_1.nullOrUndefined)({ useNullAsOptional })}` : "string";
        const pathModifier = useOptionals === "all" ? "?" : "";
        return (0, ts_poet_1.code) `
    toJSON(message: ${fullName}): ${returnType} {
      return message.paths${pathModifier}.join(',');
    }
  `;
    }
    return undefined;
}
function generateToJson(ctx, fullName, fullProtobufTypeName, messageDesc) {
    const { options, utils, typeMap } = ctx;
    const chunks = [];
    const canonicalToJson = generateCanonicalToJson(fullName, fullProtobufTypeName, options);
    if (canonicalToJson) {
        chunks.push(canonicalToJson);
        return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
    }
    // create the basic function declaration
    chunks.push((0, ts_poet_1.code) `
    toJSON(${messageDesc.field.length > 0 ? "message" : "_"}: ${fullName}): unknown {
      const obj: any = {};
  `);
    // then add a case for each field
    messageDesc.field.forEach((field) => {
        const fieldName = (0, utils_1.getFieldName)(field, options);
        const jsonName = (0, utils_1.getFieldJsonName)(field, options);
        const jsonProperty = (0, utils_1.getPropertyAccessor)("obj", jsonName);
        const messageProperty = (0, utils_1.getPropertyAccessor)("message", fieldName);
        const readSnippet = (from) => {
            var _a;
            if ((0, types_1.isEnum)(field)) {
                const toJson = (0, types_1.getEnumMethod)(ctx, field.typeName, "ToJSON");
                return (0, ts_poet_1.code) `${toJson}(${from})`;
            }
            else if ((0, types_1.isObjectId)(field) && options.useMongoObjectId) {
                return (0, ts_poet_1.code) `${from}.toString()`;
            }
            else if ((0, types_1.isTimestamp)(field) && options.useDate === options_1.DateOption.DATE) {
                return (0, ts_poet_1.code) `${from}.toISOString()`;
            }
            else if ((0, types_1.isTimestamp)(field) &&
                (options.useDate === options_1.DateOption.STRING || options.useDate === options_1.DateOption.STRING_NANO)) {
                return (0, ts_poet_1.code) `${from}`;
            }
            else if ((0, types_1.isTimestamp)(field) && options.useDate === options_1.DateOption.TIMESTAMP) {
                if (options.useJsonTimestamp === options_1.JsonTimestampOption.RAW) {
                    return (0, ts_poet_1.code) `${from}`;
                }
                return (0, ts_poet_1.code) `${utils.fromTimestamp}(${from}).toISOString()`;
            }
            else if ((0, types_1.isMapType)(ctx, messageDesc, field)) {
                // For map types, drill-in and then admittedly re-hard-code our per-value-type logic
                const valueType = typeMap.get(field.typeName)[2].field[1];
                if ((0, types_1.isEnum)(valueType)) {
                    const toJson = (0, types_1.getEnumMethod)(ctx, valueType.typeName, "ToJSON");
                    return (0, ts_poet_1.code) `${toJson}(${from})`;
                }
                else if ((0, types_1.isBytes)(valueType)) {
                    return (0, ts_poet_1.code) `${utils.base64FromBytes}(${from})`;
                }
                else if ((0, types_1.isObjectId)(valueType) && options.useMongoObjectId) {
                    return (0, ts_poet_1.code) `${from}.toString()`;
                }
                else if ((0, types_1.isTimestamp)(valueType) && options.useDate === options_1.DateOption.DATE) {
                    return (0, ts_poet_1.code) `${from}.toISOString()`;
                }
                else if ((0, types_1.isTimestamp)(valueType) &&
                    (options.useDate === options_1.DateOption.STRING || options.useDate === options_1.DateOption.STRING_NANO)) {
                    return (0, ts_poet_1.code) `${from}`;
                }
                else if ((0, types_1.isTimestamp)(valueType) && options.useDate === options_1.DateOption.TIMESTAMP) {
                    return (0, ts_poet_1.code) `${utils.fromTimestamp}(${from}).toISOString()`;
                }
                else if ((0, types_1.isLong)(valueType) && options.forceLong === options_1.LongOption.LONG) {
                    return (0, ts_poet_1.code) `${from}.toString()`;
                }
                else if ((0, types_1.isLong)(valueType) && options.forceLong === options_1.LongOption.BIGINT) {
                    return (0, ts_poet_1.code) `${from}.toString()`;
                }
                else if ((0, types_1.isWholeNumber)(valueType) && !((0, types_1.isLong)(valueType) && options.forceLong === options_1.LongOption.STRING)) {
                    return (0, ts_poet_1.code) `Math.round(${from})`;
                }
                else if ((0, types_1.isScalar)(valueType) || (0, types_1.isValueType)(ctx, valueType)) {
                    return (0, ts_poet_1.code) `${from}`;
                }
                else if ((0, types_1.isAnyValueType)(valueType)) {
                    return (0, ts_poet_1.code) `${from}`;
                }
                else {
                    const type = (0, types_1.basicTypeName)(ctx, valueType);
                    return (0, ts_poet_1.code) `${type}.toJSON(${from})`;
                }
            }
            else if ((0, types_1.isAnyValueType)(field)) {
                return (0, ts_poet_1.code) `${from}`;
            }
            else if ((0, types_1.isFieldMaskType)(field)) {
                const type = (0, types_1.basicTypeName)(ctx, field, { keepValueType: true });
                return (0, ts_poet_1.code) `${type}.toJSON(${type}.wrap(${from}))`;
            }
            else if ((0, types_1.isMessage)(field) && !(0, types_1.isValueType)(ctx, field) && !(0, types_1.isMapType)(ctx, messageDesc, field)) {
                const type = (0, types_1.basicTypeName)(ctx, field, { keepValueType: true });
                return (0, ts_poet_1.code) `${type}.toJSON(${from})`;
            }
            else if ((0, types_1.isBytes)(field)) {
                return (0, ts_poet_1.code) `${utils.base64FromBytes}(${from})`;
            }
            else if ((0, types_1.isLong)(field) && (0, types_1.isJsTypeFieldOption)(options, field)) {
                const fieldType = (_a = (0, types_1.getFieldOptionsJsType)(field, ctx.options)) !== null && _a !== void 0 ? _a : field.type;
                if (!fieldType) {
                    return (0, ts_poet_1.code) `${from}`;
                }
                const cstr = (0, case_1.capitalize)((0, types_1.basicTypeName)(ctx, { ...field, type: fieldType }, { keepValueType: true }).toCodeString([]));
                return (0, ts_poet_1.code) `${utils.globalThis}.${cstr}(${from})`;
            }
            else if ((0, types_1.isLong)(field) && options.forceLong === options_1.LongOption.LONG) {
                return (0, ts_poet_1.code) `(${from} || ${(0, types_1.defaultValue)(ctx, field)}).toString()`;
            }
            else if ((0, types_1.isLong)(field) && options.forceLong === options_1.LongOption.BIGINT) {
                return (0, ts_poet_1.code) `${from}.toString()`;
            }
            else if ((0, types_1.isWholeNumber)(field) && !((0, types_1.isLong)(field) && options.forceLong === options_1.LongOption.STRING)) {
                return (0, ts_poet_1.code) `Math.round(${from})`;
            }
            else {
                return (0, ts_poet_1.code) `${from}`;
            }
        };
        if ((0, types_1.isMapType)(ctx, messageDesc, field)) {
            // Maps might need their values transformed, i.e. bytes --> base64
            const i = convertToObjectKey(ctx, messageDesc, field, "k");
            if ((0, types_1.shouldGenerateJSMapType)(ctx, messageDesc, field)) {
                chunks.push((0, ts_poet_1.code) `
          if (${messageProperty}?.size) {
            ${jsonProperty} = {};
            ${messageProperty}.forEach((v, k) => {
              ${jsonProperty}[${i}] = ${readSnippet("v")};
            });
          }
        `);
            }
            else {
                chunks.push((0, ts_poet_1.code) `
        if (${messageProperty}) {
            const entries = Object.entries(${messageProperty});
            if (entries.length > 0) {
              ${jsonProperty} = {};
              entries.forEach(([k, v]) => {
                ${jsonProperty}[${i}] = ${readSnippet("v")};
              });
            }
          }
        `);
            }
        }
        else if ((0, types_1.isRepeated)(field)) {
            // Arrays might need their elements transformed
            const transformElement = readSnippet("e");
            const maybeMap = transformElement.toCodeString([]) !== "e" ? (0, ts_poet_1.code) `.map(e => ${transformElement})` : "";
            chunks.push((0, ts_poet_1.code) `
        if (${messageProperty}?.length) {
          ${jsonProperty} = ${messageProperty}${maybeMap};
        }
      `);
        }
        else if ((0, types_1.isWithinOneOfThatShouldBeUnion)(options, field)) {
            // oneofs in a union are only output as `oneof name = ...`
            const oneofNameWithMessage = options.useJsonName
                ? messageProperty
                : (0, utils_1.getPropertyAccessor)("message", (0, case_1.maybeSnakeToCamel)(messageDesc.oneofDecl[field.oneofIndex].name, options));
            chunks.push((0, ts_poet_1.code) `
        if (${oneofNameWithMessage}?.$case === '${fieldName}') {
          ${jsonProperty} = ${readSnippet(`${oneofNameWithMessage}.${fieldName}`)};
        }
      `);
        }
        else {
            let emitDefaultValuesForJson = ctx.options.emitDefaultValues.includes("json-methods");
            const check = ((0, types_1.isScalar)(field) || (0, types_1.isEnum)(field)) && !((0, types_1.isWithinOneOf)(field) || emitDefaultValuesForJson)
                ? (0, types_1.notDefaultCheck)(ctx, field, messageDesc.options, `${messageProperty}`)
                : `${messageProperty} !== undefined ${(0, utils_1.withAndMaybeCheckIsNotNull)(options, messageProperty)}`;
            chunks.push((0, ts_poet_1.code) `
        if (${check}) {
          ${jsonProperty} = ${readSnippet(`${messageProperty}`)};
        }
      `);
        }
    });
    chunks.push((0, ts_poet_1.code) `return obj;`);
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
function generateFromPartial(ctx, fullName, messageDesc) {
    const { options, utils } = ctx;
    const chunks = [];
    // create the create function definition
    if (ctx.options.useExactTypes) {
        chunks.push((0, ts_poet_1.code) `
      create<I extends ${utils.Exact}<${utils.DeepPartial}<${fullName}>, I>>(base?: I): ${fullName} {
        return ${fullName}.fromPartial(base ?? ({} as any));
      },
    `);
    }
    else {
        chunks.push((0, ts_poet_1.code) `
      create(base?: ${utils.DeepPartial}<${fullName}>): ${fullName} {
        return ${fullName}.fromPartial(base ?? {});
      },
    `);
    }
    // create the fromPartial function declaration
    const paramName = messageDesc.field.length > 0 ? "object" : "_";
    if (ctx.options.useExactTypes) {
        chunks.push((0, ts_poet_1.code) `
      fromPartial<I extends ${utils.Exact}<${utils.DeepPartial}<${fullName}>, I>>(${paramName}: I): ${fullName} {
    `);
    }
    else {
        chunks.push((0, ts_poet_1.code) `
      fromPartial(${paramName}: ${utils.DeepPartial}<${fullName}>): ${fullName} {
    `);
    }
    let createBase = (0, ts_poet_1.code) `createBase${fullName}()`;
    if (options.usePrototypeForDefaults) {
        createBase = (0, ts_poet_1.code) `Object.create(${createBase}) as ${fullName}`;
    }
    chunks.push((0, ts_poet_1.code) `const message = ${createBase}${maybeAsAny(options)};`);
    // add a check for each incoming field
    messageDesc.field.forEach((field) => {
        const fieldName = (0, utils_1.getFieldName)(field, options);
        const messageProperty = (0, utils_1.getPropertyAccessor)("message", fieldName);
        const objectProperty = (0, utils_1.getPropertyAccessor)("object", fieldName);
        const readSnippet = (from) => {
            if (((0, types_1.isLong)(field) || (0, types_1.isLongValueType)(field)) &&
                options.forceLong === options_1.LongOption.LONG &&
                !(0, types_1.isJsTypeFieldOption)(options, field)) {
                return (0, ts_poet_1.code) `Long.fromValue(${from})`;
            }
            else if ((0, types_1.isObjectId)(field) && options.useMongoObjectId) {
                return (0, ts_poet_1.code) `${from} as mongodb.ObjectId`;
            }
            else if ((0, types_1.isPrimitive)(field) ||
                ((0, types_1.isTimestamp)(field) &&
                    (options.useDate === options_1.DateOption.DATE ||
                        options.useDate === options_1.DateOption.STRING ||
                        options.useDate === options_1.DateOption.STRING_NANO)) ||
                (0, types_1.isValueType)(ctx, field)) {
                return (0, ts_poet_1.code) `${from}`;
            }
            else if ((0, types_1.isMessage)(field)) {
                if ((0, types_1.isRepeated)(field) && (0, types_1.isMapType)(ctx, messageDesc, field)) {
                    const { valueField, valueType } = (0, types_1.detectMapType)(ctx, messageDesc, field);
                    if ((0, types_1.isPrimitive)(valueField)) {
                        if ((0, types_1.isBytes)(valueField)) {
                            return (0, ts_poet_1.code) `${from}`;
                        }
                        else if ((0, types_1.isEnum)(valueField)) {
                            return (0, ts_poet_1.code) `${from} as ${valueType}`;
                        }
                        else if ((0, types_1.isLong)(valueField) && options.forceLong === options_1.LongOption.LONG) {
                            return (0, ts_poet_1.code) `Long.fromValue(${from})`;
                        }
                        else if ((0, types_1.isLong)(valueField) && options.forceLong === options_1.LongOption.BIGINT) {
                            return (0, ts_poet_1.code) `BigInt(${from} as string | number | bigint | boolean)`;
                        }
                        else {
                            const cstr = (0, case_1.capitalize)(valueType.toCodeString([]));
                            return (0, ts_poet_1.code) `${utils.globalThis}.${cstr}(${from})`;
                        }
                    }
                    else if ((0, types_1.isAnyValueType)(valueField)) {
                        return (0, ts_poet_1.code) `${from}`;
                    }
                    else if ((0, types_1.isObjectId)(valueField) && options.useMongoObjectId) {
                        return (0, ts_poet_1.code) `${from} as mongodb.ObjectId`;
                    }
                    else if ((0, types_1.isTimestamp)(valueField) &&
                        (options.useDate === options_1.DateOption.DATE ||
                            options.useDate === options_1.DateOption.STRING ||
                            options.useDate === options_1.DateOption.STRING_NANO)) {
                        return (0, ts_poet_1.code) `${from}`;
                    }
                    else if ((0, types_1.isValueType)(ctx, valueField)) {
                        return (0, ts_poet_1.code) `${from}`;
                    }
                    else {
                        const type = (0, types_1.basicTypeName)(ctx, valueField);
                        return (0, ts_poet_1.code) `${type}.fromPartial(${from})`;
                    }
                }
                else if ((0, types_1.isAnyValueType)(field)) {
                    return (0, ts_poet_1.code) `${from}`;
                }
                else {
                    const type = (0, types_1.basicTypeName)(ctx, field);
                    return (0, ts_poet_1.code) `${type}.fromPartial(${from})`;
                }
            }
            else {
                throw new Error(`Unhandled field ${field}`);
            }
        };
        const noDefaultValue = !options.initializeFieldsAsUndefined && (0, types_1.isOptionalProperty)(field, messageDesc.options, options, true);
        // and then use the snippet to handle repeated fields if necessary
        if ((0, types_1.isRepeated)(field)) {
            if ((0, types_1.isMapType)(ctx, messageDesc, field)) {
                const fieldType = (0, types_1.toTypeName)(ctx, messageDesc, field);
                const i = convertFromObjectKey(ctx, messageDesc, field, "key");
                const noValueSnippet = noDefaultValue
                    ? `(${objectProperty} === undefined || ${objectProperty} === null) ? undefined : `
                    : "";
                if ((0, types_1.shouldGenerateJSMapType)(ctx, messageDesc, field)) {
                    chunks.push((0, ts_poet_1.code) `
            ${messageProperty} = ${noValueSnippet} (() => {
              const m = new Map();
              (${objectProperty} as ${fieldType} ?? new Map()).forEach((value, key) => {
                if (value !== undefined) {
                  m.set(key, ${readSnippet("value")});
                }
              });
              return m;
            })();
          `);
                }
                else {
                    chunks.push((0, ts_poet_1.code) `
            ${messageProperty} = ${noValueSnippet} Object.entries(${objectProperty} ?? {}).reduce<${fieldType}>((acc, [key, value]) => {
              if (value !== undefined) {
                acc[${i}] = ${readSnippet("value")};
              }
              return acc;
            }, {});
          `);
                }
            }
            else {
                const fallback = noDefaultValue ? "undefined" : "[]";
                chunks.push((0, ts_poet_1.code) `
          ${messageProperty} = ${objectProperty}?.map((e) => ${readSnippet("e")}) || ${fallback};
        `);
            }
        }
        else if ((0, types_1.isWithinOneOfThatShouldBeUnion)(options, field)) {
            const oneofName = (0, case_1.maybeSnakeToCamel)(messageDesc.oneofDecl[field.oneofIndex].name, options);
            const oneofNameWithMessage = (0, utils_1.getPropertyAccessor)("message", oneofName);
            const oneofNameWithObject = (0, utils_1.getPropertyAccessor)("object", oneofName);
            const v = readSnippet(`${oneofNameWithObject}.${fieldName}`);
            chunks.push((0, ts_poet_1.code) `
        if (
          ${oneofNameWithObject}?.$case === '${fieldName}'
          && ${oneofNameWithObject}?.${fieldName} !== undefined
          && ${oneofNameWithObject}?.${fieldName} !== null
        ) {
          ${oneofNameWithMessage} = { $case: '${fieldName}', ${fieldName}: ${v} };
        }
      `);
        }
        else if (readSnippet(`x`).toCodeString([]) == "x") {
            // An optimized case of the else below that works when `readSnippet` returns the plain input
            const fallback = (0, types_1.isWithinOneOf)(field) || noDefaultValue ? "undefined" : (0, types_1.defaultValue)(ctx, field);
            chunks.push((0, ts_poet_1.code) `${messageProperty} = ${objectProperty} ?? ${fallback};`);
        }
        else {
            const fallback = (0, types_1.isWithinOneOf)(field) || noDefaultValue ? "undefined" : (0, types_1.defaultValue)(ctx, field);
            chunks.push((0, ts_poet_1.code) `
        ${messageProperty} = (${objectProperty} !== undefined && ${objectProperty} !== null)
          ? ${readSnippet(`${objectProperty}`)}
          : ${fallback};
      `);
        }
    });
    // and then wrap up the switch/while/return
    chunks.push((0, ts_poet_1.code) `return message;`);
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
exports.contextTypeVar = "Context extends DataLoaders";
function convertFromObjectKey(ctx, messageDesc, field, variableName) {
    const { keyType, keyField } = (0, types_1.detectMapType)(ctx, messageDesc, field);
    if (keyType.toCodeString([]) === "string") {
        return (0, ts_poet_1.code) `${variableName}`;
    }
    else if ((0, types_1.isLong)(keyField) && (0, types_1.shouldGenerateJSMapType)(ctx, messageDesc, field)) {
        if (ctx.options.forceLong === options_1.LongOption.LONG) {
            return (0, ts_poet_1.code) `${(0, case_1.capitalize)(keyType.toCodeString([]))}.fromValue(${variableName})`;
        }
        else if (ctx.options.forceLong === options_1.LongOption.BIGINT) {
            return (0, ts_poet_1.code) `BigInt(${variableName})`;
        }
        else if (ctx.options.forceLong === options_1.LongOption.STRING) {
            return (0, ts_poet_1.code) `${ctx.utils.globalThis}.String(${variableName})`;
        }
        else {
            return (0, ts_poet_1.code) `${ctx.utils.globalThis}.Number(${variableName})`;
        }
    }
    else if (keyField.type === ts_proto_descriptors_1.FieldDescriptorProto_Type.TYPE_BOOL) {
        return (0, ts_poet_1.code) `${ctx.utils.globalThis}.Boolean(${variableName})`;
    }
    else {
        return (0, ts_poet_1.code) `${ctx.utils.globalThis}.Number(${variableName})`;
    }
}
function convertToObjectKey(ctx, messageDesc, field, variableName) {
    const { keyType, keyField } = (0, types_1.detectMapType)(ctx, messageDesc, field);
    if (keyType.toCodeString([]) === "string") {
        return (0, ts_poet_1.code) `${variableName}`;
    }
    else if ((0, types_1.isLong)(keyField) && (0, types_1.shouldGenerateJSMapType)(ctx, messageDesc, field)) {
        if (ctx.options.forceLong === options_1.LongOption.LONG) {
            return (0, ts_poet_1.code) `${ctx.utils.longToNumber}(${variableName})`;
        }
        else if (ctx.options.forceLong === options_1.LongOption.BIGINT) {
            return (0, ts_poet_1.code) `${variableName}.toString()`;
        }
        else {
            return (0, ts_poet_1.code) `${variableName}`;
        }
    }
    else if (keyField.type === ts_proto_descriptors_1.FieldDescriptorProto_Type.TYPE_BOOL) {
        return (0, ts_poet_1.code) `${ctx.utils.globalThis}.String(${variableName})`;
    }
    else {
        return (0, ts_poet_1.code) `${variableName}`;
    }
}
function maybeReadonly(options) {
    return options.useReadonlyTypes ? "readonly " : "";
}
function maybeAsAny(options) {
    return options.useReadonlyTypes ? " as any" : "";
}
