"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const constants_1 = require("../../../src/clients/constants");
const faucet_client_1 = require("../../../src/clients/faucet-client");
const constants_2 = require("./constants");
describe('FaucetClient', () => {
    const client = new faucet_client_1.FaucetClient(constants_1.FaucetApiHost.TESTNET);
    describe('Faucet Endpoints', () => {
        it('Fill', async () => {
            const response = await client.fill(constants_2.DYDX_TEST_ADDRESS, 0, 2000);
            expect(response === null || response === void 0 ? void 0 : response.status).toBe(202);
        });
    });
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiRmF1Y2V0RW5kcG9pbnQudGVzdC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL19fdGVzdHNfXy9tb2R1bGVzL2NsaWVudC9GYXVjZXRFbmRwb2ludC50ZXN0LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7O0FBQUEsOERBQStEO0FBQy9ELHNFQUFrRTtBQUNsRSwyQ0FBZ0Q7QUFFaEQsUUFBUSxDQUFDLGNBQWMsRUFBRSxHQUFHLEVBQUU7SUFDNUIsTUFBTSxNQUFNLEdBQUcsSUFBSSw0QkFBWSxDQUFDLHlCQUFhLENBQUMsT0FBTyxDQUFDLENBQUM7SUFFdkQsUUFBUSxDQUFDLGtCQUFrQixFQUFFLEdBQUcsRUFBRTtRQUNoQyxFQUFFLENBQUMsTUFBTSxFQUFFLEtBQUssSUFBSSxFQUFFO1lBQ3BCLE1BQU0sUUFBUSxHQUFHLE1BQU0sTUFBTSxDQUFDLElBQUksQ0FBQyw2QkFBaUIsRUFBRSxDQUFDLEVBQUUsSUFBSSxDQUFDLENBQUM7WUFDL0QsTUFBTSxDQUFDLFFBQVEsYUFBUixRQUFRLHVCQUFSLFFBQVEsQ0FBRSxNQUFNLENBQUMsQ0FBQyxJQUFJLENBQUMsR0FBRyxDQUFDLENBQUM7UUFDckMsQ0FBQyxDQUFDLENBQUM7SUFDTCxDQUFDLENBQUMsQ0FBQztBQUNMLENBQUMsQ0FBQyxDQUFDIn0=