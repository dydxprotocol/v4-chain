import { MethodDescriptorProto, FileDescriptorProto, ServiceDescriptorProto } from "ts-proto-descriptors";
import { Code } from "ts-poet";
import { Context } from "./context";
/** Generates a client that uses the `@improbable-web/grpc-web` library. */
export declare function generateGrpcClientImpl(ctx: Context, _fileDesc: FileDescriptorProto, serviceDesc: ServiceDescriptorProto): Code;
/** Creates the service descriptor that grpc-web needs at runtime. */
export declare function generateGrpcServiceDesc(fileDesc: FileDescriptorProto, serviceDesc: ServiceDescriptorProto): Code;
/**
 * Creates the method descriptor that grpc-web needs at runtime to make `unary` calls.
 *
 * Note that we take a few liberties in the implementation give we don't 100% match
 * what grpc-web's existing output is, but it works out; see comments in the method
 * implementation.
 */
export declare function generateGrpcMethodDesc(ctx: Context, serviceDesc: ServiceDescriptorProto, methodDesc: MethodDescriptorProto): Code;
/** Adds misc top-level definitions for grpc-web functionality. */
export declare function addGrpcWebMisc(ctx: Context, hasStreamingMethods: boolean): Code;
