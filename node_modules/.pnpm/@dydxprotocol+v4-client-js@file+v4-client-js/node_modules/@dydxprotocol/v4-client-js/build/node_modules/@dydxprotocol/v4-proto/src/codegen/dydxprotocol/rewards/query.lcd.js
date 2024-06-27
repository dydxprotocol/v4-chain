"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.params = this.params.bind(this);
    }
    /* Queries the Params. */
    async params(_params = {}) {
        const endpoint = `dydxprotocol/v4/rewards/params`;
        return await this.req.get(endpoint);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vZHlkeHByb3RvY29sL3Jld2FyZHMvcXVlcnkubGNkLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7OztBQUVBLE1BQWEsY0FBYztJQUd6QixZQUFZLEVBQ1YsYUFBYSxFQUdkO1FBQ0MsSUFBSSxDQUFDLEdBQUcsR0FBRyxhQUFhLENBQUM7UUFDekIsSUFBSSxDQUFDLE1BQU0sR0FBRyxJQUFJLENBQUMsTUFBTSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztJQUN2QyxDQUFDO0lBQ0QseUJBQXlCO0lBR3pCLEtBQUssQ0FBQyxNQUFNLENBQUMsVUFBOEIsRUFBRTtRQUMzQyxNQUFNLFFBQVEsR0FBRyxnQ0FBZ0MsQ0FBQztRQUNsRCxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQTZCLFFBQVEsQ0FBQyxDQUFDO0lBQ2xFLENBQUM7Q0FFRjtBQW5CRCx3Q0FtQkMifQ==