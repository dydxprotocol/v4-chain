import { GeneratedType, Registry } from '@cosmjs/proto-signing';
import { defaultRegistryTypes } from '@cosmjs/stargate';
import {
  MsgPlaceOrder,
  MsgCancelOrder,
} from '@klyraprotocol/v4-proto/src/codegen/klyraprotocol/clob/tx';
import {
  MsgWithdrawFromSubaccount,
  MsgDepositToSubaccount,
} from '@klyraprotocol/v4-proto/src/codegen/klyraprotocol/sending/transfer';
import {
  MsgCreateTransfer,
} from '@klyraprotocol/v4-proto/src/codegen/klyraprotocol/sending/tx';

export const registry: ReadonlyArray<[string, GeneratedType]> = [];
export function generateRegistry(): Registry {
  return new Registry([
    // clob
    ['/klyraprotocol.clob.MsgPlaceOrder', MsgPlaceOrder as GeneratedType],
    ['/klyraprotocol.clob.MsgCancelOrder', MsgCancelOrder as GeneratedType],

    // sending
    ['/klyraprotocol.sending.MsgCreateTransfer', MsgCreateTransfer as GeneratedType],
    ['/klyraprotocol.sending.MsgWithdrawFromSubaccount', MsgWithdrawFromSubaccount as GeneratedType],
    ['/klyraprotocol.sending.MsgDepositToSubaccount', MsgDepositToSubaccount as GeneratedType],

    // default types
    ...defaultRegistryTypes,
  ]);
}
