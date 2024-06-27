import * as t from '@babel/types';
import { snake } from "case";
import { Field, QueryMsg, ExecuteMsg } from '../types';
import { TSTypeAnnotation, TSExpressionWithTypeArguments } from '@babel/types';

// t.TSPropertySignature - kind?
export const propertySignature = (
    name: string,
    typeAnnotation: t.TSTypeAnnotation,
    optional: boolean = false
) => {
    return {
        type: 'TSPropertySignature',
        key: t.identifier(name),
        typeAnnotation,
        optional
    }
};

export const identifier = (name: string, typeAnnotation: t.TSTypeAnnotation, optional: boolean = false): t.Identifier => {
    const type = t.identifier(name);
    type.typeAnnotation = typeAnnotation;
    type.optional = optional;
    return type;
};

export const tsTypeOperator = (typeAnnotation: t.TSType, operator: string) => {
    const obj = t.tsTypeOperator(typeAnnotation);
    obj.operator = operator;
    return obj;
};

export const getMessageProperties = (msg: QueryMsg | ExecuteMsg) => {
    if (msg.anyOf) return msg.anyOf;
    if (msg.oneOf) return msg.oneOf;
    if (msg.allOf) return msg.allOf;
    return [];
}

export const tsPropertySignature = (
    key: t.Expression,
    typeAnnotation: t.TSTypeAnnotation,
    optional: boolean
) => {
    const obj = t.tsPropertySignature(key, typeAnnotation);
    obj.optional = optional;
    return obj
};



export const tsObjectPattern = (
    properties: (t.RestElement | t.ObjectProperty)[],
    typeAnnotation: t.TSTypeAnnotation
) => {
    const obj = t.objectPattern(
        properties
    );
    obj.typeAnnotation = typeAnnotation;
    return obj;
}

export const callExpression = (
    callee: t.Expression | t.V8IntrinsicIdentifier,
    _arguments: (t.Expression | t.SpreadElement | t.ArgumentPlaceholder)[],
    typeParameters: t.TSTypeParameterInstantiation
) => {
    const callExpr = t.callExpression(callee, _arguments);
    callExpr.typeParameters = typeParameters;
    return callExpr;
};

