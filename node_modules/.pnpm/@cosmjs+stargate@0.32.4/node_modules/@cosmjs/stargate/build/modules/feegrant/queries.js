"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.setupFeegrantExtension = void 0;
const query_1 = require("cosmjs-types/cosmos/feegrant/v1beta1/query");
const queryclient_1 = require("../../queryclient");
function setupFeegrantExtension(base) {
    // Use this service to get easy typed access to query methods
    // This cannot be used for proof verification
    const rpc = (0, queryclient_1.createProtobufRpcClient)(base);
    const queryService = new query_1.QueryClientImpl(rpc);
    return {
        feegrant: {
            allowance: async (granter, grantee) => {
                const response = await queryService.Allowance({
                    granter: granter,
                    grantee: grantee,
                });
                return response;
            },
            allowances: async (grantee, paginationKey) => {
                const response = await queryService.Allowances({
                    grantee: grantee,
                    pagination: (0, queryclient_1.createPagination)(paginationKey),
                });
                return response;
            },
        },
    };
}
exports.setupFeegrantExtension = setupFeegrantExtension;
//# sourceMappingURL=queries.js.map