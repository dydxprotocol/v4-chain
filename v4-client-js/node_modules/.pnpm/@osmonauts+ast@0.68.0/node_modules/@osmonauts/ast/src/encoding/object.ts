import * as t from '@babel/types';
import { ProtoType } from '@osmonauts/types';
import { fromPartialMethod } from './proto/from-partial';
import { decodeMethod } from './proto/decode';
import { encodeMethod } from './proto/encode';
import { fromJSONMethod } from './proto/from-json';
import { toJSONMethod } from './proto/to-json';
import { toSDKMethod } from './proto/to-sdk';
import { fromSDKMethod } from './proto/from-sdk';
import { ProtoParseContext } from './context';

export const createObjectWithMethods = (
    context: ProtoParseContext,
    name: string,
    proto: ProtoType
) => {

    const methods = [
        context.pluginValue('prototypes.methods.encode') && encodeMethod(context, name, proto),
        context.pluginValue('prototypes.methods.decode') && decodeMethod(context, name, proto),
        context.pluginValue('prototypes.methods.fromJSON') && fromJSONMethod(context, name, proto),
        context.pluginValue('prototypes.methods.toJSON') && toJSONMethod(context, name, proto),
        context.pluginValue('prototypes.methods.fromPartial') && fromPartialMethod(context, name, proto),
        context.pluginValue('prototypes.methods.fromSDK') && fromSDKMethod(context, name, proto),
        context.pluginValue('prototypes.methods.toSDK') && toSDKMethod(context, name, proto)
    ].filter(Boolean);

    return t.exportNamedDeclaration(
        t.variableDeclaration('const',
            [
                t.variableDeclarator(
                    t.identifier(name),
                    t.objectExpression(
                        methods
                    )
                )
            ]
        )
    )
};