import * as t from '@babel/types';
import { ProtoParseContext } from '../../context';
import { ProtoType, ProtoField } from '@osmonauts/types';
export interface FromPartialMethod {
    context: ProtoParseContext;
    field: ProtoField;
    isOneOf: boolean;
    isOptional: boolean;
}
export declare const fromPartialMethodFields: (context: ProtoParseContext, name: string, proto: ProtoType) => t.Statement[];
export declare const fromPartialMethod: (context: ProtoParseContext, name: string, proto: ProtoType) => t.ObjectMethod;
