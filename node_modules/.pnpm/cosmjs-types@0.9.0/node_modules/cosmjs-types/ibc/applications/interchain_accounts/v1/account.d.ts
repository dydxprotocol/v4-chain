import { BaseAccount } from "../../../../cosmos/auth/v1beta1/auth";
import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "ibc.applications.interchain_accounts.v1";
/** An InterchainAccount is defined as a BaseAccount & the address of the account owner on the controller chain */
export interface InterchainAccount {
    baseAccount?: BaseAccount;
    accountOwner: string;
}
export declare const InterchainAccount: {
    typeUrl: string;
    encode(message: InterchainAccount, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): InterchainAccount;
    fromJSON(object: any): InterchainAccount;
    toJSON(message: InterchainAccount): unknown;
    fromPartial<I extends {
        baseAccount?: {
            address?: string | undefined;
            pubKey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            accountNumber?: bigint | undefined;
            sequence?: bigint | undefined;
        } | undefined;
        accountOwner?: string | undefined;
    } & {
        baseAccount?: ({
            address?: string | undefined;
            pubKey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            accountNumber?: bigint | undefined;
            sequence?: bigint | undefined;
        } & {
            address?: string | undefined;
            pubKey?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["baseAccount"]["pubKey"], keyof import("../../../../google/protobuf/any").Any>, never>) | undefined;
            accountNumber?: bigint | undefined;
            sequence?: bigint | undefined;
        } & Record<Exclude<keyof I["baseAccount"], keyof BaseAccount>, never>) | undefined;
        accountOwner?: string | undefined;
    } & Record<Exclude<keyof I, keyof InterchainAccount>, never>>(object: I): InterchainAccount;
};
