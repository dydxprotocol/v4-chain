"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.addGrpcWebMisc = exports.generateGrpcMethodDesc = exports.generateGrpcServiceDesc = exports.generateGrpcClientImpl = void 0;
const types_1 = require("./types");
const ts_poet_1 = require("ts-poet");
const utils_1 = require("./utils");
const grpc = (0, ts_poet_1.imp)("grpc@@improbable-eng/grpc-web");
const share = (0, ts_poet_1.imp)("share@rxjs/operators");
const take = (0, ts_poet_1.imp)("take@rxjs/operators");
const BrowserHeaders = (0, ts_poet_1.imp)("BrowserHeaders@browser-headers");
/** Generates a client that uses the `@improbable-web/grpc-web` library. */
function generateGrpcClientImpl(ctx, _fileDesc, serviceDesc) {
    const chunks = [];
    // Define the FooServiceImpl class
    chunks.push((0, ts_poet_1.code) `
    export class ${serviceDesc.name}ClientImpl implements ${serviceDesc.name} {
  `);
    // Create the constructor(rpc: Rpc)
    chunks.push((0, ts_poet_1.code) `
    private readonly rpc: Rpc;

    constructor(rpc: Rpc) {
  `);
    chunks.push((0, ts_poet_1.code) `this.rpc = rpc;`);
    // Bind each FooService method to the FooServiceImpl class
    for (const methodDesc of serviceDesc.method) {
        (0, utils_1.assertInstanceOf)(methodDesc, utils_1.FormattedMethodDescriptor);
        chunks.push((0, ts_poet_1.code) `this.${methodDesc.formattedName} = this.${methodDesc.formattedName}.bind(this);`);
    }
    chunks.push((0, ts_poet_1.code) `}`);
    // Create a method for each FooService method
    for (const methodDesc of serviceDesc.method) {
        chunks.push(generateRpcMethod(ctx, serviceDesc, methodDesc));
    }
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { trim: false, on: "\n" });
}
exports.generateGrpcClientImpl = generateGrpcClientImpl;
/** Creates the RPC methods that client code actually calls. */
function generateRpcMethod(ctx, serviceDesc, methodDesc) {
    (0, utils_1.assertInstanceOf)(methodDesc, utils_1.FormattedMethodDescriptor);
    const { options } = ctx;
    const { useAbortSignal } = options;
    const requestMessage = (0, types_1.requestType)(ctx, methodDesc, false);
    const inputType = (0, types_1.requestType)(ctx, methodDesc, true);
    const returns = (0, types_1.responsePromiseOrObservable)(ctx, methodDesc);
    if (methodDesc.clientStreaming) {
        return (0, ts_poet_1.code) `
    ${methodDesc.formattedName}(
      request: ${inputType},
      metadata?: grpc.Metadata,
      ${useAbortSignal ? "abortSignal?: AbortSignal," : ""}
    ): ${returns} {
      throw new Error('ts-proto does not yet support client streaming!');
    }
  `;
    }
    const method = methodDesc.serverStreaming ? "invoke" : "unary";
    return (0, ts_poet_1.code) `
    ${methodDesc.formattedName}(
      request: ${inputType},
      metadata?: grpc.Metadata,
      ${useAbortSignal ? "abortSignal?: AbortSignal," : ""}
    ): ${returns} {
      return this.rpc.${method}(
        ${methodDescName(serviceDesc, methodDesc)},
        ${requestMessage}.fromPartial(request),
        metadata,
        ${useAbortSignal ? "abortSignal," : ""}
      );
    }
  `;
}
/** Creates the service descriptor that grpc-web needs at runtime. */
function generateGrpcServiceDesc(fileDesc, serviceDesc) {
    return (0, ts_poet_1.code) `
    export const ${serviceDesc.name}Desc = {
      serviceName: "${(0, utils_1.maybePrefixPackage)(fileDesc, serviceDesc.name)}",
    };
  `;
}
exports.generateGrpcServiceDesc = generateGrpcServiceDesc;
/**
 * Creates the method descriptor that grpc-web needs at runtime to make `unary` calls.
 *
 * Note that we take a few liberties in the implementation give we don't 100% match
 * what grpc-web's existing output is, but it works out; see comments in the method
 * implementation.
 */
