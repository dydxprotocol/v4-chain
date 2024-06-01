"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateDataLoaderOptionsType = exports.generateDataLoadersType = exports.generateRpcType = exports.generateServiceClientImpl = exports.generateService = void 0;
const ts_poet_1 = require("ts-poet");
const types_1 = require("./types");
const utils_1 = require("./utils");
const sourceInfo_1 = require("./sourceInfo");
const main_1 = require("./main");
/**
 * Generates an interface for `serviceDesc`.
 *
 * Some RPC frameworks (i.e. Twirp) can use the same interface, i.e.
 * `getFoo(req): Promise<res>` for the client-side and server-side,
 * which is the intent for this interface.
 *
 * Other RPC frameworks (i.e. NestJS) that need different client-side
 * vs. server-side code/interfaces are handled separately.
 */
function generateService(ctx, fileDesc, sourceInfo, serviceDesc) {
    var _a;
    const { options } = ctx;
    const chunks = [];
    (0, utils_1.maybeAddComment)(options, sourceInfo, chunks, (_a = serviceDesc.options) === null || _a === void 0 ? void 0 : _a.deprecated);
    const maybeTypeVar = options.context ? `<${main_1.contextTypeVar}>` : "";
    chunks.push((0, ts_poet_1.code) `export interface ${(0, ts_poet_1.def)(serviceDesc.name)}${maybeTypeVar} {`);
    serviceDesc.method.forEach((methodDesc, index) => {
        var _a;
        (0, utils_1.assertInstanceOf)(methodDesc, utils_1.FormattedMethodDescriptor);
        const info = sourceInfo.lookup(sourceInfo_1.Fields.service.method, index);
        (0, utils_1.maybeAddComment)(options, info, chunks, (_a = methodDesc.options) === null || _a === void 0 ? void 0 : _a.deprecated);
        const params = [];
        if (options.context) {
            params.push((0, ts_poet_1.code) `ctx: Context`);
        }
        // the grpc-web clients auto-`fromPartial` the input before handing off to grpc-web's
        // serde runtime, so it's okay to accept partial results from the client
        const partialInput = options.outputClientImpl === "grpc-web";
        const inputType = (0, types_1.requestType)(ctx, methodDesc, partialInput);
        params.push((0, ts_poet_1.code) `request: ${inputType}`);
        // Use metadata as last argument for interface only configuration
        if (options.outputClientImpl === "grpc-web") {
            // We have to use grpc.Metadata where grpc will come from @improbable-eng
            params.push((0, ts_poet_1.code) `metadata?: grpc.Metadata`);
        }
        else if (options.metadataType) {
            // custom `metadataType` has precedence over `addGrpcMetadata` that injects Metadata from grpc-js
            const Metadata = (0, ts_poet_1.imp)(options.metadataType);
            params.push((0, ts_poet_1.code) `metadata?: ${Metadata}`);
        }
        else if (options.addGrpcMetadata) {
            const Metadata = (0, ts_poet_1.imp)("Metadata@@grpc/grpc-js");
            params.push((0, ts_poet_1.code) `metadata?: ${Metadata}`);
        }
        if (options.useAbortSignal) {
            params.push((0, ts_poet_1.code) `abortSignal?: AbortSignal`);
        }
        if (options.addNestjsRestParameter) {
            params.push((0, ts_poet_1.code) `...rest: any`);
        }
        chunks.push((0, ts_poet_1.code) `${methodDesc.formattedName}(${(0, ts_poet_1.joinCode)(params, { on: "," })}): ${(0, types_1.responsePromiseOrObservable)(ctx, methodDesc)};`);
        // If this is a batch method, auto-generate the singular version of it
        if (options.context) {
            const batchMethod = (0, types_1.detectBatchMethod)(ctx, fileDesc, serviceDesc, methodDesc);
            if (batchMethod) {
                chunks.push((0, ts_poet_1.code) `${batchMethod.singleMethodName}(
          ctx: Context,
          ${(0, utils_1.singular)(batchMethod.inputFieldName)}: ${batchMethod.inputType},
        ): Promise<${batchMethod.outputType}>;`);
            }
        }
    });
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
exports.generateService = generateService;
function generateRegularRpcMethod(ctx, methodDesc) {
    (0, utils_1.assertInstanceOf)(methodDesc, utils_1.FormattedMethodDescriptor);
    const { options } = ctx;
    const Reader = (0, utils_1.impFile)(ctx.options, "Reader@protobufjs/minimal");
    const rawInputType = (0, types_1.rawRequestType)(ctx, methodDesc, { keepValueType: true });
    const inputType = (0, types_1.requestType)(ctx, methodDesc);
    const rawOutputType = (0, types_1.responseType)(ctx, methodDesc, { keepValueType: true });
    const metadataType = options.metadataType ? (0, ts_poet_1.imp)(options.metadataType) : (0, ts_poet_1.imp)("Metadata@@grpc/grpc-js");
    const params = [
        ...(options.context ? [(0, ts_poet_1.code) `ctx: Context`] : []),
        (0, ts_poet_1.code) `request: ${inputType}`,
        ...(options.metadataType || options.addGrpcMetadata ? [(0, ts_poet_1.code) `metadata?: ${metadataType}`] : []),
        ...(options.useAbortSignal ? [(0, ts_poet_1.code) `abortSignal?: AbortSignal`] : []),
    ];
    const maybeCtx = options.context ? "ctx," : "";
    const maybeMetadata = options.addGrpcMetadata ? "metadata," : "";
    const maybeAbortSignal = options.useAbortSignal ? "abortSignal || undefined," : "";
    let errorHandler;
    if (options.rpcErrorHandler) {
        errorHandler = (0, ts_poet_1.code) `
      if (this.rpc.handleError) {
        return Promise.reject(this.rpc.handleError(this.service, "${methodDesc.name}", error));
      }
      return Promise.reject(error);
    `;
    }
    let encode = (0, ts_poet_1.code) `${rawInputType}.encode(request).finish()`;
    let beforeRequest;
    if (options.rpcBeforeRequest && !methodDesc.clientStreaming) {
        beforeRequest = generateBeforeRequest(methodDesc.name);
    }
    else if (methodDesc.clientStreaming && options.rpcBeforeRequest) {
        encode = (0, ts_poet_1.code) `{const encodedRequest = ${encode}; ${generateBeforeRequest(methodDesc.name, "encodedRequest")}; return encodedRequest}`;
    }
    let decode = (0, ts_poet_1.code) `${rawOutputType}.decode(${Reader}.create(data))`;
    if (options.rpcAfterResponse) {
        decode = (0, ts_poet_1.code) `
      const response = ${rawOutputType}.decode(${Reader}.create(data));
      if (this.rpc.afterResponse) {
        this.rpc.afterResponse(this.service, "${methodDesc.name}", response);
      }
      return response;
    `;
    }
    // if (options.useDate && rawOutputType.toString().includes("Timestamp")) {
    //   decode = code`data => ${utils.fromTimestamp}(${rawOutputType}.decode(${Reader}.create(data)))`;
    // }
    if (methodDesc.clientStreaming) {
        if (options.useAsyncIterable) {
            encode = (0, ts_poet_1.code) `${rawInputType}.encodeTransform(request)`;
        }
        else {
            encode = (0, ts_poet_1.code) `request.pipe(${(0, ts_poet_1.imp)("map@rxjs/operators")}(request => ${encode}))`;
        }
    }
    const returnStatement = createDefaultServiceReturn(ctx, methodDesc, decode, errorHandler);
    let returnVariable;
    if (options.returnObservable || methodDesc.serverStreaming) {
        returnVariable = "result";
    }
    else {
        returnVariable = "promise";
    }
    let rpcMethod;
    if (methodDesc.clientStreaming && methodDesc.serverStreaming) {
        rpcMethod = "bidirectionalStreamingRequest";
    }
    else if (methodDesc.serverStreaming) {
        rpcMethod = "serverStreamingRequest";
    }
    else if (methodDesc.clientStreaming) {
        rpcMethod = "clientStreamingRequest";
    }
    else {
        rpcMethod = "request";
    }
    return (0, ts_poet_1.code) `
    ${methodDesc.formattedName}(
      ${(0, ts_poet_1.joinCode)(params, { on: "," })}
    ): ${(0, types_1.responsePromiseOrObservable)(ctx, methodDesc)} {
      const data = ${encode}; ${beforeRequest ? beforeRequest : ""}
      const ${returnVariable} = this.rpc.${rpcMethod}(
        ${maybeCtx}
        this.service,
        "${methodDesc.name}",
        data,
        ${maybeMetadata}
        ${maybeAbortSignal}
      );
      return ${returnStatement};
    }
  `;
}
function generateBeforeRequest(methodName, requestVariableName = "request") {
    return (0, ts_poet_1.code) `
    if (this.rpc.beforeRequest) {
      this.rpc.beforeRequest(this.service, "${methodName}", ${requestVariableName});
    }`;
}
function createDefaultServiceReturn(ctx, methodDesc, decode, errorHandler) {
    const { options } = ctx;
    const rawOutputType = (0, types_1.responseType)(ctx, methodDesc, { keepValueType: true });
    const returnStatement = (0, utils_1.arrowFunction)("data", decode, !options.rpcAfterResponse);
    if (options.returnObservable || methodDesc.serverStreaming) {
        if (options.useAsyncIterable) {
            return (0, ts_poet_1.code) `${rawOutputType}.decodeTransform(result)`;
        }
        else {
            if (errorHandler) {
                const tc = (0, utils_1.arrowFunction)("data", (0, utils_1.tryCatchBlock)(decode, (0, ts_poet_1.code) `throw error`), !options.rpcAfterResponse);
                return (0, ts_poet_1.code) `result.pipe(${(0, ts_poet_1.imp)("map@rxjs/operators")}(${tc}))`;
            }
            return (0, ts_poet_1.code) `result.pipe(${(0, ts_poet_1.imp)("map@rxjs/operators")}(${returnStatement}))`;
        }
    }
    if (errorHandler) {
        if (!options.rpcAfterResponse) {
            decode = (0, ts_poet_1.code) `return ${decode}`;
        }
        return (0, ts_poet_1.code) `promise.then(${(0, utils_1.arrowFunction)("data", (0, utils_1.tryCatchBlock)(decode, (0, ts_poet_1.code) `return Promise.reject(error);`), false)}).catch(${(0, utils_1.arrowFunction)("error", errorHandler, false)})`;
    }
    return (0, ts_poet_1.code) `promise.then(${returnStatement})`;
}
function generateServiceClientImpl(ctx, fileDesc, serviceDesc) {
    const { options } = ctx;
    const chunks = [];
    // Determine information about the service.
    const { name } = serviceDesc;
    const serviceName = (0, utils_1.maybePrefixPackage)(fileDesc, serviceDesc.name);
    // Define the service name constant.
    const serviceNameConst = `${name}ServiceName`;
    chunks.push((0, ts_poet_1.code) `export const ${serviceNameConst} = "${serviceName}";`);
    // Define the FooServiceImpl class
    const i = options.context ? `${name}<Context>` : name;
    const t = options.context ? `<${main_1.contextTypeVar}>` : "";
    chunks.push((0, ts_poet_1.code) `export class ${name}ClientImpl${t} implements ${(0, ts_poet_1.def)(i)} {`);
    // Create the constructor(rpc: Rpc)
    const rpcType = options.context ? "Rpc<Context>" : "Rpc";
    chunks.push((0, ts_poet_1.code) `private readonly rpc: ${rpcType};`);
    chunks.push((0, ts_poet_1.code) `private readonly service: string;`);
    chunks.push((0, ts_poet_1.code) `constructor(rpc: ${rpcType}, opts?: {service?: string}) {`);
    chunks.push((0, ts_poet_1.code) `this.service = opts?.service || ${serviceNameConst};`);
    chunks.push((0, ts_poet_1.code) `this.rpc = rpc;`);
    // Bind each FooService method to the FooServiceImpl class
    for (const methodDesc of serviceDesc.method) {
        (0, utils_1.assertInstanceOf)(methodDesc, utils_1.FormattedMethodDescriptor);
        chunks.push((0, ts_poet_1.code) `this.${methodDesc.formattedName} = this.${methodDesc.formattedName}.bind(this);`);
    }
    chunks.push((0, ts_poet_1.code) `}`);
    // Create a method for each FooService method
    for (const methodDesc of serviceDesc.method) {
        // See if this this fuzzy matches to a batchable method
        if (options.context) {
            const batchMethod = (0, types_1.detectBatchMethod)(ctx, fileDesc, serviceDesc, methodDesc);
            if (batchMethod) {
                chunks.push(generateBatchingRpcMethod(ctx, batchMethod));
            }
        }
        if (options.context && methodDesc.name.match(/^Get[A-Z]/)) {
            chunks.push(generateCachingRpcMethod(ctx, fileDesc, serviceDesc, methodDesc));
        }
        else {
            chunks.push(generateRegularRpcMethod(ctx, methodDesc));
        }
    }
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.code) `${chunks}`;
}
exports.generateServiceClientImpl = generateServiceClientImpl;
/** We've found a BatchXxx method, create a synthetic GetXxx method that calls it. */
function generateBatchingRpcMethod(ctx, batchMethod) {
    const { methodDesc, singleMethodName, inputFieldName, inputType, outputFieldName, outputType, mapType, uniqueIdentifier, } = batchMethod;
    (0, utils_1.assertInstanceOf)(methodDesc, utils_1.FormattedMethodDescriptor);
    const { options } = ctx;
    const hash = options.esModuleInterop ? (0, ts_poet_1.imp)("hash=object-hash") : (0, ts_poet_1.imp)("hash*object-hash");
    const dataloader = options.esModuleInterop ? (0, ts_poet_1.imp)("DataLoader=dataloader") : (0, ts_poet_1.imp)("DataLoader*dataloader");
    // Create the `(keys) => ...` lambda we'll pass to the DataLoader constructor
    const lambda = [];
    lambda.push((0, ts_poet_1.code) `
    (${inputFieldName}) => {
      const request = { ${inputFieldName} };
  `);
    if (mapType) {
        // If the return type is a map, lookup each key in the result
        lambda.push((0, ts_poet_1.code) `
      return this.${methodDesc.formattedName}(ctx, request as any).then(res => {
        return ${inputFieldName}.map(key => res.${outputFieldName}[key] ?? ${ctx.utils.fail}())
      });
    `);
    }
    else {
        // Otherwise assume they come back in order
        lambda.push((0, ts_poet_1.code) `
      return this.${methodDesc.formattedName}(ctx, request as any).then(res => res.${outputFieldName})
    `);
    }
    lambda.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.code) `
    ${singleMethodName}(
      ctx: Context,
      ${(0, utils_1.singular)(inputFieldName)}: ${inputType}
    ): Promise<${outputType}> {
      const dl = ctx.getDataLoader("${uniqueIdentifier}", () => {
        return new ${dataloader}<${inputType}, ${outputType}, string>(
          ${(0, ts_poet_1.joinCode)(lambda)},
          { cacheKeyFn: ${hash}, ...ctx.rpcDataLoaderOptions }
        );
      });
      return dl.load(${(0, utils_1.singular)(inputFieldName)});
    }
  `;
}
/** We're not going to batch, but use DataLoader for per-request caching. */
function generateCachingRpcMethod(ctx, fileDesc, serviceDesc, methodDesc) {
    (0, utils_1.assertInstanceOf)(methodDesc, utils_1.FormattedMethodDescriptor);
    const { options } = ctx;
    const hash = options.esModuleInterop ? (0, ts_poet_1.imp)("hash=object-hash") : (0, ts_poet_1.imp)("hash*object-hash");
    const dataloader = options.esModuleInterop ? (0, ts_poet_1.imp)("DataLoader=dataloader") : (0, ts_poet_1.imp)("DataLoader*dataloader");
    const inputType = (0, types_1.requestType)(ctx, methodDesc);
    const outputType = (0, types_1.responseType)(ctx, methodDesc);
    const uniqueIdentifier = `${(0, utils_1.maybePrefixPackage)(fileDesc, serviceDesc.name)}.${methodDesc.name}`;
    const Reader = (0, utils_1.impFile)(ctx.options, "Reader@protobufjs/minimal");
    const lambda = (0, ts_poet_1.code) `
    (requests) => {
      const responses = requests.map(async request => {
        const data = ${inputType}.encode(request).finish()
        const response = await this.rpc.request(ctx, "${(0, utils_1.maybePrefixPackage)(fileDesc, serviceDesc.name)}", "${methodDesc.name}", data);
        return ${outputType}.decode(${Reader}.create(response));
      });
      return Promise.all(responses);
    }
  `;
    return (0, ts_poet_1.code) `
    ${methodDesc.formattedName}(
      ctx: Context,
      request: ${inputType},
    ): Promise<${outputType}> {
      const dl = ctx.getDataLoader("${uniqueIdentifier}", () => {
        return new ${dataloader}<${inputType}, ${outputType}, string>(
          ${lambda},
          { cacheKeyFn: ${hash}, ...ctx.rpcDataLoaderOptions },
        );
      });
      return dl.load(request);
    }
  `;
}
/**
 * Creates an `Rpc.request(service, method, data)` abstraction.
 *
 * This lets clients pass in their own request-promise-ish client.
 *
 * This also requires clientStreamingRequest, serverStreamingRequest and
 * bidirectionalStreamingRequest methods if any of the RPCs is streaming.
 *
 * We don't export this because if a project uses multiple `*.proto` files,
 * we don't want our the barrel imports in `index.ts` to have multiple `Rpc`
 * types.
 */
