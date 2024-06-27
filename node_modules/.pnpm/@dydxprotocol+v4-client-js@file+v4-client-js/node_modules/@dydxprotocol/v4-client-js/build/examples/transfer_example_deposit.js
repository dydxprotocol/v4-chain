"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const long_1 = __importDefault(require("long"));
const src_1 = require("../src");
const constants_1 = require("../src/clients/constants");
const local_wallet_1 = __importDefault(require("../src/clients/modules/local-wallet"));
const subaccount_1 = require("../src/clients/subaccount");
const validator_client_1 = require("../src/clients/validator-client");
const constants_2 = require("./constants");
async function test() {
    const wallet = await local_wallet_1.default.fromMnemonic(constants_2.DYDX_TEST_MNEMONIC, src_1.BECH32_PREFIX);
    console.log(wallet);
    const client = await validator_client_1.ValidatorClient.connect(constants_1.Network.testnet().validatorConfig);
    console.log('**Client**');
    console.log(client);
    const subaccount = new subaccount_1.SubaccountInfo(wallet, 0);
    const tx = await client.post.deposit(subaccount, 0, new long_1.default(10000000));
    console.log('**Deposit Tx**');
    console.log(tx);
}
test().then(() => {
}).catch((error) => {
    console.log(error.message);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidHJhbnNmZXJfZXhhbXBsZV9kZXBvc2l0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vZXhhbXBsZXMvdHJhbnNmZXJfZXhhbXBsZV9kZXBvc2l0LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7O0FBQUEsZ0RBQXdCO0FBRXhCLGdDQUF1QztBQUN2Qyx3REFBbUQ7QUFDbkQsdUZBQThEO0FBQzlELDBEQUEyRDtBQUMzRCxzRUFBa0U7QUFDbEUsMkNBQWlEO0FBRWpELEtBQUssVUFBVSxJQUFJO0lBQ2pCLE1BQU0sTUFBTSxHQUFHLE1BQU0sc0JBQVcsQ0FBQyxZQUFZLENBQUMsOEJBQWtCLEVBQUUsbUJBQWEsQ0FBQyxDQUFDO0lBQ2pGLE9BQU8sQ0FBQyxHQUFHLENBQUMsTUFBTSxDQUFDLENBQUM7SUFFcEIsTUFBTSxNQUFNLEdBQUcsTUFBTSxrQ0FBZSxDQUFDLE9BQU8sQ0FBQyxtQkFBTyxDQUFDLE9BQU8sRUFBRSxDQUFDLGVBQWUsQ0FBQyxDQUFDO0lBQ2hGLE9BQU8sQ0FBQyxHQUFHLENBQUMsWUFBWSxDQUFDLENBQUM7SUFDMUIsT0FBTyxDQUFDLEdBQUcsQ0FBQyxNQUFNLENBQUMsQ0FBQztJQUVwQixNQUFNLFVBQVUsR0FBRyxJQUFJLDJCQUFjLENBQUMsTUFBTSxFQUFFLENBQUMsQ0FBQyxDQUFDO0lBQ2pELE1BQU0sRUFBRSxHQUFHLE1BQU0sTUFBTSxDQUFDLElBQUksQ0FBQyxPQUFPLENBQ2xDLFVBQVUsRUFDVixDQUFDLEVBQ0QsSUFBSSxjQUFJLENBQUMsUUFBVSxDQUFDLENBQ3JCLENBQUM7SUFDRixPQUFPLENBQUMsR0FBRyxDQUFDLGdCQUFnQixDQUFDLENBQUM7SUFDOUIsT0FBTyxDQUFDLEdBQUcsQ0FBQyxFQUFFLENBQUMsQ0FBQztBQUNsQixDQUFDO0FBRUQsSUFBSSxFQUFFLENBQUMsSUFBSSxDQUFDLEdBQUcsRUFBRTtBQUNqQixDQUFDLENBQUMsQ0FBQyxLQUFLLENBQUMsQ0FBQyxLQUFLLEVBQUUsRUFBRTtJQUNqQixPQUFPLENBQUMsR0FBRyxDQUFDLEtBQUssQ0FBQyxPQUFPLENBQUMsQ0FBQztBQUM3QixDQUFDLENBQUMsQ0FBQyJ9