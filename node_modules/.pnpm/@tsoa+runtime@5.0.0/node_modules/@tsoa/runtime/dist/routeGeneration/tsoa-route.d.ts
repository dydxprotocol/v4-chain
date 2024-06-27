import { Tsoa } from './../metadataGeneration/tsoa';
/**
 * For Swagger, additionalProperties is implicitly allowed. So use this function to clarify that undefined should be associated with allowing additional properties
 * @param test if this is undefined then you should interpret it as a "yes"
 */
export declare function isDefaultForAdditionalPropertiesAllowed(test: TsoaRoute.RefObjectModelSchema['additionalProperties']): test is undefined;
export declare namespace TsoaRoute {
    interface Models {
        [name: string]: ModelSchema;
    }
    /**
     * This is a convenience type so you can check .properties on the items in the Record without having TypeScript throw a compiler error. That's because this Record can't have enums in it. If you want that, then just use the base interface
     */
    interface RefObjectModels extends TsoaRoute.Models {
        [refNames: string]: TsoaRoute.RefObjectModelSchema;
    }
    interface RefEnumModelSchema {
        dataType: 'refEnum';
        enums: Array<string | number>;
    }
    interface RefObjectModelSchema {
        dataType: 'refObject';
        properties: {
            [name: string]: PropertySchema;
        };
        additionalProperties?: boolean | PropertySchema;
    }
    interface RefTypeAliasModelSchema {
        dataType: 'refAlias';
        type: PropertySchema;
    }
    type ModelSchema = RefEnumModelSchema | RefObjectModelSchema | RefTypeAliasModelSchema;
    type ValidatorSchema = Tsoa.Validators;
    interface PropertySchema {
        dataType?: Tsoa.TypeStringLiteral;
        ref?: string;
        required?: boolean;
        array?: PropertySchema;
        enums?: Array<string | number | boolean | null>;
        type?: PropertySchema;
        subSchemas?: PropertySchema[];
        validators?: ValidatorSchema;
        default?: unknown;
        additionalProperties?: boolean | PropertySchema;
        nestedProperties?: {
            [name: string]: PropertySchema;
        };
    }
    interface ParameterSchema extends PropertySchema {
        name: string;
        in: string;
    }
    interface Security {
        [key: string]: string[];
    }
}
