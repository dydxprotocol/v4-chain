import { Tsoa } from '@tsoa/runtime';
/**
 * This will help us do exhaustive matching against only reference types. For example, once you have narrowed the input, you don't then have to check the case where it's a `integer` because it never will be.
 */
export declare function isRefType(metaType: Tsoa.Type): metaType is Tsoa.ReferenceType;
