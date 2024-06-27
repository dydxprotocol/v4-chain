import { GroupInfo, GroupMember, GroupPolicyInfo, Proposal, Vote } from "./types";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.group.v1";
/** GenesisState defines the group module's genesis state. */
export interface GenesisState {
    /**
     * group_seq is the group table orm.Sequence,
     * it is used to get the next group ID.
     */
    groupSeq: bigint;
    /** groups is the list of groups info. */
    groups: GroupInfo[];
    /** group_members is the list of groups members. */
    groupMembers: GroupMember[];
    /**
     * group_policy_seq is the group policy table orm.Sequence,
     * it is used to generate the next group policy account address.
     */
    groupPolicySeq: bigint;
    /** group_policies is the list of group policies info. */
    groupPolicies: GroupPolicyInfo[];
    /**
     * proposal_seq is the proposal table orm.Sequence,
     * it is used to get the next proposal ID.
     */
    proposalSeq: bigint;
    /** proposals is the list of proposals. */
    proposals: Proposal[];
    /** votes is the list of votes. */
    votes: Vote[];
}
export declare const GenesisState: {
    typeUrl: string;
    encode(message: GenesisState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial<I extends {
        groupSeq?: bigint | undefined;
        groups?: {
            id?: bigint | undefined;
            admin?: string | undefined;
            metadata?: string | undefined;
            version?: bigint | undefined;
            totalWeight?: string | undefined;
            createdAt?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[] | undefined;
        groupMembers?: {
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
        }[] | undefined;
        groupPolicySeq?: bigint | undefined;
        groupPolicies?: {
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
        }[] | undefined;
        proposalSeq?: bigint | undefined;
        proposals?: {
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
            status?: import("./types").ProposalStatus | undefined;
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
            executorResult?: import("./types").ProposalExecutorResult | undefined;
            messages?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
            title?: string | undefined;
            summary?: string | undefined;
        }[] | undefined;
        votes?: {
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            option?: import("./types").VoteOption | undefined;
            metadata?: string | undefined;
            submitTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[] | undefined;
    } & {
        groupSeq?: bigint | undefined;
        groups?: ({
            id?: bigint | undefined;
            admin?: string | undefined;
            metadata?: string | undefined;
            version?: bigint | undefined;
            totalWeight?: string | undefined;
            createdAt?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[] & ({
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
            } & Record<Exclude<keyof I["groups"][number]["createdAt"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
        } & Record<Exclude<keyof I["groups"][number], keyof GroupInfo>, never>)[] & Record<Exclude<keyof I["groups"], keyof {
            id?: bigint | undefined;
            admin?: string | undefined;
            metadata?: string | undefined;
            version?: bigint | undefined;
            totalWeight?: string | undefined;
            createdAt?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[]>, never>) | undefined;
        groupMembers?: ({
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
        }[] & ({
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
                } & Record<Exclude<keyof I["groupMembers"][number]["member"]["addedAt"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            } & Record<Exclude<keyof I["groupMembers"][number]["member"], keyof import("./types").Member>, never>) | undefined;
        } & Record<Exclude<keyof I["groupMembers"][number], keyof GroupMember>, never>)[] & Record<Exclude<keyof I["groupMembers"], keyof {
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
        }[]>, never>) | undefined;
        groupPolicySeq?: bigint | undefined;
        groupPolicies?: ({
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
        }[] & ({
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
            } & Record<Exclude<keyof I["groupPolicies"][number]["decisionPolicy"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
            createdAt?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["groupPolicies"][number]["createdAt"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
        } & Record<Exclude<keyof I["groupPolicies"][number], keyof GroupPolicyInfo>, never>)[] & Record<Exclude<keyof I["groupPolicies"], keyof {
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
        }[]>, never>) | undefined;
        proposalSeq?: bigint | undefined;
        proposals?: ({
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
            status?: import("./types").ProposalStatus | undefined;
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
            executorResult?: import("./types").ProposalExecutorResult | undefined;
            messages?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
            title?: string | undefined;
            summary?: string | undefined;
        }[] & ({
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
            status?: import("./types").ProposalStatus | undefined;
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
            executorResult?: import("./types").ProposalExecutorResult | undefined;
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
            proposers?: (string[] & string[] & Record<Exclude<keyof I["proposals"][number]["proposers"], keyof string[]>, never>) | undefined;
            submitTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["proposals"][number]["submitTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            groupVersion?: bigint | undefined;
            groupPolicyVersion?: bigint | undefined;
            status?: import("./types").ProposalStatus | undefined;
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
            } & Record<Exclude<keyof I["proposals"][number]["finalTallyResult"], keyof import("./types").TallyResult>, never>) | undefined;
            votingPeriodEnd?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["proposals"][number]["votingPeriodEnd"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            executorResult?: import("./types").ProposalExecutorResult | undefined;
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
            title?: string | undefined;
            summary?: string | undefined;
        } & Record<Exclude<keyof I["proposals"][number], keyof Proposal>, never>)[] & Record<Exclude<keyof I["proposals"], keyof {
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
            status?: import("./types").ProposalStatus | undefined;
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
            executorResult?: import("./types").ProposalExecutorResult | undefined;
            messages?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
            title?: string | undefined;
            summary?: string | undefined;
        }[]>, never>) | undefined;
        votes?: ({
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            option?: import("./types").VoteOption | undefined;
            metadata?: string | undefined;
            submitTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[] & ({
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            option?: import("./types").VoteOption | undefined;
            metadata?: string | undefined;
            submitTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } & {
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            option?: import("./types").VoteOption | undefined;
            metadata?: string | undefined;
            submitTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["votes"][number]["submitTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
        } & Record<Exclude<keyof I["votes"][number], keyof Vote>, never>)[] & Record<Exclude<keyof I["votes"], keyof {
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            option?: import("./types").VoteOption | undefined;
            metadata?: string | undefined;
            submitTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GenesisState>, never>>(object: I): GenesisState;
};
