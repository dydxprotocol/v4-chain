import { getGenericParseContextWithRef, expectCode, printCode } from '../../../test-utils'
import {
    createStargateClient,
    createStargateClientOptions,
    createStargateClientAminoRegistry,
    createStargateClientProtoRegistry
} from './stargate';
import { ProtoRef } from '@osmonauts/types';

it('createStargateClient', async () => {
    const ref: ProtoRef = {
        absolute: '/',
        filename: '/',
        proto: {
            imports: [],
            package: 'osmosis.gamm.yolo',
            root: {},
        }
    }
    const context = getGenericParseContextWithRef(ref);
    expectCode(createStargateClient({
        context,
        name: 'getSigningOsmosisClient',
        options: 'getSigningOsmosisClientOptions'
    }));
    expect(context.utils).toMatchSnapshot();
});

it('createStargateClientOptions', async () => {
    const ref: ProtoRef = {
        absolute: '/',
        filename: '/',
        proto: {
            imports: [],
            package: 'somepackage1.gamm.yolo',
            root: {},
        }
    }
    const context = getGenericParseContextWithRef(ref);
    context.options.stargateClients.includeCosmosDefaultTypes = true;
    expectCode(createStargateClientOptions({
        context,
        aminoConverters: 'aminoConverters',
        protoTypeRegistry: 'protoTypeRegistry',
        name: 'getSigningOsmosisClientOptions'
    }));
    expect(context.utils).toMatchSnapshot();
});

it('createStargateClientAminoRegistry', async () => {
    const ref: ProtoRef = {
        absolute: '/',
        filename: '/',
        proto: {
            imports: [],
            package: 'somepackage1.gamm.yolo',
            root: {},
        }
    }
    const context = getGenericParseContextWithRef(ref);
    context.options.stargateClients.includeCosmosDefaultTypes = true;
    expectCode(createStargateClientAminoRegistry({
        context,
        aminoConverters: 'aminoConverters',
        aminos: [
            'somepackage1.gamm.v1beta1',
            'somepackage1.superfluid.v1beta1',
            'somepackage1.lockup'
        ]
    }));
    expect(context.utils).toMatchSnapshot();
});

it('createStargateClientProtoRegistry', async () => {
    const ref: ProtoRef = {
        absolute: '/',
        filename: '/',
        proto: {
            imports: [],
            package: 'somepackage1.gamm.yolo',
            root: {},
        }
    }
    const context = getGenericParseContextWithRef(ref);
    context.options.stargateClients.includeCosmosDefaultTypes = true;
    expectCode(createStargateClientProtoRegistry({
        context,
        protoTypeRegistry: 'protoTypeRegistry',
        registries: [
            'somepackage1.gamm.v1beta1',
            'somepackage1.superfluid.v1beta1',
            'somepackage1.lockup'
        ]
    }));
    expect(context.utils).toMatchSnapshot();
});

it('createStargateClient w/o defaults', async () => {
    const ref: ProtoRef = {
        absolute: '/',
        filename: '/',
        proto: {
            imports: [],
            package: 'otherpackage1.gamm.yolo',
            root: {},
        }
    }
    const context = getGenericParseContextWithRef(ref);
    context.options.stargateClients.includeCosmosDefaultTypes = false;
    expectCode(createStargateClient({
        context,
        name: 'getSigningOsmosisClient',
        options: 'getSigningOsmosisClientOptions',
    }));
    expectCode(createStargateClientOptions({
        context,
        name: 'getSigningOsmosisClientOptions',
        aminoConverters: 'aminoConverters',
        protoTypeRegistry: 'protoTypeRegistry'
    }));
    expect(context.utils).toMatchSnapshot();
});