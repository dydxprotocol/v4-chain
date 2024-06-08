import { FileDescriptorProto, ServiceDescriptorProto } from "ts-proto-descriptors";
import { Code } from "ts-poet";
import SourceInfo from "./sourceInfo";
import { Context } from "./context";
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
export declare function generateService(ctx: Context, fileDesc: FileDescriptorProto, sourceInfo: SourceInfo, serviceDesc: ServiceDescriptorProto): Code;
export declare function generateServiceClientImpl(ctx: Context, fileDesc: FileDescriptorProto, serviceDesc: ServiceDescriptorProto): Code;
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
export declare function generateRpcType(ctx: Context, hasStreamingMethods: boolean): Code;
export declare function generateDataLoadersType(): Code;
export declare function generateDataLoaderOptionsType(): Code;
