"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Consumes = exports.FormField = exports.UploadedFiles = exports.UploadedFile = exports.Inject = exports.Header = exports.Queries = exports.Query = exports.Path = exports.Request = exports.BodyProp = exports.Body = void 0;
/**
 * Inject http Body
 *  @param {string} [name] properties name in body object
 */
function Body() {
    return () => {
        return;
    };
}
exports.Body = Body;
/**
 * Inject value from body
 *
 * @param {string} [name] The name of the body parameter
 */
function BodyProp(name) {
    return () => {
        return;
    };
}
exports.BodyProp = BodyProp;
/**
 * Inject http request
 */
function Request() {
    return () => {
        return;
    };
}
exports.Request = Request;
/**
 * Inject value from Path
 *
 * @param {string} [name] The name of the path parameter
 */
function Path(name) {
    return () => {
        return;
    };
}
exports.Path = Path;
/**
 * Inject value from query string
 *
 * @param {string} [name] The name of the query parameter
 */
function Query(name) {
    return () => {
        return;
    };
}
exports.Query = Query;
/**
 * Inject all query values in a single object
 */
function Queries() {
    return () => {
        return;
    };
}
exports.Queries = Queries;
/**
 * Inject value from Http header
 *
 * @param {string} [name] The name of the header parameter
 */
function Header(name) {
    return () => {
        return;
    };
}
exports.Header = Header;
/**
 * Mark parameter as manually injected, which will not be generated
 */
function Inject() {
    return () => {
        return;
    };
}
exports.Inject = Inject;
/**
 * Inject uploaded file
 *
 * @param {string} [name] The name of the uploaded file parameter
 */
function UploadedFile(name) {
    return () => {
        return;
    };
}
exports.UploadedFile = UploadedFile;
/**
 * Inject uploaded files
 *
 * @param {string} [name] The name of the uploaded files parameter
 */
function UploadedFiles(name) {
    return () => {
        return;
    };
}
exports.UploadedFiles = UploadedFiles;
/**
 * Inject uploaded files
 *
 * @param {string} [name] The name of the uploaded files parameter
 */
function FormField(name) {
    return () => {
        return;
    };
}
exports.FormField = FormField;
/**
 * Overrides the default media type of request body.
 * Can be used on specific method.
 * Can't be used on controller level.
 *
 * @link https://swagger.io/docs/specification/describing-request-body/
 */
function Consumes(value) {
    return () => {
        return;
    };
}
exports.Consumes = Consumes;
//# sourceMappingURL=parameter.js.map