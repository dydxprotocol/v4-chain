import * as t from '@babel/types';
import { arrayTypeNDim } from '../utils';
import { protoFieldsToArray } from '../utils';
import { getTSTypeForAmino } from '../../types';
import { getOneOfs, getFieldOptionality } from '../../proto';
import { RenderAminoField, renderAminoField } from '.';

export const aminoInterface = {
    defaultType(args: RenderAminoField) {
        return t.tsPropertySignature(
            t.identifier(args.context.aminoCaseField(args.field)),
            t.tsTypeAnnotation(getTSTypeForAmino(args.context, args.field))
        );
    },
    string(args: RenderAminoField) {
        return t.tsPropertySignature(
            t.identifier(args.context.aminoCaseField(args.field)),
            t.tsTypeAnnotation(t.tsStringKeyword())
        );
    },
    long(args: RenderAminoField) {
        // longs become strings...
        return t.tsPropertySignature(
            t.identifier(args.context.aminoCaseField(args.field)),
            t.tsTypeAnnotation(t.tSStringKeyword())
        )
    },
    height(args: RenderAminoField) {
        args.context.addUtil('AminoHeight');

        return t.tsPropertySignature(
            t.identifier(args.context.aminoCaseField(args.field)),
            t.tsTypeAnnotation(
                t.tsTypeReference(t.identifier('AminoHeight'))
            )
        );
    },
    duration(args: RenderAminoField) {
        const durationFormat = args.context.pluginValue('prototypes.typingsFormat.duration');
        switch (durationFormat) {
            case 'string':
                return t.tsPropertySignature(
                    t.identifier(args.context.aminoCaseField(args.field)),
                    t.tsTypeAnnotation(t.tsStringKeyword())
                );
            case 'duration':
            default:
                return aminoInterface.type(args);
        }
    },
    timestamp(args: RenderAminoField) {
        const timestampFormat = args.context.pluginValue('prototypes.typingsFormat.timestamp');
        switch (timestampFormat) {
            case 'date':
            // TODO check is date is Date for amino?
            // return t.tsPropertySignature(
            //     t.identifier(args.context.aminoCaseField(args.field)),
            //     t.tsTypeAnnotation(
            //         t.tsTypeReference(t.identifier('Date'))
            //     )
            // );
            case 'timestamp':
            default:
                return aminoInterface.type(args);
        }
    },
    enum(args: RenderAminoField) {
        return t.tsPropertySignature(
            t.identifier(args.context.aminoCaseField(args.field)),
            t.tsTypeAnnotation(t.tSNumberKeyword())
        );
    },
    enumArray(args: RenderAminoField) {
        return t.tsPropertySignature(
            t.identifier(args.context.aminoCaseField(args.field)),
            t.tsTypeAnnotation(arrayTypeNDim(t.tSNumberKeyword(), 1))
        );
    },
    type({ context, field, currentProtoPath, isOptional }: RenderAminoField) {
        const parentField = field;

        const Type = context.getTypeFromCurrentPath(field, currentProtoPath);
        const oneOfs = getOneOfs(Type);
        const properties = protoFieldsToArray(Type).map(field => {
            const isOneOf = oneOfs.includes(field.name);
            const isOptional = getFieldOptionality(context, field, isOneOf);
            // TODO how to handle isOptional from parent to child...
            if (parentField.import) currentProtoPath = parentField.import;
            return renderAminoField({
                context,
                field,
                currentProtoPath,
                isOptional // TODO how to handle nested optionality
            })
        });

        // 
        return t.tsPropertySignature(
            t.identifier(context.aminoCaseField(field)),
            t.tsTypeAnnotation(
                t.tsTypeLiteral(
                    properties
                )
            )
        );
    },
    typeArray({ context, field, currentProtoPath, isOptional }: RenderAminoField) {
        const parentField = field;
        const Type = context.getTypeFromCurrentPath(field, currentProtoPath);

        // TODO how to handle isOptional from parent to child... 
        const oneOfs = getOneOfs(Type);
        const properties = protoFieldsToArray(Type).map(field => {
            const isOneOf = oneOfs.includes(field.name);
            const isOptional = getFieldOptionality(context, field, isOneOf);

            if (parentField.import) currentProtoPath = parentField.import;
            return renderAminoField({
                context,
                field,
                currentProtoPath,
                isOptional // TODO how to handle nested optionality
            });
        });

        // 
        return t.tsPropertySignature(
            t.identifier(context.aminoCaseField(field)),
            t.tsTypeAnnotation(
                arrayTypeNDim(t.tsTypeLiteral(
                    properties
                ), 1)
            )
        );
    },

    array(args: RenderAminoField) {
        // TODO write test case 

        // return t.tsPropertySignature(
        //     t.identifier(options.aminoCasingFn(field.name)),
        //     t.tsTypeAnnotation(
        //         arrayTypeNDim(t.tsTypeLiteral(
        //             properties
        //         ), 1)
        //     )
        // );
        return t.tsPropertySignature(
            t.identifier(args.context.aminoCaseField(args.field)),
            t.tsTypeAnnotation(
                arrayTypeNDim(
                    getTSTypeForAmino(args.context, args.field),
                    1
                )
            )
        );
    }
}