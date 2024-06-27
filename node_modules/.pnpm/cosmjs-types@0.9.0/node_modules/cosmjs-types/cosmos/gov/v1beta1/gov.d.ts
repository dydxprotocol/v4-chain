import { Coin } from "../../base/v1beta1/coin";
import { Any } from "../../../google/protobuf/any";
import { Timestamp } from "../../../google/protobuf/timestamp";
import { Duration } from "../../../google/protobuf/duration";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.gov.v1beta1";
/** VoteOption enumerates the valid vote options for a given governance proposal. */
export declare enum VoteOption {
    /** VOTE_OPTION_UNSPECIFIED - VOTE_OPTION_UNSPECIFIED defines a no-op vote option. */
    VOTE_OPTION_UNSPECIFIED = 0,
    /** VOTE_OPTION_YES - VOTE_OPTION_YES defines a yes vote option. */
    VOTE_OPTION_YES = 1,
    /** VOTE_OPTION_ABSTAIN - VOTE_OPTION_ABSTAIN defines an abstain vote option. */
    VOTE_OPTION_ABSTAIN = 2,
    /** VOTE_OPTION_NO - VOTE_OPTION_NO defines a no vote option. */
    VOTE_OPTION_NO = 3,
    /** VOTE_OPTION_NO_WITH_VETO - VOTE_OPTION_NO_WITH_VETO defines a no with veto vote option. */
    VOTE_OPTION_NO_WITH_VETO = 4,
    UNRECOGNIZED = -1
}
export declare function voteOptionFromJSON(object: any): VoteOption;
export declare function voteOptionToJSON(object: VoteOption): string;
/** ProposalStatus enumerates the valid statuses of a proposal. */
export declare enum ProposalStatus {
    /** PROPOSAL_STATUS_UNSPECIFIED - PROPOSAL_STATUS_UNSPECIFIED defines the default proposal status. */
    PROPOSAL_STATUS_UNSPECIFIED = 0,
    /**
     * PROPOSAL_STATUS_DEPOSIT_PERIOD - PROPOSAL_STATUS_DEPOSIT_PERIOD defines a proposal status during the deposit
     * period.
     */
    PROPOSAL_STATUS_DEPOSIT_PERIOD = 1,
    /**
     * PROPOSAL_STATUS_VOTING_PERIOD - PROPOSAL_STATUS_VOTING_PERIOD defines a proposal status during the voting
     * period.
     */
    PROPOSAL_STATUS_VOTING_PERIOD = 2,
    /**
     * PROPOSAL_STATUS_PASSED - PROPOSAL_STATUS_PASSED defines a proposal status of a proposal that has
     * passed.
     */
    PROPOSAL_STATUS_PASSED = 3,
    /**
     * PROPOSAL_STATUS_REJECTED - PROPOSAL_STATUS_REJECTED defines a proposal status of a proposal that has
     * been rejected.
     */
    PROPOSAL_STATUS_REJECTED = 4,
    /**
     * PROPOSAL_STATUS_FAILED - PROPOSAL_STATUS_FAILED defines a proposal status of a proposal that has
     * failed.
     */
    PROPOSAL_STATUS_FAILED = 5,
    UNRECOGNIZED = -1
}
export declare function proposalStatusFromJSON(object: any): ProposalStatus;
export declare function proposalStatusToJSON(object: ProposalStatus): string;
/**
 * WeightedVoteOption defines a unit of vote for vote split.
 *
 * Since: cosmos-sdk 0.43
 */
export interface WeightedVoteOption {
    /** option defines the valid vote options, it must not contain duplicate vote options. */
    option: VoteOption;
    /** weight is the vote weight associated with the vote option. */
    weight: string;
}
/**
 * TextProposal defines a standard text proposal whose changes need to be
 * manually updated in case of approval.
 */
export interface TextProposal {
    /** title of the proposal. */
    title: string;
    /** description associated with the proposal. */
    description: string;
}
/**
 * Deposit defines an amount deposited by an account address to an active
 * proposal.
 */
