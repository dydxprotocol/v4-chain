"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const constants_1 = require("../constants");
const rest_1 = __importDefault(require("./rest"));
/**
 * @description REST endpoints for data unrelated to a particular address.
 */
class MarketsClient extends rest_1.default {
    async getPerpetualMarkets(market) {
        const uri = '/v4/perpetualMarkets';
        return this.get(uri, { ticker: market });
    }
    async getPerpetualMarketOrderbook(market) {
        const uri = `/v4/orderbooks/perpetualMarket/${market}`;
        return this.get(uri);
    }
    async getPerpetualMarketTrades(market, startingBeforeOrAtHeight, limit) {
        const uri = `/v4/trades/perpetualMarket/${market}`;
        return this.get(uri, {
            createdBeforeOrAtHeight: startingBeforeOrAtHeight,
            limit,
        });
    }
    async getPerpetualMarketCandles(market, resolution, fromISO, toISO, limit) {
        const uri = `/v4/candles/perpetualMarkets/${market}`;
        return this.get(uri, {
            resolution,
            fromISO,
            toISO,
            limit,
        });
    }
    async getPerpetualMarketHistoricalFunding(market, effectiveBeforeOrAt, effectiveBeforeOrAtHeight, limit) {
        const uri = `/v4/historicalFunding/${market}`;
        return this.get(uri, {
            effectiveBeforeOrAt,
            effectiveBeforeOrAtHeight,
            limit,
        });
    }
    async getPerpetualMarketSparklines(period = constants_1.TimePeriod.ONE_DAY) {
        const uri = '/v4/sparklines';
        return this.get(uri, {
            timePeriod: period,
        });
    }
}
exports.default = MarketsClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibWFya2V0cy5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3NyYy9jbGllbnRzL21vZHVsZXMvbWFya2V0cy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7OztBQUFBLDRDQUEwQztBQUUxQyxrREFBZ0M7QUFFaEM7O0dBRUc7QUFDSCxNQUFxQixhQUFjLFNBQVEsY0FBVTtJQUNuRCxLQUFLLENBQUMsbUJBQW1CLENBQUMsTUFBZTtRQUN2QyxNQUFNLEdBQUcsR0FBRyxzQkFBc0IsQ0FBQztRQUNuQyxPQUFPLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxFQUFFLEVBQUUsTUFBTSxFQUFFLE1BQU0sRUFBRSxDQUFDLENBQUM7SUFDM0MsQ0FBQztJQUVELEtBQUssQ0FBQywyQkFBMkIsQ0FBQyxNQUFjO1FBQzlDLE1BQU0sR0FBRyxHQUFHLGtDQUFrQyxNQUFNLEVBQUUsQ0FBQztRQUN2RCxPQUFPLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFDLENBQUM7SUFDdkIsQ0FBQztJQUVELEtBQUssQ0FBQyx3QkFBd0IsQ0FDNUIsTUFBYyxFQUNkLHdCQUF3QyxFQUN4QyxLQUFxQjtRQUVyQixNQUFNLEdBQUcsR0FBRyw4QkFBOEIsTUFBTSxFQUFFLENBQUM7UUFDbkQsT0FBTyxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsRUFBRTtZQUNuQix1QkFBdUIsRUFBRSx3QkFBd0I7WUFDakQsS0FBSztTQUNOLENBQUMsQ0FBQztJQUNMLENBQUM7SUFFRCxLQUFLLENBQUMseUJBQXlCLENBQzdCLE1BQWMsRUFDZCxVQUFrQixFQUNsQixPQUF1QixFQUN2QixLQUFxQixFQUNyQixLQUFxQjtRQUVyQixNQUFNLEdBQUcsR0FBRyxnQ0FBZ0MsTUFBTSxFQUFFLENBQUM7UUFDckQsT0FBTyxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsRUFBRTtZQUNuQixVQUFVO1lBQ1YsT0FBTztZQUNQLEtBQUs7WUFDTCxLQUFLO1NBQ04sQ0FBQyxDQUFDO0lBQ0wsQ0FBQztJQUVELEtBQUssQ0FBQyxtQ0FBbUMsQ0FDdkMsTUFBYyxFQUNkLG1CQUFtQyxFQUNuQyx5QkFBeUMsRUFDekMsS0FBcUI7UUFFckIsTUFBTSxHQUFHLEdBQUcseUJBQXlCLE1BQU0sRUFBRSxDQUFDO1FBQzlDLE9BQU8sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLEVBQUU7WUFDbkIsbUJBQW1CO1lBQ25CLHlCQUF5QjtZQUN6QixLQUFLO1NBQ04sQ0FBQyxDQUFDO0lBQ0wsQ0FBQztJQUVELEtBQUssQ0FBQyw0QkFBNEIsQ0FBQyxTQUFpQixzQkFBVSxDQUFDLE9BQU87UUFDcEUsTUFBTSxHQUFHLEdBQUcsZ0JBQWdCLENBQUM7UUFDN0IsT0FBTyxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsRUFBRTtZQUNuQixVQUFVLEVBQUUsTUFBTTtTQUNuQixDQUFDLENBQUM7SUFDTCxDQUFDO0NBQ0Y7QUEzREQsZ0NBMkRDIn0=