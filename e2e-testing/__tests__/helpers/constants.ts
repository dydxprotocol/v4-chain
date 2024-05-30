import {
  IPlaceOrder, Order_Side, Order_TimeInForce, OrderFlags,
} from '@dydxprotocol/v4-client-js';
import Long from 'long';

import { OrderDetails } from './types';

export const DYDX_LOCAL_ADDRESS = 'dydx1q90l6j6lzzgt460ehjj56azknlt5yrd4egfh9f';
export const DYDX_LOCAL_MNEMONIC = 'exile install vapor thing little toss immune notable lounge december final easy strike title end program interest quote cloth forget forward job october twenty';
export const DYDX_LOCAL_ADDRESS_2 = 'dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs';
export const DYDX_LOCAL_MNEMONIC_2 = 'grunt list hour endless observe better spoil penalty lab duck only layer vague fantasy satoshi record demise topple space shaft solar practice donor sphere';

export const MNEMONIC_TO_ADDRESS: Record<string, string> = {
  [DYDX_LOCAL_MNEMONIC]: DYDX_LOCAL_ADDRESS,
  [DYDX_LOCAL_MNEMONIC_2]: DYDX_LOCAL_ADDRESS_2,
};

export const ADDRESS_TO_MNEMONIC: Record<string, string> = {
  [DYDX_LOCAL_ADDRESS]: DYDX_LOCAL_MNEMONIC,
  [DYDX_LOCAL_ADDRESS_2]: DYDX_LOCAL_MNEMONIC_2,
};

export const PERPETUAL_PAIR_BTC_USD: number = 0;
export const quantums: Long = new Long(1_000_000_000);
export const subticks: Long = new Long(1_000_000_000);

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
};

export const orderDetails: OrderDetails[] = [
  {
    mnemonic: DYDX_LOCAL_MNEMONIC,
    timeInForce: 0,
    orderFlags: 64,
    side: 1,
    clobPairId: PERPETUAL_PAIR_BTC_USD,
    quantums: 10000000,
    subticks: 5000000000,
  },
  {
    mnemonic: DYDX_LOCAL_MNEMONIC_2,
    timeInForce: 0,
    orderFlags: 64,
    side: 2,
    clobPairId: PERPETUAL_PAIR_BTC_USD,
    quantums: 5000000,
    subticks: 5000000000,
  },
];
