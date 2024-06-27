"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __exportStar = (this && this.__exportStar) || function(m, exports) {
    for (var p in m) if (p !== "default" && !Object.prototype.hasOwnProperty.call(exports, p)) __createBinding(exports, m, p);
};
Object.defineProperty(exports, "__esModule", { value: true });
require("reflect-metadata");
__exportStar(require("./decorators/deprecated"), exports);
__exportStar(require("./decorators/example"), exports);
__exportStar(require("./decorators/parameter"), exports);
__exportStar(require("./decorators/methods"), exports);
__exportStar(require("./decorators/tags"), exports);
__exportStar(require("./decorators/operationid"), exports);
__exportStar(require("./decorators/route"), exports);
__exportStar(require("./decorators/security"), exports);
__exportStar(require("./decorators/extension"), exports);
__exportStar(require("./decorators/middlewares"), exports);
__exportStar(require("./interfaces/controller"), exports);
__exportStar(require("./interfaces/response"), exports);
__exportStar(require("./interfaces/iocModule"), exports);
__exportStar(require("./interfaces/file"), exports);
__exportStar(require("./decorators/response"), exports);
__exportStar(require("./metadataGeneration/tsoa"), exports);
__exportStar(require("./routeGeneration/templateHelpers"), exports);
__exportStar(require("./routeGeneration/tsoa-route"), exports);
__exportStar(require("./utils/assertNever"), exports);
__exportStar(require("./swagger/swagger"), exports);
__exportStar(require("./config"), exports);
__exportStar(require("./routeGeneration/additionalProps"), exports);
//# sourceMappingURL=index.js.map