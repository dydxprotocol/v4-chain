import * as t from '@babel/types';
import { ProtoParseContext } from '../../context';
import { ProtoType, ProtoField } from '@osmonauts/types';
export interface DecodeMethod {
    typeName: string;
    context: ProtoParseContext;
    field: ProtoField;
    isOptional: boolean;
}
export declare const decodeMethodFields: (context: ProtoParseContext, name: string, proto: ProtoType) => t.SwitchCase[];
export declare const decodeMethod: (context: ProtoParseContext, name: string, proto: ProtoType) => t.ObjectMethod;
