export declare type fieldType = 'Long' | 'Coin' | 'Duration' | 'Height' | string;
export interface Field {
    name: string;
    type: fieldType;
    node: any;
};
export interface Interface {
    name: string;
    fields: Field[]
};

export interface QueryMsg {
    $schema: string;
    title: "QueryMsg";
    oneOf?: any;
    allOf?: any;
    anyOf?: any;
}

export interface ExecuteMsg {
    $schema: string;
    title: "ExecuteMsg" | "ExecuteMsg_for_Empty";
    oneOf?: any;
    allOf?: any;
    anyOf?: any;
}

export interface JSONSchema {
    $ref?: string;
    $schema?: string;
    additionalProperties?: boolean;
    allOf?: JSONSchema[];
    anyOf?: JSONSchema[];
    definitions?: Record<string, JSONSchema>;
    description?: string;
    oneOf?: JSONSchema[];
    properties?: Record<string, JSONSchema>;
    patternProperties?: Record<string, JSONSchema>;
    items?: JSONSchema | JSONSchema[];
    additionalItems?: JSONSchema;
    required?: string[];
    title?: string;
    type?: string;
}