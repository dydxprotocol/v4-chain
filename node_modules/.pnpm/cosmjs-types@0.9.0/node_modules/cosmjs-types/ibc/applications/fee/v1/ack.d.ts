import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "ibc.applications.fee.v1";
/** IncentivizedAcknowledgement is the acknowledgement format to be used by applications wrapped in the fee middleware */
export interface IncentivizedAcknowledgement {
    /** the underlying app acknowledgement bytes */
    appAcknowledgement: Uint8Array;
    /** the relayer address which submits the recv packet message */
    forwardRelayerAddress: string;
    /** success flag of the base application callback */
    underlyingAppSuccess: boolean;
}
export declare const IncentivizedAcknowledgement: {
    typeUrl: string;
    encode(message: IncentivizedAcknowledgement, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): IncentivizedAcknowledgement;
    fromJSON(object: any): IncentivizedAcknowledgement;
    toJSON(message: IncentivizedAcknowledgement): unknown;
    fromPartial<I extends {
        appAcknowledgement?: Uint8Array | undefined;
        forwardRelayerAddress?: string | undefined;
        underlyingAppSuccess?: boolean | undefined;
    } & {
        appAcknowledgement?: Uint8Array | undefined;
        forwardRelayerAddress?: string | undefined;
        underlyingAppSuccess?: boolean | undefined;
    } & Record<Exclude<keyof I, keyof IncentivizedAcknowledgement>, never>>(object: I): IncentivizedAcknowledgement;
};
