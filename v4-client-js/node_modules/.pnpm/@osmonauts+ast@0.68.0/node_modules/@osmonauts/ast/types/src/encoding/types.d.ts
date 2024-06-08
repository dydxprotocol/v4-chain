import * as t from '@babel/types';
import { ProtoField } from '@osmonauts/types';
import { GenericParseContext, ProtoParseContext } from './context';
export declare const getFieldNames: (field: ProtoField) => {
    propName: string;
    origName: string;
};
export interface CreateProtoTypeOptions {
    useOriginalCase: boolean;
    typeNamePrefix?: string;
    typeNameSuffix?: string;
}
export declare const createProtoTypeOptionsDefaults: {
    useOriginalCase: boolean;
};
export declare const getMessageName: (name: string, options?: CreateProtoTypeOptions) => string;
export declare const getFieldTypeReference: (context: ProtoParseContext, field: ProtoField, options?: CreateProtoTypeOptions) => any;
export declare const getTSType: (context: GenericParseContext, type: string) => t.TSBooleanKeyword | t.TSNumberKeyword | t.TSStringKeyword | t.TSTypeReference;
export declare const getTSTypeFromGoogleType: (context: GenericParseContext, type: string, options?: CreateProtoTypeOptions) => t.TSStringKeyword | t.TSTypeReference;
export declare const getTSTypeForAmino: (context: GenericParseContext, field: ProtoField) => t.TSBooleanKeyword | t.TSNumberKeyword | t.TSStringKeyword | t.TSTypeReference;
export declare const getTSTypeForProto: (context: GenericParseContext, field: ProtoField) => t.TSBooleanKeyword | t.TSNumberKeyword | t.TSStringKeyword | t.TSTypeReference;
export declare const getDefaultTSTypeFromProtoType: (context: ProtoParseContext, field: ProtoField, isOneOf: boolean) => t.ArrayExpression | t.BooleanLiteral | t.Identifier | t.MemberExpression | t.NewExpression | t.NumericLiteral | t.ObjectExpression | t.StringLiteral;
