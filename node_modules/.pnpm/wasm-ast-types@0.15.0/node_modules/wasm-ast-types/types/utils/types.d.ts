import * as t from '@babel/types';
import { TSTypeAnnotation } from '@babel/types';
import { RenderContext } from '../context';
import { JSONSchema } from '../types';
export declare function getResponseType(context: RenderContext, underscoreName: string): string;
export declare const getTypeFromRef: ($ref: any) => t.TSTypeReference;
export declare const detectType: (type: string | string[]) => {
    type: string;
    optional: boolean;
};
export declare const getTypeInfo: (info: JSONSchema) => {
    type: any;
    optional: any;
};
export declare const getType: (type: string) => t.TSBooleanKeyword | t.TSNumberKeyword | t.TSStringKeyword;
export declare const getPropertyType: (context: RenderContext, schema: JSONSchema, prop: string) => {
    type: any;
    optional: boolean;
};
export declare function getPropertySignatureFromProp(context: RenderContext, jsonschema: JSONSchema, prop: string, camelize: boolean): {
    type: string;
    key: t.Identifier;
    typeAnnotation: t.TSTypeAnnotation;
    optional: boolean;
};
export declare const getParamsTypeAnnotation: (context: RenderContext, jsonschema: any, camelize?: boolean) => t.TSTypeAnnotation;
export declare const createTypedObjectParams: (context: RenderContext, jsonschema: JSONSchema, camelize?: boolean) => t.ObjectPattern;
