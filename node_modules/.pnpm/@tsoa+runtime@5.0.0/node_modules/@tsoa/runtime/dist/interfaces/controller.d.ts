export declare class Controller {
    private statusCode?;
    private headers;
    setStatus(statusCode: number): void;
    getStatus(): number | undefined;
    setHeader(name: string, value?: string | string[]): void;
    getHeader(name: string): string | string[] | undefined;
    getHeaders(): {
        [name: string]: string | string[] | undefined;
    };
}
