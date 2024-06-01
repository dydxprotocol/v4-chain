import * as t from '@babel/types';
import { ProtoEnum } from '@osmonauts/types';
import { ProtoParseContext } from '../context';
export declare const createProtoEnum: (context: ProtoParseContext, name: string, proto: ProtoEnum) => t.ExportNamedDeclaration;
export declare const createProtoEnumFromJSON: (context: ProtoParseContext, name: string, proto: ProtoEnum) => t.ExportNamedDeclaration;
export declare const createProtoEnumToJSON: (context: ProtoParseContext, name: string, proto: ProtoEnum) => t.ExportNamedDeclaration;
