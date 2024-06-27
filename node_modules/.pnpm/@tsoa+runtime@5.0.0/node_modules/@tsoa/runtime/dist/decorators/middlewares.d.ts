type Middleware<T extends Function | Object> = T;
/**
 * Install middlewares to the Controller or a specific method.
 * @param middlewares
 * @returns
 */
export declare function Middlewares<T extends Function | Object>(...mws: Array<Middleware<T>>): ClassDecorator & MethodDecorator;
/**
 * Internal function used to retrieve installed middlewares
 * in controller and methods (used during routes generation)
 * @param target
 * @returns list of middlewares
 */
export declare function fetchMiddlewares<T extends Function | Object>(target: any): Array<Middleware<T>>;
export {};
