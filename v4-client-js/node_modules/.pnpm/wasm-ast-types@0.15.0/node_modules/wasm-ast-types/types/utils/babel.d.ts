import * as t from '@babel/types';
import { Field, QueryMsg, ExecuteMsg } from '../types';
import { TSTypeAnnotation, TSExpressionWithTypeArguments } from '@babel/types';
export declare const propertySignature: (name: string, typeAnnotation: t.TSTypeAnnotation, optional?: boolean) => {
    type: string;
    key: t.Identifier;
    typeAnnotation: t.TSTypeAnnotation;
    optional: boolean;
};
export declare const identifier: (name: string, typeAnnotation: t.TSTypeAnnotation, optional?: boolean) => t.Identifier;
export declare const tsTypeOperator: (typeAnnotation: t.TSType, operator: string) => t.TSTypeOperator;
export declare const getMessageProperties: (msg: QueryMsg | ExecuteMsg) => any;
export declare const tsPropertySignature: (key: t.Expression, typeAnnotation: t.TSTypeAnnotation, optional: boolean) => t.TSPropertySignature;
export declare const tsObjectPattern: (properties: (t.RestElement | t.ObjectProperty)[], typeAnnotation: t.TSTypeAnnotation) => t.ObjectPattern;
export declare const callExpression: (callee: t.Expression | t.V8IntrinsicIdentifier, _arguments: (t.Expression | t.SpreadElement | t.ArgumentPlaceholder)[], typeParameters: t.TSTypeParameterInstantiation) => t.CallExpression;
export declare const bindMethod: (name: string) => t.ExpressionStatement;
export declare const typedIdentifier: (name: string, typeAnnotation: TSTypeAnnotation, optional?: boolean) => t.Identifier;
export declare const promiseTypeAnnotation: (name: any) => t.TSTypeAnnotation;
export declare const classDeclaration: (name: string, body: any[], implementsExressions?: TSExpressionWithTypeArguments[], superClass?: t.Identifier) => t.ClassDeclaration;
export declare const classProperty: (name: string, typeAnnotation?: TSTypeAnnotation, isReadonly?: boolean, isStatic?: boolean, noImplicitOverride?: boolean) => t.ClassProperty;
export declare const arrowFunctionExpression: (params: (t.Identifier | t.Pattern | t.RestElement)[], body: t.BlockStatement, returnType: t.TSTypeAnnotation, isAsync?: boolean) => t.ArrowFunctionExpression;
export declare const recursiveNamespace: (names: any, moduleBlockBody: any) => any;
export declare const arrayTypeNDimensions: (body: any, n: any) => any;
export declare const FieldTypeAsts: {
    string: () => t.TSStringKeyword;
    array: (type: any) => t.TSArrayType;
    Duration: () => t.TSTypeReference;
    Height: () => t.TSTypeReference;
    Coin: () => t.TSTypeReference;
    Long: () => t.TSTypeReference;
};
export declare const shorthandProperty: (prop: string) => t.ObjectProperty;
export declare const importStmt: (names: string[], path: string) => t.ImportDeclaration;
export declare const importAs: (name: string, importAs: string, importPath: string) => t.ImportDeclaration;
export declare const importAminoMsg: () => t.ImportDeclaration;
export declare const getFieldDimensionality: (field: Field) => {
    typeName: string;
    dimensions: number;
    isArray: boolean;
};
export declare const memberExpressionOrIdentifier: (names: any) => any;
export declare const memberExpressionOrIdentifierSnake: (names: any) => any;
/**
 * If optional, return a conditional, otherwise just the expression
 */
export declare const optionalConditionalExpression: (test: t.Expression, expression: t.Expression, alternate: t.Expression, optional?: boolean) => t.Expression;
export declare const typeRefOrUnionWithUndefined: (identifier: t.Identifier, optional?: boolean) => t.TSType;
export declare const parameterizedTypeReference: (identifier: string, from: t.TSType, omit: string | Array<string>) => t.TSTypeReference;
/**
 * omitTypeReference(t.tsTypeReference(t.identifier('Cw4UpdateMembersMutation'),),'args').....
 * Omit<Cw4UpdateMembersMutation, 'args'>
 */
export declare const omitTypeReference: (from: t.TSType, omit: string | Array<string>) => t.TSTypeReference;
export declare const pickTypeReference: (from: t.TSType, pick: string | Array<string>) => t.TSTypeReference;
