import * as t from '@babel/types';
import { arrowFunctionExpression, classDeclaration, classMethod, classProperty, cleanComment, commentBlock, identifier, tsMethodSignature } from '../../../utils';
import { ProtoService, ProtoServiceMethod } from '@osmonauts/types';
import { GenericParseContext } from '../../../encoding';
import { camel } from '@osmonauts/utils';
const cleanType = (ResponseType: string) => {
    // MARKED AS NOT DRY [google.protobuf names]
    // TODO some have google.protobuf.Any shows up... figure out the better way to handle this
    if (/\./.test(ResponseType)) {
        ResponseType = ResponseType.split('.')[ResponseType.split('.').length - 1];
    }

    return ResponseType;
}
const returnReponseType = (ResponseType: string) => {

    ResponseType = cleanType(ResponseType);

    return t.tsTypeAnnotation(
        t.tsTypeReference(
            t.identifier('Promise'),
            t.tsTypeParameterInstantiation(
                [
                    t.tsTypeReference(
                        t.identifier(ResponseType)
                        // t.identifier(ResponseType + 'SDKType')
                    )
                ]
            )
        )
    );
};

const rpcMethodDefinition = (
    name: string,
    svc: ProtoServiceMethod,
    trailingComments?: t.CommentBlock[],
    leadingComments?: t.CommentBlock[]
) => {

    const requestType = svc.requestType;
    const responseType = svc.responseType;

    const fieldNames = Object.keys(svc.fields ?? {})
    const hasParams = fieldNames.length > 0;

    let optional = false;
    // // if no params, then let's default to empty object for cleaner API
    if (!hasParams) {
        optional = true;
    } else if (hasParams && fieldNames.length === 1 && fieldNames.includes('pagination')) {
        // if only argument "required" is pagination
        // also default to empty
        optional = true;
    }

    const methodArgs: t.Identifier = identifier(
        'request',
        t.tsTypeAnnotation(
            t.tsTypeReference(
                t.identifier(requestType)
            )
        ),
        optional
    );

    return tsMethodSignature(
        t.identifier(name),
        null,
        [
            methodArgs
        ],
        returnReponseType(responseType),
        trailingComments,
        leadingComments
    );
}

// this.Accounts = this.Accounts.bind(this);
// MARKED AS NOT DRY (used in rpc/lcd)
const bindThis = (name: string) => {
    return t.expressionStatement(
        t.assignmentExpression(
            '=',
            t.memberExpression(
                t.thisExpression(),
                t.identifier(name)
            ),
            t.callExpression(
                t.memberExpression(
                    t.memberExpression(
                        t.thisExpression(),
                        t.identifier(name)
                    ),
                    t.identifier('bind')
                ),
                [
                    t.thisExpression()
                ]
            )
        )
    );
};

// const data = QueryAccountsRequest.encode(request).finish();
const encodeData = (name: string) => {
    return t.variableDeclaration(
        'const',
        [
            t.variableDeclarator(
                t.identifier('data'),
                t.callExpression(
                    t.memberExpression(
                        t.callExpression(
                            t.memberExpression(
                                t.identifier(name),
                                t.identifier('encode')
                            ),
                            [
                                t.identifier('request')
                            ]
                        ),
                        t.identifier('finish')
                    ),
                    []
                )
            )
        ]
    )
};

// const promise = this.rpc.request("cosmos.auth.v1beta1.Query", "Accounts", data);
const promiseRequest = (name: string, packageImportName: string) => {

    name = cleanType(name);

    return t.variableDeclaration(
        'const',
        [
            t.variableDeclarator(
                t.identifier('promise'),
                t.callExpression(
                    t.memberExpression(
                        t.memberExpression(
                            t.thisExpression(),
                            t.identifier('rpc')
                        ),
                        t.identifier('request')
                    ),
                    [
                        t.stringLiteral(packageImportName),
                        t.stringLiteral(name),
                        t.identifier('data')
                    ]
                )
            )
        ]
    );
};

