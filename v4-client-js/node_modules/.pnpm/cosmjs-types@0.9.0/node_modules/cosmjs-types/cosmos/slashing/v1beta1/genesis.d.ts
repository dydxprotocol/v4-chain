import { Params, ValidatorSigningInfo } from "./slashing";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.slashing.v1beta1";
/** GenesisState defines the slashing module's genesis state. */
export interface GenesisState {
    /** params defines all the parameters of the module. */
    params: Params;
    /**
     * signing_infos represents a map between validator addresses and their
     * signing infos.
     */
    signingInfos: SigningInfo[];
    /**
     * missed_blocks represents a map between validator addresses and their
     * missed blocks.
     */
    missedBlocks: ValidatorMissedBlocks[];
}
/** SigningInfo stores validator signing info of corresponding address. */
export interface SigningInfo {
    /** address is the validator address. */
    address: string;
    /** validator_signing_info represents the signing info of this validator. */
    validatorSigningInfo: ValidatorSigningInfo;
}
/**
 * ValidatorMissedBlocks contains array of missed blocks of corresponding
 * address.
 */
export interface ValidatorMissedBlocks {
    /** address is the validator address. */
    address: string;
    /** missed_blocks is an array of missed blocks by the validator. */
    missedBlocks: MissedBlock[];
}
/** MissedBlock contains height and missed status as boolean. */
export interface MissedBlock {
    /** index is the height at which the block was missed. */
    index: bigint;
    /** missed is the missed status. */
    missed: boolean;
}
export declare const GenesisState: {
    typeUrl: string;
    encode(message: GenesisState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial<I extends {
        params?: {
            signedBlocksWindow?: bigint | undefined;
            minSignedPerWindow?: Uint8Array | undefined;
            downtimeJailDuration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            slashFractionDoubleSign?: Uint8Array | undefined;
            slashFractionDowntime?: Uint8Array | undefined;
        } | undefined;
        signingInfos?: {
            address?: string | undefined;
            validatorSigningInfo?: {
                address?: string | undefined;
                startHeight?: bigint | undefined;
                indexOffset?: bigint | undefined;
                jailedUntil?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                tombstoned?: boolean | undefined;
                missedBlocksCounter?: bigint | undefined;
            } | undefined;
        }[] | undefined;
        missedBlocks?: {
            address?: string | undefined;
            missedBlocks?: {
                index?: bigint | undefined;
                missed?: boolean | undefined;
            }[] | undefined;
        }[] | undefined;
    } & {
        params?: ({
            signedBlocksWindow?: bigint | undefined;
            minSignedPerWindow?: Uint8Array | undefined;
            downtimeJailDuration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            slashFractionDoubleSign?: Uint8Array | undefined;
            slashFractionDowntime?: Uint8Array | undefined;
        } & {
            signedBlocksWindow?: bigint | undefined;
            minSignedPerWindow?: Uint8Array | undefined;
            downtimeJailDuration?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["params"]["downtimeJailDuration"], keyof import("../../../google/protobuf/duration").Duration>, never>) | undefined;
            slashFractionDoubleSign?: Uint8Array | undefined;
            slashFractionDowntime?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["params"], keyof Params>, never>) | undefined;
        signingInfos?: ({
            address?: string | undefined;
            validatorSigningInfo?: {
                address?: string | undefined;
                startHeight?: bigint | undefined;
                indexOffset?: bigint | undefined;
                jailedUntil?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                tombstoned?: boolean | undefined;
                missedBlocksCounter?: bigint | undefined;
            } | undefined;
        }[] & ({
            address?: string | undefined;
            validatorSigningInfo?: {
                address?: string | undefined;
                startHeight?: bigint | undefined;
                indexOffset?: bigint | undefined;
                jailedUntil?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                tombstoned?: boolean | undefined;
                missedBlocksCounter?: bigint | undefined;
            } | undefined;
        } & {
            address?: string | undefined;
            validatorSigningInfo?: ({
                address?: string | undefined;
                startHeight?: bigint | undefined;
                indexOffset?: bigint | undefined;
                jailedUntil?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                tombstoned?: boolean | undefined;
                missedBlocksCounter?: bigint | undefined;
            } & {
                address?: string | undefined;
                startHeight?: bigint | undefined;
                indexOffset?: bigint | undefined;
                jailedUntil?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["signingInfos"][number]["validatorSigningInfo"]["jailedUntil"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
                tombstoned?: boolean | undefined;
                missedBlocksCounter?: bigint | undefined;
            } & Record<Exclude<keyof I["signingInfos"][number]["validatorSigningInfo"], keyof ValidatorSigningInfo>, never>) | undefined;
        } & Record<Exclude<keyof I["signingInfos"][number], keyof SigningInfo>, never>)[] & Record<Exclude<keyof I["signingInfos"], keyof {
            address?: string | undefined;
            validatorSigningInfo?: {
                address?: string | undefined;
                startHeight?: bigint | undefined;
                indexOffset?: bigint | undefined;
                jailedUntil?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                tombstoned?: boolean | undefined;
                missedBlocksCounter?: bigint | undefined;
            } | undefined;
        }[]>, never>) | undefined;
        missedBlocks?: ({
            address?: string | undefined;
            missedBlocks?: {
                index?: bigint | undefined;
                missed?: boolean | undefined;
            }[] | undefined;
        }[] & ({
            address?: string | undefined;
            missedBlocks?: {
                index?: bigint | undefined;
                missed?: boolean | undefined;
            }[] | undefined;
        } & {
            address?: string | undefined;
            missedBlocks?: ({
                index?: bigint | undefined;
                missed?: boolean | undefined;
            }[] & ({
                index?: bigint | undefined;
                missed?: boolean | undefined;
            } & {
                index?: bigint | undefined;
                missed?: boolean | undefined;
            } & Record<Exclude<keyof I["missedBlocks"][number]["missedBlocks"][number], keyof MissedBlock>, never>)[] & Record<Exclude<keyof I["missedBlocks"][number]["missedBlocks"], keyof {
                index?: bigint | undefined;
                missed?: boolean | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["missedBlocks"][number], keyof ValidatorMissedBlocks>, never>)[] & Record<Exclude<keyof I["missedBlocks"], keyof {
            address?: string | undefined;
            missedBlocks?: {
                index?: bigint | undefined;
                missed?: boolean | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GenesisState>, never>>(object: I): GenesisState;
};
export declare const SigningInfo: {
    typeUrl: string;
    encode(message: SigningInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): SigningInfo;
    fromJSON(object: any): SigningInfo;
    toJSON(message: SigningInfo): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        validatorSigningInfo?: {
            address?: string | undefined;
            startHeight?: bigint | undefined;
            indexOffset?: bigint | undefined;
            jailedUntil?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            tombstoned?: boolean | undefined;
            missedBlocksCounter?: bigint | undefined;
        } | undefined;
    } & {
        address?: string | undefined;
        validatorSigningInfo?: ({
            address?: string | undefined;
            startHeight?: bigint | undefined;
            indexOffset?: bigint | undefined;
            jailedUntil?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            tombstoned?: boolean | undefined;
            missedBlocksCounter?: bigint | undefined;
        } & {
            address?: string | undefined;
            startHeight?: bigint | undefined;
            indexOffset?: bigint | undefined;
            jailedUntil?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["validatorSigningInfo"]["jailedUntil"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            tombstoned?: boolean | undefined;
            missedBlocksCounter?: bigint | undefined;
        } & Record<Exclude<keyof I["validatorSigningInfo"], keyof ValidatorSigningInfo>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof SigningInfo>, never>>(object: I): SigningInfo;
};
export declare const ValidatorMissedBlocks: {
    typeUrl: string;
    encode(message: ValidatorMissedBlocks, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ValidatorMissedBlocks;
    fromJSON(object: any): ValidatorMissedBlocks;
    toJSON(message: ValidatorMissedBlocks): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        missedBlocks?: {
            index?: bigint | undefined;
            missed?: boolean | undefined;
        }[] | undefined;
    } & {
        address?: string | undefined;
        missedBlocks?: ({
            index?: bigint | undefined;
            missed?: boolean | undefined;
        }[] & ({
            index?: bigint | undefined;
            missed?: boolean | undefined;
        } & {
            index?: bigint | undefined;
            missed?: boolean | undefined;
        } & Record<Exclude<keyof I["missedBlocks"][number], keyof MissedBlock>, never>)[] & Record<Exclude<keyof I["missedBlocks"], keyof {
            index?: bigint | undefined;
            missed?: boolean | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ValidatorMissedBlocks>, never>>(object: I): ValidatorMissedBlocks;
};
export declare const MissedBlock: {
    typeUrl: string;
    encode(message: MissedBlock, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MissedBlock;
    fromJSON(object: any): MissedBlock;
    toJSON(message: MissedBlock): unknown;
    fromPartial<I extends {
        index?: bigint | undefined;
        missed?: boolean | undefined;
    } & {
        index?: bigint | undefined;
        missed?: boolean | undefined;
    } & Record<Exclude<keyof I, keyof MissedBlock>, never>>(object: I): MissedBlock;
};
