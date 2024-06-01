import { JSONSchema } from "../types";
import { convertUtilsToImportList, getImportStatements } from "./imports";
import deepmerge from "deepmerge";

/// Plugin Types
export interface ReactQueryOptions {
    enabled?: boolean;
    optionalClient?: boolean;
    version?: 'v3' | 'v4';
    mutations?: boolean;
    camelize?: boolean;
    queryKeys?: boolean
    queryFactory?: boolean
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

/// END Plugin Types

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
};
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

export const defaultOptions: RenderOptions = {
    types: {
        enabled: true,
        aliasExecuteMsg: false
    },
    client: {
        enabled: true,
        execExtendsQuery: true,
        noImplicitOverride: false,
    },
    recoil: {
        enabled: false
    },
    messageComposer: {
        enabled: false
    },
    reactQuery: {
        enabled: false,
        optionalClient: false,
        version: 'v3',
        mutations: false,
        camelize: true,
        queryKeys: false
    }
};

export const getDefinitionSchema = (schemas: JSONSchema[]): JSONSchema => {
    const aggregateSchema = {
        definitions: {
            //
        }
    };

    schemas.forEach(schema => {
        schema.definitions = schema.definitions || {};
        aggregateSchema.definitions = {
            ...aggregateSchema.definitions,
            ...schema.definitions
        };
    });

    return aggregateSchema;
};
export class RenderContext implements RenderContext {
    contract: ContractInfo;
    utils: string[] = [];
    schema: JSONSchema;
    constructor(
        contract: ContractInfo,
        options?: RenderOptions
    ) {
        this.contract = contract;
        this.schema = getDefinitionSchema(contract.schemas);
        this.options = deepmerge(defaultOptions, options ?? {});
    }
    refLookup($ref: string) {
        const refName = $ref.replace('#/definitions/', '')
        return this.schema.definitions?.[refName];
    }
    addUtil(util: string) {
        this.utils[util] = true;
    }
    getImports() {
        return getImportStatements(
            convertUtilsToImportList(
                this,
                Object.keys(this.utils)
            )
        );
    }
}
