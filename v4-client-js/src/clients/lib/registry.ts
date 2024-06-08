import { GeneratedType, Registry } from '@cosmjs/proto-signing';
import { defaultRegistryTypes } from '@cosmjs/stargate';
import {
  MsgPlaceOrder,
  MsgCancelOrder,
} from '@dydxprotocol/v4-proto/src/codegen/dydxprotocol/clob/tx';
import {
  MsgWithdrawFromSubaccount,
  MsgDepositToSubaccount,
} from '@dydxprotocol/v4-proto/src/codegen/dydxprotocol/sending/transfer';
import {
  MsgCreateTransfer,
} from '@dydxprotocol/v4-proto/src/codegen/dydxprotocol/sending/tx';

export const registry: ReadonlyArray<[string, GeneratedType]> = [];
export function generateRegistry(): Registry {
  return new Registry([
    // clob
    ['/dydxprotocol.clob.MsgPlaceOrder', MsgPlaceOrder as GeneratedType],
    ['/dydxprotocol.clob.MsgCancelOrder', MsgCancelOrder as GeneratedType],

    // sending
    ['/dydxprotocol.sending.MsgCreateTransfer', MsgCreateTransfer as GeneratedType],
    ['/dydxprotocol.sending.MsgWithdrawFromSubaccount', MsgWithdrawFromSubaccount as GeneratedType],
    ['/dydxprotocol.sending.MsgDepositToSubaccount', MsgDepositToSubaccount as GeneratedType],

    // default types
    ...defaultRegistryTypes,
  ]);
}
