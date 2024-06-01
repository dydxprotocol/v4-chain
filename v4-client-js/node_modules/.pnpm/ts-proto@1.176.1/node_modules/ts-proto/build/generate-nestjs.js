"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateNestjsGrpcServiceMethodsDecorator = exports.generateNestjsServiceClient = exports.generateNestjsServiceController = void 0;
const ts_poet_1 = require("ts-poet");
const types_1 = require("./types");
const sourceInfo_1 = require("./sourceInfo");
const main_1 = require("./main");
const utils_1 = require("./utils");
function generateNestjsServiceController(ctx, fileDesc, sourceInfo, serviceDesc) {
    var _a;
    const { options } = ctx;
    const chunks = [];
    const Metadata = (0, ts_poet_1.imp)("Metadata@@grpc/grpc-js");
    (0, utils_1.maybeAddComment)(options, sourceInfo, chunks, (_a = serviceDesc.options) === null || _a === void 0 ? void 0 : _a.deprecated);
    const t = options.context ? `<${main_1.contextTypeVar}>` : "";
    chunks.push((0, ts_poet_1.code) `
    export interface ${serviceDesc.name}Controller${t} {
  `);
    serviceDesc.method.forEach((methodDesc, index) => {
        var _a;
        (0, utils_1.assertInstanceOf)(methodDesc, utils_1.FormattedMethodDescriptor);
        const info = sourceInfo.lookup(sourceInfo_1.Fields.service.method, index);
        (0, utils_1.maybeAddComment)(options, info, chunks, (_a = serviceDesc.options) === null || _a === void 0 ? void 0 : _a.deprecated);
        const params = [];
        if (options.context) {
            params.push((0, ts_poet_1.code) `ctx: Context`);
        }
        params.push((0, ts_poet_1.code) `request: ${(0, types_1.requestType)(ctx, methodDesc)}`);
        // Use metadata as last argument for interface only configuration
        if (options.addGrpcMetadata) {
            const q = options.addNestjsRestParameter ? "" : "?";
            params.push((0, ts_poet_1.code) `metadata${q}: ${Metadata}`);
        }
        if (options.addNestjsRestParameter) {
            params.push((0, ts_poet_1.code) `...rest: any`);
        }
        // Return observable for interface only configuration, passing returnObservable=true and methodDesc.serverStreaming=true
        let returns;
        if ((0, types_1.isEmptyType)(methodDesc.outputType)) {
            returns = (0, ts_poet_1.code) `void`;
        }
        else if (options.returnObservable || methodDesc.serverStreaming) {
            returns = (0, ts_poet_1.code) `${(0, types_1.responseObservable)(ctx, methodDesc)}`;
        }
        else {
            // generate nestjs union type
            returns = (0, ts_poet_1.code) `
        ${(0, types_1.responsePromise)(ctx, methodDesc)}
        | ${(0, types_1.responseObservable)(ctx, methodDesc)}
        | ${(0, types_1.responseType)(ctx, methodDesc)}
      `;
        }
        chunks.push((0, ts_poet_1.code) `
      ${methodDesc.formattedName}(${(0, ts_poet_1.joinCode)(params, { on: ", " })}): ${returns};
    `);
        if (options.context) {
            const batchMethod = (0, types_1.detectBatchMethod)(ctx, fileDesc, serviceDesc, methodDesc);
            if (batchMethod) {
                const maybeCtx = options.context ? "ctx: Context," : "";
                chunks.push((0, ts_poet_1.code) `
          ${batchMethod.singleMethodName}(
            ${maybeCtx}
            ${(0, utils_1.singular)(batchMethod.inputFieldName)}: ${batchMethod.inputType},
          ): Promise<${batchMethod.outputType}>;
        `);
            }
        }
    });
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n\n" });
}
exports.generateNestjsServiceController = generateNestjsServiceController;
function generateNestjsServiceClient(ctx, fileDesc, sourceInfo, serviceDesc) {
    const { options } = ctx;
    const chunks = [];
    const Metadata = (0, ts_poet_1.imp)("Metadata@@grpc/grpc-js");
    (0, utils_1.maybeAddComment)(options, sourceInfo, chunks);
    const t = options.context ? `<${main_1.contextTypeVar}>` : ``;
    chunks.push((0, ts_poet_1.code) `
    export interface ${serviceDesc.name}Client${t} {
  `);
    serviceDesc.method.forEach((methodDesc, index) => {
        var _a;
        (0, utils_1.assertInstanceOf)(methodDesc, utils_1.FormattedMethodDescriptor);
        const params = [];
        if (options.context) {
            params.push((0, ts_poet_1.code) `ctx: Context`);
        }
        params.push((0, ts_poet_1.code) `request: ${(0, types_1.requestType)(ctx, methodDesc)}`);
        // Use metadata as last argument for interface only configuration
        if (options.addGrpcMetadata) {
            const q = options.addNestjsRestParameter ? "" : "?";
            params.push((0, ts_poet_1.code) `metadata${q}: ${Metadata}`);
        }
        if (options.addNestjsRestParameter) {
            params.push((0, ts_poet_1.code) `...rest: any`);
        }
        // Return observable since nestjs client always returns an Observable
        const returns = (0, types_1.responseObservable)(ctx, methodDesc);
        const info = sourceInfo.lookup(sourceInfo_1.Fields.service.method, index);
        (0, utils_1.maybeAddComment)(options, info, chunks, (_a = methodDesc.options) === null || _a === void 0 ? void 0 : _a.deprecated);
        chunks.push((0, ts_poet_1.code) `
      ${methodDesc.formattedName}(
        ${(0, ts_poet_1.joinCode)(params, { on: "," })}
      ): ${returns};
    `);
        if (options.context) {
            const batchMethod = (0, types_1.detectBatchMethod)(ctx, fileDesc, serviceDesc, methodDesc);
            if (batchMethod) {
                const maybeContext = options.context ? `ctx: Context,` : "";
                chunks.push((0, ts_poet_1.code) `
          ${batchMethod.singleMethodName}(
            ${maybeContext}
            ${(0, utils_1.singular)(batchMethod.inputFieldName)}
          ): Promise<${batchMethod.inputType}>;
        `);
            }
        }
    });
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n\n" });
}
exports.generateNestjsServiceClient = generateNestjsServiceClient;
function generateNestjsGrpcServiceMethodsDecorator(ctx, serviceDesc) {
    const { options } = ctx;
    const GrpcMethod = (0, ts_poet_1.imp)("GrpcMethod@@nestjs/microservices");
    const GrpcStreamMethod = (0, ts_poet_1.imp)("GrpcStreamMethod@@nestjs/microservices");
    const grpcMethods = serviceDesc.method
        .filter((m) => !m.clientStreaming)
        .map((m) => {
        (0, utils_1.assertInstanceOf)(m, utils_1.FormattedMethodDescriptor);
        return m.formattedName;
    })
        .map((n) => `"${n}"`);
    const grpcStreamMethods = serviceDesc.method
        .filter((m) => m.clientStreaming)
        .map((m) => {
        (0, utils_1.assertInstanceOf)(m, utils_1.FormattedMethodDescriptor);
        return m.formattedName;
    })
        .map((n) => `"${n}"`);
    return (0, ts_poet_1.code) `
    export function ${serviceDesc.name}ControllerMethods() {
      return function(constructor: Function) {
        const grpcMethods: string[] = [${grpcMethods.join(", ")}];
        for (const method of grpcMethods) {
          const descriptor: any = Reflect.getOwnPropertyDescriptor(constructor.prototype, method);
          ${GrpcMethod}('${serviceDesc.name}', method)(constructor.prototype[method], method, descriptor);
        }
        const grpcStreamMethods: string[] = [${grpcStreamMethods.join(", ")}];
        for (const method of grpcStreamMethods) {
          const descriptor: any = Reflect.getOwnPropertyDescriptor(constructor.prototype, method);
          ${GrpcStreamMethod}('${serviceDesc.name}', method)(constructor.prototype[method], method, descriptor);
        }
      };
    }
  `;
}
exports.generateNestjsGrpcServiceMethodsDecorator = generateNestjsGrpcServiceMethodsDecorator;
