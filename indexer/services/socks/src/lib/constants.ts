import { Channel, WebsocketTopic } from '../types';

// https://developer.mozilla.org/en-US/docs/Web/API/CloseEvent/code
// 4000 - 4999 is reserved for private use
// 1006 can't be used as it's not meant to be used by endpoints
// Reference https://datatracker.ietf.org/doc/html/rfc6455#section-7.4.1
export const WS_CLOSE_GOING_AWAY: number = 1001;
export const WS_CLOSE_HEARTBEAT_TIMEOUT: number = 1011;
export const WS_CLOSE_CODE_ABNORMAL_CLOSURE: number = 4000;
export const WS_CLOSE_CODE_POLICY_VIOLATION: number = 1008;
export const WS_CLOSE_CODE_SERVICE_RESTART: number = 1012;

// https://github.com/nodejs/node/blob/master/lib/internal/errors.js#L1537
// Error code for writing to a destroyed stream
// This error is thrown when trying to send to a destroyed socket, which can happen when closing
// or writing to a websocket that was disconnected abruptly
export const ERR_WRITE_STREAM_DESTROYED: string = 'Cannot call write after a stream was destroyed';
// Error emitted by websocket connections when an invalid frame is received
export const ERR_INVALID_WEBSOCKET_FRAME: string = 'Invalid WebSocket frame';

export const WEBSOCKET_NOT_OPEN: string = 'ws not open';

export const V4_MARKETS_ID: string = 'v4_markets';
export const V4_BLOCK_HEIGHT_ID: string = 'v4_block_height';

export const TOPIC_TO_CHANNEL: Record<WebsocketTopic, Channel[]> = {
  [WebsocketTopic.TO_WEBSOCKETS_CANDLES]: [Channel.V4_CANDLES],
  [WebsocketTopic.TO_WEBSOCKETS_MARKETS]: [Channel.V4_MARKETS],
  [WebsocketTopic.TO_WEBSOCKETS_ORDERBOOKS]: [Channel.V4_ORDERBOOK],
  [WebsocketTopic.TO_WEBSOCKETS_SUBACCOUNTS]: [Channel.V4_ACCOUNTS, Channel.V4_PARENT_ACCOUNTS],
  [WebsocketTopic.TO_WEBSOCKETS_TRADES]: [Channel.V4_TRADES],
  [WebsocketTopic.TO_WEBSOCKETS_BLOCK_HEIGHT]: [Channel.V4_BLOCK_HEIGHT],
};

export const MAX_TIMEOUT_INTEGER: number = 2147483647;
