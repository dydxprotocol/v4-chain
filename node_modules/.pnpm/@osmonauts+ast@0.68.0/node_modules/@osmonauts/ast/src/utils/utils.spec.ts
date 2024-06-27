import { recursiveNamespace, renderNameSafely } from './utils';
import { createStargateClientOptions } from '../clients/stargate';
import { getGenericParseContextWithRef, expectCode } from '../../test-utils';
import { ProtoRef } from '@osmonauts/types';

it('recursiveNamespace', async () => {
    const ref: ProtoRef = {
        absolute: '/',
        filename: '/',
        proto: {
            imports: [],
            package: 'osmosis.gamm.yolo',
            root: {},
        }
    }
    expectCode(
        recursiveNamespace(['osmosis', 'gamm', 'v1beta', 'pools'].reverse(), [
            createStargateClientOptions({
                context: getGenericParseContextWithRef(ref),
                name: 'getSigningOsmosisClientOptions',
                aminoConverters: 'aminoConverters',
                protoTypeRegistry: 'protoTypeRegistry'
            })
        ])[0]
    );
});

describe('safe type names', () => {
    it('My_Name_asd.asdf.Type_rcc.dao.Yolo', () => {
        const name = 'My_Name_asd.asdf.Type_rcc.dao.Yolo';
        const filtered = renderNameSafely(name);
        expect(filtered).toEqual('My_Name_Type_Yolo');
    });
    it('dao.Yolo', () => {
        const name = 'dao.Yolo';
        const filtered = renderNameSafely(name);
        expect(filtered).toEqual('Yolo');
    });
    it('Yolo', () => {
        const name = 'Yolo';
        const filtered = renderNameSafely(name);
        expect(filtered).toEqual('Yolo');
    });
});
