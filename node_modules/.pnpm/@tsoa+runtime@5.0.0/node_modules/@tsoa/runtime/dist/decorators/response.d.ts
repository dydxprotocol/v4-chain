import { IsValidHeader } from '../utils/isHeaderType';
import { HttpStatusCodeLiteral, HttpStatusCodeStringLiteral, OtherValidOpenApiHttpStatusCode } from '../interfaces/response';
export declare function SuccessResponse<HeaderType extends IsValidHeader<HeaderType> = {}>(name: string | number, description?: string, produces?: string | string[]): Function;
export declare function Response<ExampleType, HeaderType extends IsValidHeader<HeaderType> = {}>(name: HttpStatusCodeLiteral | HttpStatusCodeStringLiteral | OtherValidOpenApiHttpStatusCode, description?: string, example?: ExampleType, produces?: string | string[]): Function;
/**
 * Inject a library-agnostic responder function that can be used to construct type-checked (usually error-) responses.
 *
 * The type of the responder function should be annotated `TsoaResponse<Status, Data, Headers>` in order to support OpenAPI documentation.
 */
export declare function Res(): Function;
/**
 * Overrides the default media type of response.
 * Can be used on controller level or only for specific method
 *
 * @link https://swagger.io/docs/specification/media-types/
 */
export declare function Produces(value: string): Function;
