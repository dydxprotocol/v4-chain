import Long from 'long';

import { generateRandomClientId } from '../../src/lib/utils';
import {
  ClobPairId,
  TransactionOptions,
  ICancelOrder,
  Transfer,
  OrderFlags,
  IPlaceOrder,
  Order_Side,
  Order_TimeInForce,
} from '../../src/types';

export const TEST_ADDRESS: string = 'dydx1vl9h9nkmau4e9v7tm30wekespu3d2qhd9404wa';
export const TEST_RECIPIENT_ADDRESS: string = 'dydx1slanxj8x9ntk9knwa6cvfv2tzlsq5gk3dshml0';
export const TEST_CHAIN_ID: string = 'dydxprotocol';
export const TEST_HOST: string = 'http://localhost:26657';

export const defaultTransactionOptions: TransactionOptions = {
  accountNumber: 1,
  sequence: 10,
  chainId: TEST_CHAIN_ID,
};

// PlaceOrder variables
export const defaultOrder: IPlaceOrder = {
  clientId: 0,
  orderFlags: OrderFlags.SHORT_TERM,
  clobPairId: ClobPairId.PERPETUAL_PAIR_BTC_USD,
  side: Order_Side.SIDE_BUY,
  quantums: Long.fromNumber(10),
  subticks: Long.fromNumber(100),
  goodTilBlock: 10000,
  timeInForce: Order_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
  reduceOnly: false,
  clientMetadata: 0,
};

// CancelOrder variables
export const defaultCancelOrder: ICancelOrder = {
  clientId: generateRandomClientId(),
  orderFlags: OrderFlags.SHORT_TERM,
  clobPairId: ClobPairId.PERPETUAL_PAIR_BTC_USD,
  goodTilBlock: 4250,
};

// Transfer variables
export const defaultTransfer: Transfer = {
  sender: {
    owner: TEST_ADDRESS,
    number: 0,
  },
  recipient: {
    owner: 'dydx14063jves4u9zhm7eja5ltf3t8zspxd92qnk23t',
    number: 0,
  },
  assetId: 0,
  amount: Long.fromNumber(1000),
};

// ------ Onboarding Constants ------ //
// Base Signature Result
export const SIGNATURE_RESULT = '0xf2006b4fc0afa08a6048b40d3d67e437ac2e20e6bacf3f947f9b33aa2756d204287f5fb2155bec067d879a25d7ddac791826cf28e2786919065e19848020f1531b';

// Derived privateKeyBytes from Base Signature Result
export const ENTROPY_FROM_SIGNATURE_RESULT: Uint8Array = new Uint8Array([
  247, 183, 226, 106, 76, 125, 241, 35, 149, 75, 103, 180, 165, 243, 80, 128,
  34, 20, 238, 201, 131, 180, 61, 76, 223, 179, 37, 211, 144, 197, 171, 251,
]);

// Derived HDKey privateKey from Base Signature Result
export const PRIVATE_KEY_FROM_SIGNATURE_RESULT: Uint8Array = new Uint8Array([
  14, 92, 178, 198, 64, 65, 27, 153, 11, 45, 118, 194, 71, 194, 10, 140, 145,
  40, 203, 107, 231, 191, 138, 220, 168, 104, 28, 69, 58, 193, 203, 213,
]);

// Derived HDKey publicKey from Base Signature Result
export const PUBLIC_KEY_FROM_SIGNATURE_RESULT: Uint8Array = new Uint8Array([
  2, 204, 94, 45, 7, 83, 196, 55, 195, 107, 10, 35, 169, 32, 24, 37, 193, 196,
  34, 210, 84, 26, 107, 89, 119, 240, 211, 63, 187, 34, 162, 197, 36,
]);

// Derived HDKey mnemonic from Base Signature Result
export const MNEMONIC_FROM_SIGNATURE_RESULT = 'waste sample once ocean tenant mushroom festival hollow regret convince stage able candy jazz champion isolate diary group under entry decorate glare quiz job';