// return promise.then((data) => QueryAccountsResponse.decode(new _m0.Reader(data)));
const returnPromise = (name: string) => {

    name = cleanType(name);

    return t.returnStatement(
        t.callExpression(
            t.memberExpression(
                t.identifier('promise'),
                t.identifier('then')
            ),
            [
                t.arrowFunctionExpression(
                    [
                        t.identifier('data')
                    ],
                    t.callExpression(
                        t.memberExpression(
                            t.identifier(name),
                            t.identifier('decode')
                        ),
                        [
                            t.newExpression(
                                t.memberExpression(
                                    t.identifier('_m0'),
                                    t.identifier('Reader')
                                ),
                                [
                                    t.identifier('data')
                                ]
                            )
                        ]
                    )
                )
            ]
        )
    )
};

const makeComment = (comment: string) => {
    return [{ type: 'CommentBlock', value: ` ${comment} ` }]
}

const rpcClassMethod = (
    context: GenericParseContext,
    name: string,
    msg: string,
    svc: ProtoServiceMethod,
    packageImport: string
) => {

    const requestType = svc.requestType;
    const responseType = svc.responseType;
    const comment = svc.comment ?? svc.name;

    let methodArgs: t.Identifier | t.AssignmentPattern = identifier(
        'request',
        t.tsTypeAnnotation(
            t.tsTypeReference(
                t.identifier(requestType)
            )
        )
    );

    const fieldNames = Object.keys(svc.fields ?? {})
    const hasParams = fieldNames.length > 0;

    // if no params, then let's default to empty object for cleaner API
    if (!hasParams) {
        methodArgs = t.assignmentPattern(
            methodArgs,
            t.objectExpression([])
        )
    } else if (hasParams && fieldNames.length === 1 && fieldNames.includes('pagination')) {
        // if only argument "required" is pagination
        // also default to empty
        methodArgs = t.assignmentPattern(
            methodArgs,
            t.objectExpression([
                t.objectProperty(
                    t.identifier('pagination'),
                    t.identifier('undefined'),
                    false,
                    false
                )
            ])
        )
    }

    const body = t.blockStatement([

        // const data = QueryAccountsRequest.encode(request).finish();
        encodeData(requestType),

        // const promise = this.rpc.request("cosmos.auth.v1beta1.Query", "Accounts", data);
        promiseRequest(msg, packageImport),

        // return promise.then((data) => QueryAccountsResponse.decode(new _m0.Reader(data)));                        
        returnPromise(responseType)

    ]);

    if (context.pluginValue('classesUseArrowFunctions')) {
        return classProperty(
            t.identifier(name),
            arrowFunctionExpression(
                [methodArgs],
                body,
                returnReponseType(responseType),
                true
            ),
            undefined,
            undefined,
            undefined,
            undefined,
            undefined,
            undefined,
            makeComment(comment) as t.CommentLine[],
        );
    }

    return classMethod(
        'method',
        t.identifier(name),
        [
            methodArgs
        ],
        body,
        returnReponseType(responseType)
    );
};

const rpcClassConstructor = (
    context: GenericParseContext,
    methods: string[]
) => {

    let bound = [];
    if (!context.pluginValue('classesUseArrowFunctions')) {
        bound = methods.map(method => bindThis(method));
    }

    return classMethod(
        'constructor',
        t.identifier('constructor'),
        [
            identifier(
                'rpc',
                t.tsTypeAnnotation(
                    t.tsTypeReference(
                        t.identifier('Rpc')
                    )
                )
            )
        ],
        t.blockStatement([
            t.expressionStatement(
                t.assignmentExpression(
                    '=',
                    t.memberExpression(
                        t.thisExpression(),
                        t.identifier('rpc')
                    ),
                    t.identifier('rpc')
                )
            ),
            /// methods
            ...bound
        ])
    );
};

const ensureOneSpace = (str) => {
    if (/^[\s\n\t]+/.test(str)) return str;
    return ` ${str}`;
}
const ensureOneSpaceEnd = (str) => {
    if (/[\s\n\t]$/.test(str)) return str;
    return `${str} `;
}

