import { InvalidMessageHandler } from '../../src/lib/invalid-message';
import { RateLimiter } from '../../src/lib/rate-limit';
import { sendMessage } from '../../src/helpers/wss';
import WebSocket from 'ws';
import { WS_CLOSE_CODE_POLICY_VIOLATION } from '../../src/lib/constants';
import { Connection, WebsocketEvents } from '../../src/types';

jest.mock('../../src/lib/rate-limit');
jest.mock('../../src/helpers/wss', () => ({
  sendMessage: jest.fn(),
}));

describe('InvalidMessageHandler', () => {
  let invalidMessageHandler: InvalidMessageHandler;
  let mockConnection: Connection;
  const connectionId = 'testConnectionId';
  const responseMessage = 'Test response message';

  beforeEach(() => {
    (RateLimiter as jest.Mock).mockImplementation(() => ({
      rateLimit: jest.fn().mockReturnValue(0),
      removeConnection: jest.fn(),
    }));

    mockConnection = {
      ws: {
        close: jest.fn(),
        removeAllListeners: jest.fn(),
      } as unknown as WebSocket,
      messageId: 1,
    };
  });

  test('should send normal response message if not rate-limited', () => {
    invalidMessageHandler = new InvalidMessageHandler();
    invalidMessageHandler.handleInvalidMessage(responseMessage, mockConnection, connectionId);

    expect(sendMessage).toHaveBeenCalled();
    expect(mockConnection.ws.close).not.toHaveBeenCalled();
  });

  test('should rate limit, close connection, remove all event listeners for messages if over limit', () => {
    (RateLimiter as jest.Mock).mockImplementation(() => ({
      rateLimit: jest.fn().mockReturnValue(1000),
      removeConnection: jest.fn(),
    }));
    invalidMessageHandler = new InvalidMessageHandler();
    invalidMessageHandler.handleInvalidMessage(responseMessage, mockConnection, connectionId);

    expect(sendMessage).toHaveBeenCalled();
    expect(mockConnection.ws.close).toHaveBeenCalledWith(
      WS_CLOSE_CODE_POLICY_VIOLATION,
      JSON.stringify({ message: 'Rate limited' }),
    );
    expect(mockConnection.ws.removeAllListeners).toHaveBeenCalledWith(WebsocketEvents.MESSAGE);
  });
});
