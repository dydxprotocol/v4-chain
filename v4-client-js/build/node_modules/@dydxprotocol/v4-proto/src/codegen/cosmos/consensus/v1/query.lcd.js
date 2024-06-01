"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LCDQueryClient = void 0;
class LCDQueryClient {
    constructor({ requestClient }) {
        this.req = requestClient;
        this.params = this.params.bind(this);
    }
    /* Params queries the parameters of x/consensus module. */
    async params(_params = {}) {
        const endpoint = `cosmos/consensus/v1/params`;
        return await this.req.get(endpoint);
    }
}
exports.LCDQueryClient = LCDQueryClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicXVlcnkubGNkLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vY29zbW9zL2NvbnNlbnN1cy92MS9xdWVyeS5sY2QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBRUEsTUFBYSxjQUFjO0lBR3pCLFlBQVksRUFDVixhQUFhLEVBR2Q7UUFDQyxJQUFJLENBQUMsR0FBRyxHQUFHLGFBQWEsQ0FBQztRQUN6QixJQUFJLENBQUMsTUFBTSxHQUFHLElBQUksQ0FBQyxNQUFNLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDO0lBQ3ZDLENBQUM7SUFDRCwwREFBMEQ7SUFHMUQsS0FBSyxDQUFDLE1BQU0sQ0FBQyxVQUE4QixFQUFFO1FBQzNDLE1BQU0sUUFBUSxHQUFHLDRCQUE0QixDQUFDO1FBQzlDLE9BQU8sTUFBTSxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBNkIsUUFBUSxDQUFDLENBQUM7SUFDbEUsQ0FBQztDQUVGO0FBbkJELHdDQW1CQyJ9