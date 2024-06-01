"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
const helpers_1 = require("../../../helpers");
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.validators = this.validators.bind(this);
        this.validator = this.validator.bind(this);
        this.validatorDelegations = this.validatorDelegations.bind(this);
        this.validatorUnbondingDelegations = this.validatorUnbondingDelegations.bind(this);
        this.delegation = this.delegation.bind(this);
        this.unbondingDelegation = this.unbondingDelegation.bind(this);
        this.delegatorDelegations = this.delegatorDelegations.bind(this);
        this.delegatorUnbondingDelegations = this.delegatorUnbondingDelegations.bind(this);
        this.redelegations = this.redelegations.bind(this);
        this.delegatorValidators = this.delegatorValidators.bind(this);
        this.delegatorValidator = this.delegatorValidator.bind(this);
        this.historicalInfo = this.historicalInfo.bind(this);
        this.pool = this.pool.bind(this);
        this.params = this.params.bind(this);
    }
    /* Validators queries all validators that match the given status.
    
     When called from another module, this query might consume a high amount of
     gas if the pagination field is incorrectly set. */
    async validators(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.status) !== "undefined") {
            options.params.status = params.status;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/staking/v1beta1/validators`;
        return await this.req.get(endpoint, options);
    }
    /* Validator queries validator info for given validator address. */
    async validator(params) {
        const endpoint = `cosmos/staking/v1beta1/validators/${params.validatorAddr}`;
        return await this.req.get(endpoint);
    }
    /* ValidatorDelegations queries delegate info for given validator.
    
     When called from another module, this query might consume a high amount of
     gas if the pagination field is incorrectly set. */
    async validatorDelegations(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/staking/v1beta1/validators/${params.validatorAddr}/delegations`;
        return await this.req.get(endpoint, options);
    }
    /* ValidatorUnbondingDelegations queries unbonding delegations of a validator.
    
     When called from another module, this query might consume a high amount of
     gas if the pagination field is incorrectly set. */
    async validatorUnbondingDelegations(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/staking/v1beta1/validators/${params.validatorAddr}/unbonding_delegations`;
        return await this.req.get(endpoint, options);
    }
    /* Delegation queries delegate info for given validator delegator pair. */
    async delegation(params) {
        const endpoint = `cosmos/staking/v1beta1/validators/${params.validatorAddr}/delegations/${params.delegatorAddr}`;
        return await this.req.get(endpoint);
    }
    /* UnbondingDelegation queries unbonding info for given validator delegator
     pair. */
    async unbondingDelegation(params) {
        const endpoint = `cosmos/staking/v1beta1/validators/${params.validatorAddr}/delegations/${params.delegatorAddr}/unbonding_delegation`;
        return await this.req.get(endpoint);
    }
    /* DelegatorDelegations queries all delegations of a given delegator address.
    
     When called from another module, this query might consume a high amount of
     gas if the pagination field is incorrectly set. */
    async delegatorDelegations(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/staking/v1beta1/delegations/${params.delegatorAddr}`;
        return await this.req.get(endpoint, options);
    }
    /* DelegatorUnbondingDelegations queries all unbonding delegations of a given
     delegator address.
    
     When called from another module, this query might consume a high amount of
     gas if the pagination field is incorrectly set. */
    async delegatorUnbondingDelegations(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/staking/v1beta1/delegators/${params.delegatorAddr}/unbonding_delegations`;
        return await this.req.get(endpoint, options);
    }
    /* Redelegations queries redelegations of given address.
    
     When called from another module, this query might consume a high amount of
     gas if the pagination field is incorrectly set. */
    async redelegations(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.srcValidatorAddr) !== "undefined") {
            options.params.src_validator_addr = params.srcValidatorAddr;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.dstValidatorAddr) !== "undefined") {
            options.params.dst_validator_addr = params.dstValidatorAddr;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/staking/v1beta1/delegators/${params.delegatorAddr}/redelegations`;
        return await this.req.get(endpoint, options);
    }
    /* DelegatorValidators queries all validators info for given delegator
     address.
    
     When called from another module, this query might consume a high amount of
     gas if the pagination field is incorrectly set. */
    async delegatorValidators(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/staking/v1beta1/delegators/${params.delegatorAddr}/validators`;
        return await this.req.get(endpoint, options);
    }
    /* DelegatorValidator queries validator info for given delegator validator
     pair. */
    async delegatorValidator(params) {
        const endpoint = `cosmos/staking/v1beta1/delegators/${params.delegatorAddr}/validators/${params.validatorAddr}`;
        return await this.req.get(endpoint);
    }
    /* HistoricalInfo queries the historical info for given height. */
    async historicalInfo(params) {
        const endpoint = `cosmos/staking/v1beta1/historical_info/${params.height}`;
        return await this.req.get(endpoint);
    }
    /* Pool queries the pool info. */
    async pool(_params = {}) {
        const endpoint = `cosmos/staking/v1beta1/pool`;
        return await this.req.get(endpoint);
    }
    /* Parameters queries the staking parameters. */
    async params(_params = {}) {
        const endpoint = `cosmos/staking/v1beta1/params`;
        return await this.req.get(endpoint);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL3N0YWtpbmcvdjFiZXRhMS9xdWVyeS5sY2QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBQUEsOENBQXVEO0FBR3ZELE1BQWEsY0FBYztJQUd6QixZQUFZLEVBQ1YsYUFBYSxFQUdkO1FBQ0MsSUFBSSxDQUFDLEdBQUcsR0FBRyxhQUFhLENBQUM7UUFDekIsSUFBSSxDQUFDLFVBQVUsR0FBRyxJQUFJLENBQUMsVUFBVSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUM3QyxJQUFJLENBQUMsU0FBUyxHQUFHLElBQUksQ0FBQyxTQUFTLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQzNDLElBQUksQ0FBQyxvQkFBb0IsR0FBRyxJQUFJLENBQUMsb0JBQW9CLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ2pFLElBQUksQ0FBQyw2QkFBNkIsR0FBRyxJQUFJLENBQUMsNkJBQTZCLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ25GLElBQUksQ0FBQyxVQUFVLEdBQUcsSUFBSSxDQUFDLFVBQVUsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDN0MsSUFBSSxDQUFDLG1CQUFtQixHQUFHLElBQUksQ0FBQyxtQkFBbUIsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDL0QsSUFBSSxDQUFDLG9CQUFvQixHQUFHLElBQUksQ0FBQyxvQkFBb0IsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDakUsSUFBSSxDQUFDLDZCQUE2QixHQUFHLElBQUksQ0FBQyw2QkFBNkIsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDbkYsSUFBSSxDQUFDLGFBQWEs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