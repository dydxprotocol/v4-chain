import * as t from '@babel/types';
import {
    BILLION,
    memberExpressionOrIdentifierAminoCaseField,
    shorthandProperty
} from '../../../utils';
import { FromAminoParseField, fromAminoParseField } from './index'
import { protoFieldsToArray } from '../utils';
import { getOneOfs, getFieldOptionality } from '../../proto';
import { ProtoField } from '@osmonauts/types';

export const fromAmino = {
    defaultType(args: FromAminoParseField) {
        if (args.field.name === args.context.aminoCaseField(args.field) && args.scope.length === 1) {
            return shorthandProperty(args.field.name);
        }
        return t.objectProperty(
            t.identifier(args.field.name),
            memberExpressionOrIdentifierAminoCaseField(args.fieldPath, args.context.aminoCaseField)
        );
    },

    string(args: FromAminoParseField) {

        if (args.field.name === args.context.aminoCaseField(args.field) && args.scope.length === 1) {
            return shorthandProperty(args.field.name);
        }
        return t.objectProperty(
            t.identifier(args.field.name),
            memberExpressionOrIdentifierAminoCaseField(args.fieldPath, args.context.aminoCaseField)
        );
    },

    rawBytes(args: FromAminoParseField) {
        args.context.addUtil('toUtf8');
        return t.objectProperty(
            t.identifier(args.field.name),
            t.callExpression(
                t.identifier('toUtf8'),
                [
                    t.callExpression(
                        t.memberExpression(
                            t.identifier('JSON'),
                            t.identifier('stringify')
                        ),
                        [
                            memberExpressionOrIdentifierAminoCaseField(args.fieldPath, args.context.aminoCaseField)
                        ]
                    )
                ]
            )
        );
    },

    wasmByteCode(args: FromAminoParseField) {
        args.context.addUtil('fromBase64');
        return t.objectProperty(
            t.identifier(args.field.name),
            t.callExpression(
                t.identifier('fromBase64'),
                [
                    memberExpressionOrIdentifierAminoCaseField(args.fieldPath, args.context.aminoCaseField)
                ]
            )
        );
    },

    long(args: FromAminoParseField) {
        args.context.addUtil('Long');

        return t.objectProperty(t.identifier(args.field.name),
            t.callExpression(
                t.memberExpression(
                    t.identifier('Long'),
                    t.identifier('fromString')
                ),
                [
                    memberExpressionOrIdentifierAminoCaseField(args.fieldPath, args.context.aminoCaseField)
                ]
            ));
    },

    duration(args: FromAminoParseField) {
        const durationFormat = args.context.pluginValue('prototypes.typingsFormat.duration');
        switch (durationFormat) {
            case 'duration':
            // TODO duration amino type
            case 'string':
            default:
                return fromAmino.durationString(args);
        }
    },

    durationString(args: FromAminoParseField) {
        args.context.addUtil('Long');

        const value = t.objectExpression(
            [
                t.objectProperty(t.identifier('seconds'), t.callExpression(
                    t.memberExpression(t.identifier('Long'), t.identifier('fromNumber')), [
                    t.callExpression(
                        t.memberExpression(
                            t.identifier('Math'),
                            t.identifier('floor')
                        ),
                        [
                            t.binaryExpression('/',
                                t.callExpression(
                                    t.identifier('parseInt'),
                                    [
                                        memberExpressionOrIdentifierAminoCaseField(args.fieldPath, args.context.aminoCaseField)
                                    ]
                                ),
                                BILLION
                            )
                        ]
                    )
                ]
                )),
                t.objectProperty(
                    t.identifier('nanos'),
                    t.binaryExpression('%',
                        t.callExpression(
                            t.identifier('parseInt'),
                            [
                                memberExpressionOrIdentifierAminoCaseField(args.fieldPath, args.context.aminoCaseField)
                            ]
                        ),
                        BILLION
                    )
                )
            ]
        );
        return t.objectProperty(t.identifier(args.field.name), value);
    },

    height(args: FromAminoParseField) {
        args.context.addUtil('Long');

        return t.objectProperty(
            t.identifier(args.field.name),
            t.conditionalExpression(
                memberExpressionOrIdentifierAminoCaseField(args.fieldPath, args.context.aminoCaseField),
                t.objectExpression([
                    t.objectProperty(t.identifier('revisionHeight'),
                        t.callExpression(
                            t.memberExpression(t.identifier('Long'), t.identifier('fromString')),
                            [
                                t.logicalExpression(
                                    '||',
                                    t.memberExpression(
                                        memberExpressionOrIdentifierAminoCaseField(args.fieldPath, args.context.aminoCaseField),
                                        t.identifier(args.context.aminoCasingFn('revision_height'))
                                    ),
                                    t.stringLiteral('0')
                                ),
                                t.booleanLiteral(true)
                            ])
                    ),
                    //
                    t.objectProperty(t.identifier('revisionNumber'),
                        t.callExpression(
                            t.memberExpression(t.identifier('Long'), t.identifier('fromString')),
                            [
                                t.logicalExpression(
                                    '||',
                                    t.memberExpression(
                                        memberExpressionOrIdentifierAminoCaseField(args.fieldPath, args.context.aminoCaseField),
                                        t.identifier(args.context.aminoCasingFn('revision_number'))
                                    ),
                                    t.stringLiteral('0')
                                ),
                                t.booleanLiteral(true)
                            ])
                    )
                ]),
                t.identifier('undefined')
            )
        )
    },

    enum({ context, field, currentProtoPath, scope, fieldPath, nested, isOptional }: FromAminoParseField) {
        const enumFunction = context.lookupEnumFromJson(field, currentProtoPath);
        const value = t.callExpression(
            t.identifier(enumFunction), [
            memberExpressionOrIdentifierAminoCaseField(fieldPath, context.aminoCaseField)
        ]);
        return t.objectProperty(t.identifier(field.name), value);
    },

    enumArray({ context, field, currentProtoPath, scope, fieldPath, nested, isOptional }: FromAminoParseField) {
        const enumFunction = context.lookupEnumFromJson(field, currentProtoPath);
        const value = t.callExpression(
            t.memberExpression(
                memberExpressionOrIdentifierAminoCaseField(fieldPath, context.aminoCaseField),
                t.identifier('map')
            ),
            [
                t.arrowFunctionExpression(
                    [
                        t.identifier('el')
                    ],
                    t.callExpression(
                        t.identifier(enumFunction),
                        [
                            t.identifier('el')
                        ]
                    )
                )
            ]
        );
        return t.objectProperty(t.identifier(field.name), value);
    },

    type({ context, field, currentProtoPath, scope, fieldPath, nested, isOptional }: FromAminoParseField) {
        const parentField = field;
        const Type = context.getTypeFromCurrentPath(field, currentProtoPath);
        const oneOfs = getOneOfs(Type);
        const properties = protoFieldsToArray(Type).map(field => {
            const isOneOf = oneOfs.includes(field.name);
            const isOptional = getFieldOptionality(context, field, isOneOf);

            if (parentField.import) currentProtoPath = parentField.import;
            return fromAminoParseField({
                context,
                field,
                currentProtoPath,
                scope: [...scope],
                fieldPath: [...fieldPath],
                nested: nested + 1,
                isOptional // TODO how to handle nested optionality?
            })
        });
        return t.objectProperty(t.identifier(field.name),
            t.objectExpression(
                properties
            )
        );
    },

    arrayFrom(args: FromAminoParseField) {
        return t.objectProperty(t.identifier(args.field.name),
            t.callExpression(
                t.memberExpression(
                    t.identifier('Array'),
                    t.identifier('from')
                ),
                [
                    memberExpressionOrIdentifierAminoCaseField(args.fieldPath, args.context.aminoCaseField)
                ]
            ));
    },

    typeArray({ context, field, currentProtoPath, scope, fieldPath, nested, isOptional }: FromAminoParseField) {
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

        const parentField = field;
        const Type = context.getTypeFromCurrentPath(field, currentProtoPath);
        const oneOfs = getOneOfs(Type);
        const properties = protoFieldsToArray(Type).map(field => {
            const isOneOf = oneOfs.includes(field.name);
            const isOptional = getFieldOptionality(context, field, isOneOf);

            if (parentField.import) currentProtoPath = parentField.import;

            return fromAminoParseField({
                context,
                field,
                currentProtoPath,
                scope: [variable],
                fieldPath: [varProto],
                nested: nested + 1,
                isOptional // TODO how to handle nested optionality?
            })
        });

        const expr = t.callExpression(
            t.memberExpression(
                memberExpressionOrIdentifierAminoCaseField(fieldPath, context.aminoCaseField),
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

        return t.objectProperty(t.identifier(field.name),
            expr
        );
    },


    scalarArray({ context, field, currentProtoPath, scope, fieldPath, nested, isOptional }: FromAminoParseField, arrayTypeAstFunc: Function) {
        const variable = 'el' + nested;

        const expr = t.callExpression(
            t.memberExpression(
                memberExpressionOrIdentifierAminoCaseField(fieldPath, context.aminoCaseField),
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

        return t.objectProperty(t.identifier(field.name),
            expr
        );
    },

    pubkey(args: FromAminoParseField) {
        args.context.addUtil('toBase64');
        args.context.addUtil('encodeBech32Pubkey');

        return t.objectProperty(
            t.identifier(args.field.name),
            t.callExpression(
                t.identifier('encodeBech32Pubkey'),
                [
                    t.objectExpression([
                        t.objectProperty(
                            t.identifier('type'),
                            t.stringLiteral('tendermint/PubKeySecp256k1')
                        ),
                        t.objectProperty(
                            t.identifier('value'),
                            t.callExpression(
                                t.identifier('toBase64'),
                                [
                                    t.memberExpression(
                                        t.identifier('pubkey'),
                                        t.identifier('value')
                                    )
                                ]
                            )
                        )
                    ]),
                    // TODO how to manage this?
                    // 1. options.prefix
                    // 2. look into prefix and how it's used across chains
                    // 3. maybe AminoConverter is a class and has this.prefix!
                    t.stringLiteral('cosmos')
                ]
            )
        )
    }
};


export const arrayTypes = {
    long(varname: string) {
        return t.callExpression(
            t.memberExpression(t.identifier('Long'), t.identifier('fromString')),
            [
                t.identifier(varname)
            ]
        )
    }
}
