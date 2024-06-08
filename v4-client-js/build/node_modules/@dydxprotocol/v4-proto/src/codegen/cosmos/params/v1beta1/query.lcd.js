"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.params = this.params.bind(this);
        this.subspaces = this.subspaces.bind(this);
    }
    /* Params queries a specific parameter of a module, given its subspace and
     key. */
    async params(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.subspace) !== "undefined") {
            options.params.subspace = params.subspace;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.key) !== "undefined") {
            options.params.key = params.key;
        }
        const endpoint = `cosmos/params/v1beta1/params`;
        return await this.req.get(endpoint, options);
    }
    /* Subspaces queries for all registered subspaces and all keys for a subspace.
    
     Since: cosmos-sdk 0.46 */
    async subspaces(_params = {}) {
        const endpoint = `cosmos/params/v1beta1/subspaces`;
        return await this.req.get(endpoint);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL3BhcmFtcy92MWJldGExL3F1ZXJ5LmxjZC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFFQSxNQUFhLGNBQWM7SUFHekIsWUFBWSxFQUNWLGFBQWEsRUFHZDtRQUNDLElBQUksQ0FBQyxHQUFHLEdBQUcsYUFBYSxDQUFDO1FBQ3pCLElBQUksQ0FBQyxNQUFNLEdBQUcsSUFBSSxDQUFDLE1BQU0sQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDckMsSUFBSSxDQUFDLFNBQVMsR0FBRyxJQUFJLENBQUMsU0FBUyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztJQUM3QyxDQUFDO0lBQ0Q7WUFDUTtJQUdSLEtBQUssQ0FBQyxNQUFNLENBQUMsTUFBMEI7UUFDckMsTUFBTSxPQUFPLEdBQVE7WUFDbkIsTUFBTSxFQUFFLEVBQUU7U0FDWCxDQUFDO1FBRUYsSUFBSSxPQUFPLENBQUEsTUFBTSxhQUFOLE1BQU0sdUJBQU4sTUFBTSxDQUFFLFFBQVEsQ0FBQSxLQUFLLFdBQVcsRUFBRTtZQUMzQyxPQUFPLENBQUMsTUFBTSxDQUFDLFFBQVEsR0FBRyxNQUFNLENBQUMsUUFBUSxDQUFDO1NBQzNDO1FBRUQsSUFBSSxPQUFPLENBQUEsTUFBTSxhQUFOLE1BQU0sdUJBQU4sTUFBTSxDQUFFLEdBQUcsQ0FBQSxLQUFLLFdBQVcsRUFBRTtZQUN0QyxPQUFPLENBQUMsTUFBTSxDQUFDLEdBQUcsR0FBRyxNQUFNLENBQUMsR0FBRyxDQUFDO1NBQ2pDO1FBRUQsTUFBTSxRQUFRLEdBQUcsOEJBQThCLENBQUM7UUFDaEQsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUE2QixRQUFRLEVBQUUsT0FBTyxDQUFDLENBQUM7SUFDM0UsQ0FBQztJQUNEOzs4QkFFMEI7SUFHMUIsS0FBSyxDQUFDLFNBQVMsQ0FBQyxVQUFpQyxFQUFFO1FBQ2pELE1BQU0sUUFBUSxHQUFHLGlDQUFpQyxDQUFDO1FBQ25ELE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBZ0MsUUFBUSxDQUFDLENBQUM7SUFDckUsQ0FBQztDQUVGO0FBMUNELHdDQTBDQyJ9