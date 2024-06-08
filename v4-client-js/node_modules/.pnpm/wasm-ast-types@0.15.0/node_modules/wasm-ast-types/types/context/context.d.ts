import { JSONSchema } from "../types";
export interface ReactQueryOptions {
    enabled?: boolean;
    optionalClient?: boolean;
    version?: 'v3' | 'v4';
    mutations?: boolean;
    camelize?: boolean;
    queryKeys?: boolean;
    queryFactory?: boolean;
}
export interface TSClientOptions {
    enabled?: boolean;
    execExtendsQuery?: boolean;
    noImplicitOverride?: boolean;
}
export interface MessageComposerOptions {
    enabled?: boolean;
}
export interface RecoilOptions {
    enabled?: boolean;
}
export interface TSTypesOptions {
    enabled?: boolean;
    aliasExecuteMsg?: boolean;
}
interface KeyedSchema {
    [key: string]: JSONSchema;
}
export interface IDLObject {
    contract_name: string;
    contract_version: string;
    idl_version: string;
    instantiate: JSONSchema;
    execute: JSONSchema;
    query: JSONSchema;
    migrate: JSONSchema;
    sudo: JSONSchema;
    responses: KeyedSchema;
}
export interface ContractInfo {
    schemas: JSONSchema[];
    responses?: Record<string, JSONSchema>;
    idlObject?: IDLObject;
}
export interface RenderOptions {
    types?: TSTypesOptions;
    recoil?: RecoilOptions;
    messageComposer?: MessageComposerOptions;
    client?: TSClientOptions;
    reactQuery?: ReactQueryOptions;
}
export interface RenderContext {
    contract: ContractInfo;
    options: RenderOptions;
}
export declare const defaultOptions: RenderOptions;
export declare const getDefinitionSchema: (schemas: JSONSchema[]) => JSONSchema;
export declare class RenderContext implements RenderContext {
    contract: ContractInfo;
    utils: string[];
    schema: JSONSchema;
    constructor(contract: ContractInfo, options?: RenderOptions);
    refLookup($ref: string): JSONSchema;
    addUtil(util: string): void;
    getImports(): any[];
}
export {};
