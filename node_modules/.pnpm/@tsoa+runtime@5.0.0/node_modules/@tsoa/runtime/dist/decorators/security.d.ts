/**
 * Can be used to indicate that a method requires no security.
 */
export declare function NoSecurity(): Function;
/**
 * @param {name} security name from securityDefinitions
 */
export declare function Security(name: string | {
    [name: string]: string[];
}, scopes?: string[]): Function;
