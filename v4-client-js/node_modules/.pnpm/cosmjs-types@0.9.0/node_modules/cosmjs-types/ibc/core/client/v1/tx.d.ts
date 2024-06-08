import { Any } from "../../../../google/protobuf/any";
import { BinaryReader, BinaryWriter } from "../../../../binary";
import { Rpc } from "../../../../helpers";
export declare const protobufPackage = "ibc.core.client.v1";
/** MsgCreateClient defines a message to create an IBC client */
export interface MsgCreateClient {
    /** light client state */
    clientState?: Any;
    /**
     * consensus state associated with the client that corresponds to a given
     * height.
     */
    consensusState?: Any;
    /** signer address */
    signer: string;
}
/** MsgCreateClientResponse defines the Msg/CreateClient response type. */
export interface MsgCreateClientResponse {
}
/**
 * MsgUpdateClient defines an sdk.Msg to update a IBC client state using
 * the given client message.
 */
export interface MsgUpdateClient {
    /** client unique identifier */
    clientId: string;
    /** client message to update the light client */
    clientMessage?: Any;
    /** signer address */
    signer: string;
}
/** MsgUpdateClientResponse defines the Msg/UpdateClient response type. */
export interface MsgUpdateClientResponse {
}
/**
 * MsgUpgradeClient defines an sdk.Msg to upgrade an IBC client to a new client
 * state
 */
export interface MsgUpgradeClient {
    /** client unique identifier */
    clientId: string;
    /** upgraded client state */
    clientState?: Any;
    /**
     * upgraded consensus state, only contains enough information to serve as a
     * basis of trust in update logic
     */
    consensusState?: Any;
    /** proof that old chain committed to new client */
    proofUpgradeClient: Uint8Array;
    /** proof that old chain committed to new consensus state */
    proofUpgradeConsensusState: Uint8Array;
    /** signer address */
    signer: string;
}
/** MsgUpgradeClientResponse defines the Msg/UpgradeClient response type. */
export interface MsgUpgradeClientResponse {
}
/**
 * MsgSubmitMisbehaviour defines an sdk.Msg type that submits Evidence for
 * light client misbehaviour.
 * Warning: DEPRECATED
 */
export interface MsgSubmitMisbehaviour {
    /** client unique identifier */
    /** @deprecated */
    clientId: string;
    /** misbehaviour used for freezing the light client */
    /** @deprecated */
    misbehaviour?: Any;
    /** signer address */
    /** @deprecated */
    signer: string;
}
/**
 * MsgSubmitMisbehaviourResponse defines the Msg/SubmitMisbehaviour response
 * type.
 */
