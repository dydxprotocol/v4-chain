import { BinaryReader, BinaryWriter } from "../../binary";
export declare const protobufPackage = "tendermint.p2p";
export interface NetAddress {
    id: string;
    ip: string;
    port: number;
}
export interface ProtocolVersion {
    p2p: bigint;
    block: bigint;
    app: bigint;
}
export interface DefaultNodeInfo {
    protocolVersion: ProtocolVersion;
    defaultNodeId: string;
    listenAddr: string;
    network: string;
    version: string;
    channels: Uint8Array;
    moniker: string;
    other: DefaultNodeInfoOther;
}
export interface DefaultNodeInfoOther {
    txIndex: string;
    rpcAddress: string;
}
export declare const NetAddress: {
    typeUrl: string;
    encode(message: NetAddress, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): NetAddress;
    fromJSON(object: any): NetAddress;
    toJSON(message: NetAddress): unknown;
    fromPartial<I extends {
        id?: string | undefined;
        ip?: string | undefined;
        port?: number | undefined;
    } & {
        id?: string | undefined;
        ip?: string | undefined;
        port?: number | undefined;
    } & Record<Exclude<keyof I, keyof NetAddress>, never>>(object: I): NetAddress;
};
export declare const ProtocolVersion: {
    typeUrl: string;
    encode(message: ProtocolVersion, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ProtocolVersion;
    fromJSON(object: any): ProtocolVersion;
    toJSON(message: ProtocolVersion): unknown;
    fromPartial<I extends {
        p2p?: bigint | undefined;
        block?: bigint | undefined;
        app?: bigint | undefined;
    } & {
        p2p?: bigint | undefined;
        block?: bigint | undefined;
        app?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof ProtocolVersion>, never>>(object: I): ProtocolVersion;
};
export declare const DefaultNodeInfo: {
    typeUrl: string;
    encode(message: DefaultNodeInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DefaultNodeInfo;
    fromJSON(object: any): DefaultNodeInfo;
    toJSON(message: DefaultNodeInfo): unknown;
    fromPartial<I extends {
        protocolVersion?: {
            p2p?: bigint | undefined;
            block?: bigint | undefined;
            app?: bigint | undefined;
        } | undefined;
        defaultNodeId?: string | undefined;
        listenAddr?: string | undefined;
        network?: string | undefined;
        version?: string | undefined;
        channels?: Uint8Array | undefined;
        moniker?: string | undefined;
        other?: {
            txIndex?: string | undefined;
            rpcAddress?: string | undefined;
        } | undefined;
    } & {
        protocolVersion?: ({
            p2p?: bigint | undefined;
            block?: bigint | undefined;
            app?: bigint | undefined;
        } & {
            p2p?: bigint | undefined;
            block?: bigint | undefined;
            app?: bigint | undefined;
        } & Record<Exclude<keyof I["protocolVersion"], keyof ProtocolVersion>, never>) | undefined;
        defaultNodeId?: string | undefined;
        listenAddr?: string | undefined;
        network?: string | undefined;
        version?: string | undefined;
        channels?: Uint8Array | undefined;
        moniker?: string | undefined;
        other?: ({
            txIndex?: string | undefined;
            rpcAddress?: string | undefined;
        } & {
            txIndex?: string | undefined;
            rpcAddress?: string | undefined;
        } & Record<Exclude<keyof I["other"], keyof DefaultNodeInfoOther>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof DefaultNodeInfo>, never>>(object: I): DefaultNodeInfo;
};
export declare const DefaultNodeInfoOther: {
    typeUrl: string;
    encode(message: DefaultNodeInfoOther, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DefaultNodeInfoOther;
    fromJSON(object: any): DefaultNodeInfoOther;
    toJSON(message: DefaultNodeInfoOther): unknown;
    fromPartial<I extends {
        txIndex?: string | undefined;
        rpcAddress?: string | undefined;
    } & {
        txIndex?: string | undefined;
        rpcAddress?: string | undefined;
    } & Record<Exclude<keyof I, keyof DefaultNodeInfoOther>, never>>(object: I): DefaultNodeInfoOther;
};
