"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
const helpers_1 = require("../../../helpers");
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.evidence = this.evidence.bind(this);
        this.allEvidence = this.allEvidence.bind(this);
    }
    /* Evidence queries evidence based on evidence hash. */
    async evidence(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.evidenceHash) !== "undefined") {
            options.params.evidence_hash = params.evidenceHash;
        }
        const endpoint = `cosmos/evidence/v1beta1/evidence/${params.hash}`;
        return await this.req.get(endpoint, options);
    }
    /* AllEvidence queries all evidence. */
    async allEvidence(params = {
        pagination: undefined
    }) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.pagination) !== "undefined") {
            (0, helpers_1.setPaginationParams)(options, params.pagination);
        }
        const endpoint = `cosmos/evidence/v1beta1/evidence`;
        return await this.req.get(endpoint, options);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL2V2aWRlbmNlL3YxYmV0YTEvcXVlcnkubGNkLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7OztBQUFBLDhDQUF1RDtBQUd2RCxNQUFhLGNBQWM7SUFHekIsWUFBWSxFQUNWLGFBQWEsRUFHZDtRQUNDLElBQUksQ0FBQyxHQUFHLEdBQUcsYUFBYSxDQUFDO1FBQ3pCLElBQUksQ0FBQyxRQUFRLEdBQUcsSUFBSSxDQUFDLFFBQVEsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDekMsSUFBSSxDQUFDLFdBQVcsR0FBRyxJQUFJLENBQUMsV0FBVyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztJQUNqRCxDQUFDO0lBQ0QsdURBQXVEO0lBR3ZELEtBQUssQ0FBQyxRQUFRLENBQUMsTUFBNEI7UUFDekMsTUFBTSxPQUFPLEdBQVE7WUFDbkIsTUFBTSxFQUFFLEVBQUU7U0FDWCxDQUFDO1FBRUYsSUFBSSxPQUFPLENBQUEsTUFBTSxhQUFOLE1BQU0sdUJBQU4sTUFBTSxDQUFFLFlBQVksQ0FBQSxLQUFLLFdBQVcsRUFBRTtZQUMvQyxPQUFPLENBQUMsTUFBTSxDQUFDLGFBQWEsR0FBRyxNQUFNLENBQUMsWUFBWSxDQUFDO1NBQ3BEO1FBRUQsTUFBTSxRQUFRLEdBQUcsb0NBQW9DLE1BQU0sQ0FBQyxJQUFJLEVBQUUsQ0FBQztRQUNuRSxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQStCLFFBQVEsRUFBRSxPQUFPLENBQUMsQ0FBQztJQUM3RSxDQUFDO0lBQ0QsdUNBQXVDO0lBR3ZDLEtBQUssQ0FBQyxXQUFXLENBQUMsU0FBa0M7UUFDbEQsVUFBVSxFQUFFLFNBQVM7S0FDdEI7UUFDQyxNQUFNLE9BQU8sR0FBUTtZQUNuQixNQUFNLEVBQUUsRUFBRTtTQUNYLENBQUM7UUFFRixJQUFJLE9BQU8sQ0FBQSxNQUFNLGFBQU4sTUFBTSx1QkFBTixNQUFNLENBQUUsVUFBVSxDQUFBLEtBQUssV0FBVyxFQUFFO1lBQzdDLElBQUEsNkJBQW1CLEVBQUMsT0FBTyxFQUFFLE1BQU0sQ0FBQyxVQUFVLENBQUMsQ0FBQztTQUNqRDtRQUVELE1BQU0sUUFBUSxHQUFHLGtDQUFrQyxDQUFDO1FBQ3BELE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBa0MsUUFBUSxFQUFFLE9BQU8sQ0FBQyxDQUFDO0lBQ2hGLENBQUM7Q0FFRjtBQTdDRCx3Q0E2Q0MifQ==