"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
const helpers_1 = require("../../../helpers");
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
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
    /* GroupInfo queries group info based on group id. */
    async groupInfo(params) {
        const endpoint = `cosmos/group/v1/group_info/${params.groupId}`;
        return await this.req.get(endpoint);
    }
    /* GroupPolicyInfo queries group policy info based on account address of group policy. */
    async groupPolicyInfo(params) {
        const endpoint = `cosmos/group/v1/group_policy_info/${params.address}`;
        return await this.req.get(endpoint);
    }
    /* GroupMembers queries members of a group by group id. */
    async groupMembers(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/group/v1/group_members/${params.groupId}`;
        return await this.req.get(endpoint, options);
    }
    /* GroupsByAdmin queries groups by admin address. */
    async groupsByAdmin(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/group/v1/groups_by_admin/${params.admin}`;
        return await this.req.get(endpoint, options);
    }
    /* GroupPoliciesByGroup queries group policies by group id. */
    async groupPoliciesByGroup(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/group/v1/group_policies_by_group/${params.groupId}`;
        return await this.req.get(endpoint, options);
    }
    /* GroupPoliciesByAdmin queries group policies by admin address. */
    async groupPoliciesByAdmin(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/group/v1/group_policies_by_admin/${params.admin}`;
        return await this.req.get(endpoint, options);
    }
    /* Proposal queries a proposal based on proposal id. */
    async proposal(params) {
        const endpoint = `cosmos/group/v1/proposal/${params.proposalId}`;
        return await this.req.get(endpoint);
    }
    /* ProposalsByGroupPolicy queries proposals based on account address of group policy. */
    async proposalsByGroupPolicy(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/group/v1/proposals_by_group_policy/${params.address}`;
        return await this.req.get(endpoint, options);
    }
    /* VoteByProposalVoter queries a vote by proposal id and voter. */
    async voteByProposalVoter(params) {
        const endpoint = `cosmos/group/v1/vote_by_proposal_voter/${params.proposalId}/${params.voter}`;
        return await this.req.get(endpoint);
    }
    /* VotesByProposal queries a vote by proposal id. */
    async votesByProposal(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/group/v1/votes_by_proposal/${params.proposalId}`;
        return await this.req.get(endpoint, options);
    }
    /* VotesByVoter queries a vote by voter. */
    async votesByVoter(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/group/v1/votes_by_voter/${params.voter}`;
        return await this.req.get(endpoint, options);
    }
    /* GroupsByMember queries groups by member address. */
    async groupsByMember(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/group/v1/groups_by_member/${params.address}`;
        return await this.req.get(endpoint, options);
    }
    /* TallyResult returns the tally result of a proposal. If the proposal is
     still in voting period, then this query computes the current tally state,
     which might not be final. On the other hand, if the proposal is final,
     then it simply returns the `final_tally_result` state stored in the
     proposal itself. */
    async tallyResult(params) {
        const endpoint = `cosmos/group/v1/proposals/${params.proposalId}/tally`;
        return await this.req.get(endpoint);
    }
    /* Groups queries all groups in state.
    
     Since: cosmos-sdk 0.47.1 */
    async groups(params = {
        pagination: undefined
    }) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/group/v1/groups`;
        return await this.req.get(endpoint, options);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL2dyb3VwL3YxL3F1ZXJ5LmxjZC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFBQSw4Q0FBdUQ7QUFHdkQsTUFBYSxjQUFjO0lBR3pCLFlBQVksRUFDVixhQUFhLEVBR2Q7UUFDQyxJQUFJLENBQUMsR0FBRyxHQUFHLGFBQWEsQ0FBQztRQUN6QixJQUFJLENBQUMsU0FBUyxHQUFHLElBQUksQ0FBQyxTQUFTLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQzNDLElBQUksQ0FBQyxlQUFlLEdBQUcsSUFBSSxDQUFDLGVBQWUsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDdkQsSUFBSSxDQUFDLFlBQVksR0FBRyxJQUFJLENBQUMsWUFBWSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNqRCxJQUFJLENBQUMsYUFBYSxHQUFHLElBQUksQ0FBQyxhQUFhLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ25ELElBQUksQ0FBQyxvQkFBb0IsR0FBRyxJQUFJLENBQUMsb0JBQW9CLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ2pFLElBQUksQ0FBQyxvQkFBb0IsR0FBRyxJQUFJLENBQUMsb0JBQW9CLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ2pFLElBQUksQ0FBQyxRQUFRLEdBQUcsSUFBSSxDQUFDLFFBQVEsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDekMsSUFBSSxDQUFDLHNCQUFzQixHQUFHLElBQUksQ0FBQyxzQkFBc0IsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDckUsSUFBSSxDQUFDLG1CQUFtQixHQUFHLElBQUksQ0FB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