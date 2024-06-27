import * as t from '@babel/types';
import { BILLION, memberExpressionOrIdentifier, shorthandProperty } from "../../../utils";
import { protoFieldsToArray } from '../utils';
import { ToAminoParseField, toAminoParseField } from './index'
import { getOneOfs, getFieldOptionality } from '../../proto';
import { ProtoField } from '@osmonauts/types';

export const toAmino = {
    defaultType(args: ToAminoParseField) {
        if (args.field.name === args.context.aminoCaseField(args.field) && args.scope.length === 1) {
            return shorthandProperty(args.field.name);
        }
        return t.objectProperty(t.identifier(args.context.aminoCaseField(args.field)), memberExpressionOrIdentifier(args.scope))
    },

    long(args: ToAminoParseField) {
        return t.objectProperty(t.identifier(args.context.aminoCaseField(args.field)),
            t.callExpression(
                t.memberExpression(memberExpressionOrIdentifier(args.scope), t.identifier('toString')),
                [])
        )
    },

    string(args: ToAminoParseField) {
        if (args.field.name === args.context.aminoCaseField(args.field) && args.scope.length === 1) {
            return shorthandProperty(args.field.name);
        }
        return t.objectProperty(t.identifier(args.context.aminoCaseField(args.field)), memberExpressionOrIdentifier(args.scope))
    },

    rawBytes(args: ToAminoParseField) {
        args.context.addUtil('fromUtf8');
        return t.objectProperty(
            t.identifier(args.context.aminoCaseField(args.field)),
            t.callExpression(
                t.memberExpression(
                    t.identifier('JSON'),
                    t.identifier('parse')
                ),
                [
                    t.callExpression(
                        t.identifier('fromUtf8'),
                        [
                            memberExpressionOrIdentifier(args.scope)
                        ]
                    )
                ]
            )

        );
    },

    wasmByteCode(args: ToAminoParseField) {
        args.context.addUtil('toBase64');
        return t.objectProperty(
            t.identifier(args.context.aminoCaseField(args.field)),
            t.callExpression(
                t.identifier('toBase64'),
                [
                    memberExpressionOrIdentifier(args.scope)
                ]
            )

        );
    },

    duration(args: ToAminoParseField) {
        const durationFormat = args.context.pluginValue('prototypes.typingsFormat.duration');
        switch (durationFormat) {
            case 'duration':
            // TODO duration amino type
            case 'string':
            default:
                return toAmino.durationString(args);
        }
    },

    durationString(args: ToAminoParseField) {
        const exp = t.binaryExpression(
            '*',
            memberExpressionOrIdentifier(args.scope),
            BILLION
        );
        exp.extra = { parenthesized: true };
        const value = t.callExpression(
            t.memberExpression(
                exp,
                t.identifier('toString')
            ),
            []
        )
        return t.objectProperty(t.identifier(args.context.aminoCaseField(args.field)), value);
    },

    height(args: ToAminoParseField) {
        args.context.addUtil('omitDefault');

        const value = t.objectExpression(
            [
                t.objectProperty(
                    t.identifier(args.context.aminoCasingFn('revision_height')),
                    t.optionalCallExpression(
                        t.optionalMemberExpression(
                            t.callExpression(
                                t.identifier('omitDefault'),
                                [
                                    t.memberExpression(
                                        memberExpressionOrIdentifier(args.scope),
                                        t.identifier('revisionHeight')
                                    )
                                ]
                            ),
                            t.identifier('toString'),
                            false,
                            true
                        ),
                        [],
                        false
                    )
                ),
                //
                t.objectProperty(
                    t.identifier(args.context.aminoCasingFn('revision_number')),
                    t.optionalCallExpression(
                        t.optionalMemberExpression(
                            t.callExpression(
                                t.identifier('omitDefault'),
                                [
                                    t.memberExpression(
                                        memberExpressionOrIdentifier(args.scope),
                                        t.identifier('revisionNumber')
                                    )
                                ]
                            ),
                            t.identifier('toString'),
                            false,
                            true
                        ),
                        [],
                        false
                    )
                )
            ]
        );

        const cond = t.conditionalExpression(
            memberExpressionOrIdentifier(args.scope),
            value,
            t.objectExpression([])
        );

        return t.objectProperty(t.identifier(args.context.aminoCaseField(args.field)), cond);
    },

    coin(args: ToAminoParseField) {
        args.context.addUtil('Long');
        const value = t.objectExpression([
            t.objectProperty(t.identifier('denom'), t.memberExpression(
                memberExpressionOrIdentifier(args.scope),
                t.identifier('denom'),
            )),
            t.objectProperty(
                t.identifier('amount'),
                t.callExpression(
                    t.memberExpression(
                        t.callExpression(
                            t.memberExpression(
                                t.identifier('Long'),
                                t.identifier('fromValue')
                            ),
                            [
                                t.memberExpression(
                                    memberExpressionOrIdentifier(args.scope),
                                    t.identifier('amount')
                                )
                            ]
                        ),
                        t.identifier('toString')
                    ),
                    []
                )
            )
        ]);
        return t.objectProperty(t.identifier(args.context.aminoCaseField(args.field)), value);
    },

    type({ context, field, currentProtoPath, scope, fieldPath, nested, isOptional }: ToAminoParseField) {
        /// TODO (can this be refactored out? e.g. no recursive calls in this file?)
        /// BEGIN
        const Type = context.getTypeFromCurrentPath(field, currentProtoPath);
        const parentField = field;
        const oneOfs = getOneOfs(Type);
        const properties = protoFieldsToArray(Type).map(field => {
            const isOneOf = oneOfs.includes(field.name);
            const isOptional = getFieldOptionality(context, field, isOneOf);

            if (parentField.import) currentProtoPath = parentField.import;
            return toAminoParseField({
                context,
                field,
                currentProtoPath,
                scope: [...scope],
                fieldPath: [...fieldPath],
                nested,
                isOptional // TODO how to handle nested optionality
            })
        });
        /// END 
        return t.objectProperty(t.identifier(context.aminoCaseField(field)),
            t.objectExpression(
                properties
            )
        );
    },

    typeArray({ context, field, currentProtoPath, scope, fieldPath, nested, isOptional }: ToAminoParseField) {
        //////
        const variable = 'el' + nested;
        const f = JSON.parse(JSON.stringify(field)); // clone
        const varProto: ProtoField = {
            ...f
        };
        varProto.name = variable;
        varProto.options['(telescope:orig)'] = variable;
        varProto.options['(telescope:name)'] = variable;
        varProto.options['(telescope:camel)'] = variable;
        //////

        if (field.parsedType.type !== 'Type') {
            throw new Error('Arrays only support types[Type] right now.');
        }

        const Type = context.getTypeFromCurrentPath(field, currentProtoPath);
        const parentField = field;
        const oneOfs = getOneOfs(Type);

        const properties = protoFieldsToArray(Type).map(field => {
            const isOneOf = oneOfs.includes(field.name);
            const isOptional = getFieldOptionality(context, field, isOneOf);

            if (parentField.import) currentProtoPath = parentField.import;

            return toAminoParseField({
                context,
                field,
                currentProtoPath,
                scope: [variable],
                fieldPath: [varProto],
                nested: nested + 1,
                isOptional // TODO how to handle nested optionality
            });
        });


        const expr = t.callExpression(
            t.memberExpression(
                memberExpressionOrIdentifier(scope),
                t.identifier('map')
            ),
            [
                t.arrowFunctionExpression(
                    [
                        t.identifier(variable)
                    ],
                    t.objectExpression(
                        properties
                    )
                )
            ]
        );

        return t.objectProperty(t.identifier(context.aminoCaseField(field)),
            expr
        );
    },

    scalarArray({ context, field, currentProtoPath, scope, nested, isOptional }: ToAminoParseField, arrayTypeAstFunc: Function) {
        const variable = 'el' + nested;

        const expr = t.callExpression(
            t.memberExpression(
                memberExpressionOrIdentifier(scope),
                t.identifier('map')
            ),
            [
                t.arrowFunctionExpression(
                    [
                        t.identifier(variable)
                    ],
                    arrayTypeAstFunc(variable)
                )
            ]
        );

        return t.objectProperty(t.identifier(context.aminoCaseField(field)),
            expr
        );
    },

    pubkey(args: ToAminoParseField) {
        args.context.addUtil('fromBase64');
        args.context.addUtil('decodeBech32Pubkey');

        return t.objectProperty(
            t.identifier(args.field.name),
            t.objectExpression([
                t.objectProperty(
                    t.identifier('typeUrl'),
                    t.stringLiteral('/cosmos.crypto.secp256k1.PubKey')
                ),
                t.objectProperty(
                    t.identifier('value'),
                    t.callExpression(
                        t.identifier('fromBase64'),
                        [
                            t.memberExpression(
                                t.callExpression(
                                    t.identifier('decodeBech32Pubkey'),
                                    [
                                        t.identifier(args.field.name)
                                    ]
                                ),
                                t.identifier('value')
                            )

                        ]
                    )
                )
            ])
        )

    }
};

export const arrayTypes = {
    long(varname: string) {
        return t.callExpression(
            t.memberExpression(memberExpressionOrIdentifier([varname]), t.identifier('toString')),
            []
        )
    }
}
