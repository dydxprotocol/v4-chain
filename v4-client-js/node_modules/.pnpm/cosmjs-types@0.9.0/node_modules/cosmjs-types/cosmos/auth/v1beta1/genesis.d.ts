import { Params } from "./auth";
import { Any } from "../../../google/protobuf/any";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.auth.v1beta1";
/** GenesisState defines the auth module's genesis state. */
export interface GenesisState {
    /** params defines all the parameters of the module. */
    params: Params;
    /** accounts are the accounts present at genesis. */
    accounts: Any[];
}
export declare const GenesisState: {
    typeUrl: string;
    encode(message: GenesisState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial<I extends {
        params?: {
            maxMemoCharacters?: bigint | undefined;
            txSigLimit?: bigint | undefined;
            txSizeCostPerByte?: bigint | undefined;
            sigVerifyCostEd25519?: bigint | undefined;
            sigVerifyCostSecp256k1?: bigint | undefined;
        } | undefined;
        accounts?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[] | undefined;
    } & {
        params?: ({
            maxMemoCharacters?: bigint | undefined;
            txSigLimit?: bigint | undefined;
            txSizeCostPerByte?: bigint | undefined;
            sigVerifyCostEd25519?: bigint | undefined;
            sigVerifyCostSecp256k1?: bigint | undefined;
        } & {
            maxMemoCharacters?: bigint | undefined;
            txSigLimit?: bigint | undefined;
            txSizeCostPerByte?: bigint | undefined;
            sigVerifyCostEd25519?: bigint | undefined;
            sigVerifyCostSecp256k1?: bigint | undefined;
        } & Record<Exclude<keyof I["params"], keyof Params>, never>) | undefined;
        accounts?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[] & ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["accounts"][number], keyof Any>, never>)[] & Record<Exclude<keyof I["accounts"], keyof {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GenesisState>, never>>(object: I): GenesisState;
};
