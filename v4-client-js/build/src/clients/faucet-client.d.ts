import { Response } from './lib/axios';
import RestClient from './modules/rest';
export declare class FaucetClient extends RestClient {
    /**
       * @description For testnet only, add USDC to an subaccount
       *
       * @returns The HTTP response.
       */
    fill(address: string, subaccountNumber: number, amount: number, headers?: {}): Promise<Response>;
}
