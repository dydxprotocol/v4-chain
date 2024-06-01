"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.params = this.params.bind(this);
        this.statsMetadata = this.statsMetadata.bind(this);
        this.globalStats = this.globalStats.bind(this);
        this.userStats = this.userStats.bind(this);
    }
    /* Queries the Params. */
    async params(_params = {}) {
        const endpoint = `dydxprotocol/v4/stats/params`;
        return await this.req.get(endpoint);
    }
    /* Queries StatsMetadata. */
    async statsMetadata(_params = {}) {
        const endpoint = `dydxprotocol/v4/stats/stats_metadata`;
        return await this.req.get(endpoint);
    }
    /* Queries GlobalStats. */
    async globalStats(_params = {}) {
        const endpoint = `dydxprotocol/v4/stats/global_stats`;
        return await this.req.get(endpoint);
    }
    /* Queries UserStats. */
    async userStats(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.user) !== "undefined") {
            options.params.user = params.user;
        }
        const endpoint = `dydxprotocol/v4/stats/user_stats`;
        return await this.req.get(endpoint, options);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vZHlkeHByb3RvY29sL3N0YXRzL3F1ZXJ5LmxjZC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFFQSxNQUFhLGNBQWM7SUFHekIsWUFBWSxFQUNWLGFBQWEsRUFHZDtRQUNDLElBQUksQ0FBQyxHQUFHLEdBQUcsYUFBYSxDQUFDO1FBQ3pCLElBQUksQ0FBQyxNQUFNLEdBQUcsSUFBSSxDQUFDLE1BQU0sQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDckMsSUFBSSxDQUFDLGFBQWEsR0FBRyxJQUFJLENBQUMsYUFBYSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNuRCxJQUFJLENBQUMsV0FBVyxHQUFHLElBQUksQ0FBQyxXQUFXLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQy9DLElBQUksQ0FBQyxTQUFTLEdBQUcsSUFBSSxDQUFDLFNBQVMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7SUFDN0MsQ0FBQztJQUNELHlCQUF5QjtJQUd6QixLQUFLLENBQUMsTUFBTSxDQUFDLFVBQThCLEVBQUU7UUFDM0MsTUFBTSxRQUFRLEdBQUcsOEJBQThCLENBQUM7UUFDaEQsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUE2QixRQUFRLENBQUMsQ0FBQztJQUNsRSxDQUFDO0lBQ0QsNEJBQTRCO0lBRzVCLEtBQUssQ0FBQyxhQUFhLENBQUMsVUFBcUMsRUFBRTtRQUN6RCxNQUFNLFFBQVEsR0FBRyxzQ0FBc0MsQ0FBQztRQUN4RCxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQW9DLFFBQVEsQ0FBQyxDQUFDO0lBQ3pFLENBQUM7SUFDRCwwQkFBMEI7SUFHMUIsS0FBSyxDQUFDLFdBQVcsQ0FBQyxVQUFtQyxFQUFFO1FBQ3JELE1BQU0sUUFBUSxHQUFHLG9DQUFvQyxDQUFDO1FBQ3RELE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBa0MsUUFBUSxDQUFDLENBQUM7SUFDdkUsQ0FBQztJQUNELHdCQUF3QjtJQUd4QixLQUFLLENBQUMsU0FBUyxDQUFDLE1BQTZCO1FBQzNDLE1BQU0sT0FBTyxHQUFRO1lBQ25CLE1BQU0sRUFBRSxFQUFFO1NBQ1gsQ0FBQztRQUVGLElBQUksT0FBTyxDQUFBLE1BQU0sYUFBTixNQUFNLHVCQUFOLE1BQU0sQ0FBRSxJQUFJLENBQUEsS0FBSyxXQUFXLEVBQUU7WUFDdkMsT0FBTyxDQUFDLE1BQU0sQ0FBQyxJQUFJLEdBQUcsTUFBTSxDQUFDLElBQUksQ0FBQztTQUNuQztRQUVELE1BQU0sUUFBUSxHQUFHLGtDQUFrQyxDQUFDO1FBQ3BELE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBZ0MsUUFBUSxFQUFFLE9BQU8sQ0FBQyxDQUFDO0lBQzlFLENBQUM7Q0FFRjtBQW5ERCx3Q0FtREMifQ==