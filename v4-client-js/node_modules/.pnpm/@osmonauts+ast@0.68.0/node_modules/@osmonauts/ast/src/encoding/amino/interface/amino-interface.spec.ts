import { makeAminoTypeInterface } from './index';
import { ProtoStore } from '@osmonauts/proto-parser'
import { snake } from 'case';
import { camel } from '@osmonauts/utils';
import { prepareContext, expectCode, printCode, getTestProtoStore } from '../../../../test-utils';
const store = getTestProtoStore();

store.traverseAll();

describe('osmosis/gamm/v1beta1/tx', () => {

    const {
        context, protos
    } = prepareContext(store, 'osmosis/gamm/v1beta1/tx.proto')

    it('Interfaces', () => {
        context.options.aminoEncoding.casingFn = camel;
        expectCode(makeAminoTypeInterface(
            {
                context,
                proto: protos.find(p => p.name === 'MsgJoinPool'),
            }
        ))
    })
});


describe('cosmos/staking/v1beta1/tx', () => {
    const {
        context, protos
    } = prepareContext(store, 'cosmos/staking/v1beta1/tx.proto')

    it('MsgCreateValidator', () => {
        context.options.aminoEncoding.casingFn = snake;

        expectCode(makeAminoTypeInterface(
            {
                context,
                proto: protos.find(p => p.name === 'MsgCreateValidator'),
            }
        ))
    })
    it('MsgEditValidator', () => {
        context.options.aminoEncoding.casingFn = snake;
        expectCode(makeAminoTypeInterface(
            {
                context,
                proto: protos.find(p => p.name === 'MsgEditValidator'),
            }
        ))
    })
    it('MsgUndelegate', () => {
        context.options.aminoEncoding.casingFn = snake;
        expectCode(makeAminoTypeInterface(
            {
                context,
                proto: protos.find(p => p.name === 'MsgUndelegate'),
            }
        ))
    })
});

