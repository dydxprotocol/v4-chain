import * as t from '@babel/types';

import { MessageSchema, Interface } from "../types"
import { FieldTypeAsts } from "../utils";

// export const mutations: Mutation[] = [
//     { methodName: 'joinPool', typeUrl: '/cosmos.pools.transfer.v1.MsgJoinPool', TypeName: 'MsgJoinPool' },
//     { methodName: 'exitPool', typeUrl: '/cosmos.pools.transfer.v1.MsgExitPool', TypeName: 'MsgExitPool' }
// ];

// export const enums: Enum[] = [
//     {
//         name: 'VoteOption',
//         node: null,
//         filename: 'myfile.ts',
//         to: {
//             convertType: 'to',
//             funcName: 'voteOptionToJSON',
//             type: 'any'
//         },
//         from: {
//             convertType: 'from',
//             funcName: 'voteOptionFromJSON',
//             type: 'string'
//         }
//     }
// ];

export const interfaces: Interface[] = [
    {
        name: 'Description',
        fields: [
            {
                type: 'string',
                name: 'moniker',
                node: FieldTypeAsts.string()
            },
            {
                type: 'string',
                name: 'identity',
                node: FieldTypeAsts.string()
            },
            {
                type: 'string',
                name: 'website',
                node: FieldTypeAsts.string()
            },
            {
                type: 'string',
                name: 'securityContact',
                node: FieldTypeAsts.string()
            },
            {
                type: 'string',
                name: 'details',
                node: FieldTypeAsts.string()
            },
            {
                type: 'CommissionRate',
                name: 'superNested',
                node: t.tsTypeReference(t.identifier('CommissionRate'))
            },
            {
                type: 'CommissionRate[]',
                name: 'manyComissions',
                node: t.tsArrayType(t.tsTypeReference(t.identifier('CommissionRate')))
            }
        ]
    },
    {
        name: 'CommissionRate',
        fields: [
            {
                type: 'string',
                name: 'rate',
                node: FieldTypeAsts.string()
            },
            {
                type: 'string',
                name: 'maxRate',
                node: FieldTypeAsts.string()
            },
            {
                type: 'string',
                name: 'maxChangeRate',
                node: FieldTypeAsts.string()
            },
            {
                type: 'SpecialType[]',
                name: 'specialPropertyHere',
                node: t.tsArrayType(t.tsTypeReference(t.identifier('SpecialType')))
            }
        ]
    },
    {
        name: 'SpecialType',
        fields: [
            {
                type: 'string',
                name: 'rate',
                node: FieldTypeAsts.string()
            },
            {
                type: 'string',
                name: 'maxRate',
                node: FieldTypeAsts.string()
            },
            {
                type: 'string',
                name: 'maxChangeRate',
                node: FieldTypeAsts.string()
            }
        ]
    }
];

export const schemata: MessageSchema[] = [
    {
        typeUrl: '/cosmos.some.MsgThing',
        name: 'MstThing',
        fields: [
            {
                type: 'Long',
                name: 'durationCamelCase',
                node: FieldTypeAsts.Long()
            },
            {
                type: 'Coin[]',
                name: 'coinsCamelCase',
                node: FieldTypeAsts.array('Coin')
            },
            {
                type: 'Duration',
                name: 'someDuration',
                node: FieldTypeAsts.array('Duration')
            },
            {
                type: 'Height',
                name: 'someHeight',
                node: FieldTypeAsts.array('Height')
            },
            {
                type: 'Coin',
                name: 'myCoin',
                node: FieldTypeAsts.Coin()
            },
            {
                type: 'string',
                name: 'camelCaseName',
                node: FieldTypeAsts.string()
            },
            {
                type: 'VoteOption',
                name: 'voteValue',
                node: t.tsTypeReference(t.identifier('VoteOption'))
            },
            {
                type: 'VoteOption[]',
                name: 'previousVotes',
                node: t.tsArrayType(t.tsTypeReference(t.identifier('VoteOption')))
            },
            {
                type: 'Description',
                name: 'desc',
                node: t.tsTypeReference(t.identifier('Description'))
            },
            {
                type: 'CommissionRate',
                name: 'commission',
                node: t.tsTypeReference(t.identifier('CommissionRate'))
            },
            {
                type: 'string',
                name: 'str',
                node: FieldTypeAsts.string()
            }
        ]
    }
];