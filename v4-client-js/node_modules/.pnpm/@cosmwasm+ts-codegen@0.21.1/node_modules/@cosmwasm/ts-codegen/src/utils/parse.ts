import babelTraverse from '@babel/traverse';
import { parse, ParserPlugin } from '@babel/parser';

export const parser = (codes) => {

    const hash = {};
    codes.forEach(code => {

        const plugins: ParserPlugin[] = [
            'typescript',
        ];

        const ast = parse(code, {
            sourceType: 'module',
            plugins
        });

        const visitor = visitorFn({
            addType(key, node) {
                hash[key] = node;
            }
        })
        babelTraverse(ast, visitor);
    });

    return hash;

}

const visitorFn = (parser) => ({
    TSTypeAliasDeclaration(path) {
        parser.addType(path.node.id.name, path.parentPath.node);
        // if (path.node.id.name.endsWith('For_Empty')) {
        //     const newName = path.node.id.name.replace(/For_Empty$/, '_for_Empty');
        //     path.parentPath.node.declaration.id.name = newName;
        //     parser.addType(newName, path.parentPath.node);
        // } else {
        //     parser.addType(path.node.id.name, path.parentPath.node);
        // }
    },
    TSInterfaceDeclaration(path) {
        parser.addType(path.node.id.name, path.parentPath.node);
    }
});

