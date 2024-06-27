/**
 * Inject http Body
 *  @param {string} [name] properties name in body object
 */
export declare function Body(): Function;
/**
 * Inject value from body
 *
 * @param {string} [name] The name of the body parameter
 */
export declare function BodyProp(name?: string): Function;
/**
 * Inject http request
 */
export declare function Request(): Function;
/**
 * Inject value from Path
 *
 * @param {string} [name] The name of the path parameter
 */
export declare function Path(name?: string): Function;
/**
 * Inject value from query string
 *
 * @param {string} [name] The name of the query parameter
 */
export declare function Query(name?: string): Function;
/**
 * Inject all query values in a single object
 */
export declare function Queries(): Function;
/**
 * Inject value from Http header
 *
 * @param {string} [name] The name of the header parameter
 */
export declare function Header(name?: string): Function;
/**
 * Mark parameter as manually injected, which will not be generated
 */
export declare function Inject(): Function;
/**
 * Inject uploaded file
 *
 * @param {string} [name] The name of the uploaded file parameter
 */
export declare function UploadedFile(name?: string): Function;
/**
 * Inject uploaded files
 *
 * @param {string} [name] The name of the uploaded files parameter
 */
export declare function UploadedFiles(name?: string): Function;
/**
 * Inject uploaded files
 *
 * @param {string} [name] The name of the uploaded files parameter
 */
export declare function FormField(name?: string): Function;
/**
 * Overrides the default media type of request body.
 * Can be used on specific method.
 * Can't be used on controller level.
 *
 * @link https://swagger.io/docs/specification/describing-request-body/
 */
export declare function Consumes(value: string): Function;
