import * as ts from 'typescript';
export declare function getDecorators(node: ts.Node, isMatching: (identifier: ts.Identifier) => boolean): ts.Identifier[];
export declare function getNodeFirstDecoratorName(node: ts.Node, isMatching: (identifier: ts.Identifier) => boolean): string | undefined;
export declare function getNodeFirstDecoratorValue(node: ts.Node, typeChecker: ts.TypeChecker, isMatching: (identifier: ts.Identifier) => boolean): any;
export declare function getDecoratorValues(decorator: ts.Identifier, typeChecker: ts.TypeChecker): any[];
export declare function getSecurites(decorator: ts.Identifier, typeChecker: ts.TypeChecker): any;
export declare function isDecorator(node: ts.Node, isMatching: (identifier: ts.Identifier) => boolean): boolean;
export declare function getPath(decorator: ts.Identifier, typeChecker: ts.TypeChecker): string;
export declare function getProduces(node: ts.Node, typeChecker: ts.TypeChecker): string[];
