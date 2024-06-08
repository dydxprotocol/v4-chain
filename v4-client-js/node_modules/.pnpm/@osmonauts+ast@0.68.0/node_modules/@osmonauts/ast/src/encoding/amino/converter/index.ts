import * as t from '@babel/types';
import { ProtoType, ProtoRoot } from '@osmonauts/types';
import { toAminoJsonMethod } from '../to-amino-json';
import { fromAminoJsonMethod } from '../from-amino-json';
import { getTypeUrl, typeUrlToAmino } from '../utils';
import { AminoParseContext } from '../../context';

interface AminoConverterItemParams {
    root: ProtoRoot,
    context: AminoParseContext,
    proto: ProtoType
}
export const createAminoConverterItem = (
    {
        root,
        context,
        proto
    }: AminoConverterItemParams
) => {

    const typeUrl = getTypeUrl(root, proto);

    return t.objectProperty(
        t.stringLiteral(typeUrl),
        t.objectExpression(
            [
                t.objectProperty(
                    t.identifier('aminoType'),
                    t.stringLiteral(
                        typeUrlToAmino(context, typeUrl)
                    )
                ),
                t.objectProperty(
                    t.identifier('toAmino'),
                    toAminoJsonMethod({
                        context,
                        proto
                    })
                ),
                t.objectProperty(
                    t.identifier('fromAmino'),
                    fromAminoJsonMethod({
                        context,
                        proto
                    })
                )
            ]
        )
    );
};




interface AminoConverterParams {
    name: string,
    root: ProtoRoot,
    context: AminoParseContext,
    protos: ProtoType[]
}
export const createAminoConverter = (
    {
        name,
        root,
        context,
        protos
    }: AminoConverterParams) => {

    const items = protos.map(proto => {
        return createAminoConverterItem({
            context,
            root,
            proto
        })
    })

    return t.exportNamedDeclaration(t.variableDeclaration('const', [
        t.variableDeclarator(t.identifier(name),
            t.objectExpression(
                items
            ))
    ]));
};

