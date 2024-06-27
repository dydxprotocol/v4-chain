"use strict";
/**
 * Simple JS example demostrating accessing public data with Indexer REST endpoints
 */
Object.defineProperty(exports, "__esModule", { value: true });
const constants_1 = require("../src/clients/constants");
const indexer_client_1 = require("../src/clients/indexer-client");
async function test() {
    const client = new indexer_client_1.IndexerClient(constants_1.Network.testnet().indexerConfig);
    // Get indexer server time
    try {
        const response = await client.utility.getTime();
        console.log(response);
        const timeIso = response.iso;
        const timeEpoch = response.epoch;
        console.log('time');
        console.log(timeIso);
        console.log(timeEpoch);
    }
    catch (error) {
        console.log(error.message);
    }
    // Get indexer server height
    try {
        const response = await client.utility.getHeight();
        console.log(response);
        const height = response.height;
        const heightTime = response.time;
        console.log('height');
        console.log(height);
        console.log(heightTime);
    }
    catch (error) {
        console.log(error.message);
    }
}
test().then(() => {
}).catch((error) => {
    console.log(error.message);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidXRpbGl0eV9lbmRwb2ludHMuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi9leGFtcGxlcy91dGlsaXR5X2VuZHBvaW50cy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiO0FBQUE7O0dBRUc7O0FBRUgsd0RBQW1EO0FBQ25ELGtFQUE4RDtBQUU5RCxLQUFLLFVBQVUsSUFBSTtJQUNqQixNQUFNLE1BQU0sR0FBRyxJQUFJLDhCQUFhLENBQUMsbUJBQU8sQ0FBQyxPQUFPLEVBQUUsQ0FBQyxhQUFhLENBQUMsQ0FBQztJQUVsRSwwQkFBMEI7SUFDMUIsSUFBSTtRQUNGLE1BQU0sUUFBUSxHQUFHLE1BQU0sTUFBTSxDQUFDLE9BQU8sQ0FBQyxPQUFPLEVBQUUsQ0FBQztRQUNoRCxPQUFPLENBQUMsR0FBRyxDQUFDLFFBQVEsQ0FBQyxDQUFDO1FBQ3RCLE1BQU0sT0FBTyxHQUFHLFFBQVEsQ0FBQyxHQUFHLENBQUM7UUFDN0IsTUFBTSxTQUFTLEdBQUcsUUFBUSxDQUFDLEtBQUssQ0FBQztRQUNqQyxPQUFPLENBQUMsR0FBRyxDQUFDLE1BQU0sQ0FBQyxDQUFDO1FBQ3BCLE9BQU8sQ0FBQyxHQUFHLENBQUMsT0FBTyxDQUFDLENBQUM7UUFDckIsT0FBTyxDQUFDLEdBQUcsQ0FBQyxTQUFTLENBQUMsQ0FBQztLQUN4QjtJQUFDLE9BQU8sS0FBSyxFQUFFO1FBQ2QsT0FBTyxDQUFDLEdBQUcsQ0FBQyxLQUFLLENBQUMsT0FBTyxDQUFDLENBQUM7S0FDNUI7SUFFRCw0QkFBNEI7SUFDNUIsSUFBSTtRQUNGLE1BQU0sUUFBUSxHQUFHLE1BQU0sTUFBTSxDQUFDLE9BQU8sQ0FBQyxTQUFTLEVBQUUsQ0FBQztRQUNsRCxPQUFPLENBQUMsR0FBRyxDQUFDLFFBQVEsQ0FBQyxDQUFDO1FBQ3RCLE1BQU0sTUFBTSxHQUFHLFFBQVEsQ0FBQyxNQUFNLENBQUM7UUFDL0IsTUFBTSxVQUFVLEdBQUcsUUFBUSxDQUFDLElBQUksQ0FBQztRQUNqQyxPQUFPLENBQUMsR0FBRyxDQUFDLFFBQVEsQ0FBQyxDQUFDO1FBQ3RCLE9BQU8sQ0FBQyxHQUFHLENBQUMsTUFBTSxDQUFDLENBQUM7UUFDcEIsT0FBTyxDQUFDLEdBQUcsQ0FBQyxVQUFVLENBQUMsQ0FBQztLQUN6QjtJQUFDLE9BQU8sS0FBSyxFQUFFO1FBQ2QsT0FBTyxDQUFDLEdBQUcsQ0FBQyxLQUFLLENBQUMsT0FBTyxDQUFDLENBQUM7S0FDNUI7QUFDSCxDQUFDO0FBRUQsSUFBSSxFQUFFLENBQUMsSUFBSSxDQUFDLEdBQUcsRUFBRTtBQUNqQixDQUFDLENBQUMsQ0FBQyxLQUFLLENBQUMsQ0FBQyxLQUFLLEVBQUUsRUFBRTtJQUNqQixPQUFPLENBQUMsR0FBRyxDQUFDLEtBQUssQ0FBQyxPQUFPLENBQUMsQ0FBQztBQUM3QixDQUFDLENBQUMsQ0FBQyJ9