export interface Deposit {
    /** proposal_id defines the unique id of the proposal. */
    proposalId: bigint;
    /** depositor defines the deposit addresses from the proposals. */
    depositor: string;
    /** amount to be deposited by depositor. */
    amount: Coin[];
}
/** Proposal defines the core field members of a governance proposal. */
export interface Proposal {
    /** proposal_id defines the unique id of the proposal. */
    proposalId: bigint;
    /** content is the proposal's content. */
    content?: Any;
    /** status defines the proposal status. */
    status: ProposalStatus;
    /**
     * final_tally_result is the final tally result of the proposal. When
     * querying a proposal via gRPC, this field is not populated until the
     * proposal's voting period has ended.
     */
    finalTallyResult: TallyResult;
    /** submit_time is the time of proposal submission. */
    submitTime: Timestamp;
    /** deposit_end_time is the end time for deposition. */
    depositEndTime: Timestamp;
    /** total_deposit is the total deposit on the proposal. */
    totalDeposit: Coin[];
    /** voting_start_time is the starting time to vote on a proposal. */
    votingStartTime: Timestamp;
    /** voting_end_time is the end time of voting on a proposal. */
    votingEndTime: Timestamp;
}
/** TallyResult defines a standard tally for a governance proposal. */
export interface TallyResult {
    /** yes is the number of yes votes on a proposal. */
    yes: string;
    /** abstain is the number of abstain votes on a proposal. */
    abstain: string;
    /** no is the number of no votes on a proposal. */
    no: string;
    /** no_with_veto is the number of no with veto votes on a proposal. */
    noWithVeto: string;
}
/**
 * Vote defines a vote on a governance proposal.
 * A Vote consists of a proposal ID, the voter, and the vote option.
 */
