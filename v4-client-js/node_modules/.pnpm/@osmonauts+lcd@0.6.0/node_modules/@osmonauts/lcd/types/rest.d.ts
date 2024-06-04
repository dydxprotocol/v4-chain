export declare class LCDClient {
    restEndpoint: string;
    private instance;
    constructor({ restEndpoint }: {
        restEndpoint: any;
    });
    get<ResponseType = unknown>(endpoint: any, opts?: {}): Promise<ResponseType>;
    post<ResponseType = unknown>(endpoint: any, body?: {}, opts?: {}): Promise<ResponseType>;
}
