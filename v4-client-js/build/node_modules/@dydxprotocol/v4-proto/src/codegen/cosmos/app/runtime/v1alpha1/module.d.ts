import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../../helpers";
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
    /**
     * order_migrations defines the order in which module migrations are performed.
     * If this is left empty, it uses the default migration order.
     * https://pkg.go.dev/github.com/cosmos/cosmos-sdk@v0.47.0-alpha2/types/module#DefaultMigrationsOrder
     */
    orderMigrations: string[];
    /**
     * precommiters specifies the module names of the precommiters
     * to call in the order in which they should be called. If this is left empty
     * no precommit function will be registered.
     */
    precommiters: string[];
    /**
     * prepare_check_staters specifies the module names of the prepare_check_staters
     * to call in the order in which they should be called. If this is left empty
     * no preparecheckstate function will be registered.
     */
    prepareCheckStaters: string[];
}
/** Module is the config object for the runtime module. */
export interface ModuleSDKType {
    app_name: string;
    begin_blockers: string[];
    end_blockers: string[];
    init_genesis: string[];
    export_genesis: string[];
    override_store_keys: StoreKeyConfigSDKType[];
    order_migrations: string[];
    precommiters: string[];
    prepare_check_staters: string[];
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
/**
 * StoreKeyConfig may be supplied to override the default module store key, which
 * is the module name.
 */
export interface StoreKeyConfigSDKType {
    module_name: string;
    kv_store_key: string;
}
export declare const Module: {
    encode(message: Module, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Module;
    fromPartial(object: DeepPartial<Module>): Module;
};
export declare const StoreKeyConfig: {
    encode(message: StoreKeyConfig, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): StoreKeyConfig;
    fromPartial(object: DeepPartial<StoreKeyConfig>): StoreKeyConfig;
};
