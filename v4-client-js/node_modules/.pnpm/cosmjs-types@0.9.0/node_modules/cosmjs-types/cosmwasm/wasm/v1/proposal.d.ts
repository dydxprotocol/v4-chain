import { AccessConfig } from "./types";
import { Coin } from "../../../cosmos/base/v1beta1/coin";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmwasm.wasm.v1";
/** StoreCodeProposal gov proposal content type to submit WASM code to the system */
export interface StoreCodeProposal {
    /** Title is a short summary */
    title: string;
    /** Description is a human readable text */
    description: string;
    /** RunAs is the address that is passed to the contract's environment as sender */
    runAs: string;
    /** WASMByteCode can be raw or gzip compressed */
    wasmByteCode: Uint8Array;
    /** InstantiatePermission to apply on contract creation, optional */
    instantiatePermission?: AccessConfig;
    /** UnpinCode code on upload, optional */
    unpinCode: boolean;
    /** Source is the URL where the code is hosted */
    source: string;
    /**
     * Builder is the docker image used to build the code deterministically, used
     * for smart contract verification
     */
    builder: string;
    /**
     * CodeHash is the SHA256 sum of the code outputted by builder, used for smart
     * contract verification
     */
    codeHash: Uint8Array;
}
/**
 * InstantiateContractProposal gov proposal content type to instantiate a
 * contract.
 */
export interface InstantiateContractProposal {
    /** Title is a short summary */
    title: string;
    /** Description is a human readable text */
    description: string;
    /** RunAs is the address that is passed to the contract's environment as sender */
    runAs: string;
    /** Admin is an optional address that can execute migrations */
    admin: string;
    /** CodeID is the reference to the stored WASM code */
    codeId: bigint;
    /** Label is optional metadata to be stored with a constract instance. */
    label: string;
    /** Msg json encoded message to be passed to the contract on instantiation */
    msg: Uint8Array;
    /** Funds coins that are transferred to the contract on instantiation */
    funds: Coin[];
}
/**
 * InstantiateContract2Proposal gov proposal content type to instantiate
 * contract 2
 */
export interface InstantiateContract2Proposal {
    /** Title is a short summary */
    title: string;
    /** Description is a human readable text */
    description: string;
    /** RunAs is the address that is passed to the contract's enviroment as sender */
    runAs: string;
    /** Admin is an optional address that can execute migrations */
    admin: string;
    /** CodeID is the reference to the stored WASM code */
    codeId: bigint;
    /** Label is optional metadata to be stored with a constract instance. */
    label: string;
    /** Msg json encode message to be passed to the contract on instantiation */
    msg: Uint8Array;
    /** Funds coins that are transferred to the contract on instantiation */
    funds: Coin[];
    /** Salt is an arbitrary value provided by the sender. Size can be 1 to 64. */
    salt: Uint8Array;
    /**
     * FixMsg include the msg value into the hash for the predictable address.
     * Default is false
     */
    fixMsg: boolean;
}
/** MigrateContractProposal gov proposal content type to migrate a contract. */
export interface MigrateContractProposal {
    /** Title is a short summary */
    title: string;
    /** Description is a human readable text */
    description: string;
    /** Contract is the address of the smart contract */
    contract: string;
    /** CodeID references the new WASM code */
    codeId: bigint;
    /** Msg json encoded message to be passed to the contract on migration */
    msg: Uint8Array;
}
/** SudoContractProposal gov proposal content type to call sudo on a contract. */
export interface SudoContractProposal {
    /** Title is a short summary */
    title: string;
    /** Description is a human readable text */
    description: string;
    /** Contract is the address of the smart contract */
    contract: string;
    /** Msg json encoded message to be passed to the contract as sudo */
    msg: Uint8Array;
}
/**
 * ExecuteContractProposal gov proposal content type to call execute on a
 * contract.
 */