export const bindMethod = (name: string) => {
    return t.expressionStatement(
        t.assignmentExpression('=', t.memberExpression(
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
    )
}

export const typedIdentifier = (name: string, typeAnnotation: TSTypeAnnotation, optional: boolean = false) => {
    const type = t.identifier(name);
    type.typeAnnotation = typeAnnotation;
    type.optional = optional;
    return type;
}

export const promiseTypeAnnotation = (name) => {
    return t.tsTypeAnnotation(
        t.tsTypeReference(
            t.identifier('Promise'),
            t.tsTypeParameterInstantiation(
                [
                    t.tsTypeReference(t.identifier(name))
                ]
            )
        )
    );
}

export const classDeclaration = (name: string, body: any[], implementsExressions: TSExpressionWithTypeArguments[] = [], superClass: t.Identifier = null) => {
    const declaration = t.classDeclaration(
        t.identifier(name),
        superClass,
        t.classBody(body)
    );
    if (implementsExressions.length) {
        declaration.implements = implementsExressions;
    }
    return declaration;
};


export const classProperty = (
    name: string,
    typeAnnotation: TSTypeAnnotation = null,
    isReadonly: boolean = false,
    isStatic: boolean = false,
    noImplicitOverride: boolean = false
) => {
    const prop = t.classProperty(t.identifier(name));
    if (isReadonly) prop.readonly = true;
    if (isStatic) prop.static = true;
    if (typeAnnotation) prop.typeAnnotation = typeAnnotation;
    if (noImplicitOverride) prop.override = true;
    return prop;
};

export const arrowFunctionExpression = (
    params: (t.Identifier | t.Pattern | t.RestElement)[],
    body: t.BlockStatement,
    returnType: t.TSTypeAnnotation,
    isAsync: boolean = false
) => {
    const func = t.arrowFunctionExpression(params, body, isAsync);
    if (returnType) func.returnType = returnType;
    return func;
};


export const recursiveNamespace = (names, moduleBlockBody) => {
    if (!names || !names.length) return moduleBlockBody;
    const name = names.pop();
    const body = [
        t.exportNamedDeclaration(
            t.tsModuleDeclaration(
                t.identifier(name),
                t.tsModuleBlock(recursiveNamespace(names, moduleBlockBody))
            )
        )
    ];
    return body;
};

export const arrayTypeNDimensions = (body, n) => {
    if (!n) return t.tsArrayType(body);
    return t.tsArrayType(
        arrayTypeNDimensions(body, n - 1)
    );
};

export const FieldTypeAsts = {
    string: () => {
        return t.tsStringKeyword()
    },
    array: (type) => {
        return t.tsArrayType(FieldTypeAsts[type]())
    },
    Duration: () => {
        return t.tsTypeReference(t.identifier('Duration'))
    },
    Height: () => {
        return t.tsTypeReference(t.identifier('Height'))
    },
    Coin: () => {
        return t.tsTypeReference(t.identifier('Coin'))
    },
    Long: () => {
        return t.tsTypeReference(t.identifier('Long'))
    }
};

export const shorthandProperty = (prop: string) => {
    return t.objectProperty(t.identifier(prop), t.identifier(prop), false, true);
};

export const importStmt = (names: string[], path: string) => {
    return t.importDeclaration(
        names.map(name => t.importSpecifier(t.identifier(name), t.identifier(name))),
        t.stringLiteral(path));
};

export const importAs = (name: string, importAs: string, importPath: string) => {
    return t.importDeclaration(
        [
            t.importSpecifier(
                t.identifier(importAs),
                t.identifier(name)
            )
        ],
        t.stringLiteral(importPath)
    )
};

export const importAminoMsg = () => {
    return importStmt(['AminoMsg'], '@cosmjs/amino');
};

export const getFieldDimensionality = (field: Field) => {
    let typeName = field.type;
    const isArray = typeName.endsWith('[]');
    let dimensions = 0;
    if (isArray) {
        dimensions = typeName.match(/\[\]/g).length - 1;
        typeName = typeName.replace(/\[\]/g, '');
    }
    return {
        typeName,
        dimensions,
        isArray
    };
}

export const memberExpressionOrIdentifier = (names) => {
    if (names.length === 1) {
        return t.identifier(names[0])
    }
    if (names.length === 2) {
        const [b, a] = names;
        return t.memberExpression(
            t.identifier(a),
            t.identifier(b)
        );
    }
    const [name, ...rest] = names;

    return t.memberExpression(
        memberExpressionOrIdentifier(rest),
        t.identifier(name)
    )
};

export const memberExpressionOrIdentifierSnake = (names) => {
    if (names.length === 1) {
        return t.identifier(snake(names[0]))
    }
    if (names.length === 2) {
        const [b, a] = names;
        return t.memberExpression(
            t.identifier(snake(a)),
            t.identifier(snake(b))
        );
    }
    const [name, ...rest] = names;

    return t.memberExpression(
        memberExpressionOrIdentifierSnake(rest),
        t.identifier(snake(name))
    )
};

/**
 * If optional, return a conditional, otherwise just the expression
 */
export const optionalConditionalExpression = (test: t.Expression, expression: t.Expression, alternate: t.Expression, optional: boolean = false): t.Expression => {
    return optional
        ? t.conditionalExpression(
            test,
            expression,
            alternate
        )
        : expression
}

export const typeRefOrUnionWithUndefined = (identifier: t.Identifier, optional: boolean = false): t.TSType => {
    const typeReference = t.tsTypeReference(identifier)
    return optional
        ? t.tsUnionType([
            typeReference,
            t.tsUndefinedKeyword()
        ])
        : typeReference
}

export const parameterizedTypeReference = (identifier: string, from: t.TSType, omit: string | Array<string>) => {
    return t.tsTypeReference(
        t.identifier(identifier),
        t.tsTypeParameterInstantiation([
            from,
            typeof omit === 'string'
                ? t.tsLiteralType(t.stringLiteral(omit))
                : t.tsUnionType(omit.map(o => t.tsLiteralType(t.stringLiteral(o))))
        ])
    )
}

/**
 * omitTypeReference(t.tsTypeReference(t.identifier('Cw4UpdateMembersMutation'),),'args').....
 * Omit<Cw4UpdateMembersMutation, 'args'>
 */
export const omitTypeReference = (from: t.TSType, omit: string | Array<string>) => {
    return parameterizedTypeReference('Omit', from, omit)
}


export const pickTypeReference = (from: t.TSType, pick: string | Array<string>) => {
    return parameterizedTypeReference('Pick', from, pick)
}
