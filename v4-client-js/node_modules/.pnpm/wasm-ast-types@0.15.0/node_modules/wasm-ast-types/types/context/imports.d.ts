import { RenderContext } from './context';
export interface ImportObj {
    type: 'import' | 'default' | 'namespace';
    name: string;
    path: string;
    importAs?: string;
}
export declare const UTILS: {
    MsgExecuteContract: string;
    MsgExecuteContractEncodeObject: string;
    Coin: string;
    toUtf8: string;
    selectorFamily: string;
    StdFee: string;
    CosmWasmClient: string;
    ExecuteResult: string;
    SigningCosmWasmClient: string;
    useQuery: (context: RenderContext) => {
        type: string;
        path: string;
        name: any;
    };
    UseQueryOptions: (context: RenderContext) => {
        type: string;
        path: string;
        name: any;
    };
    useMutation: (context: RenderContext) => {
        type: string;
        path: string;
        name: any;
    };
    UseMutationOptions: (context: RenderContext) => {
        type: string;
        path: string;
        name: any;
    };
};
export declare const convertUtilsToImportList: (context: RenderContext, utils: string[]) => ImportObj[];
export declare const getImportStatements: (list: ImportObj[]) => any[];
