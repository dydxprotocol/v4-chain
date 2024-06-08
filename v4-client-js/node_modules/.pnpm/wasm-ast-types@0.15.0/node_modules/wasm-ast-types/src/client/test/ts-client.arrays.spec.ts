import message from '../../../../../__fixtures__/misc/schema/arrays.json';

import {
    createQueryClass,
    createExecuteClass,
    createExecuteInterface,
    createTypeInterface
} from '../client'
import { expectCode, printCode, makeContext } from '../../../test-utils';
import { getPropertyType } from '../../utils';

const ctx = makeContext(message);

it('getPropertyType', () => {
    const ast = getPropertyType(
        ctx,
        message.oneOf[0].properties.update_edges,
        'edges3'
    );
    expectCode(ast.type)
    // printCode(ast.type)
})

it('execute_msg_for__empty', () => {
    expectCode(createTypeInterface(
        ctx,
        message
    ))
})


it('query classes', () => {
    expectCode(createQueryClass(
        ctx,
        'SG721QueryClient',
        'SG721ReadOnlyInstance',
        message
    ))
});

it('execute classes array types', () => {
    expectCode(createExecuteClass(
        ctx,
        'SG721Client',
        'SG721Instance',
        null,
        message
    ))
});

it('execute interfaces no extends', () => {
    expectCode(createExecuteInterface(
        ctx,
        'SG721Instance',
        null,
        message
    ))
});
