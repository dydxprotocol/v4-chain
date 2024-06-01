"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const src_1 = require("../../../src");
const constants_1 = require("../../../src/clients/constants");
const constants_2 = require("./constants");
describe('IndexerClient', () => {
    const client = new src_1.IndexerClient(constants_1.Network.testnet().indexerConfig);
    describe('Utility Endpoints', () => {
        it('getTime', async () => {
            const response = await client.utility.getTime();
            const iso = response.iso;
            expect(iso).not.toBeUndefined();
        });
        it('getHeight', async () => {
            const response = await client.utility.getHeight();
            const height = response.height;
            const time = response.time;
            expect(height).not.toBeUndefined();
            expect(time).not.toBeUndefined();
        });
        it('Screen Address', async () => {
            const response = await client.utility.screen(constants_2.DYDX_TEST_ADDRESS);
            const { restricted } = response !== null && response !== void 0 ? response : {};
            expect(restricted).toBeDefined();
        });
    });
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiVXRpbGl0eUVuZHBvaW50cy50ZXN0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vX190ZXN0c19fL21vZHVsZXMvY2xpZW50L1V0aWxpdHlFbmRwb2ludHMudGVzdC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOztBQUFBLHNDQUE2QztBQUM3Qyw4REFBeUQ7QUFDekQsMkNBQWdEO0FBRWhELFFBQVEsQ0FBQyxlQUFlLEVBQUUsR0FBRyxFQUFFO0lBQzdCLE1BQU0sTUFBTSxHQUFHLElBQUksbUJBQWEsQ0FBQyxtQkFBTyxDQUFDLE9BQU8sRUFBRSxDQUFDLGFBQWEsQ0FBQyxDQUFDO0lBRWxFLFFBQVEsQ0FBQyxtQkFBbUIsRUFBRSxHQUFHLEVBQUU7UUFDakMsRUFBRSxDQUFDLFNBQVMsRUFBRSxLQUFLLElBQUksRUFBRTtZQUN2QixNQUFNLFFBQVEsR0FBRyxNQUFNLE1BQU0sQ0FBQyxPQUFPLENBQUMsT0FBTyxFQUFFLENBQUM7WUFDaEQsTUFBTSxHQUFHLEdBQUcsUUFBUSxDQUFDLEdBQUcsQ0FBQztZQUN6QixNQUFNLENBQUMsR0FBRyxDQUFDLENBQUMsR0FBRyxDQUFDLGFBQWEsRUFBRSxDQUFDO1FBQ2xDLENBQUMsQ0FBQyxDQUFDO1FBRUgsRUFBRSxDQUFDLFdBQVcsRUFBRSxLQUFLLElBQUksRUFBRTtZQUN6QixNQUFNLFFBQVEsR0FBRyxNQUFNLE1BQU0sQ0FBQyxPQUFPLENBQUMsU0FBUyxFQUFFLENBQUM7WUFDbEQsTUFBTSxNQUFNLEdBQUcsUUFBUSxDQUFDLE1BQU0sQ0FBQztZQUMvQixNQUFNLElBQUksR0FBRyxRQUFRLENBQUMsSUFBSSxDQUFDO1lBQzNCLE1BQU0sQ0FBQyxNQUFNLENBQUMsQ0FBQyxHQUFHLENBQUMsYUFBYSxFQUFFLENBQUM7WUFDbkMsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDLEdBQUcsQ0FBQyxhQUFhLEVBQUUsQ0FBQztRQUNuQyxDQUFDLENBQUMsQ0FBQztRQUVILEVBQUUsQ0FBQyxnQkFBZ0IsRUFBRSxLQUFLLElBQUksRUFBRTtZQUM5QixNQUFNLFFBQVEsR0FBRyxNQUFNLE1BQU0sQ0FBQyxPQUFPLENBQUMsTUFBTSxDQUFDLDZCQUFpQixDQUFDLENBQUM7WUFDaEUsTUFBTSxFQUFFLFVBQVUsRUFBRSxHQUFHLFFBQVEsYUFBUixRQUFRLGNBQVIsUUFBUSxHQUFJLEVBQUUsQ0FBQztZQUN0QyxNQUFNLENBQUMsVUFBVSxDQUFDLENBQUMsV0FBVyxFQUFFLENBQUM7UUFDbkMsQ0FBQyxDQUFDLENBQUM7SUFDTCxDQUFDLENBQUMsQ0FBQztBQUNMLENBQUMsQ0FBQyxDQUFDIn0=