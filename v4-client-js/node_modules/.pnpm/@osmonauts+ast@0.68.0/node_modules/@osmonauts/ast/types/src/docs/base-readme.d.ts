interface ReadmeArgs {
    imgSrc: string;
    description: string;
    libName: string;
    libPrettyName: string;
    baseModule: string;
    exampleAddr: string;
    signingBaseClient: string;
    chainName: string;
    denom: string;
}
export declare const readme: ({ imgSrc, description, libName, libPrettyName, baseModule, exampleAddr, signingBaseClient, chainName, denom }: ReadmeArgs) => string;
export {};
