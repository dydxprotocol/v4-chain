import * as t from '@babel/types';
import { ProtoParseContext } from '../../context';
import { ProtoField, ProtoType } from '@osmonauts/types';
export interface ToJSONMethod {
    context: ProtoParseContext;
    field: ProtoField;
    isOneOf: boolean;
    isOptional: boolean;
}
export declare const toJSONMethodFields: (context: ProtoParseContext, name: string, proto: ProtoType) => any[];
export declare const toJSONMethod: (context: ProtoParseContext, name: string, proto: ProtoType) => t.ObjectMethod;
