/* eslint-disable @typescript-eslint/no-empty-function */
export class BaseTendermintClient {
  async block(): Promise<void> {}
  async broadcastTxSync(): Promise<void> {}
  async broadcastTxAsync(): Promise<void> {}
  async txSearchAll(): Promise<void> {}
}

export class BaseQueryClient {
  tx = {
    async simulate(): Promise<void> {},
  };

  async queryUnverified(): Promise<void> {}
}

export class BaseStargateSigningClient {
  async sign(): Promise<void> {}
}

export class BaseWallet {
  async getAccounts(): Promise<void> {}
}
