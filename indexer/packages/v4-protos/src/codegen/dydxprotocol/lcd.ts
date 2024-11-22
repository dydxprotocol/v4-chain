import { LCDClient } from "@osmonauts/lcd";
export const createLCDClient = async ({
  restEndpoint
}: {
  restEndpoint: string;
}) => {
  const requestClient = new LCDClient({
    restEndpoint
  });
  return {
    dydxprotocol: {
      accountplus: new (await import("./accountplus/query.lcd")).LCDQueryClient({
        requestClient
      }),
      affiliates: new (await import("./affiliates/query.lcd")).LCDQueryClient({
        requestClient
      }),
      assets: new (await import("./assets/query.lcd")).LCDQueryClient({
        requestClient
      }),
      blocktime: new (await import("./blocktime/query.lcd")).LCDQueryClient({
        requestClient
      }),
      bridge: new (await import("./bridge/query.lcd")).LCDQueryClient({
        requestClient
      }),
      clob: new (await import("./clob/query.lcd")).LCDQueryClient({
        requestClient
      }),
      delaymsg: new (await import("./delaymsg/query.lcd")).LCDQueryClient({
        requestClient
      }),
      epochs: new (await import("./epochs/query.lcd")).LCDQueryClient({
        requestClient
      }),
      feetiers: new (await import("./feetiers/query.lcd")).LCDQueryClient({
        requestClient
      }),
      listing: new (await import("./listing/query.lcd")).LCDQueryClient({
        requestClient
      }),
      perpetuals: new (await import("./perpetuals/query.lcd")).LCDQueryClient({
        requestClient
      }),
      prices: new (await import("./prices/query.lcd")).LCDQueryClient({
        requestClient
      }),
      ratelimit: new (await import("./ratelimit/query.lcd")).LCDQueryClient({
        requestClient
      }),
      revshare: new (await import("./revshare/query.lcd")).LCDQueryClient({
        requestClient
      }),
      rewards: new (await import("./rewards/query.lcd")).LCDQueryClient({
        requestClient
      }),
      stats: new (await import("./stats/query.lcd")).LCDQueryClient({
        requestClient
      }),
      subaccounts: new (await import("./subaccounts/query.lcd")).LCDQueryClient({
        requestClient
      }),
      vault: new (await import("./vault/query.lcd")).LCDQueryClient({
        requestClient
      }),
      vest: new (await import("./vest/query.lcd")).LCDQueryClient({
        requestClient
      })
    }
  };
};