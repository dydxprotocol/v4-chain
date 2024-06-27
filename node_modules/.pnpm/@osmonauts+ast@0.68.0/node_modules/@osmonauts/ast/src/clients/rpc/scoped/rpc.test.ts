import { expectCode, getGenericParseContext, printCode } from '../../../../test-utils'

import {
    rpcRecursiveObjectProps,
    createScopedRpcFactory,
    createScopedRpcTmFactory
} from './rpc';

it('rpcRecursiveObjectProps', async () => {
    expectCode(rpcRecursiveObjectProps(['cosmos', 'bank', 'v1beta1']))
});

it('createScopedRpcFactory', async () => {
    expectCode(createScopedRpcFactory({
        cosmos: {
            bank: {
                v1beta1: "./proto/cosmos/bank/v1beta1/query.lcd"
            },
            gov: {
                v1beta1: "./proto/cosmos/bank/v1beta1/query.lcd"
            },
        },
        osmosis: {
            gamm: {
                v1beta1: "./proto/cosmos/bank/v1beta1/query.lcd"
            }
        }
    },
        'createRpc',
        'RPCQueryClient'
    ))
});

it('createScopedRpcTmFactory', async () => {
    const context = getGenericParseContext();
    expectCode(createScopedRpcTmFactory(context, {
        cosmos: {
            bank: {
                v1beta1: "./proto/cosmos/bank/v1beta1/query.lcd"
            },
            gov: {
                v1beta1: "./proto/cosmos/bank/v1beta1/query.lcd"
            },
        },
        osmosis: {
            gamm: {
                v1beta1: "./proto/cosmos/bank/v1beta1/query.lcd"
            }
        }
    },
        'createRpcFactorizzle',
    ))
});
