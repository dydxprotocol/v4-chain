import { Timestamp } from "../../../google/protobuf/timestamp";
import { Duration } from "../../../google/protobuf/duration";
import { Any } from "../../../google/protobuf/any";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.group.v1";
/** VoteOption enumerates the valid vote options for a given proposal. */
export declare enum VoteOption {
    /**
     * VOTE_OPTION_UNSPECIFIED - VOTE_OPTION_UNSPECIFIED defines an unspecified vote option which will
     * return an error.
     */
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
/** ProposalStatus defines proposal statuses. */
export declare enum ProposalStatus {
    /** PROPOSAL_STATUS_UNSPECIFIED - An empty value is invalid and not allowed. */
    PROPOSAL_STATUS_UNSPECIFIED = 0,
    /** PROPOSAL_STATUS_SUBMITTED - Initial status of a proposal when submitted. */
    PROPOSAL_STATUS_SUBMITTED = 1,
    /**
     * PROPOSAL_STATUS_ACCEPTED - Final status of a proposal when the final tally is done and the outcome
     * passes the group policy's decision policy.
     */
    PROPOSAL_STATUS_ACCEPTED = 2,
    /**
     * PROPOSAL_STATUS_REJECTED - Final status of a proposal when the final tally is done and the outcome
     * is rejected by the group policy's decision policy.
     */
    PROPOSAL_STATUS_REJECTED = 3,
    /**
     * PROPOSAL_STATUS_ABORTED - Final status of a proposal when the group policy is modified before the
     * final tally.
     */
    PROPOSAL_STATUS_ABORTED = 4,
    /**
     * PROPOSAL_STATUS_WITHDRAWN - A proposal can be withdrawn before the voting start time by the owner.
     * When this happens the final status is Withdrawn.
     */
    PROPOSAL_STATUS_WITHDRAWN = 5,
    UNRECOGNIZED = -1
}
export declare function proposalStatusFromJSON(object: any): ProposalStatus;
export declare function proposalStatusToJSON(object: ProposalStatus): string;
/** ProposalExecutorResult defines types of proposal executor results. */
export declare enum ProposalExecutorResult {
    /** PROPOSAL_EXECUTOR_RESULT_UNSPECIFIED - An empty value is not allowed. */
    PROPOSAL_EXECUTOR_RESULT_UNSPECIFIED = 0,
    /** PROPOSAL_EXECUTOR_RESULT_NOT_RUN - We have not yet run the executor. */
    PROPOSAL_EXECUTOR_RESULT_NOT_RUN = 1,
    /** PROPOSAL_EXECUTOR_RESULT_SUCCESS - The executor was successful and proposed action updated state. */
    PROPOSAL_EXECUTOR_RESULT_SUCCESS = 2,
    /** PROPOSAL_EXECUTOR_RESULT_FAILURE - The executor returned an error and proposed action didn't update state. */
    PROPOSAL_EXECUTOR_RESULT_FAILURE = 3,
    UNRECOGNIZED = -1
}
export declare function proposalExecutorResultFromJSON(object: any): ProposalExecutorResult;
export declare function proposalExecutorResultToJSON(object: ProposalExecutorResult): string;
/**
 * Member represents a group member with an account address,
 * non-zero weight, metadata and added_at timestamp.
 */
export interface Member {
    /** address is the member's account address. */
    address: string;
    /** weight is the member's voting weight that should be greater than 0. */
    weight: string;
    /** metadata is any arbitrary metadata attached to the member. */
    metadata: string;
    /** added_at is a timestamp specifying when a member was added. */
    addedAt: Timestamp;
}
/**
 * MemberRequest represents a group member to be used in Msg server requests.
 * Contrary to `Member`, it doesn't have any `added_at` field
 * since this field cannot be set as part of requests.
 */
export interface MemberRequest {
    /** address is the member's account address. */
    address: string;
    /** weight is the member's voting weight that should be greater than 0. */
    weight: string;
    /** metadata is any arbitrary metadata attached to the member. */
    metadata: string;
}
/**
 * ThresholdDecisionPolicy is a decision policy where a proposal passes when it
 * satisfies the two following conditions:
 * 1. The sum of all `YES` voter's weights is greater or equal than the defined
 *    `threshold`.
 * 2. The voting and execution periods of the proposal respect the parameters
 *    given by `windows`.
 */
export interface ThresholdDecisionPolicy {
    /**
     * threshold is the minimum weighted sum of `YES` votes that must be met or
     * exceeded for a proposal to succeed.
     */
    threshold: string;
    /** windows defines the different windows for voting and execution. */
    windows?: DecisionPolicyWindows;
}
/**
 * PercentageDecisionPolicy is a decision policy where a proposal passes when
 * it satisfies the two following conditions:
 * 1. The percentage of all `YES` voters' weights out of the total group weight
 *    is greater or equal than the given `percentage`.
 * 2. The voting and execution periods of the proposal respect the parameters
 *    given by `windows`.
 */
export interface PercentageDecisionPolicy {
    /**
     * percentage is the minimum percentage of the weighted sum of `YES` votes must
     * meet for a proposal to succeed.
     */
    percentage: string;
    /** windows defines the different windows for voting and execution. */
    windows?: DecisionPolicyWindows;
}
/** DecisionPolicyWindows defines the different windows for voting and execution. */
export interface DecisionPolicyWindows {
    /**
     * voting_period is the duration from submission of a proposal to the end of voting period
     * Within this times votes can be submitted with MsgVote.
     */
    votingPeriod: Duration;
    /**
     * min_execution_period is the minimum duration after the proposal submission
     * where members can start sending MsgExec. This means that the window for
     * sending a MsgExec transaction is:
     * `[ submission + min_execution_period ; submission + voting_period + max_execution_period]`
     * where max_execution_period is a app-specific config, defined in the keeper.
     * If not set, min_execution_period will default to 0.
     *
     * Please make sure to set a `min_execution_period` that is smaller than
     * `voting_period + max_execution_period`, or else the above execution window
     * is empty, meaning that all proposals created with this decision policy
     * won't be able to be executed.
     */
    minExecutionPeriod: Duration;
}
/** GroupInfo represents the high-level on-chain information for a group. */
export interface GroupInfo {
    /** id is the unique ID of the group. */
    id: bigint;
    /** admin is the account address of the group's admin. */
    admin: string;
    /** metadata is any arbitrary metadata to attached to the group. */
    metadata: string;
    /**
     * version is used to track changes to a group's membership structure that
     * would break existing proposals. Whenever any members weight is changed,
     * or any member is added or removed this version is incremented and will
     * cause proposals based on older versions of this group to fail
     */
    version: bigint;
    /** total_weight is the sum of the group members' weights. */
    totalWeight: string;
    /** created_at is a timestamp specifying when a group was created. */
    createdAt: Timestamp;
}
/** GroupMember represents the relationship between a group and a member. */
export interface GroupMember {
    /** group_id is the unique ID of the group. */
    groupId: bigint;
    /** member is the member data. */
    member?: Member;
}
/** GroupPolicyInfo represents the high-level on-chain information for a group policy. */
export interface GroupPolicyInfo {
    /** address is the account address of group policy. */
    address: string;
    /** group_id is the unique ID of the group. */
    groupId: bigint;
    /** admin is the account address of the group admin. */
    admin: string;
    /** metadata is any arbitrary metadata attached to the group policy. */
    metadata: string;
    /**
     * version is used to track changes to a group's GroupPolicyInfo structure that
     * would create a different result on a running proposal.
     */
    version: bigint;
    /** decision_policy specifies the group policy's decision policy. */
    decisionPolicy?: Any;
    /** created_at is a timestamp specifying when a group policy was created. */
    createdAt: Timestamp;
}
/**
 * Proposal defines a group proposal. Any member of a group can submit a proposal
 * for a group policy to decide upon.
 * A proposal consists of a set of `sdk.Msg`s that will be executed if the proposal
 * passes as well as some optional metadata associated with the proposal.
 */
export interface Proposal {
    /** id is the unique id of the proposal. */
    id: bigint;
    /** group_policy_address is the account address of group policy. */
    groupPolicyAddress: string;
    /** metadata is any arbitrary metadata attached to the proposal. */
    metadata: string;
    /** proposers are the account addresses of the proposers. */
    proposers: string[];
    /** submit_time is a timestamp specifying when a proposal was submitted. */
    submitTime: Timestamp;
    /**
     * group_version tracks the version of the group at proposal submission.
     * This field is here for informational purposes only.
     */
    groupVersion: bigint;
    /**
     * group_policy_version tracks the version of the group policy at proposal submission.
     * When a decision policy is changed, existing proposals from previous policy
     * versions will become invalid with the `ABORTED` status.
     * This field is here for informational purposes only.
     */
    groupPolicyVersion: bigint;
    /** status represents the high level position in the life cycle of the proposal. Initial value is Submitted. */
    status: ProposalStatus;
    /**
     * final_tally_result contains the sums of all weighted votes for this
     * proposal for each vote option. It is empty at submission, and only
     * populated after tallying, at voting period end or at proposal execution,
     * whichever happens first.
     */
    finalTallyResult: TallyResult;
    /**
     * voting_period_end is the timestamp before which voting must be done.
     * Unless a successful MsgExec is called before (to execute a proposal whose
     * tally is successful before the voting period ends), tallying will be done
     * at this point, and the `final_tally_result`and `status` fields will be
     * accordingly updated.
     */
    votingPeriodEnd: Timestamp;
    /** executor_result is the final result of the proposal execution. Initial value is NotRun. */
    executorResult: ProposalExecutorResult;
    /** messages is a list of `sdk.Msg`s that will be executed if the proposal passes. */
    messages: Any[];
    /**
     * title is the title of the proposal
     *
     * Since: cosmos-sdk 0.47
     */
    title: string;
    /**
     * summary is a short summary of the proposal
     *
     * Since: cosmos-sdk 0.47
     */
    summary: string;
}
/** TallyResult represents the sum of weighted votes for each vote option. */
export interface TallyResult {
    /** yes_count is the weighted sum of yes votes. */
    yesCount: string;
    /** abstain_count is the weighted sum of abstainers. */
    abstainCount: string;
    /** no_count is the weighted sum of no votes. */
    noCount: string;
    /** no_with_veto_count is the weighted sum of veto. */
    noWithVetoCount: string;
}
/** Vote represents a vote for a proposal. */
export interface Vote {
    /** proposal is the unique ID of the proposal. */
    proposalId: bigint;
    /** voter is the account address of the voter. */
    voter: string;
    /** option is the voter's choice on the proposal. */
    option: VoteOption;
    /** metadata is any arbitrary metadata attached to the vote. */
    metadata: string;
    /** submit_time is the timestamp when the vote was submitted. */
    submitTime: Timestamp;
}
export declare const Member: {
    typeUrl: string;
    encode(message: Member, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Member;
    fromJSON(object: any): Member;
    toJSON(message: Member): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        weight?: string | undefined;
        metadata?: string | undefined;
        addedAt?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
    } & {
        address?: string | undefined;
        weight?: string | undefined;
        metadata?: string | undefined;
        addedAt?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["addedAt"], keyof Timestamp>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Member>, never>>(object: I): Member;
};
export declare const MemberRequest: {
    typeUrl: string;
    encode(message: MemberRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MemberRequest;
    fromJSON(object: any): MemberRequest;
    toJSON(message: MemberRequest): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        weight?: string | undefined;
        metadata?: string | undefined;
    } & {
        address?: string | undefined;
        weight?: string | undefined;
        metadata?: string | undefined;
    } & Record<Exclude<keyof I, keyof MemberRequest>, never>>(object: I): MemberRequest;
};
export declare const ThresholdDecisionPolicy: {
    typeUrl: string;
    encode(message: ThresholdDecisionPolicy, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ThresholdDecisionPolicy;
    fromJSON(object: any): ThresholdDecisionPolicy;
    toJSON(message: ThresholdDecisionPolicy): unknown;
    fromPartial<I extends {
        threshold?: string | undefined;
        windows?: {
            votingPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            minExecutionPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } | undefined;
    } & {
        threshold?: string | undefined;
        windows?: ({
            votingPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            minExecutionPeriod?: {
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
            } & Record<Exclude<keyof I["windows"]["votingPeriod"], keyof Duration>, never>) | undefined;
            minExecutionPeriod?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["windows"]["minExecutionPeriod"], keyof Duration>, never>) | undefined;
        } & Record<Exclude<keyof I["windows"], keyof DecisionPolicyWindows>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ThresholdDecisionPolicy>, never>>(object: I): ThresholdDecisionPolicy;
};
export declare const PercentageDecisionPolicy: {
    typeUrl: string;
    encode(message: PercentageDecisionPolicy, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): PercentageDecisionPolicy;
    fromJSON(object: any): PercentageDecisionPolicy;
    toJSON(message: PercentageDecisionPolicy): unknown;
    fromPartial<I extends {
        percentage?: string | undefined;
        windows?: {
            votingPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            minExecutionPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } | undefined;
    } & {
        percentage?: string | undefined;
        windows?: ({
            votingPeriod?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            minExecutionPeriod?: {
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
            } & Record<Exclude<keyof I["windows"]["votingPeriod"], keyof Duration>, never>) | undefined;
            minExecutionPeriod?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["windows"]["minExecutionPeriod"], keyof Duration>, never>) | undefined;
        } & Record<Exclude<keyof I["windows"], keyof DecisionPolicyWindows>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof PercentageDecisionPolicy>, never>>(object: I): PercentageDecisionPolicy;
};
export declare const DecisionPolicyWindows: {
    typeUrl: string;
    encode(message: DecisionPolicyWindows, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DecisionPolicyWindows;
    fromJSON(object: any): DecisionPolicyWindows;
    toJSON(message: DecisionPolicyWindows): unknown;
    fromPartial<I extends {
        votingPeriod?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        minExecutionPeriod?: {
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
        minExecutionPeriod?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["minExecutionPeriod"], keyof Duration>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof DecisionPolicyWindows>, never>>(object: I): DecisionPolicyWindows;
};
export declare const GroupInfo: {
    typeUrl: string;
    encode(message: GroupInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GroupInfo;
    fromJSON(object: any): GroupInfo;
    toJSON(message: GroupInfo): unknown;
    fromPartial<I extends {
        id?: bigint | undefined;
        admin?: string | undefined;
        metadata?: string | undefined;
        version?: bigint | undefined;
        totalWeight?: string | undefined;
        createdAt?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
    } & {
        id?: bigint | undefined;
        admin?: string | undefined;
        metadata?: string | undefined;
        version?: bigint | undefined;
        totalWeight?: string | undefined;
        createdAt?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["createdAt"], keyof Timestamp>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GroupInfo>, never>>(object: I): GroupInfo;
};
export declare const GroupMember: {
    typeUrl: string;
    encode(message: GroupMember, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GroupMember;
    fromJSON(object: any): GroupMember;
    toJSON(message: GroupMember): unknown;
    fromPartial<I extends {
        groupId?: bigint | undefined;
        member?: {
            address?: string | undefined;
            weight?: string | undefined;
            metadata?: string | undefined;
            addedAt?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } | undefined;
    } & {
        groupId?: bigint | undefined;
        member?: ({
            address?: string | undefined;
            weight?: string | undefined;
            metadata?: string | undefined;
            addedAt?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } & {
            address?: string | undefined;
            weight?: string | undefined;
            metadata?: string | undefined;
            addedAt?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["member"]["addedAt"], keyof Timestamp>, never>) | undefined;
        } & Record<Exclude<keyof I["member"], keyof Member>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GroupMember>, never>>(object: I): GroupMember;
};
export declare const GroupPolicyInfo: {
    typeUrl: string;
    encode(message: GroupPolicyInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GroupPolicyInfo;
    fromJSON(object: any): GroupPolicyInfo;
    toJSON(message: GroupPolicyInfo): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        groupId?: bigint | undefined;
        admin?: string | undefined;
        metadata?: string | undefined;
        version?: bigint | undefined;
        decisionPolicy?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        createdAt?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
    } & {
        address?: string | undefined;
        groupId?: bigint | undefined;
        admin?: string | undefined;
        metadata?: string | undefined;
        version?: bigint | undefined;
        decisionPolicy?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["decisionPolicy"], keyof Any>, never>) | undefined;
        createdAt?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["createdAt"], keyof Timestamp>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GroupPolicyInfo>, never>>(object: I): GroupPolicyInfo;
};
export declare const Proposal: {
    typeUrl: string;
    encode(message: Proposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Proposal;
    fromJSON(object: any): Proposal;
    toJSON(message: Proposal): unknown;
    fromPartial<I extends {
        id?: bigint | undefined;
        groupPolicyAddress?: string | undefined;
        metadata?: string | undefined;
        proposers?: string[] | undefined;
        submitTime?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        groupVersion?: bigint | undefined;
        groupPolicyVersion?: bigint | undefined;
        status?: ProposalStatus | undefined;
        finalTallyResult?: {
            yesCount?: string | undefined;
            abstainCount?: string | undefined;
            noCount?: string | undefined;
            noWithVetoCount?: string | undefined;
        } | undefined;
        votingPeriodEnd?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        executorResult?: ProposalExecutorResult | undefined;
        messages?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[] | undefined;
        title?: string | undefined;
        summary?: string | undefined;
    } & {
        id?: bigint | undefined;
        groupPolicyAddress?: string | undefined;
        metadata?: string | undefined;
        proposers?: (string[] & string[] & Record<Exclude<keyof I["proposers"], keyof string[]>, never>) | undefined;
        submitTime?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["submitTime"], keyof Timestamp>, never>) | undefined;
        groupVersion?: bigint | undefined;
        groupPolicyVersion?: bigint | undefined;
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
        } & Record<Exclude<keyof I["finalTallyResult"], keyof TallyResult>, never>) | undefined;
        votingPeriodEnd?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["votingPeriodEnd"], keyof Timestamp>, never>) | undefined;
        executorResult?: ProposalExecutorResult | undefined;
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
        title?: string | undefined;
        summary?: string | undefined;
    } & Record<Exclude<keyof I, keyof Proposal>, never>>(object: I): Proposal;
};
export declare const TallyResult: {
    typeUrl: string;
    encode(message: TallyResult, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): TallyResult;
    fromJSON(object: any): TallyResult;
    toJSON(message: TallyResult): unknown;
    fromPartial<I extends {
        yesCount?: string | undefined;
        abstainCount?: string | undefined;
        noCount?: string | undefined;
        noWithVetoCount?: string | undefined;
    } & {
        yesCount?: string | undefined;
        abstainCount?: string | undefined;
        noCount?: string | undefined;
        noWithVetoCount?: string | undefined;
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
        metadata?: string | undefined;
        submitTime?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
    } & {
        proposalId?: bigint | undefined;
        voter?: string | undefined;
        option?: VoteOption | undefined;
        metadata?: string | undefined;
        submitTime?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["submitTime"], keyof Timestamp>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Vote>, never>>(object: I): Vote;
};
