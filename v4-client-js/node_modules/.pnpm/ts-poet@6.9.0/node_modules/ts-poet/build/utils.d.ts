export declare function groupBy<K extends PropertyKey, T, Y = T>(list: T[], fn: (x: T) => K, valueFn?: (x: T) => Y): Record<K, Y[]>;
export declare function last<T>(list: T[]): T | undefined;