export interface ExecuteContractProposal {
    /** Title is a short summary */
    title: string;
    /** Description is a human readable text */
    description: string;
    /** RunAs is the address that is passed to the contract's environment as sender */
    runAs: string;
    /** Contract is the address of the smart contract */
    contract: string;
    /** Msg json encoded message to be passed to the contract as execute */
    msg: Uint8Array;
    /** Funds coins that are transferred to the contract on instantiation */
    funds: Coin[];
}
/** UpdateAdminProposal gov proposal content type to set an admin for a contract. */
export interface UpdateAdminProposal {
    /** Title is a short summary */
    title: string;
    /** Description is a human readable text */
    description: string;
    /** NewAdmin address to be set */
    newAdmin: string;
    /** Contract is the address of the smart contract */
    contract: string;
}
/**
 * ClearAdminProposal gov proposal content type to clear the admin of a
 * contract.
 */
export interface ClearAdminProposal {
    /** Title is a short summary */
    title: string;
    /** Description is a human readable text */
    description: string;
    /** Contract is the address of the smart contract */
    contract: string;
}
/**
 * PinCodesProposal gov proposal content type to pin a set of code ids in the
 * wasmvm cache.
 */
export interface PinCodesProposal {
    /** Title is a short summary */
    title: string;
    /** Description is a human readable text */
    description: string;
    /** CodeIDs references the new WASM codes */
    codeIds: bigint[];
}
/**
 * UnpinCodesProposal gov proposal content type to unpin a set of code ids in
 * the wasmvm cache.
 */
export interface UnpinCodesProposal {
    /** Title is a short summary */
    title: string;
    /** Description is a human readable text */
    description: string;
    /** CodeIDs references the WASM codes */
    codeIds: bigint[];
}
/**
 * AccessConfigUpdate contains the code id and the access config to be
 * applied.
 */
export interface AccessConfigUpdate {
    /** CodeID is the reference to the stored WASM code to be updated */
    codeId: bigint;
    /** InstantiatePermission to apply to the set of code ids */
    instantiatePermission: AccessConfig;
}
/**
 * UpdateInstantiateConfigProposal gov proposal content type to update
 * instantiate config to a  set of code ids.
 */
export interface UpdateInstantiateConfigProposal {
    /** Title is a short summary */
    title: string;
    /** Description is a human readable text */
    description: string;
    /**
     * AccessConfigUpdate contains the list of code ids and the access config
     * to be applied.
     */
    accessConfigUpdates: AccessConfigUpdate[];
}
/**
 * StoreAndInstantiateContractProposal gov proposal content type to store
 * and instantiate the contract.
 */