const processRpcComment = (e: ProtoServiceMethod) => {
    const comment = e.comment;
    if (!comment) return '';
    if (!/[\n]+/.test(comment)) {
        return `*${ensureOneSpaceEnd(ensureOneSpace(cleanComment(comment)))}`
    }
    let lines = comment.split('\n');
    lines = ['*', ...lines, ' '];
    const comments = lines.map((line, i) => {
        if (i == 0) return line;
        if (i == (lines.length - 1)) return cleanComment(line);
        return ` *${ensureOneSpace(cleanComment(line))}`
    });
    return comments.join('\n');
};

export const createRpcClientInterface = (
    context: GenericParseContext,
    service: ProtoService
) => {
    const camelRpcMethods = context.pluginValue('rpcClients.camelCase');
    const keys = Object.keys(service.methods ?? {});
    const methods = keys
        .map((key) => {
            const method = service.methods[key];
            const name = camelRpcMethods ? camel(key) : key;
            const leadingComments = method.comment ? [commentBlock(processRpcComment(method))] : [];
            let trailingComments = [];
            return rpcMethodDefinition(
                name,
                method,
                trailingComments,
                leadingComments
            )
        });

    const obj = t.exportNamedDeclaration(
        t.tsInterfaceDeclaration(
            t.identifier(service.name),
            null,
            [],
            t.tsInterfaceBody(
                [
                    ...methods
                ]
            )
        )
    );
    if (service.comment) {
        obj.leadingComments = [commentBlock(`* ${service.comment} `)];
    }
    return obj;
};

export const getRpcClassName = (service: ProtoService) => {
    return `${service.name}ClientImpl`;
}

export const createRpcClientClass = (
    context: GenericParseContext,
    service: ProtoService
) => {

    context.addUtil('Rpc');
    context.addUtil('_m0');

    const camelRpcMethods = context.pluginValue('rpcClients.camelCase');
    const name = getRpcClassName(service);
    const implementsName = service.name;
    const methodNames = Object.keys(service.methods ?? {})
        .map(key => {
            return camelRpcMethods ? camel(key) : key
        });
    const methods = Object.keys(service.methods ?? {})
        .map(key => {
            const method = service.methods[key];
            const name = camelRpcMethods ? camel(key) : key;
            return rpcClassMethod(
                context,
                name,
                key,
                method,
                context.ref.proto.package + '.' + service.name
            )
        });

    return t.exportNamedDeclaration(
        classDeclaration(
            t.identifier(name),
            null,
            t.classBody([
                classProperty(
                    t.identifier('rpc'),
                    null,
                    t.tsTypeAnnotation(
                        t.tsTypeReference(
                            t.identifier('Rpc')
                        )
                    ),
                    null,
                    false,
                    false,
                    true,
                    'private'
                ),

                // CONSTRUCTOR
                rpcClassConstructor(context, methodNames),

                // METHODS
                ...methods
            ]),
            null,
            [
                t.tsExpressionWithTypeArguments(
                    t.identifier(implementsName)
                )
            ]
        )
    );
};


export const createRpcInterface = (context: GenericParseContext, service: ProtoService) => {
    return t.tsInterfaceDeclaration(
        t.identifier('Rpc'),
        null,
        [],
        t.tsInterfaceBody(
            [
                t.tsMethodSignature(
                    t.identifier('request'),
                    null,
                    [
                        identifier('service',
                            t.tsTypeAnnotation(
                                t.tsStringKeyword()
                            )
                        ),
                        identifier('method',
                            t.tsTypeAnnotation(
                                t.tsStringKeyword()
                            )
                        ),
                        identifier('data',
                            t.tsTypeAnnotation(
                                t.tsTypeReference(
                                    t.identifier('Uint8Array')
                                )
                            )
                        )
                    ],
                    t.tsTypeAnnotation(
                        t.tsTypeReference(
                            t.identifier('Promise'),
                            t.tsTypeParameterInstantiation(
                                [
                                    t.tsTypeReference(
                                        t.identifier('Uint8Array')
                                    )
                                ]
                            )
                        )
                    )
                )
            ]
        )
    )
}