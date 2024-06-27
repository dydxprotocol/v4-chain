import { Sanitizers } from '../chain/sanitizers';
import { Validators } from '../chain/validators';
import { CustomValidator, DynamicMessageCreator, Location, Request } from '../base';
import { ValidationChain } from '../chain';
import { Optional } from '../context';
import { ResultWithContext } from '../chain/context-runner-impl';
declare type ValidatorSchemaOptions<K extends keyof Validators<any>> = true | {
    /**
     * Options to pass to the validator.
     * Not used with custom validators.
     */
    options?: Parameters<Validators<any>[K]> | Parameters<Validators<any>[K]>[0];
    /**
     * The error message if there's a validation error,
     * or a function for creating an error message dynamically.
     */
    errorMessage?: DynamicMessageCreator | any;
    /**
     * Whether the validation should be reversed.
     */
    negated?: boolean;
    /**
     * Whether the validation should bail after running this validator
     */
    bail?: boolean;
    /**
     * Specify a condition upon which this validator should run.
     * Can either be a validation chain, or a custom validator function.
     */
    if?: CustomValidator | ValidationChain;
};
export declare type ValidatorsSchema = {
    [K in keyof Validators<any>]?: ValidatorSchemaOptions<K>;
};
declare type SanitizerSchemaOptions<K extends keyof Sanitizers<any>> = true | {
    /**
     * Options to pass to the sanitizer.
     * Not used with custom sanitizers.
     */
    options?: Parameters<Sanitizers<any>[K]> | Parameters<Sanitizers<any>[K]>[0];
};
export declare type SanitizersSchema = {
    [K in keyof Sanitizers<any>]?: SanitizerSchemaOptions<K>;
};
declare type InternalParamSchema = ValidatorsSchema & SanitizersSchema;
/**
 * Defines a schema of validations/sanitizations for a field
 */
export declare type ParamSchema = InternalParamSchema & {
    /**
     * Which request location(s) the field to validate is.
     * If unset, the field will be checked in every request location.
     */
    in?: Location | Location[];
    /**
     * The general error message in case a validator doesn't specify one,
     * or a function for creating the error message dynamically.
     */
    errorMessage?: DynamicMessageCreator | any;
    /**
     * Whether this field should be considered optional
     */
    optional?: true | {
        options?: Partial<Optional>;
    };
};
/**
 * @deprecated  Only here for v5 compatibility. Please use ParamSchema instead.
 */
export declare type ValidationParamSchema = ParamSchema;
/**
 * Defines a mapping from field name to a validations/sanitizations schema.
 */
export declare type Schema = Record<string, ParamSchema>;
/**
 * @deprecated  Only here for v5 compatibility. Please use Schema instead.
 */
export declare type ValidationSchema = Schema;
/**
 * Creates an express middleware with validations for multiple fields at once in the form of
 * a schema object.
 *
 * @param schema the schema to validate.
 * @param defaultLocations
 * @returns
 */
export declare function checkSchema(schema: Schema, defaultLocations?: Location[]): ValidationChain[] & {
    run: (req: Request) => Promise<ResultWithContext[]>;
};
export {};
