import { Code } from "ts-poet";
import { FileDescriptorProto, ServiceDescriptorProto } from "ts-proto-descriptors";
import { Context } from "./context";
import SourceInfo from "./sourceInfo";
/**
 * Generates a service definition and server / client stubs for the
 * `@grpc/grpc-js` library.
 */
export declare function generateGrpcJsService(ctx: Context, fileDesc: FileDescriptorProto, sourceInfo: SourceInfo, serviceDesc: ServiceDescriptorProto): Code;
