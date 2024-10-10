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