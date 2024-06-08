import { AccessConfig } from "./types";
import { Coin } from "../../../cosmos/base/v1beta1/coin";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmwasm.wasm.v1";
/** MsgStoreCode submit Wasm code to the system */
export interface MsgStoreCode {
    /** Sender is the actor that signed the messages */
    sender: string;
    /** WASMByteCode can be raw or gzip compressed */
    wasmByteCode: Uint8Array;
    /**
     * InstantiatePermission access control to apply on contract creation,
     * optional
     */
    instantiatePermission?: AccessConfig;
}
/** MsgStoreCodeResponse returns store result data. */
export interface MsgStoreCodeResponse {
    /** CodeID is the reference to the stored WASM code */
    codeId: bigint;
    /** Checksum is the sha256 hash of the stored code */
    checksum: Uint8Array;
}
/**
 * MsgInstantiateContract create a new smart contract instance for the given
 * code id.
 */
export interface MsgInstantiateContract {
    /** Sender is the that actor that signed the messages */
    sender: string;
    /** Admin is an optional address that can execute migrations */
    admin: string;
    /** CodeID is the reference to the stored WASM code */
    codeId: bigint;
    /** Label is optional metadata to be stored with a contract instance. */
    label: string;
    /** Msg json encoded message to be passed to the contract on instantiation */
    msg: Uint8Array;
    /** Funds coins that are transferred to the contract on instantiation */
    funds: Coin[];
}
/**
 * MsgInstantiateContract2 create a new smart contract instance for the given
 * code id with a predicable address.
 */
