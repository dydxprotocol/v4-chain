"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.perpetualFeeParams = this.perpetualFeeParams.bind(this);
        this.userFeeTier = this.userFeeTier.bind(this);
    }
    /* Queries the PerpetualFeeParams. */
    async perpetualFeeParams(_params = {}) {
        const endpoint = `dydxprotocol/v4/feetiers/perpetual_fee_params`;
        return await this.req.get(endpoint);
    }
    /* Queries a user's fee tier */
    async userFeeTier(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.user) !== "undefined") {
            options.params.user = params.user;
        }
        const endpoint = `dydxprotocol/v4/feetiers/user_fee_tier`;
        return await this.req.get(endpoint, options);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vZHlkeHByb3RvY29sL2ZlZXRpZXJzL3F1ZXJ5LmxjZC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFFQSxNQUFhLGNBQWM7SUFHekIsWUFBWSxFQUNWLGFBQWEsRUFHZDtRQUNDLElBQUksQ0FBQyxHQUFHLEdBQUcsYUFBYSxDQUFDO1FBQ3pCLElBQUksQ0FBQyxrQkFBa0IsR0FBRyxJQUFJLENBQUMsa0JBQWtCLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQzdELElBQUksQ0FBQyxXQUFXLEdBQUcsSUFBSSxDQUFDLFdBQVcsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7SUFDakQsQ0FBQztJQUNELHFDQUFxQztJQUdyQyxLQUFLLENBQUMsa0JBQWtCLENBQUMsVUFBMEMsRUFBRTtRQUNuRSxNQUFNLFFBQVEsR0FBRywrQ0FBK0MsQ0FBQztRQUNqRSxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQXlDLFFBQVEsQ0FBQyxDQUFDO0lBQzlFLENBQUM7SUFDRCwrQkFBK0I7SUFHL0IsS0FBSyxDQUFDLFdBQVcsQ0FBQyxNQUErQjtRQUMvQyxNQUFNLE9BQU8sR0FBUTtZQUNuQixNQUFNLEVBQUUsRUFBRTtTQUNYLENBQUM7UUFFRixJQUFJLE9BQU8sQ0FBQSxNQUFNLGFBQU4sTUFBTSx1QkFBTixNQUFNLENBQUUsSUFBSSxDQUFBLEtBQUssV0FBVyxFQUFFO1lBQ3ZDLE9BQU8sQ0FBQyxNQUFNLENBQUMsSUFBSSxHQUFHLE1BQU0sQ0FBQyxJQUFJLENBQUM7U0FDbkM7UUFFRCxNQUFNLFFBQVEsR0FBRyx3Q0FBd0MsQ0FBQztRQUMxRCxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQWtDLFFBQVEsRUFBRSxPQUFPLENBQUMsQ0FBQztJQUNoRixDQUFDO0NBRUY7QUFuQ0Qsd0NBbUNDIn0=