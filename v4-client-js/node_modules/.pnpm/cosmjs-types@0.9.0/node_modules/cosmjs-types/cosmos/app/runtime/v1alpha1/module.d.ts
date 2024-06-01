import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "cosmos.app.runtime.v1alpha1";
/** Module is the config object for the runtime module. */
export interface Module {
    /** app_name is the name of the app. */
    appName: string;
    /**
     * begin_blockers specifies the module names of begin blockers
     * to call in the order in which they should be called. If this is left empty
     * no begin blocker will be registered.
     */
    beginBlockers: string[];
    /**
     * end_blockers specifies the module names of the end blockers
     * to call in the order in which they should be called. If this is left empty
     * no end blocker will be registered.
     */
    endBlockers: string[];
    /**
     * init_genesis specifies the module names of init genesis functions
     * to call in the order in which they should be called. If this is left empty
     * no init genesis function will be registered.
     */
    initGenesis: string[];
    /**
     * export_genesis specifies the order in which to export module genesis data.
     * If this is left empty, the init_genesis order will be used for export genesis
     * if it is specified.
     */
    exportGenesis: string[];
    /**
     * override_store_keys is an optional list of overrides for the module store keys
     * to be used in keeper construction.
     */
    overrideStoreKeys: StoreKeyConfig[];
}
/**
 * StoreKeyConfig may be supplied to override the default module store key, which
 * is the module name.
 */
export interface StoreKeyConfig {
    /** name of the module to override the store key of */
    moduleName: string;
    /** the kv store key to use instead of the module name. */
    kvStoreKey: string;
}
export declare const Module: {
    typeUrl: string;
    encode(message: Module, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Module;
    fromJSON(object: any): Module;
    toJSON(message: Module): unknown;
    fromPartial<I extends {
        appName?: string | undefined;
        beginBlockers?: string[] | undefined;
        endBlockers?: string[] | undefined;
        initGenesis?: string[] | undefined;
        exportGenesis?: string[] | undefined;
        overrideStoreKeys?: {
            moduleName?: string | undefined;
            kvStoreKey?: string | undefined;
        }[] | undefined;
    } & {
        appName?: string | undefined;
        beginBlockers?: (string[] & string[] & Record<Exclude<keyof I["beginBlockers"], keyof string[]>, never>) | undefined;
        endBlockers?: (string[] & string[] & Record<Exclude<keyof I["endBlockers"], keyof string[]>, never>) | undefined;
        initGenesis?: (string[] & string[] & Record<Exclude<keyof I["initGenesis"], keyof string[]>, never>) | undefined;
        exportGenesis?: (string[] & string[] & Record<Exclude<keyof I["exportGenesis"], keyof string[]>, never>) | undefined;
        overrideStoreKeys?: ({
            moduleName?: string | undefined;
            kvStoreKey?: string | undefined;
        }[] & ({
            moduleName?: string | undefined;
            kvStoreKey?: string | undefined;
        } & {
            moduleName?: string | undefined;
            kvStoreKey?: string | undefined;
        } & Record<Exclude<keyof I["overrideStoreKeys"][number], keyof StoreKeyConfig>, never>)[] & Record<Exclude<keyof I["overrideStoreKeys"], keyof {
            moduleName?: string | undefined;
            kvStoreKey?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Module>, never>>(object: I): Module;
};
export declare const StoreKeyConfig: {
    typeUrl: string;
    encode(message: StoreKeyConfig, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): StoreKeyConfig;
    fromJSON(object: any): StoreKeyConfig;
    toJSON(message: StoreKeyConfig): unknown;
    fromPartial<I extends {
        moduleName?: string | undefined;
        kvStoreKey?: string | undefined;
    } & {
        moduleName?: string | undefined;
        kvStoreKey?: string | undefined;
    } & Record<Exclude<keyof I, keyof StoreKeyConfig>, never>>(object: I): StoreKeyConfig;
};
