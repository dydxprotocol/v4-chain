"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
const helpers_1 = require("../../../helpers");
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.params = this.params.bind(this);
        this.validatorDistributionInfo = this.validatorDistributionInfo.bind(this);
        this.validatorOutstandingRewards = this.validatorOutstandingRewards.bind(this);
        this.validatorCommission = this.validatorCommission.bind(this);
        this.validatorSlashes = this.validatorSlashes.bind(this);
        this.delegationRewards = this.delegationRewards.bind(this);
        this.delegationTotalRewards = this.delegationTotalRewards.bind(this);
        this.delegatorValidators = this.delegatorValidators.bind(this);
        this.delegatorWithdrawAddress = this.delegatorWithdrawAddress.bind(this);
        this.communityPool = this.communityPool.bind(this);
    }
    /* Params queries params of the distribution module. */
    async params(_params = {}) {
        const endpoint = `cosmos/distribution/v1beta1/params`;
        return await this.req.get(endpoint);
    }
    /* ValidatorDistributionInfo queries validator commission and self-delegation rewards for validator */
    async validatorDistributionInfo(params) {
        const endpoint = `cosmos/distribution/v1beta1/validators/${params.validatorAddress}`;
        return await this.req.get(endpoint);
    }
    /* ValidatorOutstandingRewards queries rewards of a validator address. */
    async validatorOutstandingRewards(params) {
        const endpoint = `cosmos/distribution/v1beta1/validators/${params.validatorAddress}/outstanding_rewards`;
        return await this.req.get(endpoint);
    }
    /* ValidatorCommission queries accumulated commission for a validator. */
    async validatorCommission(params) {
        const endpoint = `cosmos/distribution/v1beta1/validators/${params.validatorAddress}/commission`;
        return await this.req.get(endpoint);
    }
    /* ValidatorSlashes queries slash events of a validator. */
    async validatorSlashes(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.startingHeight) !== "undefined") {
            options.params.starting_height = params.startingHeight;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.endingHeight) !== "undefined") {
            options.params.ending_height = params.endingHeight;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/distribution/v1beta1/validators/${params.validatorAddress}/slashes`;
        return await this.req.get(endpoint, options);
    }
    /* DelegationRewards queries the total rewards accrued by a delegation. */
    async delegationRewards(params) {
        const endpoint = `cosmos/distribution/v1beta1/delegators/${params.delegatorAddress}/rewards/${params.validatorAddress}`;
        return await this.req.get(endpoint);
    }
    /* DelegationTotalRewards queries the total rewards accrued by each
     validator. */
    async delegationTotalRewards(params) {
        const endpoint = `cosmos/distribution/v1beta1/delegators/${params.delegatorAddress}/rewards`;
        return await this.req.get(endpoint);
    }
    /* DelegatorValidators queries the validators of a delegator. */
    async delegatorValidators(params) {
        const endpoint = `cosmos/distribution/v1beta1/delegators/${params.delegatorAddress}/validators`;
        return await this.req.get(endpoint);
    }
    /* DelegatorWithdrawAddress queries withdraw address of a delegator. */
    async delegatorWithdrawAddress(params) {
        const endpoint = `cosmos/distribution/v1beta1/delegators/${params.delegatorAddress}/withdraw_address`;
        return await this.req.get(endpoint);
    }
    /* CommunityPool queries the community pool coins. */
    async communityPool(_params = {}) {
        const endpoint = `cosmos/distribution/v1beta1/community_pool`;
        return await this.req.get(endpoint);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL2Rpc3RyaWJ1dGlvbi92MWJldGExL3F1ZXJ5LmxjZC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFBQSw4Q0FBdUQ7QUFHdkQsTUFBYSxjQUFjO0lBR3pCLFlBQVksRUFDVixhQUFhLEVBR2Q7UUFDQyxJQUFJLENBQUMsR0FBRyxHQUFHLGFBQWEsQ0FBQztRQUN6QixJQUFJLENBQUMsTUFBTSxHQUFHLElBQUksQ0FBQyxNQUFNLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ3JDLElBQUksQ0FBQyx5QkFBeUIsR0FBRyxJQUFJLENBQUMseUJBQXlCLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQzNFLElBQUksQ0FBQywyQkFBMkIsR0FBRyxJQUFJLENBQUMsMkJBQTJCLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQy9FLElBQUksQ0FBQyxtQkFBbUIsR0FBRyxJQUFJLENBQUMsbUJBQW1CLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQy9ELElBQUksQ0FBQyxnQkFBZ0IsR0FBRyxJQUFJLENBQUMsZ0JBQWdCLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ3pELElBQUksQ0FBQyxpQkFBaUIsR0FBRyxJQUFJLENBQUMsaUJBQWlCLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQzNELElBQUksQ0FBQyxzQkFBc0IsR0FBRyxJQUFJLENBQUMsc0JBQXNCLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ3JFLElBQUksQ0FBQyxtQkFBbUIsR0FBRyxJQUFJLENBQUMsbUJBQW1CLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQy9ELElB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