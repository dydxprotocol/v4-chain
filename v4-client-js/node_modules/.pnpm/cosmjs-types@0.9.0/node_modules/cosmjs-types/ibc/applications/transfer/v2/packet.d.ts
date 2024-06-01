import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "ibc.applications.transfer.v2";
/**
 * FungibleTokenPacketData defines a struct for the packet payload
 * See FungibleTokenPacketData spec:
 * https://github.com/cosmos/ibc/tree/master/spec/app/ics-020-fungible-token-transfer#data-structures
 */
export interface FungibleTokenPacketData {
    /** the token denomination to be transferred */
    denom: string;
    /** the token amount to be transferred */
    amount: string;
    /** the sender address */
    sender: string;
    /** the recipient address on the destination chain */
    receiver: string;
    /** optional memo */
    memo: string;
}
export declare const FungibleTokenPacketData: {
    typeUrl: string;
    encode(message: FungibleTokenPacketData, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): FungibleTokenPacketData;
    fromJSON(object: any): FungibleTokenPacketData;
    toJSON(message: FungibleTokenPacketData): unknown;
    fromPartial<I extends {
        denom?: string | undefined;
        amount?: string | undefined;
        sender?: string | undefined;
        receiver?: string | undefined;
        memo?: string | undefined;
    } & {
        denom?: string | undefined;
        amount?: string | undefined;
        sender?: string | undefined;
        receiver?: string | undefined;
        memo?: string | undefined;
    } & Record<Exclude<keyof I, keyof FungibleTokenPacketData>, never>>(object: I): FungibleTokenPacketData;
};
