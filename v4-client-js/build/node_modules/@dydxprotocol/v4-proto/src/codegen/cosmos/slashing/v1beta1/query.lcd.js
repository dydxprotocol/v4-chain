"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
const helpers_1 = require("../../../helpers");
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.params = this.params.bind(this);
        this.signingInfo = this.signingInfo.bind(this);
        this.signingInfos = this.signingInfos.bind(this);
    }
    /* Params queries the parameters of slashing module */
    async params(_params = {}) {
        const endpoint = `cosmos/slashing/v1beta1/params`;
        return await this.req.get(endpoint);
    }
    /* SigningInfo queries the signing info of given cons address */
    async signingInfo(params) {
        const endpoint = `cosmos/slashing/v1beta1/signing_infos/${params.consAddress}`;
        return await this.req.get(endpoint);
    }
    /* SigningInfos queries signing info of all validators */
    async signingInfos(params = {
        pagination: undefined
    }) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/slashing/v1beta1/signing_infos`;
        return await this.req.get(endpoint, options);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL3NsYXNoaW5nL3YxYmV0YTEvcXVlcnkubGNkLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7OztBQUFBLDhDQUF1RDtBQUd2RCxNQUFhLGNBQWM7SUFHekIsWUFBWSxFQUNWLGFBQWEsRUFHZDtRQUNDLElBQUksQ0FBQyxHQUFHLEdBQUcsYUFBYSxDQUFDO1FBQ3pCLElBQUksQ0FBQyxNQUFNLEdBQUcsSUFBSSxDQUFDLE1BQU0sQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDckMsSUFBSSxDQUFDLFdBQVcsR0FBRyxJQUFJLENBQUMsV0FBVyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUMvQyxJQUFJLENBQUMsWUFBWSxHQUFHLElBQUksQ0FBQyxZQUFZLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO0lBQ25ELENBQUM7SUFDRCxzREFBc0Q7SUFHdEQsS0FBSyxDQUFDLE1BQU0sQ0FBQyxVQUE4QixFQUFFO1FBQzNDLE1BQU0sUUFBUSxHQUFHLGdDQUFnQyxDQUFDO1FBQ2xELE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBNkIsUUFBUSxDQUFDLENBQUM7SUFDbEUsQ0FBQztJQUNELGdFQUFnRTtJQUdoRSxLQUFLLENBQUMsV0FBVyxDQUFDLE1BQStCO1FBQy9DLE1BQU0sUUFBUSxHQUFHLHlDQUF5QyxNQUFNLENBQUMsV0FBVyxFQUFFLENBQUM7UUFDL0UsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFrQyxRQUFRLENBQUMsQ0FBQztJQUN2RSxDQUFDO0lBQ0QseURBQXlEO0lBR3pELEtBQUssQ0FBQyxZQUFZLENBQUMsU0FBbUM7UUFDcEQsVUFBVSxFQUFFLFNBQVM7S0FDdEI7UUFDQyxNQUFNLE9BQU8sR0FBUTtZQUNuQixNQUFNLEVBQUUsRUFBRTtTQUNYLENBQUM7UUFFRixJQUFJLE9BQU8sQ0FBQSxNQUFNLGFBQU4sTUFBTSx1QkFBTixNQUFNLENBQUUsVUFBVSxDQUFBLEtBQUssV0FBVyxFQUFFO1lBQzdDLElBQUEsNkJBQW1CLEVBQUMsT0FBTyxFQUFFLE1BQU0sQ0FBQyxVQUFVLENBQUMsQ0FBQztTQUNqRDtRQUVELE1BQU0sUUFBUSxHQUFHLHVDQUF1QyxDQUFDO1FBQ3pELE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBbUMsUUFBUSxFQUFFLE9BQU8sQ0FBQyxDQUFDO0lBQ2pGLENBQUM7Q0FFRjtBQTdDRCx3Q0E2Q0MifQ==