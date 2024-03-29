import { createPongMessage } from '../helpers/message';
import { sendMessage } from '../helpers/wss';
import { Connection, PingMessage } from '../types';

export class PingHandler {

  public handlePing(
    pingMessage: PingMessage,
    connection: Connection,
    connectionId: string,
  ): void {
    sendMessage(
      connection.ws,
      connectionId,
      createPongMessage(
        connectionId,
        connection.messageId,
        pingMessage.id,
      ),
    );
  }
}