export interface Vote {
    /** proposal_id defines the unique id of the proposal. */
    proposalId: bigint;
    /** voter is the voter address of the proposal. */
    voter: string;
    /**
     * Deprecated: Prefer to use `options` instead. This field is set in queries
     * if and only if `len(options) == 1` and that option has weight 1. In all
     * other cases, this field will default to VOTE_OPTION_UNSPECIFIED.
     */
    /** @deprecated */
    option: VoteOption;
    /**
     * options is the weighted vote options.
     *
     * Since: cosmos-sdk 0.43
     */
    options: WeightedVoteOption[];
}
/** DepositParams defines the params for deposits on governance proposals. */
export interface DepositParams {
    /** Minimum deposit for a proposal to enter voting period. */
    minDeposit: Coin[];
    /**
     * Maximum period for Atom holders to deposit on a proposal. Initial value: 2
     * months.
     */
    maxDepositPeriod: Duration;
}
/** VotingParams defines the params for voting on governance proposals. */
export interface VotingParams {
    /** Duration of the voting period. */
    votingPeriod: Duration;
}
/** TallyParams defines the params for tallying votes on governance proposals. */
export interface TallyParams {
    /**
     * Minimum percentage of total stake needed to vote for a result to be
     * considered valid.
     */
    quorum: Uint8Array;
    /** Minimum proportion of Yes votes for proposal to pass. Default value: 0.5. */
    threshold: Uint8Array;
    /**
     * Minimum value of Veto votes to Total votes ratio for proposal to be
     * vetoed. Default value: 1/3.
     */
    vetoThreshold: Uint8Array;
}
export declare const WeightedVoteOption: {
    typeUrl: string;
    encode(message: WeightedVoteOption, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): WeightedVoteOption;
    fromJSON(object: any): WeightedVoteOption;
    toJSON(message: WeightedVoteOption): unknown;
    fromPartial<I extends {
        option?: VoteOption | undefined;
        weight?: string | undefined;
    } & {
        option?: VoteOption | undefined;
        weight?: string | undefined;
    } & Record<Exclude<keyof I, keyof WeightedVoteOption>, never>>(object: I): WeightedVoteOption;
};
export declare const TextProposal: {
    typeUrl: string;
    encode(message: TextProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): TextProposal;
    fromJSON(object: any): TextProposal;
    toJSON(message: TextProposal): unknown;
    fromPartial<I extends {
        title?: string | undefined;
        description?: string | undefined;
    } & {
        title?: string | undefined;
        description?: string | undefined;
    } & Record<Exclude<keyof I, keyof TextProposal>, never>>(object: I): TextProposal;
};
export declare const Deposit: {
    typeUrl: string;
    encode(message: Deposit, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Deposit;
    fromJSON(object: any): Deposit;
    toJSON(message: Deposit): unknown;
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
    } & Record<Exclude<keyof I, keyof Deposit>, never>>(object: I): Deposit;
};
export declare const Proposal: {
    typeUrl: string;
    encode(message: Proposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Proposal;
    fromJSON(object: any): Proposal;
    toJSON(message: Proposal): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
        content?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        status?: ProposalStatus | undefined;
        finalTallyResult?: {
            yes?: string | undefined;
            abstain?: string | undefined;
            no?: string | undefined;
            noWithVeto?: string | undefined;
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
    } & {
        proposalId?: bigint | undefined;
        content?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["content"], keyof Any>, never>) | undefined;
        status?: ProposalStatus | undefined;
        finalTallyResult?: ({
            yes?: string | undefined;
            abstain?: string | undefined;
            no?: string | undefined;
            noWithVeto?: string | undefined;
        } & {
            yes?: string | undefined;
            abstain?: string | undefined;
            no?: string | undefined;
            noWithVeto?: string | undefined;
        } & Record<Exclude<keyof I["finalTallyResult"], keyof TallyResult>, never>) | undefined;
        submitTime?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["submitTime"], keyof Timestamp>, never>) | undefined;
        depositEndTime?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["depositEndTime"], keyof Timestamp>, never>) | undefined;
        totalDeposit?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        }[] & ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["totalDeposit"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["totalDeposit"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
        votingStartTime?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["votingStartTime"], keyof Timestamp>, never>) | undefined;
        votingEndTime?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["votingEndTime"], keyof Timestamp>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Proposal>, never>>(object: I): Proposal;
};
export declare const TallyResult: {
    typeUrl: string;
    encode(message: TallyResult, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): TallyResult;
    fromJSON(object: any): TallyResult;
    toJSON(message: TallyResult): unknown;
    fromPartial<I extends {
        yes?: string | undefined;
        abstain?: string | undefined;
        no?: string | undefined;
        noWithVeto?: string | undefined;
    } & {
        yes?: string | undefined;
        abstain?: string | undefined;
        no?: string | undefined;
        noWithVeto?: string | undefined;
    } & Record<Exclude<keyof I, keyof TallyResult>, never>>(object: I): TallyResult;
};
export declare const Vote: {
    typeUrl: string;
    encode(message: Vote, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Vote;
    fromJSON(object: any): Vote;
    toJSON(message: Vote): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
        voter?: string | undefined;
        option?: VoteOption | undefined;
        options?: {
            option?: VoteOption | undefined;
            weight?: string | undefined;
        }[] | undefined;
    } & {
        proposalId?: bigint | undefined;
        voter?: string | undefined;
        option?: VoteOption | undefined;
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
    } & Record<Exclude<keyof I, keyof Vote>, never>>(object: I): Vote;
};
export declare const DepositParams: {
    typeUrl: string;
    encode(message: DepositParams, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DepositParams;
    fromJSON(object: any): DepositParams;
    toJSON(message: DepositParams): unknown;
    fromPartial<I extends {
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
        } & Record<Exclude<keyof I["minDeposit"][number], keyof Coin>, never>)[] & Record<Exclude<keyof I["minDeposit"], keyof {
            denom?: string | undefined;
            amount?: string | undefined;
        }[]>, never>) | undefined;
        maxDepositPeriod?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["maxDepositPeriod"], keyof Duration>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof DepositParams>, never>>(object: I): DepositParams;
};
export declare const VotingParams: {
    typeUrl: string;
    encode(message: VotingParams, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): VotingParams;
    fromJSON(object: any): VotingParams;
    toJSON(message: VotingParams): unknown;
    fromPartial<I extends {
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
        } & Record<Exclude<keyof I["votingPeriod"], keyof Duration>, never>) | undefined;
    } & Record<Exclude<keyof I, "votingPeriod">, never>>(object: I): VotingParams;
};
export declare const TallyParams: {
    typeUrl: string;
    encode(message: TallyParams, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): TallyParams;
    fromJSON(object: any): TallyParams;
    toJSON(message: TallyParams): unknown;
    fromPartial<I extends {
        quorum?: Uint8Array | undefined;
        threshold?: Uint8Array | undefined;
        vetoThreshold?: Uint8Array | undefined;
    } & {
        quorum?: Uint8Array | undefined;
        threshold?: Uint8Array | undefined;
        vetoThreshold?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof TallyParams>, never>>(object: I): TallyParams;
};
