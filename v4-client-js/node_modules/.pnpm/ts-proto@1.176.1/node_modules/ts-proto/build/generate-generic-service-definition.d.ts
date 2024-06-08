import { Code } from "ts-poet";
import { FileDescriptorProto, ServiceDescriptorProto } from "ts-proto-descriptors";
import { Context } from "./context";
import SourceInfo from "./sourceInfo";
/**
 * Generates a framework-agnostic service descriptor.
 */
export declare function generateGenericServiceDefinition(ctx: Context, fileDesc: FileDescriptorProto, sourceInfo: SourceInfo, serviceDesc: ServiceDescriptorProto): Code;
