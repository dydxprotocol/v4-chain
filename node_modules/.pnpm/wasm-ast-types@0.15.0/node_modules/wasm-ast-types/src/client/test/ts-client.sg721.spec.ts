import execute_msg_for__empty from '../../../../../__fixtures__/sg721/execute_msg_for__empty.json';
import {
    createQueryClass,
    createExecuteClass,
    createExecuteInterface,
    createTypeInterface
} from '../client'
import { expectCode, makeContext } from '../../../test-utils';

const ctx = makeContext(execute_msg_for__empty);

it('execute_msg_for__empty', () => {
    expectCode(createTypeInterface(
        ctx,
        execute_msg_for__empty
    ))
})


it('query classes', () => {
    expectCode(createQueryClass(
        ctx,
        'SG721QueryClient',
        'SG721ReadOnlyInstance',
        execute_msg_for__empty
    ))
});

it('execute classes array types', () => {
    expectCode(createExecuteClass(
        ctx,
        'SG721Client',
        'SG721Instance',
        null,
        execute_msg_for__empty
    ))
});

it('execute interfaces no extends', () => {
    expectCode(createExecuteInterface(
        ctx,
        'SG721Instance',
        null,
        execute_msg_for__empty
    ))
});
