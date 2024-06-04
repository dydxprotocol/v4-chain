"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
const helpers_1 = require("../../helpers");
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.clobPair = this.clobPair.bind(this);
        this.clobPairAll = this.clobPairAll.bind(this);
        this.equityTierLimitConfiguration = this.equityTierLimitConfiguration.bind(this);
        this.blockRateLimitConfiguration = this.blockRateLimitConfiguration.bind(this);
        this.liquidationsConfiguration = this.liquidationsConfiguration.bind(this);
    }
    /* Queries a ClobPair by id. */
    async clobPair(params) {
        const endpoint = `dydxprotocol/clob/clob_pair/${params.id}`;
        return await this.req.get(endpoint);
    }
    /* Queries a list of ClobPair items. */
    async clobPairAll(params = {
        pagination: undefined
    }) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `dydxprotocol/clob/clob_pair`;
        return await this.req.get(endpoint, options);
    }
    /* Queries EquityTierLimitConfiguration. */
    async equityTierLimitConfiguration(_params = {}) {
        const endpoint = `dydxprotocol/clob/equity_tier`;
        return await this.req.get(endpoint);
    }
    /* Queries BlockRateLimitConfiguration. */
    async blockRateLimitConfiguration(_params = {}) {
        const endpoint = `dydxprotocol/clob/block_rate`;
        return await this.req.get(endpoint);
    }
    /* Queries LiquidationsConfiguration. */
    async liquidationsConfiguration(_params = {}) {
        const endpoint = `dydxprotocol/clob/liquidations_config`;
        return await this.req.get(endpoint);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vZHlkeHByb3RvY29sL2Nsb2IvcXVlcnkubGNkLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7OztBQUFBLDJDQUFvRDtBQUdwRCxNQUFhLGNBQWM7SUFHekIsWUFBWSxFQUNWLGFBQWEsRUFHZDtRQUNDLElBQUksQ0FBQyxHQUFHLEdBQUcsYUFBYSxDQUFDO1FBQ3pCLElBQUksQ0FBQyxRQUFRLEdBQUcsSUFBSSxDQUFDLFFBQVEsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDekMsSUFBSSxDQUFDLFdBQVcsR0FBRyxJQUFJLENBQUMsV0FBVyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUMvQyxJQUFJLENBQUMsNEJBQTRCLEdBQUcsSUFBSSxDQUFDLDRCQUE0QixDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNqRixJQUFJLENBQUMsMkJBQTJCLEdBQUcsSUFBSSxDQUFDLDJCQUEyQixDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUMvRSxJQUFJLENBQUMseUJBQXlCLEdBQUcsSUFBSSxDQUFDLHlCQUF5QixDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztJQUM3RSxDQUFDO0lBQ0QsK0JBQStCO0lBRy9CLEtBQUssQ0FBQyxRQUFRLENBQUMsTUFBK0I7UUFDNUMsTUFBTSxRQUFRLEdBQUcsK0JBQStCLE1BQU0sQ0FBQyxFQUFFLEVBQUUsQ0FBQztRQUM1RCxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQStCLFFBQVEsQ0FBQyxDQUFDO0lBQ3BFLENBQUM7SUFDRCx1Q0FBdUM7SUFHdkMsS0FBSyxDQUFDLFdBQVcsQ0FBQyxTQUFrQztRQUNsRCxVQUFVLEVBQUUsU0FBUztLQUN0QjtRQUNDLE1BQU0sT0FBTyxHQUFRO1lBQ25CLE1BQU0sRUFBRSxFQUFFO1NBQ1gsQ0FBQztRQUVGLElBQUksT0FBTyxDQUFBLE1BQU0sYUFBTixNQUFNLHVCQUFOLE1BQU0sQ0FBRSxVQUFVLENBQUEsS0FBSyxXQUFXLEVBQUU7WUFDN0MsSUFBQSw2QkFBbUIsRUFBQyxPQUFPLEVBQUUsTUFBTSxDQUFDLFVBQVUsQ0FBQyxDQUFDO1NBQ2pEO1FBRUQsTUFBTSxRQUFRLEdBQUcsNkJBQTZCLENBQUM7UUFDL0MsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFrQyxRQUFRLEVBQUUsT0FBTyxDQUFDLENBQUM7SUFDaEYsQ0FBQztJQUNELDJDQUEyQztJQUczQyxLQUFLLENBQUMsNEJBQTRCLENBQUMsVUFBb0QsRUFBRTtRQUN2RixNQUFNLFFBQVEsR0FBRywrQkFBK0IsQ0FBQztRQUNqRCxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQW1ELFFBQVEsQ0FBQyxDQUFDO0lBQ3hGLENBQUM7SUFDRCwwQ0FBMEM7SUFHMUMsS0FBSyxDQUFDLDJCQUEyQixDQUFDLFVBQW1ELEVBQUU7UUFDckYsTUFBTSxRQUFRLEdBQUcsOEJBQThCLENBQUM7UUFDaEQsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFrRCxRQUFRLENBQUMsQ0FBQztJQUN2RixDQUFDO0lBQ0Qsd0NBQXdDO0lBR3hDLEtBQUssQ0FBQyx5QkFBeUIsQ0FBQyxVQUFpRCxFQUFFO1FBQ2pGLE1BQU0sUUFBUSxHQUFHLHVDQUF1QyxDQUFDO1FBQ3pELE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBZ0QsUUFBUSxDQUFDLENBQUM7SUFDckYsQ0FBQztDQUVGO0FBN0RELHdDQTZEQyJ9