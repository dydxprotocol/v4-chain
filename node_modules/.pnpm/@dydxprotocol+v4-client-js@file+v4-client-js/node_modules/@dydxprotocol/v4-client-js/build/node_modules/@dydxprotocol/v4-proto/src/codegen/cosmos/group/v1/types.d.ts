/// <reference types="long" />
import { Duration, DurationSDKType } from "../../../google/protobuf/duration";
import { Any, AnySDKType } from "../../../google/protobuf/any";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../../helpers";
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
export declare const VoteOptionSDKType: typeof VoteOption;
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
export declare const ProposalStatusSDKType: typeof ProposalStatus;
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
export declare const ProposalExecutorResultSDKType: typeof ProposalExecutorResult;
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
    addedAt?: Date;
}
/**
 * Member represents a group member with an account address,
 * non-zero weight, metadata and added_at timestamp.
 */
export interface MemberSDKType {
    address: string;
    weight: string;
    metadata: string;
    added_at?: Date;
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
 * MemberRequest represents a group member to be used in Msg server requests.
 * Contrary to `Member`, it doesn't have any `added_at` field
 * since this field cannot be set as part of requests.
 */
export interface MemberRequestSDKType {
    address: string;
    weight: string;
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
 * ThresholdDecisionPolicy is a decision policy where a proposal passes when it
 * satisfies the two following conditions:
 * 1. The sum of all `YES` voter's weights is greater or equal than the defined
 *    `threshold`.
 * 2. The voting and execution periods of the proposal respect the parameters
 *    given by `windows`.
 */
export interface ThresholdDecisionPolicySDKType {
    threshold: string;
    windows?: DecisionPolicyWindowsSDKType;
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
/**
 * PercentageDecisionPolicy is a decision policy where a proposal passes when
 * it satisfies the two following conditions:
 * 1. The percentage of all `YES` voters' weights out of the total group weight
 *    is greater or equal than the given `percentage`.
 * 2. The voting and execution periods of the proposal respect the parameters
 *    given by `windows`.
 */
export interface PercentageDecisionPolicySDKType {
    percentage: string;
    windows?: DecisionPolicyWindowsSDKType;
}
/** DecisionPolicyWindows defines the different windows for voting and execution. */
export interface DecisionPolicyWindows {
    /**
     * voting_period is the duration from submission of a proposal to the end of voting period
     * Within this times votes can be submitted with MsgVote.
     */
    votingPeriod?: Duration;
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
    minExecutionPeriod?: Duration;
}
/** DecisionPolicyWindows defines the different windows for voting and execution. */
export interface DecisionPolicyWindowsSDKType {
    voting_period?: DurationSDKType;
    min_execution_period?: DurationSDKType;
}
/** GroupInfo represents the high-level on-chain information for a group. */
export interface GroupInfo {
    /** id is the unique ID of the group. */
    id: Long;
    /** admin is the account address of the group's admin. */
    admin: string;
    /**
     * metadata is any arbitrary metadata to attached to the group.
     * the recommended format of the metadata is to be found here: https://docs.cosmos.network/v0.47/modules/group#group-1
     */
    metadata: string;
    /**
     * version is used to track changes to a group's membership structure that
     * would break existing proposals. Whenever any members weight is changed,
     * or any member is added or removed this version is incremented and will
     * cause proposals based on older versions of this group to fail
     */
    version: Long;
    /** total_weight is the sum of the group members' weights. */
    totalWeight: string;
    /** created_at is a timestamp specifying when a group was created. */
    createdAt?: Date;
}
/** GroupInfo represents the high-level on-chain information for a group. */
export interface GroupInfoSDKType {
    id: Long;
    admin: string;
    metadata: string;
    version: Long;
    total_weight: string;
    created_at?: Date;
}
/** GroupMember represents the relationship between a group and a member. */
export interface GroupMember {
    /** group_id is the unique ID of the group. */
    groupId: Long;
    /** member is the member data. */
    member?: Member;
}
/** GroupMember represents the relationship between a group and a member. */
export interface GroupMemberSDKType {
    group_id: Long;
    member?: MemberSDKType;
}
/** GroupPolicyInfo represents the high-level on-chain information for a group policy. */
export interface GroupPolicyInfo {
    /** address is the account address of group policy. */
    address: string;
    /** group_id is the unique ID of the group. */
    groupId: Long;
    /** admin is the account address of the group admin. */
    admin: string;
    /**
     * metadata is any arbitrary metadata attached to the group policy.
     * the recommended format of the metadata is to be found here:
     * https://docs.cosmos.network/v0.47/modules/group#decision-policy-1
     */
    metadata: string;
    /**
     * version is used to track changes to a group's GroupPolicyInfo structure that
     * would create a different result on a running proposal.
     */
    version: Long;
    /** decision_policy specifies the group policy's decision policy. */
    decisionPolicy?: Any;
    /** created_at is a timestamp specifying when a group policy was created. */
    createdAt?: Date;
}
/** GroupPolicyInfo represents the high-level on-chain information for a group policy. */
export interface GroupPolicyInfoSDKType {
    address: string;
    group_id: Long;
    admin: string;
    metadata: string;
    version: Long;
    decision_policy?: AnySDKType;
    created_at?: Date;
}
/**
 * Proposal defines a group proposal. Any member of a group can submit a proposal
 * for a group policy to decide upon.
 * A proposal consists of a set of `sdk.Msg`s that will be executed if the proposal
 * passes as well as some optional metadata associated with the proposal.
 */
export interface Proposal {
    /** id is the unique id of the proposal. */
    id: Long;
    /** group_policy_address is the account address of group policy. */
    groupPolicyAddress: string;
    /**
     * metadata is any arbitrary metadata attached to the proposal.
     * the recommended format of the metadata is to be found here:
     * https://docs.cosmos.network/v0.47/modules/group#proposal-4
     */
    metadata: string;
    /** proposers are the account addresses of the proposers. */
    proposers: string[];
    /** submit_time is a timestamp specifying when a proposal was submitted. */
    submitTime?: Date;
    /**
     * group_version tracks the version of the group at proposal submission.
     * This field is here for informational purposes only.
     */
    groupVersion: Long;
    /**
     * group_policy_version tracks the version of the group policy at proposal submission.
     * When a decision policy is changed, existing proposals from previous policy
     * versions will become invalid with the `ABORTED` status.
     * This field is here for informational purposes only.
     */
    groupPolicyVersion: Long;
    /** status represents the high level position in the life cycle of the proposal. Initial value is Submitted. */
    status: ProposalStatus;
    /**
     * final_tally_result contains the sums of all weighted votes for this
     * proposal for each vote option. It is empty at submission, and only
     * populated after tallying, at voting period end or at proposal execution,
     * whichever happens first.
     */
    finalTallyResult?: TallyResult;
    /**
     * voting_period_end is the timestamp before which voting must be done.
     * Unless a successful MsgExec is called before (to execute a proposal whose
     * tally is successful before the voting period ends), tallying will be done
     * at this point, and the `final_tally_result`and `status` fields will be
     * accordingly updated.
     */
    votingPeriodEnd?: Date;
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
/**
 * Proposal defines a group proposal. Any member of a group can submit a proposal
 * for a group policy to decide upon.
 * A proposal consists of a set of `sdk.Msg`s that will be executed if the proposal
 * passes as well as some optional metadata associated with the proposal.
 */
export interface ProposalSDKType {
    id: Long;
    group_policy_address: string;
    metadata: string;
    proposers: string[];
    submit_time?: Date;
    group_version: Long;
    group_policy_version: Long;
    status: ProposalStatus;
    final_tally_result?: TallyResultSDKType;
    voting_period_end?: Date;
    executor_result: ProposalExecutorResult;
    messages: AnySDKType[];
    title: string;
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
/** TallyResult represents the sum of weighted votes for each vote option. */
export interface TallyResultSDKType {
    yes_count: string;
    abstain_count: string;
    no_count: string;
    no_with_veto_count: string;
}
/** Vote represents a vote for a proposal.string metadata */
export interface Vote {
    /** proposal is the unique ID of the proposal. */
    proposalId: Long;
    /** voter is the account address of the voter. */
    voter: string;
    /** option is the voter's choice on the proposal. */
    option: VoteOption;
    /**
     * metadata is any arbitrary metadata attached to the vote.
     * the recommended format of the metadata is to be found here: https://docs.cosmos.network/v0.47/modules/group#vote-2
     */
    metadata: string;
    /** submit_time is the timestamp when the vote was submitted. */
    submitTime?: Date;
}
/** Vote represents a vote for a proposal.string metadata */
export interface VoteSDKType {
    proposal_id: Long;
    voter: string;
    option: VoteOption;
    metadata: string;
    submit_time?: Date;
}
export declare const Member: {
    encode(message: Member, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Member;
    fromPartial(object: DeepPartial<Member>): Member;
};
export declare const MemberRequest: {
    encode(message: MemberRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MemberRequest;
    fromPartial(object: DeepPartial<MemberRequest>): MemberRequest;
};
export declare const ThresholdDecisionPolicy: {
    encode(message: ThresholdDecisionPolicy, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ThresholdDecisionPolicy;
    fromPartial(object: DeepPartial<ThresholdDecisionPolicy>): ThresholdDecisionPolicy;
};
export declare const PercentageDecisionPolicy: {
    encode(message: PercentageDecisionPolicy, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): PercentageDecisionPolicy;
    fromPartial(object: DeepPartial<PercentageDecisionPolicy>): PercentageDecisionPolicy;
};
export declare const DecisionPolicyWindows: {
    encode(message: DecisionPolicyWindows, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): DecisionPolicyWindows;
    fromPartial(object: DeepPartial<DecisionPolicyWindows>): DecisionPolicyWindows;
};
export declare const GroupInfo: {
    encode(message: GroupInfo, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GroupInfo;
    fromPartial(object: DeepPartial<GroupInfo>): GroupInfo;
};
export declare const GroupMember: {
    encode(message: GroupMember, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GroupMember;
    fromPartial(object: DeepPartial<GroupMember>): GroupMember;
};
export declare const GroupPolicyInfo: {
    encode(message: GroupPolicyInfo, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GroupPolicyInfo;
    fromPartial(object: DeepPartial<GroupPolicyInfo>): GroupPolicyInfo;
};
export declare const Proposal: {
    encode(message: Proposal, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Proposal;
    fromPartial(object: DeepPartial<Proposal>): Proposal;
};
export declare const TallyResult: {
    encode(message: TallyResult, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): TallyResult;
    fromPartial(object: DeepPartial<TallyResult>): TallyResult;
};
export declare const Vote: {
    encode(message: Vote, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Vote;
    fromPartial(object: DeepPartial<Vote>): Vote;
};
