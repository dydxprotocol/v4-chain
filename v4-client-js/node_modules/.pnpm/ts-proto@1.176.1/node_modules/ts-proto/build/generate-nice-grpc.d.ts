import { Code } from "ts-poet";
import { FileDescriptorProto, ServiceDescriptorProto } from "ts-proto-descriptors";
import { Context } from "./context";
import SourceInfo from "./sourceInfo";
/**
 * Generates server / client stubs for `nice-grpc` library.
 */
export declare function generateNiceGrpcService(ctx: Context, fileDesc: FileDescriptorProto, sourceInfo: SourceInfo, serviceDesc: ServiceDescriptorProto): Code;
