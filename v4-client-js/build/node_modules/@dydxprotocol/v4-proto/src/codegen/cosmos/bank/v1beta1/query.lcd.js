"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
const helpers_1 = require("../../../helpers");
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.balance = this.balance.bind(this);
        this.allBalances = this.allBalances.bind(this);
        this.spendableBalances = this.spendableBalances.bind(this);
        this.spendableBalanceByDenom = this.spendableBalanceByDenom.bind(this);
        this.totalSupply = this.totalSupply.bind(this);
        this.supplyOf = this.supplyOf.bind(this);
        this.params = this.params.bind(this);
        this.denomMetadata = this.denomMetadata.bind(this);
        this.denomMetadataByQueryString = this.denomMetadataByQueryString.bind(this);
        this.denomsMetadata = this.denomsMetadata.bind(this);
        this.denomOwners = this.denomOwners.bind(this);
        this.sendEnabled = this.sendEnabled.bind(this);
    }
    /* Balance queries the balance of a single coin for a single account. */
    async balance(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.denom) !== "undefined") {
            options.params.denom = params.denom;
        }
        const endpoint = `cosmos/bank/v1beta1/balances/${params.address}/by_denom`;
        return await this.req.get(endpoint, options);
    }
    /* AllBalances queries the balance of all coins for a single account.
    
     When called from another module, this query might consume a high amount of
     gas if the pagination field is incorrectly set. */
    async allBalances(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.resolveDenom) !== "undefined") {
            options.params.resolve_denom = params.resolveDenom;
        }
        const endpoint = `cosmos/bank/v1beta1/balances/${params.address}`;
        return await this.req.get(endpoint, options);
    }
    /* SpendableBalances queries the spendable balance of all coins for a single
     account.
    
     When called from another module, this query might consume a high amount of
     gas if the pagination field is incorrectly set.
    
     Since: cosmos-sdk 0.46 */
    async spendableBalances(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/bank/v1beta1/spendable_balances/${params.address}`;
        return await this.req.get(endpoint, options);
    }
    /* SpendableBalanceByDenom queries the spendable balance of a single denom for
     a single account.
    
     When called from another module, this query might consume a high amount of
     gas if the pagination field is incorrectly set.
    
     Since: cosmos-sdk 0.47 */
    async spendableBalanceByDenom(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.denom) !== "undefined") {
            options.params.denom = params.denom;
        }
        const endpoint = `cosmos/bank/v1beta1/spendable_balances/${params.address}/by_denom`;
        return await this.req.get(endpoint, options);
    }
    /* TotalSupply queries the total supply of all coins.
    
     When called from another module, this query might consume a high amount of
     gas if the pagination field is incorrectly set. */
    async totalSupply(params = {
        pagination: undefined
    }) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/bank/v1beta1/supply`;
        return await this.req.get(endpoint, options);
    }
    /* SupplyOf queries the supply of a single coin.
    
     When called from another module, this query might consume a high amount of
     gas if the pagination field is incorrectly set. */
    async supplyOf(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.denom) !== "undefined") {
            options.params.denom = params.denom;
        }
        const endpoint = `cosmos/bank/v1beta1/supply/by_denom`;
        return await this.req.get(endpoint, options);
    }
    /* Params queries the parameters of x/bank module. */
    async params(_params = {}) {
        const endpoint = `cosmos/bank/v1beta1/params`;
        return await this.req.get(endpoint);
    }
    /* DenomsMetadata queries the client metadata of a given coin denomination. */
    async denomMetadata(params) {
        const endpoint = `cosmos/bank/v1beta1/denoms_metadata/${params.denom}`;
        return await this.req.get(endpoint);
    }
    /* DenomsMetadata queries the client metadata of a given coin denomination. */
    async denomMetadataByQueryString(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.denom) !== "undefined") {
            options.params.denom = params.denom;
        }
        const endpoint = `cosmos/bank/v1beta1/denoms_metadata_by_query_string`;
        return await this.req.get(endpoint, options);
    }
    /* DenomsMetadata queries the client metadata for all registered coin
     denominations. */
    async denomsMetadata(params = {
        pagination: undefined
    }) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/bank/v1beta1/denoms_metadata`;
        return await this.req.get(endpoint, options);
    }
    /* DenomOwners queries for all account addresses that own a particular token
     denomination.
    
     When called from another module, this query might consume a high amount of
     gas if the pagination field is incorrectly set.
    
     Since: cosmos-sdk 0.46 */
    async denomOwners(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/bank/v1beta1/denom_owners/${params.denom}`;
        return await this.req.get(endpoint, options);
    }
    /* SendEnabled queries for SendEnabled entries.
    
     This query only returns denominations that have specific SendEnabled settings.
     Any denomination that does not have a specific setting will use the default
     params.default_send_enabled, and will not be returned by this query.
    
     Since: cosmos-sdk 0.47 */
    async sendEnabled(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.denoms) !== "undefined") {
            options.params.denoms = params.denoms;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/bank/v1beta1/send_enabled`;
        return await this.req.get(endpoint, options);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL2JhbmsvdjFiZXRhMS9xdWVyeS5sY2QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBQUEsOENBQXVEO0FBR3ZELE1BQWEsY0FBYztJQUd6QixZQUFZLEVBQ1YsYUFBYSxFQUdkO1FBQ0MsSUFBSSxDQUFDLEdBQUcsR0FBRyxhQUFhLENBQUM7UUFDekIsSUFBSSxDQUFDLE9BQU8sR0FBRyxJQUFJLENBQUMsT0FBTyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUN2QyxJQUFJLENBQUMsV0FBVyxHQUFHLElBQUksQ0FBQyxXQUFXLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQy9DLElBQUksQ0FBQyxpQkFBaUIsR0FBRyxJQUFJLENBQUMsaUJBQWlCLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQzNELElBQUksQ0FBQyx1QkFBdUIsR0FBRyxJQUFJLENBQUMsdUJBQXVCLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ3ZFLElBQUksQ0FBQyxXQUFXLEdBQUcsSUFBSSxDQUFDLFdBQVcsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDL0MsSUFBSSxDQUFDLFFBQVEsR0FBRyxJQUFJLENBQUMsUUFBUSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUN6QyxJQUFJLENBQUMsTUFBTSxHQUFHLElBQUksQ0FBQyxNQUFNLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ3JDLElBQUksQ0FBQyxhQUFhLEdBQUcsSUFBSSxDQUFDLGFBQWEsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDbkQsSUFBSSxDQUFDLDBCQUEwQixHQUFHLElBQUksQ0FB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