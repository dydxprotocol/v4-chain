import * as t from '@babel/types';
import { ToAminoParseField } from './index';
export declare const toAmino: {
    defaultType(args: ToAminoParseField): t.ObjectProperty;
    long(args: ToAminoParseField): t.ObjectProperty;
    string(args: ToAminoParseField): t.ObjectProperty;
    rawBytes(args: ToAminoParseField): t.ObjectProperty;
    wasmByteCode(args: ToAminoParseField): t.ObjectProperty;
    duration(args: ToAminoParseField): t.ObjectProperty;
    durationString(args: ToAminoParseField): t.ObjectProperty;
    height(args: ToAminoParseField): t.ObjectProperty;
    coin(args: ToAminoParseField): t.ObjectProperty;
    type({ context, field, currentProtoPath, scope, fieldPath, nested, isOptional }: ToAminoParseField): any;
    typeArray({ context, field, currentProtoPath, scope, fieldPath, nested, isOptional }: ToAminoParseField): t.ObjectProperty;
    scalarArray({ context, field, currentProtoPath, scope, nested, isOptional }: ToAminoParseField, arrayTypeAstFunc: Function): t.ObjectProperty;
    pubkey(args: ToAminoParseField): t.ObjectProperty;
};
export declare const arrayTypes: {
    long(varname: string): t.CallExpression;
};
