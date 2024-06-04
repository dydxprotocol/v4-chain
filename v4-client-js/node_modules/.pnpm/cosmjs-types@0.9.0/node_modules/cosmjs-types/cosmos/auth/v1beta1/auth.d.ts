import { Any } from "../../../google/protobuf/any";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.auth.v1beta1";
/**
 * BaseAccount defines a base account type. It contains all the necessary fields
 * for basic account functionality. Any custom account type should extend this
 * type for additional functionality (e.g. vesting).
 */
export interface BaseAccount {
    address: string;
    pubKey?: Any;
    accountNumber: bigint;
    sequence: bigint;
}
/** ModuleAccount defines an account for modules that holds coins on a pool. */
export interface ModuleAccount {
    baseAccount?: BaseAccount;
    name: string;
    permissions: string[];
}
/**
 * ModuleCredential represents a unclaimable pubkey for base accounts controlled by modules.
 *
 * Since: cosmos-sdk 0.47
 */
export interface ModuleCredential {
    /** module_name is the name of the module used for address derivation (passed into address.Module). */
    moduleName: string;
    /**
     * derivation_keys is for deriving a module account address (passed into address.Module)
     * adding more keys creates sub-account addresses (passed into address.Derive)
     */
    derivationKeys: Uint8Array[];
}
/** Params defines the parameters for the auth module. */
export interface Params {
    maxMemoCharacters: bigint;
    txSigLimit: bigint;
    txSizeCostPerByte: bigint;
    sigVerifyCostEd25519: bigint;
    sigVerifyCostSecp256k1: bigint;
}
export declare const BaseAccount: {
    typeUrl: string;
    encode(message: BaseAccount, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): BaseAccount;
    fromJSON(object: any): BaseAccount;
    toJSON(message: BaseAccount): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        pubKey?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        accountNumber?: bigint | undefined;
        sequence?: bigint | undefined;
    } & {
        address?: string | undefined;
        pubKey?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["pubKey"], keyof Any>, never>) | undefined;
        accountNumber?: bigint | undefined;
        sequence?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof BaseAccount>, never>>(object: I): BaseAccount;
};
export declare const ModuleAccount: {
    typeUrl: string;
    encode(message: ModuleAccount, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ModuleAccount;
    fromJSON(object: any): ModuleAccount;
    toJSON(message: ModuleAccount): unknown;
    fromPartial<I extends {
        baseAccount?: {
            address?: string | undefined;
            pubKey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            accountNumber?: bigint | undefined;
            sequence?: bigint | undefined;
        } | undefined;
        name?: string | undefined;
        permissions?: string[] | undefined;
    } & {
        baseAccount?: ({
            address?: string | undefined;
            pubKey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            accountNumber?: bigint | undefined;
            sequence?: bigint | undefined;
        } & {
            address?: string | undefined;
            pubKey?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["baseAccount"]["pubKey"], keyof Any>, never>) | undefined;
            accountNumber?: bigint | undefined;
            sequence?: bigint | undefined;
        } & Record<Exclude<keyof I["baseAccount"], keyof BaseAccount>, never>) | undefined;
        name?: string | undefined;
        permissions?: (string[] & string[] & Record<Exclude<keyof I["permissions"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ModuleAccount>, never>>(object: I): ModuleAccount;
};
export declare const ModuleCredential: {
    typeUrl: string;
    encode(message: ModuleCredential, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ModuleCredential;
    fromJSON(object: any): ModuleCredential;
    toJSON(message: ModuleCredential): unknown;
    fromPartial<I extends {
        moduleName?: string | undefined;
        derivationKeys?: Uint8Array[] | undefined;
    } & {
        moduleName?: string | undefined;
        derivationKeys?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["derivationKeys"], keyof Uint8Array[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ModuleCredential>, never>>(object: I): ModuleCredential;
};
export declare const Params: {
    typeUrl: string;
    encode(message: Params, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Params;
    fromJSON(object: any): Params;
    toJSON(message: Params): unknown;
    fromPartial<I extends {
        maxMemoCharacters?: bigint | undefined;
        txSigLimit?: bigint | undefined;
        txSizeCostPerByte?: bigint | undefined;
        sigVerifyCostEd25519?: bigint | undefined;
        sigVerifyCostSecp256k1?: bigint | undefined;
    } & {
        maxMemoCharacters?: bigint | undefined;
        txSigLimit?: bigint | undefined;
        txSizeCostPerByte?: bigint | undefined;
        sigVerifyCostEd25519?: bigint | undefined;
        sigVerifyCostSecp256k1?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof Params>, never>>(object: I): Params;
};
