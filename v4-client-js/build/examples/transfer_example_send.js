"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const tendermint_rpc_1 = require("@cosmjs/tendermint-rpc");
const long_1 = __importDefault(require("long"));
const constants_1 = require("../__tests__/helpers/constants");
const src_1 = require("../src");
const constants_2 = require("../src/clients/constants");
const local_wallet_1 = __importDefault(require("../src/clients/modules/local-wallet"));
const subaccount_1 = require("../src/clients/subaccount");
const validator_client_1 = require("../src/clients/validator-client");
const constants_3 = require("./constants");
async function test() {
    const wallet = await local_wallet_1.default.fromMnemonic(constants_3.DYDX_TEST_MNEMONIC, src_1.BECH32_PREFIX);
    console.log(wallet);
    const client = await validator_client_1.ValidatorClient.connect(constants_2.Network.testnet().validatorConfig);
    console.log('**Client**');
    console.log(client);
    const subaccount = new subaccount_1.SubaccountInfo(wallet, 0);
    const amount = new long_1.default(100000000);
    const msgs = new Promise((resolve) => {
        const msg = client.post.composer.composeMsgSendToken(subaccount.address, constants_1.TEST_RECIPIENT_ADDRESS, client.config.denoms.CHAINTOKEN_DENOM, amount.toString());
        resolve([msg]);
    });
    const totalFee = await client.post.simulate(subaccount.wallet, () => msgs, undefined, undefined);
    console.log('**Total Fee**');
    console.log(totalFee);
    const amountAfterFee = amount.sub(long_1.default.fromString(totalFee.amount[0].amount));
    console.log('**Amount after fee**');
    console.log(amountAfterFee);
    const tx = await client.post.sendToken(subaccount, constants_1.TEST_RECIPIENT_ADDRESS, client.config.denoms.CHAINTOKEN_DENOM, amountAfterFee.toString(), false, tendermint_rpc_1.Method.BroadcastTxCommit);
    console.log('**Send**');
    console.log(tx);
}
test()
    .then(() => { })
    .catch((error) => {
    console.log(error.message);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidHJhbnNmZXJfZXhhbXBsZV9zZW5kLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vZXhhbXBsZXMvdHJhbnNmZXJfZXhhbXBsZV9zZW5kLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7O0FBQ0EsMkRBQWdEO0FBQ2hELGdEQUF3QjtBQUV4Qiw4REFBd0U7QUFDeEUsZ0NBQXVDO0FBQ3ZDLHdEQUFtRDtBQUNuRCx1RkFBOEQ7QUFDOUQsMERBQTJEO0FBQzNELHNFQUFrRTtBQUNsRSwyQ0FBaUQ7QUFFakQsS0FBSyxVQUFVLElBQUk7SUFDakIsTUFBTSxNQUFNLEdBQUcsTUFBTSxzQkFBVyxDQUFDLFlBQVksQ0FDM0MsOEJBQWtCLEVBQ2xCLG1CQUFhLENBQ2QsQ0FBQztJQUNGLE9BQU8sQ0FBQyxHQUFHLENBQUMsTUFBTSxDQUFDLENBQUM7SUFFcEIsTUFBTSxNQUFNLEdBQUcsTUFBTSxrQ0FBZSxDQUFDLE9BQU8sQ0FBQyxtQkFBTyxDQUFDLE9BQU8sRUFBRSxDQUFDLGVBQWUsQ0FBQyxDQUFDO0lBQ2hGLE9BQU8sQ0FBQyxHQUFHLENBQUMsWUFBWSxDQUFDLENBQUM7SUFDMUIsT0FBTyxDQUFDLEdBQUcsQ0FBQyxNQUFNLENBQUMsQ0FBQztJQUVwQixNQUFNLFVBQVUsR0FBRyxJQUFJLDJCQUFjLENBQUMsTUFBTSxFQUFFLENBQUMsQ0FBQyxDQUFDO0lBRWpELE1BQU0sTUFBTSxHQUFHLElBQUksY0FBSSxDQUFDLFNBQVcsQ0FBQyxDQUFDO0lBRXJDLE1BQU0sSUFBSSxHQUE0QixJQUFJLE9BQU8sQ0FBQyxDQUFDLE9BQU8sRUFBRSxFQUFFO1FBQzVELE1BQU0sR0FBRyxHQUFHLE1BQU0sQ0FBQyxJQUFJLENBQUMsUUFBUSxDQUFDLG1CQUFtQixDQUNsRCxVQUFVLENBQUMsT0FBTyxFQUNsQixrQ0FBc0IsRUFDdEIsTUFBTSxDQUFDLE1BQU0sQ0FBQyxNQUFNLENBQUMsZ0JBQWdCLEVBQ3JDLE1BQU0sQ0FBQyxRQUFRLEVBQUUsQ0FDbEIsQ0FBQztRQUVGLE9BQU8sQ0FBQyxDQUFDLEdBQUcsQ0FBQyxDQUFDLENBQUM7SUFDakIsQ0FBQyxDQUFDLENBQUM7SUFFSCxNQUFNLFFBQVEsR0FBRyxNQUFNLE1BQU0sQ0FBQyxJQUFJLENBQUMsUUFBUSxDQUN6QyxVQUFVLENBQUMsTUFBTSxFQUNqQixHQUFHLEVBQUUsQ0FBQyxJQUFJLEVBQ1YsU0FBUyxFQUNULFNBQVMsQ0FDVixDQUFDO0lBQ0YsT0FBTyxDQUFDLEdBQUcsQ0FBQyxlQUFlLENBQUMsQ0FBQztJQUM3QixPQUFPLENBQUMsR0FBRyxDQUFDLFFBQVEsQ0FBQyxDQUFDO0lBRXRCLE1BQU0sY0FBYyxHQUFHLE1BQU0sQ0FBQyxHQUFHLENBQUMsY0FBSSxDQUFDLFVBQVUsQ0FBQyxRQUFRLENBQUMsTUFBTSxDQUFDLENBQUMsQ0FBQyxDQUFDLE1BQU0sQ0FBQyxDQUFDLENBQUM7SUFDOUUsT0FBTyxDQUFDLEdBQUcsQ0FBQyxzQkFBc0IsQ0FBQyxDQUFDO0lBQ3BDLE9BQU8sQ0FBQyxHQUFHLENBQUMsY0FBYyxDQUFDLENBQUM7SUFFNUIsTUFBTSxFQUFFLEdBQUcsTUFBTSxNQUFNLENBQUMsSUFBSSxDQUFDLFNBQVMsQ0FDcEMsVUFBVSxFQUNWLGtDQUFzQixFQUN0QixNQUFNLENBQUMsTUFBTSxDQUFDLE1BQU0sQ0FBQyxnQkFBZ0IsRUFDckMsY0FBYyxDQUFDLFFBQVEsRUFBRSxFQUN6QixLQUFLLEVBQ0wsdUJBQU0sQ0FBQyxpQkFBaUIsQ0FDekIsQ0FBQztJQUNGLE9BQU8sQ0FBQyxHQUFHLENBQUMsVUFBVSxDQUFDLENBQUM7SUFDeEIsT0FBTyxDQUFDLEdBQUcsQ0FBQyxFQUFFLENBQUMsQ0FBQztBQUNsQixDQUFDO0FBRUQsSUFBSSxFQUFFO0tBQ0gsSUFBSSxDQUFDLEdBQUcsRUFBRSxHQUFFLENBQUMsQ0FBQztLQUNkLEtBQUssQ0FBQyxDQUFDLEtBQUssRUFBRSxFQUFFO0lBQ2YsT0FBTyxDQUFDLEdBQUcsQ0FBQyxLQUFLLENBQUMsT0FBTyxDQUFDLENBQUM7QUFDN0IsQ0FBQyxDQUFDLENBQUMifQ==