export interface StoreAndInstantiateContractProposal {
    /** Title is a short summary */
    title: string;
    /** Description is a human readable text */
    description: string;
    /** RunAs is the address that is passed to the contract's environment as sender */
    runAs: string;
    /** WASMByteCode can be raw or gzip compressed */
    wasmByteCode: Uint8Array;
    /** InstantiatePermission to apply on contract creation, optional */
    instantiatePermission?: AccessConfig;
    /** UnpinCode code on upload, optional */
    unpinCode: boolean;
    /** Admin is an optional address that can execute migrations */
    admin: string;
    /** Label is optional metadata to be stored with a constract instance. */
    label: string;
    /** Msg json encoded message to be passed to the contract on instantiation */
    msg: Uint8Array;
    /** Funds coins that are transferred to the contract on instantiation */
    funds: Coin[];
    /** Source is the URL where the code is hosted */
    source: string;
    /**
     * Builder is the docker image used to build the code deterministically, used
     * for smart contract verification
     */
    builder: string;
    /**
     * CodeHash is the SHA256 sum of the code outputted by builder, used for smart
     * contract verification
     */
    codeHash: Uint8Array;
}
export declare const StoreCodeProposal: {
    typeUrl: string;
    encode(message: StoreCodeProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): StoreCodeProposal;
    fromJSON(object: any): StoreCodeProposal;
    toJSON(message: StoreCodeProposal): unknown;
    fromPartial<I extends {
        title?: string | undefined;
        description?: string | undefined;
        runAs?: string | undefined;
        wasmByteCode?: Uint8Array | undefined;
        instantiatePermission?: {
            permission?: import("./types").AccessType | undefined;
            address?: string | undefined;
            addresses?: string[] | undefined;
        } | undefined;
        unpinCode?: boolean | undefined;
        source?: string | undefined;
        builder?: string | undefined;
        codeHash?: Uint8Array | undefined;
    } & {
        title?: string | undefined;
        description?: string | undefined;
        runAs?: string | undefined;
        wasmByteCode?: Uint8Array | undefined;
        instantiatePermission?: ({
            permission?: import("./types").AccessType | undefined;
            address?: string | undefined;
            addresses?: string[] | undefined;
        } & {
            permission?: import("./types").AccessType | undefined;
            address?: string | undefined;
            addresses?: (string[] & string[] & Record<Exclude<keyof I["instantiatePermission"]["addresses"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["instantiatePermission"], keyof AccessConfig>, never>) | undefined;
        unpinCode?: boolean | undefined;
        source?: string | undefined;
        builder?: string | undefined;
        codeHash?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof StoreCodeProposal>, never>>(object: I): StoreCodeProposal;
};
export declare const InstantiateContractProposal: {
    typeUrl: string;
    encode(message: InstantiateContractProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): InstantiateContractProposal;
    fromJSON(object: any): InstantiateContractProposal;
    toJSON(message: InstantiateContractProposal): unknown;
    fromPartial<I extends {
        title?: string | undefined;
        description?: string | undefined;
        runAs?: string | undefined;
        admin?: string | undefined;
        codeId?: bigint | undefined;
        label?: string | undefined;
        msg?: Uint8Array | undefined;
        funds?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
    } & {
        title?: string | undefined;
        description?: string | undefined;
        runAs?: string | undefined;
        admin?: string | undefined;
        codeId?: bigint | undefined;
        label?: string | undefined;
        msg?: Uint8Array | undefined;
        funds?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["funds"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["funds"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof InstantiateContractProposal>, never>>(object: I): InstantiateContractProposal;
};
export declare const InstantiateContract2Proposal: {
    typeUrl: string;
    encode(message: InstantiateContract2Proposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): InstantiateContract2Proposal;
    fromJSON(object: any): InstantiateContract2Proposal;
    toJSON(message: InstantiateContract2Proposal): unknown;
    fromPartial<I extends {
        title?: string | undefined;
        description?: string | undefined;
        runAs?: string | undefined;
        admin?: string | undefined;
        codeId?: bigint | undefined;
        label?: string | undefined;
        msg?: Uint8Array | undefined;
        funds?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
        salt?: Uint8Array | undefined;
        fixMsg?: boolean | undefined;
    } & {
        title?: string | undefined;
        description?: string | undefined;
        runAs?: string | undefined;
        admin?: string | undefined;
        codeId?: bigint | undefined;
        label?: string | undefined;
        msg?: Uint8Array | undefined;
        funds?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["funds"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["funds"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
        salt?: Uint8Array | undefined;
        fixMsg?: boolean | undefined;
    } & Record<Exclude<keyof I, keyof InstantiateContract2Proposal>, never>>(object: I): InstantiateContract2Proposal;
};
export declare const MigrateContractProposal: {
    typeUrl: string;
    encode(message: MigrateContractProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MigrateContractProposal;
    fromJSON(object: any): MigrateContractProposal;
    toJSON(message: MigrateContractProposal): unknown;
    fromPartial<I extends {
        title?: string | undefined;
        description?: string | undefined;
        contract?: string | undefined;
        codeId?: bigint | undefined;
        msg?: Uint8Array | undefined;
    } & {
        title?: string | undefined;
        description?: string | undefined;
        contract?: string | undefined;
        codeId?: bigint | undefined;
        msg?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof MigrateContractProposal>, never>>(object: I): MigrateContractProposal;
};
export declare const SudoContractProposal: {
    typeUrl: string;
    encode(message: SudoContractProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): SudoContractProposal;
    fromJSON(object: any): SudoContractProposal;
    toJSON(message: SudoContractProposal): unknown;
    fromPartial<I extends {
        title?: string | undefined;
        description?: string | undefined;
        contract?: string | undefined;
        msg?: Uint8Array | undefined;
    } & {
        title?: string | undefined;
        description?: string | undefined;
        contract?: string | undefined;
        msg?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof SudoContractProposal>, never>>(object: I): SudoContractProposal;
};
export declare const ExecuteContractProposal: {
    typeUrl: string;
    encode(message: ExecuteContractProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ExecuteContractProposal;
    fromJSON(object: any): ExecuteContractProposal;
    toJSON(message: ExecuteContractProposal): unknown;
    fromPartial<I extends {
        title?: string | undefined;
        description?: string | undefined;
        runAs?: string | undefined;
        contract?: string | undefined;
        msg?: Uint8Array | undefined;
        funds?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
    } & {
        title?: string | undefined;
        description?: string | undefined;
        runAs?: string | undefined;
        contract?: string | undefined;
        msg?: Uint8Array | undefined;
        funds?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["funds"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["funds"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ExecuteContractProposal>, never>>(object: I): ExecuteContractProposal;
};
export declare const UpdateAdminProposal: {
    typeUrl: string;
    encode(message: UpdateAdminProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): UpdateAdminProposal;
    fromJSON(object: any): UpdateAdminProposal;
    toJSON(message: UpdateAdminProposal): unknown;
    fromPartial<I extends {
        title?: string | undefined;
        description?: string | undefined;
        newAdmin?: string | undefined;
        contract?: string | undefined;
    } & {
        title?: string | undefined;
        description?: string | undefined;
        newAdmin?: string | undefined;
        contract?: string | undefined;
    } & Record<Exclude<keyof I, keyof UpdateAdminProposal>, never>>(object: I): UpdateAdminProposal;
};
export declare const ClearAdminProposal: {
    typeUrl: string;
    encode(message: ClearAdminProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ClearAdminProposal;
    fromJSON(object: any): ClearAdminProposal;
    toJSON(message: ClearAdminProposal): unknown;
    fromPartial<I extends {
        title?: string | undefined;
        description?: string | undefined;
        contract?: string | undefined;
    } & {
        title?: string | undefined;
        description?: string | undefined;
        contract?: string | undefined;
    } & Record<Exclude<keyof I, keyof ClearAdminProposal>, never>>(object: I): ClearAdminProposal;
};
export declare const PinCodesProposal: {
    typeUrl: string;
    encode(message: PinCodesProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): PinCodesProposal;
    fromJSON(object: any): PinCodesProposal;
    toJSON(message: PinCodesProposal): unknown;
    fromPartial<I extends {
        title?: string | undefined;
        description?: string | undefined;
        codeIds?: bigint[] | undefined;
    } & {
        title?: string | undefined;
        description?: string | undefined;
        codeIds?: (bigint[] & bigint[] & Record<Exclude<keyof I["codeIds"], keyof bigint[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof PinCodesProposal>, never>>(object: I): PinCodesProposal;
};
export declare const UnpinCodesProposal: {
    typeUrl: string;
    encode(message: UnpinCodesProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): UnpinCodesProposal;
    fromJSON(object: any): UnpinCodesProposal;
    toJSON(message: UnpinCodesProposal): unknown;
    fromPartial<I extends {
        title?: string | undefined;
        description?: string | undefined;
        codeIds?: bigint[] | undefined;
    } & {
        title?: string | undefined;
        description?: string | undefined;
        codeIds?: (bigint[] & bigint[] & Record<Exclude<keyof I["codeIds"], keyof bigint[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof UnpinCodesProposal>, never>>(object: I): UnpinCodesProposal;
};
export declare const AccessConfigUpdate: {
    typeUrl: string;
    encode(message: AccessConfigUpdate, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): AccessConfigUpdate;
    fromJSON(object: any): AccessConfigUpdate;
    toJSON(message: AccessConfigUpdate): unknown;
    fromPartial<I extends {
        codeId?: bigint | undefined;
        instantiatePermission?: {
            permission?: import("./types").AccessType | undefined;
            address?: string | undefined;
            addresses?: string[] | undefined;
        } | undefined;
    } & {
        codeId?: bigint | undefined;
        instantiatePermission?: ({
            permission?: import("./types").AccessType | undefined;
            address?: string | undefined;
            addresses?: string[] | undefined;
        } & {
            permission?: import("./types").AccessType | undefined;
            address?: string | undefined;
            addresses?: (string[] & string[] & Record<Exclude<keyof I["instantiatePermission"]["addresses"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["instantiatePermission"], keyof AccessConfig>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof AccessConfigUpdate>, never>>(object: I): AccessConfigUpdate;
};
export declare const UpdateInstantiateConfigProposal: {
    typeUrl: string;
    encode(message: UpdateInstantiateConfigProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): UpdateInstantiateConfigProposal;
    fromJSON(object: any): UpdateInstantiateConfigProposal;
    toJSON(message: UpdateInstantiateConfigProposal): unknown;
    fromPartial<I extends {
        title?: string | undefined;
        description?: string | undefined;
        accessConfigUpdates?: {
            codeId?: bigint | undefined;
            instantiatePermission?: {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } | undefined;
        }[] | undefined;
    } & {
        title?: string | undefined;
        description?: string | undefined;
        accessConfigUpdates?: ({
            codeId?: bigint | undefined;
            instantiatePermission?: {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } | undefined;
        }[] & ({
            codeId?: bigint | undefined;
            instantiatePermission?: {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } | undefined;
        } & {
            codeId?: bigint | undefined;
            instantiatePermission?: ({
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } & {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: (string[] & string[] & Record<Exclude<keyof I["accessConfigUpdates"][number]["instantiatePermission"]["addresses"], keyof string[]>, never>) | undefined;
            } & Record<Exclude<keyof I["accessConfigUpdates"][number]["instantiatePermission"], keyof AccessConfig>, never>) | undefined;
        } & Record<Exclude<keyof I["accessConfigUpdates"][number], keyof AccessConfigUpdate>, never>)[] & Record<Exclude<keyof I["accessConfigUpdates"], keyof {
            codeId?: bigint | undefined;
            instantiatePermission?: {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof UpdateInstantiateConfigProposal>, never>>(object: I): UpdateInstantiateConfigProposal;
};
export declare const StoreAndInstantiateContractProposal: {
    typeUrl: string;
    encode(message: StoreAndInstantiateContractProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): StoreAndInstantiateContractProposal;
    fromJSON(object: any): StoreAndInstantiateContractProposal;
    toJSON(message: StoreAndInstantiateContractProposal): unknown;
    fromPartial<I extends {
        title?: string | undefined;
        description?: string | undefined;
        runAs?: string | undefined;
        wasmByteCode?: Uint8Array | undefined;
        instantiatePermission?: {
            permission?: import("./types").AccessType | undefined;
            address?: string | undefined;
            addresses?: string[] | undefined;
        } | undefined;
        unpinCode?: boolean | undefined;
        admin?: string | undefined;
        label?: string | undefined;
        msg?: Uint8Array | undefined;
        funds?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
        source?: string | undefined;
        builder?: string | undefined;
        codeHash?: Uint8Array | undefined;
    } & {
        title?: string | undefined;
        description?: string | undefined;
        runAs?: string | undefined;
        wasmByteCode?: Uint8Array | undefined;
        instantiatePermission?: ({
            permission?: import("./types").AccessType | undefined;
            address?: string | undefined;
            addresses?: string[] | undefined;
        } & {
            permission?: import("./types").AccessType | undefined;
            address?: string | undefined;
            addresses?: (string[] & string[] & Record<Exclude<keyof I["instantiatePermission"]["addresses"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["instantiatePermission"], keyof AccessConfig>, never>) | undefined;
        unpinCode?: boolean | undefined;
        admin?: string | undefined;
        label?: string | undefined;
        msg?: Uint8Array | undefined;
        funds?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["funds"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["funds"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
        source?: string | undefined;
        builder?: string | undefined;
        codeHash?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof StoreAndInstantiateContractProposal>, never>>(object: I): StoreAndInstantiateContractProposal;
};
