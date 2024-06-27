"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateRegistry = exports.registry = void 0;
const proto_signing_1 = require("@cosmjs/proto-signing");
const stargate_1 = require("@cosmjs/stargate");
const tx_1 = require("@dydxprotocol/v4-proto/src/codegen/dydxprotocol/clob/tx");
const transfer_1 = require("@dydxprotocol/v4-proto/src/codegen/dydxprotocol/sending/transfer");
const tx_2 = require("@dydxprotocol/v4-proto/src/codegen/dydxprotocol/sending/tx");
exports.registry = [];
function generateRegistry() {
    return new proto_signing_1.Registry([
        // clob
        ['/dydxprotocol.clob.MsgPlaceOrder', tx_1.MsgPlaceOrder],
        ['/dydxprotocol.clob.MsgCancelOrder', tx_1.MsgCancelOrder],
        // sending
        ['/dydxprotocol.sending.MsgCreateTransfer', tx_2.MsgCreateTransfer],
        ['/dydxprotocol.sending.MsgWithdrawFromSubaccount', transfer_1.MsgWithdrawFromSubaccount],
        ['/dydxprotocol.sending.MsgDepositToSubaccount', transfer_1.MsgDepositToSubaccount],
        // default types
        ...stargate_1.defaultRegistryTypes,
    ]);
}
exports.generateRegistry = generateRegistry;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicmVnaXN0cnkuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi9zcmMvY2xpZW50cy9saWIvcmVnaXN0cnkudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBQUEseURBQWdFO0FBQ2hFLCtDQUF3RDtBQUN4RCxnRkFHaUU7QUFDakUsK0ZBRzBFO0FBQzFFLG1GQUVvRTtBQUV2RCxRQUFBLFFBQVEsR0FBMkMsRUFBRSxDQUFDO0FBQ25FLFNBQWdCLGdCQUFnQjtJQUM5QixPQUFPLElBQUksd0JBQVEsQ0FBQztRQUNsQixPQUFPO1FBQ1AsQ0FBQyxrQ0FBa0MsRUFBRSxrQkFBOEIsQ0FBQztRQUNwRSxDQUFDLG1DQUFtQyxFQUFFLG1CQUErQixDQUFDO1FBRXRFLFVBQVU7UUFDVixDQUFDLHlDQUF5QyxFQUFFLHNCQUFrQyxDQUFDO1FBQy9FLENBQUMsaURBQWlELEVBQUUsb0NBQTBDLENBQUM7UUFDL0YsQ0FBQyw4Q0FBOEMsRUFBRSxpQ0FBdUMsQ0FBQztRQUV6RixnQkFBZ0I7UUFDaEIsR0FBRywrQkFBb0I7S0FDeEIsQ0FBQyxDQUFDO0FBQ0wsQ0FBQztBQWRELDRDQWNDIn0=