"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.sleep = exports.clientIdFromString = exports.generateRandomClientId = exports.randomInt = void 0;
const constants_1 = require("../clients/constants");
/**
 * Returns a random integer value between 0 and (n-1).
 */
function randomInt(n) {
    return Math.floor(Math.random() * n);
}
exports.randomInt = randomInt;
/**
 * Generate a random clientId.
 */
function generateRandomClientId() {
    return randomInt(constants_1.MAX_UINT_32 + 1);
}
exports.generateRandomClientId = generateRandomClientId;
/**
 * Deterministically generate a valid clientId from an arbitrary string by performing a
 * quick hashing function on the string.
 */
function clientIdFromString(input) {
    let hash = 0;
    if (input.length === 0)
        return hash;
    for (let i = 0; i < input.length; i++) {
        hash = ((hash << 5) - hash) + input.charCodeAt(i); // eslint-disable-line no-bitwise
        hash |= 0; // eslint-disable-line no-bitwise
    }
    // Bitwise operators covert the value to a 32-bit integer.
    // We must coerce this into a 32-bit unsigned integer.
    return hash + (2 ** 31);
}
exports.clientIdFromString = clientIdFromString;
/**
 * Pauses the execution of the program for a specified time.
 * @param ms - The number of milliseconds to pause the program.
 * @returns A promise that resolves after the specified number of milliseconds.
 */
async function sleep(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}
exports.sleep = sleep;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidXRpbHMuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9zcmMvbGliL3V0aWxzLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7OztBQUFBLG9EQUFtRDtBQUVuRDs7R0FFRztBQUNILFNBQWdCLFNBQVMsQ0FDdkIsQ0FBUztJQUVULE9BQU8sSUFBSSxDQUFDLEtBQUssQ0FBQyxJQUFJLENBQUMsTUFBTSxFQUFFLEdBQUcsQ0FBQyxDQUFDLENBQUM7QUFDdkMsQ0FBQztBQUpELDhCQUlDO0FBRUQ7O0dBRUc7QUFDSCxTQUFnQixzQkFBc0I7SUFDcEMsT0FBTyxTQUFTLENBQUMsdUJBQVcsR0FBRyxDQUFDLENBQUMsQ0FBQztBQUNwQyxDQUFDO0FBRkQsd0RBRUM7QUFFRDs7O0dBR0c7QUFDSCxTQUFnQixrQkFBa0IsQ0FDaEMsS0FBYTtJQUViLElBQUksSUFBSSxHQUFXLENBQUMsQ0FBQztJQUNyQixJQUFJLEtBQUssQ0FBQyxNQUFNLEtBQUssQ0FBQztRQUFFLE9BQU8sSUFBSSxDQUFDO0lBQ3BDLEtBQUssSUFBSSxDQUFDLEdBQUcsQ0FBQyxFQUFFLENBQUMsR0FBRyxLQUFLLENBQUMsTUFBTSxFQUFFLENBQUMsRUFBRSxFQUFFO1FBQ3JDLElBQUksR0FBRyxDQUFDLENBQUMsSUFBSSxJQUFJLENBQUMsQ0FBQyxHQUFHLElBQUksQ0FBQyxHQUFHLEtBQUssQ0FBQyxVQUFVLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxpQ0FBaUM7UUFDcEYsSUFBSSxJQUFJLENBQUMsQ0FBQyxDQUFDLGlDQUFpQztLQUM3QztJQUVELDBEQUEwRDtJQUMxRCxzREFBc0Q7SUFDdEQsT0FBTyxJQUFJLEdBQUcsQ0FBQyxDQUFDLElBQUksRUFBRSxDQUFDLENBQUM7QUFDMUIsQ0FBQztBQWJELGdEQWFDO0FBRUQ7Ozs7R0FJRztBQUNJLEtBQUssVUFBVSxLQUFLLENBQUMsRUFBVTtJQUNwQyxPQUFPLElBQUksT0FBTyxDQUFDLENBQUMsT0FBTyxFQUFFLEVBQUUsQ0FBQyxVQUFVLENBQUMsT0FBTyxFQUFFLEVBQUUsQ0FBQyxDQUFDLENBQUM7QUFDM0QsQ0FBQztBQUZELHNCQUVDIn0=