export interface MsgSubmitMisbehaviourResponse {
}
export declare const MsgCreateClient: {
    typeUrl: string;
    encode(message: MsgCreateClient, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgCreateClient;
    fromJSON(object: any): MsgCreateClient;
    toJSON(message: MsgCreateClient): unknown;
    fromPartial<I extends {
        clientState?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        consensusState?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        signer?: string | undefined;
    } & {
        clientState?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["clientState"], keyof Any>, never>) | undefined;
        consensusState?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["consensusState"], keyof Any>, never>) | undefined;
        signer?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgCreateClient>, never>>(object: I): MsgCreateClient;
};
export declare const MsgCreateClientResponse: {
    typeUrl: string;
    encode(_: MsgCreateClientResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgCreateClientResponse;
    fromJSON(_: any): MsgCreateClientResponse;
    toJSON(_: MsgCreateClientResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgCreateClientResponse;
};
export declare const MsgUpdateClient: {
    typeUrl: string;
    encode(message: MsgUpdateClient, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateClient;
    fromJSON(object: any): MsgUpdateClient;
    toJSON(message: MsgUpdateClient): unknown;
    fromPartial<I extends {
        clientId?: string | undefined;
        clientMessage?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        signer?: string | undefined;
    } & {
        clientId?: string | undefined;
        clientMessage?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["clientMessage"], keyof Any>, never>) | undefined;
        signer?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgUpdateClient>, never>>(object: I): MsgUpdateClient;
};
export declare const MsgUpdateClientResponse: {
    typeUrl: string;
    encode(_: MsgUpdateClientResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateClientResponse;
    fromJSON(_: any): MsgUpdateClientResponse;
    toJSON(_: MsgUpdateClientResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgUpdateClientResponse;
};
export declare const MsgUpgradeClient: {
    typeUrl: string;
    encode(message: MsgUpgradeClient, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpgradeClient;
    fromJSON(object: any): MsgUpgradeClient;
    toJSON(message: MsgUpgradeClient): unknown;
    fromPartial<I extends {
        clientId?: string | undefined;
        clientState?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        consensusState?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        proofUpgradeClient?: Uint8Array | undefined;
        proofUpgradeConsensusState?: Uint8Array | undefined;
        signer?: string | undefined;
    } & {
        clientId?: string | undefined;
        clientState?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["clientState"], keyof Any>, never>) | undefined;
        consensusState?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["consensusState"], keyof Any>, never>) | undefined;
        proofUpgradeClient?: Uint8Array | undefined;
        proofUpgradeConsensusState?: Uint8Array | undefined;
        signer?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgUpgradeClient>, never>>(object: I): MsgUpgradeClient;
};
export declare const MsgUpgradeClientResponse: {
    typeUrl: string;
    encode(_: MsgUpgradeClientResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpgradeClientResponse;
    fromJSON(_: any): MsgUpgradeClientResponse;
    toJSON(_: MsgUpgradeClientResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgUpgradeClientResponse;
};
export declare const MsgSubmitMisbehaviour: {
    typeUrl: string;
    encode(message: MsgSubmitMisbehaviour, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgSubmitMisbehaviour;
    fromJSON(object: any): MsgSubmitMisbehaviour;
    toJSON(message: MsgSubmitMisbehaviour): unknown;
    fromPartial<I extends {
        clientId?: string | undefined;
        misbehaviour?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        signer?: string | undefined;
    } & {
        clientId?: string | undefined;
        misbehaviour?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["misbehaviour"], keyof Any>, never>) | undefined;
        signer?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgSubmitMisbehaviour>, never>>(object: I): MsgSubmitMisbehaviour;
};
export declare const MsgSubmitMisbehaviourResponse: {
    typeUrl: string;
    encode(_: MsgSubmitMisbehaviourResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgSubmitMisbehaviourResponse;
    fromJSON(_: any): MsgSubmitMisbehaviourResponse;
    toJSON(_: MsgSubmitMisbehaviourResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgSubmitMisbehaviourResponse;
};
/** Msg defines the ibc/client Msg service. */
export interface Msg {
    /** CreateClient defines a rpc handler method for MsgCreateClient. */
    CreateClient(request: MsgCreateClient): Promise<MsgCreateClientResponse>;
    /** UpdateClient defines a rpc handler method for MsgUpdateClient. */
    UpdateClient(request: MsgUpdateClient): Promise<MsgUpdateClientResponse>;
    /** UpgradeClient defines a rpc handler method for MsgUpgradeClient. */
    UpgradeClient(request: MsgUpgradeClient): Promise<MsgUpgradeClientResponse>;
    /** SubmitMisbehaviour defines a rpc handler method for MsgSubmitMisbehaviour. */
    SubmitMisbehaviour(request: MsgSubmitMisbehaviour): Promise<MsgSubmitMisbehaviourResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    CreateClient(request: MsgCreateClient): Promise<MsgCreateClientResponse>;
    UpdateClient(request: MsgUpdateClient): Promise<MsgUpdateClientResponse>;
    UpgradeClient(request: MsgUpgradeClient): Promise<MsgUpgradeClientResponse>;
    SubmitMisbehaviour(request: MsgSubmitMisbehaviour): Promise<MsgSubmitMisbehaviourResponse>;
}
