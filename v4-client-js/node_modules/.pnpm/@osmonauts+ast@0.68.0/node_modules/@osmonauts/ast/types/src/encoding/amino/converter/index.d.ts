import * as t from '@babel/types';
import { ProtoType, ProtoRoot } from '@osmonauts/types';
import { AminoParseContext } from '../../context';
interface AminoConverterItemParams {
    root: ProtoRoot;
    context: AminoParseContext;
    proto: ProtoType;
}
export declare const createAminoConverterItem: ({ root, context, proto }: AminoConverterItemParams) => t.ObjectProperty;
interface AminoConverterParams {
    name: string;
    root: ProtoRoot;
    context: AminoParseContext;
    protos: ProtoType[];
}
export declare const createAminoConverter: ({ name, root, context, protos }: AminoConverterParams) => t.ExportNamedDeclaration;
export {};
