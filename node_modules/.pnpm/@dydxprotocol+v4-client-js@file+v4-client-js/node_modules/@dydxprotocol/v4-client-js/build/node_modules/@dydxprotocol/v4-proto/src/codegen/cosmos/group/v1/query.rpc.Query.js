"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.createRpcQueryExtension = exports.QueryClientImpl = void 0;
const _m0 = __importStar(require("protobufjs/minimal"));
const stargate_1 = require("@cosmjs/stargate");
const query_1 = require("./query");
class QueryClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
        this.groupInfo = this.groupInfo.bind(this);
        this.groupPolicyInfo = this.groupPolicyInfo.bind(this);
        this.groupMembers = this.groupMembers.bind(this);
        this.groupsByAdmin = this.groupsByAdmin.bind(this);
        this.groupPoliciesByGroup = this.groupPoliciesByGroup.bind(this);
        this.groupPoliciesByAdmin = this.groupPoliciesByAdmin.bind(this);
        this.proposal = this.proposal.bind(this);
        this.proposalsByGroupPolicy = this.proposalsByGroupPolicy.bind(this);
        this.voteByProposalVoter = this.voteByProposalVoter.bind(this);
        this.votesByProposal = this.votesByProposal.bind(this);
        this.votesByVoter = this.votesByVoter.bind(this);
        this.groupsByMember = this.groupsByMember.bind(this);
        this.tallyResult = this.tallyResult.bind(this);
        this.groups = this.groups.bind(this);
    }
    groupInfo(request) {
        const data = query_1.QueryGroupInfoRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.group.v1.Query", "GroupInfo", data);
        return promise.then(data => query_1.QueryGroupInfoResponse.decode(new _m0.Reader(data)));
    }
    groupPolicyInfo(request) {
        const data = query_1.QueryGroupPolicyInfoRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.group.v1.Query", "GroupPolicyInfo", data);
        return promise.then(data => query_1.QueryGroupPolicyInfoResponse.decode(new _m0.Reader(data)));
    }
    groupMembers(request) {
        const data = query_1.QueryGroupMembersRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.group.v1.Query", "GroupMembers", data);
        return promise.then(data => query_1.QueryGroupMembersResponse.decode(new _m0.Reader(data)));
    }
    groupsByAdmin(request) {
        const data = query_1.QueryGroupsByAdminRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.group.v1.Query", "GroupsByAdmin", data);
        return promise.then(data => query_1.QueryGroupsByAdminResponse.decode(new _m0.Reader(data)));
    }
    groupPoliciesByGroup(request) {
        const data = query_1.QueryGroupPoliciesByGroupRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.group.v1.Query", "GroupPoliciesByGroup", data);
        return promise.then(data => query_1.QueryGroupPoliciesByGroupResponse.decode(new _m0.Reader(data)));
    }
    groupPoliciesByAdmin(request) {
        const data = query_1.QueryGroupPoliciesByAdminRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.group.v1.Query", "GroupPoliciesByAdmin", data);
        return promise.then(data => query_1.QueryGroupPoliciesByAdminResponse.decode(new _m0.Reader(data)));
    }
    proposal(request) {
        const data = query_1.QueryProposalRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.group.v1.Query", "Proposal", data);
        return promise.then(data => query_1.QueryProposalResponse.decode(new _m0.Reader(data)));
    }
    proposalsByGroupPolicy(request) {
        const data = query_1.QueryProposalsByGroupPolicyRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.group.v1.Query", "ProposalsByGroupPolicy", data);
        return promise.then(data => query_1.QueryProposalsByGroupPolicyResponse.decode(new _m0.Reader(data)));
    }
    voteByProposalVoter(request) {
        const data = query_1.QueryVoteByProposalVoterRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.group.v1.Query", "VoteByProposalVoter", data);
        return promise.then(data => query_1.QueryVoteByProposalVoterResponse.decode(new _m0.Reader(data)));
    }
    votesByProposal(request) {
        const data = query_1.QueryVotesByProposalRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.group.v1.Query", "VotesByProposal", data);
        return promise.then(data => query_1.QueryVotesByProposalResponse.decode(new _m0.Reader(data)));
    }
    votesByVoter(request) {
        const data = query_1.QueryVotesByVoterRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.group.v1.Query", "VotesByVoter", data);
        return promise.then(data => query_1.QueryVotesByVoterResponse.decode(new _m0.Reader(data)));
    }
    groupsByMember(request) {
        const data = query_1.QueryGroupsByMemberRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.group.v1.Query", "GroupsByMember", data);
        return promise.then(data => query_1.QueryGroupsByMemberResponse.decode(new _m0.Reader(data)));
    }
    tallyResult(request) {
        const data = query_1.QueryTallyResultRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.group.v1.Query", "TallyResult", data);
        return promise.then(data => query_1.QueryTallyResultResponse.decode(new _m0.Reader(data)));
    }
    groups(request = {
        pagination: undefined
    }) {
        const data = query_1.QueryGroupsRequest.encode(request).finish();
        const promise = this.rpc.request("cosmos.group.v1.Query", "Groups", data);
        return promise.then(data => query_1.QueryGroupsResponse.decode(new _m0.Reader(data)));
    }
}
exports.QueryClientImpl = QueryClientImpl;
const createRpcQueryExtension = (base) => {
    const rpc = (0, stargate_1.createProtobufRpcClient)(base);
    const queryService = new QueryClientImpl(rpc);
    return {
        groupInfo(request) {
            return queryService.groupInfo(request);
        },
        groupPolicyInfo(request) {
            return queryService.groupPolicyInfo(request);
        },
        groupMembers(request) {
            return queryService.groupMembers(request);
        },
        groupsByAdmin(request) {
            return queryService.groupsByAdmin(request);
        },
        groupPoliciesByGroup(request) {
            return queryService.groupPoliciesByGroup(request);
        },
        groupPoliciesByAdmin(request) {
            return queryService.groupPoliciesByAdmin(request);
        },
        proposal(request) {
            return queryService.proposal(request);
        },
        proposalsByGroupPolicy(request) {
            return queryService.proposalsByGroupPolicy(request);
        },
        voteByProposalVoter(request) {
            return queryService.voteByProposalVoter(request);
        },
        votesByProposal(request) {
            return queryService.votesByProposal(request);
        },
        votesByVoter(request) {
            return queryService.votesByVoter(request);
        },
        groupsByMember(request) {
            return queryService.groupsByMember(request);
        },
        tallyResult(request) {
            return queryService.tallyResult(request);
        },
        groups(request) {
            return queryService.groups(request);
        }
    };
};
exports.createRpcQueryExtension = createRpcQueryExtension;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkucnBjLlF1ZXJ5LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL2dyb3VwL3YxL3F1ZXJ5LnJwYy5RdWVyeS50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztBQUNBLHdEQUEwQztBQUMxQywrQ0FBd0U7QUFDeEUsbUNBQXV6QjtBQXdEdnpCLE1BQWEsZUFBZTtJQUcxQixZQUFZLEdBQVE7UUFDbEIsSUFBSSxDQUFDLEdBQUcsR0FBRyxHQUFHLENBQUM7UUFDZixJQUFJLENBQUMsU0FBUyxHQUFHLElBQUksQ0FBQyxTQUFTLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQzNDLElBQUksQ0FBQyxlQUFlLEdBQUcsSUFBSSxDQUFDLGVBQWUsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDdkQsSUFBSSxDQUFDLFlBQVksR0FBRyxJQUFJLENBQUMsWUFBWSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNqRCxJQUFJLENBQUMsYUFBYSxHQUFHLElBQUksQ0FBQyxhQUFhLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ25ELElBQUksQ0FBQyxvQkFBb0IsR0FBRyxJQUFJLENBQUMsb0JBQW9CLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ2pFLElBQUksQ0FBQyxvQkFBb0IsR0FBRyxJQUFJLENBQUMsb0JBQW9CLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ2pFLElBQUksQ0FBQyxRQUFRLEdBQUcsSUFBSSxDQUFDLFFBQVEsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDekMsSUFBSSxDQUFDLHNCQUFzQixHQUFHLElBQUksQ0FBQyxzQkFBc0IsQ0FBQyxJQUFJLENB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