function generateRpcType(ctx, hasStreamingMethods) {
    const { options } = ctx;
    const metadata = options.metadataType ? (0, ts_poet_1.imp)(options.metadataType) : (0, ts_poet_1.imp)("Metadata@@grpc/grpc-js");
    const metadataType = metadata.symbol;
    const maybeContext = options.context ? "<Context>" : "";
    const maybeContextParam = options.context ? "ctx: Context," : "";
    const maybeMetadataParam = options.metadataType || options.addGrpcMetadata ? `metadata?: ${metadataType},` : "";
    const maybeAbortSignalParam = options.useAbortSignal ? "abortSignal?: AbortSignal," : "";
    const methods = [[(0, ts_poet_1.code) `request`, (0, ts_poet_1.code) `Uint8Array`, (0, ts_poet_1.code) `Promise<Uint8Array>`]];
    const additionalMethods = [];
    if (options.rpcBeforeRequest) {
        additionalMethods.push((0, ts_poet_1.code) `beforeRequest?<T extends { [k in keyof T]: unknown }>(service: string, method: string, request: T): void;`);
    }
    if (options.rpcAfterResponse) {
        additionalMethods.push((0, ts_poet_1.code) `afterResponse?<T extends { [k in keyof T]: unknown }>(service: string, method: string, response: T): void;`);
    }
    if (options.rpcErrorHandler) {
        additionalMethods.push((0, ts_poet_1.code) `handleError?(service: string, method: string, error: globalThis.Error): globalThis.Error;`);
    }
    if (hasStreamingMethods) {
        const observable = (0, types_1.observableType)(ctx, true);
        methods.push([(0, ts_poet_1.code) `clientStreamingRequest`, (0, ts_poet_1.code) `${observable}<Uint8Array>`, (0, ts_poet_1.code) `Promise<Uint8Array>`]);
        methods.push([(0, ts_poet_1.code) `serverStreamingRequest`, (0, ts_poet_1.code) `Uint8Array`, (0, ts_poet_1.code) `${observable}<Uint8Array>`]);
        methods.push([
            (0, ts_poet_1.code) `bidirectionalStreamingRequest`,
            (0, ts_poet_1.code) `${observable}<Uint8Array>`,
            (0, ts_poet_1.code) `${observable}<Uint8Array>`,
        ]);
    }
    const chunks = [];
    chunks.push((0, ts_poet_1.code) `    interface Rpc${maybeContext} {`);
    methods.forEach((method) => {
        chunks.push((0, ts_poet_1.code) `
      ${method[0]}(
        ${maybeContextParam}
        service: string,
        method: string,
        data: ${method[1]},
        ${maybeMetadataParam}
        ${maybeAbortSignalParam}
      ): ${method[2]};`);
    });
    additionalMethods.forEach((method) => chunks.push(method));
    chunks.push((0, ts_poet_1.code) `    }`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
exports.generateRpcType = generateRpcType;
function generateDataLoadersType() {
    // TODO Maybe should be a generic `Context.get<T>(id, () => T): T` method
    return (0, ts_poet_1.code) `
    export interface DataLoaders {
      rpcDataLoaderOptions?: DataLoaderOptions;
      getDataLoader<T>(identifier: string, constructorFn: () => T): T;
    }
  `;
}
exports.generateDataLoadersType = generateDataLoadersType;
function generateDataLoaderOptionsType() {
    return (0, ts_poet_1.code) `
    export interface DataLoaderOptions {
      cache?: boolean;
    }
  `;
}
exports.generateDataLoaderOptionsType = generateDataLoaderOptionsType;
