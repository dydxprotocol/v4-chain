"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateGrpcJsService = void 0;
const ts_poet_1 = require("ts-poet");
const sourceInfo_1 = require("./sourceInfo");
const types_1 = require("./types");
const utils_1 = require("./utils");
const encode_1 = require("./encode");
const CallOptions = (0, ts_poet_1.imp)("t:CallOptions@@grpc/grpc-js");
const ChannelCredentials = (0, ts_poet_1.imp)("ChannelCredentials@@grpc/grpc-js");
const ClientOptions = (0, ts_poet_1.imp)("t:ClientOptions@@grpc/grpc-js");
const Client = (0, ts_poet_1.imp)("Client@@grpc/grpc-js");
const ClientDuplexStream = (0, ts_poet_1.imp)("ClientDuplexStream@@grpc/grpc-js");
const ClientReadableStream = (0, ts_poet_1.imp)("ClientReadableStream@@grpc/grpc-js");
const ClientUnaryCall = (0, ts_poet_1.imp)("t:ClientUnaryCall@@grpc/grpc-js");
const ClientWritableStream = (0, ts_poet_1.imp)("ClientWritableStream@@grpc/grpc-js");
const handleBidiStreamingCall = (0, ts_poet_1.imp)("handleBidiStreamingCall@@grpc/grpc-js");
const handleClientStreamingCall = (0, ts_poet_1.imp)("handleClientStreamingCall@@grpc/grpc-js");
const handleServerStreamingCall = (0, ts_poet_1.imp)("handleServerStreamingCall@@grpc/grpc-js");
const handleUnaryCall = (0, ts_poet_1.imp)("t:handleUnaryCall@@grpc/grpc-js");
const UntypedServiceImplementation = (0, ts_poet_1.imp)("t:UntypedServiceImplementation@@grpc/grpc-js");
const makeGenericClientConstructor = (0, ts_poet_1.imp)("makeGenericClientConstructor@@grpc/grpc-js");
const Metadata = (0, ts_poet_1.imp)("Metadata@@grpc/grpc-js");
const ServiceError = (0, ts_poet_1.imp)("t:ServiceError@@grpc/grpc-js");
/**
 * Generates a service definition and server / client stubs for the
 * `@grpc/grpc-js` library.
 */
