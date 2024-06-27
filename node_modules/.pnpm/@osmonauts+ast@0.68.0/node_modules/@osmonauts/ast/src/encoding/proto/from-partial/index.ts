import * as t from '@babel/types';
import { getFieldOptionality, getOneOfs } from '..';
import { identifier, objectMethod } from '../../../utils';
import { ProtoParseContext } from '../../context';
import { ProtoType, ProtoField } from '@osmonauts/types';
import { getBaseCreateTypeFuncName } from '../types';
import { arrayTypes, fromPartial } from './utils';

const needsImplementation = (name: string, field: ProtoField) => {
    throw new Error(`need to implement fromPartial (${field.type} rules[${field.rule}] name[${name}])`);
}
export interface FromPartialMethod {
    context: ProtoParseContext;
    field: ProtoField;
    isOneOf: boolean;
    isOptional: boolean;
}

export const fromPartialMethodFields = (context: ProtoParseContext, name: string, proto: ProtoType) => {
    const oneOfs = getOneOfs(proto);
    const fields = Object.keys(proto.fields ?? {}).map(fieldName => {
        const field = {
            name: fieldName,
            ...proto.fields[fieldName]
        };

        const isOneOf = oneOfs.includes(fieldName);
        const isOptional = getFieldOptionality(context, field, isOneOf);

        const args: FromPartialMethod = {
            context,
            field,
            isOneOf,
            isOptional
        };

        if (field.rule === 'repeated') {
            switch (field.type) {
                case 'string':
                    return fromPartial.array(args, arrayTypes.string());
                case 'bool':
                    return fromPartial.array(args, arrayTypes.bool());
                case 'bytes':
                    return fromPartial.array(args, arrayTypes.bytes());
                case 'float':
                    return fromPartial.array(args, arrayTypes.float());
                case 'double':
                    return fromPartial.array(args, arrayTypes.double());
                case 'int32':
                    return fromPartial.array(args, arrayTypes.int32());
                case 'sint32':
                    return fromPartial.array(args, arrayTypes.sint32());
                case 'uint32':
                    return fromPartial.array(args, arrayTypes.uint32());
                case 'fixed32':
                    return fromPartial.array(args, arrayTypes.fixed32());
                case 'sfixed32':
                    return fromPartial.array(args, arrayTypes.sfixed32());
                case 'int64':
                    return fromPartial.array(args, arrayTypes.int64());
                case 'sint64':
                    return fromPartial.array(args, arrayTypes.sint64());
                case 'uint64':
                    return fromPartial.array(args, arrayTypes.uint64());
                case 'fixed64':
                    return fromPartial.array(args, arrayTypes.fixed64());
                case 'sfixed64':
                    return fromPartial.array(args, arrayTypes.sfixed64());
                default:
                    switch (field.parsedType.type) {
                        case 'Enum':
                            return fromPartial.array(args, arrayTypes.enum());
                        case 'Type':
                            return fromPartial.array(args, arrayTypes.type(args));
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
                    return fromPartial.keyHash(args);
                default:
                    return needsImplementation(fieldName, field);
            }
        }

        switch (field.type) {
            case 'string':
                return fromPartial.string(args);
            case 'bytes':
                return fromPartial.bytes(args);
            case 'bool':
                return fromPartial.bool(args);
            case 'double':
                return fromPartial.double(args);
            case 'float':
                return fromPartial.float(args);
            case 'int32':
                return fromPartial.int32(args);
            case 'sint32':
                return fromPartial.sint32(args);
            case 'uint32':
                return fromPartial.uint32(args);
            case 'fixed32':
                return fromPartial.fixed32(args);
            case 'sfixed32':
                return fromPartial.sfixed32(args);
            case 'int64':
                return fromPartial.int64(args);
            case 'sint64':
                return fromPartial.sint64(args);
            case 'uint64':
                return fromPartial.uint64(args);
            case 'fixed64':
                return fromPartial.fixed64(args);
            case 'sfixed64':
                return fromPartial.sfixed64(args);
            case 'google.protobuf.Duration':
            case 'Duration':
                return fromPartial.duration(args);
            case 'google.protobuf.Timestamp':
            case 'Timestamp':
                return fromPartial.timestamp(args);
            default:
                switch (field.parsedType.type) {
                    case 'Enum':
                        return fromPartial.enum(args);
                    case 'Type':
                        return fromPartial.type(args);
                }
                return needsImplementation(fieldName, field);
        }
    });
    return fields;
};

export const fromPartialMethod = (context: ProtoParseContext, name: string, proto: ProtoType) => {
    const useDeepPartial = context.pluginValue('prototypes.typingsFormat.useDeepPartial');

    let partialName: 'Partial' | 'DeepPartial' = 'Partial';
    if (useDeepPartial) {
        context.addUtil('DeepPartial');
        partialName = 'DeepPartial';
    }


    const fields = fromPartialMethodFields(context, name, proto);
    let varName = 'object';
    if (!fields.length) {
        varName = '_';
    }

    let typeParameters = undefined;
    let param = null;

    const useExact = context.pluginValue('prototypes.typingsFormat.useExact');
    if (useExact === true) {
        context.addUtil('Exact');

        // type params
        typeParameters = t.tsTypeParameterDeclaration([
            t.tsTypeParameter(
                t.tsTypeReference(
                    t.identifier('Exact'),
                    t.tsTypeParameterInstantiation([
                        t.tsTypeReference(
                            t.identifier(partialName),
                            t.tsTypeParameterInstantiation([
                                t.tsTypeReference(
                                    t.identifier(name)
                                )
                            ])
                        ),
                        t.tsTypeReference(
                            t.identifier('I')
                        )
                    ])
                ),
                null,
                'I'
            )
        ]);
        // param
        param = identifier(
            varName,
            t.tsTypeAnnotation(
                t.tsTypeReference(
                    t.identifier('I')
                )
            )
        );
    } else {
        // param
        param = identifier(
            varName,
            t.tsTypeAnnotation(
                t.tsTypeReference(
                    t.identifier(partialName),
                    t.tsTypeParameterInstantiation(
                        [
                            t.tsTypeReference(
                                t.identifier(name)
                            )
                        ]
                    )
                )
            )
        );
    }

    return objectMethod(
        'method',
        t.identifier('fromPartial'),
        [
            param
        ],
        t.blockStatement([

            // init 
            t.variableDeclaration(
                'const',
                [
                    t.variableDeclarator(
                        t.identifier('message'),
                        t.callExpression(
                            t.identifier(getBaseCreateTypeFuncName(name)),
                            []
                        )
                    )
                ]
            ),

            ...fields,

            // RETURN 
            t.returnStatement(
                t.identifier('message')
            )

        ]),
        false,
        false,
        false,
        t.tsTypeAnnotation(
            t.tsTypeReference(
                t.identifier(name)
            )
        ),
        typeParameters
    )
};

