"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.fetchMiddlewares = exports.Middlewares = void 0;
const TSOA_MIDDLEWARES = Symbol('@tsoa:middlewares');
/**
 * Helper function to create a decorator
 * that can act as a class and method decorator.
 * @param fn a callback function that accepts
 *           the subject of the decorator
 *           either the constructor or the
 *           method
 * @returns
 */
function decorator(fn) {
    return (...args) => {
        // class decorator
        if (args.length === 1) {
            fn(args[0]);
            // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
        }
        else if (args.length === 3 && args[2].value) {
            // method decorator
            const descriptor = args[2];
            if (descriptor.value) {
                fn(descriptor.value);
            }
        }
    };
}
/**
 * Install middlewares to the Controller or a specific method.
 * @param middlewares
 * @returns
 */
function Middlewares(...mws) {
    return decorator(target => {
        if (mws) {
            const current = fetchMiddlewares(target);
            Reflect.defineMetadata(TSOA_MIDDLEWARES, [...current, ...mws], target);
        }
    });
}
exports.Middlewares = Middlewares;
/**
 * Internal function used to retrieve installed middlewares
 * in controller and methods (used during routes generation)
 * @param target
 * @returns list of middlewares
 */
function fetchMiddlewares(target) {
    return Reflect.getMetadata(TSOA_MIDDLEWARES, target) || [];
}
exports.fetchMiddlewares = fetchMiddlewares;
//# sourceMappingURL=middlewares.js.map