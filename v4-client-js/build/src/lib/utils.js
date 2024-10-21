"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.randomInt = randomInt;
exports.generateRandomClientId = generateRandomClientId;
exports.clientIdFromString = clientIdFromString;
exports.sleep = sleep;
const constants_1 = require("../clients/constants");
/**
 * Returns a random integer value between 0 and (n-1).
 */
function randomInt(n) {
    return Math.floor(Math.random() * n);
}
/**
 * Generate a random clientId.
 */
function generateRandomClientId() {
    return randomInt(constants_1.MAX_UINT_32 + 1);
}
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
/**
 * Pauses the execution of the program for a specified time.
 * @param ms - The number of milliseconds to pause the program.
 * @returns A promise that resolves after the specified number of milliseconds.
 */
async function sleep(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidXRpbHMuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9zcmMvbGliL3V0aWxzLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7O0FBS0EsOEJBSUM7QUFLRCx3REFFQztBQU1ELGdEQWFDO0FBT0Qsc0JBRUM7QUE1Q0Qsb0RBQW1EO0FBRW5EOztHQUVHO0FBQ0gsU0FBZ0IsU0FBUyxDQUN2QixDQUFTO0lBRVQsT0FBTyxJQUFJLENBQUMsS0FBSyxDQUFDLElBQUksQ0FBQyxNQUFNLEVBQUUsR0FBRyxDQUFDLENBQUMsQ0FBQztBQUN2QyxDQUFDO0FBRUQ7O0dBRUc7QUFDSCxTQUFnQixzQkFBc0I7SUFDcEMsT0FBTyxTQUFTLENBQUMsdUJBQVcsR0FBRyxDQUFDLENBQUMsQ0FBQztBQUNwQyxDQUFDO0FBRUQ7OztHQUdHO0FBQ0gsU0FBZ0Isa0JBQWtCLENBQ2hDLEtBQWE7SUFFYixJQUFJLElBQUksR0FBVyxDQUFDLENBQUM7SUFDckIsSUFBSSxLQUFLLENBQUMsTUFBTSxLQUFLLENBQUM7UUFBRSxPQUFPLElBQUksQ0FBQztJQUNwQyxLQUFLLElBQUksQ0FBQyxHQUFHLENBQUMsRUFBRSxDQUFDLEdBQUcsS0FBSyxDQUFDLE1BQU0sRUFBRSxDQUFDLEVBQUUsRUFBRSxDQUFDO1FBQ3RDLElBQUksR0FBRyxDQUFDLENBQUMsSUFBSSxJQUFJLENBQUMsQ0FBQyxHQUFHLElBQUksQ0FBQyxHQUFHLEtBQUssQ0FBQyxVQUFVLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxpQ0FBaUM7UUFDcEYsSUFBSSxJQUFJLENBQUMsQ0FBQyxDQUFDLGlDQUFpQztJQUM5QyxDQUFDO0lBRUQsMERBQTBEO0lBQzFELHNEQUFzRDtJQUN0RCxPQUFPLElBQUksR0FBRyxDQUFDLENBQUMsSUFBSSxFQUFFLENBQUMsQ0FBQztBQUMxQixDQUFDO0FBRUQ7Ozs7R0FJRztBQUNJLEtBQUssVUFBVSxLQUFLLENBQUMsRUFBVTtJQUNwQyxPQUFPLElBQUksT0FBTyxDQUFDLENBQUMsT0FBTyxFQUFFLEVBQUUsQ0FBQyxVQUFVLENBQUMsT0FBTyxFQUFFLEVBQUUsQ0FBQyxDQUFDLENBQUM7QUFDM0QsQ0FBQyJ9