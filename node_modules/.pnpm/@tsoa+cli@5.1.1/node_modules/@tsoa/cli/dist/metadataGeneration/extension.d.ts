import * as ts from 'typescript';
import { MetadataGenerator } from './metadataGenerator';
import { Tsoa } from '@tsoa/runtime';
export declare function getExtensions(decorators: ts.Identifier[], metadataGenerator: MetadataGenerator): Tsoa.Extension[];
export declare function getExtensionsFromJSDocComments(comments: string[]): Tsoa.Extension[];
