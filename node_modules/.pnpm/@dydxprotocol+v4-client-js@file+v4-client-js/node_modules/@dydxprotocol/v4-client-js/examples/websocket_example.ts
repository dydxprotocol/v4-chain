import { Network } from '../src/clients/constants';
import { CandlesResolution, IncomingMessageTypes, SocketClient } from '../src/clients/socket-client';
import { DYDX_TEST_ADDRESS } from './constants';

function test(): void {
  const mySocket = new SocketClient(
    Network.testnet().indexerConfig,
    () => {
      console.log('socket opened');
    },
    () => {
      console.log('socket closed');
    },
    (message) => {
      console.log(message);
      if (typeof message.data === 'string') {
        const jsonString = message.data as string;
        try {
          const data = JSON.parse(jsonString);
          if (data.type === IncomingMessageTypes.CONNECTED) {
            mySocket.subscribeToMarkets();
            mySocket.subscribeToOrderbook('ETH-USD');
            mySocket.subscribeToTrades('ETH-USD');
            mySocket.subscribeToCandles('ETH-USD', CandlesResolution.FIFTEEN_MINUTES);
            mySocket.subscribeToSubaccount(DYDX_TEST_ADDRESS, 0);
          }
          console.log(data);
        } catch (e) {
          console.error('Error parsing JSON message:', e);
        }
      }
    },
  );
  mySocket.connect();
}

test();
