"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const amino_1 = require("@cosmjs/amino");
const proto_signing_1 = require("@cosmjs/proto-signing");
const stargate_1 = require("@cosmjs/stargate");
const long_1 = __importDefault(require("long"));
const protobufjs_1 = __importDefault(require("protobufjs"));
const registry_1 = require("../lib/registry");
const signer_1 = require("./signer");
// Required for encoding and decoding queries that are of type Long.
protobufjs_1.default.util.Long = long_1.default;
protobufjs_1.default.configure();
class LocalWallet {
    static async fromOfflineSigner(signer) {
        const wallet = new LocalWallet();
        await wallet.setSigner(signer);
        return wallet;
    }
    static async fromMnemonic(mnemonic, prefix) {
        const wallet = new LocalWallet();
        await wallet.setMnemonic(mnemonic, prefix);
        return wallet;
    }
    async setSigner(signer) {
        this.offlineSigner = signer;
        const stargateClient = await stargate_1.SigningStargateClient.offline(signer, {
            registry: (0, registry_1.generateRegistry)(),
        });
        const accountData = await signer.getAccounts();
        const firstAccount = accountData[0];
        this.accounts = [...accountData];
        this.address = firstAccount.address;
        this.pubKey = (0, amino_1.encodeSecp256k1Pubkey)(firstAccount.pubkey);
        this.signer = new signer_1.TransactionSigner(this.address, stargateClient);
    }
    async setMnemonic(mnemonic, prefix) {
        const signer = await proto_signing_1.DirectSecp256k1HdWallet.fromMnemonic(mnemonic, { prefix });
        return this.setSigner(signer);
    }
    async signTransaction(messages, transactionOptions, fee, memo = '') {
        return this.signer.signTransaction(messages, transactionOptions, fee, memo);
    }
}
exports.default = LocalWallet;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibG9jYWwtd2FsbGV0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2NsaWVudHMvbW9kdWxlcy9sb2NhbC13YWxsZXQudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7QUFBQSx5Q0FJdUI7QUFDdkIseURBSytCO0FBQy9CLCtDQUUwQjtBQUMxQixnREFBd0I7QUFDeEIsNERBQWtDO0FBRWxDLDhDQUFtRDtBQUVuRCxxQ0FBNkM7QUFFN0Msb0VBQW9FO0FBQ3BFLG9CQUFRLENBQUMsSUFBSSxDQUFDLElBQUksR0FBRyxjQUFJLENBQUM7QUFDMUIsb0JBQVEsQ0FBQyxTQUFTLEVBQUUsQ0FBQztBQUVyQixNQUFxQixXQUFXO0lBTzVCLE1BQU0sQ0FBQyxLQUFLLENBQUMsaUJBQWlCLENBQUMsTUFBb0I7UUFDakQsTUFBTSxNQUFNLEdBQUcsSUFBSSxXQUFXLEVBQUUsQ0FBQztRQUNqQyxNQUFNLE1BQU0sQ0FBQyxTQUFTLENBQUMsTUFBTSxDQUFDLENBQUM7UUFDL0IsT0FBTyxNQUFNLENBQUM7SUFDaEIsQ0FBQztJQUVELE1BQU0sQ0FBQyxLQUFLLENBQUMsWUFBWSxDQUFDLFFBQWdCLEVBQUUsTUFBZTtRQUN6RCxNQUFNLE1BQU0sR0FBRyxJQUFJLFdBQVcsRUFBRSxDQUFDO1FBQ2pDLE1BQU0sTUFBTSxDQUFDLFdBQVcsQ0FBQyxRQUFRLEVBQUUsTUFBTSxDQUFDLENBQUM7UUFDM0MsT0FBTyxNQUFNLENBQUM7SUFDaEIsQ0FBQztJQUVELEtBQUssQ0FBQyxTQUFTLENBQUMsTUFBcUI7UUFDbkMsSUFBSSxDQUFDLGFBQWEsR0FBRyxNQUFNLENBQUM7UUFDNUIsTUFBTSxjQUFjLEdBQUcsTUFBTSxnQ0FBcUIsQ0FBQyxPQUFPLENBQ3hELE1BQU0sRUFDTjtZQUNFLFFBQVEsRUFBRSxJQUFBLDJCQUFnQixHQUFFO1NBQzdCLENBQ0YsQ0FBQztRQUNGLE1BQU0sV0FBVyxHQUFHLE1BQU0sTUFBTSxDQUFDLFdBQVcsRUFBRSxDQUFDO1FBQy9DLE1BQU0sWUFBWSxHQUFHLFdBQVcsQ0FBQyxDQUFDLENBQUMsQ0FBQztRQUNwQyxJQUFJLENBQUMsUUFBUSxHQUFHLENBQUMsR0FBRyxXQUFXLENBQUMsQ0FBQztRQUNqQyxJQUFJLENBQUMsT0FBTyxHQUFHLFlBQVksQ0FBQyxPQUFPLENBQUM7UUFDcEMsSUFBSSxDQUFDLE1BQU0sR0FBRyxJQUFBLDZCQUFxQixFQUFDLFlBQVksQ0FBQyxNQUFNLENBQUMsQ0FBQztRQUN6RCxJQUFJLENBQUMsTUFBTSxHQUFHLElBQUksMEJBQWlCLENBQ2pDLElBQUksQ0FBQyxPQUFPLEVBQ1osY0FBYyxDQUNmLENBQUM7SUFDSixDQUFDO0lBRUQsS0FBSyxDQUFDLFdBQVcsQ0FBQyxRQUFnQixFQUFFLE1BQWU7UUFDakQsTUFBTSxNQUFNLEdBQUcsTUFBTSx1Q0FBdUIsQ0FBQyxZQUFZLENBQ3ZELFFBQVEsRUFDUixFQUFFLE1BQU0sRUFBRSxDQUNYLENBQUM7UUFDRixPQUFPLElBQUksQ0FBQyxTQUFTLENBQUMsTUFBTSxDQUFDLENBQUM7SUFDaEMsQ0FBQztJQUVNLEtBQUssQ0FBQyxlQUFlLENBQzFCLFFBQXdCLEVBQ3hCLGtCQUFzQyxFQUN0QyxHQUFZLEVBQ1osT0FBZSxFQUFFO1FBRWpCLE9BQU8sSUFBSSxDQUFDLE1BQU8sQ0FBQyxlQUFlLENBQUMsUUFBUSxFQUFFLGtCQUFrQixFQUFFLEdBQUcsRUFBRSxJQUFJLENBQUMsQ0FBQztJQUMvRSxDQUFDO0NBQ0o7QUF0REQsOEJBc0RDIn0=