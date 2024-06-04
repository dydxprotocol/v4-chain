"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
const helpers_1 = require("../../../../helpers");
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.getNodeInfo = this.getNodeInfo.bind(this);
        this.getSyncing = this.getSyncing.bind(this);
        this.getLatestBlock = this.getLatestBlock.bind(this);
        this.getBlockByHeight = this.getBlockByHeight.bind(this);
        this.getLatestValidatorSet = this.getLatestValidatorSet.bind(this);
        this.getValidatorSetByHeight = this.getValidatorSetByHeight.bind(this);
        this.aBCIQuery = this.aBCIQuery.bind(this);
    }
    /* GetNodeInfo queries the current node info. */
    async getNodeInfo(_params = {}) {
        const endpoint = `cosmos/base/tendermint/v1beta1/node_info`;
        return await this.req.get(endpoint);
    }
    /* GetSyncing queries node syncing. */
    async getSyncing(_params = {}) {
        const endpoint = `cosmos/base/tendermint/v1beta1/syncing`;
        return await this.req.get(endpoint);
    }
    /* GetLatestBlock returns the latest block. */
    async getLatestBlock(_params = {}) {
        const endpoint = `cosmos/base/tendermint/v1beta1/blocks/latest`;
        return await this.req.get(endpoint);
    }
    /* GetBlockByHeight queries block for given height. */
    async getBlockByHeight(params) {
        const endpoint = `cosmos/base/tendermint/v1beta1/blocks/${params.height}`;
        return await this.req.get(endpoint);
    }
    /* GetLatestValidatorSet queries latest validator-set. */
    async getLatestValidatorSet(params = {
        pagination: undefined
    }) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/base/tendermint/v1beta1/validatorsets/latest`;
        return await this.req.get(endpoint, options);
    }
    /* GetValidatorSetByHeight queries validator-set at a given height. */
    async getValidatorSetByHeight(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/base/tendermint/v1beta1/validatorsets/${params.height}`;
        return await this.req.get(endpoint, options);
    }
    /* ABCIQuery defines a query handler that supports ABCI queries directly to the
     application, bypassing Tendermint completely. The ABCI query must contain
     a valid and supported path, including app, custom, p2p, and store.
    
     Since: cosmos-sdk 0.46 */
    async aBCIQuery(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.data) !== "undefined") {
            options.params.data = params.data;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.path) !== "undefined") {
            options.params.path = params.path;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.height) !== "undefined") {
            options.params.height = params.height;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.prove) !== "undefined") {
            options.params.prove = params.prove;
        }
        const endpoint = `cosmos/base/tendermint/v1beta1/abci_query`;
        return await this.req.get(endpoint, options);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL2Jhc2UvdGVuZGVybWludC92MWJldGExL3F1ZXJ5LmxjZC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFBQSxpREFBMEQ7QUFHMUQsTUFBYSxjQUFjO0lBR3pCLFlBQVksRUFDVixhQUFhLEVBR2Q7UUFDQyxJQUFJLENBQUMsR0FBRyxHQUFHLGFBQWEsQ0FBQztRQUN6QixJQUFJLENBQUMsV0FBVyxHQUFHLElBQUksQ0FBQyxXQUFXLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQy9DLElBQUksQ0FBQyxVQUFVLEdBQUcsSUFBSSxDQUFDLFVBQVUsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDN0MsSUFBSSxDQUFDLGNBQWMsR0FBRyxJQUFJLENBQUMsY0FBYyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNyRCxJQUFJLENBQUMsZ0JBQWdCLEdBQUcsSUFBSSxDQUFDLGdCQUFnQixDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUN6RCxJQUFJLENBQUMscUJBQXFCLEdBQUcsSUFBSSxDQUFDLHFCQUFxQixDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNuRSxJQUFJLENBQUMsdUJBQXVCLEdBQUcsSUFBSSxDQUFDLHVCQUF1QixDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUN2RSxJQUFJLENBQUMsU0FBUyxHQUFHLElBQUksQ0FBQyxTQUFTLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO0lBQzdDLENBQUM7SUFDRCxnREFBZ0Q7SUFHaEQsS0FBSyxDQUFDLFdBQVcsQ0FBQyxVQUE4QixFQUFFO1FBQ2hELE1BQU0sUUFBUSxHQUFHLDBDQUEwQyxDQUFDO1FBQzVELE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBNkIsUUFBUSxDQUFDLENBQUM7SUFDbEUsQ0FBQztJQUNELHNDQUFzQztJQUd0QyxLQUFLLENBQUMsVUFBVSxDQUFDLFVBQTZCLEVBQUU7UUFDOUMsTUFBTSxRQUFRLEdBQUcsd0NBQXdDLENBQUM7UUFDMUQsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUE0QixRQUFRLENBQUMsQ0FBQztJQUNqRSxDQUFDO0lBQ0QsOENBQThDO0lBRzlDLEtBQUssQ0FBQyxjQUFjLENBQUMsVUFBaUMsRUFBRTtRQUN0RCxNQUFNLFFBQVEsR0FBRyw4Q0FBOEMsQ0FBQztRQUNoRSxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQWdDLFFBQVEsQ0FBQyxDQUFDO0lBQ3JFLENBQUM7SUFDRCxzREFBc0Q7SUFHdEQsS0FBSyxDQUFDLGdCQUFnQixDQUFDLE1BQStCO1FBQ3BELE1BQU0sUUFBUSxHQUFHLHlDQUF5QyxNQUFNLENBQUMsTUFBTSxFQUFFLENBQUM7UUFDMUUsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFrQyxRQUFRLENBQUMsQ0FBQztJQUN2RSxDQUFDO0lBQ0QseURBQXlEO0lBR3pELEtBQUssQ0FBQyxxQkFBcUIsQ0FBQyxTQUF1QztRQUNqRSxVQUFVLEVBQUUsU0FBUztLQUN0QjtRQUNDLE1BQU0sT0FBTyxHQUFRO1lBQ25CLE1BQU0sRUFBRSxFQUFFO1NBQ1gsQ0FBQztRQUVGLElBQUksT0FBTyxDQUFBLE1BQU0sYUFBTixNQUFNLHVCQUFOLE1BQU0sQ0FBRSxVQUFVLENBQUEsS0FBSyxXQUFXLEVBQUU7WUFDN0MsSUFBQSw2QkFBbUIsRUFBQyxPQUFPLEVBQUUsTUFBTSxDQUFDLFVBQVUsQ0FBQyxDQUFDO1NBQ2pEO1FBRUQsTUFBTSxRQUFRLEdBQUcscURBQXFELENBQUM7UUFDdkUsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUF1QyxRQUFRLEVBQUUsT0FBTyxDQUFDLENBQUM7SUFDckYsQ0FBQztJQUNELHNFQUFzRTtJQUd0RSxLQUFLLENBQUMsdUJBQXVCLENBQUMsTUFBc0M7UUFDbEUsTUFBTSxPQUFPLEdBQVE7WUFDbkIsTUFBTSxFQUFFLEVBQUU7U0FDWCxDQUFDO1FBRUYsSUFBSSxPQUFPLENBQUEsTUFBTSxhQUFOLE1BQU0sdUJBQU4sTUFBTSxDQUFFLFVBQVUsQ0FBQSxLQUFLLFdBQVcsRUFBRTtZQUM3QyxJQUFBLDZCQUFtQixFQUFDLE9BQU8sRUFBRSxNQUFNLENBQUMsVUFBVSxDQUFDLENBQUM7U0FDakQ7UUFFRCxNQUFNLFFBQVEsR0FBRyxnREFBZ0QsTUFBTSxDQUFDLE1BQU0sRUFBRSxDQUFDO1FBQ2pGLE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBeUMsUUFBUSxFQUFFLE9BQU8sQ0FBQyxDQUFDO0lBQ3ZGLENBQUM7SUFDRDs7Ozs4QkFJMEI7SUFHMUIsS0FBSyxDQUFDLFNBQVMsQ0FBQyxNQUF3QjtRQUN0QyxNQUFNLE9BQU8sR0FBUTtZQUNuQixNQUFNLEVBQUUsRUFBRTtTQUNYLENBQUM7UUFFRixJQUFJLE9BQU8sQ0FBQSxNQUFNLGFBQU4sTUFBTSx1QkFBTixNQUFNLENBQUUsSUFBSSxDQUFBLEtBQUssV0FBVyxFQUFFO1lBQ3ZDLE9BQU8sQ0FBQyxNQUFNLENBQUMsSUFBSSxHQUFHLE1BQU0sQ0FBQyxJQUFJLENBQUM7U0FDbkM7UUFFRCxJQUFJLE9BQU8sQ0FBQSxNQUFNLGFBQU4sTUFBTSx1QkFBTixNQUFNLENBQUUsSUFBSSxDQUFBLEtBQUssV0FBVyxFQUFFO1lBQ3ZDLE9BQU8sQ0FBQyxNQUFNLENBQUMsSUFBSSxHQUFHLE1BQU0sQ0FBQyxJQUFJLENBQUM7U0FDbkM7UUFFRCxJQUFJLE9BQU8sQ0FBQSxNQUFNLGFBQU4sTUFBTSx1QkFBTixNQUFNLENBQUUsTUFBTSxDQUFBLEtBQUssV0FBVyxFQUFFO1lBQ3pDLE9BQU8sQ0FBQyxNQUFNLENBQUMsTUFBTSxHQUFHLE1BQU0sQ0FBQyxNQUFNLENBQUM7U0FDdkM7UUFFRCxJQUFJLE9BQU8sQ0FBQSxNQUFNLGFBQU4sTUFBTSx1QkFBTixNQUFNLENBQUUsS0FBSyxDQUFBLEtBQUssV0FBVyxFQUFFO1lBQ3hDLE9BQU8sQ0FBQyxNQUFNLENBQUMsS0FBSyxHQUFHLE1BQU0sQ0FBQyxLQUFLLENBQUM7U0FDckM7UUFFRCxNQUFNLFFBQVEsR0FBRywyQ0FBMkMsQ0FBQztRQUM3RCxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQTJCLFFBQVEsRUFBRSxPQUFPLENBQUMsQ0FBQztJQUN6RSxDQUFDO0NBRUY7QUE3R0Qsd0NBNkdDIn0=