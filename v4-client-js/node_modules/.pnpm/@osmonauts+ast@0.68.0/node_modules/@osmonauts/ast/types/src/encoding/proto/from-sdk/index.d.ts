import * as t from '@babel/types';
import { ProtoType, ProtoField } from '@osmonauts/types';
import { ProtoParseContext } from '../../context';
export interface FromSDKMethod {
    context: ProtoParseContext;
    field: ProtoField;
    isOptional: boolean;
}
export declare const fromSDKMethodFields: (context: ProtoParseContext, name: string, proto: ProtoType) => t.ObjectProperty[];
export declare const fromSDKMethod: (context: ProtoParseContext, name: string, proto: ProtoType) => t.ObjectMethod;
