import * as t from '@babel/types';
import { ProtoType } from '@osmonauts/types';
import { ProtoParseContext } from '../../context';
import { CreateProtoTypeOptions } from '../../types';
export declare const createProtoType: (context: ProtoParseContext, name: string, proto: ProtoType, options?: CreateProtoTypeOptions) => t.ExportNamedDeclaration;
export declare const createCreateProtoType: (context: ProtoParseContext, name: string, proto: ProtoType) => t.FunctionDeclaration;
