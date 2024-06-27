interface ServiceOptions {
    "(google.api.http).get"?: string;
    "(google.api.http)"?: {
        post: string;
        body: string;
    };
}
export declare const parseServiceUrl: (options: ServiceOptions) => {
    method: string;
    url: string;
    pathParams: string[];
};
export declare const parseService: (obj: any) => {
    queryParams: string[];
    paramMap: {};
    casing: {};
    method: string;
    url: string;
    pathParams: string[];
};
export {};
