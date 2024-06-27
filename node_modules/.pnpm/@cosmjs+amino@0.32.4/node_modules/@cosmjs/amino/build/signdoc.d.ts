import { Coin } from "./coins";
export interface AminoMsg {
    readonly type: string;
    readonly value: any;
}
export interface StdFee {
    readonly amount: readonly Coin[];
    readonly gas: string;
    /** The granter address that is used for paying with feegrants */
    readonly granter?: string;
    /** The fee payer address. The payer must have signed the transaction. */
    readonly payer?: string;
}
/**
 * The document to be signed
 *
 * @see https://docs.cosmos.network/master/modules/auth/03_types.html#stdsigndoc
 */
export interface StdSignDoc {
    readonly chain_id: string;
    readonly account_number: string;
    readonly sequence: string;
    readonly fee: StdFee;
    readonly msgs: readonly AminoMsg[];
    readonly memo: string;
    readonly timeout_height?: string;
}
/** Returns a JSON string with objects sorted by key */
export declare function sortedJsonStringify(obj: any): string;
export declare function makeSignDoc(msgs: readonly AminoMsg[], fee: StdFee, chainId: string, memo: string | undefined, accountNumber: number | string, sequence: number | string, timeout_height?: bigint): StdSignDoc;
/**
 * Takes a valid JSON document and performs the following escapings in string values:
 *
 * `&` -> `\u0026`
 * `<` -> `\u003c`
 * `>` -> `\u003e`
 *
 * Since those characters do not occur in other places of the JSON document, only
 * string values are affected.
 *
 * If the input is invalid JSON, the behaviour is undefined.
 */
export declare function escapeCharacters(input: string): string;
export declare function serializeSignDoc(signDoc: StdSignDoc): Uint8Array;
