import { makeCommentBlock, memberExpressionOrIdentifier } from "../utils";
import * as t from '@babel/types';
import generate from '@babel/generator';
import { ServiceMutation } from "@osmonauts/types";
import { camel } from "@osmonauts/utils";

export const documentWithTypeUrl = (
    mutations: ServiceMutation[]
) => {
    const path = mutations[0].package.split('.');
    return t.variableDeclaration('const', [
        t.variableDeclarator(
            t.objectPattern(mutations.map(mutation => {
                const obj = t.objectProperty(
                    t.identifier(camel(mutation.methodName)),
                    t.identifier(camel(mutation.methodName)),
                    false,
                    true
                );
                // typeUrl: `/${mutation.package}.${mutation.message}`,
                obj.leadingComments = mutation.comment ? [makeCommentBlock(mutation.comment)] : []
                return obj;
            })),
            memberExpressionOrIdentifier([
                'withTypeUrl', 'MessageComposer', ...(path.reverse())
            ])
        )
    ]);
};

export const documentWithTypeUrlReadme = (
    mutations: ServiceMutation[]
) => {
    const ast = documentWithTypeUrl(mutations);
    const code = generate(ast).code;
    return `
#### \`${mutations[0].package}\` messages

\`\`\`js
${code}
\`\`\`
    `;
};