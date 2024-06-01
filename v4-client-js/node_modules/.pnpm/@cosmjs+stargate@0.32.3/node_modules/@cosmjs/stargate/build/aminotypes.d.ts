import { AminoMsg } from "@cosmjs/amino";
import { EncodeObject } from "@cosmjs/proto-signing";
export interface AminoConverter {
    readonly aminoType: string;
    readonly toAmino: (value: any) => any;
    readonly fromAmino: (value: any) => any;
}
/** A map from protobuf type URL to the AminoConverter implementation if supported on chain */
export type AminoConverters = Record<string, AminoConverter>;
/**
 * A map from Stargate message types as used in the messages's `Any` type
 * to Amino types.
 */
export declare class AminoTypes {
    private readonly register;
    constructor(types: AminoConverters);
    toAmino({ typeUrl, value }: EncodeObject): AminoMsg;
    fromAmino({ type, value }: AminoMsg): EncodeObject;
}
