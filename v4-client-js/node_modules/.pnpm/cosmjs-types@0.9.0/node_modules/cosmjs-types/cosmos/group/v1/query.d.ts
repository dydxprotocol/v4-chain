import { PageRequest, PageResponse } from "../../base/query/v1beta1/pagination";
import { GroupInfo, GroupPolicyInfo, GroupMember, Proposal, Vote, TallyResult } from "./types";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.group.v1";
/** QueryGroupInfoRequest is the Query/GroupInfo request type. */
export interface QueryGroupInfoRequest {
    /** group_id is the unique ID of the group. */
    groupId: bigint;
}
/** QueryGroupInfoResponse is the Query/GroupInfo response type. */
export interface QueryGroupInfoResponse {
    /** info is the GroupInfo of the group. */
    info?: GroupInfo;
}
/** QueryGroupPolicyInfoRequest is the Query/GroupPolicyInfo request type. */
export interface QueryGroupPolicyInfoRequest {
    /** address is the account address of the group policy. */
    address: string;
}
/** QueryGroupPolicyInfoResponse is the Query/GroupPolicyInfo response type. */
export interface QueryGroupPolicyInfoResponse {
    /** info is the GroupPolicyInfo of the group policy. */
    info?: GroupPolicyInfo;
}
/** QueryGroupMembersRequest is the Query/GroupMembers request type. */
export interface QueryGroupMembersRequest {
    /** group_id is the unique ID of the group. */
    groupId: bigint;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/** QueryGroupMembersResponse is the Query/GroupMembersResponse response type. */
export interface QueryGroupMembersResponse {
    /** members are the members of the group with given group_id. */
    members: GroupMember[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/** QueryGroupsByAdminRequest is the Query/GroupsByAdmin request type. */
export interface QueryGroupsByAdminRequest {
    /** admin is the account address of a group's admin. */
    admin: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/** QueryGroupsByAdminResponse is the Query/GroupsByAdminResponse response type. */
export interface QueryGroupsByAdminResponse {
    /** groups are the groups info with the provided admin. */
    groups: GroupInfo[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/** QueryGroupPoliciesByGroupRequest is the Query/GroupPoliciesByGroup request type. */
export interface QueryGroupPoliciesByGroupRequest {
    /** group_id is the unique ID of the group policy's group. */
    groupId: bigint;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/** QueryGroupPoliciesByGroupResponse is the Query/GroupPoliciesByGroup response type. */
export interface QueryGroupPoliciesByGroupResponse {
    /** group_policies are the group policies info associated with the provided group. */
    groupPolicies: GroupPolicyInfo[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/** QueryGroupPoliciesByAdminRequest is the Query/GroupPoliciesByAdmin request type. */
export interface QueryGroupPoliciesByAdminRequest {
    /** admin is the admin address of the group policy. */
    admin: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/** QueryGroupPoliciesByAdminResponse is the Query/GroupPoliciesByAdmin response type. */
export interface QueryGroupPoliciesByAdminResponse {
    /** group_policies are the group policies info with provided admin. */
    groupPolicies: GroupPolicyInfo[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/** QueryProposalRequest is the Query/Proposal request type. */
export interface QueryProposalRequest {
    /** proposal_id is the unique ID of a proposal. */
    proposalId: bigint;
}
/** QueryProposalResponse is the Query/Proposal response type. */
export interface QueryProposalResponse {
    /** proposal is the proposal info. */
    proposal?: Proposal;
}
/** QueryProposalsByGroupPolicyRequest is the Query/ProposalByGroupPolicy request type. */
export interface QueryProposalsByGroupPolicyRequest {
    /** address is the account address of the group policy related to proposals. */
    address: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/** QueryProposalsByGroupPolicyResponse is the Query/ProposalByGroupPolicy response type. */
export interface QueryProposalsByGroupPolicyResponse {
    /** proposals are the proposals with given group policy. */
    proposals: Proposal[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/** QueryVoteByProposalVoterRequest is the Query/VoteByProposalVoter request type. */
export interface QueryVoteByProposalVoterRequest {
    /** proposal_id is the unique ID of a proposal. */
    proposalId: bigint;
    /** voter is a proposal voter account address. */
    voter: string;
}
/** QueryVoteByProposalVoterResponse is the Query/VoteByProposalVoter response type. */
export interface QueryVoteByProposalVoterResponse {
    /** vote is the vote with given proposal_id and voter. */
    vote?: Vote;
}
/** QueryVotesByProposalRequest is the Query/VotesByProposal request type. */
export interface QueryVotesByProposalRequest {
    /** proposal_id is the unique ID of a proposal. */
    proposalId: bigint;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/** QueryVotesByProposalResponse is the Query/VotesByProposal response type. */
export interface QueryVotesByProposalResponse {
    /** votes are the list of votes for given proposal_id. */
    votes: Vote[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/** QueryVotesByVoterRequest is the Query/VotesByVoter request type. */
export interface QueryVotesByVoterRequest {
    /** voter is a proposal voter account address. */
    voter: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/** QueryVotesByVoterResponse is the Query/VotesByVoter response type. */
export interface QueryVotesByVoterResponse {
    /** votes are the list of votes by given voter. */
    votes: Vote[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/** QueryGroupsByMemberRequest is the Query/GroupsByMember request type. */
export interface QueryGroupsByMemberRequest {
    /** address is the group member address. */
    address: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/** QueryGroupsByMemberResponse is the Query/GroupsByMember response type. */
export interface QueryGroupsByMemberResponse {
    /** groups are the groups info with the provided group member. */
    groups: GroupInfo[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/** QueryTallyResultRequest is the Query/TallyResult request type. */
export interface QueryTallyResultRequest {
    /** proposal_id is the unique id of a proposal. */
    proposalId: bigint;
}
/** QueryTallyResultResponse is the Query/TallyResult response type. */
export interface QueryTallyResultResponse {
    /** tally defines the requested tally. */
    tally: TallyResult;
}
/**
 * QueryGroupsRequest is the Query/Groups request type.
 *
 * Since: cosmos-sdk 0.47.1
 */
export interface QueryGroupsRequest {
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/**
 * QueryGroupsResponse is the Query/Groups response type.
 *
 * Since: cosmos-sdk 0.47.1
 */
export interface QueryGroupsResponse {
    /** `groups` is all the groups present in state. */
    groups: GroupInfo[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
export declare const QueryGroupInfoRequest: {
    typeUrl: string;
    encode(message: QueryGroupInfoRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGroupInfoRequest;
    fromJSON(object: any): QueryGroupInfoRequest;
    toJSON(message: QueryGroupInfoRequest): unknown;
    fromPartial<I extends {
        groupId?: bigint | undefined;
    } & {
        groupId?: bigint | undefined;
    } & Record<Exclude<keyof I, "groupId">, never>>(object: I): QueryGroupInfoRequest;
};
export declare const QueryGroupInfoResponse: {
    typeUrl: string;
    encode(message: QueryGroupInfoResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGroupInfoResponse;
    fromJSON(object: any): QueryGroupInfoResponse;
    toJSON(message: QueryGroupInfoResponse): unknown;
    fromPartial<I extends {
        info?: {
            id?: bigint | undefined;
            admin?: string | undefined;
            metadata?: string | undefined;
            version?: bigint | undefined;
            totalWeight?: string | undefined;
            createdAt?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } | undefined;
    } & {
        info?: ({
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
            } & Record<Exclude<keyof I["info"]["createdAt"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
        } & Record<Exclude<keyof I["info"], keyof GroupInfo>, never>) | undefined;
    } & Record<Exclude<keyof I, "info">, never>>(object: I): QueryGroupInfoResponse;
};
export declare const QueryGroupPolicyInfoRequest: {
    typeUrl: string;
    encode(message: QueryGroupPolicyInfoRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGroupPolicyInfoRequest;
    fromJSON(object: any): QueryGroupPolicyInfoRequest;
    toJSON(message: QueryGroupPolicyInfoRequest): unknown;
    fromPartial<I extends {
        address?: string | undefined;
    } & {
        address?: string | undefined;
    } & Record<Exclude<keyof I, "address">, never>>(object: I): QueryGroupPolicyInfoRequest;
};
export declare const QueryGroupPolicyInfoResponse: {
    typeUrl: string;
    encode(message: QueryGroupPolicyInfoResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGroupPolicyInfoResponse;
    fromJSON(object: any): QueryGroupPolicyInfoResponse;
    toJSON(message: QueryGroupPolicyInfoResponse): unknown;
    fromPartial<I extends {
        info?: {
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
        } | undefined;
    } & {
        info?: ({
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
            } & Record<Exclude<keyof I["info"]["decisionPolicy"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
            createdAt?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["info"]["createdAt"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
        } & Record<Exclude<keyof I["info"], keyof GroupPolicyInfo>, never>) | undefined;
    } & Record<Exclude<keyof I, "info">, never>>(object: I): QueryGroupPolicyInfoResponse;
};
export declare const QueryGroupMembersRequest: {
    typeUrl: string;
    encode(message: QueryGroupMembersRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGroupMembersRequest;
    fromJSON(object: any): QueryGroupMembersRequest;
    toJSON(message: QueryGroupMembersRequest): unknown;
    fromPartial<I extends {
        groupId?: bigint | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        groupId?: bigint | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryGroupMembersRequest>, never>>(object: I): QueryGroupMembersRequest;
};
export declare const QueryGroupMembersResponse: {
    typeUrl: string;
    encode(message: QueryGroupMembersResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGroupMembersResponse;
    fromJSON(object: any): QueryGroupMembersResponse;
    toJSON(message: QueryGroupMembersResponse): unknown;
    fromPartial<I extends {
        members?: {
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
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        members?: ({
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
                } & Record<Exclude<keyof I["members"][number]["member"]["addedAt"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            } & Record<Exclude<keyof I["members"][number]["member"], keyof import("./types").Member>, never>) | undefined;
        } & Record<Exclude<keyof I["members"][number], keyof GroupMember>, never>)[] & Record<Exclude<keyof I["members"], keyof {
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
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryGroupMembersResponse>, never>>(object: I): QueryGroupMembersResponse;
};
export declare const QueryGroupsByAdminRequest: {
    typeUrl: string;
    encode(message: QueryGroupsByAdminRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGroupsByAdminRequest;
    fromJSON(object: any): QueryGroupsByAdminRequest;
    toJSON(message: QueryGroupsByAdminRequest): unknown;
    fromPartial<I extends {
        admin?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        admin?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryGroupsByAdminRequest>, never>>(object: I): QueryGroupsByAdminRequest;
};
export declare const QueryGroupsByAdminResponse: {
    typeUrl: string;
    encode(message: QueryGroupsByAdminResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGroupsByAdminResponse;
    fromJSON(object: any): QueryGroupsByAdminResponse;
    toJSON(message: QueryGroupsByAdminResponse): unknown;
    fromPartial<I extends {
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
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
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
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryGroupsByAdminResponse>, never>>(object: I): QueryGroupsByAdminResponse;
};
export declare const QueryGroupPoliciesByGroupRequest: {
    typeUrl: string;
    encode(message: QueryGroupPoliciesByGroupRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGroupPoliciesByGroupRequest;
    fromJSON(object: any): QueryGroupPoliciesByGroupRequest;
    toJSON(message: QueryGroupPoliciesByGroupRequest): unknown;
    fromPartial<I extends {
        groupId?: bigint | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        groupId?: bigint | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryGroupPoliciesByGroupRequest>, never>>(object: I): QueryGroupPoliciesByGroupRequest;
};
export declare const QueryGroupPoliciesByGroupResponse: {
    typeUrl: string;
    encode(message: QueryGroupPoliciesByGroupResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGroupPoliciesByGroupResponse;
    fromJSON(object: any): QueryGroupPoliciesByGroupResponse;
    toJSON(message: QueryGroupPoliciesByGroupResponse): unknown;
    fromPartial<I extends {
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
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
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
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryGroupPoliciesByGroupResponse>, never>>(object: I): QueryGroupPoliciesByGroupResponse;
};
export declare const QueryGroupPoliciesByAdminRequest: {
    typeUrl: string;
    encode(message: QueryGroupPoliciesByAdminRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGroupPoliciesByAdminRequest;
    fromJSON(object: any): QueryGroupPoliciesByAdminRequest;
    toJSON(message: QueryGroupPoliciesByAdminRequest): unknown;
    fromPartial<I extends {
        admin?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        admin?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryGroupPoliciesByAdminRequest>, never>>(object: I): QueryGroupPoliciesByAdminRequest;
};
export declare const QueryGroupPoliciesByAdminResponse: {
    typeUrl: string;
    encode(message: QueryGroupPoliciesByAdminResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGroupPoliciesByAdminResponse;
    fromJSON(object: any): QueryGroupPoliciesByAdminResponse;
    toJSON(message: QueryGroupPoliciesByAdminResponse): unknown;
    fromPartial<I extends {
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
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
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
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryGroupPoliciesByAdminResponse>, never>>(object: I): QueryGroupPoliciesByAdminResponse;
};
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
        } | undefined;
    } & {
        proposal?: ({
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
            proposers?: (string[] & string[] & Record<Exclude<keyof I["proposal"]["proposers"], keyof string[]>, never>) | undefined;
            submitTime?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["proposal"]["submitTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
            } & Record<Exclude<keyof I["proposal"]["finalTallyResult"], keyof TallyResult>, never>) | undefined;
            votingPeriodEnd?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["proposal"]["votingPeriodEnd"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
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
            } & Record<Exclude<keyof I["proposal"]["messages"][number], keyof import("../../../google/protobuf/any").Any>, never>)[] & Record<Exclude<keyof I["proposal"]["messages"], keyof {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            }[]>, never>) | undefined;
            title?: string | undefined;
            summary?: string | undefined;
        } & Record<Exclude<keyof I["proposal"], keyof Proposal>, never>) | undefined;
    } & Record<Exclude<keyof I, "proposal">, never>>(object: I): QueryProposalResponse;
};
export declare const QueryProposalsByGroupPolicyRequest: {
    typeUrl: string;
    encode(message: QueryProposalsByGroupPolicyRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryProposalsByGroupPolicyRequest;
    fromJSON(object: any): QueryProposalsByGroupPolicyRequest;
    toJSON(message: QueryProposalsByGroupPolicyRequest): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        address?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryProposalsByGroupPolicyRequest>, never>>(object: I): QueryProposalsByGroupPolicyRequest;
};
export declare const QueryProposalsByGroupPolicyResponse: {
    typeUrl: string;
    encode(message: QueryProposalsByGroupPolicyResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryProposalsByGroupPolicyResponse;
    fromJSON(object: any): QueryProposalsByGroupPolicyResponse;
    toJSON(message: QueryProposalsByGroupPolicyResponse): unknown;
    fromPartial<I extends {
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
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
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
            } & Record<Exclude<keyof I["proposals"][number]["finalTallyResult"], keyof TallyResult>, never>) | undefined;
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
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryProposalsByGroupPolicyResponse>, never>>(object: I): QueryProposalsByGroupPolicyResponse;
};
export declare const QueryVoteByProposalVoterRequest: {
    typeUrl: string;
    encode(message: QueryVoteByProposalVoterRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryVoteByProposalVoterRequest;
    fromJSON(object: any): QueryVoteByProposalVoterRequest;
    toJSON(message: QueryVoteByProposalVoterRequest): unknown;
    fromPartial<I extends {
        proposalId?: bigint | undefined;
        voter?: string | undefined;
    } & {
        proposalId?: bigint | undefined;
        voter?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryVoteByProposalVoterRequest>, never>>(object: I): QueryVoteByProposalVoterRequest;
};
export declare const QueryVoteByProposalVoterResponse: {
    typeUrl: string;
    encode(message: QueryVoteByProposalVoterResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryVoteByProposalVoterResponse;
    fromJSON(object: any): QueryVoteByProposalVoterResponse;
    toJSON(message: QueryVoteByProposalVoterResponse): unknown;
    fromPartial<I extends {
        vote?: {
            proposalId?: bigint | undefined;
            voter?: string | undefined;
            option?: import("./types").VoteOption | undefined;
            metadata?: string | undefined;
            submitTime?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } | undefined;
    } & {
        vote?: ({
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
            } & Record<Exclude<keyof I["vote"]["submitTime"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
        } & Record<Exclude<keyof I["vote"], keyof Vote>, never>) | undefined;
    } & Record<Exclude<keyof I, "vote">, never>>(object: I): QueryVoteByProposalVoterResponse;
};
export declare const QueryVotesByProposalRequest: {
    typeUrl: string;
    encode(message: QueryVotesByProposalRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryVotesByProposalRequest;
    fromJSON(object: any): QueryVotesByProposalRequest;
    toJSON(message: QueryVotesByProposalRequest): unknown;
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
    } & Record<Exclude<keyof I, keyof QueryVotesByProposalRequest>, never>>(object: I): QueryVotesByProposalRequest;
};
export declare const QueryVotesByProposalResponse: {
    typeUrl: string;
    encode(message: QueryVotesByProposalResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryVotesByProposalResponse;
    fromJSON(object: any): QueryVotesByProposalResponse;
    toJSON(message: QueryVotesByProposalResponse): unknown;
    fromPartial<I extends {
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
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
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
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryVotesByProposalResponse>, never>>(object: I): QueryVotesByProposalResponse;
};
export declare const QueryVotesByVoterRequest: {
    typeUrl: string;
    encode(message: QueryVotesByVoterRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryVotesByVoterRequest;
    fromJSON(object: any): QueryVotesByVoterRequest;
    toJSON(message: QueryVotesByVoterRequest): unknown;
    fromPartial<I extends {
        voter?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        voter?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryVotesByVoterRequest>, never>>(object: I): QueryVotesByVoterRequest;
};
export declare const QueryVotesByVoterResponse: {
    typeUrl: string;
    encode(message: QueryVotesByVoterResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryVotesByVoterResponse;
    fromJSON(object: any): QueryVotesByVoterResponse;
    toJSON(message: QueryVotesByVoterResponse): unknown;
    fromPartial<I extends {
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
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
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
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryVotesByVoterResponse>, never>>(object: I): QueryVotesByVoterResponse;
};
export declare const QueryGroupsByMemberRequest: {
    typeUrl: string;
    encode(message: QueryGroupsByMemberRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGroupsByMemberRequest;
    fromJSON(object: any): QueryGroupsByMemberRequest;
    toJSON(message: QueryGroupsByMemberRequest): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        address?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryGroupsByMemberRequest>, never>>(object: I): QueryGroupsByMemberRequest;
};
export declare const QueryGroupsByMemberResponse: {
    typeUrl: string;
    encode(message: QueryGroupsByMemberResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGroupsByMemberResponse;
    fromJSON(object: any): QueryGroupsByMemberResponse;
    toJSON(message: QueryGroupsByMemberResponse): unknown;
    fromPartial<I extends {
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
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
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
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryGroupsByMemberResponse>, never>>(object: I): QueryGroupsByMemberResponse;
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
export declare const QueryGroupsRequest: {
    typeUrl: string;
    encode(message: QueryGroupsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGroupsRequest;
    fromJSON(object: any): QueryGroupsRequest;
    toJSON(message: QueryGroupsRequest): unknown;
    fromPartial<I extends {
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
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
    } & Record<Exclude<keyof I, "pagination">, never>>(object: I): QueryGroupsRequest;
};
export declare const QueryGroupsResponse: {
    typeUrl: string;
    encode(message: QueryGroupsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGroupsResponse;
    fromJSON(object: any): QueryGroupsResponse;
    toJSON(message: QueryGroupsResponse): unknown;
    fromPartial<I extends {
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
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
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
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryGroupsResponse>, never>>(object: I): QueryGroupsResponse;
};
/** Query is the cosmos.group.v1 Query service. */
export interface Query {
    /** GroupInfo queries group info based on group id. */
    GroupInfo(request: QueryGroupInfoRequest): Promise<QueryGroupInfoResponse>;
    /** GroupPolicyInfo queries group policy info based on account address of group policy. */
    GroupPolicyInfo(request: QueryGroupPolicyInfoRequest): Promise<QueryGroupPolicyInfoResponse>;
    /** GroupMembers queries members of a group by group id. */
    GroupMembers(request: QueryGroupMembersRequest): Promise<QueryGroupMembersResponse>;
    /** GroupsByAdmin queries groups by admin address. */
    GroupsByAdmin(request: QueryGroupsByAdminRequest): Promise<QueryGroupsByAdminResponse>;
    /** GroupPoliciesByGroup queries group policies by group id. */
    GroupPoliciesByGroup(request: QueryGroupPoliciesByGroupRequest): Promise<QueryGroupPoliciesByGroupResponse>;
    /** GroupPoliciesByAdmin queries group policies by admin address. */
    GroupPoliciesByAdmin(request: QueryGroupPoliciesByAdminRequest): Promise<QueryGroupPoliciesByAdminResponse>;
    /** Proposal queries a proposal based on proposal id. */
    Proposal(request: QueryProposalRequest): Promise<QueryProposalResponse>;
    /** ProposalsByGroupPolicy queries proposals based on account address of group policy. */
    ProposalsByGroupPolicy(request: QueryProposalsByGroupPolicyRequest): Promise<QueryProposalsByGroupPolicyResponse>;
    /** VoteByProposalVoter queries a vote by proposal id and voter. */
    VoteByProposalVoter(request: QueryVoteByProposalVoterRequest): Promise<QueryVoteByProposalVoterResponse>;
    /** VotesByProposal queries a vote by proposal id. */
    VotesByProposal(request: QueryVotesByProposalRequest): Promise<QueryVotesByProposalResponse>;
    /** VotesByVoter queries a vote by voter. */
    VotesByVoter(request: QueryVotesByVoterRequest): Promise<QueryVotesByVoterResponse>;
    /** GroupsByMember queries groups by member address. */
    GroupsByMember(request: QueryGroupsByMemberRequest): Promise<QueryGroupsByMemberResponse>;
    /**
     * TallyResult returns the tally result of a proposal. If the proposal is
     * still in voting period, then this query computes the current tally state,
     * which might not be final. On the other hand, if the proposal is final,
     * then it simply returns the `final_tally_result` state stored in the
     * proposal itself.
     */
    TallyResult(request: QueryTallyResultRequest): Promise<QueryTallyResultResponse>;
    /**
     * Groups queries all groups in state.
     *
     * Since: cosmos-sdk 0.47.1
     */
    Groups(request?: QueryGroupsRequest): Promise<QueryGroupsResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    GroupInfo(request: QueryGroupInfoRequest): Promise<QueryGroupInfoResponse>;
    GroupPolicyInfo(request: QueryGroupPolicyInfoRequest): Promise<QueryGroupPolicyInfoResponse>;
    GroupMembers(request: QueryGroupMembersRequest): Promise<QueryGroupMembersResponse>;
    GroupsByAdmin(request: QueryGroupsByAdminRequest): Promise<QueryGroupsByAdminResponse>;
    GroupPoliciesByGroup(request: QueryGroupPoliciesByGroupRequest): Promise<QueryGroupPoliciesByGroupResponse>;
    GroupPoliciesByAdmin(request: QueryGroupPoliciesByAdminRequest): Promise<QueryGroupPoliciesByAdminResponse>;
    Proposal(request: QueryProposalRequest): Promise<QueryProposalResponse>;
    ProposalsByGroupPolicy(request: QueryProposalsByGroupPolicyRequest): Promise<QueryProposalsByGroupPolicyResponse>;
    VoteByProposalVoter(request: QueryVoteByProposalVoterRequest): Promise<QueryVoteByProposalVoterResponse>;
    VotesByProposal(request: QueryVotesByProposalRequest): Promise<QueryVotesByProposalResponse>;
    VotesByVoter(request: QueryVotesByVoterRequest): Promise<QueryVotesByVoterResponse>;
    GroupsByMember(request: QueryGroupsByMemberRequest): Promise<QueryGroupsByMemberResponse>;
    TallyResult(request: QueryTallyResultRequest): Promise<QueryTallyResultResponse>;
    Groups(request?: QueryGroupsRequest): Promise<QueryGroupsResponse>;
}
