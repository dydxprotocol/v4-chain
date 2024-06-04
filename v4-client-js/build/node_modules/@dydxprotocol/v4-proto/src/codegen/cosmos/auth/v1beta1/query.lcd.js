"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
const helpers_1 = require("../../../helpers");
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.accounts = this.accounts.bind(this);
        this.account = this.account.bind(this);
        this.accountAddressByID = this.accountAddressByID.bind(this);
        this.params = this.params.bind(this);
        this.moduleAccounts = this.moduleAccounts.bind(this);
        this.moduleAccountByName = this.moduleAccountByName.bind(this);
        this.bech32Prefix = this.bech32Prefix.bind(this);
        this.addressBytesToString = this.addressBytesToString.bind(this);
        this.addressStringToBytes = this.addressStringToBytes.bind(this);
        this.accountInfo = this.accountInfo.bind(this);
    }
    /* Accounts returns all the existing accounts.
    
     When called from another module, this query might consume a high amount of
     gas if the pagination field is incorrectly set.
    
     Since: cosmos-sdk 0.43 */
    async accounts(params = {
        pagination: undefined
    }) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/auth/v1beta1/accounts`;
        return await this.req.get(endpoint, options);
    }
    /* Account returns account details based on address. */
    async account(params) {
        const endpoint = `cosmos/auth/v1beta1/accounts/${params.address}`;
        return await this.req.get(endpoint);
    }
    /* AccountAddressByID returns account address based on account number.
    
     Since: cosmos-sdk 0.46.2 */
    async accountAddressByID(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.accountId) !== "undefined") {
            options.params.account_id = params.accountId;
        }
        const endpoint = `cosmos/auth/v1beta1/address_by_id/${params.id}`;
        return await this.req.get(endpoint, options);
    }
    /* Params queries all parameters. */
    async params(_params = {}) {
        const endpoint = `cosmos/auth/v1beta1/params`;
        return await this.req.get(endpoint);
    }
    /* ModuleAccounts returns all the existing module accounts.
    
     Since: cosmos-sdk 0.46 */
    async moduleAccounts(_params = {}) {
        const endpoint = `cosmos/auth/v1beta1/module_accounts`;
        return await this.req.get(endpoint);
    }
    /* ModuleAccountByName returns the module account info by module name */
    async moduleAccountByName(params) {
        const endpoint = `cosmos/auth/v1beta1/module_accounts/${params.name}`;
        return await this.req.get(endpoint);
    }
    /* Bech32Prefix queries bech32Prefix
    
     Since: cosmos-sdk 0.46 */
    async bech32Prefix(_params = {}) {
        const endpoint = `cosmos/auth/v1beta1/bech32`;
        return await this.req.get(endpoint);
    }
    /* AddressBytesToString converts Account Address bytes to string
    
     Since: cosmos-sdk 0.46 */
    async addressBytesToString(params) {
        const endpoint = `cosmos/auth/v1beta1/bech32/${params.addressBytes}`;
        return await this.req.get(endpoint);
    }
    /* AddressStringToBytes converts Address string to bytes
    
     Since: cosmos-sdk 0.46 */
    async addressStringToBytes(params) {
        const endpoint = `cosmos/auth/v1beta1/bech32/${params.addressString}`;
        return await this.req.get(endpoint);
    }
    /* AccountInfo queries account info which is common to all account types.
    
     Since: cosmos-sdk 0.47 */
    async accountInfo(params) {
        const endpoint = `cosmos/auth/v1beta1/account_info/${params.address}`;
        return await this.req.get(endpoint);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL2F1dGgvdjFiZXRhMS9xdWVyeS5sY2QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBQUEsOENBQXVEO0FBR3ZELE1BQWEsY0FBYztJQUd6QixZQUFZLEVBQ1YsYUFBYSxFQUdkO1FBQ0MsSUFBSSxDQUFDLEdBQUcsR0FBRyxhQUFhLENBQUM7UUFDekIsSUFBSSxDQUFDLFFBQVEsR0FBRyxJQUFJLENBQUMsUUFBUSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUN6QyxJQUFJLENBQUMsT0FBTyxHQUFHLElBQUksQ0FBQyxPQUFPLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ3ZDLElBQUksQ0FBQyxrQkFBa0IsR0FBRyxJQUFJLENBQUMsa0JBQWtCLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQzdELElBQUksQ0FBQyxNQUFNLEdBQUcsSUFBSSxDQUFDLE1BQU0sQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDckMsSUFBSSxDQUFDLGNBQWMsR0FBRyxJQUFJLENBQUMsY0FBYyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNyRCxJQUFJLENBQUMsbUJBQW1CLEdBQUcsSUFBSSxDQUFDLG1CQUFtQixDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUMvRCxJQUFJLENBQUMsWUFBWSxHQUFHLElBQUksQ0FBQyxZQUFZLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ2pELElBQUksQ0FBQyxvQkFBb0IsR0FBRyxJQUFJLENBQUMsb0JBQW9CLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ2pFLElBQUksQ0FBQyxvQkFBb0IsR0FBRyxJQUFJ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