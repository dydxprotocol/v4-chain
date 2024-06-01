import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "ibc.applications.interchain_accounts.v1";
/**
 * Metadata defines a set of protocol specific data encoded into the ICS27 channel version bytestring
 * See ICS004: https://github.com/cosmos/ibc/tree/master/spec/core/ics-004-channel-and-packet-semantics#Versioning
 */
export interface Metadata {
    /** version defines the ICS27 protocol version */
    version: string;
    /** controller_connection_id is the connection identifier associated with the controller chain */
    controllerConnectionId: string;
    /** host_connection_id is the connection identifier associated with the host chain */
    hostConnectionId: string;
    /**
     * address defines the interchain account address to be fulfilled upon the OnChanOpenTry handshake step
     * NOTE: the address field is empty on the OnChanOpenInit handshake step
     */
    address: string;
    /** encoding defines the supported codec format */
    encoding: string;
    /** tx_type defines the type of transactions the interchain account can execute */
    txType: string;
}
export declare const Metadata: {
    typeUrl: string;
    encode(message: Metadata, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Metadata;
    fromJSON(object: any): Metadata;
    toJSON(message: Metadata): unknown;
    fromPartial<I extends {
        version?: string | undefined;
        controllerConnectionId?: string | undefined;
        hostConnectionId?: string | undefined;
        address?: string | undefined;
        encoding?: string | undefined;
        txType?: string | undefined;
    } & {
        version?: string | undefined;
        controllerConnectionId?: string | undefined;
        hostConnectionId?: string | undefined;
        address?: string | undefined;
        encoding?: string | undefined;
        txType?: string | undefined;
    } & Record<Exclude<keyof I, keyof Metadata>, never>>(object: I): Metadata;
};
