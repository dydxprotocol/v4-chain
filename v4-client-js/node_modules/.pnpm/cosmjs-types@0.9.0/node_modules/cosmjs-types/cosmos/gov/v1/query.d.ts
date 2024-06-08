import { ProposalStatus, Proposal, Vote, VotingParams, DepositParams, TallyParams, Params, Deposit, TallyResult } from "./gov";
import { PageRequest, PageResponse } from "../../base/query/v1beta1/pagination";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.gov.v1";
/** QueryProposalRequest is the request type for the Query/Proposal RPC method. */
export interface QueryProposalRequest {
    /** proposal_id defines the unique id of the proposal. */
    proposalId: bigint;
}
/** QueryProposalResponse is the response type for the Query/Proposal RPC method. */
export interface QueryProposalResponse {
    /** proposal is the requested governance proposal. */
    proposal?: Proposal;
}
/** QueryProposalsRequest is the request type for the Query/Proposals RPC method. */
export interface QueryProposalsRequest {
    /** proposal_status defines the status of the proposals. */
    proposalStatus: ProposalStatus;
    /** voter defines the voter address for the proposals. */
    voter: string;
    /** depositor defines the deposit addresses from the proposals. */
    depositor: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/**
 * QueryProposalsResponse is the response type for the Query/Proposals RPC
 * method.
 */
export interface QueryProposalsResponse {
    /** proposals defines all the requested governance proposals. */
    proposals: Proposal[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/** QueryVoteRequest is the request type for the Query/Vote RPC method. */
export interface QueryVoteRequest {
    /** proposal_id defines the unique id of the proposal. */
    proposalId: bigint;
    /** voter defines the voter address for the proposals. */
    voter: string;
}
/** QueryVoteResponse is the response type for the Query/Vote RPC method. */
export interface QueryVoteResponse {
    /** vote defines the queried vote. */
    vote?: Vote;
}
/** QueryVotesRequest is the request type for the Query/Votes RPC method. */
export interface QueryVotesRequest {
    /** proposal_id defines the unique id of the proposal. */
    proposalId: bigint;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/** QueryVotesResponse is the response type for the Query/Votes RPC method. */
export interface QueryVotesResponse {
    /** votes defines the queried votes. */
    votes: Vote[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/** QueryParamsRequest is the request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
    /**
     * params_type defines which parameters to query for, can be one of "voting",
     * "tallying" or "deposit".
     */
    paramsType: string;
}
/** QueryParamsResponse is the response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
    /**
     * Deprecated: Prefer to use `params` instead.
     * voting_params defines the parameters related to voting.
     */
    /** @deprecated */
    votingParams?: VotingParams;
    /**
     * Deprecated: Prefer to use `params` instead.
     * deposit_params defines the parameters related to deposit.
     */
    /** @deprecated */
    depositParams?: DepositParams;
    /**
     * Deprecated: Prefer to use `params` instead.
     * tally_params defines the parameters related to tally.
     */
    /** @deprecated */
    tallyParams?: TallyParams;
    /**
     * params defines all the paramaters of x/gov module.
     *
     * Since: cosmos-sdk 0.47
     */
    params?: Params;
}
/** QueryDepositRequest is the request type for the Query/Deposit RPC method. */
export interface QueryDepositRequest {
    /** proposal_id defines the unique id of the proposal. */
    proposalId: bigint;
    /** depositor defines the deposit addresses from the proposals. */
    depositor: string;
}
/** QueryDepositResponse is the response type for the Query/Deposit RPC method. */
export interface QueryDepositResponse {
    /** deposit defines the requested deposit. */
    deposit?: Deposit;
}
/** QueryDepositsRequest is the request type for the Query/Deposits RPC method. */
export interface QueryDepositsRequest {
    /** proposal_id defines the unique id of the proposal. */
    proposalId: bigint;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/** QueryDepositsResponse is the response type for the Query/Deposits RPC method. */
export interface QueryDepositsResponse {
    /** deposits defines the requested deposits. */
    deposits: Deposit[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/** QueryTallyResultRequest is the request type for the Query/Tally RPC method. */
export interface QueryTallyResultRequest {
    /** proposal_id defines the unique id of the proposal. */
    proposalId: bigint;
}
/** QueryTallyResultResponse is the response type for the Query/Tally RPC method. */
export interface QueryTallyResultResponse {
    /** tally defines the requested tally. */
    tally?: TallyResult;
}
export declare const QueryProposalRequest: {
    typeUrl: string;
    encode(message: QueryProposalRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryProposalRequest;
    fromJSON(object: any): QueryProposalRequest;
    toJSON(message: QueryProposalRequest): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
    } & {
        proposalId?: bigint | undefined;
    } & Record<Exclude<keyof I, "proposalId">, never>>(object: I): QueryProposalRequest;
};
export declare const QueryProposalResponse: {
    typeUrl: string;
    encode(message: QueryProposalResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryProposalResponse;
    fromJSON(object: any): QueryProposalResponse;
    toJSON(message: QueryProposalResponse): unknown;
    fromPartial<I extends {
        proposal?: {
            id?: bigint | undefined;
            messages?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
            status?: ProposalStatus | undefined;
            finalTallyResult?: {
                yesCount?: string | undefined;
                abstainCount?: string | undefined;
                noCount?: string | undefined;
                noWithVetoCount?: string | undefined;
            } | undefined;
            submitTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            depositEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalDeposit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            votingStartTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            votingEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            metadata?: string | undefined;
            title?: string | undefined;
            summary?: string | undefined;
            proposer?: string | undefined;
        } | undefined;
    } & {
        proposal?: ({
            id?: bigint | undefined;
            messages?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
            status?: ProposalStatus | undefined;
            finalTallyResult?: {
                yesCount?: string | undefined;
                abstainCount?: string | undefined;
                noCount?: string | undefined;
                noWithVetoCount?: string | undefined;
            } | undefined;
            submitTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            depositEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalDeposit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            votingStartTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            votingEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            metadata?: string | undefined;
            title?: string | undefined;
            summary?: string | undefined;
            proposer?: string | undefined;
        } & {
            id?: bigint | undefined;
            messages?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] & ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["proposal"]["messages"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["proposal"]["messages"], keyof {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[]>, never>) | undefined;
            status?: ProposalStatus | undefined;
            finalTallyResult?: ({
                yesCount?: string | undefined;
                abstainCount?: string | undefined;
                noCount?: string | undefined;
                noWithVetoCount?: string | undefined;
            } & {
                yesCount?: string | undefined;
                abstainCount?: string | undefined;
                noCount?: string | undefined;
                noWithVetoCount?: string | undefined;
            } & Record<Exclude<keyof I["proposal"]["finalTallyResult"], keyof TallyResult>, never>) | undefined;
            submitTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["proposal"]["submitTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            depositEndTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["proposal"]["depositEndTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            totalDeposit?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["proposal"]["totalDeposit"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["proposal"]["totalDeposit"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            votingStartTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["proposal"]["votingStartTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            votingEndTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["proposal"]["votingEndTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            metadata?: string | undefined;
            title?: string | undefined;
            summary?: string | undefined;
            proposer?: string | undefined;
        } & Record<Exclude<keyof I["proposal"], keyof Proposal>, never>) | undefined;
    } & Record<Exclude<keyof I, "proposal">, never>>(object: I): QueryProposalResponse;
};
export declare const QueryProposalsRequest: {
    typeUrl: string;
    encode(message: QueryProposalsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryProposalsRequest;
    fromJSON(object: any): QueryProposalsRequest;
    toJSON(message: QueryProposalsRequest): unknown;
    fromPartial<I extends {
        proposalStatus?: ProposalStatus | undefined;
        voter?: string | undefined;
        depositor?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        proposalStatus?: ProposalStatus | undefined;
        voter?: string | undefined;
        depositor?: string | undefined;
        pagination?: ({
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageRequest>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryProposalsRequest>, never>>(object: I): QueryProposalsRequest;
};
export declare const QueryProposalsResponse: {
    typeUrl: string;
    encode(message: QueryProposalsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryProposalsResponse;
    fromJSON(object: any): QueryProposalsResponse;
    toJSON(message: QueryProposalsResponse): unknown;
    fromPartial<I extends {
        proposals?: {
            id?: bigint | undefined;
            messages?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
            status?: ProposalStatus | undefined;
            finalTallyResult?: {
                yesCount?: string | undefined;
                abstainCount?: string | undefined;
                noCount?: string | undefined;
                noWithVetoCount?: string | undefined;
            } | undefined;
            submitTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            depositEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalDeposit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            votingStartTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            votingEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            metadata?: string | undefined;
            title?: string | undefined;
            summary?: string | undefined;
            proposer?: string | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        proposals?: ({
            id?: bigint | undefined;
            messages?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
            status?: ProposalStatus | undefined;
            finalTallyResult?: {
                yesCount?: string | undefined;
                abstainCount?: string | undefined;
                noCount?: string | undefined;
                noWithVetoCount?: string | undefined;
            } | undefined;
            submitTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            depositEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalDeposit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            votingStartTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            votingEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            metadata?: string | undefined;
            title?: string | undefined;
            summary?: string | undefined;
            proposer?: string | undefined;
        }[] & ({
            id?: bigint | undefined;
            messages?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
            status?: ProposalStatus | undefined;
            finalTallyResult?: {
                yesCount?: string | undefined;
                abstainCount?: string | undefined;
                noCount?: string | undefined;
                noWithVetoCount?: string | undefined;
            } | undefined;
            submitTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            depositEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalDeposit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            votingStartTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            votingEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            metadata?: string | undefined;
            title?: string | undefined;
            summary?: string | undefined;
            proposer?: string | undefined;
        } & {
            id?: bigint | undefined;
            messages?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] & ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["proposals"][number]["messages"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["proposals"][number]["messages"], keyof {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[]>, never>) | undefined;
            status?: ProposalStatus | undefined;
            finalTallyResult?: ({
                yesCount?: string | undefined;
                abstainCount?: string | undefined;
                noCount?: string | undefined;
                noWithVetoCount?: string | undefined;
            } & {
                yesCount?: string | undefined;
                abstainCount?: string | undefined;
                noCount?: string | undefined;
                noWithVetoCount?: string | undefined;
            } & Record<Exclude<keyof I["proposals"][number]["finalTallyResult"], keyof TallyResult>, never>) | undefined;
            submitTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["proposals"][number]["submitTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            depositEndTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["proposals"][number]["depositEndTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            totalDeposit?: ({
                denom?: string | undefined;
                amount?: string | undefined;
            }[] & ({
                denom?: string | undefined;
                amount?: string | undefined;
            } & {
                denom?: string | undefined;
                amount?: string | undefined;
            } & Record<Exclude<keyof I["proposals"][number]["totalDeposit"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["proposals"][number]["totalDeposit"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            votingStartTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["proposals"][number]["votingStartTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            votingEndTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["proposals"][number]["votingEndTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            metadata?: string | undefined;
            title?: string | undefined;
            summary?: string | undefined;
            proposer?: string | undefined;
        } & Record<Exclude<keyof I["proposals"][number], keyof Proposal>, never>)[] & Record<Exclude<keyof I["proposals"], keyof {
            id?: bigint | undefined;
            messages?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
            status?: ProposalStatus | undefined;
            finalTallyResult?: {
                yesCount?: string | undefined;
                abstainCount?: string | undefined;
                noCount?: string | undefined;
                noWithVetoCount?: string | undefined;
            } | undefined;
            submitTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            depositEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            totalDeposit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            votingStartTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            votingEndTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            metadata?: string | undefined;
            title?: string | undefined;
            summary?: string | undefined;
            proposer?: string | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryProposalsResponse>, never>>(object: I): QueryProposalsResponse;
};
export declare const QueryVoteRequest: {
    typeUrl: string;
    encode(message: QueryVoteRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryVoteRequest;
    fromJSON(object: any): QueryVoteRequest;
    toJSON(message: QueryVoteRequest): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
        voter?: string | undefined;
    } & {
        proposalId?: bigint | undefined;
        voter?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryVoteRequest>, never>>(object: I): QueryVoteRequest;
};
export declare const QueryVoteResponse: {
    typeUrl: string;
    encode(message: QueryVoteResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryVoteResponse;
    fromJSON(object: any): QueryVoteResponse;
    toJSON(message: QueryVoteResponse): unknown;
    fromPartial<I extends {
        vote?: {
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            options?: {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[] | undefined;
            metadata?: string | undefined;
        } | undefined;
    } & {
        vote?: ({
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            options?: {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[] | undefined;
            metadata?: string | undefined;
        } & {
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            options?: ({
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[] & ({
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            } & {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            } & Record<Exclude<keyof I["vote"]["options"][number], keyof import("./gov").WeightedVoteOption>, never>)[] & Record<Exclude<keyof I["vote"]["options"], keyof {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[]>, never>) | undefined;
            metadata?: string | undefined;
        } & Record<Exclude<keyof I["vote"], keyof Vote>, never>) | undefined;
    } & Record<Exclude<keyof I, "vote">, never>>(object: I): QueryVoteResponse;
};
export declare const QueryVotesRequest: {
    typeUrl: string;
    encode(message: QueryVotesRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryVotesRequest;
    fromJSON(object: any): QueryVotesRequest;
    toJSON(message: QueryVotesRequest): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        proposalId?: bigint | undefined;
        pagination?: ({
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageRequest>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryVotesRequest>, never>>(object: I): QueryVotesRequest;
};
export declare const QueryVotesResponse: {
    typeUrl: string;
    encode(message: QueryVotesResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryVotesResponse;
    fromJSON(object: any): QueryVotesResponse;
    toJSON(message: QueryVotesResponse): unknown;
    fromPartial<I extends {
        votes?: {
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            options?: {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[] | undefined;
            metadata?: string | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        votes?: ({
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            options?: {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[] | undefined;
            metadata?: string | undefined;
        }[] & ({
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            options?: {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[] | undefined;
            metadata?: string | undefined;
        } & {
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            options?: ({
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[] & ({
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            } & {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            } & Record<Exclude<keyof I["votes"][number]["options"][number], keyof import("./gov").WeightedVoteOption>, never>)[] & Record<Exclude<keyof I["votes"][number]["options"], keyof {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[]>, never>) | undefined;
            metadata?: string | undefined;
        } & Record<Exclude<keyof I["votes"][number], keyof Vote>, never>)[] & Record<Exclude<keyof I["votes"], keyof {
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            options?: {
                option?: import("./gov").VoteOption | undefined;
                weight?: string | undefined;
            }[] | undefined;
            metadata?: string | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryVotesResponse>, never>>(object: I): QueryVotesResponse;
};
export declare const QueryParamsRequest: {
    typeUrl: string;
    encode(message: QueryParamsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryParamsRequest;
    fromJSON(object: any): QueryParamsRequest;
    toJSON(message: QueryParamsRequest): unknown;
    fromPartial<I extends {
        paramsType?: string | undefined;
    } & {
        paramsType?: string | undefined;
    } & Record<Exclude<keyof I, "paramsType">, never>>(object: I): QueryParamsRequest;
};
export declare const QueryParamsResponse: {
    typeUrl: string;
    encode(message: QueryParamsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryParamsResponse;
    fromJSON(object: any): QueryParamsResponse;
    toJSON(message: QueryParamsResponse): unknown;
    fromPartial<I extends {
        votingParams?: {
            votingPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } | undefined;
        depositParams?: {
            minDeposit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            maxDepositPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } | undefined;
        tallyParams?: {
            quorum?: string | undefined;
            threshold?: string | undefined;
            vetoThreshold?: string | undefined;
        } | undefined;
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
        votingParams?: ({
            votingPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } & {
            votingPeriod?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["votingParams"]["votingPeriod"], keyof import("../../../google/protobuf/duration").Duration>, never>) | undefined;
        } & Record<Exclude<keyof I["votingParams"], "votingPeriod">, never>) | undefined;
        depositParams?: ({
            minDeposit?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
            maxDepositPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
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
            } & Record<Exclude<keyof I["depositParams"]["minDeposit"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["depositParams"]["minDeposit"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
            maxDepositPeriod?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["depositParams"]["maxDepositPeriod"], keyof import("../../../google/protobuf/duration").Duration>, never>) | undefined;
        } & Record<Exclude<keyof I["depositParams"], keyof DepositParams>, never>) | undefined;
        tallyParams?: ({
            quorum?: string | undefined;
            threshold?: string | undefined;
            vetoThreshold?: string | undefined;
        } & {
            quorum?: string | undefined;
            threshold?: string | undefined;
            vetoThreshold?: string | undefined;
        } & Record<Exclude<keyof I["tallyParams"], keyof TallyParams>, never>) | undefined;
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
            } & Record<Exclude<keyof I["params"]["minDeposit"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["params"]["minDeposit"], keyof {
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
    } & Record<Exclude<keyof I, keyof QueryParamsResponse>, never>>(object: I): QueryParamsResponse;
};
export declare const QueryDepositRequest: {
    typeUrl: string;
    encode(message: QueryDepositRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDepositRequest;
    fromJSON(object: any): QueryDepositRequest;
    toJSON(message: QueryDepositRequest): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
        depositor?: string | undefined;
    } & {
        proposalId?: bigint | undefined;
        depositor?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryDepositRequest>, never>>(object: I): QueryDepositRequest;
};
export declare const QueryDepositResponse: {
    typeUrl: string;
    encode(message: QueryDepositResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDepositResponse;
    fromJSON(object: any): QueryDepositResponse;
    toJSON(message: QueryDepositResponse): unknown;
    fromPartial<I extends {
        deposit?: {
            proposalId?: bigint | undefined;
            depositor?: string | undefined;
            amount?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
        } | undefined;
    } & {
        deposit?: ({
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
            } & Record<Exclude<keyof I["deposit"]["amount"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["deposit"]["amount"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["deposit"], keyof Deposit>, never>) | undefined;
    } & Record<Exclude<keyof I, "deposit">, never>>(object: I): QueryDepositResponse;
};
export declare const QueryDepositsRequest: {
    typeUrl: string;
    encode(message: QueryDepositsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDepositsRequest;
    fromJSON(object: any): QueryDepositsRequest;
    toJSON(message: QueryDepositsRequest): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        proposalId?: bigint | undefined;
        pagination?: ({
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageRequest>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryDepositsRequest>, never>>(object: I): QueryDepositsRequest;
};
export declare const QueryDepositsResponse: {
    typeUrl: string;
    encode(message: QueryDepositsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryDepositsResponse;
    fromJSON(object: any): QueryDepositsResponse;
    toJSON(message: QueryDepositsResponse): unknown;
    fromPartial<I extends {
        deposits?: {
            proposalId?: bigint | undefined;
            depositor?: string | undefined;
            amount?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        deposits?: ({
            proposalId?: bigint | undefined;
            depositor?: string | undefined;
            amount?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
        }[] & ({
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
            } & Record<Exclude<keyof I["deposits"][number]["amount"][number], keyof import("../../base/v1beta1/coin").Coin>, never>)[] & Record<Exclude<keyof I["deposits"][number]["amount"], keyof {
                denom?: string | undefined;
                amount?: string | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["deposits"][number], keyof Deposit>, never>)[] & Record<Exclude<keyof I["deposits"], keyof {
            proposalId?: bigint | undefined;
            depositor?: string | undefined;
            amount?: {
                denom?: string | undefined;
                amount?: string | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryDepositsResponse>, never>>(object: I): QueryDepositsResponse;
};
export declare const QueryTallyResultRequest: {
    typeUrl: string;
    encode(message: QueryTallyResultRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryTallyResultRequest;
    fromJSON(object: any): QueryTallyResultRequest;
    toJSON(message: QueryTallyResultRequest): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
    } & {
        proposalId?: bigint | undefined;
    } & Record<Exclude<keyof I, "proposalId">, never>>(object: I): QueryTallyResultRequest;
};
export declare const QueryTallyResultResponse: {
    typeUrl: string;
    encode(message: QueryTallyResultResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryTallyResultResponse;
    fromJSON(object: any): QueryTallyResultResponse;
    toJSON(message: QueryTallyResultResponse): unknown;
    fromPartial<I extends {
        tally?: {
            yesCount?: string | undefined;
            abstainCount?: string | undefined;
            noCount?: string | undefined;
            noWithVetoCount?: string | undefined;
        } | undefined;
    } & {
        tally?: ({
            yesCount?: string | undefined;
            abstainCount?: string | undefined;
            noCount?: string | undefined;
            noWithVetoCount?: string | undefined;
        } & {
            yesCount?: string | undefined;
            abstainCount?: string | undefined;
            noCount?: string | undefined;
            noWithVetoCount?: string | undefined;
        } & Record<Exclude<keyof I["tally"], keyof TallyResult>, never>) | undefined;
    } & Record<Exclude<keyof I, "tally">, never>>(object: I): QueryTallyResultResponse;
};
/** Query defines the gRPC querier service for gov module */
export interface Query {
    /** Proposal queries proposal details based on ProposalID. */
    Proposal(request: QueryProposalRequest): Promise<QueryProposalResponse>;
    /** Proposals queries all proposals based on given status. */
    Proposals(request: QueryProposalsRequest): Promise<QueryProposalsResponse>;
    /** Vote queries voted information based on proposalID, voterAddr. */
    Vote(request: QueryVoteRequest): Promise<QueryVoteResponse>;
    /** Votes queries votes of a given proposal. */
    Votes(request: QueryVotesRequest): Promise<QueryVotesResponse>;
    /** Params queries all parameters of the gov module. */
    Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
    /** Deposit queries single deposit information based proposalID, depositAddr. */
    Deposit(request: QueryDepositRequest): Promise<QueryDepositResponse>;
    /** Deposits queries all deposits of a single proposal. */
    Deposits(request: QueryDepositsRequest): Promise<QueryDepositsResponse>;
    /** TallyResult queries the tally of a proposal vote. */
    TallyResult(request: QueryTallyResultRequest): Promise<QueryTallyResultResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    Proposal(request: QueryProposalRequest): Promise<QueryProposalResponse>;
    Proposals(request: QueryProposalsRequest): Promise<QueryProposalsResponse>;
    Vote(request: QueryVoteRequest): Promise<QueryVoteResponse>;
    Votes(request: QueryVotesRequest): Promise<QueryVotesResponse>;
    Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
    Deposit(request: QueryDepositRequest): Promise<QueryDepositResponse>;
    Deposits(request: QueryDepositsRequest): Promise<QueryDepositsResponse>;
    TallyResult(request: QueryTallyResultRequest): Promise<QueryTallyResultResponse>;
}