function generateGrpcMethodDesc(ctx, serviceDesc, methodDesc) {
    const inputType = (0, types_1.requestType)(ctx, methodDesc);
    const outputType = (0, types_1.responseType)(ctx, methodDesc);
    // grpc-web expects this to be a class, but the ts-proto messages are just interfaces.
    //
    // That said, grpc-web's runtime doesn't really use this (at least so far for what ts-proto
    // does), so we could potentially set it to `null!`.
    //
    // However, grpc-web does want messages to have a `.serializeBinary()` method, which again
    // due to the class-less nature of ts-proto's messages, we don't have. So we appropriate
    // this `requestType` as a placeholder for our GrpcWebImpl to Object.assign-in this request
    // message's `serializeBinary` method into the data before handing it off to grpc-web.
    //
    // This makes our data look enough like an object/class that grpc-web works just fine.
    const requestFn = (0, ts_poet_1.code) `{
    serializeBinary() {
      return ${inputType}.encode(this).finish();
    },
  }`;
    // grpc-web also expects this to be a class, but with a static `deserializeBinary` method to
    // create new instances of messages. We again don't have an actual class constructor/symbol
    // to pass to it, but we can make up a lambda that has a `deserializeBinary` that does what
    // we want/what grpc-web's runtime needs.
    const responseFn = (0, ts_poet_1.code) `{
    deserializeBinary(data: Uint8Array) {
      const value = ${outputType}.decode(data);
      return {
        ...value,
        toObject() { return value; },
      };
    }
  }`;
    return (0, ts_poet_1.code) `
    export const ${methodDescName(serviceDesc, methodDesc)}: UnaryMethodDefinitionish = {
      methodName: "${methodDesc.name}",
      service: ${serviceDesc.name}Desc,
      requestStream: false,
      responseStream: ${methodDesc.serverStreaming ? "true" : "false"},
      requestType: ${requestFn} as any,
      responseType: ${responseFn} as any,
    };
  `;
}
exports.generateGrpcMethodDesc = generateGrpcMethodDesc;
function methodDescName(serviceDesc, methodDesc) {
    return `${serviceDesc.name}${methodDesc.name}Desc`;
}
/** Adds misc top-level definitions for grpc-web functionality. */
function addGrpcWebMisc(ctx, hasStreamingMethods) {
    const { options } = ctx;
    const chunks = [];
    chunks.push((0, ts_poet_1.code) `
    interface UnaryMethodDefinitionishR extends ${grpc}.UnaryMethodDefinition<any, any> { requestStream: any; responseStream: any; }
  `);
    chunks.push((0, ts_poet_1.code) `type UnaryMethodDefinitionish = UnaryMethodDefinitionishR;`);
    chunks.push(generateGrpcWebRpcType(ctx, options.returnObservable, hasStreamingMethods));
    chunks.push(generateGrpcWebImpl(ctx, options.returnObservable, hasStreamingMethods));
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n\n" });
}
exports.addGrpcWebMisc = addGrpcWebMisc;
/** Makes an `Rpc` interface to decouple from the low-level grpc-web `grpc.invoke and grpc.unary`/etc. methods. */
function generateGrpcWebRpcType(ctx, returnObservable, hasStreamingMethods) {
    const chunks = [];
    const { options } = ctx;
    const { useAbortSignal } = options;
    chunks.push((0, ts_poet_1.code) `interface Rpc {`);
    const wrapper = returnObservable ? (0, types_1.observableType)(ctx) : "Promise";
    chunks.push((0, ts_poet_1.code) `
    unary<T extends UnaryMethodDefinitionish>(
      methodDesc: T,
      request: any,
      metadata: grpc.Metadata | undefined,
      ${useAbortSignal ? "abortSignal?: AbortSignal," : ""}
    ): ${wrapper}<any>;
  `);
    if (hasStreamingMethods) {
        chunks.push((0, ts_poet_1.code) `
      invoke<T extends UnaryMethodDefinitionish>(
        methodDesc: T,
        request: any,
        metadata: grpc.Metadata | undefined,
        ${useAbortSignal ? "abortSignal?: AbortSignal," : ""}
      ): ${(0, types_1.observableType)(ctx)}<any>;
    `);
    }
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { on: "\n" });
}
/** Implements the `Rpc` interface by making calls using the `grpc.unary` method. */
function generateGrpcWebImpl(ctx, returnObservable, hasStreamingMethods) {
    const options = (0, ts_poet_1.code) `
    {
      transport?: grpc.TransportFactory,
      ${hasStreamingMethods ? "streamingTransport?: grpc.TransportFactory," : ``}
      debug?: boolean,
      metadata?: grpc.Metadata,
      upStreamRetryCodes?: number[],
    }
  `;
    const chunks = [];
    chunks.push((0, ts_poet_1.code) `
    export class GrpcWebImpl {
      private host: string;
      private options: ${options};
      
      constructor(host: string, options: ${options}) {
        this.host = host;
        this.options = options;
      }
  `);
    if (returnObservable) {
        chunks.push(createObservableUnaryMethod(ctx));
    }
    else {
        chunks.push(createPromiseUnaryMethod(ctx));
    }
    if (hasStreamingMethods) {
        chunks.push(createInvokeMethod(ctx));
    }
    chunks.push((0, ts_poet_1.code) `}`);
    return (0, ts_poet_1.joinCode)(chunks, { trim: false });
}
function createPromiseUnaryMethod(ctx) {
    const { options } = ctx;
    const { useAbortSignal } = options;
    const maybeAbortSignal = useAbortSignal
        ? `
      if (abortSignal) abortSignal.addEventListener("abort", () => {
        client.close();
        reject(abortSignal.reason);
      });`
        : "";
    return (0, ts_poet_1.code) `
    unary<T extends UnaryMethodDefinitionish>(
      methodDesc: T,
      _request: any,
      metadata: grpc.Metadata | undefined,
      ${useAbortSignal ? "abortSignal?: AbortSignal," : ""}
    ): Promise<any> {
      const request = { ..._request, ...methodDesc.requestType };
      const maybeCombinedMetadata = metadata && this.options.metadata
        ? new ${BrowserHeaders}({ ...this.options?.metadata.headersMap, ...metadata?.headersMap })
        : metadata ?? this.options.metadata;
      return new Promise((resolve, reject) => {
        ${useAbortSignal ? `const client =` : ""} ${grpc}.unary(methodDesc, {
          request,
          host: this.host,
          metadata: maybeCombinedMetadata ?? {},
          ...(this.options.transport !== undefined ? {transport: this.options.transport} : {}),
          debug: this.options.debug ?? false,
          onEnd: function (response) {
            if (response.status === grpc.Code.OK) {
              resolve(response.message!.toObject());
            } else {
              const err = new ${ctx.utils.GrpcWebError}(response.statusMessage, response.status, response.trailers);
              reject(err);
            }
          },
        });

        ${maybeAbortSignal}
      });
    }
  `;
}
function createObservableUnaryMethod(ctx) {
    const { options } = ctx;
    const { useAbortSignal } = options;
    const maybeAbortSignal = useAbortSignal
        ? `
      if (abortSignal) abortSignal.addEventListener("abort", () => {
        observer.error(abortSignal.reason);
        client.close();
      });`
        : "";
    return (0, ts_poet_1.code) `
    unary<T extends UnaryMethodDefinitionish>(
      methodDesc: T,
      _request: any,
      metadata: grpc.Metadata | undefined,
      ${useAbortSignal ? "abortSignal?: AbortSignal," : ""}
    ): ${(0, types_1.observableType)(ctx)}<any> {
      const request = { ..._request, ...methodDesc.requestType };
      const maybeCombinedMetadata = metadata && this.options.metadata
        ? new ${BrowserHeaders}({ ...this.options?.metadata.headersMap, ...metadata?.headersMap })
        : metadata ?? this.options.metadata;
      return new Observable(observer => {
        ${useAbortSignal ? `const client =` : ""} ${grpc}.unary(methodDesc, {
          request,
          host: this.host,
          metadata: maybeCombinedMetadata ?? {},
          ...(this.options.transport !== undefined ? {transport: this.options.transport} : {}),
          debug: this.options.debug ?? false,
          onEnd: (next) => {
            if (next.status !== 0) {
              const err = new ${ctx.utils.GrpcWebError}(next.statusMessage, next.status, next.trailers);
              observer.error(err);
            } else {
              observer.next(next.message as any);
              observer.complete();
            }
          },
        });


      ${maybeAbortSignal}

      }).pipe(${take}(1));
    } 
  `;
}
function createInvokeMethod(ctx) {
    const { options } = ctx;
    const { useAbortSignal } = options;
    return (0, ts_poet_1.code) `
    invoke<T extends UnaryMethodDefinitionish>(
      methodDesc: T,
      _request: any,
      metadata: grpc.Metadata | undefined,
      ${useAbortSignal ? "abortSignal?: AbortSignal," : ""}
    ): ${(0, types_1.observableType)(ctx)}<any> {
      const upStreamCodes = this.options.upStreamRetryCodes ?? [];
      const DEFAULT_TIMEOUT_TIME: number = 3_000;
      const request = { ..._request, ...methodDesc.requestType };
      const transport = this.options.streamingTransport ?? this.options.transport;
      const maybeCombinedMetadata = metadata && this.options.metadata
        ? new ${BrowserHeaders}({ ...this.options?.metadata.headersMap, ...metadata?.headersMap })
        : metadata ?? this.options.metadata;
      return new Observable(observer => {
        const upStream = (() => {
          const client = ${grpc}.invoke(methodDesc, {
            host: this.host,
            request,
            ...(transport !== undefined ? {transport} : {}),
            metadata: maybeCombinedMetadata ?? {},
            debug: this.options.debug ?? false,
            onMessage: (next) => observer.next(next),
            onEnd: (code: ${grpc}.Code, message: string, trailers: ${grpc}.Metadata) => {
              if (code === 0) {
                observer.complete();
              } else if (upStreamCodes.includes(code)) {
                setTimeout(upStream, DEFAULT_TIMEOUT_TIME);
              } else {
                const err = new Error(message) as any;
                err.code = code;
                err.metadata = trailers;
                observer.error(err);
              }
            },
          });
          ${useAbortSignal
        ? `
          if (abortSignal) {
            const abort = () => {
              observer.error(abortSignal.reason);
              client.close();
            };
            abortSignal.addEventListener("abort", abort);
            observer.add(() => {
              if (abortSignal.aborted) {
                return;
              }

              abortSignal.removeEventListener('abort', abort); 
              client.close();
            });
          } else {
            observer.add(() => client.close());
          }
          `
        : `observer.add(() => client.close());`}
        });
        upStream();
      }).pipe(${share}());
    }
  `;
}
