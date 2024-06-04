import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "cosmos.auth.module.v1";
/** Module is the config object for the auth module. */
export interface Module {
    /** bech32_prefix is the bech32 account prefix for the app. */
    bech32Prefix: string;
    /** module_account_permissions are module account permissions. */
    moduleAccountPermissions: ModuleAccountPermission[];
    /** authority defines the custom module authority. If not set, defaults to the governance module. */
    authority: string;
}
/** ModuleAccountPermission represents permissions for a module account. */
export interface ModuleAccountPermission {
    /** account is the name of the module. */
    account: string;
    /**
     * permissions are the permissions this module has. Currently recognized
     * values are minter, burner and staking.
     */
    permissions: string[];
}
export declare const Module: {
    typeUrl: string;
    encode(message: Module, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Module;
    fromJSON(object: any): Module;
    toJSON(message: Module): unknown;
    fromPartial<I extends {
        bech32Prefix?: string | undefined;
        moduleAccountPermissions?: {
            account?: string | undefined;
            permissions?: string[] | undefined;
        }[] | undefined;
        authority?: string | undefined;
    } & {
        bech32Prefix?: string | undefined;
        moduleAccountPermissions?: ({
            account?: string | undefined;
            permissions?: string[] | undefined;
        }[] & ({
            account?: string | undefined;
            permissions?: string[] | undefined;
        } & {
            account?: string | undefined;
            permissions?: (string[] & string[] & Record<Exclude<keyof I["moduleAccountPermissions"][number]["permissions"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["moduleAccountPermissions"][number], keyof ModuleAccountPermission>, never>)[] & Record<Exclude<keyof I["moduleAccountPermissions"], keyof {
            account?: string | undefined;
            permissions?: string[] | undefined;
        }[]>, never>) | undefined;
        authority?: string | undefined;
    } & Record<Exclude<keyof I, keyof Module>, never>>(object: I): Module;
};
export declare const ModuleAccountPermission: {
    typeUrl: string;
    encode(message: ModuleAccountPermission, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ModuleAccountPermission;
    fromJSON(object: any): ModuleAccountPermission;
    toJSON(message: ModuleAccountPermission): unknown;
    fromPartial<I extends {
        account?: string | undefined;
        permissions?: string[] | undefined;
    } & {
        account?: string | undefined;
        permissions?: (string[] & string[] & Record<Exclude<keyof I["permissions"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ModuleAccountPermission>, never>>(object: I): ModuleAccountPermission;
};
