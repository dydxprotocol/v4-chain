"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.defaultOrder = exports.MAX_CLIENT_ID = exports.PERPETUAL_PAIR_BTC_USD = exports.MARKET_BTC_USD = exports.DYDX_LOCAL_MNEMONIC = exports.DYDX_LOCAL_ADDRESS = exports.DYDX_TEST_MNEMONIC = exports.DYDX_TEST_PRIVATE_KEY = exports.DYDX_TEST_ADDRESS = void 0;
const order_1 = require("@dydxprotocol/v4-proto/src/codegen/dydxprotocol/clob/order");
const long_1 = __importDefault(require("long"));
const types_1 = require("../src/clients/types");
exports.DYDX_TEST_ADDRESS = 'dydx14zzueazeh0hj67cghhf9jypslcf9sh2n5k6art';
exports.DYDX_TEST_PRIVATE_KEY = 'e92a6595c934c991d3b3e987ea9b3125bf61a076deab3a9cb519787b7b3e8d77';
exports.DYDX_TEST_MNEMONIC = 'mirror actor skill push coach wait confirm orchard lunch mobile athlete gossip awake miracle matter bus reopen team ladder lazy list timber render wait';
exports.DYDX_LOCAL_ADDRESS = 'dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4';
exports.DYDX_LOCAL_MNEMONIC = 'merge panther lobster crazy road hollow amused security before critic about cliff exhibit cause coyote talent happy where lion river tobacco option coconut small';
exports.MARKET_BTC_USD = 'BTC-USD';
exports.PERPETUAL_PAIR_BTC_USD = 0;
const quantums = new long_1.default(1000000000);
const subticks = new long_1.default(1000000000);
exports.MAX_CLIENT_ID = 2 ** 32 - 1;
// PlaceOrder variables
exports.defaultOrder = {
    clientId: 0,
    orderFlags: types_1.OrderFlags.SHORT_TERM,
    clobPairId: exports.PERPETUAL_PAIR_BTC_USD,
    side: order_1.Order_Side.SIDE_BUY,
    quantums,
    subticks,
    timeInForce: order_1.Order_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
    reduceOnly: false,
    clientMetadata: 0,
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY29uc3RhbnRzLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vZXhhbXBsZXMvY29uc3RhbnRzLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7OztBQUFBLHNGQUEyRztBQUMzRyxnREFBd0I7QUFFeEIsZ0RBQStEO0FBRWxELFFBQUEsaUJBQWlCLEdBQUcsNkNBQTZDLENBQUM7QUFDbEUsUUFBQSxxQkFBcUIsR0FBRyxrRUFBa0UsQ0FBQztBQUMzRixRQUFBLGtCQUFrQixHQUFHLHlKQUF5SixDQUFDO0FBQy9LLFFBQUEsa0JBQWtCLEdBQUcsNkNBQTZDLENBQUM7QUFDbkUsUUFBQSxtQkFBbUIsR0FBRyxtS0FBbUssQ0FBQztBQUUxTCxRQUFBLGNBQWMsR0FBVyxTQUFTLENBQUM7QUFDbkMsUUFBQSxzQkFBc0IsR0FBVyxDQUFDLENBQUM7QUFFaEQsTUFBTSxRQUFRLEdBQVMsSUFBSSxjQUFJLENBQUMsVUFBYSxDQUFDLENBQUM7QUFDL0MsTUFBTSxRQUFRLEdBQVMsSUFBSSxjQUFJLENBQUMsVUFBYSxDQUFDLENBQUM7QUFFbEMsUUFBQSxhQUFhLEdBQUcsQ0FBQyxJQUFJLEVBQUUsR0FBRyxDQUFDLENBQUM7QUFFekMsdUJBQXVCO0FBQ1YsUUFBQSxZQUFZLEdBQWdCO0lBQ3ZDLFFBQVEsRUFBRSxDQUFDO0lBQ1gsVUFBVSxFQUFFLGtCQUFVLENBQUMsVUFBVTtJQUNqQyxVQUFVLEVBQUUsOEJBQXNCO0lBQ2xDLElBQUksRUFBRSxrQkFBVSxDQUFDLFFBQVE7SUFDekIsUUFBUTtJQUNSLFFBQVE7SUFDUixXQUFXLEVBQUUseUJBQWlCLENBQUMseUJBQXlCO0lBQ3hELFVBQVUsRUFBRSxLQUFLO0lBQ2pCLGNBQWMsRUFBRSxDQUFDO0NBQ2xCLENBQUMifQ==