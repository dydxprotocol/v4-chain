export interface ServiceInfo {
    methodName: string;
    package: string;
    message: string;
    messageImport: string;
    response: string;
    responseImport: string;
    comment?: string;
}
export interface ServiceMutation extends ServiceInfo {
}
export interface ServiceQuery extends ServiceInfo {
}
