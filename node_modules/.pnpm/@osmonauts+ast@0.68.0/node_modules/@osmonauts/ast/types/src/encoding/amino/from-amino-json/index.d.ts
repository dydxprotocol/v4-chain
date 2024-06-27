import * as t from '@babel/types';
import { ProtoType, ProtoField } from '@osmonauts/types';
import { AminoParseContext } from '../../context';
export interface FromAminoParseField {
    context: AminoParseContext;
    field: ProtoField;
    currentProtoPath: string;
    scope: string[];
    fieldPath: ProtoField[];
    nested: number;
    isOptional: boolean;
}
export declare const fromAminoParseField: ({ context, field, currentProtoPath, scope: previousScope, fieldPath: previousFieldPath, nested, isOptional }: FromAminoParseField) => any;
interface fromAminoJSON {
    context: AminoParseContext;
    proto: ProtoType;
}
export declare const fromAminoJsonMethod: ({ context, proto }: fromAminoJSON) => t.ArrowFunctionExpression;
export {};
