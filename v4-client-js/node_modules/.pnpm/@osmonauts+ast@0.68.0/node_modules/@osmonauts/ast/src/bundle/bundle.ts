import { extname } from 'path';
import * as t from '@babel/types';
import { TelescopeOptions } from '@osmonauts/types';

export const recursiveModuleBundle = (
    options: TelescopeOptions,
    obj
) => {
    return Object.keys(obj).map(key => {
        if (obj[key]?.__export) {
            // e.g. abci
            // 1. create variable for abci
            // 2. splat ALL _0, parms into abci
            // 3. export that variable

            const nmspc = t.variableDeclaration('const',
                [t.variableDeclarator(
                    t.identifier(key),
                    t.objectExpression(
                        Object.keys(obj[key])
                            .filter(a => a !== '__export')
                            .filter(a => a.startsWith('_'))
                            .map(a => t.spreadElement(t.identifier(a)))
                    )
                )]
            );

            const others = Object.keys(obj[key])
                .filter(a => a !== '__export')
                .filter(a => !a.startsWith('_'));

            if (others.length) {
                if (!options.experimentalGlobalProtoNamespace) {
                    throw new Error('namespace and package not supported, yet.')
                }
            }

            // return nmspc;
            return t.exportNamedDeclaration(nmspc, []);



        } else {
            // you can make a namespace for obj[key]
            // e.g. libs
            return t.exportNamedDeclaration(
                t.tsModuleDeclaration(
                    t.identifier(key),
                    t.tsModuleBlock(recursiveModuleBundle(options, obj[key]))
                )
            )
        }
    });
};

export const importNamespace = (ident: string, path: string) => t.importDeclaration(
    [
        t.importNamespaceSpecifier(t.identifier(ident))
    ],
    t.stringLiteral(path.replace(extname(path), ''))
);