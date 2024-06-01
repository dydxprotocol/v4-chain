"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.SocketClient = exports.CandlesResolution = exports.IncomingMessageTypes = void 0;
const ws_1 = __importDefault(require("ws"));
var OutgoingMessageTypes;
(function (OutgoingMessageTypes) {
    OutgoingMessageTypes["PING"] = "ping";
    OutgoingMessageTypes["SUBSCRIBE"] = "subscribe";
    OutgoingMessageTypes["UNSUBSCRIBE"] = "unsubscribe";
})(OutgoingMessageTypes || (OutgoingMessageTypes = {}));
var SocketChannels;
(function (SocketChannels) {
    SocketChannels["SUBACCOUNTS"] = "v4_subaccounts";
    SocketChannels["ORDERBOOK"] = "v4_orderbook";
    SocketChannels["TRADES"] = "v4_trades";
    SocketChannels["MARKETS"] = "v4_markets";
    SocketChannels["CANDLES"] = "v4_candles";
})(SocketChannels || (SocketChannels = {}));
var IncomingMessageTypes;
(function (IncomingMessageTypes) {
    IncomingMessageTypes["CONNECTED"] = "connected";
    IncomingMessageTypes["SUBSCRIBED"] = "subscribed";
    IncomingMessageTypes["ERROR"] = "error";
    IncomingMessageTypes["CHANNEL_DATA"] = "channel_data";
    IncomingMessageTypes["CHANNEL_BATCH_DATA"] = "channel_batch_data";
    IncomingMessageTypes["PONG"] = "pong";
})(IncomingMessageTypes = exports.IncomingMessageTypes || (exports.IncomingMessageTypes = {}));
var CandlesResolution;
(function (CandlesResolution) {
    CandlesResolution["ONE_MINUTE"] = "1MIN";
    CandlesResolution["FIVE_MINUTES"] = "5MINS";
    CandlesResolution["FIFTEEN_MINUTES"] = "15MINS";
    CandlesResolution["THIRTY_MINUTES"] = "30MINS";
    CandlesResolution["ONE_HOUR"] = "1HOUR";
    CandlesResolution["FOUR_HOURS"] = "4HOURS";
    CandlesResolution["ONE_DAY"] = "1DAY";
})(CandlesResolution = exports.CandlesResolution || (exports.CandlesResolution = {}));
class SocketClient {
    constructor(config, onOpenCallback, onCloseCallback, onMessageCallback) {
        this.pingInterval = 30000;
        this.lastMessageTime = Date.now();
        this.url = config.websocketEndpoint;
        this.onOpenCallback = onOpenCallback;
        this.onCloseCallback = onCloseCallback;
        this.onMessageCallback = onMessageCallback;
    }
    connect() {
        this.ws = new ws_1.default(this.url);
        this.ws.addEventListener('open', this.handleOpen.bind(this));
        this.ws.addEventListener('close', this.handleClose.bind(this));
        this.ws.addEventListener('message', this.handleMessage.bind(this));
    }
    /**
     * @description Close the websocket connection.
     *
     */
    close() {
        var _a;
        (_a = this.ws) === null || _a === void 0 ? void 0 : _a.close();
        this.ws = undefined;
    }
    /**
     * @description Send data to the websocket connection.
     *
     */
    send(data) {
        var _a;
        (_a = this.ws) === null || _a === void 0 ? void 0 : _a.send(data);
    }
    handleOpen() {
        if (this.onOpenCallback) {
            this.onOpenCallback();
        }
        this.restartPingInterval();
    }
    handleClose() {
        if (this.onCloseCallback) {
            this.onCloseCallback();
        }
        clearInterval(this.pingIntervalId);
    }
    handleMessage(event) {
        if (event.data === 'PING') {
            this.send('PONG');
        }
        else {
            this.lastMessageTime = Date.now();
            if (this.onMessageCallback) {
                this.onMessageCallback(event);
            }
        }
        this.restartPingInterval();
    }
    restartPingInterval() {
        clearInterval(this.pingIntervalId);
        this.pingIntervalId = setInterval(() => {
            const elapsedTime = Date.now() - this.lastMessageTime;
            if (elapsedTime > this.pingInterval) {
                this.send('PING');
            }
        }, this.pingInterval);
    }
    /**
     * @description Set callback when the socket is opened.
     *
     */
    set onOpen(callback) {
        this.onOpenCallback = callback;
    }
    /**
     * @description Set callback when the socket is closed.
     *
     */
    set onClose(callback) {
        this.onCloseCallback = callback;
    }
    /**
     * @description Set callback when the socket receives a message.
     *
     */
    set onMessage(callback) {
        this.onMessageCallback = callback;
    }
    /**
     * @description Send a subscribe message to the websocket connection.
     *
     */
    subscribe(channel, params) {
        const message = {
            type: OutgoingMessageTypes.SUBSCRIBE,
            channel,
            ...params,
        };
        this.send(JSON.stringify(message));
    }
    /**
     * @description Send an unsubscribe message to the websocket connection.
     *
     */
    unsubscribe(channel, params) {
        const message = {
            type: OutgoingMessageTypes.UNSUBSCRIBE,
            channel,
            ...params,
        };
        this.send(JSON.stringify(message));
    }
    /**
     * @description Subscribe to markets channel.
     *
     */
    subscribeToMarkets() {
        const channel = SocketChannels.MARKETS;
        const params = {
            batched: true,
        };
        this.subscribe(channel, params);
    }
    /**
     * @description Unsubscribe from markets channel
     *
     */
    unsubscribeFromMarkets() {
        const channel = SocketChannels.MARKETS;
        this.unsubscribe(channel);
    }
    /**
     * @description Subscribe to trade channel
     * for a specific market.
     *
     */
    subscribeToTrades(market) {
        const channel = SocketChannels.TRADES;
        const params = {
            id: market,
            batched: true,
        };
        this.subscribe(channel, params);
    }
    /**
     * @description Unscribed from trade channel
     * for a specific market.
     *
     */
    unsubscribeFromTrades(market) {
        const channel = SocketChannels.TRADES;
        const params = {
            id: market,
        };
        this.unsubscribe(channel, params);
    }
    /**
     * @description Subscribe to orderbook channel
     * for a specific market.
     *
     */
    subscribeToOrderbook(market) {
        const channel = SocketChannels.ORDERBOOK;
        const params = {
            id: market,
            batched: true,
        };
        this.subscribe(channel, params);
    }
    /**
     * @description Unsubscribe from orderbook channel
     * for a specific market.
     */
    unsubscribeFromOrderbook(market) {
        const channel = SocketChannels.ORDERBOOK;
        const params = {
            id: market,
        };
        this.unsubscribe(channel, params);
    }
    /**
     * @description Subscribe to candles channel
     * for a specific market and resolution.
     *
     */
    subscribeToCandles(market, resolution) {
        const channel = SocketChannels.CANDLES;
        const params = {
            id: `${market}/${resolution}`,
            batched: true,
        };
        this.subscribe(channel, params);
    }
    /**
     * @description Unsubscribe from candles channel
     * for a specific market and resolution.
     */
    unsubscribeFromCandles(market, resolution) {
        const channel = SocketChannels.CANDLES;
        const params = {
            id: `${market}/${resolution}`,
        };
        this.unsubscribe(channel, params);
    }
    /**
     * @description Subscribe to subaccount channel
     * for a specific address and subaccount number.
     */
    subscribeToSubaccount(address, subaccountNumber) {
        const channel = SocketChannels.SUBACCOUNTS;
        const subaccountId = `${address}/${subaccountNumber}`;
        const params = {
            id: subaccountId,
        };
        this.subscribe(channel, params);
    }
    /**
     * @description Unsubscribe from subaccount channel
     * for a specific address and subaccount number.
     *
     */
    unsubscribeFromSubaccount(address, subaccountNumber) {
        const channel = SocketChannels.SUBACCOUNTS;
        const subaccountId = `${address}/${subaccountNumber}`;
        const params = {
            id: subaccountId,
        };
        this.unsubscribe(channel, params);
    }
}
exports.SocketClient = SocketClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic29ja2V0LWNsaWVudC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3NyYy9jbGllbnRzL3NvY2tldC1jbGllbnQudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7O0FBQUEsNENBQTZDO0FBSTdDLElBQUssb0JBSUo7QUFKRCxXQUFLLG9CQUFvQjtJQUN2QixxQ0FBYSxDQUFBO0lBQ2IsK0NBQXVCLENBQUE7SUFDdkIsbURBQTJCLENBQUE7QUFDN0IsQ0FBQyxFQUpJLG9CQUFvQixLQUFwQixvQkFBb0IsUUFJeEI7QUFFRCxJQUFLLGNBTUo7QUFORCxXQUFLLGNBQWM7SUFDakIsZ0RBQThCLENBQUE7SUFDOUIsNENBQTBCLENBQUE7SUFDMUIsc0NBQW9CLENBQUE7SUFDcEIsd0NBQXNCLENBQUE7SUFDdEIsd0NBQXNCLENBQUE7QUFDeEIsQ0FBQyxFQU5JLGNBQWMsS0FBZCxjQUFjLFFBTWxCO0FBRUQsSUFBWSxvQkFPWDtBQVBELFdBQVksb0JBQW9CO0lBQzlCLCtDQUF1QixDQUFBO0lBQ3ZCLGlEQUF5QixDQUFBO0lBQ3pCLHVDQUFlLENBQUE7SUFDZixxREFBNkIsQ0FBQTtJQUM3QixpRUFBeUMsQ0FBQTtJQUN6QyxxQ0FBYSxDQUFBO0FBQ2YsQ0FBQyxFQVBXLG9CQUFvQixHQUFwQiw0QkFBb0IsS0FBcEIsNEJBQW9CLFFBTy9CO0FBRUQsSUFBWSxpQkFRWDtBQVJELFdBQVksaUJBQWlCO0lBQzNCLHdDQUFtQixDQUFBO0lBQ25CLDJDQUFzQixDQUFBO0lBQ3RCLCtDQUEwQixDQUFBO0lBQzFCLDhDQUF5QixDQUFBO0lBQ3pCLHVDQUFrQixDQUFBO0lBQ2xCLDBDQUFxQixDQUFBO0lBQ3JCLHFDQUFnQixDQUFBO0FBQ2xCLENBQUMsRUFSVyxpQkFBaUIsR0FBakIseUJBQWlCLEtBQWpCLHlCQUFpQixRQVE1QjtBQUVELE1BQWEsWUFBWTtJQVVyQixZQUNFLE1BQXFCLEVBQ3JCLGNBQTBCLEVBQzFCLGVBQTBCLEVBQzFCLGlCQUErQztRQVJ6QyxpQkFBWSxHQUFXLEtBQU0sQ0FBQztRQUM5QixvQkFBZSxHQUFXLElBQUksQ0FBQyxHQUFHLEVBQUUsQ0FBQztRQVMzQyxJQUFJLENBQUMsR0FBRyxHQUFHLE1BQU0sQ0FBQyxpQkFBaUIsQ0FBQztRQUNwQyxJQUFJLENBQUMsY0FBYyxHQUFHLGNBQWMsQ0FBQztRQUNyQyxJQUFJLENBQUMsZUFBZSxHQUFHLGVBQWUsQ0FBQztRQUN2QyxJQUFJLENBQUMsaUJBQWlCLEdBQUcsaUJBQWlCLENBQUM7SUFDN0MsQ0FBQztJQUVELE9BQU87UUFDTCxJQUFJLENBQUMsRUFBRSxHQUFHLElBQUksWUFBUyxDQUFDLElBQUksQ0FBQyxHQUFHLENBQUMsQ0FBQztRQUNsQyxJQUFJLENBQUMsRUFBRSxDQUFDLGdCQUFnQixDQUFDLE1BQU0sRUFBRSxJQUFJLENBQUMsVUFBVSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQyxDQUFDO1FBQzdELElBQUksQ0FBQyxFQUFFLENBQUMsZ0JBQWdCLENBQUMsT0FBTyxFQUFFLElBQUksQ0FBQyxXQUFXLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUM7UUFDL0QsSUFBSSxDQUFDLEVBQUUsQ0FBQyxnQkFBZ0IsQ0FBQyxTQUFTLEVBQUUsSUFBSSxDQUFDLGFBQWEsQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLENBQUMsQ0FBQztJQUNyRSxDQUFDO0lBRUQ7OztPQUdHO0lBQ0gsS0FBSzs7UUFDSCxNQUFBLElBQUksQ0FBQyxFQUFFLDBDQUFFLEtBQUssRUFBRSxDQUFDO1FBQ2pCLElBQUksQ0FBQyxFQUFFLEdBQUcsU0FBUyxDQUFDO0lBQ3RCLENBQUM7SUFFRDs7O09BR0c7SUFDSCxJQUFJLENBQUMsSUFBWTs7UUFDZixNQUFBLElBQUksQ0FBQyxFQUFFLDBDQUFFLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQztJQUN0QixDQUFDO0lBRU8sVUFBVTtRQUNoQixJQUFJLElBQUksQ0FBQyxjQUFjLEVBQUU7WUFDdkIsSUFBSSxDQUFDLGNBQWMsRUFBRSxDQUFDO1NBQ3ZCO1FBQ0QsSUFBSSxDQUFDLG1CQUFtQixFQUFFLENBQUM7SUFDN0IsQ0FBQztJQUVPLFdBQVc7UUFDakIsSUFBSSxJQUFJLENBQUMsZUFBZSxFQUFFO1lBQ3hCLElBQUksQ0FBQyxlQUFlLEVBQUUsQ0FBQztTQUN4QjtRQUNELGFBQWEsQ0FBQyxJQUFJLENBQUMsY0FBYyxDQUFDLENBQUM7SUFDckMsQ0FBQztJQUVPLGFBQWEsQ0FBQyxLQUFtQjtRQUN2QyxJQUFJLEtBQUssQ0FBQyxJQUFJLEtBQUssTUFBTSxFQUFFO1lBQ3pCLElBQUksQ0FBQyxJQUFJLENBQUMsTUFBTSxDQUFDLENBQUM7U0FDbkI7YUFBTTtZQUNMLElBQUksQ0FBQyxlQUFlLEdBQUcsSUFBSSxDQUFDLEdBQUcsRUFBRSxDQUFDO1lBQ2xDLElBQUksSUFBSSxDQUFDLGlCQUFpQixFQUFFO2dCQUMxQixJQUFJLENBQUMsaUJBQWlCLENBQUMsS0FBSyxDQUFDLENBQUM7YUFDL0I7U0FDRjtRQUNELElBQUksQ0FBQyxtQkFBbUIsRUFBRSxDQUFDO0lBQzdCLENBQUM7SUFFTyxtQkFBbUI7UUFDekIsYUFBYSxDQUFDLElBQUksQ0FBQyxjQUFjLENBQUMsQ0FBQztRQUNuQyxJQUFJLENBQUMsY0FBYyxHQUFHLFdBQVcsQ0FBQyxHQUFHLEVBQUU7WUFDckMsTUFBTSxXQUFXLEdBQUcsSUFBSSxDQUFDLEdBQUcsRUFBRSxHQUFHLElBQUksQ0FBQyxlQUFlLENBQUM7WUFDdEQsSUFBSSxXQUFXLEdBQUcsSUFBSSxDQUFDLFlBQVksRUFBRTtnQkFDbkMsSUFBSSxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUMsQ0FBQzthQUNuQjtRQUNILENBQUMsRUFBRSxJQUFJLENBQUMsWUFBWSxDQUFDLENBQUM7SUFDeEIsQ0FBQztJQUVEOzs7T0FHRztJQUNILElBQUksTUFBTSxDQUFDLFFBQW9CO1FBQzdCLElBQUksQ0FBQyxjQUFjLEdBQUcsUUFBUSxDQUFDO0lBQ2pDLENBQUM7SUFFRDs7O09BR0c7SUFDSCxJQUFJLE9BQU8sQ0FBQyxRQUFvQjtRQUM5QixJQUFJLENBQUMsZUFBZSxHQUFHLFFBQVEsQ0FBQztJQUNsQyxDQUFDO0lBRUQ7OztPQUdHO0lBQ0gsSUFBSSxTQUFTLENBQUMsUUFBdUM7UUFDbkQsSUFBSSxDQUFDLGlCQUFpQixHQUFHLFFBQVEsQ0FBQztJQUNwQyxDQUFDO0lBRUQ7OztPQUdHO0lBQ0gsU0FBUyxDQUFDLE9BQWUsRUFBRSxNQUFlO1FBQ3hDLE1BQU0sT0FBTyxHQUFHO1lBQ2QsSUFBSSxFQUFFLG9CQUFvQixDQUFDLFNBQVM7WUFDcEMsT0FBTztZQUNQLEdBQUcsTUFBTTtTQUNWLENBQUM7UUFDRixJQUFJLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxTQUFTLENBQUMsT0FBTyxDQUFDLENBQUMsQ0FBQztJQUNyQyxDQUFDO0lBRUQ7OztPQUdHO0lBQ0gsV0FBVyxDQUFDLE9BQWUsRUFBRSxNQUFlO1FBQzFDLE1BQU0sT0FBTyxHQUFHO1lBQ2QsSUFBSSxFQUFFLG9CQUFvQixDQUFDLFdBQVc7WUFDdEMsT0FBTztZQUNQLEdBQUcsTUFBTTtTQUNWLENBQUM7UUFDRixJQUFJLENBQUMsSUFBSSxDQUFDLElBQUksQ0FBQyxTQUFTLENBQUMsT0FBTyxDQUFDLENBQUMsQ0FBQztJQUNyQyxDQUFDO0lBRUQ7OztPQUdHO0lBQ0gsa0JBQWtCO1FBQ2hCLE1BQU0sT0FBTyxHQUFHLGNBQWMsQ0FBQyxPQUFPLENBQUM7UUFDdkMsTUFBTSxNQUFNLEdBQUc7WUFDYixPQUFPLEVBQUUsSUFBSTtTQUNkLENBQUM7UUFDRixJQUFJLENBQUMsU0FBUyxDQUFDLE9BQU8sRUFBRSxNQUFNLENBQUMsQ0FBQztJQUNsQyxDQUFDO0lBRUQ7OztPQUdHO0lBQ0gsc0JBQXNCO1FBQ3BCLE1BQU0sT0FBTyxHQUFHLGNBQWMsQ0FBQyxPQUFPLENBQUM7UUFDdkMsSUFBSSxDQUFDLFdBQVcsQ0FBQyxPQUFPLENBQUMsQ0FBQztJQUM1QixDQUFDO0lBRUQ7Ozs7T0FJRztJQUNILGlCQUFpQixDQUFDLE1BQWM7UUFDOUIsTUFBTSxPQUFPLEdBQUcsY0FBYyxDQUFDLE1BQU0sQ0FBQztRQUN0QyxNQUFNLE1BQU0sR0FBRztZQUNiLEVBQUUsRUFBRSxNQUFNO1lBQ1YsT0FBTyxFQUFFLElBQUk7U0FDZCxDQUFDO1FBQ0YsSUFBSSxDQUFDLFNBQVMsQ0FBQyxPQUFPLEVBQUUsTUFBTSxDQUFDLENBQUM7SUFDbEMsQ0FBQztJQUVEOzs7O09BSUc7SUFDSCxxQkFBcUIsQ0FBQyxNQUFjO1FBQ2xDLE1BQU0sT0FBTyxHQUFHLGNBQWMsQ0FBQyxNQUFNLENBQUM7UUFDdEMsTUFBTSxNQUFNLEdBQUc7WUFDYixFQUFFLEVBQUUsTUFBTTtTQUNYLENBQUM7UUFDRixJQUFJLENBQUMsV0FBVyxDQUFDLE9BQU8sRUFBRSxNQUFNLENBQUMsQ0FBQztJQUNwQyxDQUFDO0lBRUQ7Ozs7T0FJRztJQUNILG9CQUFvQixDQUFDLE1BQWM7UUFDakMsTUFBTSxPQUFPLEdBQUcsY0FBYyxDQUFDLFNBQVMsQ0FBQztRQUN6QyxNQUFNLE1BQU0sR0FBRztZQUNiLEVBQUUsRUFBRSxNQUFNO1lBQ1YsT0FBTyxFQUFFLElBQUk7U0FDZCxDQUFDO1FBQ0YsSUFBSSxDQUFDLFNBQVMsQ0FBQyxPQUFPLEVBQUUsTUFBTSxDQUFDLENBQUM7SUFDbEMsQ0FBQztJQUVEOzs7T0FHRztJQUNILHdCQUF3QixDQUFDLE1BQWM7UUFDckMsTUFBTSxPQUFPLEdBQUcsY0FBYyxDQUFDLFNBQVMsQ0FBQztRQUN6QyxNQUFNLE1BQU0sR0FBRztZQUNiLEVBQUUsRUFBRSxNQUFNO1NBQ1gsQ0FBQztRQUNGLElBQUksQ0FBQyxXQUFXLENBQUMsT0FBTyxFQUFFLE1BQU0sQ0FBQyxDQUFDO0lBQ3BDLENBQUM7SUFFRDs7OztPQUlHO0lBQ0gsa0JBQWtCLENBQUMsTUFBYyxFQUFFLFVBQTZCO1FBQzlELE1BQU0sT0FBTyxHQUFHLGNBQWMsQ0FBQyxPQUFPLENBQUM7UUFDdkMsTUFBTSxNQUFNLEdBQUc7WUFDYixFQUFFLEVBQUUsR0FBRyxNQUFNLElBQUksVUFBVSxFQUFFO1lBQzdCLE9BQU8sRUFBRSxJQUFJO1NBQ2QsQ0FBQztRQUNGLElBQUksQ0FBQyxTQUFTLENBQUMsT0FBTyxFQUFFLE1BQU0sQ0FBQyxDQUFDO0lBQ2xDLENBQUM7SUFFRDs7O09BR0c7SUFDSCxzQkFBc0IsQ0FBQyxNQUFjLEVBQUUsVUFBNkI7UUFDbEUsTUFBTSxPQUFPLEdBQUcsY0FBYyxDQUFDLE9BQU8sQ0FBQztRQUN2QyxNQUFNLE1BQU0sR0FBRztZQUNiLEVBQUUsRUFBRSxHQUFHLE1BQU0sSUFBSSxVQUFVLEVBQUU7U0FDOUIsQ0FBQztRQUNGLElBQUksQ0FBQyxXQUFXLENBQUMsT0FBTyxFQUFFLE1BQU0sQ0FBQyxDQUFDO0lBQ3BDLENBQUM7SUFFRDs7O09BR0c7SUFDSCxxQkFBcUIsQ0FBQyxPQUFlLEVBQUUsZ0JBQXdCO1FBQzdELE1BQU0sT0FBTyxHQUFHLGNBQWMsQ0FBQyxXQUFXLENBQUM7UUFDM0MsTUFBTSxZQUFZLEdBQUcsR0FBRyxPQUFPLElBQUksZ0JBQWdCLEVBQUUsQ0FBQztRQUN0RCxNQUFNLE1BQU0sR0FBRztZQUNiLEVBQUUsRUFBRSxZQUFZO1NBQ2pCLENBQUM7UUFDRixJQUFJLENBQUMsU0FBUyxDQUFDLE9BQU8sRUFBRSxNQUFNLENBQUMsQ0FBQztJQUNsQyxDQUFDO0lBRUQ7Ozs7T0FJRztJQUNILHlCQUF5QixDQUFDLE9BQWUsRUFBRSxnQkFBd0I7UUFDakUsTUFBTSxPQUFPLEdBQUcsY0FBYyxDQUFDLFdBQVcsQ0FBQztRQUMzQyxNQUFNLFlBQVksR0FBRyxHQUFHLE9BQU8sSUFBSSxnQkFBZ0IsRUFBRSxDQUFDO1FBQ3RELE1BQU0sTUFBTSxHQUFHO1lBQ2IsRUFBRSxFQUFFLFlBQVk7U0FDakIsQ0FBQztRQUNGLElBQUksQ0FBQyxXQUFXLENBQUMsT0FBTyxFQUFFLE1BQU0sQ0FBQyxDQUFDO0lBQ3BDLENBQUM7Q0FDSjtBQWxRRCxvQ0FrUUMifQ==