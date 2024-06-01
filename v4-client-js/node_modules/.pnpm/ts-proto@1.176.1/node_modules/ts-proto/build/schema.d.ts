import { FileDescriptorProto } from "ts-proto-descriptors";
import { Code } from "ts-poet";
import { Context } from "./context";
import SourceInfo from "./sourceInfo";
export declare function generateSchema(ctx: Context, fileDesc: FileDescriptorProto, sourceInfo: SourceInfo): Code[];
