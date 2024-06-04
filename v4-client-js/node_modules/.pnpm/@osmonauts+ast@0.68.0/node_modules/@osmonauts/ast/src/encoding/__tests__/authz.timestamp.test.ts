import { getNestedProto } from '@osmonauts/proto-parser';
import { defaultTelescopeOptions, expectCode, getTestProtoStore } from '../../../test-utils/'
import { AminoParseContext, ProtoParseContext } from '../context';
import { createProtoType } from '..';
import { createObjectWithMethods } from '../object';
import { createAminoConverter, makeAminoTypeInterface } from '../amino/index';
import cases from 'jest-in-case';
const store = getTestProtoStore();
store.traverseAll();

cases('cosmos/authz/v1beta1/authz', opts => {
    const ref = store.findProto('cosmos/authz/v1beta1/authz.proto');
    const context = new ProtoParseContext(ref, store, defaultTelescopeOptions);
    const aminoCtx = new AminoParseContext(ref, store, defaultTelescopeOptions);
    context.options.prototypes.typingsFormat.timestamp = opts.name;
    aminoCtx.options.prototypes.typingsFormat.timestamp = opts.name;
    expectCode(createProtoType(context, 'Grant',
        getNestedProto(ref.traversed).Grant
    ));
    expectCode(createObjectWithMethods(context, 'Grant',
        getNestedProto(ref.traversed).Grant
    ));
    expectCode(makeAminoTypeInterface({
        context: aminoCtx,
        proto: getNestedProto(ref.traversed).Grant
    }));
    expectCode(createAminoConverter({
        context: aminoCtx,
        root: ref.proto,
        name: 'AminoConverter',
        protos: [getNestedProto(ref.traversed).Grant]
    }));
}, [
    { name: 'date' },
    { name: 'timestamp' }
]);