import { Any } from "../../../google/protobuf/any";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmwasm.wasm.v1";
/** AccessType permission types */
export declare enum AccessType {
    /** ACCESS_TYPE_UNSPECIFIED - AccessTypeUnspecified placeholder for empty value */
    ACCESS_TYPE_UNSPECIFIED = 0,
    /** ACCESS_TYPE_NOBODY - AccessTypeNobody forbidden */
    ACCESS_TYPE_NOBODY = 1,
    /**
     * ACCESS_TYPE_ONLY_ADDRESS - AccessTypeOnlyAddress restricted to a single address
     * Deprecated: use AccessTypeAnyOfAddresses instead
     */
    ACCESS_TYPE_ONLY_ADDRESS = 2,
    /** ACCESS_TYPE_EVERYBODY - AccessTypeEverybody unrestricted */
    ACCESS_TYPE_EVERYBODY = 3,
    /** ACCESS_TYPE_ANY_OF_ADDRESSES - AccessTypeAnyOfAddresses allow any of the addresses */
    ACCESS_TYPE_ANY_OF_ADDRESSES = 4,
    UNRECOGNIZED = -1
}
export declare function accessTypeFromJSON(object: any): AccessType;
export declare function accessTypeToJSON(object: AccessType): string;
/** ContractCodeHistoryOperationType actions that caused a code change */
export declare enum ContractCodeHistoryOperationType {
    /** CONTRACT_CODE_HISTORY_OPERATION_TYPE_UNSPECIFIED - ContractCodeHistoryOperationTypeUnspecified placeholder for empty value */
    CONTRACT_CODE_HISTORY_OPERATION_TYPE_UNSPECIFIED = 0,
    /** CONTRACT_CODE_HISTORY_OPERATION_TYPE_INIT - ContractCodeHistoryOperationTypeInit on chain contract instantiation */
    CONTRACT_CODE_HISTORY_OPERATION_TYPE_INIT = 1,
    /** CONTRACT_CODE_HISTORY_OPERATION_TYPE_MIGRATE - ContractCodeHistoryOperationTypeMigrate code migration */
    CONTRACT_CODE_HISTORY_OPERATION_TYPE_MIGRATE = 2,
    /** CONTRACT_CODE_HISTORY_OPERATION_TYPE_GENESIS - ContractCodeHistoryOperationTypeGenesis based on genesis data */
    CONTRACT_CODE_HISTORY_OPERATION_TYPE_GENESIS = 3,
    UNRECOGNIZED = -1
}
export declare function contractCodeHistoryOperationTypeFromJSON(object: any): ContractCodeHistoryOperationType;
export declare function contractCodeHistoryOperationTypeToJSON(object: ContractCodeHistoryOperationType): string;
/** AccessTypeParam */
export interface AccessTypeParam {
    value: AccessType;
}
/** AccessConfig access control type. */
export interface AccessConfig {
    permission: AccessType;
    /**
     * Address
     * Deprecated: replaced by addresses
     */
    address: string;
    addresses: string[];
}
/** Params defines the set of wasm parameters. */
export interface Params {
    codeUploadAccess: AccessConfig;
    instantiateDefaultPermission: AccessType;
}
/** CodeInfo is data for the uploaded contract WASM code */
export interface CodeInfo {
    /** CodeHash is the unique identifier created by wasmvm */
    codeHash: Uint8Array;
    /** Creator address who initially stored the code */
    creator: string;
    /** InstantiateConfig access control to apply on contract creation, optional */
    instantiateConfig: AccessConfig;
}
/** ContractInfo stores a WASM contract instance */
export interface ContractInfo {
    /** CodeID is the reference to the stored Wasm code */
    codeId: bigint;
    /** Creator address who initially instantiated the contract */
    creator: string;
    /** Admin is an optional address that can execute migrations */
    admin: string;
    /** Label is optional metadata to be stored with a contract instance. */
    label: string;
    /** Created Tx position when the contract was instantiated. */
    created?: AbsoluteTxPosition;
    ibcPortId: string;
    /**
     * Extension is an extension point to store custom metadata within the
     * persistence model.
     */
    extension?: Any;
}
/** ContractCodeHistoryEntry metadata to a contract. */
export interface ContractCodeHistoryEntry {
    operation: ContractCodeHistoryOperationType;
    /** CodeID is the reference to the stored WASM code */
    codeId: bigint;
    /** Updated Tx position when the operation was executed. */
    updated?: AbsoluteTxPosition;
    msg: Uint8Array;
}
/**
 * AbsoluteTxPosition is a unique transaction position that allows for global
 * ordering of transactions.
 */
