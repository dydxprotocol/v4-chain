"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.ValidatorClient = void 0;
const stargate_1 = require("@cosmjs/stargate");
const tendermint_rpc_1 = require("@cosmjs/tendermint-rpc");
const long_1 = __importDefault(require("long"));
const protobufjs_1 = __importDefault(require("protobufjs"));
const constants_1 = require("./constants");
const get_1 = require("./modules/get");
const post_1 = require("./modules/post");
const tendermintClient_1 = require("./modules/tendermintClient");
// Required for encoding and decoding queries that are of type Long.
// Must be done once but since the individal modules should be usable
// - must be set in each module that encounters encoding/decoding Longs.
// Reference: https://github.com/protobufjs/protobuf.js/issues/921
protobufjs_1.default.util.Long = long_1.default;
protobufjs_1.default.configure();
class ValidatorClient {
    /**
       * @description Connect to a validator client
       *
       * @returns The validator client
       */
    static async connect(config) {
        const client = new ValidatorClient(config);
        await client.initialize();
        return client;
    }
    constructor(config) {
        this.config = config;
    }
    /**
       * @description Get the query module, used for retrieving on-chain data.
       *
       * @returns The query module
       */
    get get() {
        return this._get;
    }
    /**
       * @description transaction module, used for sending transactions.
       *
       * @returns The transaction module
       */
    get post() {
        return this._post;
    }
    async initialize() {
        const tendermint37Client = await tendermint_rpc_1.Tendermint37Client.connect(this.config.restEndpoint);
        const tendermintClient = new tendermintClient_1.TendermintClient(tendermint37Client, {
            broadcastPollIntervalMs: constants_1.BROADCAST_POLL_INTERVAL_MS,
            broadcastTimeoutMs: constants_1.BROADCAST_TIMEOUT_MS,
        });
        const queryClient = stargate_1.QueryClient.withExtensions(tendermint37Client, stargate_1.setupTxExtension);
        this._get = new get_1.Get(tendermintClient, queryClient);
        this._post = new post_1.Post(this._get, this.config.chainId, this.config.denoms);
    }
}
exports.ValidatorClient = ValidatorClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidmFsaWRhdG9yLWNsaWVudC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3NyYy9jbGllbnRzL3ZhbGlkYXRvci1jbGllbnQudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7O0FBQUEsK0NBQThFO0FBQzlFLDJEQUE0RDtBQUM1RCxnREFBd0I7QUFDeEIsNERBQWtDO0FBRWxDLDJDQUFnRztBQUNoRyx1Q0FBb0M7QUFDcEMseUNBQXNDO0FBQ3RDLGlFQUE4RDtBQUU5RCxvRUFBb0U7QUFDcEUscUVBQXFFO0FBQ3JFLHdFQUF3RTtBQUN4RSxrRUFBa0U7QUFDbEUsb0JBQVEsQ0FBQyxJQUFJLENBQUMsSUFBSSxHQUFHLGNBQUksQ0FBQztBQUMxQixvQkFBUSxDQUFDLFNBQVMsRUFBRSxDQUFDO0FBRXJCLE1BQWEsZUFBZTtJQUsxQjs7OztTQUlLO0lBQ0wsTUFBTSxDQUFDLEtBQUssQ0FBQyxPQUFPLENBQUMsTUFBdUI7UUFDMUMsTUFBTSxNQUFNLEdBQUcsSUFBSSxlQUFlLENBQUMsTUFBTSxDQUFDLENBQUM7UUFDM0MsTUFBTSxNQUFNLENBQUMsVUFBVSxFQUFFLENBQUM7UUFDMUIsT0FBTyxNQUFNLENBQUM7SUFDaEIsQ0FBQztJQUVELFlBQ0UsTUFBdUI7UUFFdkIsSUFBSSxDQUFDLE1BQU0sR0FBRyxNQUFNLENBQUM7SUFDdkIsQ0FBQztJQUVEOzs7O1NBSUs7SUFDTCxJQUFJLEdBQUc7UUFDTCxPQUFPLElBQUksQ0FBQyxJQUFLLENBQUM7SUFDcEIsQ0FBQztJQUVEOzs7O1NBSUs7SUFDTCxJQUFJLElBQUk7UUFDTixPQUFPLElBQUksQ0FBQyxLQUFNLENBQUM7SUFDckIsQ0FBQztJQUVPLEtBQUssQ0FBQyxVQUFVO1FBQ3RCLE1BQU0sa0JBQWtCLEdBQXVCLE1BQU0sbUNBQWtCLENBQUMsT0FBTyxDQUM3RSxJQUFJLENBQUMsTUFBTSxDQUFDLFlBQVksQ0FDekIsQ0FBQztRQUVGLE1BQU0sZ0JBQWdCLEdBQUcsSUFBSSxtQ0FBZ0IsQ0FBQyxrQkFBa0IsRUFBRTtZQUNoRSx1QkFBdUIsRUFBRSxzQ0FBMEI7WUFDbkQsa0JBQWtCLEVBQUUsZ0NBQW9CO1NBQ3pDLENBQUMsQ0FBQztRQUNILE1BQU0sV0FBVyxHQUFnQyxzQkFBVyxDQUFDLGNBQWMsQ0FDekUsa0JBQWtCLEVBQ2xCLDJCQUFnQixDQUNqQixDQUFDO1FBQ0YsSUFBSSxDQUFDLElBQUksR0FBRyxJQUFJLFNBQUcsQ0FBQyxnQkFBZ0IsRUFBRSxXQUFXLENBQUMsQ0FBQztRQUNuRCxJQUFJLENBQUMsS0FBSyxHQUFHLElBQUksV0FBSSxDQUFDLElBQUksQ0FBQyxJQUFLLEVBQUUsSUFBSSxDQUFDLE1BQU0sQ0FBQyxPQUFPLEVBQUUsSUFBSSxDQUFDLE1BQU0sQ0FBQyxNQUFNLENBQUMsQ0FBQztJQUM3RSxDQUFDO0NBQ0Y7QUF4REQsMENBd0RDIn0=