import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "cosmos.bank.module.v1";
/** Module is the config object of the bank module. */
export interface Module {
    /**
     * blocked_module_accounts configures exceptional module accounts which should be blocked from receiving funds.
     * If left empty it defaults to the list of account names supplied in the auth module configuration as
     * module_account_permissions
     */
    blockedModuleAccountsOverride: string[];
    /** authority defines the custom module authority. If not set, defaults to the governance module. */
    authority: string;
}
export declare const Module: {
    typeUrl: string;
    encode(message: Module, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Module;
    fromJSON(object: any): Module;
    toJSON(message: Module): unknown;
    fromPartial<I extends {
        blockedModuleAccountsOverride?: string[] | undefined;
        authority?: string | undefined;
    } & {
        blockedModuleAccountsOverride?: (string[] & string[] & Record<Exclude<keyof I["blockedModuleAccountsOverride"], keyof string[]>, never>) | undefined;
        authority?: string | undefined;
    } & Record<Exclude<keyof I, keyof Module>, never>>(object: I): Module;
};