export interface AbsoluteTxPosition {
    /** BlockHeight is the block the contract was created at */
    blockHeight: bigint;
    /**
     * TxIndex is a monotonic counter within the block (actual transaction index,
     * or gas consumed)
     */
    txIndex: bigint;
}
/** Model is a struct that holds a KV pair */
export interface Model {
    /** hex-encode key to read it better (this is often ascii) */
    key: Uint8Array;
    /** base64-encode raw value */
    value: Uint8Array;
}
export declare const AccessTypeParam: {
    typeUrl: string;
    encode(message: AccessTypeParam, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): AccessTypeParam;
    fromJSON(object: any): AccessTypeParam;
    toJSON(message: AccessTypeParam): unknown;
    fromPartial<I extends {
        value?: AccessType | undefined;
    } & {
        value?: AccessType | undefined;
    } & Record<Exclude<keyof I, "value">, never>>(object: I): AccessTypeParam;
};
export declare const AccessConfig: {
    typeUrl: string;
    encode(message: AccessConfig, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): AccessConfig;
    fromJSON(object: any): AccessConfig;
    toJSON(message: AccessConfig): unknown;
    fromPartial<I extends {
        permission?: AccessType | undefined;
        address?: string | undefined;
        addresses?: string[] | undefined;
    } & {
        permission?: AccessType | undefined;
        address?: string | undefined;
        addresses?: (string[] & string[] & Record<Exclude<keyof I["addresses"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof AccessConfig>, never>>(object: I): AccessConfig;
};
export declare const Params: {
    typeUrl: string;
    encode(message: Params, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Params;
    fromJSON(object: any): Params;
    toJSON(message: Params): unknown;
    fromPartial<I extends {
        codeUploadAccess?: {
            permission?: AccessType | undefined;
            address?: string | undefined;
            addresses?: string[] | undefined;
        } | undefined;
        instantiateDefaultPermission?: AccessType | undefined;
    } & {
        codeUploadAccess?: ({
            permission?: AccessType | undefined;
            address?: string | undefined;
            addresses?: string[] | undefined;
        } & {
            permission?: AccessType | undefined;
            address?: string | undefined;
            addresses?: (string[] & string[] & Record<Exclude<keyof I["codeUploadAccess"]["addresses"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["codeUploadAccess"], keyof AccessConfig>, never>) | undefined;
        instantiateDefaultPermission?: AccessType | undefined;
    } & Record<Exclude<keyof I, keyof Params>, never>>(object: I): Params;
};
export declare const CodeInfo: {
    typeUrl: string;
    encode(message: CodeInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): CodeInfo;
    fromJSON(object: any): CodeInfo;
    toJSON(message: CodeInfo): unknown;
    fromPartial<I extends {
        codeHash?: Uint8Array | undefined;
        creator?: string | undefined;
        instantiateConfig?: {
            permission?: AccessType | undefined;
            address?: string | undefined;
            addresses?: string[] | undefined;
        } | undefined;
    } & {
        codeHash?: Uint8Array | undefined;
        creator?: string | undefined;
        instantiateConfig?: ({
            permission?: AccessType | undefined;
            address?: string | undefined;
            addresses?: string[] | undefined;
        } & {
            permission?: AccessType | undefined;
            address?: string | undefined;
            addresses?: (string[] & string[] & Record<Exclude<keyof I["instantiateConfig"]["addresses"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["instantiateConfig"], keyof AccessConfig>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof CodeInfo>, never>>(object: I): CodeInfo;
};
export declare const ContractInfo: {
    typeUrl: string;
    encode(message: ContractInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ContractInfo;
    fromJSON(object: any): ContractInfo;
    toJSON(message: ContractInfo): unknown;
    fromPartial<I extends {
        codeId?: bigint | undefined;
        creator?: string | undefined;
        admin?: string | undefined;
        label?: string | undefined;
        created?: {
            blockHeight?: bigint | undefined;
            txIndex?: bigint | undefined;
        } | undefined;
        ibcPortId?: string | undefined;
        extension?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
    } & {
        codeId?: bigint | undefined;
        creator?: string | undefined;
        admin?: string | undefined;
        label?: string | undefined;
        created?: ({
            blockHeight?: bigint | undefined;
            txIndex?: bigint | undefined;
        } & {
            blockHeight?: bigint | undefined;
            txIndex?: bigint | undefined;
        } & Record<Exclude<keyof I["created"], keyof AbsoluteTxPosition>, never>) | undefined;
        ibcPortId?: string | undefined;
        extension?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["extension"], keyof Any>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ContractInfo>, never>>(object: I): ContractInfo;
};
export declare const ContractCodeHistoryEntry: {
    typeUrl: string;
    encode(message: ContractCodeHistoryEntry, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ContractCodeHistoryEntry;
    fromJSON(object: any): ContractCodeHistoryEntry;
    toJSON(message: ContractCodeHistoryEntry): unknown;
    fromPartial<I extends {
        operation?: ContractCodeHistoryOperationType | undefined;
        codeId?: bigint | undefined;
        updated?: {
            blockHeight?: bigint | undefined;
            txIndex?: bigint | undefined;
        } | undefined;
        msg?: Uint8Array | undefined;
    } & {
        operation?: ContractCodeHistoryOperationType | undefined;
        codeId?: bigint | undefined;
        updated?: ({
            blockHeight?: bigint | undefined;
            txIndex?: bigint | undefined;
        } & {
            blockHeight?: bigint | undefined;
            txIndex?: bigint | undefined;
        } & Record<Exclude<keyof I["updated"], keyof AbsoluteTxPosition>, never>) | undefined;
        msg?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof ContractCodeHistoryEntry>, never>>(object: I): ContractCodeHistoryEntry;
};
export declare const AbsoluteTxPosition: {
    typeUrl: string;
    encode(message: AbsoluteTxPosition, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): AbsoluteTxPosition;
    fromJSON(object: any): AbsoluteTxPosition;
    toJSON(message: AbsoluteTxPosition): unknown;
    fromPartial<I extends {
        blockHeight?: bigint | undefined;
        txIndex?: bigint | undefined;
    } & {
        blockHeight?: bigint | undefined;
        txIndex?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof AbsoluteTxPosition>, never>>(object: I): AbsoluteTxPosition;
};
export declare const Model: {
    typeUrl: string;
    encode(message: Model, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Model;
    fromJSON(object: any): Model;
    toJSON(message: Model): unknown;
    fromPartial<I extends {
        key?: Uint8Array | undefined;
        value?: Uint8Array | undefined;
    } & {
        key?: Uint8Array | undefined;
        value?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof Model>, never>>(object: I): Model;
};
