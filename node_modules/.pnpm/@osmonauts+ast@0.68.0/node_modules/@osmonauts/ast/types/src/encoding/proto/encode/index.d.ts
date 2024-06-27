import * as t from '@babel/types';
import { ProtoType, ProtoField } from '@osmonauts/types';
import { ProtoParseContext } from '../../context';
export interface EncodeMethod {
    typeName: string;
    context: ProtoParseContext;
    field: ProtoField;
    isOptional: boolean;
    isOneOf: boolean;
}
export declare const encodeMethodFields: (context: ProtoParseContext, name: string, proto: ProtoType) => any[];
export declare const encodeMethod: (context: ProtoParseContext, name: string, proto: ProtoType) => t.ObjectMethod;
