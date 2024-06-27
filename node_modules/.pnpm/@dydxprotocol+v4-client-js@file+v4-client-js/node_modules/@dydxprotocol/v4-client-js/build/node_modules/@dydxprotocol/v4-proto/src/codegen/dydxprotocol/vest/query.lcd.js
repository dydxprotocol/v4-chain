"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.vestEntry = this.vestEntry.bind(this);
    }
    /* Queries the VestEntry. */
    async vestEntry(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.vesterAccount) !== "undefined") {
            options.params.vester_account = params.vesterAccount;
        }
        const endpoint = `dydxprotocol/v4/vest/vest_entry`;
        return await this.req.get(endpoint, options);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vZHlkeHByb3RvY29sL3Zlc3QvcXVlcnkubGNkLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7OztBQUVBLE1BQWEsY0FBYztJQUd6QixZQUFZLEVBQ1YsYUFBYSxFQUdkO1FBQ0MsSUFBSSxDQUFDLEdBQUcsR0FBRyxhQUFhLENBQUM7UUFDekIsSUFBSSxDQUFDLFNBQVMsR0FBRyxJQUFJLENBQUMsU0FBUyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztJQUM3QyxDQUFDO0lBQ0QsNEJBQTRCO0lBRzVCLEtBQUssQ0FBQyxTQUFTLENBQUMsTUFBNkI7UUFDM0MsTUFBTSxPQUFPLEdBQVE7WUFDbkIsTUFBTSxFQUFFLEVBQUU7U0FDWCxDQUFDO1FBRUYsSUFBSSxPQUFPLENBQUEsTUFBTSxhQUFOLE1BQU0sdUJBQU4sTUFBTSxDQUFFLGFBQWEsQ0FBQSxLQUFLLFdBQVcsRUFBRTtZQUNoRCxPQUFPLENBQUMsTUFBTSxDQUFDLGNBQWMsR0FBRyxNQUFNLENBQUMsYUFBYSxDQUFDO1NBQ3REO1FBRUQsTUFBTSxRQUFRLEdBQUcsaUNBQWlDLENBQUM7UUFDbkQsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFnQyxRQUFRLEVBQUUsT0FBTyxDQUFDLENBQUM7SUFDOUUsQ0FBQztDQUVGO0FBM0JELHdDQTJCQyJ9