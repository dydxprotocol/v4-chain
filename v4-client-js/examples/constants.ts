import { Order_Side, Order_TimeInForce } from '@klyraprotocol/v4-proto/src/codegen/klyraprotocol/clob/order';
import Long from 'long';

import { IPlaceOrder, OrderFlags } from '../src/clients/types';

export const KLYRA_TEST_ADDRESS = 'klyra14zzueazeh0hj67cghhf9jypslcf9sh2nt80jtq';
export const KLYRA_TEST_PRIVATE_KEY = 'e92a6595c934c991d3b3e987ea9b3125bf61a076deab3a9cb519787b7b3e8d77';
export const KLYRA_TEST_MNEMONIC = 'mirror actor skill push coach wait confirm orchard lunch mobile athlete gossip awake miracle matter bus reopen team ladder lazy list timber render wait';
export const KLYRA_LOCAL_ADDRESS = 'klyra199tqg4wdlnu4qjlxchpd7seg454937hju8xa57';
export const KLYRA_LOCAL_MNEMONIC = 'merge panther lobster crazy road hollow amused security before critic about cliff exhibit cause coyote talent happy where lion river tobacco option coconut small';

export const MARKET_BTC_USD: string = 'BTC-USD';
export const PERPETUAL_PAIR_BTC_USD: number = 0;

const quantums: Long = new Long(1_000_000_000);
const subticks: Long = new Long(1_000_000_000);

export const MAX_CLIENT_ID = 2 ** 32 - 1;

// PlaceOrder variables
export const defaultOrder: IPlaceOrder = {
  clientId: 0,
  orderFlags: OrderFlags.SHORT_TERM,
  clobPairId: PERPETUAL_PAIR_BTC_USD,
  side: Order_Side.SIDE_BUY,
  quantums,
  subticks,
  timeInForce: Order_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
  reduceOnly: false,
  clientMetadata: 0,
  routerFeePpm: 0,
};
