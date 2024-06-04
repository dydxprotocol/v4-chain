"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const constants_1 = require("../constants");
const rest_1 = __importDefault(require("./rest"));
/**
 * @description REST endpoints for data related to a particular address.
 */
class AccountClient extends rest_1.default {
    async getSubaccounts(address, limit) {
        const uri = `/v4/addresses/${address}`;
        return this.get(uri, { limit });
    }
    async getSubaccount(address, subaccountNumber) {
        const uri = `/v4/addresses/${address}/subaccountNumber/${subaccountNumber}`;
        return this.get(uri);
    }
    async getSubaccountPerpetualPositions(address, subaccountNumber, status, limit, createdBeforeOrAtHeight, createdBeforeOrAt) {
        const uri = '/v4/perpetualPositions';
        return this.get(uri, {
            address,
            subaccountNumber,
            status,
            limit,
            createdBeforeOrAtHeight,
            createdBeforeOrAt,
        });
    }
    async getSubaccountAssetPositions(address, subaccountNumber, status, limit, createdBeforeOrAtHeight, createdBeforeOrAt) {
        const uri = '/v4/assetPositions';
        return this.get(uri, {
            address,
            subaccountNumber,
            status,
            limit,
            createdBeforeOrAtHeight,
            createdBeforeOrAt,
        });
    }
    async getSubaccountTransfers(address, subaccountNumber, limit, createdBeforeOrAtHeight, createdBeforeOrAt) {
        const uri = '/v4/transfers';
        return this.get(uri, {
            address,
            subaccountNumber,
            limit,
            createdBeforeOrAtHeight,
            createdBeforeOrAt,
        });
    }
    async getSubaccountOrders(address, subaccountNumber, ticker, tickerType = constants_1.TickerType.PERPETUAL, side, status, type, limit, goodTilBlockBeforeOrAt, goodTilBlockTimeBeforeOrAt, returnLatestOrders) {
        const uri = '/v4/orders';
        return this.get(uri, {
            address,
            subaccountNumber,
            ticker,
            tickerType,
            side,
            status,
            type,
            limit,
            goodTilBlockBeforeOrAt,
            goodTilBlockTimeBeforeOrAt,
            returnLatestOrders,
        });
    }
    async getOrder(orderId) {
        const uri = `/v4/orders${orderId}`;
        return this.get(uri);
    }
    async getSubaccountFills(address, subaccountNumber, ticker, tickerType = constants_1.TickerType.PERPETUAL, limit, createdBeforeOrAtHeight, createdBeforeOrAt) {
        const uri = '/v4/fills';
        return this.get(uri, {
            address,
            subaccountNumber,
            ticker,
            tickerType,
            limit,
            createdBeforeOrAtHeight,
            createdBeforeOrAt,
        });
    }
    async getSubaccountHistoricalPNLs(address, subaccountNumber, effectiveBeforeOrAt, effectiveAtOrAfter) {
        const uri = '/v4/historical-pnl';
        return this.get(uri, {
            address,
            subaccountNumber,
            effectiveBeforeOrAt,
            effectiveAtOrAfter,
        });
    }
}
exports.default = AccountClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYWNjb3VudC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3NyYy9jbGllbnRzL21vZHVsZXMvYWNjb3VudC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7OztBQUFBLDRDQUVzQjtBQUV0QixrREFBZ0M7QUFFaEM7O0dBRUc7QUFDSCxNQUFxQixhQUFjLFNBQVEsY0FBVTtJQUNuRCxLQUFLLENBQUMsY0FBYyxDQUNsQixPQUFlLEVBQ2YsS0FBYztRQUVkLE1BQU0sR0FBRyxHQUFHLGlCQUFpQixPQUFPLEVBQUUsQ0FBQztRQUN2QyxPQUFPLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxFQUFFLEVBQUUsS0FBSyxFQUFFLENBQUMsQ0FBQztJQUNsQyxDQUFDO0lBRUQsS0FBSyxDQUFDLGFBQWEsQ0FDakIsT0FBZSxFQUNmLGdCQUF3QjtRQUV4QixNQUFNLEdBQUcsR0FBRyxpQkFBaUIsT0FBTyxxQkFBcUIsZ0JBQWdCLEVBQUUsQ0FBQztRQUM1RSxPQUFPLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFDLENBQUM7SUFDdkIsQ0FBQztJQUVELEtBQUssQ0FBQywrQkFBK0IsQ0FDbkMsT0FBZSxFQUNmLGdCQUF3QixFQUN4QixNQUE4QixFQUM5QixLQUFxQixFQUNyQix1QkFBdUMsRUFDdkMsaUJBQWlDO1FBRWpDLE1BQU0sR0FBRyxHQUFHLHdCQUF3QixDQUFDO1FBQ3JDLE9BQU8sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLEVBQUU7WUFDbkIsT0FBTztZQUNQLGdCQUFnQjtZQUNoQixNQUFNO1lBQ04sS0FBSztZQUNMLHVCQUF1QjtZQUN2QixpQkFBaUI7U0FDbEIsQ0FBQyxDQUFDO0lBQ0wsQ0FBQztJQUVELEtBQUssQ0FBQywyQkFBMkIsQ0FDL0IsT0FBZSxFQUNmLGdCQUF3QixFQUN4QixNQUE4QixFQUM5QixLQUFxQixFQUNyQix1QkFBdUMsRUFDdkMsaUJBQWlDO1FBRWpDLE1BQU0sR0FBRyxHQUFHLG9CQUFvQixDQUFDO1FBQ2pDLE9BQU8sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLEVBQUU7WUFDbkIsT0FBTztZQUNQLGdCQUFnQjtZQUNoQixNQUFNO1lBQ04sS0FBSztZQUNMLHVCQUF1QjtZQUN2QixpQkFBaUI7U0FDbEIsQ0FBQyxDQUFDO0lBQ0wsQ0FBQztJQUVELEtBQUssQ0FBQyxzQkFBc0IsQ0FDMUIsT0FBZSxFQUNmLGdCQUF3QixFQUN4QixLQUFxQixFQUNyQix1QkFBdUMsRUFDdkMsaUJBQWlDO1FBRWpDLE1BQU0sR0FBRyxHQUFHLGVBQWUsQ0FBQztRQUM1QixPQUFPLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxFQUFFO1lBQ25CLE9BQU87WUFDUCxnQkFBZ0I7WUFDaEIsS0FBSztZQUNMLHVCQUF1QjtZQUN2QixpQkFBaUI7U0FDbEIsQ0FBQyxDQUFDO0lBQ0wsQ0FBQztJQUVELEtBQUssQ0FBQyxtQkFBbUIsQ0FDdkIsT0FBZSxFQUNmLGdCQUF3QixFQUN4QixNQUFzQixFQUN0QixhQUF5QixzQkFBVSxDQUFDLFNBQVMsRUFDN0MsSUFBdUIsRUFDdkIsTUFBMkIsRUFDM0IsSUFBdUIsRUFDdkIsS0FBcUIsRUFDckIsc0JBQXNDLEVBQ3RDLDBCQUEwQyxFQUMxQyxrQkFBbUM7UUFFbkMsTUFBTSxHQUFHLEdBQUcsWUFBWSxDQUFDO1FBQ3pCLE9BQU8sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLEVBQUU7WUFDbkIsT0FBTztZQUNQLGdCQUFnQjtZQUNoQixNQUFNO1lBQ04sVUFBVTtZQUNWLElBQUk7WUFDSixNQUFNO1lBQ04sSUFBSTtZQUNKLEtBQUs7WUFDTCxzQkFBc0I7WUFDdEIsMEJBQTBCO1lBQzFCLGtCQUFrQjtTQUNuQixDQUFDLENBQUM7SUFDTCxDQUFDO0lBRUQsS0FBSyxDQUFDLFFBQVEsQ0FBQyxPQUFlO1FBQzVCLE1BQU0sR0FBRyxHQUFHLGFBQWEsT0FBTyxFQUFFLENBQUM7UUFDbkMsT0FBTyxJQUFJLENBQUMsR0FBRyxDQUFDLEdBQUcsQ0FBQyxDQUFDO0lBQ3ZCLENBQUM7SUFFRCxLQUFLLENBQUMsa0JBQWtCLENBQ3RCLE9BQWUsRUFDZixnQkFBd0IsRUFDeEIsTUFBc0IsRUFDdEIsYUFBeUIsc0JBQVUsQ0FBQyxTQUFTLEVBQzdDLEtBQXFCLEVBQ3JCLHVCQUF1QyxFQUN2QyxpQkFBaUM7UUFFakMsTUFBTSxHQUFHLEdBQUcsV0FBVyxDQUFDO1FBQ3hCLE9BQU8sSUFBSSxDQUFDLEdBQUcsQ0FBQyxHQUFHLEVBQUU7WUFDbkIsT0FBTztZQUNQLGdCQUFnQjtZQUNoQixNQUFNO1lBQ04sVUFBVTtZQUNWLEtBQUs7WUFDTCx1QkFBdUI7WUFDdkIsaUJBQWlCO1NBQ2xCLENBQUMsQ0FBQztJQUNMLENBQUM7SUFFRCxLQUFLLENBQUMsMkJBQTJCLENBQy9CLE9BQWUsRUFDZixnQkFBd0IsRUFDeEIsbUJBQW1DLEVBQ25DLGtCQUFrQztRQUVsQyxNQUFNLEdBQUcsR0FBRyxvQkFBb0IsQ0FBQztRQUNqQyxPQUFPLElBQUksQ0FBQyxHQUFHLENBQUMsR0FBRyxFQUFFO1lBQ25CLE9BQU87WUFDUCxnQkFBZ0I7WUFDaEIsbUJBQW1CO1lBQ25CLGtCQUFrQjtTQUNuQixDQUFDLENBQUM7SUFDTCxDQUFDO0NBQ0Y7QUE3SUQsZ0NBNklDIn0=