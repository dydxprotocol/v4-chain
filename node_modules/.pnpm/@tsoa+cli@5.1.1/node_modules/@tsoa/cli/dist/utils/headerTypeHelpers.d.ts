import { Tsoa } from '@tsoa/runtime';
import { NodeArray, TypeNode } from 'typescript';
import { MetadataGenerator } from '../metadataGeneration/metadataGenerator';
export declare function getHeaderType(typeArgumentNodes: NodeArray<TypeNode> | undefined, index: number, metadataGenerator: MetadataGenerator): Tsoa.HeaderType | undefined;
export declare function isSupportedHeaderDataType(parameterType: Tsoa.Type): parameterType is Tsoa.HeaderType;
