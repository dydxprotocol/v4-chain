"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Swagger = void 0;
// eslint-disable-next-line @typescript-eslint/no-namespace
var Swagger;
(function (Swagger) {
    function isQueryParameter(parameter) {
        return parameter.in === 'query';
    }
    Swagger.isQueryParameter = isQueryParameter;
})(Swagger = exports.Swagger || (exports.Swagger = {}));
//# sourceMappingURL=swagger.js.map