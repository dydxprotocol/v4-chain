"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
const helpers_1 = require("../../../helpers");
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.getTx = this.getTx.bind(this);
        this.getTxsEvent = this.getTxsEvent.bind(this);
        this.getBlockWithTxs = this.getBlockWithTxs.bind(this);
    }
    /* GetTx fetches a tx by hash. */
    async getTx(params) {
        const endpoint = `cosmos/tx/v1beta1/txs/${params.hash}`;
        return await this.req.get(endpoint);
    }
    /* GetTxsEvent fetches txs by event. */
    async getTxsEvent(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.events) !== "undefined") {
            options.params.events = params.events;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.orderBy) !== "undefined") {
            options.params.order_by = params.orderBy;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.page) !== "undefined") {
            options.params.page = params.page;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.limit) !== "undefined") {
            options.params.limit = params.limit;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.query) !== "undefined") {
            options.params.query = params.query;
        }
        const endpoint = `cosmos/tx/v1beta1/txs`;
        return await this.req.get(endpoint, options);
    }
    /* GetBlockWithTxs fetches a block with decoded txs.
    
     Since: cosmos-sdk 0.45.2 */
    async getBlockWithTxs(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/tx/v1beta1/txs/block/${params.height}`;
        return await this.req.get(endpoint, options);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic2VydmljZS5sY2QuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi8uLi9ub2RlX21vZHVsZXMvQGR5ZHhwcm90b2NvbC92NC1wcm90by9zcmMvY29kZWdlbi9jb3Ntb3MvdHgvdjFiZXRhMS9zZXJ2aWNlLmxjZC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFBQSw4Q0FBdUQ7QUFHdkQsTUFBYSxjQUFjO0lBR3pCLFlBQVksRUFDVixhQUFhLEVBR2Q7UUFDQyxJQUFJLENBQUMsR0FBRyxHQUFHLGFBQWEsQ0FBQztRQUN6QixJQUFJLENBQUMsS0FBSyxHQUFHLElBQUksQ0FBQyxLQUFLLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ25DLElBQUksQ0FBQyxXQUFXLEdBQUcsSUFBSSxDQUFDLFdBQVcsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDL0MsSUFBSSxDQUFDLGVBQWUsR0FBRyxJQUFJLENBQUMsZUFBZSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztJQUN6RCxDQUFDO0lBQ0QsaUNBQWlDO0lBR2pDLEtBQUssQ0FBQyxLQUFLLENBQUMsTUFBb0I7UUFDOUIsTUFBTSxRQUFRLEdBQUcseUJBQXlCLE1BQU0sQ0FBQyxJQUFJLEVBQUUsQ0FBQztRQUN4RCxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQXVCLFFBQVEsQ0FBQyxDQUFDO0lBQzVELENBQUM7SUFDRCx1Q0FBdUM7SUFHdkMsS0FBSyxDQUFDLFdBQVcsQ0FBQyxNQUEwQjtRQUMxQyxNQUFNLE9BQU8sR0FBUTtZQUNuQixNQUFNLEVBQUUsRUFBRTtTQUNYLENBQUM7UUFFRixJQUFJLE9BQU8sQ0FBQSxNQUFNLGFBQU4sTUFBTSx1QkFBTixNQUFNLENBQUUsTUFBTSxDQUFBLEtBQUssV0FBVyxFQUFFO1lBQ3pDLE9BQU8sQ0FBQyxNQUFNLENBQUMsTUFBTSxHQUFHLE1BQU0sQ0FBQyxNQUFNLENBQUM7U0FDdkM7UUFFRCxJQUFJLE9BQU8sQ0FBQSxNQUFNLGFBQU4sTUFBTSx1QkFBTixNQUFNLENBQUUsVUFBVSxDQUFBLEtBQUssV0FBVyxFQUFFO1lBQzdDLElBQUEsNkJBQW1CLEVBQUMsT0FBTyxFQUFFLE1BQU0sQ0FBQyxVQUFVLENBQUMsQ0FBQztTQUNqRDtRQUVELElBQUksT0FBTyxDQUFBLE1BQU0sYUFBTixNQUFNLHVCQUFOLE1BQU0sQ0FBRSxPQUFPLENBQUEsS0FBSyxXQUFXLEVBQUU7WUFDMUMsT0FBTyxDQUFDLE1BQU0sQ0FBQyxRQUFRLEdBQUcsTUFBTSxDQUFDLE9BQU8sQ0FBQztTQUMxQztRQUVELElBQUksT0FBTyxDQUFBLE1BQU0sYUFBTixNQUFNLHVCQUFOLE1BQU0sQ0FBRSxJQUFJLENBQUEsS0FBSyxXQUFXLEVBQUU7WUFDdkMsT0FBTyxDQUFDLE1BQU0sQ0FBQyxJQUFJLEdBQUcsTUFBTSxDQUFDLElBQUksQ0FBQztTQUNuQztRQUVELElBQUksT0FBTyxDQUFBLE1BQU0sYUFBTixNQUFNLHVCQUFOLE1BQU0sQ0FBRSxLQUFLLENBQUEsS0FBSyxXQUFXLEVBQUU7WUFDeEMsT0FBTyxDQUFDLE1BQU0sQ0FBQyxLQUFLLEdBQUcsTUFBTSxDQUFDLEtBQUssQ0FBQztTQUNyQztRQUVELElBQUksT0FBTyxDQUFBLE1BQU0sYUFBTixNQUFNLHVCQUFOLE1BQU0sQ0FBRSxLQUFLLENBQUEsS0FBSyxXQUFXLEVBQUU7WUFDeEMsT0FBTyxDQUFDLE1BQU0sQ0FBQyxLQUFLLEdBQUcsTUFBTSxDQUFDLEtBQUssQ0FBQztTQUNyQztRQUVELE1BQU0sUUFBUSxHQUFHLHVCQUF1QixDQUFDO1FBQ3pDLE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBNkIsUUFBUSxFQUFFLE9BQU8sQ0FBQyxDQUFDO0lBQzNFLENBQUM7SUFDRDs7Z0NBRTRCO0lBRzVCLEtBQUssQ0FBQyxlQUFlLENBQUMsTUFBOEI7UUFDbEQsTUFBTSxPQUFPLEdBQVE7WUFDbkIsTUFBTSxFQUFFLEVBQUU7U0FDWCxDQUFDO1FBRUYsSUFBSSxPQUFPLENBQUEsTUFBTSxhQUFOLE1BQU0sdUJBQU4sTUFBTSxDQUFFLFVBQVUsQ0FBQSxLQUFLLFdBQVcsRUFBRTtZQUM3QyxJQUFBLDZCQUFtQixFQUFDLE9BQU8sRUFBRSxNQUFNLENBQUMsVUFBVSxDQUFDLENBQUM7U0FDakQ7UUFFRCxNQUFNLFFBQVEsR0FBRywrQkFBK0IsTUFBTSxDQUFDLE1BQU0sRUFBRSxDQUFDO1FBQ2hFLE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBaUMsUUFBUSxFQUFFLE9BQU8sQ0FBQyxDQUFDO0lBQy9FLENBQUM7Q0FFRjtBQXpFRCx3Q0F5RUMifQ==