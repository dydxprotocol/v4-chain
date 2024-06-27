import { FieldInstance, Location, Meta, ValidationError } from './base';
import { ContextItem } from './context-items';
export declare type Optional = {
    /**
     * Whether a field whose value is `null` or `undefined` is to be considered optional.
     * @default false
     */
    nullable: boolean;
    /**
     * Whether a field whose value is falsy (that is, `0`, `false`, `null`, `undefined` or an empty
     * string) is to be considered optional.
     * @default false
     */
    checkFalsy: boolean;
} | false;
export declare class Context {
    readonly fields: string[];
    readonly locations: Location[];
    readonly stack: ReadonlyArray<ContextItem>;
    readonly optional: Optional;
    readonly message?: any;
    private readonly _errors;
    get errors(): ReadonlyArray<ValidationError>;
    private readonly dataMap;
    constructor(fields: string[], locations: Location[], stack: ReadonlyArray<ContextItem>, optional: Optional, message?: any);
    getData(options?: {
        requiredOnly: boolean;
    }): FieldInstance[];
    addFieldInstances(instances: FieldInstance[]): void;
    setData(path: string, value: any, location: Location): void;
    addError(message: any, value: any, meta: Meta): void;
    addError(message: any, nestedErrors: ValidationError[]): void;
}
export declare type ReadonlyContext = Pick<Context, Exclude<keyof Context, 'setData' | 'addFieldInstances' | 'addError'>>;
