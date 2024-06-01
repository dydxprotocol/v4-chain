import { Any } from "../../../google/protobuf/any";
import { Coin } from "../../base/v1beta1/coin";
import { VoteOption, WeightedVoteOption, Params } from "./gov";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.gov.v1";
/**
 * MsgSubmitProposal defines an sdk.Msg type that supports submitting arbitrary
 * proposal Content.
 */
export interface MsgSubmitProposal {
    /** messages are the arbitrary messages to be executed if proposal passes. */
    messages: Any[];
    /** initial_deposit is the deposit value that must be paid at proposal submission. */
    initialDeposit: Coin[];
    /** proposer is the account address of the proposer. */
    proposer: string;
    /** metadata is any arbitrary metadata attached to the proposal. */
    metadata: string;
    /**
     * title is the title of the proposal.
     *
     * Since: cosmos-sdk 0.47
     */
    title: string;
    /**
     * summary is the summary of the proposal
     *
     * Since: cosmos-sdk 0.47
     */
    summary: string;
}
/** MsgSubmitProposalResponse defines the Msg/SubmitProposal response type. */
export interface MsgSubmitProposalResponse {
    /** proposal_id defines the unique id of the proposal. */
    proposalId: bigint;
}
/**
 * MsgExecLegacyContent is used to wrap the legacy content field into a message.
 * This ensures backwards compatibility with v1beta1.MsgSubmitProposal.
 */
export interface MsgExecLegacyContent {
    /** content is the proposal's content. */
    content?: Any;
    /** authority must be the gov module address. */
    authority: string;
}
/** MsgExecLegacyContentResponse defines the Msg/ExecLegacyContent response type. */
export interface MsgExecLegacyContentResponse {
}
/** MsgVote defines a message to cast a vote. */
export interface MsgVote {
    /** proposal_id defines the unique id of the proposal. */
    proposalId: bigint;
    /** voter is the voter address for the proposal. */
    voter: string;
    /** option defines the vote option. */
    option: VoteOption;
    /** metadata is any arbitrary metadata attached to the Vote. */
    metadata: string;
}
/** MsgVoteResponse defines the Msg/Vote response type. */
export interface MsgVoteResponse {
}
/** MsgVoteWeighted defines a message to cast a vote. */
export interface MsgVoteWeighted {
    /** proposal_id defines the unique id of the proposal. */
    proposalId: bigint;
    /** voter is the voter address for the proposal. */
    voter: string;
    /** options defines the weighted vote options. */
    options: WeightedVoteOption[];
    /** metadata is any arbitrary metadata attached to the VoteWeighted. */
    metadata: string;
}
/** MsgVoteWeightedResponse defines the Msg/VoteWeighted response type. */
export interface MsgVoteWeightedResponse {
}
/** MsgDeposit defines a message to submit a deposit to an existing proposal. */
export interface MsgDeposit {
    /** proposal_id defines the unique id of the proposal. */
    proposalId: bigint;
    /** depositor defines the deposit addresses from the proposals. */
    depositor: string;
    /** amount to be deposited by depositor. */
    amount: Coin[];
}
/** MsgDepositResponse defines the Msg/Deposit response type. */
export interface MsgDepositResponse {
}
/**
 * MsgUpdateParams is the Msg/UpdateParams request type.
 *
 * Since: cosmos-sdk 0.47
 */
export interface MsgUpdateParams {
    /** authority is the address that controls the module (defaults to x/gov unless overwritten). */
    authority: string;
    /**
     * params defines the x/gov parameters to update.
     *
     * NOTE: All parameters must be supplied.
     */
    params: Params;
}
/**
 * MsgUpdateParamsResponse defines the response structure for executing a
 * MsgUpdateParams message.
 *
 * Since: cosmos-sdk 0.47
 */
