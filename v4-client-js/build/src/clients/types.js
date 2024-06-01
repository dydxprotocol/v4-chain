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
exports.OrderFlags = void 0;
// OrderFlags, just a number in proto, defined as enum for convenience
var OrderFlags;
(function (OrderFlags) {
    OrderFlags[OrderFlags["SHORT_TERM"] = 0] = "SHORT_TERM";
    OrderFlags[OrderFlags["LONG_TERM"] = 64] = "LONG_TERM";
    OrderFlags[OrderFlags["CONDITIONAL"] = 32] = "CONDITIONAL";
})(OrderFlags = exports.OrderFlags || (exports.OrderFlags = {}));
__exportStar(require("./modules/proto-includes"), exports);
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidHlwZXMuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9zcmMvY2xpZW50cy90eXBlcy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7Ozs7OztBQTJCQSxzRUFBc0U7QUFDdEUsSUFBWSxVQUlYO0FBSkQsV0FBWSxVQUFVO0lBQ3BCLHVEQUFjLENBQUE7SUFDZCxzREFBYyxDQUFBO0lBQ2QsMERBQWdCLENBQUE7QUFDbEIsQ0FBQyxFQUpXLFVBQVUsR0FBVixrQkFBVSxLQUFWLGtCQUFVLFFBSXJCO0FBaUVELDJEQUF5QyJ9