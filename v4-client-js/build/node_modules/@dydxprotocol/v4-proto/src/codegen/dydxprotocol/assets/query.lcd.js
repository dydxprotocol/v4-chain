"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
const helpers_1 = require("../../helpers");
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.asset = this.asset.bind(this);
        this.allAssets = this.allAssets.bind(this);
    }
    /* Queries a Asset by id. */
    async asset(params) {
        const endpoint = `dydxprotocol/assets/asset/${params.id}`;
        return await this.req.get(endpoint);
    }
    /* Queries a list of Asset items. */
    async allAssets(params = {
        pagination: undefined
    }) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `dydxprotocol/assets/asset`;
        return await this.req.get(endpoint, options);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vZHlkeHByb3RvY29sL2Fzc2V0cy9xdWVyeS5sY2QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBQUEsMkNBQW9EO0FBR3BELE1BQWEsY0FBYztJQUd6QixZQUFZLEVBQ1YsYUFBYSxFQUdkO1FBQ0MsSUFBSSxDQUFDLEdBQUcsR0FBRyxhQUFhLENBQUM7UUFDekIsSUFBSSxDQUFDLEtBQUssR0FBRyxJQUFJLENBQUMsS0FBSyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNuQyxJQUFJLENBQUMsU0FBUyxHQUFHLElBQUksQ0FBQyxTQUFTLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO0lBQzdDLENBQUM7SUFDRCw0QkFBNEI7SUFHNUIsS0FBSyxDQUFDLEtBQUssQ0FBQyxNQUF5QjtRQUNuQyxNQUFNLFFBQVEsR0FBRyw2QkFBNkIsTUFBTSxDQUFDLEVBQUUsRUFBRSxDQUFDO1FBQzFELE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBNEIsUUFBUSxDQUFDLENBQUM7SUFDakUsQ0FBQztJQUNELG9DQUFvQztJQUdwQyxLQUFLLENBQUMsU0FBUyxDQUFDLFNBQWdDO1FBQzlDLFVBQVUsRUFBRSxTQUFTO0tBQ3RCO1FBQ0MsTUFBTSxPQUFPLEdBQVE7WUFDbkIsTUFBTSxFQUFFLEVBQUU7U0FDWCxDQUFDO1FBRUYsSUFBSSxPQUFPLENBQUEsTUFBTSxhQUFOLE1BQU0sdUJBQU4sTUFBTSxDQUFFLFVBQVUsQ0FBQSxLQUFLLFdBQVcsRUFBRTtZQUM3QyxJQUFBLDZCQUFtQixFQUFDLE9BQU8sRUFBRSxNQUFNLENBQUMsVUFBVSxDQUFDLENBQUM7U0FDakQ7UUFFRCxNQUFNLFFBQVEsR0FBRywyQkFBMkIsQ0FBQztRQUM3QyxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQWdDLFFBQVEsRUFBRSxPQUFPLENBQUMsQ0FBQztJQUM5RSxDQUFDO0NBRUY7QUFyQ0Qsd0NBcUNDIn0=