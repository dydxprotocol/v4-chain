import { ProposalExecutorResult } from "./types";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.group.v1";
/** EventCreateGroup is an event emitted when a group is created. */
export interface EventCreateGroup {
    /** group_id is the unique ID of the group. */
    groupId: bigint;
}
/** EventUpdateGroup is an event emitted when a group is updated. */
export interface EventUpdateGroup {
    /** group_id is the unique ID of the group. */
    groupId: bigint;
}
/** EventCreateGroupPolicy is an event emitted when a group policy is created. */
export interface EventCreateGroupPolicy {
    /** address is the account address of the group policy. */
    address: string;
}
/** EventUpdateGroupPolicy is an event emitted when a group policy is updated. */
export interface EventUpdateGroupPolicy {
    /** address is the account address of the group policy. */
    address: string;
}
/** EventSubmitProposal is an event emitted when a proposal is created. */
export interface EventSubmitProposal {
    /** proposal_id is the unique ID of the proposal. */
    proposalId: bigint;
}
/** EventWithdrawProposal is an event emitted when a proposal is withdrawn. */
export interface EventWithdrawProposal {
    /** proposal_id is the unique ID of the proposal. */
    proposalId: bigint;
}
/** EventVote is an event emitted when a voter votes on a proposal. */
export interface EventVote {
    /** proposal_id is the unique ID of the proposal. */
    proposalId: bigint;
}
/** EventExec is an event emitted when a proposal is executed. */
export interface EventExec {
    /** proposal_id is the unique ID of the proposal. */
    proposalId: bigint;
    /** result is the proposal execution result. */
    result: ProposalExecutorResult;
    /** logs contains error logs in case the execution result is FAILURE. */
    logs: string;
}
/** EventLeaveGroup is an event emitted when group member leaves the group. */
export interface EventLeaveGroup {
    /** group_id is the unique ID of the group. */
    groupId: bigint;
    /** address is the account address of the group member. */
    address: string;
}
export declare const EventCreateGroup: {
    typeUrl: string;
    encode(message: EventCreateGroup, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): EventCreateGroup;
    fromJSON(object: any): EventCreateGroup;
    toJSON(message: EventCreateGroup): unknown;
    fromPartial<I extends {
        groupId?: bigint | undefined;
    } & {
        groupId?: bigint | undefined;
    } & Record<Exclude<keyof I, "groupId">, never>>(object: I): EventCreateGroup;
};
export declare const EventUpdateGroup: {
    typeUrl: string;
    encode(message: EventUpdateGroup, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): EventUpdateGroup;
    fromJSON(object: any): EventUpdateGroup;
    toJSON(message: EventUpdateGroup): unknown;
    fromPartial<I extends {
        groupId?: bigint | undefined;
    } & {
        groupId?: bigint | undefined;
    } & Record<Exclude<keyof I, "groupId">, never>>(object: I): EventUpdateGroup;
};
export declare const EventCreateGroupPolicy: {
    typeUrl: string;
    encode(message: EventCreateGroupPolicy, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): EventCreateGroupPolicy;
    fromJSON(object: any): EventCreateGroupPolicy;
    toJSON(message: EventCreateGroupPolicy): unknown;
    fromPartial<I extends {
        address?: string | undefined;
    } & {
        address?: string | undefined;
    } & Record<Exclude<keyof I, "address">, never>>(object: I): EventCreateGroupPolicy;
};
export declare const EventUpdateGroupPolicy: {
    typeUrl: string;
    encode(message: EventUpdateGroupPolicy, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): EventUpdateGroupPolicy;
    fromJSON(object: any): EventUpdateGroupPolicy;
    toJSON(message: EventUpdateGroupPolicy): unknown;
    fromPartial<I extends {
        address?: string | undefined;
    } & {
        address?: string | undefined;
    } & Record<Exclude<keyof I, "address">, never>>(object: I): EventUpdateGroupPolicy;
};
export declare const EventSubmitProposal: {
    typeUrl: string;
    encode(message: EventSubmitProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): EventSubmitProposal;
    fromJSON(object: any): EventSubmitProposal;
    toJSON(message: EventSubmitProposal): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
    } & {
        proposalId?: bigint | undefined;
    } & Record<Exclude<keyof I, "proposalId">, never>>(object: I): EventSubmitProposal;
};
export declare const EventWithdrawProposal: {
    typeUrl: string;
    encode(message: EventWithdrawProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): EventWithdrawProposal;
    fromJSON(object: any): EventWithdrawProposal;
    toJSON(message: EventWithdrawProposal): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
    } & {
        proposalId?: bigint | undefined;
    } & Record<Exclude<keyof I, "proposalId">, never>>(object: I): EventWithdrawProposal;
};
export declare const EventVote: {
    typeUrl: string;
    encode(message: EventVote, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): EventVote;
    fromJSON(object: any): EventVote;
    toJSON(message: EventVote): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
    } & {
        proposalId?: bigint | undefined;
    } & Record<Exclude<keyof I, "proposalId">, never>>(object: I): EventVote;
};
export declare const EventExec: {
    typeUrl: string;
    encode(message: EventExec, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): EventExec;
    fromJSON(object: any): EventExec;
    toJSON(message: EventExec): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
        result?: ProposalExecutorResult | undefined;
        logs?: string | undefined;
    } & {
        proposalId?: bigint | undefined;
        result?: ProposalExecutorResult | undefined;
        logs?: string | undefined;
    } & Record<Exclude<keyof I, keyof EventExec>, never>>(object: I): EventExec;
};
export declare const EventLeaveGroup: {
    typeUrl: string;
    encode(message: EventLeaveGroup, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): EventLeaveGroup;
    fromJSON(object: any): EventLeaveGroup;
    toJSON(message: EventLeaveGroup): unknown;
    fromPartial<I extends {
        groupId?: bigint | undefined;
        address?: string | undefined;
    } & {
        groupId?: bigint | undefined;
        address?: string | undefined;
    } & Record<Exclude<keyof I, keyof EventLeaveGroup>, never>>(object: I): EventLeaveGroup;
};
