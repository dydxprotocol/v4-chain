import { BinaryReader, BinaryWriter } from "../../binary";
export declare const protobufPackage = "tendermint.version";
/**
 * App includes the protocol and software version for the application.
 * This information is included in ResponseInfo. The App.Protocol can be
 * updated in ResponseEndBlock.
 */
export interface App {
    protocol: bigint;
    software: string;
}
/**
 * Consensus captures the consensus rules for processing a block in the blockchain,
 * including all blockchain data structures and the rules of the application's
 * state transition machine.
 */
export interface Consensus {
    block: bigint;
    app: bigint;
}
export declare const App: {
    typeUrl: string;
    encode(message: App, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): App;
    fromJSON(object: any): App;
    toJSON(message: App): unknown;
    fromPartial<I extends {
        protocol?: bigint | undefined;
        software?: string | undefined;
    } & {
        protocol?: bigint | undefined;
        software?: string | undefined;
    } & Record<Exclude<keyof I, keyof App>, never>>(object: I): App;
};
export declare const Consensus: {
    typeUrl: string;
    encode(message: Consensus, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Consensus;
    fromJSON(object: any): Consensus;
    toJSON(message: Consensus): unknown;
    fromPartial<I extends {
        block?: bigint | undefined;
        app?: bigint | undefined;
    } & {
        block?: bigint | undefined;
        app?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof Consensus>, never>>(object: I): Consensus;
};
