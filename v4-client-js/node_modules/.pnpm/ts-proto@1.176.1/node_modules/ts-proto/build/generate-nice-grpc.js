"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateNiceGrpcService = void 0;
const ts_poet_1 = require("ts-poet");
const case_1 = require("./case");
const sourceInfo_1 = require("./sourceInfo");
const types_1 = require("./types");
const utils_1 = require("./utils");
const CallOptions = (0, ts_poet_1.imp)("t:CallOptions@nice-grpc-common");
const CallContext = (0, ts_poet_1.imp)("t:CallContext@nice-grpc-common");
/**
 * Generates server / client stubs for `nice-grpc` library.
 */
function generateNiceGrpcService(ctx, fileDesc, sourceInfo, serviceDesc) {
    const chunks = [];
    chunks.push(generateServerStub(ctx, sourceInfo, serviceDesc));
    chunks.push(generateClientStub(ctx, sourceInfo, serviceDesc));
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n\n" });
}
exports.generateNiceGrpcService = generateNiceGrpcService;
function generateServerStub(ctx, sourceInfo, serviceDesc) {
    var _a;
    const chunks = [];
    const maybeSuffix = serviceDesc.name.endsWith("Service") ? "" : "Service";
    chunks.push((0, ts_poet_1.code) `export interface ${(0, ts_poet_1.def)(`${serviceDesc.name}${maybeSuffix}Implementation`)}<CallContextExt = {}> {`);
    for (const [index, methodDesc] of serviceDesc.method.entries()) {
        (0, utils_1.assertInstanceOf)(methodDesc, utils_1.FormattedMethodDescriptor);
        const inputType = (0, types_1.messageToTypeName)(ctx, methodDesc.inputType, { keepValueType: true });
        let outputType = (0, types_1.messageToTypeName)(ctx, methodDesc.outputType, { keepValueType: true });
        if (ctx.options.outputPartialMethods) {
            outputType = (0, ts_poet_1.code) `${ctx.utils.DeepPartial}<${outputType}>`;
        }
        const ServerStreamingMethodResult = ctx.utils.NiceGrpcServerStreamingMethodResult;
        const info = sourceInfo.lookup(sourceInfo_1.Fields.service.method, index);
        (0, utils_1.maybeAddComment)(ctx.options, info, chunks, (_a = methodDesc.options) === null || _a === void 0 ? void 0 : _a.deprecated);
        if (methodDesc.clientStreaming) {
            if (methodDesc.serverStreaming) {
                // bidi streaming
                chunks.push((0, ts_poet_1.code) `
          ${(0, case_1.uncapitalize)(methodDesc.name)}(
            request: AsyncIterable<${inputType}>,
            context: ${CallContext} & CallContextExt,
          ): ${ServerStreamingMethodResult}<${outputType}>;
        `);
            }
            else {
                // client streaming
                chunks.push((0, ts_poet_1.code) `
          ${(0, case_1.uncapitalize)(methodDesc.name)}(
            request: AsyncIterable<${inputType}>,
            context: ${CallContext} & CallContextExt,
          ): Promise<${outputType}>;
        `);
            }
        }
        else {
            if (methodDesc.serverStreaming) {
                // server streaming
                chunks.push((0, ts_poet_1.code) `
          ${(0, case_1.uncapitalize)(methodDesc.name)}(
            request: ${inputType},
            context: ${CallContext} & CallContextExt,
          ): ${ServerStreamingMethodResult}<${outputType}>;
        `);
            }
            else {
                // unary
                chunks.push((0, ts_poet_1.code) `
          ${(0, case_1.uncapitalize)(methodDesc.name)}(
            request: ${inputType},
            context: ${CallContext} & CallContextExt,
          ): Promise<${outputType}>;
        `);
            }
        }
    }
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
function generateClientStub(ctx, sourceInfo, serviceDesc) {
    var _a;
    const chunks = [];
    chunks.push((0, ts_poet_1.code) `export interface ${(0, ts_poet_1.def)(`${serviceDesc.name}Client`)}<CallOptionsExt = {}> {`);
    for (const [index, methodDesc] of serviceDesc.method.entries()) {
        (0, utils_1.assertInstanceOf)(methodDesc, utils_1.FormattedMethodDescriptor);
        let inputType = (0, types_1.messageToTypeName)(ctx, methodDesc.inputType, { keepValueType: true });
        if (ctx.options.outputPartialMethods) {
            inputType = (0, ts_poet_1.code) `${ctx.utils.DeepPartial}<${inputType}>`;
        }
        const outputType = (0, types_1.messageToTypeName)(ctx, methodDesc.outputType, { keepValueType: true });
        const info = sourceInfo.lookup(sourceInfo_1.Fields.service.method, index);
        (0, utils_1.maybeAddComment)(ctx.options, info, chunks, (_a = methodDesc.options) === null || _a === void 0 ? void 0 : _a.deprecated);
        if (methodDesc.clientStreaming) {
            if (methodDesc.serverStreaming) {
                // bidi streaming
                chunks.push((0, ts_poet_1.code) `
          ${(0, case_1.uncapitalize)(methodDesc.name)}(
            request: AsyncIterable<${inputType}>,
            options?: ${CallOptions} & CallOptionsExt,
          ): AsyncIterable<${outputType}>;
        `);
            }
            else {
                // client streaming
                chunks.push((0, ts_poet_1.code) `
          ${(0, case_1.uncapitalize)(methodDesc.name)}(
            request: AsyncIterable<${inputType}>,
            options?: ${CallOptions} & CallOptionsExt,
          ): Promise<${outputType}>;
        `);
            }
        }
        else {
            if (methodDesc.serverStreaming) {
                // server streaming
                chunks.push((0, ts_poet_1.code) `
          ${(0, case_1.uncapitalize)(methodDesc.name)}(
            request: ${inputType},
            options?: ${CallOptions} & CallOptionsExt,
          ): AsyncIterable<${outputType}>;
        `);
            }
            else {
                // unary
                chunks.push((0, ts_poet_1.code) `
          ${(0, case_1.uncapitalize)(methodDesc.name)}(
            request: ${inputType},
            options?: ${CallOptions} & CallOptionsExt,
          ): Promise<${outputType}>;
        `);
            }
        }
    }
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
