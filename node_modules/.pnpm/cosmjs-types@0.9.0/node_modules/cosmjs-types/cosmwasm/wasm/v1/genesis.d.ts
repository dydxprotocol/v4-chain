import { Params, CodeInfo, ContractInfo, Model, ContractCodeHistoryEntry } from "./types";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmwasm.wasm.v1";
/** GenesisState - genesis state of x/wasm */
export interface GenesisState {
    params: Params;
    codes: Code[];
    contracts: Contract[];
    sequences: Sequence[];
}
/** Code struct encompasses CodeInfo and CodeBytes */
export interface Code {
    codeId: bigint;
    codeInfo: CodeInfo;
    codeBytes: Uint8Array;
    /** Pinned to wasmvm cache */
    pinned: boolean;
}
/** Contract struct encompasses ContractAddress, ContractInfo, and ContractState */
export interface Contract {
    contractAddress: string;
    contractInfo: ContractInfo;
    contractState: Model[];
    contractCodeHistory: ContractCodeHistoryEntry[];
}
/** Sequence key and value of an id generation counter */
export interface Sequence {
    idKey: Uint8Array;
    value: bigint;
}
export declare const GenesisState: {
    typeUrl: string;
    encode(message: GenesisState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial<I extends {
        params?: {
            codeUploadAccess?: {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } | undefined;
            instantiateDefaultPermission?: import("./types").AccessType | undefined;
        } | undefined;
        codes?: {
            codeId?: bigint | undefined;
            codeInfo?: {
                codeHash?: Uint8Array | undefined;
                creator?: string | undefined;
                instantiateConfig?: {
                    permission?: import("./types").AccessType | undefined;
                    address?: string | undefined;
                    addresses?: string[] | undefined;
                } | undefined;
            } | undefined;
            codeBytes?: Uint8Array | undefined;
            pinned?: boolean | undefined;
        }[] | undefined;
        contracts?: {
            contractAddress?: string | undefined;
            contractInfo?: {
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
            } | undefined;
            contractState?: {
                key?: Uint8Array | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
            contractCodeHistory?: {
                operation?: import("./types").ContractCodeHistoryOperationType | undefined;
                codeId?: bigint | undefined;
                updated?: {
                    blockHeight?: bigint | undefined;
                    txIndex?: bigint | undefined;
                } | undefined;
                msg?: Uint8Array | undefined;
            }[] | undefined;
        }[] | undefined;
        sequences?: {
            idKey?: Uint8Array | undefined;
            value?: bigint | undefined;
        }[] | undefined;
    } & {
        params?: ({
            codeUploadAccess?: {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } | undefined;
            instantiateDefaultPermission?: import("./types").AccessType | undefined;
        } & {
            codeUploadAccess?: ({
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } & {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: (string[] & string[] & Record<Exclude<keyof I["params"]["codeUploadAccess"]["addresses"], keyof string[]>, never>) | undefined;
            } & Record<Exclude<keyof I["params"]["codeUploadAccess"], keyof import("./types").AccessConfig>, never>) | undefined;
            instantiateDefaultPermission?: import("./types").AccessType | undefined;
        } & Record<Exclude<keyof I["params"], keyof Params>, never>) | undefined;
        codes?: ({
            codeId?: bigint | undefined;
            codeInfo?: {
                codeHash?: Uint8Array | undefined;
                creator?: string | undefined;
                instantiateConfig?: {
                    permission?: import("./types").AccessType | undefined;
                    address?: string | undefined;
                    addresses?: string[] | undefined;
                } | undefined;
            } | undefined;
            codeBytes?: Uint8Array | undefined;
            pinned?: boolean | undefined;
        }[] & ({
            codeId?: bigint | undefined;
            codeInfo?: {
                codeHash?: Uint8Array | undefined;
                creator?: string | undefined;
                instantiateConfig?: {
                    permission?: import("./types").AccessType | undefined;
                    address?: string | undefined;
                    addresses?: string[] | undefined;
                } | undefined;
            } | undefined;
            codeBytes?: Uint8Array | undefined;
            pinned?: boolean | undefined;
        } & {
            codeId?: bigint | undefined;
            codeInfo?: ({
                codeHash?: Uint8Array | undefined;
                creator?: string | undefined;
                instantiateConfig?: {
                    permission?: import("./types").AccessType | undefined;
                    address?: string | undefined;
                    addresses?: string[] | undefined;
                } | undefined;
            } & {
                codeHash?: Uint8Array | undefined;
                creator?: string | undefined;
                instantiateConfig?: ({
                    permission?: import("./types").AccessType | undefined;
                    address?: string | undefined;
                    addresses?: string[] | undefined;
                } & {
                    permission?: import("./types").AccessType | undefined;
                    address?: string | undefined;
                    addresses?: (string[] & string[] & Record<Exclude<keyof I["codes"][number]["codeInfo"]["instantiateConfig"]["addresses"], keyof string[]>, never>) | undefined;
                } & Record<Exclude<keyof I["codes"][number]["codeInfo"]["instantiateConfig"], keyof import("./types").AccessConfig>, never>) | undefined;
            } & Record<Exclude<keyof I["codes"][number]["codeInfo"], keyof CodeInfo>, never>) | undefined;
            codeBytes?: Uint8Array | undefined;
            pinned?: boolean | undefined;
        } & Record<Exclude<keyof I["codes"][number], keyof Code>, never>)[] & Record<Exclude<keyof I["codes"], keyof {
            codeId?: bigint | undefined;
            codeInfo?: {
                codeHash?: Uint8Array | undefined;
                creator?: string | undefined;
                instantiateConfig?: {
                    permission?: import("./types").AccessType | undefined;
                    address?: string | undefined;
                    addresses?: string[] | undefined;
                } | undefined;
            } | undefined;
            codeBytes?: Uint8Array | undefined;
            pinned?: boolean | undefined;
        }[]>, never>) | undefined;
        contracts?: ({
            contractAddress?: string | undefined;
            contractInfo?: {
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
            } | undefined;
            contractState?: {
                key?: Uint8Array | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
            contractCodeHistory?: {
                operation?: import("./types").ContractCodeHistoryOperationType | undefined;
                codeId?: bigint | undefined;
                updated?: {
                    blockHeight?: bigint | undefined;
                    txIndex?: bigint | undefined;
                } | undefined;
                msg?: Uint8Array | undefined;
            }[] | undefined;
        }[] & ({
            contractAddress?: string | undefined;
            contractInfo?: {
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
            } | undefined;
            contractState?: {
                key?: Uint8Array | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
            contractCodeHistory?: {
                operation?: import("./types").ContractCodeHistoryOperationType | undefined;
                codeId?: bigint | undefined;
                updated?: {
                    blockHeight?: bigint | undefined;
                    txIndex?: bigint | undefined;
                } | undefined;
                msg?: Uint8Array | undefined;
            }[] | undefined;
        } & {
            contractAddress?: string | undefined;
            contractInfo?: ({
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
                } & Record<Exclude<keyof I["contracts"][number]["contractInfo"]["created"], keyof import("./types").AbsoluteTxPosition>, never>) | undefined;
                ibcPortId?: string | undefined;
                extension?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["contracts"][number]["contractInfo"]["extension"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
            } & Record<Exclude<keyof I["contracts"][number]["contractInfo"], keyof ContractInfo>, never>) | undefined;
            contractState?: ({
                key?: Uint8Array | undefined;
                value?: Uint8Array | undefined;
            }[] & ({
                key?: Uint8Array | undefined;
                value?: Uint8Array | undefined;
            } & {
                key?: Uint8Array | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["contracts"][number]["contractState"][number], keyof Model>, never>)[] & Record<Exclude<keyof I["contracts"][number]["contractState"], keyof {
                key?: Uint8Array | undefined;
                value?: Uint8Array | undefined;
            }[]>, never>) | undefined;
            contractCodeHistory?: ({
                operation?: import("./types").ContractCodeHistoryOperationType | undefined;
                codeId?: bigint | undefined;
                updated?: {
                    blockHeight?: bigint | undefined;
                    txIndex?: bigint | undefined;
                } | undefined;
                msg?: Uint8Array | undefined;
            }[] & ({
                operation?: import("./types").ContractCodeHistoryOperationType | undefined;
                codeId?: bigint | undefined;
                updated?: {
                    blockHeight?: bigint | undefined;
                    txIndex?: bigint | undefined;
                } | undefined;
                msg?: Uint8Array | undefined;
            } & {
                operation?: import("./types").ContractCodeHistoryOperationType | undefined;
                codeId?: bigint | undefined;
                updated?: ({
                    blockHeight?: bigint | undefined;
                    txIndex?: bigint | undefined;
                } & {
                    blockHeight?: bigint | undefined;
                    txIndex?: bigint | undefined;
                } & Record<Exclude<keyof I["contracts"][number]["contractCodeHistory"][number]["updated"], keyof import("./types").AbsoluteTxPosition>, never>) | undefined;
                msg?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["contracts"][number]["contractCodeHistory"][number], keyof ContractCodeHistoryEntry>, never>)[] & Record<Exclude<keyof I["contracts"][number]["contractCodeHistory"], keyof {
                operation?: import("./types").ContractCodeHistoryOperationType | undefined;
                codeId?: bigint | undefined;
                updated?: {
                    blockHeight?: bigint | undefined;
                    txIndex?: bigint | undefined;
                } | undefined;
                msg?: Uint8Array | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["contracts"][number], keyof Contract>, never>)[] & Record<Exclude<keyof I["contracts"], keyof {
            contractAddress?: string | undefined;
            contractInfo?: {
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
            } | undefined;
            contractState?: {
                key?: Uint8Array | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
            contractCodeHistory?: {
                operation?: import("./types").ContractCodeHistoryOperationType | undefined;
                codeId?: bigint | undefined;
                updated?: {
                    blockHeight?: bigint | undefined;
                    txIndex?: bigint | undefined;
                } | undefined;
                msg?: Uint8Array | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
        sequences?: ({
            idKey?: Uint8Array | undefined;
            value?: bigint | undefined;
        }[] & ({
            idKey?: Uint8Array | undefined;
            value?: bigint | undefined;
        } & {
            idKey?: Uint8Array | undefined;
            value?: bigint | undefined;
        } & Record<Exclude<keyof I["sequences"][number], keyof Sequence>, never>)[] & Record<Exclude<keyof I["sequences"], keyof {
            idKey?: Uint8Array | undefined;
            value?: bigint | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GenesisState>, never>>(object: I): GenesisState;
};
export declare const Code: {
    typeUrl: string;
    encode(message: Code, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Code;
    fromJSON(object: any): Code;
    toJSON(message: Code): unknown;
    fromPartial<I extends {
        codeId?: bigint | undefined;
        codeInfo?: {
            codeHash?: Uint8Array | undefined;
            creator?: string | undefined;
            instantiateConfig?: {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } | undefined;
        } | undefined;
        codeBytes?: Uint8Array | undefined;
        pinned?: boolean | undefined;
    } & {
        codeId?: bigint | undefined;
        codeInfo?: ({
            codeHash?: Uint8Array | undefined;
            creator?: string | undefined;
            instantiateConfig?: {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } | undefined;
        } & {
            codeHash?: Uint8Array | undefined;
            creator?: string | undefined;
            instantiateConfig?: ({
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } & {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: (string[] & string[] & Record<Exclude<keyof I["codeInfo"]["instantiateConfig"]["addresses"], keyof string[]>, never>) | undefined;
            } & Record<Exclude<keyof I["codeInfo"]["instantiateConfig"], keyof import("./types").AccessConfig>, never>) | undefined;
        } & Record<Exclude<keyof I["codeInfo"], keyof CodeInfo>, never>) | undefined;
        codeBytes?: Uint8Array | undefined;
        pinned?: boolean | undefined;
    } & Record<Exclude<keyof I, keyof Code>, never>>(object: I): Code;
};
export declare const Contract: {
    typeUrl: string;
    encode(message: Contract, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Contract;
    fromJSON(object: any): Contract;
    toJSON(message: Contract): unknown;
    fromPartial<I extends {
        contractAddress?: string | undefined;
        contractInfo?: {
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
        } | undefined;
        contractState?: {
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
        }[] | undefined;
        contractCodeHistory?: {
            operation?: import("./types").ContractCodeHistoryOperationType | undefined;
            codeId?: bigint | undefined;
            updated?: {
                blockHeight?: bigint | undefined;
                txIndex?: bigint | undefined;
            } | undefined;
            msg?: Uint8Array | undefined;
        }[] | undefined;
    } & {
        contractAddress?: string | undefined;
        contractInfo?: ({
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
            } & Record<Exclude<keyof I["contractInfo"]["created"], keyof import("./types").AbsoluteTxPosition>, never>) | undefined;
            ibcPortId?: string | undefined;
            extension?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["contractInfo"]["extension"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
        } & Record<Exclude<keyof I["contractInfo"], keyof ContractInfo>, never>) | undefined;
        contractState?: ({
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
        }[] & ({
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
        } & {
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["contractState"][number], keyof Model>, never>)[] & Record<Exclude<keyof I["contractState"], keyof {
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
        }[]>, never>) | undefined;
        contractCodeHistory?: ({
            operation?: import("./types").ContractCodeHistoryOperationType | undefined;
            codeId?: bigint | undefined;
            updated?: {
                blockHeight?: bigint | undefined;
                txIndex?: bigint | undefined;
            } | undefined;
            msg?: Uint8Array | undefined;
        }[] & ({
            operation?: import("./types").ContractCodeHistoryOperationType | undefined;
            codeId?: bigint | undefined;
            updated?: {
                blockHeight?: bigint | undefined;
                txIndex?: bigint | undefined;
            } | undefined;
            msg?: Uint8Array | undefined;
        } & {
            operation?: import("./types").ContractCodeHistoryOperationType | undefined;
            codeId?: bigint | undefined;
            updated?: ({
                blockHeight?: bigint | undefined;
                txIndex?: bigint | undefined;
            } & {
                blockHeight?: bigint | undefined;
                txIndex?: bigint | undefined;
            } & Record<Exclude<keyof I["contractCodeHistory"][number]["updated"], keyof import("./types").AbsoluteTxPosition>, never>) | undefined;
            msg?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["contractCodeHistory"][number], keyof ContractCodeHistoryEntry>, never>)[] & Record<Exclude<keyof I["contractCodeHistory"], keyof {
            operation?: import("./types").ContractCodeHistoryOperationType | undefined;
            codeId?: bigint | undefined;
            updated?: {
                blockHeight?: bigint | undefined;
                txIndex?: bigint | undefined;
            } | undefined;
            msg?: Uint8Array | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Contract>, never>>(object: I): Contract;
};
export declare const Sequence: {
    typeUrl: string;
    encode(message: Sequence, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Sequence;
    fromJSON(object: any): Sequence;
    toJSON(message: Sequence): unknown;
    fromPartial<I extends {
        idKey?: Uint8Array | undefined;
        value?: bigint | undefined;
    } & {
        idKey?: Uint8Array | undefined;
        value?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof Sequence>, never>>(object: I): Sequence;
};
