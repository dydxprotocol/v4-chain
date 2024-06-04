import execute_msg from '../../../../__fixtures__/basic/execute_msg_for__empty.json';
import {
    createMessageComposerClass,
    createMessageComposerInterface
} from './message-composer'
import { expectCode, makeContext } from '../../test-utils';

it('execute classes', () => {
    const ctx = makeContext(execute_msg);
    expectCode(createMessageComposerClass(
        ctx,
        'SG721MessageComposer',
        'SG721Message',
        execute_msg
    ))
});

it('createMessageComposerInterface', () => {
    const ctx = makeContext(execute_msg);
    expectCode(createMessageComposerInterface(
        ctx,
        'SG721Message',
        execute_msg
    ))
});
