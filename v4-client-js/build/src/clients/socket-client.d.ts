import { MessageEvent } from 'ws';
import { IndexerConfig } from './constants';
export declare enum IncomingMessageTypes {
    CONNECTED = "connected",
    SUBSCRIBED = "subscribed",
    ERROR = "error",
    CHANNEL_DATA = "channel_data",
    CHANNEL_BATCH_DATA = "channel_batch_data",
    PONG = "pong"
}
export declare enum CandlesResolution {
    ONE_MINUTE = "1MIN",
    FIVE_MINUTES = "5MINS",
    FIFTEEN_MINUTES = "15MINS",
    THIRTY_MINUTES = "30MINS",
    ONE_HOUR = "1HOUR",
    FOUR_HOURS = "4HOURS",
    ONE_DAY = "1DAY"
}
export declare class SocketClient {
    private url;
    private ws?;
    private onOpenCallback?;
    private onCloseCallback?;
    private onMessageCallback?;
    private pingInterval;
    private lastMessageTime;
    private pingIntervalId?;
    constructor(config: IndexerConfig, onOpenCallback: () => void, onCloseCallback: () => void, onMessageCallback: (event: MessageEvent) => void);
    connect(): void;
    /**
     * @description Close the websocket connection.
     *
     */
    close(): void;
    /**
     * @description Send data to the websocket connection.
     *
     */
    send(data: string): void;
    private handleOpen;
    private handleClose;
    private handleMessage;
    private restartPingInterval;
    /**
     * @description Set callback when the socket is opened.
     *
     */
    set onOpen(callback: () => void);
    /**
     * @description Set callback when the socket is closed.
     *
     */
    set onClose(callback: () => void);
    /**
     * @description Set callback when the socket receives a message.
     *
     */
    set onMessage(callback: (event: MessageEvent) => void);
    /**
     * @description Send a subscribe message to the websocket connection.
     *
     */
    subscribe(channel: string, params?: object): void;
    /**
     * @description Send an unsubscribe message to the websocket connection.
     *
     */
    unsubscribe(channel: string, params?: object): void;
    /**
     * @description Subscribe to markets channel.
     *
     */
    subscribeToMarkets(): void;
    /**
     * @description Unsubscribe from markets channel
     *
     */
    unsubscribeFromMarkets(): void;
    /**
     * @description Subscribe to trade channel
     * for a specific market.
     *
     */
    subscribeToTrades(market: string): void;
    /**
     * @description Unscribed from trade channel
     * for a specific market.
     *
     */
    unsubscribeFromTrades(market: string): void;
    /**
     * @description Subscribe to orderbook channel
     * for a specific market.
     *
     */
    subscribeToOrderbook(market: string): void;
    /**
     * @description Unsubscribe from orderbook channel
     * for a specific market.
     */
    unsubscribeFromOrderbook(market: string): void;
    /**
     * @description Subscribe to candles channel
     * for a specific market and resolution.
     *
     */
    subscribeToCandles(market: string, resolution: CandlesResolution): void;
    /**
     * @description Unsubscribe from candles channel
     * for a specific market and resolution.
     */
    unsubscribeFromCandles(market: string, resolution: CandlesResolution): void;
    /**
     * @description Subscribe to subaccount channel
     * for a specific address and subaccount number.
     */
    subscribeToSubaccount(address: string, subaccountNumber: number): void;
    /**
     * @description Unsubscribe from subaccount channel
     * for a specific address and subaccount number.
     *
     */
    unsubscribeFromSubaccount(address: string, subaccountNumber: number): void;
}
