"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const constants_1 = require("../../../src/clients/constants");
const indexer_client_1 = require("../../../src/clients/indexer-client");
const constants_2 = require("./constants");
describe('IndexerClient', () => {
    const client = new indexer_client_1.IndexerClient(constants_1.Network.testnet().indexerConfig);
    describe('Private Endpoints', () => {
        it('Subaccounts', async () => {
            const response = await client.account.getSubaccounts(constants_2.DYDX_TEST_ADDRESS);
            const subaccounts = response.subaccounts;
            expect(subaccounts.size).not.toBeNull();
            const subaccount0 = subaccounts[0];
            const subaccountNumber = subaccount0.subaccountNumber;
            expect(subaccountNumber).not.toBeNull();
        });
        it('Subaccount 0', async () => {
            const response = await client.account.getSubaccount(constants_2.DYDX_TEST_ADDRESS, 0);
            const subaccount = response.subaccount;
            expect(subaccount).not.toBeNull();
            const subaccountNumber = subaccount.subaccountNumber;
            expect(subaccountNumber).not.toBeNull();
        });
        it('Asset Positions', async () => {
            const response = await client.account.getSubaccountAssetPositions(constants_2.DYDX_TEST_ADDRESS, 0);
            expect(response).not.toBeNull();
            const positions = response.positions;
            expect(positions).not.toBeNull();
            if (positions.length > 0) {
                const position = positions[0];
                expect(position).not.toBeNull();
            }
        });
        it('Perpetual Positions', async () => {
            const response = await client.account.getSubaccountPerpetualPositions(constants_2.DYDX_TEST_ADDRESS, 0);
            expect(response).not.toBeNull();
            const positions = response.positions;
            expect(positions).not.toBeNull();
            if (positions.length > 0) {
                const position = positions[0];
                expect(position).not.toBeNull();
            }
        });
        it('Transfers', async () => {
            const response = await client.account.getSubaccountTransfers(constants_2.DYDX_TEST_ADDRESS, 0);
            expect(response).not.toBeNull();
            const transfers = response.transfers;
            expect(transfers).not.toBeNull();
            if (transfers.length > 0) {
                const transfer = transfers[0];
                expect(transfer).not.toBeNull();
            }
        });
        it('Orders', async () => {
            const response = await client.account.getSubaccountOrders(constants_2.DYDX_TEST_ADDRESS, 0);
            expect(response).not.toBeNull();
            const orders = response;
            expect(orders).not.toBeNull();
            if (orders.length > 0) {
                const order = orders[0];
                expect(order).not.toBeNull();
            }
        });
        it('Fills', async () => {
            const response = await client.account.getSubaccountFills(constants_2.DYDX_TEST_ADDRESS, 0);
            expect(response).not.toBeNull();
            const fills = response.fills;
            expect(fills).not.toBeNull();
            if (fills.length > 0) {
                const fill = fills[0];
                expect(fill).not.toBeNull();
            }
        });
        it('Historical PNL', async () => {
            const response = await client.account.getSubaccountHistoricalPNLs(constants_2.DYDX_TEST_ADDRESS, 0);
            expect(response).not.toBeNull();
            const historicalPnl = response.historicalPnl;
            expect(historicalPnl).not.toBeNull();
            if (historicalPnl.length > 0) {
                const historicalPnl0 = historicalPnl[0];
                expect(historicalPnl0).not.toBeNull();
            }
        });
    });
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiQWNjb3VudEVuZHBvaW50cy50ZXN0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vX190ZXN0c19fL21vZHVsZXMvY2xpZW50L0FjY291bnRFbmRwb2ludHMudGVzdC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOztBQUFBLDhEQUF5RDtBQUN6RCx3RUFBb0U7QUFDcEUsMkNBQWdEO0FBRWhELFFBQVEsQ0FBQyxlQUFlLEVBQUUsR0FBRyxFQUFFO0lBQzdCLE1BQU0sTUFBTSxHQUFHLElBQUksOEJBQWEsQ0FBQyxtQkFBTyxDQUFDLE9BQU8sRUFBRSxDQUFDLGFBQWEsQ0FBQyxDQUFDO0lBRWxFLFFBQVEsQ0FBQyxtQkFBbUIsRUFBRSxHQUFHLEVBQUU7UUFDakMsRUFBRSxDQUFDLGFBQWEsRUFBRSxLQUFLLElBQUksRUFBRTtZQUMzQixNQUFNLFFBQVEsR0FBRyxNQUFNLE1BQU0sQ0FBQyxPQUFPLENBQUMsY0FBYyxDQUFDLDZCQUFpQixDQUFDLENBQUM7WUFDeEUsTUFBTSxXQUFXLEdBQUcsUUFBUSxDQUFDLFdBQVcsQ0FBQztZQUN6QyxNQUFNLENBQUMsV0FBVyxDQUFDLElBQUksQ0FBQyxDQUFDLEdBQUcsQ0FBQyxRQUFRLEVBQUUsQ0FBQztZQUN4QyxNQUFNLFdBQVcsR0FBRyxXQUFXLENBQUMsQ0FBQyxDQUFDLENBQUM7WUFDbkMsTUFBTSxnQkFBZ0IsR0FBRyxXQUFXLENBQUMsZ0JBQWdCLENBQUM7WUFDdEQsTUFBTSxDQUFDLGdCQUFnQixDQUFDLENBQUMsR0FBRyxDQUFDLFFBQVEsRUFBRSxDQUFDO1FBQzFDLENBQUMsQ0FBQyxDQUFDO1FBRUgsRUFBRSxDQUFDLGNBQWMsRUFBRSxLQUFLLElBQUksRUFBRTtZQUM1QixNQUFNLFFBQVEsR0FBRyxNQUFNLE1BQU0sQ0FBQyxPQUFPLENBQUMsYUFBYSxDQUFDLDZCQUFpQixFQUFFLENBQUMsQ0FBQyxDQUFDO1lBQzFFLE1BQU0sVUFBVSxHQUFHLFFBQVEsQ0FBQyxV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