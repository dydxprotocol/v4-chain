export interface Logger {
    error(...args: Array<{}>): void;
    trace(...args: Array<{}>): void;
    debug(...args: Array<{}>): void;
    info(...args: Array<{}>): void;
    warn(...args: Array<{}>): void;
    fatal(...args: Array<{}>): void;
}
export declare class NullLogger implements Logger {
    info(...args: Array<{}>): void;
    error(...args: Array<{}>): void;
    trace(...args: Array<{}>): void;
    warn(...args: Array<{}>): void;
    fatal(...args: Array<{}>): void;
    debug(...args: Array<{}>): void;
}
export declare let logger: NullLogger;
export declare function setLogger(newLogger: Logger): void;
