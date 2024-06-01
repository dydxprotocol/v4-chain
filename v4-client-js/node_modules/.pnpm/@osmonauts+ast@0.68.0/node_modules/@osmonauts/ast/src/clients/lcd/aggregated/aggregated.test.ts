import {
    createAggregatedLCDClient,
} from '../class';
import { ProtoStore, traverse, getNestedProto } from '@osmonauts/proto-parser'
import { defaultTelescopeOptions, ProtoRef, ProtoService } from '@osmonauts/types';
import generate from '@babel/generator';
import { GenericParseContext } from '../../../encoding';
import { getTestProtoStore } from '../../../../test-utils';

const store = getTestProtoStore();
store.traverseAll();

const expectCode = (ast) => {
    expect(
        generate(ast).code
    ).toMatchSnapshot();
}
const printCode = (ast) => {
    console.log(
        generate(ast).code
    );
}

it('AggregatedLCDClient', () => {
    const ref1 = store.findProto('cosmos/bank/v1beta1/query.proto');
    const ref2 = store.findProto('osmosis/gamm/v1beta1/query.proto');
    const res1 = traverse(store, ref1);
    const res2 = traverse(store, ref2);
    const service1: ProtoService = getNestedProto(res1).Query;
    const service2: ProtoService = getNestedProto(res2).Query;
    const context = new GenericParseContext(ref1, store, defaultTelescopeOptions);
    const ast = createAggregatedLCDClient(context, [service1, service2], 'QueryClient');
    expectCode(ast);
});

// TODO - use package names to shape the class
// e.g. osmosis.gamm.v1beta1.pools()

it('options', () => {

    const packages = [
        'cosmos.bank.v1beta1',
        'osmosis.gamm.v1beta1',
    ];

    const refs = store.filterProtoWhere((ref: ProtoRef) => {
        return packages.includes(ref.proto.package)
    });

    const services: ProtoService[] = refs.map(ref => {
        const proto = getNestedProto(ref.traversed);
        if (!proto?.Query || proto.Query?.type !== 'Service') {
            return;
        }
        return proto.Query;
    }).filter(Boolean);

    const context = new GenericParseContext(refs[0], store, defaultTelescopeOptions);
    const ast = createAggregatedLCDClient(context, services, 'QueryClient');
    expectCode(ast);
});