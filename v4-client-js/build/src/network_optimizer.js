"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.NetworkOptimizer = exports.isTruthy = void 0;
const indexer_client_1 = require("./clients/indexer-client");
const validator_client_1 = require("./clients/validator-client");
const helpers_1 = require("./lib/helpers");
const types_1 = require("./types");
class PingResponse {
    constructor(height) {
        this.height = height;
        this.responseTime = new Date();
    }
}
const isTruthy = (n) => Boolean(n);
exports.isTruthy = isTruthy;
class NetworkOptimizer {
    async validatorClients(endpointUrls, chainId) {
        return (await Promise.all(endpointUrls.map((endpointUrl) => validator_client_1.ValidatorClient.connect(new types_1.ValidatorConfig(endpointUrl, chainId, {
            CHAINTOKEN_DENOM: 'placeholder',
            CHAINTOKEN_DECIMALS: 18,
            TDAI_DENOM: 'utdai',
            TDAI_DECIMALS: 6,
        }))
            .catch((_) => undefined)))).filter(exports.isTruthy);
    }
    indexerClients(endpointUrls) {
        return endpointUrls.map((endpointUrl) => new indexer_client_1.IndexerClient(
        // socket is not used for finding optimal indexer, but required as a parameter to the config
        new types_1.IndexerConfig(endpointUrl, endpointUrl.replace('https://', 'wss://').replace('http://', 'ws://')))).filter(exports.isTruthy);
    }
    async findOptimalNode(endpointUrls, chainId) {
        if (endpointUrls.length === 0) {
            const errorResponse = {
                error: {
                    message: 'No nodes provided',
                },
            };
            return (0, helpers_1.encodeJson)(errorResponse);
        }
        const clients = await this.validatorClients(endpointUrls, chainId);
        const responses = (await Promise.all(clients
            .map(async (client) => {
            const block = await client.get.latestBlock();
            const response = new PingResponse(block.header.height);
            return {
                endpoint: client.config.restEndpoint,
                height: response.height,
                time: response.responseTime.getTime(),
            };
        })
            .map((promise) => promise.catch((_) => undefined)))).filter(exports.isTruthy);
        if (responses.length === 0) {
            throw new Error('Could not connect to endpoints');
        }
        const maxHeight = Math.max(...responses.map(({ height }) => height));
        return responses
            // Only consider nodes at `maxHeight` or `maxHeight - 1`
            .filter(({ height }) => height === maxHeight || height === maxHeight - 1)
            // Return the endpoint with the fastest response time
            .sort((a, b) => a.time - b.time)[0]
            .endpoint;
    }
    async findOptimalIndexer(endpointUrls) {
        if (endpointUrls.length === 0) {
            const errorResponse = {
                error: {
                    message: 'No URL provided',
                },
            };
            return (0, helpers_1.encodeJson)(errorResponse);
        }
        const clients = this.indexerClients(endpointUrls);
        const responses = (await Promise.all(clients
            .map(async (client) => {
            const block = await client.utility.getHeight();
            const response = new PingResponse(+block.height);
            return {
                endpoint: client.config.restEndpoint,
                height: response.height,
                time: response.responseTime.getTime(),
            };
        })
            .map((promise) => promise.catch((_) => undefined)))).filter(exports.isTruthy);
        if (responses.length === 0) {
            throw new Error('Could not connect to endpoints');
        }
        const maxHeight = Math.max(...responses.map(({ height }) => height));
        return responses
            // Only consider nodes at `maxHeight` or `maxHeight - 1`
            .filter(({ height }) => height === maxHeight || height === maxHeight - 1)
            // Return the endpoint with the fastest response time
            .sort((a, b) => a.time - b.time)[0]
            .endpoint;
    }
}
exports.NetworkOptimizer = NetworkOptimizer;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibmV0d29ya19vcHRpbWl6ZXIuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi9zcmMvbmV0d29ya19vcHRpbWl6ZXIudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBQUEsNkRBQXlEO0FBQ3pELGlFQUE2RDtBQUM3RCwyQ0FBMkM7QUFDM0MsbUNBQXlEO0FBRXpELE1BQU0sWUFBWTtJQUtkLFlBQ0UsTUFBYztRQUVkLElBQUksQ0FBQyxNQUFNLEdBQUcsTUFBTSxDQUFDO1FBQ3JCLElBQUksQ0FBQyxZQUFZLEdBQUcsSUFBSSxJQUFJLEVBQUUsQ0FBQztJQUNqQyxDQUFDO0NBQ0o7QUFFTSxNQUFNLFFBQVEsR0FBRyxDQUFJLENBQW9DLEVBQVUsRUFBRSxDQUFDLE9BQU8sQ0FBQyxDQUFDLENBQUMsQ0FBQztBQUEzRSxRQUFBLFFBQVEsWUFBbUU7QUFFeEYsTUFBYSxnQkFBZ0I7SUFDbkIsS0FBSyxDQUFDLGdCQUFnQixDQUM1QixZQUFzQixFQUN0QixPQUFlO1FBRWYsT0FBTyxDQUFDLE1BQU0sT0FBTyxDQUFDLEdBQUcsQ0FDdkIsWUFBWSxDQUFDLEdBQUcsQ0FBQyxDQUFDLFdBQVcsRUFBRSxFQUFFLENBQUMsa0NBQWUsQ0FBQyxPQUFPLENBQ3ZELElBQUksdUJBQWUsQ0FBQyxXQUFXLEVBQUUsT0FBTyxFQUFFO1lBQ3hDLGdCQUFnQixFQUFFLGFBQWE7WUFDL0IsbUJBQW1CLEVBQUUsRUFBRTtZQUN2QixVQUFVLEVBQUUsT0FBTztZQUNuQixhQUFhLEVBQUUsQ0FBQztTQUNqQixDQUFDLENBQUM7YUFDRixLQUFLLENBQUMsQ0FBQyxDQUFDLEVBQUUsRUFBRSxDQUFDLFNBQVMsQ0FBQyxDQUN6QixDQUNGLENBQUMsQ0FBQyxNQUFNLENBQUMsZ0JBQVEsQ0FBQyxDQUFDO0lBQ3RCLENBQUM7SUFFTyxjQUFjLENBQ3BCLFlBQXNCO1FBRXRCLE9BQU8sWUFBWSxDQUFDLEdBQUcsQ0FBQyxDQUFDLFdBQVcsRUFBRSxFQUFFLENBQUMs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