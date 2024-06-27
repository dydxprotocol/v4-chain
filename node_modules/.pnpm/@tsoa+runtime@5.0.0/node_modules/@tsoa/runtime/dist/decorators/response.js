"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Produces = exports.Res = exports.Response = exports.SuccessResponse = void 0;
function SuccessResponse(name, description, produces) {
    return () => {
        return;
    };
}
exports.SuccessResponse = SuccessResponse;
function Response(name, description, example, produces) {
    return () => {
        return;
    };
}
exports.Response = Response;
/**
 * Inject a library-agnostic responder function that can be used to construct type-checked (usually error-) responses.
 *
 * The type of the responder function should be annotated `TsoaResponse<Status, Data, Headers>` in order to support OpenAPI documentation.
 */
function Res() {
    return () => {
        return;
    };
}
exports.Res = Res;
/**
 * Overrides the default media type of response.
 * Can be used on controller level or only for specific method
 *
 * @link https://swagger.io/docs/specification/media-types/
 */
function Produces(value) {
    return () => {
        return;
    };
}
exports.Produces = Produces;
//# sourceMappingURL=response.js.map