"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const constants_1 = require("../src/clients/constants");
const socket_client_1 = require("../src/clients/socket-client");
const constants_2 = require("./constants");
function test() {
    const mySocket = new socket_client_1.SocketClient(constants_1.Network.testnet().indexerConfig, () => {
        console.log('socket opened');
    }, () => {
        console.log('socket closed');
    }, (message) => {
        console.log(message);
        if (typeof message.data === 'string') {
            const jsonString = message.data;
            try {
                const data = JSON.parse(jsonString);
                if (data.type === socket_client_1.IncomingMessageTypes.CONNECTED) {
                    mySocket.subscribeToMarkets();
                    mySocket.subscribeToOrderbook('ETH-USD');
                    mySocket.subscribeToTrades('ETH-USD');
                    mySocket.subscribeToCandles('ETH-USD', socket_client_1.CandlesResolution.FIFTEEN_MINUTES);
                    mySocket.subscribeToSubaccount(constants_2.DYDX_TEST_ADDRESS, 0);
                }
                console.log(data);
            }
            catch (e) {
                console.error('Error parsing JSON message:', e);
            }
        }
    });
    mySocket.connect();
}
test();
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoid2Vic29ja2V0X2V4YW1wbGUuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi9leGFtcGxlcy93ZWJzb2NrZXRfZXhhbXBsZS50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOztBQUFBLHdEQUFtRDtBQUNuRCxnRUFBcUc7QUFDckcsMkNBQWdEO0FBRWhELFNBQVMsSUFBSTtJQUNYLE1BQU0sUUFBUSxHQUFHLElBQUksNEJBQVksQ0FDL0IsbUJBQU8sQ0FBQyxPQUFPLEVBQUUsQ0FBQyxhQUFhLEVBQy9CLEdBQUcsRUFBRTtRQUNILE9BQU8sQ0FBQyxHQUFHLENBQUMsZUFBZSxDQUFDLENBQUM7SUFDL0IsQ0FBQyxFQUNELEdBQUcsRUFBRTtRQUNILE9BQU8sQ0FBQyxHQUFHLENBQUMsZUFBZSxDQUFDLENBQUM7SUFDL0IsQ0FBQyxFQUNELENBQUMsT0FBTyxFQUFFLEVBQUU7UUFDVixPQUFPLENBQUMsR0FBRyxDQUFDLE9BQU8sQ0FBQyxDQUFDO1FBQ3JCLElBQUksT0FBTyxPQUFPLENBQUMsSUFBSSxLQUFLLFFBQVEsRUFBRTtZQUNwQyxNQUFNLFVBQVUsR0FBRyxPQUFPLENBQUMsSUFBYyxDQUFDO1lBQzFDLElBQUk7Z0JBQ0YsTUFBTSxJQUFJLEdBQUcsSUFBSSxDQUFDLEtBQUssQ0FBQyxVQUFVLENBQUMsQ0FBQztnQkFDcEMsSUFBSSxJQUFJLENBQUMsSUFBSSxLQUFLLG9DQUFvQixDQUFDLFNBQVMsRUFBRTtvQkFDaEQsUUFBUSxDQUFDLGtCQUFrQixFQUFFLENBQUM7b0JBQzlCLFFBQVEsQ0FBQyxvQkFBb0IsQ0FBQyxTQUFTLENBQUMsQ0FBQztvQkFDekMsUUFBUSxDQUFDLGlCQUFpQixDQUFDLFNBQVMsQ0FBQyxDQUFDO29CQUN0QyxRQUFRLENBQUMsa0JBQWtCLENBQUMsU0FBUyxFQUFFLGlDQUFpQixDQUFDLGVBQWUsQ0FBQyxDQUFDO29CQUMxRSxRQUFRLENBQUMscUJBQXFCLENBQUMsNkJBQWlCLEVBQUUsQ0FBQyxDQUFDLENBQUM7aUJBQ3REO2dCQUNELE9BQU8sQ0FBQyxHQUFHLENBQUMsSUFBSSxDQUFDLENBQUM7YUFDbkI7WUFBQyxPQUFPLENBQUMsRUFBRTtnQkFDVixPQUFPLENBQUMsS0FBSyxDQUFDLDZCQUE2QixFQUFFLENBQUMsQ0FBQyxDQUFDO2FBQ2pEO1NBQ0Y7SUFDSCxDQUFDLENBQ0YsQ0FBQztJQUNGLFFBQVEsQ0FBQyxPQUFPLEVBQUUsQ0FBQztBQUNyQixDQUFDO0FBRUQsSUFBSSxFQUFFLENBQUMifQ==