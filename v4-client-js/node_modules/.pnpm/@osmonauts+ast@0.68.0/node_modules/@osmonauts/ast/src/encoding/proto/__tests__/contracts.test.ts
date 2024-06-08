import { getNestedProto } from '@osmonauts/proto-parser';
import { defaultTelescopeOptions, expectCode, getTestProtoStore } from '../../../../test-utils/'
import { ProtoParseContext } from '../../context';
import { createProtoType } from '..';
import { createObjectWithMethods } from '../../object';

const store = getTestProtoStore();
store.traverseAll();

describe('MsgExecuteContract', () => {
    const ref = store.findProto('cosmwasm/wasm/v1/tx.proto');
    const context = new ProtoParseContext(ref, store, defaultTelescopeOptions);
    it('interface', () => {
        expectCode(createProtoType(context, 'MsgExecuteContract',
            getNestedProto(ref.traversed).MsgExecuteContract
        ));
    });
    it('interface', () => {
        expectCode(createObjectWithMethods(context, 'MsgExecuteContract',
            getNestedProto(ref.traversed).MsgExecuteContract
        ));
    });
});