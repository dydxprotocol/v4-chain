import { MemberRequest, VoteOption, ProposalExecutorResult } from "./types";
import { Any } from "../../../google/protobuf/any";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.group.v1";
/** Exec defines modes of execution of a proposal on creation or on new vote. */
export declare enum Exec {
    /**
     * EXEC_UNSPECIFIED - An empty value means that there should be a separate
     * MsgExec request for the proposal to execute.
     */
    EXEC_UNSPECIFIED = 0,
    /**
     * EXEC_TRY - Try to execute the proposal immediately.
     * If the proposal is not allowed per the DecisionPolicy,
     * the proposal will still be open and could
     * be executed at a later point.
     */
    EXEC_TRY = 1,
    UNRECOGNIZED = -1
}
export declare function execFromJSON(object: any): Exec;
export declare function execToJSON(object: Exec): string;
/** MsgCreateGroup is the Msg/CreateGroup request type. */
export interface MsgCreateGroup {
    /** admin is the account address of the group admin. */
    admin: string;
    /** members defines the group members. */
    members: MemberRequest[];
    /** metadata is any arbitrary metadata to attached to the group. */
    metadata: string;
}
/** MsgCreateGroupResponse is the Msg/CreateGroup response type. */
export interface MsgCreateGroupResponse {
    /** group_id is the unique ID of the newly created group. */
    groupId: bigint;
}
/** MsgUpdateGroupMembers is the Msg/UpdateGroupMembers request type. */
export interface MsgUpdateGroupMembers {
    /** admin is the account address of the group admin. */
    admin: string;
    /** group_id is the unique ID of the group. */
    groupId: bigint;
    /**
     * member_updates is the list of members to update,
     * set weight to 0 to remove a member.
     */
    memberUpdates: MemberRequest[];
}
/** MsgUpdateGroupMembersResponse is the Msg/UpdateGroupMembers response type. */
export interface MsgUpdateGroupMembersResponse {
}
/** MsgUpdateGroupAdmin is the Msg/UpdateGroupAdmin request type. */
export interface MsgUpdateGroupAdmin {
    /** admin is the current account address of the group admin. */
    admin: string;
    /** group_id is the unique ID of the group. */
    groupId: bigint;
    /** new_admin is the group new admin account address. */
    newAdmin: string;
}
/** MsgUpdateGroupAdminResponse is the Msg/UpdateGroupAdmin response type. */
export interface MsgUpdateGroupAdminResponse {
}
/** MsgUpdateGroupMetadata is the Msg/UpdateGroupMetadata request type. */
export interface MsgUpdateGroupMetadata {
    /** admin is the account address of the group admin. */
    admin: string;
    /** group_id is the unique ID of the group. */
    groupId: bigint;
    /** metadata is the updated group's metadata. */
    metadata: string;
}
/** MsgUpdateGroupMetadataResponse is the Msg/UpdateGroupMetadata response type. */
export interface MsgUpdateGroupMetadataResponse {
}
/** MsgCreateGroupPolicy is the Msg/CreateGroupPolicy request type. */
export interface MsgCreateGroupPolicy {
    /** admin is the account address of the group admin. */
    admin: string;
    /** group_id is the unique ID of the group. */
    groupId: bigint;
    /** metadata is any arbitrary metadata attached to the group policy. */
    metadata: string;
    /** decision_policy specifies the group policy's decision policy. */
    decisionPolicy?: Any;
}
/** MsgCreateGroupPolicyResponse is the Msg/CreateGroupPolicy response type. */
export interface MsgCreateGroupPolicyResponse {
    /** address is the account address of the newly created group policy. */
    address: string;
}
/** MsgUpdateGroupPolicyAdmin is the Msg/UpdateGroupPolicyAdmin request type. */
export interface MsgUpdateGroupPolicyAdmin {
    /** admin is the account address of the group admin. */
    admin: string;
    /** group_policy_address is the account address of the group policy. */
    groupPolicyAddress: string;
    /** new_admin is the new group policy admin. */
    newAdmin: string;
}
/** MsgUpdateGroupPolicyAdminResponse is the Msg/UpdateGroupPolicyAdmin response type. */
export interface MsgUpdateGroupPolicyAdminResponse {
}
/** MsgCreateGroupWithPolicy is the Msg/CreateGroupWithPolicy request type. */
export interface MsgCreateGroupWithPolicy {
    /** admin is the account address of the group and group policy admin. */
    admin: string;
    /** members defines the group members. */
    members: MemberRequest[];
    /** group_metadata is any arbitrary metadata attached to the group. */
    groupMetadata: string;
    /** group_policy_metadata is any arbitrary metadata attached to the group policy. */
    groupPolicyMetadata: string;
    /**
     * group_policy_as_admin is a boolean field, if set to true, the group policy account address will be used as group
     * and group policy admin.
     */
    groupPolicyAsAdmin: boolean;
    /** decision_policy specifies the group policy's decision policy. */
    decisionPolicy?: Any;
}
/** MsgCreateGroupWithPolicyResponse is the Msg/CreateGroupWithPolicy response type. */
export interface MsgCreateGroupWithPolicyResponse {
    /** group_id is the unique ID of the newly created group with policy. */
    groupId: bigint;
    /** group_policy_address is the account address of the newly created group policy. */
    groupPolicyAddress: string;
}
/** MsgUpdateGroupPolicyDecisionPolicy is the Msg/UpdateGroupPolicyDecisionPolicy request type. */
export interface MsgUpdateGroupPolicyDecisionPolicy {
    /** admin is the account address of the group admin. */
    admin: string;
    /** group_policy_address is the account address of group policy. */
    groupPolicyAddress: string;
    /** decision_policy is the updated group policy's decision policy. */
    decisionPolicy?: Any;
}
/** MsgUpdateGroupPolicyDecisionPolicyResponse is the Msg/UpdateGroupPolicyDecisionPolicy response type. */
export interface MsgUpdateGroupPolicyDecisionPolicyResponse {
}
/** MsgUpdateGroupPolicyMetadata is the Msg/UpdateGroupPolicyMetadata request type. */
export interface MsgUpdateGroupPolicyMetadata {
    /** admin is the account address of the group admin. */
    admin: string;
    /** group_policy_address is the account address of group policy. */
    groupPolicyAddress: string;
    /** metadata is the group policy metadata to be updated. */
    metadata: string;
}
/** MsgUpdateGroupPolicyMetadataResponse is the Msg/UpdateGroupPolicyMetadata response type. */
export interface MsgUpdateGroupPolicyMetadataResponse {
}
/** MsgSubmitProposal is the Msg/SubmitProposal request type. */
export interface MsgSubmitProposal {
    /** group_policy_address is the account address of group policy. */
    groupPolicyAddress: string;
    /**
     * proposers are the account addresses of the proposers.
     * Proposers signatures will be counted as yes votes.
     */
    proposers: string[];
    /** metadata is any arbitrary metadata attached to the proposal. */
    metadata: string;
    /** messages is a list of `sdk.Msg`s that will be executed if the proposal passes. */
    messages: Any[];
    /**
     * exec defines the mode of execution of the proposal,
     * whether it should be executed immediately on creation or not.
     * If so, proposers signatures are considered as Yes votes.
     */
    exec: Exec;
    /**
     * title is the title of the proposal.
     *
     * Since: cosmos-sdk 0.47
     */
    title: string;
    /**
     * summary is the summary of the proposal.
     *
     * Since: cosmos-sdk 0.47
     */
    summary: string;
}
/** MsgSubmitProposalResponse is the Msg/SubmitProposal response type. */
export interface MsgSubmitProposalResponse {
    /** proposal is the unique ID of the proposal. */
    proposalId: bigint;
}
/** MsgWithdrawProposal is the Msg/WithdrawProposal request type. */
export interface MsgWithdrawProposal {
    /** proposal is the unique ID of the proposal. */
    proposalId: bigint;
    /** address is the admin of the group policy or one of the proposer of the proposal. */
    address: string;
}
/** MsgWithdrawProposalResponse is the Msg/WithdrawProposal response type. */
export interface MsgWithdrawProposalResponse {
}
/** MsgVote is the Msg/Vote request type. */
export interface MsgVote {
    /** proposal is the unique ID of the proposal. */
    proposalId: bigint;
    /** voter is the voter account address. */
    voter: string;
    /** option is the voter's choice on the proposal. */
    option: VoteOption;
    /** metadata is any arbitrary metadata attached to the vote. */
    metadata: string;
    /**
     * exec defines whether the proposal should be executed
     * immediately after voting or not.
     */
    exec: Exec;
}
/** MsgVoteResponse is the Msg/Vote response type. */
export interface MsgVoteResponse {
}
/** MsgExec is the Msg/Exec request type. */
export interface MsgExec {
    /** proposal is the unique ID of the proposal. */
    proposalId: bigint;
    /** executor is the account address used to execute the proposal. */
    executor: string;
}
/** MsgExecResponse is the Msg/Exec request type. */
export interface MsgExecResponse {
    /** result is the final result of the proposal execution. */
    result: ProposalExecutorResult;
}
/** MsgLeaveGroup is the Msg/LeaveGroup request type. */
export interface MsgLeaveGroup {
    /** address is the account address of the group member. */
    address: string;
    /** group_id is the unique ID of the group. */
    groupId: bigint;
}
/** MsgLeaveGroupResponse is the Msg/LeaveGroup response type. */
export interface MsgLeaveGroupResponse {
}
export declare const MsgCreateGroup: {
    typeUrl: string;
    encode(message: MsgCreateGroup, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgCreateGroup;
    fromJSON(object: any): MsgCreateGroup;
    toJSON(message: MsgCreateGroup): unknown;
    fromPartial<I extends {
        admin?: string | undefined;
        members?: {
            address?: string | undefined;
            weight?: string | undefined;
            metadata?: string | undefined;
        }[] | undefined;
        metadata?: string | undefined;
    } & {
        admin?: string | undefined;
        members?: ({
            address?: string | undefined;
            weight?: string | undefined;
            metadata?: string | undefined;
        }[] & ({
            address?: string | undefined;
            weight?: string | undefined;
            metadata?: string | undefined;
        } & {
            address?: string | undefined;
            weight?: string | undefined;
            metadata?: string | undefined;
        } & Record<Exclude<keyof I["members"][number], keyof MemberRequest>, never>)[] & Record<Exclude<keyof I["members"], keyof {
            address?: string | undefined;
            weight?: string | undefined;
            metadata?: string | undefined;
        }[]>, never>) | undefined;
        metadata?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgCreateGroup>, never>>(object: I): MsgCreateGroup;
};
export declare const MsgCreateGroupResponse: {
    typeUrl: string;
    encode(message: MsgCreateGroupResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgCreateGroupResponse;
    fromJSON(object: any): MsgCreateGroupResponse;
    toJSON(message: MsgCreateGroupResponse): unknown;
    fromPartial<I extends {
        groupId?: bigint | undefined;
    } & {
        groupId?: bigint | undefined;
    } & Record<Exclude<keyof I, "groupId">, never>>(object: I): MsgCreateGroupResponse;
};
export declare const MsgUpdateGroupMembers: {
    typeUrl: string;
    encode(message: MsgUpdateGroupMembers, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateGroupMembers;
    fromJSON(object: any): MsgUpdateGroupMembers;
    toJSON(message: MsgUpdateGroupMembers): unknown;
    fromPartial<I extends {
        admin?: string | undefined;
        groupId?: bigint | undefined;
        memberUpdates?: {
            address?: string | undefined;
            weight?: string | undefined;
            metadata?: string | undefined;
        }[] | undefined;
    } & {
        admin?: string | undefined;
        groupId?: bigint | undefined;
        memberUpdates?: ({
            address?: string | undefined;
            weight?: string | undefined;
            metadata?: string | undefined;
        }[] & ({
            address?: string | undefined;
            weight?: string | undefined;
            metadata?: string | undefined;
        } & {
            address?: string | undefined;
            weight?: string | undefined;
            metadata?: string | undefined;
        } & Record<Exclude<keyof I["memberUpdates"][number], keyof MemberRequest>, never>)[] & Record<Exclude<keyof I["memberUpdates"], keyof {
            address?: string | undefined;
            weight?: string | undefined;
            metadata?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof MsgUpdateGroupMembers>, never>>(object: I): MsgUpdateGroupMembers;
};
export declare const MsgUpdateGroupMembersResponse: {
    typeUrl: string;
    encode(_: MsgUpdateGroupMembersResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateGroupMembersResponse;
    fromJSON(_: any): MsgUpdateGroupMembersResponse;
    toJSON(_: MsgUpdateGroupMembersResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgUpdateGroupMembersResponse;
};
export declare const MsgUpdateGroupAdmin: {
    typeUrl: string;
    encode(message: MsgUpdateGroupAdmin, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateGroupAdmin;
    fromJSON(object: any): MsgUpdateGroupAdmin;
    toJSON(message: MsgUpdateGroupAdmin): unknown;
    fromPartial<I extends {
        admin?: string | undefined;
        groupId?: bigint | undefined;
        newAdmin?: string | undefined;
    } & {
        admin?: string | undefined;
        groupId?: bigint | undefined;
        newAdmin?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgUpdateGroupAdmin>, never>>(object: I): MsgUpdateGroupAdmin;
};
export declare const MsgUpdateGroupAdminResponse: {
    typeUrl: string;
    encode(_: MsgUpdateGroupAdminResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateGroupAdminResponse;
    fromJSON(_: any): MsgUpdateGroupAdminResponse;
    toJSON(_: MsgUpdateGroupAdminResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgUpdateGroupAdminResponse;
};
export declare const MsgUpdateGroupMetadata: {
    typeUrl: string;
    encode(message: MsgUpdateGroupMetadata, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateGroupMetadata;
    fromJSON(object: any): MsgUpdateGroupMetadata;
    toJSON(message: MsgUpdateGroupMetadata): unknown;
    fromPartial<I extends {
        admin?: string | undefined;
        groupId?: bigint | undefined;
        metadata?: string | undefined;
    } & {
        admin?: string | undefined;
        groupId?: bigint | undefined;
        metadata?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgUpdateGroupMetadata>, never>>(object: I): MsgUpdateGroupMetadata;
};
export declare const MsgUpdateGroupMetadataResponse: {
    typeUrl: string;
    encode(_: MsgUpdateGroupMetadataResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateGroupMetadataResponse;
    fromJSON(_: any): MsgUpdateGroupMetadataResponse;
    toJSON(_: MsgUpdateGroupMetadataResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgUpdateGroupMetadataResponse;
};
export declare const MsgCreateGroupPolicy: {
    typeUrl: string;
    encode(message: MsgCreateGroupPolicy, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgCreateGroupPolicy;
    fromJSON(object: any): MsgCreateGroupPolicy;
    toJSON(message: MsgCreateGroupPolicy): unknown;
    fromPartial<I extends {
        admin?: string | undefined;
        groupId?: bigint | undefined;
        metadata?: string | undefined;
        decisionPolicy?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
    } & {
        admin?: string | undefined;
        groupId?: bigint | undefined;
        metadata?: string | undefined;
        decisionPolicy?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["decisionPolicy"], keyof Any>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof MsgCreateGroupPolicy>, never>>(object: I): MsgCreateGroupPolicy;
};
export declare const MsgCreateGroupPolicyResponse: {
    typeUrl: string;
    encode(message: MsgCreateGroupPolicyResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgCreateGroupPolicyResponse;
    fromJSON(object: any): MsgCreateGroupPolicyResponse;
    toJSON(message: MsgCreateGroupPolicyResponse): unknown;
    fromPartial<I extends {
        address?: string | undefined;
    } & {
        address?: string | undefined;
    } & Record<Exclude<keyof I, "address">, never>>(object: I): MsgCreateGroupPolicyResponse;
};
export declare const MsgUpdateGroupPolicyAdmin: {
    typeUrl: string;
    encode(message: MsgUpdateGroupPolicyAdmin, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateGroupPolicyAdmin;
    fromJSON(object: any): MsgUpdateGroupPolicyAdmin;
    toJSON(message: MsgUpdateGroupPolicyAdmin): unknown;
    fromPartial<I extends {
        admin?: string | undefined;
        groupPolicyAddress?: string | undefined;
        newAdmin?: string | undefined;
    } & {
        admin?: string | undefined;
        groupPolicyAddress?: string | undefined;
        newAdmin?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgUpdateGroupPolicyAdmin>, never>>(object: I): MsgUpdateGroupPolicyAdmin;
};
export declare const MsgUpdateGroupPolicyAdminResponse: {
    typeUrl: string;
    encode(_: MsgUpdateGroupPolicyAdminResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateGroupPolicyAdminResponse;
    fromJSON(_: any): MsgUpdateGroupPolicyAdminResponse;
    toJSON(_: MsgUpdateGroupPolicyAdminResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgUpdateGroupPolicyAdminResponse;
};
export declare const MsgCreateGroupWithPolicy: {
    typeUrl: string;
    encode(message: MsgCreateGroupWithPolicy, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgCreateGroupWithPolicy;
    fromJSON(object: any): MsgCreateGroupWithPolicy;
    toJSON(message: MsgCreateGroupWithPolicy): unknown;
    fromPartial<I extends {
        admin?: string | undefined;
        members?: {
            address?: string | undefined;
            weight?: string | undefined;
            metadata?: string | undefined;
        }[] | undefined;
        groupMetadata?: string | undefined;
        groupPolicyMetadata?: string | undefined;
        groupPolicyAsAdmin?: boolean | undefined;
        decisionPolicy?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
    } & {
        admin?: string | undefined;
        members?: ({
            address?: string | undefined;
            weight?: string | undefined;
            metadata?: string | undefined;
        }[] & ({
            address?: string | undefined;
            weight?: string | undefined;
            metadata?: string | undefined;
        } & {
            address?: string | undefined;
            weight?: string | undefined;
            metadata?: string | undefined;
        } & Record<Exclude<keyof I["members"][number], keyof MemberRequest>, never>)[] & Record<Exclude<keyof I["members"], keyof {
            address?: string | undefined;
            weight?: string | undefined;
            metadata?: string | undefined;
        }[]>, never>) | undefined;
        groupMetadata?: string | undefined;
        groupPolicyMetadata?: string | undefined;
        groupPolicyAsAdmin?: boolean | undefined;
        decisionPolicy?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["decisionPolicy"], keyof Any>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof MsgCreateGroupWithPolicy>, never>>(object: I): MsgCreateGroupWithPolicy;
};
export declare const MsgCreateGroupWithPolicyResponse: {
    typeUrl: string;
    encode(message: MsgCreateGroupWithPolicyResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgCreateGroupWithPolicyResponse;
    fromJSON(object: any): MsgCreateGroupWithPolicyResponse;
    toJSON(message: MsgCreateGroupWithPolicyResponse): unknown;
    fromPartial<I extends {
        groupId?: bigint | undefined;
        groupPolicyAddress?: string | undefined;
    } & {
        groupId?: bigint | undefined;
        groupPolicyAddress?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgCreateGroupWithPolicyResponse>, never>>(object: I): MsgCreateGroupWithPolicyResponse;
};
export declare const MsgUpdateGroupPolicyDecisionPolicy: {
    typeUrl: string;
    encode(message: MsgUpdateGroupPolicyDecisionPolicy, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateGroupPolicyDecisionPolicy;
    fromJSON(object: any): MsgUpdateGroupPolicyDecisionPolicy;
    toJSON(message: MsgUpdateGroupPolicyDecisionPolicy): unknown;
    fromPartial<I extends {
        admin?: string | undefined;
        groupPolicyAddress?: string | undefined;
        decisionPolicy?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
    } & {
        admin?: string | undefined;
        groupPolicyAddress?: string | undefined;
        decisionPolicy?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["decisionPolicy"], keyof Any>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof MsgUpdateGroupPolicyDecisionPolicy>, never>>(object: I): MsgUpdateGroupPolicyDecisionPolicy;
};
export declare const MsgUpdateGroupPolicyDecisionPolicyResponse: {
    typeUrl: string;
    encode(_: MsgUpdateGroupPolicyDecisionPolicyResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateGroupPolicyDecisionPolicyResponse;
    fromJSON(_: any): MsgUpdateGroupPolicyDecisionPolicyResponse;
    toJSON(_: MsgUpdateGroupPolicyDecisionPolicyResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgUpdateGroupPolicyDecisionPolicyResponse;
};
export declare const MsgUpdateGroupPolicyMetadata: {
    typeUrl: string;
    encode(message: MsgUpdateGroupPolicyMetadata, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateGroupPolicyMetadata;
    fromJSON(object: any): MsgUpdateGroupPolicyMetadata;
    toJSON(message: MsgUpdateGroupPolicyMetadata): unknown;
    fromPartial<I extends {
        admin?: string | undefined;
        groupPolicyAddress?: string | undefined;
        metadata?: string | undefined;
    } & {
        admin?: string | undefined;
        groupPolicyAddress?: string | undefined;
        metadata?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgUpdateGroupPolicyMetadata>, never>>(object: I): MsgUpdateGroupPolicyMetadata;
};
export declare const MsgUpdateGroupPolicyMetadataResponse: {
    typeUrl: string;
    encode(_: MsgUpdateGroupPolicyMetadataResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateGroupPolicyMetadataResponse;
    fromJSON(_: any): MsgUpdateGroupPolicyMetadataResponse;
    toJSON(_: MsgUpdateGroupPolicyMetadataResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgUpdateGroupPolicyMetadataResponse;
};
export declare const MsgSubmitProposal: {
    typeUrl: string;
    encode(message: MsgSubmitProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgSubmitProposal;
    fromJSON(object: any): MsgSubmitProposal;
    toJSON(message: MsgSubmitProposal): unknown;
    fromPartial<I extends {
        groupPolicyAddress?: string | undefined;
        proposers?: string[] | undefined;
        metadata?: string | undefined;
        messages?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[] | undefined;
        exec?: Exec | undefined;
        title?: string | undefined;
        summary?: string | undefined;
    } & {
        groupPolicyAddress?: string | undefined;
        proposers?: (string[] & string[] & Record<Exclude<keyof I["proposers"], keyof string[]>, never>) | undefined;
        metadata?: string | undefined;
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
        exec?: Exec | undefined;
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
export declare const MsgWithdrawProposal: {
    typeUrl: string;
    encode(message: MsgWithdrawProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgWithdrawProposal;
    fromJSON(object: any): MsgWithdrawProposal;
    toJSON(message: MsgWithdrawProposal): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
        address?: string | undefined;
    } & {
        proposalId?: bigint | undefined;
        address?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgWithdrawProposal>, never>>(object: I): MsgWithdrawProposal;
};
export declare const MsgWithdrawProposalResponse: {
    typeUrl: string;
    encode(_: MsgWithdrawProposalResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgWithdrawProposalResponse;
    fromJSON(_: any): MsgWithdrawProposalResponse;
    toJSON(_: MsgWithdrawProposalResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgWithdrawProposalResponse;
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
        exec?: Exec | undefined;
    } & {
        proposalId?: bigint | undefined;
        voter?: string | undefined;
        option?: VoteOption | undefined;
        metadata?: string | undefined;
        exec?: Exec | undefined;
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
export declare const MsgExec: {
    typeUrl: string;
    encode(message: MsgExec, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgExec;
    fromJSON(object: any): MsgExec;
    toJSON(message: MsgExec): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
        executor?: string | undefined;
    } & {
        proposalId?: bigint | undefined;
        executor?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgExec>, never>>(object: I): MsgExec;
};
export declare const MsgExecResponse: {
    typeUrl: string;
    encode(message: MsgExecResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgExecResponse;
    fromJSON(object: any): MsgExecResponse;
    toJSON(message: MsgExecResponse): unknown;
    fromPartial<I extends {
        result?: ProposalExecutorResult | undefined;
    } & {
        result?: ProposalExecutorResult | undefined;
    } & Record<Exclude<keyof I, "result">, never>>(object: I): MsgExecResponse;
};
export declare const MsgLeaveGroup: {
    typeUrl: string;
    encode(message: MsgLeaveGroup, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgLeaveGroup;
    fromJSON(object: any): MsgLeaveGroup;
    toJSON(message: MsgLeaveGroup): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        groupId?: bigint | undefined;
    } & {
        address?: string | undefined;
        groupId?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof MsgLeaveGroup>, never>>(object: I): MsgLeaveGroup;
};
export declare const MsgLeaveGroupResponse: {
    typeUrl: string;
    encode(_: MsgLeaveGroupResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgLeaveGroupResponse;
    fromJSON(_: any): MsgLeaveGroupResponse;
    toJSON(_: MsgLeaveGroupResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgLeaveGroupResponse;
};
/** Msg is the cosmos.group.v1 Msg service. */
export interface Msg {
    /** CreateGroup creates a new group with an admin account address, a list of members and some optional metadata. */
    CreateGroup(request: MsgCreateGroup): Promise<MsgCreateGroupResponse>;
    /** UpdateGroupMembers updates the group members with given group id and admin address. */
    UpdateGroupMembers(request: MsgUpdateGroupMembers): Promise<MsgUpdateGroupMembersResponse>;
    /** UpdateGroupAdmin updates the group admin with given group id and previous admin address. */
    UpdateGroupAdmin(request: MsgUpdateGroupAdmin): Promise<MsgUpdateGroupAdminResponse>;
    /** UpdateGroupMetadata updates the group metadata with given group id and admin address. */
    UpdateGroupMetadata(request: MsgUpdateGroupMetadata): Promise<MsgUpdateGroupMetadataResponse>;
    /** CreateGroupPolicy creates a new group policy using given DecisionPolicy. */
    CreateGroupPolicy(request: MsgCreateGroupPolicy): Promise<MsgCreateGroupPolicyResponse>;
    /** CreateGroupWithPolicy creates a new group with policy. */
    CreateGroupWithPolicy(request: MsgCreateGroupWithPolicy): Promise<MsgCreateGroupWithPolicyResponse>;
    /** UpdateGroupPolicyAdmin updates a group policy admin. */
    UpdateGroupPolicyAdmin(request: MsgUpdateGroupPolicyAdmin): Promise<MsgUpdateGroupPolicyAdminResponse>;
    /** UpdateGroupPolicyDecisionPolicy allows a group policy's decision policy to be updated. */
    UpdateGroupPolicyDecisionPolicy(request: MsgUpdateGroupPolicyDecisionPolicy): Promise<MsgUpdateGroupPolicyDecisionPolicyResponse>;
    /** UpdateGroupPolicyMetadata updates a group policy metadata. */
    UpdateGroupPolicyMetadata(request: MsgUpdateGroupPolicyMetadata): Promise<MsgUpdateGroupPolicyMetadataResponse>;
    /** SubmitProposal submits a new proposal. */
    SubmitProposal(request: MsgSubmitProposal): Promise<MsgSubmitProposalResponse>;
    /** WithdrawProposal withdraws a proposal. */
    WithdrawProposal(request: MsgWithdrawProposal): Promise<MsgWithdrawProposalResponse>;
    /** Vote allows a voter to vote on a proposal. */
    Vote(request: MsgVote): Promise<MsgVoteResponse>;
    /** Exec executes a proposal. */
    Exec(request: MsgExec): Promise<MsgExecResponse>;
    /** LeaveGroup allows a group member to leave the group. */
    LeaveGroup(request: MsgLeaveGroup): Promise<MsgLeaveGroupResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    CreateGroup(request: MsgCreateGroup): Promise<MsgCreateGroupResponse>;
    UpdateGroupMembers(request: MsgUpdateGroupMembers): Promise<MsgUpdateGroupMembersResponse>;
    UpdateGroupAdmin(request: MsgUpdateGroupAdmin): Promise<MsgUpdateGroupAdminResponse>;
    UpdateGroupMetadata(request: MsgUpdateGroupMetadata): Promise<MsgUpdateGroupMetadataResponse>;
    CreateGroupPolicy(request: MsgCreateGroupPolicy): Promise<MsgCreateGroupPolicyResponse>;
    CreateGroupWithPolicy(request: MsgCreateGroupWithPolicy): Promise<MsgCreateGroupWithPolicyResponse>;
    UpdateGroupPolicyAdmin(request: MsgUpdateGroupPolicyAdmin): Promise<MsgUpdateGroupPolicyAdminResponse>;
    UpdateGroupPolicyDecisionPolicy(request: MsgUpdateGroupPolicyDecisionPolicy): Promise<MsgUpdateGroupPolicyDecisionPolicyResponse>;
    UpdateGroupPolicyMetadata(request: MsgUpdateGroupPolicyMetadata): Promise<MsgUpdateGroupPolicyMetadataResponse>;
    SubmitProposal(request: MsgSubmitProposal): Promise<MsgSubmitProposalResponse>;
    WithdrawProposal(request: MsgWithdrawProposal): Promise<MsgWithdrawProposalResponse>;
    Vote(request: MsgVote): Promise<MsgVoteResponse>;
    Exec(request: MsgExec): Promise<MsgExecResponse>;
    LeaveGroup(request: MsgLeaveGroup): Promise<MsgLeaveGroupResponse>;
}
