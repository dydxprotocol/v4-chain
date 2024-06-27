"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.setLogger = exports.logger = exports.NullLogger = void 0;
class NullLogger {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    info(...args) {
        return;
    }
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    error(...args) {
        return;
    }
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    trace(...args) {
        return;
    }
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    warn(...args) {
        return;
    }
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    fatal(...args) {
        return;
    }
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    debug(...args) {
        return;
    }
}
exports.NullLogger = NullLogger;
exports.logger = new NullLogger();
function setLogger(newLogger) {
    exports.logger = newLogger;
}
exports.setLogger = setLogger;
//# sourceMappingURL=logger.js.map