export interface MsgUpdateParamsResponse {
}
export declare const MsgSubmitProposal: {
    typeUrl: string;
    encode(message: MsgSubmitProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgSubmitProposal;
    fromJSON(object: any): MsgSubmitProposal;
    toJSON(message: MsgSubmitProposal): unknown;
    fromPartial<I extends {
        messages?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[] | undefined;
        initialDeposit?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
        proposer?: string | undefined;
        metadata?: string | undefined;
        title?: string | undefined;
        summary?: string | undefined;
    } & {
        messages?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[] & ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["messages"][number], keyof Any>, never>)[] & Record<Exclude<keyof I["messages"], keyof {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[]>, never>) | undefined;
        initialDeposit?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["initialDeposit"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["initialDeposit"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
        proposer?: string | undefined;
        metadata?: string | undefined;
        title?: string | undefined;
        summary?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgSubmitProposal>, never>>(object: I): MsgSubmitProposal;
};
export declare const MsgSubmitProposalResponse: {
    typeUrl: string;
    encode(message: MsgSubmitProposalResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgSubmitProposalResponse;
    fromJSON(object: any): MsgSubmitProposalResponse;
    toJSON(message: MsgSubmitProposalResponse): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
    } & {
        proposalId?: bigint | undefined;
    } & Record<Exclude<keyof I, "proposalId">, never>>(object: I): MsgSubmitProposalResponse;
};
export declare const MsgExecLegacyContent: {
    typeUrl: string;
    encode(message: MsgExecLegacyContent, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgExecLegacyContent;
    fromJSON(object: any): MsgExecLegacyContent;
    toJSON(message: MsgExecLegacyContent): unknown;
    fromPartial<I extends {
        content?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        authority?: string | undefined;
    } & {
        content?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["content"], keyof Any>, never>) | undefined;
        authority?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgExecLegacyContent>, never>>(object: I): MsgExecLegacyContent;
};
export declare const MsgExecLegacyContentResponse: {
    typeUrl: string;
    encode(_: MsgExecLegacyContentResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgExecLegacyContentResponse;
    fromJSON(_: any): MsgExecLegacyContentResponse;
    toJSON(_: MsgExecLegacyContentResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgExecLegacyContentResponse;
};
export declare const MsgVote: {
    typeUrl: string;
    encode(message: MsgVote, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgVote;
    fromJSON(object: any): MsgVote;
    toJSON(message: MsgVote): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
        voter?: string | undefined;
        option?: VoteOption | undefined;
        metadata?: string | undefined;
    } & {
        proposalId?: bigint | undefined;
        voter?: string | undefined;
        option?: VoteOption | undefined;
        metadata?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgVote>, never>>(object: I): MsgVote;
};
export declare const MsgVoteResponse: {
    typeUrl: string;
    encode(_: MsgVoteResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgVoteResponse;
    fromJSON(_: any): MsgVoteResponse;
    toJSON(_: MsgVoteResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgVoteResponse;
};
export declare const MsgVoteWeighted: {
    typeUrl: string;
    encode(message: MsgVoteWeighted, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgVoteWeighted;
    fromJSON(object: any): MsgVoteWeighted;
    toJSON(message: MsgVoteWeighted): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
        voter?: string | undefined;
        options?: {
            option?: VoteOption | undefined;
            weight?: string | undefined;
        }[] | undefined;
        metadata?: string | undefined;
    } & {
        proposalId?: bigint | undefined;
        voter?: string | undefined;
        options?: ({
            option?: VoteOption | undefined;
            weight?: string | undefined;
        }[] & ({
            option?: VoteOption | undefined;
            weight?: string | undefined;
        } & {
            option?: VoteOption | undefined;
            weight?: string | undefined;
        } & Record<Exclude<keyof I["options"][number], keyof WeightedVoteOption>, never>)[] & Record<Exclude<keyof I["options"], keyof {
            option?: VoteOption | undefined;
            weight?: string | undefined;
        }[]>, never>) | undefined;
        metadata?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgVoteWeighted>, never>>(object: I): MsgVoteWeighted;
};
export declare const MsgVoteWeightedResponse: {
    typeUrl: string;
    encode(_: MsgVoteWeightedResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgVoteWeightedResponse;
    fromJSON(_: any): MsgVoteWeightedResponse;
    toJSON(_: MsgVoteWeightedResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgVoteWeightedResponse;
};
export declare const MsgDeposit: {
    typeUrl: string;
    encode(message: MsgDeposit, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgDeposit;
    fromJSON(object: any): MsgDeposit;
    toJSON(message: MsgDeposit): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
        depositor?: string | undefined;
        amount?: {
            denom?: string | undefined;
            amount?: string | undefined;
        }[] | undefined;
    } & {
        proposalId?: bigint | undefined;
        depositor?: string | undefined;
        amount?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["amount"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["amount"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof MsgDeposit>, never>>(object: I): MsgDeposit;
};
export declare const MsgDepositResponse: {
    typeUrl: string;
    encode(_: MsgDepositResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgDepositResponse;
    fromJSON(_: any): MsgDepositResponse;
    toJSON(_: MsgDepositResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgDepositResponse;
};
export declare const MsgUpdateParams: {
    typeUrl: string;
    encode(message: MsgUpdateParams, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateParams;
    fromJSON(object: any): MsgUpdateParams;
    toJSON(message: MsgUpdateParams): unknown;
    fromPartial<I extends {
        authority?: string | undefined;
        params?: {
            minDeposit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            maxDepositPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            votingPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            quorum?: string | undefined;
            threshold?: string | undefined;
            vetoThreshold?: string | undefined;
            minInitialDepositRatio?: string | undefined;
            burnVoteQuorum?: boolean | undefined;
            burnProposalDepositPrevote?: boolean | undefined;
            burnVoteVeto?: boolean | undefined;
        } | undefined;
    } & {
        authority?: string | undefined;
        params?: ({
            minDeposit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            maxDepositPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            votingPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            quorum?: string | undefined;
            threshold?: string | undefined;
            vetoThreshold?: string | undefined;
            minInitialDepositRatio?: string | undefined;
            burnVoteQuorum?: boolean | undefined;
            burnProposalDepositPrevote?: boolean | undefined;
            burnVoteVeto?: boolean | undefined;
        } & {
            minDeposit?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["params"]["minDeposit"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["params"]["minDeposit"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            maxDepositPeriod?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["params"]["maxDepositPeriod"], keyof import("../../../google/protobuf/duration").Duration>, never>) | undefined;
            votingPeriod?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["params"]["votingPeriod"], keyof import("../../../google/protobuf/duration").Duration>, never>) | undefined;
            quorum?: string | undefined;
            threshold?: string | undefined;
            vetoThreshold?: string | undefined;
            minInitialDepositRatio?: string | undefined;
            burnVoteQuorum?: boolean | undefined;
            burnProposalDepositPrevote?: boolean | undefined;
            burnVoteVeto?: boolean | undefined;
        } & Record<Exclude<keyof I["params"], keyof Params>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof MsgUpdateParams>, never>>(object: I): MsgUpdateParams;
};
export declare const MsgUpdateParamsResponse: {
    typeUrl: string;
    encode(_: MsgUpdateParamsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateParamsResponse;
    fromJSON(_: any): MsgUpdateParamsResponse;
    toJSON(_: MsgUpdateParamsResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgUpdateParamsResponse;
};
/** Msg defines the gov Msg service. */
export interface Msg {
    /** SubmitProposal defines a method to create new proposal given the messages. */
    SubmitProposal(request: MsgSubmitProposal): Promise<MsgSubmitProposalResponse>;
    /**
     * ExecLegacyContent defines a Msg to be in included in a MsgSubmitProposal
     * to execute a legacy content-based proposal.
     */
    ExecLegacyContent(request: MsgExecLegacyContent): Promise<MsgExecLegacyContentResponse>;
    /** Vote defines a method to add a vote on a specific proposal. */
    Vote(request: MsgVote): Promise<MsgVoteResponse>;
    /** VoteWeighted defines a method to add a weighted vote on a specific proposal. */
    VoteWeighted(request: MsgVoteWeighted): Promise<MsgVoteWeightedResponse>;
    /** Deposit defines a method to add deposit on a specific proposal. */
    Deposit(request: MsgDeposit): Promise<MsgDepositResponse>;
    /**
     * UpdateParams defines a governance operation for updating the x/gov module
     * parameters. The authority is defined in the keeper.
     *
     * Since: cosmos-sdk 0.47
     */
    UpdateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    SubmitProposal(request: MsgSubmitProposal): Promise<MsgSubmitProposalResponse>;
    ExecLegacyContent(request: MsgExecLegacyContent): Promise<MsgExecLegacyContentResponse>;
    Vote(request: MsgVote): Promise<MsgVoteResponse>;
    VoteWeighted(request: MsgVoteWeighted): Promise<MsgVoteWeightedResponse>;
    Deposit(request: MsgDeposit): Promise<MsgDepositResponse>;
    UpdateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
}
