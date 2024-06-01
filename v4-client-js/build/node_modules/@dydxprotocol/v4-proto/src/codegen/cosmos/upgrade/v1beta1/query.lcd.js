"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.currentPlan = this.currentPlan.bind(this);
        this.appliedPlan = this.appliedPlan.bind(this);
        this.upgradedConsensusState = this.upgradedConsensusState.bind(this);
        this.moduleVersions = this.moduleVersions.bind(this);
        this.authority = this.authority.bind(this);
    }
    /* CurrentPlan queries the current upgrade plan. */
    async currentPlan(_params = {}) {
        const endpoint = `cosmos/upgrade/v1beta1/current_plan`;
        return await this.req.get(endpoint);
    }
    /* AppliedPlan queries a previously applied upgrade plan by its name. */
    async appliedPlan(params) {
        const endpoint = `cosmos/upgrade/v1beta1/applied_plan/${params.name}`;
        return await this.req.get(endpoint);
    }
    /* UpgradedConsensusState queries the consensus state that will serve
     as a trusted kernel for the next version of this chain. It will only be
     stored at the last height of this chain.
     UpgradedConsensusState RPC not supported with legacy querier
     This rpc is deprecated now that IBC has its own replacement
     (https://github.com/cosmos/ibc-go/blob/2c880a22e9f9cc75f62b527ca94aa75ce1106001/proto/ibc/core/client/v1/query.proto#L54) */
    async upgradedConsensusState(params) {
        const endpoint = `cosmos/upgrade/v1beta1/upgraded_consensus_state/${params.lastHeight}`;
        return await this.req.get(endpoint);
    }
    /* ModuleVersions queries the list of module versions from state.
    
     Since: cosmos-sdk 0.43 */
    async moduleVersions(params) {
        const options = {
            params: {}
        };
        if (typeof (params === null || params === void 0 ? void 0 : params.moduleName) !== "undefined") {
            options.params.module_name = params.moduleName;
        }
        const endpoint = `cosmos/upgrade/v1beta1/module_versions`;
        return await this.req.get(endpoint, options);
    }
    /* Returns the account with authority to conduct upgrades
    
     Since: cosmos-sdk 0.46 */
    async authority(_params = {}) {
        const endpoint = `cosmos/upgrade/v1beta1/authority`;
        return await this.req.get(endpoint);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL3VwZ3JhZGUvdjFiZXRhMS9xdWVyeS5sY2QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBRUEsTUFBYSxjQUFjO0lBR3pCLFlBQVksRUFDVixhQUFhLEVBR2Q7UUFDQyxJQUFJLENBQUMsR0FBRyxHQUFHLGFBQWEsQ0FBQztRQUN6QixJQUFJLENBQUMsV0FBVyxHQUFHLElBQUksQ0FBQyxXQUFXLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQy9DLElBQUksQ0FBQyxXQUFXLEdBQUcsSUFBSSxDQUFDLFdBQVcsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDL0MsSUFBSSxDQUFDLHNCQUFzQixHQUFHLElBQUksQ0FBQyxzQkFBc0IsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDckUsSUFBSSxDQUFDLGNBQWMsR0FBRyxJQUFJLENBQUMsY0FBYyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztRQUNyRCxJQUFJLENBQUMsU0FBUyxHQUFHLElBQUksQ0FBQyxTQUFTLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO0lBQzdDLENBQUM7SUFDRCxtREFBbUQ7SUFHbkQsS0FBSyxDQUFDLFdBQVcsQ0FBQyxVQUFtQyxFQUFFO1FBQ3JELE1BQU0sUUFBUSxHQUFHLHFDQUFxQyxDQUFDO1FBQ3ZELE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBa0MsUUFBUSxDQUFDLENBQUM7SUFDdkUsQ0FBQztJQUNELHdFQUF3RTtJQUd4RSxLQUFLLENBQUMsV0FBVyxDQUFDLE1BQStCO1FBQy9DLE1BQU0sUUFBUSxHQUFHLHVDQUF1QyxNQUFNLENBQUMsSUFBSSxFQUFFLENBQUM7UUFDdEUsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFrQyxRQUFRLENBQUMsQ0FBQztJQUN2RSxDQUFDO0lBQ0Q7Ozs7O2lJQUs2SDtJQUc3SCxLQUFLLENBQUMsc0JBQXNCLENBQUMsTUFBMEM7UUFDckUsTUFBTSxRQUFRLEdBQUcsbURBQW1ELE1BQU0sQ0FBQyxVQUFVLEVBQUUsQ0FBQztRQUN4RixPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQTZDLFFBQVEsQ0FBQyxDQUFDO0lBQ2xGLENBQUM7SUFDRDs7OEJBRTBCO0lBRzFCLEtBQUssQ0FBQyxjQUFjLENBQUMsTUFBa0M7UUFDckQsTUFBTSxPQUFPLEdBQVE7WUFDbkIsTUFBTSxFQUFFLEVBQUU7U0FDWCxDQUFDO1FBRUYsSUFBSSxPQUFPLENBQUEsTUFBTSxhQUFOLE1BQU0sdUJBQU4sTUFBTSxDQUFFLFVBQVUsQ0FBQSxLQUFLLFdBQVcsRUFBRTtZQUM3QyxPQUFPLENBQUMsTUFBTSxDQUFDLFdBQVcsR0FBRyxNQUFNLENBQUMsVUFBVSxDQUFDO1NBQ2hEO1FBRUQsTUFBTSxRQUFRLEdBQUcsd0NBQXdDLENBQUM7UUFDMUQsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFxQyxRQUFRLEVBQUUsT0FBTyxDQUFDLENBQUM7SUFDbkYsQ0FBQztJQUNEOzs4QkFFMEI7SUFHMUIsS0FBSyxDQUFDLFNBQVMsQ0FBQyxVQUFpQyxFQUFFO1FBQ2pELE1BQU0sUUFBUSxHQUFHLGtDQUFrQyxDQUFDO1FBQ3BELE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBZ0MsUUFBUSxDQUFDLENBQUM7SUFDckUsQ0FBQztDQUVGO0FBcEVELHdDQW9FQyJ9