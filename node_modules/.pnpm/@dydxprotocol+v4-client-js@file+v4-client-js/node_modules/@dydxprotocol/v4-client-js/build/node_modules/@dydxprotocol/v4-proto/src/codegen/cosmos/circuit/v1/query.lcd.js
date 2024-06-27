"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
const helpers_1 = require("../../../helpers");
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.account = this.account.bind(this);
        this.accounts = this.accounts.bind(this);
        this.disabledList = this.disabledList.bind(this);
    }
    /* Account returns account permissions. */
    async account(params) {
        const endpoint = `cosmos/circuit/v1/accounts/${params.address}`;
        return await this.req.get(endpoint);
    }
    /* Account returns account permissions. */
    async accounts(params = {
        pagination: undefined
    }) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/circuit/v1/accounts`;
        return await this.req.get(endpoint, options);
    }
    /* DisabledList returns a list of disabled message urls */
    async disabledList(_params = {}) {
        const endpoint = `cosmos/circuit/v1/disable_list`;
        return await this.req.get(endpoint);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL2NpcmN1aXQvdjEvcXVlcnkubGNkLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7OztBQUFBLDhDQUF1RDtBQUd2RCxNQUFhLGNBQWM7SUFHekIsWUFBWSxFQUNWLGFBQWEsRUFHZDtRQUNDLElBQUksQ0FBQyxHQUFHLEdBQUcsYUFBYSxDQUFDO1FBQ3pCLElBQUksQ0FBQyxPQUFPLEdBQUcsSUFBSSxDQUFDLE9BQU8sQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDdkMsSUFBSSxDQUFDLFFBQVEsR0FBRyxJQUFJLENBQUMsUUFBUSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUN6QyxJQUFJLENBQUMsWUFBWSxHQUFHLElBQUksQ0FBQyxZQUFZLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO0lBQ25ELENBQUM7SUFDRCwwQ0FBMEM7SUFHMUMsS0FBSyxDQUFDLE9BQU8sQ0FBQyxNQUEyQjtRQUN2QyxNQUFNLFFBQVEsR0FBRyw4QkFBOEIsTUFBTSxDQUFDLE9BQU8sRUFBRSxDQUFDO1FBQ2hFLE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBeUIsUUFBUSxDQUFDLENBQUM7SUFDOUQsQ0FBQztJQUNELDBDQUEwQztJQUcxQyxLQUFLLENBQUMsUUFBUSxDQUFDLFNBQStCO1FBQzVDLFVBQVUsRUFBRSxTQUFTO0tBQ3RCO1FBQ0MsTUFBTSxPQUFPLEdBQVE7WUFDbkIsTUFBTSxFQUFFLEVBQUU7U0FDWCxDQUFDO1FBRUYsSUFBSSxPQUFPLENBQUEsTUFBTSxhQUFOLE1BQU0sdUJBQU4sTUFBTSxDQUFFLFVBQVUsQ0FBQSxLQUFLLFdBQVcsRUFBRTtZQUM3QyxJQUFBLDZCQUFtQixFQUFDLE9BQU8sRUFBRSxNQUFNLENBQUMsVUFBVSxDQUFDLENBQUM7U0FDakQ7UUFFRCxNQUFNLFFBQVEsR0FBRyw0QkFBNEIsQ0FBQztRQUM5QyxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQTBCLFFBQVEsRUFBRSxPQUFPLENBQUMsQ0FBQztJQUN4RSxDQUFDO0lBQ0QsMERBQTBEO0lBRzFELEtBQUssQ0FBQyxZQUFZLENBQUMsVUFBb0MsRUFBRTtRQUN2RCxNQUFNLFFBQVEsR0FBRyxnQ0FBZ0MsQ0FBQztRQUNsRCxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQThCLFFBQVEsQ0FBQyxDQUFDO0lBQ25FLENBQUM7Q0FFRjtBQTdDRCx3Q0E2Q0MifQ==