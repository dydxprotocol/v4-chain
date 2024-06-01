import { identifier, makeCommentBlock, memberExpressionOrIdentifier } from "../utils";
import generate from '@babel/generator';
import * as t from '@babel/types';
import { ProtoRef, ProtoService, ProtoServiceMethod, ServiceMutation } from "@osmonauts/types";
import { camel } from "@osmonauts/utils";
import { getNestedProto, ProtoStore } from "@osmonauts/proto-parser";
import { ProtoParseContext } from "../encoding";
import { ServiceMethod } from "../registry";

interface DocumentRpcClient {
    service: DocumentService;
    method: ProtoServiceMethod;
    methodName: string;
    asts: t.Statement[]
}

export const documentRpcClient = (
    context: ProtoParseContext,
    service: DocumentService
): DocumentRpcClient[] => {

    const methods = Object.entries(service.svc.methods).reduce((m, [key, method]) => {
        const variable = t.variableDeclaration('const', [
            t.variableDeclarator(
                identifier('request', t.tsTypeAnnotation(
                    t.tsTypeReference(
                        t.identifier(method.requestType)
                    )
                ))
            )
        ]);
        if (method.comment) {
            variable.leadingComments = [makeCommentBlock(method.comment)];
        } else {
            variable.leadingComments = [makeCommentBlock(method.name)];
        }
        const methodName = context.pluginValue('rpcClients.camelCase') ? camel(method.name) : method.name;
        return [
            ...m,
            {
                service,
                method,
                methodName,
                asts: [
                    variable,
                    //
                    t.variableDeclaration('const', [
                        t.variableDeclarator(
                            t.identifier('result'),
                            t.awaitExpression(
                                t.callExpression(
                                    memberExpressionOrIdentifier(
                                        [methodName, ...service.ref.proto.package.split('.').reverse()]
                                    ),
                                    [t.identifier('request')]
                                )
                            )
                        )
                    ])
                ]
            }
        ];
    }, []);
    return methods;

};

interface DocumentService {
    svc: ProtoService;
    ref: ProtoRef;
}
export const documentRpcClients = (
    context: ProtoParseContext,
    myBase: string,
    store: ProtoStore
): DocumentRpcClient[] => {
    const svcs = store.getServices(myBase);
    const services = Object.entries(svcs).reduce((m, [pkg, refs]) => {
        const res = refs.reduce((m2, ref) => {
            const proto = getNestedProto(ref.proto)
            // TODO generic service types...
            if (proto.Query) {
                return [
                    ...m2, { svc: proto.Query, ref }
                ]
            }
            if (proto.Service) {
                return [
                    ...m2, { svc: proto.Service, ref }
                ]
            }
            return m2;
        }, []);
        return [...m, ...res];
    }, [])

    //////
    return services.reduce((m, svc: DocumentService) => {
        return [...m, ...documentRpcClient(context, svc)];
    }, []);
};

const replaceChars = (str: string) => {
    return str.split(' ').map(s => {
        return s.replace(/\W/g, '')
    }).join('-').toLowerCase();
};

export const documentRpcClientsReadme = (
    context: ProtoParseContext,
    myBase: string,
    store: ProtoStore
) => {

    const results = documentRpcClients(context, myBase, store);

    const toc = results.map(res => {
        const pkg = res.service.ref.proto.package
        const slug = replaceChars(`${pkg}.${res.methodName} RPC`);
        return `[\`${pkg}.${res.methodName}()\` RPC](#${slug})`
    });

    const lines = results.map(res => {
        const pkg = res.service.ref.proto.package
        const ast = t.program(res.asts)
        const code = generate(ast).code;

        return `##### \`${pkg}.${res.methodName}()\` RPC
        
${res.method.name}

\`\`\`js
${code}
\`\`\`
`;
    });

    const pkg = results[0].service.ref.proto.package

    return `
## Table of Contents

${toc.join('\n')}

### \`${pkg}\` RPC

${lines.join('\n')}
`;
};