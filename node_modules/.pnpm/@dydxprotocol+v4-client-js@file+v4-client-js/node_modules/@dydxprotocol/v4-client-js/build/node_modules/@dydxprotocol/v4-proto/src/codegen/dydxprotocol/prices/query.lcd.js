"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
const helpers_1 = require("../../helpers");
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.marketPrice = this.marketPrice.bind(this);
        this.allMarketPrices = this.allMarketPrices.bind(this);
        this.marketParam = this.marketParam.bind(this);
        this.allMarketParams = this.allMarketParams.bind(this);
    }
    /* Queries a MarketPrice by id. */
    async marketPrice(params) {
        const endpoint = `dydxprotocol/prices/market/${params.id}`;
        return await this.req.get(endpoint);
    }
    /* Queries a list of MarketPrice items. */
    async allMarketPrices(params = {
        pagination: undefined
    }) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `dydxprotocol/prices/market`;
        return await this.req.get(endpoint, options);
    }
    /* Queries a MarketParam by id. */
    async marketParam(params) {
        const endpoint = `dydxprotocol/prices/params/market/${params.id}`;
        return await this.req.get(endpoint);
    }
    /* Queries a list of MarketParam items. */
    async allMarketParams(params = {
        pagination: undefined
    }) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `dydxprotocol/prices/params/market`;
        return await this.req.get(endpoint, options);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vZHlkeHByb3RvY29sL3ByaWNlcy9xdWVyeS5sY2QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBQUEsMkNBQW9EO0FBR3BELE1BQWEsY0FBYztJQUd6QixZQUFZLEVBQ1YsYUFBYSxFQUdkO1FBQ0MsSUFBSSxDQUFDLEdBQUcsR0FBRyxhQUFhLENBQUM7UUFDekIsSUFBSSxDQUFDLFdBQVcsR0FBRyxJQUFJLENBQUMsV0FBVyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUMvQyxJQUFJLENBQUMsZUFBZSxHQUFHLElBQUksQ0FBQyxlQUFlLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ3ZELElBQUksQ0FBQyxXQUFXLEdBQUcsSUFBSSxDQUFDLFdBQVcsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDL0MsSUFBSSxDQUFDLGVBQWUsR0FBRyxJQUFJLENBQUMsZUFBZSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztJQUN6RCxDQUFDO0lBQ0Qsa0NBQWtDO0lBR2xDLEtBQUssQ0FBQyxXQUFXLENBQUMsTUFBK0I7UUFDL0MsTUFBTSxRQUFRLEdBQUcsOEJBQThCLE1BQU0sQ0FBQyxFQUFFLEVBQUUsQ0FBQztRQUMzRCxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQWtDLFFBQVEsQ0FBQyxDQUFDO0lBQ3ZFLENBQUM7SUFDRCwwQ0FBMEM7SUFHMUMsS0FBSyxDQUFDLGVBQWUsQ0FBQyxTQUFzQztRQUMxRCxVQUFVLEVBQUUsU0FBUztLQUN0QjtRQUNDLE1BQU0sT0FBTyxHQUFRO1lBQ25CLE1BQU0sRUFBRSxFQUFFO1NBQ1gsQ0FBQztRQUVGLElBQUksT0FBTyxDQUFBLE1BQU0sYUFBTixNQUFNLHVCQUFOLE1BQU0sQ0FBRSxVQUFVLENBQUEsS0FBSyxXQUFXLEVBQUU7WUFDN0MsSUFBQSw2QkFBbUIsRUFBQyxPQUFPLEVBQUUsTUFBTSxDQUFDLFVBQVUsQ0FBQyxDQUFDO1NBQ2pEO1FBRUQsTUFBTSxRQUFRLEdBQUcsNEJBQTRCLENBQUM7UUFDOUMsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFzQyxRQUFRLEVBQUUsT0FBTyxDQUFDLENBQUM7SUFDcEYsQ0FBQztJQUNELGtDQUFrQztJQUdsQyxLQUFLLENBQUMsV0FBVyxDQUFDLE1BQStCO1FBQy9DLE1BQU0sUUFBUSxHQUFHLHFDQUFxQyxNQUFNLENBQUMsRUFBRSxFQUFFLENBQUM7UUFDbEUsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFrQyxRQUFRLENBQUMsQ0FBQztJQUN2RSxDQUFDO0lBQ0QsMENBQTBDO0lBRzFDLEtBQUssQ0FBQyxlQUFlLENBQUMsU0FBc0M7UUFDMUQsVUFBVSxFQUFFLFNBQVM7S0FDdEI7UUFDQyxNQUFNLE9BQU8sR0FBUTtZQUNuQixNQUFNLEVBQUUsRUFBRTtTQUNYLENBQUM7UUFFRixJQUFJLE9BQU8sQ0FBQSxNQUFNLGFBQU4sTUFBTSx1QkFBTixNQUFNLENBQUUsVUFBVSxDQUFBLEtBQUssV0FBVyxFQUFFO1lBQzdDLElBQUEsNkJBQW1CLEVBQUMsT0FBTyxFQUFFLE1BQU0sQ0FBQyxVQUFVLENBQUMsQ0FBQztTQUNqRDtRQUVELE1BQU0sUUFBUSxHQUFHLG1DQUFtQyxDQUFDO1FBQ3JELE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBc0MsUUFBUSxFQUFFLE9BQU8sQ0FBQyxDQUFDO0lBQ3BGLENBQUM7Q0FFRjtBQS9ERCx3Q0ErREMifQ==