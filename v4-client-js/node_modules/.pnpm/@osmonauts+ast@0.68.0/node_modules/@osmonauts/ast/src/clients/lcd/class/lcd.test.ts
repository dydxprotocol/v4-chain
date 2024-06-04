import {
    getUrlTemplateString,
    createAggregatedLCDClient,
    createLCDClient,
    makeTemplateTagLegacy
} from './lcd';
import { ProtoStore, traverse, getNestedProto } from '@osmonauts/proto-parser'
import { defaultTelescopeOptions, ProtoService } from '@osmonauts/types';
import generate from '@babel/generator';
import { GenericParseContext } from '../../../encoding';
import { getTestProtoStore, expectCode, printCode } from '../../../../test-utils';
const store = getTestProtoStore();
store.traverseAll();

it('service info template', () => {
    expect(getUrlTemplateString('/osmosis/{gamm}/v1beta1/estimate/swap_exact_amount_in')).toMatchSnapshot();
    expect(getUrlTemplateString('/osmosis/{gamm}/v1beta1/{estimate}/swap_exact_amount_in')).toMatchSnapshot();
    expect(getUrlTemplateString('/osmosis/{gamm}/{v1beta1}/{estimate}/{swap_exact_amount_in}')).toMatchSnapshot();
    expect(getUrlTemplateString('/osmosis/gamm/v1beta1/estimate/{swap_exact_amount_in}')).toMatchSnapshot();
    expect(getUrlTemplateString('/cosmos/feegrant/v1beta1/allowance/{granter}/{grantee}')).toMatchSnapshot();
});

it('template tags', () => {
    const info = {
        url: '/{cosmos}/feegrant/v1beta1/{allowance}/{granter}/{grantee}',
        pathParams: [
            'cosmos',
            'allowance',
            'granter',
            'grantee'
        ]
    };
    expectCode(makeTemplateTagLegacy(info));
})

it('osmosis LCDClient', () => {
    const ref = store.findProto('osmosis/gamm/v1beta1/query.proto');
    const res = traverse(store, ref);
    const service: ProtoService = getNestedProto(res).Query;
    const context = new GenericParseContext(ref, store, defaultTelescopeOptions);
    const ast = createLCDClient(context, service);
    expectCode(ast);
});
it('cosmos LCDClient', () => {
    const ref = store.findProto('cosmos/bank/v1beta1/query.proto');
    const res = traverse(store, ref);
    const service: ProtoService = getNestedProto(res).Query;
    const context = new GenericParseContext(ref, store, defaultTelescopeOptions);
    const ast = createLCDClient(context, service);
    expectCode(ast);
});
it('cosmos fee LCDClient', () => {
    const ref = store.findProto('cosmos/feegrant/v1beta1/query.proto');
    const res = traverse(store, ref);
    const service: ProtoService = getNestedProto(res).Query;
    const context = new GenericParseContext(ref, store, defaultTelescopeOptions);
    const ast = createLCDClient(context, service);
    expectCode(ast);
});
it('cosmos/staking/v1beta1/query.proto', () => {
    const ref = store.findProto('cosmos/staking/v1beta1/query.proto');
    const res = traverse(store, ref);
    const service: ProtoService = getNestedProto(res).Query;
    const context = new GenericParseContext(ref, store, defaultTelescopeOptions);
    const ast = createLCDClient(context, service);
    expectCode(ast);
});
it('cosmos/app/v1alpha1/query.proto', () => {
    const ref = store.findProto('cosmos/app/v1alpha1/query.proto');
    const res = traverse(store, ref);
    const service: ProtoService = getNestedProto(res).Query;
    const context = new GenericParseContext(ref, store, defaultTelescopeOptions);
    const ast = createLCDClient(context, service);
    expectCode(ast);
});
it('cosmos/group/v1/query.proto', () => {
    const ref = store.findProto('cosmos/group/v1/query.proto');
    store.options.prototypes.parser.keepCase = true;
    const res = traverse(store, ref);
    const service: ProtoService = getNestedProto(res).Query;
    const context = new GenericParseContext(ref, store, defaultTelescopeOptions);
    const ast = createLCDClient(context, service);
    expectCode(ast);
});
it('cosmos/gov/v1beta1/query.proto', () => {
    const ref = store.findProto('cosmos/gov/v1beta1/query.proto');
    store.options.prototypes.parser.keepCase = true;
    const res = traverse(store, ref);
    const service: ProtoService = getNestedProto(res).Query;
    const context = new GenericParseContext(ref, store, defaultTelescopeOptions);
    const ast = createLCDClient(context, service);
    expectCode(ast);
});

