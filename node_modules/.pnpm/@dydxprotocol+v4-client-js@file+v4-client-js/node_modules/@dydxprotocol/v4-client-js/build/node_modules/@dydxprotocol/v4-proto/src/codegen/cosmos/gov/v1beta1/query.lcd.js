"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
const helpers_1 = require("../../../helpers");
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.proposal = this.proposal.bind(this);
        this.proposals = this.proposals.bind(this);
        this.vote = this.vote.bind(this);
        this.votes = this.votes.bind(this);
        this.params = this.params.bind(this);
        this.deposit = this.deposit.bind(this);
        this.deposits = this.deposits.bind(this);
        this.tallyResult = this.tallyResult.bind(this);
    }
    /* Proposal queries proposal details based on ProposalID. */
    async proposal(params) {
        const endpoint = `cosmos/gov/v1beta1/proposals/${params.proposalId}`;
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
        const endpoint = `cosmos/gov/v1beta1/proposals`;
        return await this.req.get(endpoint, options);
    }
    /* Vote queries voted information based on proposalID, voterAddr. */
    async vote(params) {
        const endpoint = `cosmos/gov/v1beta1/proposals/${params.proposalId}/votes/${params.voter}`;
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
        const endpoint = `cosmos/gov/v1beta1/proposals/${params.proposalId}/votes`;
        return await this.req.get(endpoint, options);
    }
    /* Params queries all parameters of the gov module. */
    async params(params) {
        const endpoint = `cosmos/gov/v1beta1/params/${params.paramsType}`;
        return await this.req.get(endpoint);
    }
    /* Deposit queries single deposit information based on proposalID, depositor address. */
    async deposit(params) {
        const endpoint = `cosmos/gov/v1beta1/proposals/${params.proposalId}/deposits/${params.depositor}`;
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
        const endpoint = `cosmos/gov/v1beta1/proposals/${params.proposalId}/deposits`;
        return await this.req.get(endpoint, options);
    }
    /* TallyResult queries the tally of a proposal vote. */
    async tallyResult(params) {
        const endpoint = `cosmos/gov/v1beta1/proposals/${params.proposalId}/tally`;
        return await this.req.get(endpoint);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL2dvdi92MWJldGExL3F1ZXJ5LmxjZC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFBQSw4Q0FBdUQ7QUFHdkQsTUFBYSxjQUFjO0lBR3pCLFlBQVksRUFDVixhQUFhLEVBR2Q7UUFDQyxJQUFJLENBQUMsR0FBRyxHQUFHLGFBQWEsQ0FBQztRQUN6QixJQUFJLENBQUMsUUFBUSxHQUFHLElBQUksQ0FBQyxRQUFRLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ3pDLElBQUksQ0FBQyxTQUFTLEdBQUcsSUFBSSxDQUFDLFNBQVMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDM0MsSUFBSSxDQUFDLElBQUksR0FBRyxJQUFJLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNqQyxJQUFJLENBQUMsS0FBSyxHQUFHLElBQUksQ0FBQyxLQUFLLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ25DLElBQUksQ0FBQyxNQUFNLEdBQUcsSUFBSSxDQUFDLE1BQU0sQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDckMsSUFBSSxDQUFDLE9BQU8sR0FBRyxJQUFJLENBQUMsT0FBTyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUN2QyxJQUFJLENBQUMsUUFBUSxHQUFHLElBQUksQ0FBQyxRQUFRLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ3pDLElBQUksQ0FBQyxXQUFXLEdBQUcsSUFBSSxDQUFDLFdBQVcsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7SUFDakQsQ0FBQztJQUNELDREQUE0RDtJQUc1RCxLQUFLLENBQUMsUUFBUSxD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