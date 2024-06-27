export declare function Extension(_name: string, _value: ExtensionType | ExtensionType[]): Function;
export type ExtensionType = string | number | boolean | null | ExtensionType[] | {
    [name: string]: ExtensionType | ExtensionType[];
};
