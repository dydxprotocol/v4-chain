"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const rest_1 = __importDefault(require("./rest"));
class UtilityClient extends rest_1.default {
    /**
     * @description Get the current time of the Indexer
     * @returns {TimeResponse} isoString and epoch
     */
    async getTime() {
        const uri = '/v4/time';
        return this.get(uri);
    }
    /**
     * @description Get the block height of the most recent block processed by the Indexer
     * @returns {HeightResponse} block height and time
     */
    async getHeight() {
        const uri = '/v4/height';
        return this.get(uri);
    }
    /**
     * @description Screen an address to see if it is restricted
     * @param {string} address evm or dydx address
     * @returns {ComplianceResponse} whether the specified address is restricted
     */
    async screen(address) {
        const uri = '/v4/screen';
        return this.get(uri, { address });
    }
}
exports.default = UtilityClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidXRpbGl0eS5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3NyYy9jbGllbnRzL21vZHVsZXMvdXRpbGl0eS50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7OztBQUNBLGtEQUFnQztBQUVoQyxNQUFxQixhQUFjLFNBQVEsY0FBVTtJQUNuRDs7O09BR0c7SUFDSCxLQUFLLENBQUMsT0FBTztRQUNYLE1BQU0sR0FBRyxHQUFHLFVBQVUsQ0FBQztRQUN2QixPQUFPLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFDLENBQUM7SUFDdkIsQ0FBQztJQUVEOzs7T0FHRztJQUNILEtBQUssQ0FBQyxTQUFTO1FBQ2IsTUFBTSxHQUFHLEdBQUcsWUFBWSxDQUFDO1FBQ3pCLE9BQU8sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQUMsQ0FBQztJQUN2QixDQUFDO0lBRUQ7Ozs7T0FJRztJQUNILEtBQUssQ0FBQyxNQUFNLENBQUMsT0FBZTtRQUMxQixNQUFNLEdBQUcsR0FBRyxZQUFZLENBQUM7UUFDekIsT0FBTyxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsRUFBRSxFQUFFLE9BQU8sRUFBRSxDQUFDLENBQUM7SUFDcEMsQ0FBQztDQUNGO0FBNUJELGdDQTRCQyJ9