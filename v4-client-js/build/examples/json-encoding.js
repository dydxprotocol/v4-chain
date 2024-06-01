"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const bignumber_js_1 = __importDefault(require("bignumber.js"));
const long_1 = __importDefault(require("long"));
const helpers_1 = require("../src/lib/helpers");
const longValue = long_1.default.fromInt(1, true);
const objWithLong = {
    longValue,
};
const jsonedWithLong = (0, helpers_1.encodeJson)(objWithLong);
console.log(jsonedWithLong);
const bigNumberValue = (0, bignumber_js_1.default)(1);
const objWithBigNumber = {
    bigNumberValue,
};
const jsonedWithBigNumber = (0, helpers_1.encodeJson)(objWithBigNumber);
console.log(jsonedWithBigNumber);
const buffer = Buffer.from('this is a tést');
const objWithBuffer = {
    buffer,
};
// Object.keys(objWithBuffer).forEach((k) => {
//   console.log(k, objWithBuffer[k]); // ERROR!!
// });
const isBuffer2 = objWithBuffer.buffer instanceof Uint8Array;
console.log(isBuffer2);
const jsonedWithBuffer = (0, helpers_1.encodeJson)(objWithBuffer);
console.log(jsonedWithBuffer);
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoianNvbi1lbmNvZGluZy5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uL2V4YW1wbGVzL2pzb24tZW5jb2RpbmcudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7QUFBQSxnRUFBcUM7QUFDckMsZ0RBQXdCO0FBRXhCLGdEQUFnRDtBQUVoRCxNQUFNLFNBQVMsR0FBRyxjQUFJLENBQUMsT0FBTyxDQUFDLENBQUMsRUFBRSxJQUFJLENBQUMsQ0FBQztBQUN4QyxNQUFNLFdBQVcsR0FBRztJQUNsQixTQUFTO0NBQ1YsQ0FBQztBQUVGLE1BQU0sY0FBYyxHQUFHLElBQUEsb0JBQVUsRUFBQyxXQUFXLENBQUMsQ0FBQztBQUMvQyxPQUFPLENBQUMsR0FBRyxDQUFDLGNBQWMsQ0FBQyxDQUFDO0FBRTVCLE1BQU0sY0FBYyxHQUFHLElBQUEsc0JBQVMsRUFBQyxDQUFDLENBQUMsQ0FBQztBQUNwQyxNQUFNLGdCQUFnQixHQUFHO0lBQ3ZCLGNBQWM7Q0FDZixDQUFDO0FBQ0YsTUFBTSxtQkFBbUIsR0FBRyxJQUFBLG9CQUFVLEVBQUMsZ0JBQWdCLENBQUMsQ0FBQztBQUN6RCxPQUFPLENBQUMsR0FBRyxDQUFDLG1CQUFtQixDQUFDLENBQUM7QUFFakMsTUFBTSxNQUFNLEdBQUcsTUFBTSxDQUFDLElBQUksQ0FBQyxnQkFBZ0IsQ0FBQyxDQUFDO0FBQzdDLE1BQU0sYUFBYSxHQUFHO0lBQ3BCLE1BQU07Q0FDUCxDQUFDO0FBRUYsOENBQThDO0FBQzlDLGlEQUFpRDtBQUNqRCxNQUFNO0FBQ04sTUFBTSxTQUFTLEdBQUcsYUFBYSxDQUFDLE1BQU0sWUFBWSxVQUFVLENBQUM7QUFDN0QsT0FBTyxDQUFDLEdBQUcsQ0FBQyxTQUFTLENBQUMsQ0FBQztBQUV2QixNQUFNLGdCQUFnQixHQUFHLElBQUEsb0JBQVUsRUFBQyxhQUFhLENBQUMsQ0FBQztBQUNuRCxPQUFPLENBQUMsR0FBRyxDQUFDLGdCQUFnQixDQUFDLENBQUMifQ==