"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
const helpers_1 = require("../../helpers");
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.perpetual = this.perpetual.bind(this);
        this.allPerpetuals = this.allPerpetuals.bind(this);
        this.allLiquidityTiers = this.allLiquidityTiers.bind(this);
        this.premiumVotes = this.premiumVotes.bind(this);
        this.premiumSamples = this.premiumSamples.bind(this);
        this.params = this.params.bind(this);
    }
    /* Queries a Perpetual by id. */
    async perpetual(params) {
        const endpoint = `dydxprotocol/perpetuals/perpetual/${params.id}`;
        return await this.req.get(endpoint);
    }
    /* Queries a list of Perpetual items. */
    async allPerpetuals(params = {
        pagination: undefined
    }) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `dydxprotocol/perpetuals/perpetual`;
        return await this.req.get(endpoint, options);
    }
    /* Queries a list of LiquidityTiers. */
    async allLiquidityTiers(params = {
        pagination: undefined
    }) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `dydxprotocol/perpetuals/liquidity_tiers`;
        return await this.req.get(endpoint, options);
    }
    /* Queries a list of premium votes. */
    async premiumVotes(_params = {}) {
        const endpoint = `dydxprotocol/perpetuals/premium_votes`;
        return await this.req.get(endpoint);
    }
    /* Queries a list of premium samples. */
    async premiumSamples(_params = {}) {
        const endpoint = `dydxprotocol/perpetuals/premium_samples`;
        return await this.req.get(endpoint);
    }
    /* Queries the perpetual params. */
    async params(_params = {}) {
        const endpoint = `dydxprotocol/perpetuals/params`;
        return await this.req.get(endpoint);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vZHlkeHByb3RvY29sL3BlcnBldHVhbHMvcXVlcnkubGNkLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7OztBQUFBLDJDQUFvRDtBQUdwRCxNQUFhLGNBQWM7SUFHekIsWUFBWSxFQUNWLGFBQWEsRUFHZDtRQUNDLElBQUksQ0FBQyxHQUFHLEdBQUcsYUFBYSxDQUFDO1FBQ3pCLElBQUksQ0FBQyxTQUFTLEdBQUcsSUFBSSxDQUFDLFNBQVMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDM0MsSUFBSSxDQUFDLGFBQWEsR0FBRyxJQUFJLENBQUMsYUFBYSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNuRCxJQUFJLENBQUMsaUJBQWlCLEdBQUcsSUFBSSxDQUFDLGlCQUFpQixDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUMzRCxJQUFJLENBQUMsWUFBWSxHQUFHLElBQUksQ0FBQyxZQUFZLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ2pELElBQUksQ0FBQyxjQUFjLEdBQUcsSUFBSSxDQUFDLGNBQWMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDckQsSUFBSSxDQUFDLE1BQU0sR0FBRyxJQUFJLENBQUMsTUFBTSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztJQUN2QyxDQUFDO0lBQ0QsZ0NBQWdDO0lBR2hDLEtBQUssQ0FBQyxTQUFTLENBQUMsTUFBNkI7UUFDM0MsTUFBTSxRQUFRLEdBQUcscUNBQXFDLE1BQU0sQ0FBQyxFQUFFLEVBQUUsQ0FBQztRQUNsRSxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQWdDLFFBQVEsQ0FBQyxDQUFDO0lBQ3JFLENBQUM7SUFDRCx3Q0FBd0M7SUFHeEMsS0FBSyxDQUFDLGFBQWEsQ0FBQyxTQUFvQztRQUN0RCxVQUFVLEVBQUUsU0FBUztLQUN0QjtRQUNDLE1BQU0sT0FBTyxHQUFRO1lBQ25CLE1BQU0sRUFBRSxFQUFFO1NBQ1gsQ0FBQztRQUVGLElBQUksT0FBTyxDQUFBLE1BQU0sYUFBTixNQUFNLHVCQUFOLE1BQU0sQ0FBRSxVQUFVLENBQUEsS0FBSyxXQUFXLEVBQUU7WUFDN0MsSUFBQSw2QkFBbUIsRUFBQyxPQUFPLEVBQUUsTUFBTSxDQUFDLFVBQVUsQ0FBQyxDQUFDO1NBQ2pEO1FBRUQsTUFBTSxRQUFRLEdBQUcsbUNBQW1DLENBQUM7UUFDckQsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFvQyxRQUFRLEVBQUUsT0FBTyxDQUFDLENBQUM7SUFDbEYsQ0FBQztJQUNELHVDQUF1QztJQUd2QyxLQUFLLENBQUMsaUJBQWlCLENBQUMsU0FBd0M7UUFDOUQsVUFBVSxFQUFFLFNBQVM7S0FDdEI7UUFDQyxNQUFNLE9BQU8sR0FBUTtZQUNuQixNQUFNLEVBQUUsRUFBRTtTQUNYLENBQUM7UUFFRixJQUFJLE9BQU8sQ0FBQSxNQUFNLGFBQU4sTUFBTSx1QkFBTixNQUFNLENBQUUsVUFBVSxDQUFBLEtBQUssV0FBVyxFQUFFO1lBQzdDLElBQUEsNkJBQW1CLEVBQUMsT0FBTyxFQUFFLE1BQU0sQ0FBQyxVQUFVLENBQUMsQ0FBQztTQUNqRDtRQUVELE1BQU0sUUFBUSxHQUFHLHlDQUF5QyxDQUFDO1FBQzNELE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBd0MsUUFBUSxFQUFFLE9BQU8sQ0FBQyxDQUFDO0lBQ3RGLENBQUM7SUFDRCxzQ0FBc0M7SUFHdEMsS0FBSyxDQUFDLFlBQVksQ0FBQyxVQUFvQyxFQUFFO1FBQ3ZELE1BQU0sUUFBUSxHQUFHLHVDQUF1QyxDQUFDO1FBQ3pELE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBbUMsUUFBUSxDQUFDLENBQUM7SUFDeEUsQ0FBQztJQUNELHdDQUF3QztJQUd4QyxLQUFLLENBQUMsY0FBYyxDQUFDLFVBQXNDLEVBQUU7UUFDM0QsTUFBTSxRQUFRLEdBQUcseUNBQXlDLENBQUM7UUFDM0QsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFxQyxRQUFRLENBQUMsQ0FBQztJQUMxRSxDQUFDO0lBQ0QsbUNBQW1DO0lBR25DLEtBQUssQ0FBQyxNQUFNLENBQUMsVUFBOEIsRUFBRTtRQUMzQyxNQUFNLFFBQVEsR0FBRyxnQ0FBZ0MsQ0FBQztRQUNsRCxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQTZCLFFBQVEsQ0FBQyxDQUFDO0lBQ2xFLENBQUM7Q0FFRjtBQS9FRCx3Q0ErRUMifQ==