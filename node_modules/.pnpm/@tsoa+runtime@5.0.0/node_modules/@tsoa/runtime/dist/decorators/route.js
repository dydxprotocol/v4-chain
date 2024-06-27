"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Hidden = exports.Route = void 0;
function Route(name) {
    return () => {
        return;
    };
}
exports.Route = Route;
/**
 * can be used to entirely hide an method from documentation
 */
function Hidden() {
    return () => {
        return;
    };
}
exports.Hidden = Hidden;
//# sourceMappingURL=route.js.map