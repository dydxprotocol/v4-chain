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
store.options.prototypes.parser.keepCase = true;
store.traverseAll();

it('cosmos/group/v1/query.proto', () => {
    const ref = store.findProto('cosmos/group/v1/query.proto');
    const res = traverse(store, ref);
    const service: ProtoService = getNestedProto(res).Query;
    const context = new GenericParseContext(ref, store, defaultTelescopeOptions);
    const ast = createLCDClient(context, service);
    expectCode(ast);
});
