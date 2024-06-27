import { RenderOptions } from "wasm-ast-types";
export interface TSBuilderInput {
    contracts: Array<ContractFile | string>;
    outPath: string;
    options?: TSBuilderOptions;
}
export interface BundleOptions {
    enabled?: boolean;
    scope?: string;
    bundleFile?: string;
}
export declare type TSBuilderOptions = {
    bundle?: BundleOptions;
} & RenderOptions;
export interface BuilderFile {
    type: 'type' | 'client' | 'recoil' | 'react-query' | 'message-composer';
    contract: string;
    localname: string;
    filename: string;
}
export interface ContractFile {
    name: string;
    dir: string;
}
export declare class TSBuilder {
    contracts: Array<ContractFile | string>;
    outPath: string;
    options?: TSBuilderOptions;
    protected files: BuilderFile[];
    constructor({ contracts, outPath, options }: TSBuilderInput);
    getContracts(): ContractFile[];
    renderTypes(contract: ContractFile): Promise<void>;
    renderClient(contract: ContractFile): Promise<void>;
    renderRecoil(contract: ContractFile): Promise<void>;
    renderReactQuery(contract: ContractFile): Promise<void>;
    renderMessageComposer(contract: ContractFile): Promise<void>;
    build(): Promise<void>;
    bundle(): Promise<void>;
}
