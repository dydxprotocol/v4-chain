"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const constants_1 = require("../src/clients/constants");
const network_optimizer_1 = require("../src/network_optimizer");
async function testNodes() {
    // all valid endpoints
    try {
        const optimizer = new network_optimizer_1.NetworkOptimizer();
        const endpoints = [
            'https://validator.v4testnet1.dydx.exchange',
            'https://dydx-testnet.nodefleet.org',
            'https://dydx-testnet-archive.allthatnode.com:26657/XZvMM41hESf8PJrEQiTzbCOMVyFca79R',
        ];
        const optimal = await optimizer.findOptimalNode(endpoints, constants_1.TESTNET_CHAIN_ID);
        console.log(optimal);
    }
    catch (error) {
        console.log(error.message);
    }
    // one invalid endpoint
    try {
        const optimizer = new network_optimizer_1.NetworkOptimizer();
        const endpoints = [
            'https://validator.v4testnet1.dydx.exchange',
            'https://dydx-testnet.nodefleet.org',
            'https://dydx-testnet-archive.allthatnode.com:26657/XZvMM41hESf8PJrEQiTzbCOMVyFca79R',
            'https://example.com',
        ];
        const optimal = await optimizer.findOptimalNode(endpoints, constants_1.TESTNET_CHAIN_ID);
        console.log(optimal);
    }
    catch (error) {
        console.log(error.message);
    }
    // all invalid endpoints
    try {
        const optimizer = new network_optimizer_1.NetworkOptimizer();
        const endpoints = [
            'https://example.com',
            'https://example.org',
        ];
        const optimal = await optimizer.findOptimalNode(endpoints, constants_1.TESTNET_CHAIN_ID);
        console.log(optimal);
    }
    catch (error) {
        console.log(error.message);
    }
}
async function testIndexers() {
    // all valid endpoints
    try {
        const optimizer = new network_optimizer_1.NetworkOptimizer();
        const endpoints = [
            'https://indexer.v4testnet2.dydx.exchange',
        ];
        const optimal = await optimizer.findOptimalIndexer(endpoints);
        console.log(optimal);
    }
    catch (error) {
        console.log(error.message);
    }
    // all invalid endpoints
    try {
        const optimizer = new network_optimizer_1.NetworkOptimizer();
        const endpoints = [
            'https://example.com',
            'https://example.org',
        ];
        const optimal = await optimizer.findOptimalIndexer(endpoints);
        console.log(optimal);
    }
    catch (error) {
        console.log(error.message);
    }
}
testNodes().catch(console.log);
testIndexers().catch(console.log);
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoib3B0aW1hbF9ub2RlLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vZXhhbXBsZXMvb3B0aW1hbF9ub2RlLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7O0FBQUEsd0RBRWtDO0FBQ2xDLGdFQUE0RDtBQUU1RCxLQUFLLFVBQVUsU0FBUztJQUN0QixzQkFBc0I7SUFDdEIsSUFBSSxDQUFDO1FBQ0gsTUFBTSxTQUFTLEdBQUcsSUFBSSxvQ0FBZ0IsRUFBRSxDQUFDO1FBQ3pDLE1BQU0sU0FBUyxHQUFHO1lBQ2hCLDRDQUE0QztZQUM1QyxvQ0FBb0M7WUFDcEMscUZBQXFGO1NBQ3RGLENBQUM7UUFDRixNQUFNLE9BQU8sR0FBRyxNQUFNLFNBQVMsQ0FBQyxlQUFlLENBQUMsU0FBUyxFQUFFLDRCQUFnQixDQUFDLENBQUM7UUFDN0UsT0FBTyxDQUFDLEdBQUcsQ0FBQyxPQUFPLENBQUMsQ0FBQztJQUN2QixDQUFDO0lBQUMsT0FBTyxLQUFLLEVBQUUsQ0FBQztRQUNmLE9BQU8sQ0FBQyxHQUFHLENBQUMsS0FBSyxDQUFDLE9BQU8sQ0FBQyxDQUFDO0lBQzdCLENBQUM7SUFFRCx1QkFBdUI7SUFDdkIsSUFBSSxDQUFDO1FBQ0gsTUFBTSxTQUFTLEdBQUcsSUFBSSxvQ0FBZ0IsRUFBRSxDQUFDO1FBQ3pDLE1BQU0sU0FBUyxHQUFHO1lBQ2hCLDRDQUE0QztZQUM1QyxvQ0FBb0M7WUFDcEMscUZBQXFGO1lBQ3JGLHFCQUFxQjtTQUN0QixDQUFDO1FBQ0YsTUFBTSxPQUFPLEdBQUcsTUFBTSxTQUFTLENBQUMsZUFBZSxDQUFDLFNBQVMsRUFBRSw0QkFBZ0IsQ0FBQyxDQUFDO1FBQzdFLE9BQU8sQ0FBQyxHQUFHLENBQUMsT0FBTyxDQUFDLENBQUM7SUFDdkIsQ0FBQztJQUFDLE9BQU8sS0FBSyxFQUFFLENBQUM7UUFDZixPQUFPLENBQUMsR0FBRyxDQUFDLEtBQUssQ0FBQyxPQUFPLENBQUMsQ0FBQztJQUM3QixDQUFDO0lBRUQsd0JBQXdCO0lBRXhCLElBQUksQ0FBQztRQUNILE1BQU0sU0FBUyxHQUFHLElBQUksb0NBQWdCLEVBQUUsQ0FBQztRQUN6QyxNQUFNLFNBQVMsR0FBRztZQUNoQixxQkFBcUI7WUFDckIscUJBQXFCO1NBQ3RCLENBQUM7UUFDRixNQUFNLE9BQU8sR0FBRyxNQUFNLFNBQVMsQ0FBQyxlQUFlLENBQUMsU0FBUyxFQUFFLDRCQUFnQixDQUFDLENBQUM7UUFDN0UsT0FBTyxDQUFDLEdBQUcsQ0FBQyxPQUFPLENBQUMsQ0FBQztJQUN2QixDQUFDO0lBQUMsT0FBTyxLQUFLLEVBQUUsQ0FBQztRQUNmLE9BQU8sQ0FBQyxHQUFHLENBQUMsS0FBSyxDQUFDLE9BQU8sQ0FBQyxDQUFDO0lBQzdCLENBQUM7QUFDSCxDQUFDO0FBRUQsS0FBSyxVQUFVLFlBQVk7SUFDekIsc0JBQXNCO0lBQ3RCLElBQUksQ0FBQztRQUNILE1BQU0sU0FBUyxHQUFHLElBQUksb0NBQWdCLEVBQUUsQ0FBQztRQUN6QyxNQUFNLFNBQVMsR0FBRztZQUNoQiwwQ0FBMEM7U0FDM0MsQ0FBQztRQUNGLE1BQU0sT0FBTyxHQUFHLE1BQU0sU0FBUyxDQUFDLGtCQUFrQixDQUFDLFNBQVMsQ0FBQyxDQUFDO1FBQzlELE9BQU8sQ0FBQyxHQUFHLENBQUMsT0FBTyxDQUFDLENBQUM7SUFDdkIsQ0FBQztJQUFDLE9BQU8sS0FBSyxFQUFFLENBQUM7UUFDZixPQUFPLENBQUMsR0FBRyxDQUFDLEtBQUssQ0FBQyxPQUFPLENBQUMsQ0FBQztJQUM3QixDQUFDO0lBRUQsd0JBQXdCO0lBRXhCLElBQUksQ0FBQztRQUNILE1BQU0sU0FBUyxHQUFHLElBQUksb0NBQWdCLEVBQUUsQ0FBQztRQUN6QyxNQUFNLFNBQVMsR0FBRztZQUNoQixxQkFBcUI7WUFDckIscUJBQXFCO1NBQ3RCLENBQUM7UUFDRixNQUFNLE9BQU8sR0FBRyxNQUFNLFNBQVMsQ0FBQyxrQkFBa0IsQ0FBQyxTQUFTLENBQUMsQ0FBQztRQUM5RCxPQUFPLENBQUMsR0FBRyxDQUFDLE9BQU8sQ0FBQyxDQUFDO0lBQ3ZCLENBQUM7SUFBQyxPQUFPLEtBQUssRUFBRSxDQUFDO1FBQ2YsT0FBTyxDQUFDLEdBQUcsQ0FBQyxLQUFLLENBQUMsT0FBTyxDQUFDLENBQUM7SUFDN0IsQ0FBQztBQUNILENBQUM7QUFFRCxTQUFTLEVBQUUsQ0FBQyxLQUFLLENBQUMsT0FBTyxDQUFDLEdBQUcsQ0FBQyxDQUFDO0FBQy9CLFlBQVksRUFBRSxDQUFDLEtBQUssQ0FBQyxPQUFPLENBQUMsR0FBRyxDQUFDLENBQUMifQ==