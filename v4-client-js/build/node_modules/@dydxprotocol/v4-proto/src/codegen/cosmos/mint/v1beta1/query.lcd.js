"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.params = this.params.bind(this);
        this.inflation = this.inflation.bind(this);
        this.annualProvisions = this.annualProvisions.bind(this);
    }
    /* Params returns the total set of minting parameters. */
    async params(_params = {}) {
        const endpoint = `cosmos/mint/v1beta1/params`;
        return await this.req.get(endpoint);
    }
    /* Inflation returns the current minting inflation value. */
    async inflation(_params = {}) {
        const endpoint = `cosmos/mint/v1beta1/inflation`;
        return await this.req.get(endpoint);
    }
    /* AnnualProvisions current minting annual provisions value. */
    async annualProvisions(_params = {}) {
        const endpoint = `cosmos/mint/v1beta1/annual_provisions`;
        return await this.req.get(endpoint);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL21pbnQvdjFiZXRhMS9xdWVyeS5sY2QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBRUEsTUFBYSxjQUFjO0lBR3pCLFlBQVksRUFDVixhQUFhLEVBR2Q7UUFDQyxJQUFJLENBQUMsR0FBRyxHQUFHLGFBQWEsQ0FBQztRQUN6QixJQUFJLENBQUMsTUFBTSxHQUFHLElBQUksQ0FBQyxNQUFNLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO1FBQ3JDLElBQUksQ0FBQyxTQUFTLEdBQUcsSUFBSSxDQUFDLFNBQVMsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDM0MsSUFBSSxDQUFDLGdCQUFnQixHQUFHLElBQUksQ0FBQyxnQkFBZ0IsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7SUFDM0QsQ0FBQztJQUNELHlEQUF5RDtJQUd6RCxLQUFLLENBQUMsTUFBTSxDQUFDLFVBQThCLEVBQUU7UUFDM0MsTUFBTSxRQUFRLEdBQUcsNEJBQTRCLENBQUM7UUFDOUMsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUE2QixRQUFRLENBQUMsQ0FBQztJQUNsRSxDQUFDO0lBQ0QsNERBQTREO0lBRzVELEtBQUssQ0FBQyxTQUFTLENBQUMsVUFBaUMsRUFBRTtRQUNqRCxNQUFNLFFBQVEsR0FBRywrQkFBK0IsQ0FBQztRQUNqRCxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQWdDLFFBQVEsQ0FBQyxDQUFDO0lBQ3JFLENBQUM7SUFDRCwrREFBK0Q7SUFHL0QsS0FBSyxDQUFDLGdCQUFnQixDQUFDLFVBQXdDLEVBQUU7UUFDL0QsTUFBTSxRQUFRLEdBQUcsdUNBQXVDLENBQUM7UUFDekQsT0FBTyxNQUFNLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUF1QyxRQUFRLENBQUMsQ0FBQztJQUM1RSxDQUFDO0NBRUY7QUFuQ0Qsd0NBbUNDIn0=