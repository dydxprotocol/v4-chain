import { DescriptorProto, EnumDescriptorProto, FileDescriptorProto, ServiceDescriptorProto } from "ts-proto-descriptors";
import SourceInfo from "./sourceInfo";
import { Options } from "./options";
type MessageVisitor = (fullName: string, desc: DescriptorProto, sourceInfo: SourceInfo, fullProtoTypeName: string) => void;
type EnumVisitor = (fullName: string, desc: EnumDescriptorProto, sourceInfo: SourceInfo, fullProtoTypeName: string) => void;
export declare function visit(proto: FileDescriptorProto | DescriptorProto, sourceInfo: SourceInfo, messageFn: MessageVisitor, options: Options, enumFn?: EnumVisitor, tsPrefix?: string, protoPrefix?: string): void;
export declare function visitServices(proto: FileDescriptorProto, sourceInfo: SourceInfo, serviceFn: (desc: ServiceDescriptorProto, sourceInfo: SourceInfo) => void): void;
export {};
