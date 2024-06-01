import { ProtoEnum, ProtoType } from '@osmonauts/types';
import { ProtoParseContext } from '../../context';
import { createProtoEnum } from '../enums';
import { createProtoType } from './interface';

export const createSDKType = (
    context: ProtoParseContext,
    name: string,
    proto: ProtoType
) => {
    return createProtoType(
        context,
        name,
        proto,
        {
            useOriginalCase: true,
            typeNameSuffix: 'SDKType',
        }
    );
};

export const createEnumSDKType = (
    context: ProtoParseContext,
    name: string,
    proto: ProtoEnum
) => {
    return createProtoEnum(
        context,
        name + "SDKType",
        proto
    );
};
