import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "ibc.applications.fee.v1";
/**
 * Metadata defines the ICS29 channel specific metadata encoded into the channel version bytestring
 * See ICS004: https://github.com/cosmos/ibc/tree/master/spec/core/ics-004-channel-and-packet-semantics#Versioning
 */
export interface Metadata {
    /** fee_version defines the ICS29 fee version */
    feeVersion: string;
    /** app_version defines the underlying application version, which may or may not be a JSON encoded bytestring */
    appVersion: string;
}
export declare const Metadata: {
    typeUrl: string;
    encode(message: Metadata, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Metadata;
    fromJSON(object: any): Metadata;
    toJSON(message: Metadata): unknown;
    fromPartial<I extends {
        feeVersion?: string | undefined;
        appVersion?: string | undefined;
    } & {
        feeVersion?: string | undefined;
        appVersion?: string | undefined;
    } & Record<Exclude<keyof I, keyof Metadata>, never>>(object: I): Metadata;
};
