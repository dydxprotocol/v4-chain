import * as t from '@babel/types';
import { ProtoType } from '@osmonauts/types';
import { ProtoParseContext } from './context';
export declare const createObjectWithMethods: (context: ProtoParseContext, name: string, proto: ProtoType) => t.ExportNamedDeclaration;
