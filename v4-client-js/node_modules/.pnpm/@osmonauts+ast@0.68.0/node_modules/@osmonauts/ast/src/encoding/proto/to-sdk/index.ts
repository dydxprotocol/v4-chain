import * as t from '@babel/types';
import { getFieldOptionality, getOneOfs } from '..';
import { identifier, objectMethod } from '../../../utils';
import { ProtoParseContext } from '../../context';
import { ProtoField, ProtoType } from '@osmonauts/types';
import { arrayTypes, toSDK } from './utils';
import { pascal } from 'case';

const needsImplementation = (name: string, field: ProtoField) => {
    throw new Error(`need to implement toSDK (${field.type} rules[${field.rule}] name[${name}])`);
}

export interface ToSDKMethod {
    context: ProtoParseContext;
    field: ProtoField;
    isOptional: boolean;
}

export const toSDKMethodFields = (context: ProtoParseContext, name: string, proto: ProtoType) => {
    const oneOfs = getOneOfs(proto);
    const fields = Object.keys(proto.fields ?? {}).reduce((m, fieldName) => {
        const field = {
            name: fieldName,
            ...proto.fields[fieldName]
        };

        const isOneOf = oneOfs.includes(fieldName);
        const isOptional = getFieldOptionality(context, field, isOneOf);

        const args: ToSDKMethod = {
            context,
            field,
            isOptional
        };

        if (field.rule === 'repeated') {
            switch (field.type) {
                case 'string':
                    return [...m, toSDK.array(args, arrayTypes.string())];
                case 'bytes':
                    return [...m, toSDK.array(args, arrayTypes.bytes(args))];
                case 'bool':
                    return [...m, toSDK.array(args, arrayTypes.bool())];
                case 'double':
                    return [...m, toSDK.array(args, arrayTypes.double())];
                case 'float':
                    return [...m, toSDK.array(args, arrayTypes.float())];
                case 'int32':
                    return [...m, toSDK.array(args, arrayTypes.int32())];
                case 'sint32':
                    return [...m, toSDK.array(args, arrayTypes.sint32())];
                case 'uint32':
                    return [...m, toSDK.array(args, arrayTypes.uint32())];
                case 'fixed32':
                    return [...m, toSDK.array(args, arrayTypes.fixed32())];
                case 'sfixed32':
                    return [...m, toSDK.array(args, arrayTypes.sfixed32())];
                case 'int64':
                    return [...m, toSDK.array(args, arrayTypes.int64(args))];
                case 'sint64':
                    return [...m, toSDK.array(args, arrayTypes.sint64(args))];
                case 'uint64':
                    return [...m, toSDK.array(args, arrayTypes.uint64(args))];
                case 'fixed64':
                    return [...m, toSDK.array(args, arrayTypes.fixed64(args))];
                case 'sfixed64':
                    return [...m, toSDK.array(args, arrayTypes.sfixed64(args))];
                default:
                    switch (field.parsedType.type) {
                        case 'Enum':
                            return [...m, toSDK.array(args, arrayTypes.enum(args))];
                        case 'Type':
                            return [...m, toSDK.array(args, arrayTypes.type(args))];
                    }
                    return needsImplementation(fieldName, field);
            }

        }

        if (field.keyType) {
            switch (field.keyType) {
                case 'string':
                case 'int32':
                case 'sint32':
                case 'uint32':
                case 'fixed32':
                case 'sfixed32':
                case 'int64':
                case 'sint64':
                case 'uint64':
                case 'fixed64':
                case 'sfixed64':
                    return [...m, ...toSDK.keyHash(args)];
                default:
                    return needsImplementation(fieldName, field);
            }
        }

        switch (field.type) {
            case 'string':
                return [...m, toSDK.string(args)];
            case 'double':
                return [...m, toSDK.double(args)];
            case 'float':
                return [...m, toSDK.float(args)];
            case 'bytes':
                return [...m, toSDK.bytes(args)];
            case 'bool':
                return [...m, toSDK.bool(args)];
            case 'int32':
                return [...m, toSDK.int32(args)];
            case 'sint32':
                return [...m, toSDK.sint32(args)];
            case 'uint32':
                return [...m, toSDK.uint32(args)];
            case 'fixed32':
                return [...m, toSDK.fixed32(args)];
            case 'sfixed32':
                return [...m, toSDK.sfixed32(args)];
            case 'int64':
                return [...m, toSDK.int64(args)];
            case 'sint64':
                return [...m, toSDK.sint64(args)];
            case 'uint64':
                return [...m, toSDK.uint64(args)];
            case 'fixed64':
                return [...m, toSDK.fixed64(args)];
            case 'sfixed64':
                return [...m, toSDK.sfixed64(args)];
            case 'google.protobuf.Duration':
            case 'Duration':
                return [...m, toSDK.duration(args)];
            case 'google.protobuf.Timestamp':
            case 'Timestamp':
                return [...m, toSDK.timestamp(args)];
            default:
                switch (field.parsedType.type) {
                    case 'Enum':
                        return [...m, toSDK.enum(args)];
                    case 'Type':
                        return [...m, toSDK.type(args)];
                }
                return needsImplementation(fieldName, field);
        }
    }, []);
    return fields;
};

export const toSDKMethod = (context: ProtoParseContext, name: string, proto: ProtoType) => {
    const fields = toSDKMethodFields(context, name, proto);
    let varName = 'message';
    if (!fields.length) {
        varName = '_';
    }

    const SDKTypeName =
        [name, 'SDKType']
            .filter(Boolean).join('')


    return objectMethod('method',
        t.identifier('toSDK'),
        [
            identifier(
                varName,
                t.tsTypeAnnotation(
                    t.tsTypeReference(
                        t.identifier(name)
                    )
                )
            )
        ],
        t.blockStatement([
            t.variableDeclaration(
                'const',
                [
                    t.variableDeclarator(
                        identifier('obj', t.tsTypeAnnotation(t.tsAnyKeyword())),
                        t.objectExpression([])
                    )
                ]
            ),

            ...fields,

            // RETURN 
            t.returnStatement(t.identifier('obj'))

        ]),
        false,
        false,
        false,
        t.tsTypeAnnotation(
            t.tsTypeReference(
                t.identifier(SDKTypeName)
            )
        )
    );
};
