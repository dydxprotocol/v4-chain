import * as t from '@babel/types';
import { RenderAminoField } from '.';
export declare const aminoInterface: {
    defaultType(args: RenderAminoField): t.TSPropertySignature;
    string(args: RenderAminoField): t.TSPropertySignature;
    long(args: RenderAminoField): t.TSPropertySignature;
    height(args: RenderAminoField): t.TSPropertySignature;
    duration(args: RenderAminoField): any;
    timestamp(args: RenderAminoField): any;
    enum(args: RenderAminoField): t.TSPropertySignature;
    enumArray(args: RenderAminoField): t.TSPropertySignature;
    type({ context, field, currentProtoPath, isOptional }: RenderAminoField): any;
    typeArray({ context, field, currentProtoPath, isOptional }: RenderAminoField): any;
    array(args: RenderAminoField): t.TSPropertySignature;
};