function generateGrpcJsService(ctx, fileDesc, sourceInfo, serviceDesc) {
    const { options } = ctx;
    const chunks = [];
    chunks.push(generateServiceDefinition(ctx, fileDesc, sourceInfo, serviceDesc));
    chunks.push(generateServerStub(ctx, sourceInfo, serviceDesc));
    if (options.outputClientImpl) {
        chunks.push(generateClientStub(ctx, sourceInfo, serviceDesc));
        chunks.push(generateClientConstructor(fileDesc, serviceDesc));
    }
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n\n" });
}
exports.generateGrpcJsService = generateGrpcJsService;
function generateServiceDefinition(ctx, fileDesc, sourceInfo, serviceDesc) {
    var _a, _b;
    const chunks = [];
    (0, utils_1.maybeAddComment)(ctx.options, sourceInfo, chunks, (_a = serviceDesc.options) === null || _a === void 0 ? void 0 : _a.deprecated);
    // Service definition type
    const name = (0, ts_poet_1.def)(`${serviceDesc.name}Service`);
    chunks.push((0, ts_poet_1.code) `
    export type ${name} = typeof ${name};
  `);
    // Service definition
    chunks.push((0, ts_poet_1.code) `
    export const ${name} = {
  `);
    for (const [index, methodDesc] of serviceDesc.method.entries()) {
        (0, utils_1.assertInstanceOf)(methodDesc, utils_1.FormattedMethodDescriptor);
        const inputType = (0, types_1.messageToTypeName)(ctx, methodDesc.inputType);
        const outputType = (0, types_1.messageToTypeName)(ctx, methodDesc.outputType);
        const info = sourceInfo.lookup(sourceInfo_1.Fields.service.method, index);
        (0, utils_1.maybeAddComment)(ctx.options, info, chunks, (_b = methodDesc.options) === null || _b === void 0 ? void 0 : _b.deprecated);
        const inputEncoder = (0, encode_1.generateEncoder)(ctx, methodDesc.inputType);
        const outputEncoder = (0, encode_1.generateEncoder)(ctx, methodDesc.outputType);
        const inputDecoder = (0, encode_1.generateDecoder)(ctx, methodDesc.inputType);
        const outputDecoder = (0, encode_1.generateDecoder)(ctx, methodDesc.outputType);
        chunks.push((0, ts_poet_1.code) `
      ${methodDesc.formattedName}: {
        path: '/${(0, utils_1.maybePrefixPackage)(fileDesc, serviceDesc.name)}/${methodDesc.name}',
        requestStream: ${methodDesc.clientStreaming},
        responseStream: ${methodDesc.serverStreaming},
        requestSerialize: (value: ${inputType}) =>
          Buffer.from(${inputEncoder}),
        requestDeserialize: (value: Buffer) => ${inputDecoder},
        responseSerialize: (value: ${outputType}) =>
          Buffer.from(${outputEncoder}),
        responseDeserialize: (value: Buffer) => ${outputDecoder},
      },
    `);
    }
    chunks.push((0, ts_poet_1.code) `} as const;`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
function generateServerStub(ctx, sourceInfo, serviceDesc) {
    var _a;
    const chunks = [];
    chunks.push((0, ts_poet_1.code) `export interface ${(0, ts_poet_1.def)(`${serviceDesc.name}Server`)} extends ${UntypedServiceImplementation} {`);
    for (const [index, methodDesc] of serviceDesc.method.entries()) {
        (0, utils_1.assertInstanceOf)(methodDesc, utils_1.FormattedMethodDescriptor);
        const inputType = (0, types_1.messageToTypeName)(ctx, methodDesc.inputType);
        const outputType = (0, types_1.messageToTypeName)(ctx, methodDesc.outputType);
        const info = sourceInfo.lookup(sourceInfo_1.Fields.service.method, index);
        (0, utils_1.maybeAddComment)(ctx.options, info, chunks, (_a = methodDesc.options) === null || _a === void 0 ? void 0 : _a.deprecated);
        const callType = methodDesc.clientStreaming
            ? methodDesc.serverStreaming
                ? handleBidiStreamingCall
                : handleClientStreamingCall
            : methodDesc.serverStreaming
                ? handleServerStreamingCall
                : handleUnaryCall;
        chunks.push((0, ts_poet_1.code) `
      ${methodDesc.formattedName}: ${callType}<${inputType}, ${outputType}>;
    `);
    }
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
function generateClientStub(ctx, sourceInfo, serviceDesc) {
    var _a;
    const chunks = [];
    chunks.push((0, ts_poet_1.code) `export interface ${(0, ts_poet_1.def)(`${serviceDesc.name}Client`)} extends ${Client} {`);
    for (const [index, methodDesc] of serviceDesc.method.entries()) {
        (0, utils_1.assertInstanceOf)(methodDesc, utils_1.FormattedMethodDescriptor);
        const inputType = (0, types_1.messageToTypeName)(ctx, methodDesc.inputType);
        const outputType = (0, types_1.messageToTypeName)(ctx, methodDesc.outputType);
        const info = sourceInfo.lookup(sourceInfo_1.Fields.service.method, index);
        (0, utils_1.maybeAddComment)(ctx.options, info, chunks, (_a = methodDesc.options) === null || _a === void 0 ? void 0 : _a.deprecated);
        const responseCallback = (0, ts_poet_1.code) `(error: ${ServiceError} | null, response: ${outputType}) => void`;
        if (methodDesc.clientStreaming) {
            if (methodDesc.serverStreaming) {
                // bidi streaming
                chunks.push((0, ts_poet_1.code) `
          ${methodDesc.formattedName}(): ${ClientDuplexStream}<${inputType}, ${outputType}>;
          ${methodDesc.formattedName}(
            options: Partial<${CallOptions}>,
          ): ${ClientDuplexStream}<${inputType}, ${outputType}>;
          ${methodDesc.formattedName}(
            metadata: ${Metadata},
            options?: Partial<${CallOptions}>,
          ): ${ClientDuplexStream}<${inputType}, ${outputType}>;
        `);
            }
            else {
                // client streaming
                chunks.push((0, ts_poet_1.code) `
          ${methodDesc.formattedName}(
            callback: ${responseCallback},
          ): ${ClientWritableStream}<${inputType}>;
          ${methodDesc.formattedName}(
            metadata: ${Metadata},
            callback: ${responseCallback},
          ): ${ClientWritableStream}<${inputType}>;
          ${methodDesc.formattedName}(
            options: Partial<${CallOptions}>,
            callback: ${responseCallback},
          ): ${ClientWritableStream}<${inputType}>;
          ${methodDesc.formattedName}(
            metadata: ${Metadata},
            options: Partial<${CallOptions}>,
            callback: ${responseCallback},
          ): ${ClientWritableStream}<${inputType}>;
        `);
            }
        }
        else {
            if (methodDesc.serverStreaming) {
                // server streaming
                chunks.push((0, ts_poet_1.code) `
          ${methodDesc.formattedName}(
            request: ${inputType},
            options?: Partial<${CallOptions}>,
          ): ${ClientReadableStream}<${outputType}>;
          ${methodDesc.formattedName}(
            request: ${inputType},
            metadata?: ${Metadata},
            options?: Partial<${CallOptions}>,
          ): ${ClientReadableStream}<${outputType}>;
        `);
            }
            else {
                // unary
                chunks.push((0, ts_poet_1.code) `
          ${methodDesc.formattedName}(
            request: ${inputType},
            callback: ${responseCallback},
          ): ${ClientUnaryCall};
          ${methodDesc.formattedName}(
            request: ${inputType},
            metadata: ${Metadata},
            callback: ${responseCallback},
          ): ${ClientUnaryCall};
          ${methodDesc.formattedName}(
            request: ${inputType},
            metadata: ${Metadata},
            options: Partial<${CallOptions}>,
            callback: ${responseCallback},
          ): ${ClientUnaryCall};
        `);
            }
        }
    }
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
function generateClientConstructor(fileDesc, serviceDesc) {
    return (0, ts_poet_1.code) `
    export const ${(0, ts_poet_1.def)(`${serviceDesc.name}Client`)} = ${makeGenericClientConstructor}(
      ${serviceDesc.name}Service,
      '${(0, utils_1.maybePrefixPackage)(fileDesc, serviceDesc.name)}'
    ) as unknown as {
      new (
        address: string,
        credentials: ${ChannelCredentials},
        options?: Partial<${ClientOptions}>,
      ): ${serviceDesc.name}Client;
      service: typeof ${serviceDesc.name}Service;
      serviceName: string;
    }
  `;
}
