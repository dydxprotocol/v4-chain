import { Response } from './lib/axios';
import RestClient from './modules/rest';

export class FaucetClient extends RestClient {
  /**
     * @description For testnet only, add USDC to an subaccount
     *
     * @returns The HTTP response.
     */
  public async fill(
    address: string,
    subaccountNumber: number,
    amount: number,
    headers?: {},
  ): Promise<Response> {
    const uri = '/faucet/tokens';

    return this.post(
      uri,
      {},
      {
        address,
        subaccountNumber,
        amount,
      },
      headers,
    );
  }
}
