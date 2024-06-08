import { FileDescriptorProto, ServiceDescriptorProto } from "ts-proto-descriptors";
import { Code } from "ts-poet";
import SourceInfo from "./sourceInfo";
import { Context } from "./context";
export declare function generateNestjsServiceController(ctx: Context, fileDesc: FileDescriptorProto, sourceInfo: SourceInfo, serviceDesc: ServiceDescriptorProto): Code;
export declare function generateNestjsServiceClient(ctx: Context, fileDesc: FileDescriptorProto, sourceInfo: SourceInfo, serviceDesc: ServiceDescriptorProto): Code;
export declare function generateNestjsGrpcServiceMethodsDecorator(ctx: Context, serviceDesc: ServiceDescriptorProto): Code;
