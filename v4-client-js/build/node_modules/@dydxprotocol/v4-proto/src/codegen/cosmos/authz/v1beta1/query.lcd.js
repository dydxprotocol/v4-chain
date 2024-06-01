"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
const helpers_1 = require("../../../helpers");
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.grants = this.grants.bind(this);
        this.granterGrants = this.granterGrants.bind(this);
        this.granteeGrants = this.granteeGrants.bind(this);
    }
    /* Returns list of `Authorization`, granted to the grantee by the granter. */
    async grants(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.granter) !== "undefined") {
            options.params.granter = params.granter;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.grantee) !== "undefined") {
            options.params.grantee = params.grantee;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.msgTypeUrl) !== "undefined") {
            options.params.msg_type_url = params.msgTypeUrl;
        }
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/authz/v1beta1/grants`;
        return await this.req.get(endpoint, options);
    }
    /* GranterGrants returns list of `GrantAuthorization`, granted by granter.
    
     Since: cosmos-sdk 0.46 */
    async granterGrants(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/authz/v1beta1/grants/granter/${params.granter}`;
        return await this.req.get(endpoint, options);
    }
    /* GranteeGrants returns a list of `GrantAuthorization` by grantee.
    
     Since: cosmos-sdk 0.46 */
    async granteeGrants(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/authz/v1beta1/grants/grantee/${params.grantee}`;
        return await this.req.get(endpoint, options);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL2F1dGh6L3YxYmV0YTEvcXVlcnkubGNkLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7OztBQUFBLDhDQUF1RDtBQUd2RCxNQUFhLGNBQWM7SUFHekIsWUFBWSxFQUNWLGFBQWEsRUFHZDtRQUNDLElBQUksQ0FBQyxHQUFHLEdBQUcsYUFBYSxDQUFDO1FBQ3pCLElBQUksQ0FBQyxNQUFNLEdBQUcsSUFBSSxDQUFDLE1BQU0sQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDckMsSUFBSSxDQUFDLGFBQWEsR0FBRyxJQUFJLENBQUMsYUFBYSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNuRCxJQUFJLENBQUMsYUFBYSxHQUFHLElBQUksQ0FBQyxhQUFhLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO0lBQ3JELENBQUM7SUFDRCw2RUFBNkU7SUFHN0UsS0FBSyxDQUFDLE1BQU0sQ0FBQyxNQUEwQjtRQUNyQyxNQUFNLE9BQU8sR0FBUTtZQUNuQixNQUFNLEVBQUUsRUFBRTtTQUNYLENBQUM7UUFFRixJQUFJLE9BQU8sQ0FBQSxNQUFNLGFBQU4sTUFBTSx1QkFBTixNQUFNLENBQUUsT0FBTyxDQUFBLEtBQUssV0FBVyxFQUFFO1lBQzFDLE9BQU8sQ0FBQyxNQUFNLENBQUMsT0FBTyxHQUFHLE1BQU0sQ0FBQyxPQUFPLENBQUM7U0FDekM7UUFFRCxJQUFJLE9BQU8sQ0FBQSxNQUFNLGFBQU4sTUFBTSx1QkFBTixNQUFNLENBQUUsT0FBTyxDQUFBLEtBQUssV0FBVyxFQUFFO1lBQzFDLE9BQU8sQ0FBQyxNQUFNLENBQUMsT0FBTyxHQUFHLE1BQU0sQ0FBQyxPQUFPLENBQUM7U0FDekM7UUFFRCxJQUFJLE9BQU8sQ0FBQSxNQUFNLGFBQU4sTUFBTSx1QkFBTixNQUFNLENBQUUsVUFBVSxDQUFBLEtBQUssV0FBVyxFQUFFO1lBQzdDLE9BQU8sQ0FBQyxNQUFNLENBQUMsWUFBWSxHQUFHLE1BQU0sQ0FBQyxVQUFVLENBQUM7U0FDakQ7UUFFRCxJQUFJLE9BQU8sQ0FBQSxNQUFNLGFBQU4sTUFBTSx1QkFBTixNQUFNLENBQUUsVUFBVSxDQUFBLEtBQUssV0FBVyxFQUFFO1lBQzdDLElBQUEsNkJBQW1CLEVBQUMsT0FBTyxFQUFFLE1BQU0sQ0FBQyxVQUFVLENBQUMsQ0FBQztTQUNqRDtRQUVELE1BQU0sUUFBUSxHQUFHLDZCQUE2QixDQUFDO1FBQy9DLE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBNkIsUUFBUSxFQUFFLE9BQU8sQ0FBQyxDQUFDO0lBQzNFLENBQUM7SUFDRDs7OEJBRTBCO0lBRzFCLEtBQUssQ0FBQyxhQUFhLENBQUMsTUFBaUM7UUFDbkQsTUFBTSxPQUFPLEdBQVE7WUFDbkIsTUFBTSxFQUFFLEVBQUU7U0FDWCxDQUFDO1FBRUYsSUFBSSxPQUFPLENBQUEsTUFBTSxhQUFOLE1BQU0sdUJBQU4sTUFBTSxDQUFFLFVBQVUsQ0FBQSxLQUFLLFdBQVcsRUFBRTtZQUM3QyxJQUFBLDZCQUFtQixFQUFDLE9BQU8sRUFBRSxNQUFNLENBQUMsVUFBVSxDQUFDLENBQUM7U0FDakQ7UUFFRCxNQUFNLFFBQVEsR0FBRyx1Q0FBdUMsTUFBTSxDQUFDLE9BQU8sRUFBRSxDQUFDO1FBQ3pFLE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBb0MsUUFBUSxFQUFFLE9BQU8sQ0FBQyxDQUFDO0lBQ2xGLENBQUM7SUFDRDs7OEJBRTBCO0lBRzFCLEtBQUssQ0FBQyxhQUFhLENBQUMsTUFBaUM7UUFDbkQsTUFBTSxPQUFPLEdBQVE7WUFDbkIsTUFBTSxFQUFFLEVBQUU7U0FDWCxDQUFDO1FBRUYsSUFBSSxPQUFPLENBQUEsTUFBTSxhQUFOLE1BQU0sdUJBQU4sTUFBTSxDQUFFLFVBQVUsQ0FBQSxLQUFLLFdBQVcsRUFBRTtZQUM3QyxJQUFBLDZCQUFtQixFQUFDLE9BQU8sRUFBRSxNQUFNLENBQUMsVUFBVSxDQUFDLENBQUM7U0FDakQ7UUFFRCxNQUFNLFFBQVEsR0FBRyx1Q0FBdUMsTUFBTSxDQUFDLE9BQU8sRUFBRSxDQUFDO1FBQ3pFLE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBb0MsUUFBUSxFQUFFLE9BQU8sQ0FBQyxDQUFDO0lBQ2xGLENBQUM7Q0FFRjtBQTNFRCx3Q0EyRUMifQ==