export interface MsgInstantiateContract2 {
    /** Sender is the that actor that signed the messages */
    sender: string;
    /** Admin is an optional address that can execute migrations */
    admin: string;
    /** CodeID is the reference to the stored WASM code */
    codeId: bigint;
    /** Label is optional metadata to be stored with a contract instance. */
    label: string;
    /** Msg json encoded message to be passed to the contract on instantiation */
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
/** MsgInstantiateContractResponse return instantiation result data */
export interface MsgInstantiateContractResponse {
    /** Address is the bech32 address of the new contract instance. */
    address: string;
    /** Data contains bytes to returned from the contract */
    data: Uint8Array;
}
/** MsgInstantiateContract2Response return instantiation result data */
export interface MsgInstantiateContract2Response {
    /** Address is the bech32 address of the new contract instance. */
    address: string;
    /** Data contains bytes to returned from the contract */
    data: Uint8Array;
}
/** MsgExecuteContract submits the given message data to a smart contract */
export interface MsgExecuteContract {
    /** Sender is the that actor that signed the messages */
    sender: string;
    /** Contract is the address of the smart contract */
    contract: string;
    /** Msg json encoded message to be passed to the contract */
    msg: Uint8Array;
    /** Funds coins that are transferred to the contract on execution */
    funds: Coin[];
}
/** MsgExecuteContractResponse returns execution result data. */
export interface MsgExecuteContractResponse {
    /** Data contains bytes to returned from the contract */
    data: Uint8Array;
}
/** MsgMigrateContract runs a code upgrade/ downgrade for a smart contract */
export interface MsgMigrateContract {
    /** Sender is the that actor that signed the messages */
    sender: string;
    /** Contract is the address of the smart contract */
    contract: string;
    /** CodeID references the new WASM code */
    codeId: bigint;
    /** Msg json encoded message to be passed to the contract on migration */
    msg: Uint8Array;
}
/** MsgMigrateContractResponse returns contract migration result data. */
export interface MsgMigrateContractResponse {
    /**
     * Data contains same raw bytes returned as data from the wasm contract.
     * (May be empty)
     */
    data: Uint8Array;
}
/** MsgUpdateAdmin sets a new admin for a smart contract */
export interface MsgUpdateAdmin {
    /** Sender is the that actor that signed the messages */
    sender: string;
    /** NewAdmin address to be set */
    newAdmin: string;
    /** Contract is the address of the smart contract */
    contract: string;
}
/** MsgUpdateAdminResponse returns empty data */
export interface MsgUpdateAdminResponse {
}
/** MsgClearAdmin removes any admin stored for a smart contract */
export interface MsgClearAdmin {
    /** Sender is the actor that signed the messages */
    sender: string;
    /** Contract is the address of the smart contract */
    contract: string;
}
/** MsgClearAdminResponse returns empty data */
export interface MsgClearAdminResponse {
}
/** MsgUpdateInstantiateConfig updates instantiate config for a smart contract */
export interface MsgUpdateInstantiateConfig {
    /** Sender is the that actor that signed the messages */
    sender: string;
    /** CodeID references the stored WASM code */
    codeId: bigint;
    /** NewInstantiatePermission is the new access control */
    newInstantiatePermission?: AccessConfig;
}
/** MsgUpdateInstantiateConfigResponse returns empty data */
export interface MsgUpdateInstantiateConfigResponse {
}
export declare const MsgStoreCode: {
    typeUrl: string;
    encode(message: MsgStoreCode, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgStoreCode;
    fromJSON(object: any): MsgStoreCode;
    toJSON(message: MsgStoreCode): unknown;
    fromPartial<I extends {
        sender?: string | undefined;
        wasmByteCode?: Uint8Array | undefined;
        instantiatePermission?: {
            permission?: import("./types").AccessType | undefined;
            address?: string | undefined;
            addresses?: string[] | undefined;
        } | undefined;
    } & {
        sender?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof MsgStoreCode>, never>>(object: I): MsgStoreCode;
};
export declare const MsgStoreCodeResponse: {
    typeUrl: string;
    encode(message: MsgStoreCodeResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgStoreCodeResponse;
    fromJSON(object: any): MsgStoreCodeResponse;
    toJSON(message: MsgStoreCodeResponse): unknown;
    fromPartial<I extends {
        codeId?: bigint | undefined;
        checksum?: Uint8Array | undefined;
    } & {
        codeId?: bigint | undefined;
        checksum?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof MsgStoreCodeResponse>, never>>(object: I): MsgStoreCodeResponse;
};
export declare const MsgInstantiateContract: {
    typeUrl: string;
    encode(message: MsgInstantiateContract, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgInstantiateContract;
    fromJSON(object: any): MsgInstantiateContract;
    toJSON(message: MsgInstantiateContract): unknown;
    fromPartial<I extends {
        sender?: string | undefined;
        admin?: string | undefined;
        codeId?: bigint | undefined;
        label?: string | undefined;
        msg?: Uint8Array | undefined;
        funds?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
    } & {
        sender?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof MsgInstantiateContract>, never>>(object: I): MsgInstantiateContract;
};
export declare const MsgInstantiateContract2: {
    typeUrl: string;
    encode(message: MsgInstantiateContract2, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgInstantiateContract2;
    fromJSON(object: any): MsgInstantiateContract2;
    toJSON(message: MsgInstantiateContract2): unknown;
    fromPartial<I extends {
        sender?: string | undefined;
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
        sender?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof MsgInstantiateContract2>, never>>(object: I): MsgInstantiateContract2;
};
export declare const MsgInstantiateContractResponse: {
    typeUrl: string;
    encode(message: MsgInstantiateContractResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgInstantiateContractResponse;
    fromJSON(object: any): MsgInstantiateContractResponse;
    toJSON(message: MsgInstantiateContractResponse): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        data?: Uint8Array | undefined;
    } & {
        address?: string | undefined;
        data?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof MsgInstantiateContractResponse>, never>>(object: I): MsgInstantiateContractResponse;
};
export declare const MsgInstantiateContract2Response: {
    typeUrl: string;
    encode(message: MsgInstantiateContract2Response, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgInstantiateContract2Response;
    fromJSON(object: any): MsgInstantiateContract2Response;
    toJSON(message: MsgInstantiateContract2Response): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        data?: Uint8Array | undefined;
    } & {
        address?: string | undefined;
        data?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof MsgInstantiateContract2Response>, never>>(object: I): MsgInstantiateContract2Response;
};
export declare const MsgExecuteContract: {
    typeUrl: string;
    encode(message: MsgExecuteContract, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgExecuteContract;
    fromJSON(object: any): MsgExecuteContract;
    toJSON(message: MsgExecuteContract): unknown;
    fromPartial<I extends {
        sender?: string | undefined;
        contract?: string | undefined;
        msg?: Uint8Array | undefined;
        funds?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
    } & {
        sender?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof MsgExecuteContract>, never>>(object: I): MsgExecuteContract;
};
export declare const MsgExecuteContractResponse: {
    typeUrl: string;
    encode(message: MsgExecuteContractResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgExecuteContractResponse;
    fromJSON(object: any): MsgExecuteContractResponse;
    toJSON(message: MsgExecuteContractResponse): unknown;
    fromPartial<I extends {
        data?: Uint8Array | undefined;
    } & {
        data?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, "data">, never>>(object: I): MsgExecuteContractResponse;
};
export declare const MsgMigrateContract: {
    typeUrl: string;
    encode(message: MsgMigrateContract, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgMigrateContract;
    fromJSON(object: any): MsgMigrateContract;
    toJSON(message: MsgMigrateContract): unknown;
    fromPartial<I extends {
        sender?: string | undefined;
        contract?: string | undefined;
        codeId?: bigint | undefined;
        msg?: Uint8Array | undefined;
    } & {
        sender?: string | undefined;
        contract?: string | undefined;
        codeId?: bigint | undefined;
        msg?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof MsgMigrateContract>, never>>(object: I): MsgMigrateContract;
};
export declare const MsgMigrateContractResponse: {
    typeUrl: string;
    encode(message: MsgMigrateContractResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgMigrateContractResponse;
    fromJSON(object: any): MsgMigrateContractResponse;
    toJSON(message: MsgMigrateContractResponse): unknown;
    fromPartial<I extends {
        data?: Uint8Array | undefined;
    } & {
        data?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, "data">, never>>(object: I): MsgMigrateContractResponse;
};
export declare const MsgUpdateAdmin: {
    typeUrl: string;
    encode(message: MsgUpdateAdmin, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateAdmin;
    fromJSON(object: any): MsgUpdateAdmin;
    toJSON(message: MsgUpdateAdmin): unknown;
    fromPartial<I extends {
        sender?: string | undefined;
        newAdmin?: string | undefined;
        contract?: string | undefined;
    } & {
        sender?: string | undefined;
        newAdmin?: string | undefined;
        contract?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgUpdateAdmin>, never>>(object: I): MsgUpdateAdmin;
};
export declare const MsgUpdateAdminResponse: {
    typeUrl: string;
    encode(_: MsgUpdateAdminResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateAdminResponse;
    fromJSON(_: any): MsgUpdateAdminResponse;
    toJSON(_: MsgUpdateAdminResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgUpdateAdminResponse;
};
export declare const MsgClearAdmin: {
    typeUrl: string;
    encode(message: MsgClearAdmin, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgClearAdmin;
    fromJSON(object: any): MsgClearAdmin;
    toJSON(message: MsgClearAdmin): unknown;
    fromPartial<I extends {
        sender?: string | undefined;
        contract?: string | undefined;
    } & {
        sender?: string | undefined;
        contract?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgClearAdmin>, never>>(object: I): MsgClearAdmin;
};
export declare const MsgClearAdminResponse: {
    typeUrl: string;
    encode(_: MsgClearAdminResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgClearAdminResponse;
    fromJSON(_: any): MsgClearAdminResponse;
    toJSON(_: MsgClearAdminResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgClearAdminResponse;
};
export declare const MsgUpdateInstantiateConfig: {
    typeUrl: string;
    encode(message: MsgUpdateInstantiateConfig, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateInstantiateConfig;
    fromJSON(object: any): MsgUpdateInstantiateConfig;
    toJSON(message: MsgUpdateInstantiateConfig): unknown;
    fromPartial<I extends {
        sender?: string | undefined;
        codeId?: bigint | undefined;
        newInstantiatePermission?: {
            permission?: import("./types").AccessType | undefined;
            address?: string | undefined;
            addresses?: string[] | undefined;
        } | undefined;
    } & {
        sender?: string | undefined;
        codeId?: bigint | undefined;
        newInstantiatePermission?: ({
            permission?: import("./types").AccessType | undefined;
            address?: string | undefined;
            addresses?: string[] | undefined;
        } & {
            permission?: import("./types").AccessType | undefined;
            address?: string | undefined;
            addresses?: (string[] & string[] & Record<Exclude<keyof I["newInstantiatePermission"]["addresses"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["newInstantiatePermission"], keyof AccessConfig>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof MsgUpdateInstantiateConfig>, never>>(object: I): MsgUpdateInstantiateConfig;
};
export declare const MsgUpdateInstantiateConfigResponse: {
    typeUrl: string;
    encode(_: MsgUpdateInstantiateConfigResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateInstantiateConfigResponse;
    fromJSON(_: any): MsgUpdateInstantiateConfigResponse;
    toJSON(_: MsgUpdateInstantiateConfigResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgUpdateInstantiateConfigResponse;
};
/** Msg defines the wasm Msg service. */
export interface Msg {
    /** StoreCode to submit Wasm code to the system */
    StoreCode(request: MsgStoreCode): Promise<MsgStoreCodeResponse>;
    /**
     * InstantiateContract creates a new smart contract instance for the given
     *  code id.
     */
    InstantiateContract(request: MsgInstantiateContract): Promise<MsgInstantiateContractResponse>;
    /**
     * InstantiateContract2 creates a new smart contract instance for the given
     *  code id with a predictable address
     */
    InstantiateContract2(request: MsgInstantiateContract2): Promise<MsgInstantiateContract2Response>;
    /** Execute submits the given message data to a smart contract */
    ExecuteContract(request: MsgExecuteContract): Promise<MsgExecuteContractResponse>;
    /** Migrate runs a code upgrade/ downgrade for a smart contract */
    MigrateContract(request: MsgMigrateContract): Promise<MsgMigrateContractResponse>;
    /** UpdateAdmin sets a new   admin for a smart contract */
    UpdateAdmin(request: MsgUpdateAdmin): Promise<MsgUpdateAdminResponse>;
    /** ClearAdmin removes any admin stored for a smart contract */
    ClearAdmin(request: MsgClearAdmin): Promise<MsgClearAdminResponse>;
    /** UpdateInstantiateConfig updates instantiate config for a smart contract */
    UpdateInstantiateConfig(request: MsgUpdateInstantiateConfig): Promise<MsgUpdateInstantiateConfigResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    StoreCode(request: MsgStoreCode): Promise<MsgStoreCodeResponse>;
    InstantiateContract(request: MsgInstantiateContract): Promise<MsgInstantiateContractResponse>;
    InstantiateContract2(request: MsgInstantiateContract2): Promise<MsgInstantiateContract2Response>;
    ExecuteContract(request: MsgExecuteContract): Promise<MsgExecuteContractResponse>;
    MigrateContract(request: MsgMigrateContract): Promise<MsgMigrateContractResponse>;
    UpdateAdmin(request: MsgUpdateAdmin): Promise<MsgUpdateAdminResponse>;
    ClearAdmin(request: MsgClearAdmin): Promise<MsgClearAdminResponse>;
    UpdateInstantiateConfig(request: MsgUpdateInstantiateConfig): Promise<MsgUpdateInstantiateConfigResponse>;
}
