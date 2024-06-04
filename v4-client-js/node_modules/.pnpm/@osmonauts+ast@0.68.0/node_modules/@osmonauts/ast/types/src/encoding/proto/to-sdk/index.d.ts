import * as t from '@babel/types';
import { ProtoParseContext } from '../../context';
import { ProtoField, ProtoType } from '@osmonauts/types';
export interface ToSDKMethod {
    context: ProtoParseContext;
    field: ProtoField;
    isOptional: boolean;
}
export declare const toSDKMethodFields: (context: ProtoParseContext, name: string, proto: ProtoType) => any[];
export declare const toSDKMethod: (context: ProtoParseContext, name: string, proto: ProtoType) => t.ObjectMethod;
