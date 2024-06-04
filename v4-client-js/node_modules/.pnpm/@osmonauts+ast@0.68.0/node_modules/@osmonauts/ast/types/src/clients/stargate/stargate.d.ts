import * as t from '@babel/types';
import { GenericParseContext } from '../../encoding';
interface CreateStargateClient {
    name: string;
    options: string;
    context: GenericParseContext;
}
interface CreateStargateClientProtoRegistry {
    registries: string[];
    protoTypeRegistry: string;
    context: GenericParseContext;
}
interface CreateStargateClientOptions {
    name: string;
    aminoConverters: string;
    protoTypeRegistry: string;
    context: GenericParseContext;
}
interface CreateStargateClientAminoConverters {
    aminos: string[];
    aminoConverters: string;
    context: GenericParseContext;
}
export declare const createStargateClient: ({ name, options, context }: CreateStargateClient) => t.ExportNamedDeclaration;
export declare const createStargateClientAminoRegistry: ({ aminos, aminoConverters, context }: CreateStargateClientAminoConverters) => t.ExportNamedDeclaration;
export declare const createStargateClientProtoRegistry: ({ registries, protoTypeRegistry, context }: CreateStargateClientProtoRegistry) => t.ExportNamedDeclaration;
export declare const createStargateClientOptions: ({ name, aminoConverters, protoTypeRegistry, context }: CreateStargateClientOptions) => t.ExportNamedDeclaration;
export {};
