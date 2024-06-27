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
// TODO: Test after testnet deploy latest transfer contracts.
async function test() {
    const wallet = await local_wallet_1.default.fromMnemonic(constants_2.DYDX_TEST_MNEMONIC, src_1.BECH32_PREFIX);
    console.log(wallet);
    const client = await validator_client_1.ValidatorClient.connect(constants_1.Network.testnet().validatorConfig);
    console.log('**Client**');
    console.log(client);
    const subaccount = new subaccount_1.SubaccountInfo(wallet, 0);
    const tx = await client.post.withdraw(subaccount, 0, new long_1.default(100000000));
    console.log('**Withdraw Tx**');
    console.log(tx);
}
test().then(() => {
}).catch((error) => {
    console.log(error.message);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidHJhbnNmZXJfZXhhbXBsZV93aXRoZHJhdy5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uL2V4YW1wbGVzL3RyYW5zZmVyX2V4YW1wbGVfd2l0aGRyYXcudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7QUFBQSxnREFBd0I7QUFFeEIsZ0NBQXVDO0FBQ3ZDLHdEQUFtRDtBQUNuRCx1RkFBOEQ7QUFDOUQsMERBQTJEO0FBQzNELHNFQUFrRTtBQUNsRSwyQ0FBaUQ7QUFFakQsNkRBQTZEO0FBRTdELEtBQUssVUFBVSxJQUFJO0lBQ2pCLE1BQU0sTUFBTSxHQUFHLE1BQU0sc0JBQVcsQ0FBQyxZQUFZLENBQUMsOEJBQWtCLEVBQUUsbUJBQWEsQ0FBQyxDQUFDO0lBQ2pGLE9BQU8sQ0FBQyxHQUFHLENBQUMsTUFBTSxDQUFDLENBQUM7SUFFcEIsTUFBTSxNQUFNLEdBQUcsTUFBTSxrQ0FBZSxDQUFDLE9BQU8sQ0FBQyxtQkFBTyxDQUFDLE9BQU8sRUFBRSxDQUFDLGVBQWUsQ0FBQyxDQUFDO0lBQ2hGLE9BQU8sQ0FBQyxHQUFHLENBQUMsWUFBWSxDQUFDLENBQUM7SUFDMUIsT0FBTyxDQUFDLEdBQUcsQ0FBQyxNQUFNLENBQUMsQ0FBQztJQUVwQixNQUFNLFVBQVUsR0FBRyxJQUFJLDJCQUFjLENBQUMsTUFBTSxFQUFFLENBQUMsQ0FBQyxDQUFDO0lBRWpELE1BQU0sRUFBRSxHQUFHLE1BQU0sTUFBTSxDQUFDLElBQUksQ0FBQyxRQUFRLENBQ25DLFVBQVUsRUFDVixDQUFDLEVBQ0QsSUFBSSxjQUFJLENBQUMsU0FBWSxDQUFDLENBQ3ZCLENBQUM7SUFDRixPQUFPLENBQUMsR0FBRyxDQUFDLGlCQUFpQixDQUFDLENBQUM7SUFDL0IsT0FBTyxDQUFDLEdBQUcsQ0FBQyxFQUFFLENBQUMsQ0FBQztBQUNsQixDQUFDO0FBRUQsSUFBSSxFQUFFLENBQUMsSUFBSSxDQUFDLEdBQUcsRUFBRTtBQUNqQixDQUFDLENBQUMsQ0FBQyxLQUFLLENBQUMsQ0FBQyxLQUFLLEVBQUUsRUFBRTtJQUNqQixPQUFPLENBQUMsR0FBRyxDQUFDLEtBQUssQ0FBQyxPQUFPLENBQUMsQ0FBQztBQUM3QixDQUFDLENBQUMsQ0FBQyJ9