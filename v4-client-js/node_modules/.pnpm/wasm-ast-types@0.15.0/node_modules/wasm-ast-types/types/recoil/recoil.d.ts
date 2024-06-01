import * as t from '@babel/types';
import { QueryMsg } from '../types';
import { RenderContext } from '../context';
export declare const createRecoilSelector: (context: RenderContext, keyPrefix: string, QueryClient: string, methodName: string, responseType: string) => t.ExportNamedDeclaration;
export declare const createRecoilSelectors: (context: RenderContext, keyPrefix: string, QueryClient: string, queryMsg: QueryMsg) => any;
export declare const createRecoilQueryClientType: () => {
    type: string;
    id: {
        type: string;
        name: string;
    };
    typeAnnotation: {
        type: string;
        members: {
            type: string;
            key: {
                type: string;
                name: string;
            };
            computed: boolean;
            typeAnnotation: {
                type: string;
                typeAnnotation: {
                    type: string;
                };
            };
        }[];
    };
};
export declare const createRecoilQueryClient: (context: RenderContext, keyPrefix: string, QueryClient: string) => t.ExportNamedDeclaration;
