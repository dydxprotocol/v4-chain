import { createAminoConverter } from '../index';
import { snake } from 'case';
import { prepareContext, expectCode, getTestProtoStore } from '../../../../../test-utils/';

const store = getTestProtoStore();
store.traverseAll();

describe('cosmos/gov/v1beta1/tx', () => {

    const {
        context, root, protos
    } = prepareContext(store, 'cosmos/gov/v1beta1/tx.proto')

    it('AminoConverter', () => {
        context.options.aminoEncoding.casingFn = snake;
        expectCode(createAminoConverter({
            context,
            root,
            name: 'AminoConverter',
            protos
        }))
    })
});

