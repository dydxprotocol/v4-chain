"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.NobleClient = void 0;
const proto_signing_1 = require("@cosmjs/proto-signing");
const stargate_1 = require("@cosmjs/stargate");
const constants_1 = require("./constants");
const cctpProto_1 = require("./lib/cctpProto");
class NobleClient {
    constructor(restEndpoint) {
        this.restEndpoint = restEndpoint;
    }
    get isConnected() {
        return Boolean(this.stargateClient);
    }
    async connect(wallet) {
        if ((wallet === null || wallet === void 0 ? void 0 : wallet.offlineSigner) === undefined) {
            throw new Error('Wallet signer not found');
        }
        this.wallet = wallet;
        this.stargateClient = await stargate_1.SigningStargateClient.connectWithSigner(this.restEndpoint, wallet.offlineSigner, {
            registry: new proto_signing_1.Registry([
                ['/circle.cctp.v1.MsgDepositForBurn', cctpProto_1.MsgDepositForBurn],
                ...stargate_1.defaultRegistryTypes,
            ]),
        });
    }
    getAccountBalances() {
        var _a;
        if (!this.stargateClient || ((_a = this.wallet) === null || _a === void 0 ? void 0 : _a.address) === undefined) {
            throw new Error('stargateClient not initialized');
        }
        return this.stargateClient.getAllBalances(this.wallet.address);
    }
    getAccountBalance(denom) {
        var _a;
        if (!this.stargateClient || ((_a = this.wallet) === null || _a === void 0 ? void 0 : _a.address) === undefined) {
            throw new Error('stargateClient not initialized');
        }
        return this.stargateClient.getBalance(this.wallet.address, denom);
    }
    async send(messages, gasPrice = stargate_1.GasPrice.fromString('0.025utdai'), memo) {
        var _a;
        if (!this.stargateClient) {
            throw new Error('NobleClient stargateClient not initialized');
        }
        if (((_a = this.wallet) === null || _a === void 0 ? void 0 : _a.address) === undefined) {
            throw new Error('NobleClient wallet not initialized');
        }
        // Simulate to get the gas estimate
        const fee = await this.simulateTransaction(messages, gasPrice, memo);
        // Sign and broadcast the transaction
        return this.stargateClient.signAndBroadcast(this.wallet.address, messages, fee, memo !== null && memo !== void 0 ? memo : '');
    }
    async simulateTransaction(messages, gasPrice = stargate_1.GasPrice.fromString('0.025utdai'), memo) {
        var _a, _b;
        if (!this.stargateClient) {
            throw new Error('NobleClient stargateClient not initialized');
        }
        if (((_a = this.wallet) === null || _a === void 0 ? void 0 : _a.address) === undefined) {
            throw new Error('NobleClient wallet not initialized');
        }
        // Get simulated response
        const gasEstimate = await this.stargateClient.simulate((_b = this.wallet) === null || _b === void 0 ? void 0 : _b.address, messages, memo);
        // Calculate and return the fee
        return (0, stargate_1.calculateFee)(Math.floor(gasEstimate * constants_1.GAS_MULTIPLIER), gasPrice);
    }
}
exports.NobleClient = NobleClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibm9ibGUtY2xpZW50LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL2NsaWVudHMvbm9ibGUtY2xpZW50LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7OztBQUFBLHlEQUFxRTtBQUNyRSwrQ0FPMEI7QUFFMUIsMkNBQTZDO0FBQzdDLCtDQUFvRDtBQUdwRCxNQUFhLFdBQVc7SUFLdEIsWUFBWSxZQUFvQjtRQUM5QixJQUFJLENBQUMsWUFBWSxHQUFHLFlBQVksQ0FBQztJQUNuQyxDQUFDO0lBRUQsSUFBSSxXQUFXO1FBQ2IsT0FBTyxPQUFPLENBQUMsSUFBSSxDQUFDLGNBQWMsQ0FBQyxDQUFDO0lBQ3RDLENBQUM7SUFFRCxLQUFLLENBQUMsT0FBTyxDQUFDLE1BQW1CO1FBQy9CLElBQUksQ0FBQSxNQUFNLGFBQU4sTUFBTSx1QkFBTixNQUFNLENBQUUsYUFBYSxNQUFLLFNBQVMsRUFBRSxDQUFDO1lBQ3hDLE1BQU0sSUFBSSxLQUFLLENBQUMseUJBQXlCLENBQUMsQ0FBQztRQUM3QyxDQUFDO1FBQ0QsSUFBSSxDQUFDLE1BQU0sR0FBRyxNQUFNLENBQUM7UUFDckIsSUFBSSxDQUFDLGNBQWMsR0FBRyxNQUFNLGdDQUFxQixDQUFDLGlCQUFpQixDQUNqRSxJQUFJLENBQUMsWUFBWSxFQUNqQixNQUFNLENBQUMsYUFBYSxFQUNwQjtZQUNFLFFBQVEsRUFBRSxJQUFJLHdCQUFRLENBQUM7Z0JBQ3JCLENBQUMsbUNBQW1DLEVBQUUsNkJBQWlCLENBQUM7Z0JBQ3hELEdBQUcsK0JBQW9CO2FBQ3hCLENBQUM7U0FDSCxDQUNGLENBQUM7SUFDSixDQUFDO0lBRUQsa0JBQWtCOztRQUNoQixJQUFJLENBQUMsSUFBSSxDQUFDLGNBQWMsSUFBSSxDQUFBLE1BQUEsSUFBSSxDQUFDLE1BQU0sMENBQUUsT0FBTyxNQUFLLFNBQVMsRUFBRSxDQUFDO1lBQy9ELE1BQU0sSUFBSSxLQUFLLENBQUMsZ0NBQWdDLENBQUMsQ0FBQztRQUNwRCxDQUFDO1FBQ0QsT0FBTyxJQUFJLENBQUMsY0FBYyxDQUFDLGNBQWMsQ0FBQyxJQUFJLENBQUMsTUFBTSxDQUFDLE9BQU8sQ0FBQyxDQUFDO0lBQ2pFLENBQUM7SUFFRCxpQkFBaUIsQ0FBQyxLQUFhOztRQUM3QixJQUFJLENBQUMsSUFBSSxDQUFDLGNBQWMsSUFBSSxDQUFBLE1BQUEsSUFBSSxDQUFDLE1BQU0sMENBQUUsT0FBTyxNQUFLLFNBQVMsRUFBRSxDQUFDO1lBQy9ELE1BQU0sSUFBSSxLQUFLLENBQUMsZ0NBQWdDLENBQUMsQ0FBQztRQUNwRCxDQUFDO1FBQ0QsT0FBTyxJQUFJLENBQUMsY0FBYyxDQUFDLFVBQVUsQ0FBQyxJQUFJLENBQUMsTUFBTSxDQUFDLE9BQU8sRUFBRSxLQUFLLENBQUMsQ0FBQztJQUNwRSxDQUFDO0lBRUQsS0FBSyxDQUFDLElBQUksQ0FDUixRQUF3QixFQUN4QixXQUFxQixtQkFBUSxDQUFDLFVBQVUsQ0FBQyxZQUFZLENBQUMsRUFDdEQsSUFBYTs7UUFFYixJQUFJLENBQUMsSUFBSSxDQUFDLGNBQWMsRUFBRSxDQUFDO1lBQ3pCLE1BQU0sSUFBSSxLQUFLLENBQUMsNENBQTRDLENBQUMsQ0FBQztRQUNoRSxDQUFDO1FBQ0QsSUFBSSxDQUFBLE1BQUEsSUFBSSxDQUFDLE1BQU0sMENBQUUsT0FBTyxNQUFLLFNBQVMsRUFBRSxDQUFDO1lBQ3ZDLE1BQU0sSUFBSSxLQUFLLENBQUMsb0NBQW9DLENBQUMsQ0FBQztRQUN4RCxDQUFDO1FBQ0QsbUNBQW1DO1FBQ25DLE1BQU0sR0FBRyxHQUFHLE1BQU0sSUFBSSxDQUFDLG1CQUFtQixDQUFDLFFBQVEsRUFBRSxRQUFRLEVBQUUsSUFBSSxDQUFDLENBQUM7UUFFckUscUNBQXFDO1FBQ3JDLE9BQU8sSUFBSSxDQUFDLGNBQWMsQ0FBQyxnQkFBZ0IsQ0FDekMsSUFBSSxDQUFDLE1BQU0sQ0FBQyxPQUFPLEVBQ25CLFFBQVEsRUFDUixHQUFHLEVBQ0gsSUFBSSxhQUFKLElBQUksY0FBSixJQUFJLEdBQUksRUFBRSxDQUNYLENBQUM7SUFDSixDQUFDO0lBRUQsS0FBSyxDQUFDLG1CQUFtQixDQUN2QixRQUFpQyxFQUNqQyxXQUFxQixtQkFBUSxDQUFDLFVBQVUsQ0FBQyxZQUFZLENBQUMsRUFDdEQsSUFBYTs7UUFFYixJQUFJLENBQUMsSUFBSSxDQUFDLGNBQWMsRUFBRSxDQUFDO1lBQ3pCLE1BQU0sSUFBSSxLQUFLLENBQUMsNENBQTRDLENBQUMsQ0FBQztRQUNoRSxDQUFDO1FBQ0QsSUFBSSxDQUFBLE1BQUEsSUFBSSxDQUFDLE1BQU0sMENBQUUsT0FBTyxNQUFLLFNBQVMsRUFBRSxDQUFDO1lBQ3ZDLE1BQU0sSUFBSSxLQUFLLENBQUMsb0NBQW9DLENBQUMsQ0FBQztRQUN4RCxDQUFDO1FBQ0QseUJBQXlCO1FBQ3pCLE1BQU0sV0FBVyxHQUFHLE1BQU0sSUFBSSxDQUFDLGNBQWMsQ0FBQyxRQUFRLENBQ3BELE1BQUEsSUFBSSxDQUFDLE1BQU0sMENBQUUsT0FBTyxFQUNwQixRQUFRLEVBQ1IsSUFBSSxDQUNMLENBQUM7UUFFRiwrQkFBK0I7UUFDL0IsT0FBTyxJQUFBLHVCQUFZLEVBQUMsSUFBSSxDQUFDLEtBQUssQ0FBQyxXQUFXLEdBQUcsMEJBQWMsQ0FBQyxFQUFFLFFBQVEsQ0FBQyxDQUFDO0lBQzFFLENBQUM7Q0FDRjtBQXhGRCxrQ0F3RkMifQ==