import { ProtoRef, TelescopeOptions } from '@osmonauts/types';
import { ProtoStore } from '@osmonauts/proto-parser';
import { AminoParseContext, GenericParseContext, ProtoParseContext } from '@osmonauts/ast';
import { ServiceMutation, ServiceQuery } from '@osmonauts/types';
export declare const getMutations: (mutations: ServiceMutation[]) => {
    typeUrl: string;
    TypeName: string;
    methodName: string;
}[];
export declare const getAminoProtos: (mutations: ServiceMutation[], store: ProtoStore) => any[];
export declare const buildBaseTypeScriptClass: (context: TelescopeParseContext, name: string, obj: any) => void;
export declare const buildBaseTypeScriptInterface: (context: TelescopeParseContext, name: string, obj: any) => void;
export declare const buildEnums: (context: TelescopeParseContext, name: string, obj: any) => void;
export interface TelescopeParseContext {
    options: TelescopeOptions;
    generic: GenericParseContext;
    proto: ProtoParseContext;
    amino: AminoParseContext;
    store: ProtoStore;
    ref: ProtoRef;
    parsedImports: Record<string, any>;
    body: any[];
    mutations: ServiceMutation[];
    queries: any[];
    services: any[];
    types: any[];
}
export declare class TelescopeParseContext implements TelescopeParseContext {
    constructor(ref: ProtoRef, store: ProtoStore, options: TelescopeOptions);
    hasMutations(): boolean;
    addType(name: string, obj: any, isNested: boolean): void;
    addMutation(mutation: ServiceMutation): void;
    addQuery(query: ServiceQuery): void;
    addService(query: any): void;
    buildBase(): void;
    buildRegistry(): void;
    buildRegistryLoader(): void;
    buildAminoInterfaces(): void;
    buildAminoConverter(): void;
    buildHelperObject(): void;
}
