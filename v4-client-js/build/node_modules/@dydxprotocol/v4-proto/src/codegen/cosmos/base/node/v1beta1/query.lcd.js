"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.config = this.config.bind(this);
        this.status = this.status.bind(this);
    }
    /* Config queries for the operator configuration. */
    async config(_params = {}) {
        const endpoint = `cosmos/base/node/v1beta1/config`;
        return await this.req.get(endpoint);
    }
    /* Status queries for the node status. */
    async status(_params = {}) {
        const endpoint = `cosmos/base/node/v1beta1/status`;
        return await this.req.get(endpoint);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL2Jhc2Uvbm9kZS92MWJldGExL3F1ZXJ5LmxjZC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFFQSxNQUFhLGNBQWM7SUFHekIsWUFBWSxFQUNWLGFBQWEsRUFHZDtRQUNDLElBQUksQ0FBQyxHQUFHLEdBQUcsYUFBYSxDQUFDO1FBQ3pCLElBQUksQ0FBQyxNQUFNLEdBQUcsSUFBSSxDQUFDLE1BQU0sQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUM7UUFDckMsSUFBSSxDQUFDLE1BQU0sR0FBRyxJQUFJLENBQUMsTUFBTSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztJQUN2QyxDQUFDO0lBQ0Qsb0RBQW9EO0lBR3BELEtBQUssQ0FBQyxNQUFNLENBQUMsVUFBeUIsRUFBRTtRQUN0QyxNQUFNLFFBQVEsR0FBRyxpQ0FBaUMsQ0FBQztRQUNuRCxPQUFPLE1BQU0sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQXdCLFFBQVEsQ0FBQyxDQUFDO0lBQzdELENBQUM7SUFDRCx5Q0FBeUM7SUFHekMsS0FBSyxDQUFDLE1BQU0sQ0FBQyxVQUF5QixFQUFFO1FBQ3RDLE1BQU0sUUFBUSxHQUFHLGlDQUFpQyxDQUFDO1FBQ25ELE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBd0IsUUFBUSxDQUFDLENBQUM7SUFDN0QsQ0FBQztDQUVGO0FBM0JELHdDQTJCQyJ9