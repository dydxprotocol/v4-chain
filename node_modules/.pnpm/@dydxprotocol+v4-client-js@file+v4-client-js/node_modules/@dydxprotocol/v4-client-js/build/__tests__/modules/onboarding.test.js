"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const onboarding_1 = require("../../src/lib/onboarding");
const proto_signing_1 = require("@cosmjs/proto-signing");
const constants_1 = require("../helpers/constants");
describe('Onboarding', () => {
    describe('exportMnemonicAndPrivateKey', () => {
        it('Successfully creates HDKey', () => {
            expect((0, onboarding_1.exportMnemonicAndPrivateKey)(constants_1.ENTROPY_FROM_SIGNATURE_RESULT))
                .toEqual({
                mnemonic: constants_1.MNEMONIC_FROM_SIGNATURE_RESULT,
                privateKey: constants_1.PRIVATE_KEY_FROM_SIGNATURE_RESULT,
                publicKey: constants_1.PUBLIC_KEY_FROM_SIGNATURE_RESULT,
            });
        });
        it('Expect mnemonic and private key to generate the same address', async () => {
            const { mnemonic, privateKey } = (0, onboarding_1.exportMnemonicAndPrivateKey)(constants_1.ENTROPY_FROM_SIGNATURE_RESULT);
            const wallet = await proto_signing_1.DirectSecp256k1HdWallet.fromMnemonic(mnemonic, {
                prefix: 'dydx',
            });
            const pkWallet = await proto_signing_1.DirectSecp256k1Wallet.fromKey(privateKey, 'dydx');
            const mnemonicAddress = (await wallet.getAccounts())[0].address;
            const pkAddress = (await pkWallet.getAccounts())[0].address;
            expect(mnemonicAddress).toEqual(pkAddress);
        });
    });
    describe('deriveHDKeyFromEthereumSignature', () => {
        it('Throw error if signature not valid', () => {
            expect(() => (0, onboarding_1.deriveHDKeyFromEthereumSignature)('0x04db21dfsa321b')).toThrow();
        });
        it('Successfully creates HDKey', () => {
            expect((0, onboarding_1.deriveHDKeyFromEthereumSignature)(constants_1.SIGNATURE_RESULT)).toEqual({
                mnemonic: constants_1.MNEMONIC_FROM_SIGNATURE_RESULT,
                privateKey: constants_1.PRIVATE_KEY_FROM_SIGNATURE_RESULT,
                publicKey: constants_1.PUBLIC_KEY_FROM_SIGNATURE_RESULT,
            });
        });
        it('Successfully creates HDKey from signature with different v values', () => {
            const suffixes = ['00', '1b', '01', '1c'];
            const testSignatures = suffixes.map((suffix) => constants_1.SIGNATURE_RESULT.slice(0, -2) + suffix);
            testSignatures.forEach((sig) => {
                expect((0, onboarding_1.deriveHDKeyFromEthereumSignature)(sig)).toEqual({
                    mnemonic: constants_1.MNEMONIC_FROM_SIGNATURE_RESULT,
                    privateKey: constants_1.PRIVATE_KEY_FROM_SIGNATURE_RESULT,
                    publicKey: constants_1.PUBLIC_KEY_FROM_SIGNATURE_RESULT,
                });
            });
        });
    });
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoib25ib2FyZGluZy50ZXN0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vX190ZXN0c19fL21vZHVsZXMvb25ib2FyZGluZy50ZXN0LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7O0FBQUEseURBQXlHO0FBQ3pHLHlEQUkrQjtBQUMvQixvREFNOEI7QUFFOUIsUUFBUSxDQUFDLFlBQVksRUFBRSxHQUFHLEVBQUU7SUFDMUIsUUFBUSxDQUFDLDZCQUE2QixFQUFFLEdBQUcsRUFBRTtRQUMzQyxFQUFFLENBQUMsNEJBQTRCLEVBQUUsR0FBRyxFQUFFO1lBQ3BDLE1BQU0sQ0FBQyxJQUFBLHdDQUEyQixFQUFDLHlDQUE2QixDQUFDLENBQUM7aUJBQy9ELE9BQU8sQ0FBQztnQkFDUCxRQUFRLEVBQUUsMENBQThCO2dCQUN4QyxVQUFVLEVBQUUsNkNBQWlDO2dCQUM3QyxTQUFTLEVBQUUsNENBQWdDO2FBQzVDLENBQUMsQ0FBQztRQUNQLENBQUMsQ0FBQyxDQUFDO1FBRUgsRUFBRSxDQUFDLDhEQUE4RCxFQUFFLEtBQUssSUFBSSxFQUFFO1lBQzVFLE1BQU0sRUFBRSxRQUFRLEVBQUUsVUFBVSxFQUFFLEdBQUcsSUFBQSx3Q0FBMkIsRUFBQyx5Q0FBNkIsQ0FBQyxDQUFDO1lBRTVGLE1BQU0sTUFBTSxHQUFrQixNQUFNLHVDQUF1QixDQUFDLFlBQVksQ0FBQyxRQUFTLEVBQUU7Z0JBQ2xGLE1BQU0sRUFBRSxNQUFNO2FBQ2YsQ0FBQyxDQUFDO1lBQ0gsTUFBTSxRQUFRLEdBQUcsTUFBTSxxQ0FBcUIsQ0FBQyxPQUFPLENBQUMsVUFBVyxFQUFFLE1BQU0sQ0FBQyxDQUFDO1lBQzFFLE1BQU0sZUFBZSxHQUFHLENBQUMsTUFBTSxNQUFNLENBQUMsV0FBVyxFQUFFLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxPQUFPLENBQUM7WUFDaEUsTUFBTSxTQUFTLEdBQUcsQ0FBQyxNQUFNLFFBQVEsQ0FBQyxXQUFXLEVBQUUsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLE9BQU8sQ0FBQztZQUM1RCxNQUFNLENBQUMsZUFBZSxDQUFDLENBQUMsT0FBTyxDQUFDLFNBQVMsQ0FBQyxDQUFDO1FBQzdDLENBQUMsQ0FBQyxDQUFDO0lBQ0wsQ0FBQyxDQUFDLENBQUM7SUFFSCxRQUFRLENBQUMsa0NBQWtDLEVBQUUsR0FBRyxFQUFFO1FBQ2hELEVBQUUsQ0FBQyxvQ0FBb0MsRUFBRSxHQUFHLEVBQUU7WUFDNUMsTUFBTSxDQUFDLEdBQUcsRUFBRSxDQUFDLElBQUEsNkNBQWdDLEVBQUMsa0JBQWtCLENBQUMsQ0FBQyxDQUFDLE9BQU8sRUFBRSxDQUFDO1FBQy9FLENBQUMsQ0FBQyxDQUFDO1FBQ0gsRUFBRSxDQUFDLDRCQUE0QixFQUFFLEdBQUcsRUFBRTtZQUNwQyxNQUFNLENBQUMsSUFBQSw2Q0FBZ0MsRUFBQyw0QkFBZ0IsQ0FBQyxDQUFDLENBQUMsT0FBTyxDQUFDO2dCQUNqRSxRQUFRLEVBQUUsMENBQThCO2dCQUN4QyxVQUFVLEVBQUUsNkNBQWlDO2dCQUM3QyxTQUFTLEVBQUUsNENBQWdDO2FBQzVDLENBQUMsQ0FBQztRQUNMLENBQUMsQ0FBQyxDQUFDO1FBRUgsRUFBRSxDQUFDLG1FQUFtRSxFQUFFLEdBQUcsRUFBRTtZQUMzRSxNQUFNLFFBQVEsR0FBRyxDQUFDLElBQUksRUFBRSxJQUFJLEVBQUUsSUFBSSxFQUFFLElBQUksQ0FBQyxDQUFDO1lBRTFDLE1BQU0sY0FBYyxHQUFHLFFBQVEsQ0FBQyxHQUFHLENBQUMsQ0FBQyxNQUFNLEVBQUUsRUFBRSxDQUFDLDRCQUFnQixDQUFDLEtBQUssQ0FBQyxDQUFDLEVBQUUsQ0FBQyxDQUFDLENBQUMsR0FBRyxNQUFNLENBQUMsQ0FBQztZQUN4RixjQUFjLENBQUMsT0FBTyxDQUFDLENBQUMsR0FBRyxFQUFFLEVBQUU7Z0JBQzdCLE1BQU0sQ0FBQyxJQUFBLDZDQUFnQyxFQUFDLEdBQUcsQ0FBQyxDQUFDLENBQUMsT0FBTyxDQUFDO29CQUNwRCxRQUFRLEVBQUUsMENBQThCO29CQUN4QyxVQUFVLEVBQUUsNkNBQWlDO29CQUM3QyxTQUFTLEVBQUUsNENBQWdDO2lCQUM1QyxDQUFDLENBQUM7WUFDTCxDQUFDLENBQUMsQ0FBQztRQUNMLENBQUMsQ0FBQyxDQUFDO0lBQ0wsQ0FBQyxDQUFDLENBQUM7QUFDTCxDQUFDLENBQUMsQ0FBQyJ9