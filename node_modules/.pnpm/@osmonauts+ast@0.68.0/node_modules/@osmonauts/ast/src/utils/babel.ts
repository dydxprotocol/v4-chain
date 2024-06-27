import * as t from '@babel/types';
import { makeCommentBlock } from './utils';

// TODO move to @osmonauts/utils package

export const commentBlock = (value: string): t.CommentBlock => {
    return {
        type: 'CommentBlock',
        value,
        start: null,
        end: null,
        loc: null
    };
};

export const commentLine = (value: string): t.CommentLine => {
    return {
        type: 'CommentLine',
        value,
        start: null,
        end: null,
        loc: null
    };
};

export function tsMethodSignature(
    key: t.Expression,
    typeParameters: t.TSTypeParameterDeclaration | null | undefined = null,
    parameters: Array<t.Identifier | t.RestElement>,
    typeAnnotation: t.TSTypeAnnotation | null = null,
    trailingComments?: t.Comment[],
    leadingComments?: t.Comment[]
): t.TSMethodSignature {
    const obj = t.tsMethodSignature(
        key, typeParameters, parameters, typeAnnotation
    );
    obj.kind = 'method';
    if (trailingComments && trailingComments.length) {
        obj.trailingComments = trailingComments;
    }
    if (leadingComments && leadingComments.length) {
        obj.leadingComments = leadingComments;
    }
    return obj;
}

export const classMethod = (
    kind: "get" | "set" | "method" | "constructor" | undefined,
    key: t.Identifier | t.StringLiteral | t.NumericLiteral | t.Expression,
    params: Array<t.Identifier | t.Pattern | t.RestElement | t.TSParameterProperty>,
    body: t.BlockStatement,
    returnType?: t.TSTypeAnnotation,
    leadingComments: t.CommentLine[] = [],
    computed: boolean = false,
    _static: boolean = false,
    generator: boolean = false,
    async: boolean = false,
) => {
    const obj = t.classMethod(kind, key, params, body, computed, _static, generator, async)
    if (returnType) {
        obj.returnType = returnType;
    }
    if (leadingComments) {
        obj.leadingComments = leadingComments;
    }
    return obj;
};

export const tsEnumMember = (
    id: t.Identifier | t.StringLiteral,
    initializer?: t.Expression,
    leadingComments?: any[]
) => {
    const obj = t.tsEnumMember(id, initializer);
    obj.leadingComments = leadingComments;
    return obj;
};

export const tsPropertySignature = (
    key: t.Expression,
    typeAnnotation: t.TSTypeAnnotation,
    optional: boolean
) => {
    const obj = t.tsPropertySignature(key, typeAnnotation);
    obj.optional = optional;
    return obj
};

export const functionDeclaration = (
    id: t.Identifier,
    params: (t.Identifier | t.Pattern | t.RestElement)[],
    body: t.BlockStatement,
    generator?: boolean,
    async?: boolean,
    returnType?: t.TSTypeAnnotation
): t.FunctionDeclaration => {
    const func = t.functionDeclaration(id, params, body, generator, async);
    func.returnType = returnType;
    return func;
};
export const callExpression = (
    callee: t.Expression | t.V8IntrinsicIdentifier,
    _arguments: (t.Expression | t.SpreadElement | t.ArgumentPlaceholder)[],
    typeParameters: t.TSTypeParameterInstantiation
) => {
    const callExpr = t.callExpression(callee, _arguments);
    callExpr.typeParameters = typeParameters;
    return callExpr;
};


export const identifier = (name: string, typeAnnotation: t.TSTypeAnnotation, optional: boolean = false) => {
    const type = t.identifier(name);
    type.typeAnnotation = typeAnnotation;
    type.optional = optional;
    return type;
}

export const classDeclaration = (
    id: t.Identifier,
    superClass: t.Expression | null | undefined = null,
    body: t.ClassBody,
    decorators: Array<t.Decorator> | null = null,
    vImplements?: t.TSExpressionWithTypeArguments[]
): t.ClassDeclaration => {
    const obj = t.classDeclaration(id, superClass, body, decorators);
    if (vImplements) {
        obj.implements = vImplements;
    }
    return obj;
}

export function classProperty(
    key: t.Identifier | t.StringLiteral | t.NumericLiteral | t.Expression,
    value: t.Expression | null = null,
    typeAnnotation: t.TypeAnnotation | t.TSTypeAnnotation | t.Noop | null = null,
    decorators: Array<t.Decorator> | null = null,
    computed: boolean = false,
    _static: boolean = false,
    _readonly: boolean = false,
    accessibility?: 'private' | 'protected' | 'public',
    leadingComments: t.CommentLine[] = []
): t.ClassProperty {
    const obj = t.classProperty(key, value, typeAnnotation, decorators, computed, _static);
    if (accessibility) obj.accessibility = accessibility;
    if (_readonly) obj.readonly = _readonly;
    if (leadingComments.length) obj.leadingComments = leadingComments;
    return obj;
};

export const arrowFunctionExpression = (
    params: (t.Identifier | t.Pattern | t.RestElement)[],
    body: t.BlockStatement,
    returnType?: t.TSTypeAnnotation,
    isAsync: boolean = false,
    typeParameters?: t.TypeParameterDeclaration | t.TSTypeParameterDeclaration | t.Noop
) => {
    const func = t.arrowFunctionExpression(params, body, isAsync);
    func.returnType = returnType;
    func.typeParameters = typeParameters;
    return func;
};

export const tsTypeParameterDeclaration = (params: Array<t.TSTypeParameter>): t.TSTypeParameterDeclaration => {
    const obj = t.tsTypeParameterDeclaration(params);
    delete obj.extra;
    return obj;
};

export const objectPattern = (
    properties: (t.RestElement | t.ObjectProperty)[],
    typeAnnotation: t.TSTypeAnnotation
) => {
    const obj = t.objectPattern(properties);
    obj.typeAnnotation = typeAnnotation;
    return obj;
}

export const objectMethod =
    (
        kind: "method" | "get" | "set",
        key: t.Expression,
        params: (t.Identifier | t.Pattern | t.RestElement)[],
        body: t.BlockStatement,
        computed?: boolean,
        generator?: boolean,
        async?: boolean,
        returnType?: t.TSTypeAnnotation | t.TypeAnnotation | t.Noop,
        typeParameters?: t.TypeParameterDeclaration | t.TSTypeParameterDeclaration | t.Noop
    ): t.ObjectMethod => {
        const obj = t.objectMethod(kind, key, params, body, computed, generator, async);
        obj.returnType = returnType;
        obj.typeParameters = typeParameters;
        return obj;
    };

export const objectProperty = (
    key: t.Expression | t.Identifier | t.StringLiteral | t.NumericLiteral | t.BigIntLiteral | t.DecimalLiteral | t.PrivateName,
    value: t.Expression | t.PatternLike,
    computed?: boolean,
    shorthand?: boolean,
    decorators?: Array<t.Decorator> | null,
    leadingComments: t.CommentLine[] = []
): t.ObjectProperty => {
    const obj = t.objectProperty(key, value, computed, shorthand, decorators);
    if (leadingComments.length) obj.leadingComments = leadingComments;
    return obj;
};

export const makeCommentLineWithBlocks = (comment: string): t.CommentLine[] => {
    if (!comment) return [];
    // NOTE using blocks instead of lines here...
    // @ts-ignore
    return [makeCommentBlock(comment)];
}


