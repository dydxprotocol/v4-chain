"use strict";
/**
 * Simple JS example demostrating filling subaccount with Faucet API
 */
Object.defineProperty(exports, "__esModule", { value: true });
const constants_1 = require("../src/clients/constants");
const faucet_client_1 = require("../src/clients/faucet-client");
const constants_2 = require("./constants");
async function test() {
    const client = new faucet_client_1.FaucetClient(constants_1.FaucetApiHost.TESTNET);
    const address = constants_2.DYDX_TEST_ADDRESS;
    // Use faucet to fill subaccount
    const faucetResponse = await (client === null || client === void 0 ? void 0 : client.fill(address, 0, 2000));
    console.log(faucetResponse);
    const status = faucetResponse.status;
    console.log(status);
}
test().then(() => {
}).catch((error) => {
    console.log(error.message);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZmF1Y2V0X2VuZHBvaW50LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vZXhhbXBsZXMvZmF1Y2V0X2VuZHBvaW50LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7QUFBQTs7R0FFRzs7QUFFSCx3REFBeUQ7QUFDekQsZ0VBQTREO0FBQzVELDJDQUFnRDtBQUVoRCxLQUFLLFVBQVUsSUFBSTtJQUNqQixNQUFNLE1BQU0sR0FBRyxJQUFJLDRCQUFZLENBQUMseUJBQWEsQ0FBQyxPQUFPLENBQUMsQ0FBQztJQUN2RCxNQUFNLE9BQU8sR0FBRyw2QkFBaUIsQ0FBQztJQUVsQyxnQ0FBZ0M7SUFDaEMsTUFBTSxjQUFjLEdBQUcsTUFBTSxDQUFBLE1BQU0sYUFBTixNQUFNLHVCQUFOLE1BQU0sQ0FBRSxJQUFJLENBQUMsT0FBTyxFQUFFLENBQUMsRUFBRSxJQUFJLENBQUMsQ0FBQSxDQUFDO0lBQzVELE9BQU8sQ0FBQyxHQUFHLENBQUMsY0FBYyxDQUFDLENBQUM7SUFDNUIsTUFBTSxNQUFNLEdBQUcsY0FBYyxDQUFDLE1BQU0sQ0FBQztJQUNyQyxPQUFPLENBQUMsR0FBRyxDQUFDLE1BQU0sQ0FBQyxDQUFDO0FBQ3RCLENBQUM7QUFFRCxJQUFJLEVBQUUsQ0FBQyxJQUFJLENBQUMsR0FBRyxFQUFFO0FBQ2pCLENBQUMsQ0FBQyxDQUFDLEtBQUssQ0FBQyxDQUFDLEtBQUssRUFBRSxFQUFFO0lBQ2pCLE9BQU8sQ0FBQyxHQUFHLENBQUMsS0FBSyxDQUFDLE9BQU8sQ0FBQyxDQUFDO0FBQzdCLENBQUMsQ0FBQyxDQUFDIn0=