import { globContracts, makeContext } from '../../../test-utils'
import {
    createQueryClass,
    createExecuteClass,
    createExecuteInterface,
    createTypeInterface
} from '../client'
import { expectCode } from '../../../test-utils';
import cases from 'jest-in-case';

const contracts = globContracts('issues/71');

cases('execute_msg_for__empty', async opts => {
    const ctx = makeContext(opts.content);
    expectCode(createTypeInterface(
        ctx,
        opts.content
    ));
}, contracts);

cases('query classes', async opts => {
    const ctx = makeContext(opts.content);
    expectCode(createQueryClass(
        ctx,
        'SG721QueryClient',
        'SG721ReadOnlyInstance',
        opts.content
    ))
}, contracts);

cases('execute class', async opts => {
    const ctx = makeContext(opts.content);
    expectCode(createExecuteClass(
        ctx,
        'SG721Client',
        'SG721Instance',
        null,
        opts.content
    ))
}, contracts);

cases('execute interface', async opts => {
    const ctx = makeContext(opts.content);
    expectCode(createExecuteInterface(
        ctx,
        'SG721Instance',
        null,
        opts.content
    ))
}, contracts);

