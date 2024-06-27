"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.safeFromJson = void 0;
function safeFromJson(json) {
    try {
        return JSON.parse(json);
    }
    catch {
        return undefined;
    }
}
exports.safeFromJson = safeFromJson;
//# sourceMappingURL=jsonUtils.js.map