import * as t from '@babel/types';
import { ProtoField, ProtoType } from '@osmonauts/types';
import { AminoParseContext } from '../../context';
export interface RenderAminoField {
    context: AminoParseContext;
    field: ProtoField;
    currentProtoPath: string;
    isOptional: boolean;
}
export declare const renderAminoField: ({ context, field, currentProtoPath, isOptional }: RenderAminoField) => any;
export interface MakeAminoTypeInterface {
    context: AminoParseContext;
    proto: ProtoType;
}
export declare const makeAminoTypeInterface: ({ context, proto }: MakeAminoTypeInterface) => t.ExportNamedDeclaration;
