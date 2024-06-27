"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
const helpers_1 = require("../../../helpers");
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.constitution = this.constitution.bind(this);
        this.proposal = this.proposal.bind(this);
        this.proposals = this.proposals.bind(this);
        this.vote = this.vote.bind(this);
        this.votes = this.votes.bind(this);
        this.params = this.params.bind(this);
        this.deposit = this.deposit.bind(this);
        this.deposits = this.deposits.bind(this);
        this.tallyResult = this.tallyResult.bind(this);
    }
    /* Constitution queries the chain's constitution. */
    async constitution(_params = {}) {
        const endpoint = `cosmos/gov/v1/constitution`;
        return await this.req.get(endpoint);
    }
    /* Proposal queries proposal details based on ProposalID. */
    async proposal(params) {
        const endpoint = `cosmos/gov/v1/proposals/${params.proposalId}`;
        return await this.req.get(endpoint);
    }
    /* Proposals queries all proposals based on given status. */
    async proposals(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.proposalStatus) !== "undefined") {
            options.params.proposal_status = params.proposalStatus;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.voter) !== "undefined") {
            options.params.voter = params.voter;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.depositor) !== "undefined") {
            options.params.depositor = params.depositor;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/gov/v1/proposals`;
        return await this.req.get(endpoint, options);
    }
    /* Vote queries voted information based on proposalID, voterAddr. */
    async vote(params) {
        const endpoint = `cosmos/gov/v1/proposals/${params.proposalId}/votes/${params.voter}`;
        return await this.req.get(endpoint);
    }
    /* Votes queries votes of a given proposal. */
    async votes(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/gov/v1/proposals/${params.proposalId}/votes`;
        return await this.req.get(endpoint, options);
    }
    /* Params queries all parameters of the gov module. */
    async params(params) {
        const endpoint = `cosmos/gov/v1/params/${params.paramsType}`;
        return await this.req.get(endpoint);
    }
    /* Deposit queries single deposit information based on proposalID, depositAddr. */
    async deposit(params) {
        const endpoint = `cosmos/gov/v1/proposals/${params.proposalId}/deposits/${params.depositor}`;
        return await this.req.get(endpoint);
    }
    /* Deposits queries all deposits of a single proposal. */
    async deposits(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/gov/v1/proposals/${params.proposalId}/deposits`;
        return await this.req.get(endpoint, options);
    }
    /* TallyResult queries the tally of a proposal vote. */
    async tallyResult(params) {
        const endpoint = `cosmos/gov/v1/proposals/${params.proposalId}/tally`;
        return await this.req.get(endpoint);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL2dvdi92MS9xdWVyeS5sY2QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBQUEsOENBQXVEO0FBR3ZELE1BQWEsY0FBYztJQUd6QixZQUFZLEVBQ1YsYUFBYSxFQUdkO1FBQ0MsSUFBSSxDQUFDLEdBQUcsR0FBRyxhQUFhLENBQUM7UUFDekIsSUFBSSxDQUFDLFlBQVksR0FBRyxJQUFJLENBQUMsWUFBWSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNqRCxJQUFJLENBQUMsUUFBUSxHQUFHLElBQUksQ0FBQyxRQUFRLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ3pDLElBQUksQ0FBQyxTQUFTLEdBQUcsSUFBSSxDQUFDLFNBQVMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDM0MsSUFBSSxDQUFDLElBQUksR0FBRyxJQUFJLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNqQyxJQUFJLENBQUMsS0FBSyxHQUFHLElBQUksQ0FBQyxLQUFLLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ25DLElBQUksQ0FBQyxNQUFNLEdBQUcsSUFBSSxDQUFDLE1BQU0sQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDckMsSUFBSSxDQUFDLE9BQU8sR0FBRyxJQUFJLENBQUMsT0FBTyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUN2QyxJQUFJLENBQUMsUUFBUSxHQUFHLElBQUksQ0FBQyxRQUFRLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ3pDLElBQUksQ0FBQyxXQUFXLEdBQUcsSUFBSSxDQUFDLFdBQVcsQ0FBQyxJQUFJ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