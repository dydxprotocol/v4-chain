import { Any, AnySDKType } from "../../../google/protobuf/any";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
/**
 * Config represents the configuration for a Cosmos SDK ABCI app.
 * It is intended that all state machine logic including the version of
 * baseapp and tx handlers (and possibly even Tendermint) that an app needs
 * can be described in a config object. For compatibility, the framework should
 * allow a mixture of declarative and imperative app wiring, however, apps
 * that strive for the maximum ease of maintainability should be able to describe
 * their state machine with a config object alone.
 */
export interface Config {
    /** modules are the module configurations for the app. */
    modules: ModuleConfig[];
    /**
     * golang_bindings specifies explicit interface to implementation type bindings which
     * depinject uses to resolve interface inputs to provider functions.  The scope of this
     * field's configuration is global (not module specific).
     */
    golangBindings: GolangBinding[];
}
/**
 * Config represents the configuration for a Cosmos SDK ABCI app.
 * It is intended that all state machine logic including the version of
 * baseapp and tx handlers (and possibly even Tendermint) that an app needs
 * can be described in a config object. For compatibility, the framework should
 * allow a mixture of declarative and imperative app wiring, however, apps
 * that strive for the maximum ease of maintainability should be able to describe
 * their state machine with a config object alone.
 */
export interface ConfigSDKType {
    modules: ModuleConfigSDKType[];
    golang_bindings: GolangBindingSDKType[];
}
/** ModuleConfig is a module configuration for an app. */
export interface ModuleConfig {
    /**
     * name is the unique name of the module within the app. It should be a name
     * that persists between different versions of a module so that modules
     * can be smoothly upgraded to new versions.
     *
     * For example, for the module cosmos.bank.module.v1.Module, we may chose
     * to simply name the module "bank" in the app. When we upgrade to
     * cosmos.bank.module.v2.Module, the app-specific name "bank" stays the same
     * and the framework knows that the v2 module should receive all the same state
     * that the v1 module had. Note: modules should provide info on which versions
     * they can migrate from in the ModuleDescriptor.can_migration_from field.
     */
    name: string;
    /**
     * config is the config object for the module. Module config messages should
     * define a ModuleDescriptor using the cosmos.app.v1alpha1.is_module extension.
     */
    config?: Any;
    /**
     * golang_bindings specifies explicit interface to implementation type bindings which
     * depinject uses to resolve interface inputs to provider functions.  The scope of this
     * field's configuration is module specific.
     */
    golangBindings: GolangBinding[];
}
/** ModuleConfig is a module configuration for an app. */
export interface ModuleConfigSDKType {
    name: string;
    config?: AnySDKType;
    golang_bindings: GolangBindingSDKType[];
}
/** GolangBinding is an explicit interface type to implementing type binding for dependency injection. */
export interface GolangBinding {
    /** interface_type is the interface type which will be bound to a specific implementation type */
    interfaceType: string;
    /** implementation is the implementing type which will be supplied when an input of type interface is requested */
    implementation: string;
}
/** GolangBinding is an explicit interface type to implementing type binding for dependency injection. */
export interface GolangBindingSDKType {
    interface_type: string;
    implementation: string;
}
export declare const Config: {
    encode(message: Config, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Config;
    fromPartial(object: DeepPartial<Config>): Config;
};
export declare const ModuleConfig: {
    encode(message: ModuleConfig, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ModuleConfig;
    fromPartial(object: DeepPartial<ModuleConfig>): ModuleConfig;
};
export declare const GolangBinding: {
    encode(message: GolangBinding, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GolangBinding;
    fromPartial(object: DeepPartial<GolangBinding>): GolangBinding;
};
