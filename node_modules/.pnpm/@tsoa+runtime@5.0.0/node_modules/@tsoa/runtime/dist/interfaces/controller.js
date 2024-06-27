"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Controller = void 0;
class Controller {
    constructor() {
        this.statusCode = undefined;
        this.headers = {};
    }
    setStatus(statusCode) {
        this.statusCode = statusCode;
    }
    getStatus() {
        return this.statusCode;
    }
    setHeader(name, value) {
        this.headers[name] = value;
    }
    getHeader(name) {
        return this.headers[name];
    }
    getHeaders() {
        return this.headers;
    }
}
exports.Controller = Controller;
//# sourceMappingURL=controller.js.map