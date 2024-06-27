import * as ts from 'typescript';
export declare function getJSDocDescription(node: ts.Node): string | undefined;
export declare function getJSDocComment(node: ts.Node, tagName: string): string | undefined;
export declare function getJSDocComments(node: ts.Node, tagName: string): string[] | undefined;
export declare function getJSDocTagNames(node: ts.Node, requireTagName?: boolean): string[];
export declare function getJSDocTags(node: ts.Node, isMatching: (tag: ts.JSDocTag) => boolean): ts.JSDocTag[];
export declare function isExistJSDocTag(node: ts.Node, isMatching: (tag: ts.JSDocTag) => boolean): boolean;
export declare function commentToString(comment?: string | ts.NodeArray<ts.JSDocText | ts.JSDocLink | ts.JSDocComment>): string | undefined;
