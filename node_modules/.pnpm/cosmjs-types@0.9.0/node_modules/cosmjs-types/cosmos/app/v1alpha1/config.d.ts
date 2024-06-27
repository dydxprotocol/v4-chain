import { Any } from "../../../google/protobuf/any";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.app.v1alpha1";
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
/** GolangBinding is an explicit interface type to implementing type binding for dependency injection. */
export interface GolangBinding {
    /** interface_type is the interface type which will be bound to a specific implementation type */
    interfaceType: string;
    /** implementation is the implementing type which will be supplied when an input of type interface is requested */
    implementation: string;
}
export declare const Config: {
    typeUrl: string;
    encode(message: Config, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Config;
    fromJSON(object: any): Config;
    toJSON(message: Config): unknown;
    fromPartial<I extends {
        modules?: {
            name?: string | undefined;
            config?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            golangBindings?: {
                interfaceType?: string | undefined;
                implementation?: string | undefined;
            }[] | undefined;
        }[] | undefined;
        golangBindings?: {
            interfaceType?: string | undefined;
            implementation?: string | undefined;
        }[] | undefined;
    } & {
        modules?: ({
            name?: string | undefined;
            config?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            golangBindings?: {
                interfaceType?: string | undefined;
                implementation?: string | undefined;
            }[] | undefined;
        }[] & ({
            name?: string | undefined;
            config?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            golangBindings?: {
                interfaceType?: string | undefined;
                implementation?: string | undefined;
            }[] | undefined;
        } & {
            name?: string | undefined;
            config?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["modules"][number]["config"], keyof Any>, never>) | undefined;
            golangBindings?: ({
                interfaceType?: string | undefined;
                implementation?: string | undefined;
            }[] & ({
                interfaceType?: string | undefined;
                implementation?: string | undefined;
            } & {
                interfaceType?: string | undefined;
                implementation?: string | undefined;
            } & Record<Exclude<keyof I["modules"][number]["golangBindings"][number], keyof GolangBinding>, never>)[] & Record<Exclude<keyof I["modules"][number]["golangBindings"], keyof {
                interfaceType?: string | undefined;
                implementation?: string | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["modules"][number], keyof ModuleConfig>, never>)[] & Record<Exclude<keyof I["modules"], keyof {
            name?: string | undefined;
            config?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            golangBindings?: {
                interfaceType?: string | undefined;
                implementation?: string | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
        golangBindings?: ({
            interfaceType?: string | undefined;
            implementation?: string | undefined;
        }[] & ({
            interfaceType?: string | undefined;
            implementation?: string | undefined;
        } & {
            interfaceType?: string | undefined;
            implementation?: string | undefined;
        } & Record<Exclude<keyof I["golangBindings"][number], keyof GolangBinding>, never>)[] & Record<Exclude<keyof I["golangBindings"], keyof {
            interfaceType?: string | undefined;
            implementation?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Config>, never>>(object: I): Config;
};
export declare const ModuleConfig: {
    typeUrl: string;
    encode(message: ModuleConfig, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ModuleConfig;
    fromJSON(object: any): ModuleConfig;
    toJSON(message: ModuleConfig): unknown;
    fromPartial<I extends {
        name?: string | undefined;
        config?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        golangBindings?: {
            interfaceType?: string | undefined;
            implementation?: string | undefined;
        }[] | undefined;
    } & {
        name?: string | undefined;
        config?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["config"], keyof Any>, never>) | undefined;
        golangBindings?: ({
            interfaceType?: string | undefined;
            implementation?: string | undefined;
        }[] & ({
            interfaceType?: string | undefined;
            implementation?: string | undefined;
        } & {
            interfaceType?: string | undefined;
            implementation?: string | undefined;
        } & Record<Exclude<keyof I["golangBindings"][number], keyof GolangBinding>, never>)[] & Record<Exclude<keyof I["golangBindings"], keyof {
            interfaceType?: string | undefined;
            implementation?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ModuleConfig>, never>>(object: I): ModuleConfig;
};
export declare const GolangBinding: {
    typeUrl: string;
    encode(message: GolangBinding, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GolangBinding;
    fromJSON(object: any): GolangBinding;
    toJSON(message: GolangBinding): unknown;
    fromPartial<I extends {
        interfaceType?: string | undefined;
        implementation?: string | undefined;
    } & {
        interfaceType?: string | undefined;
        implementation?: string | undefined;
    } & Record<Exclude<keyof I, keyof GolangBinding>, never>>(object: I): GolangBinding;
};
