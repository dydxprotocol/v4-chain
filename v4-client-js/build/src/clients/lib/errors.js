"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.WrappedError = exports.CustomError = exports.UserError = exports.BroadcastErrorObject = exports.UnexpectedClientError = void 0;
/**
 * An edge-case was hit in the client that should never have been reached.
 */
class UnexpectedClientError extends Error {
    constructor() {
        super('An unexpected error occurred on the client');
        this.name = 'UnexpectedClientError';
    }
}
exports.UnexpectedClientError = UnexpectedClientError;
/**
 * An error occurred during the broadcasting process.
 */
class BroadcastErrorObject extends Error {
    constructor(message, result) {
        super(message);
        this.name = 'BroadcastError';
        this.result = result;
        this.code = result.code;
        this.codespace = result.codespace;
    }
}
exports.BroadcastErrorObject = BroadcastErrorObject;
/**
 * User error occurred during a client operation.
 */
class UserError extends Error {
    constructor(message) {
        super(message);
        this.name = 'UserError';
    }
}
exports.UserError = UserError;
/**
 * @description Base class for custom errors.
 */
class CustomError extends Error {
    constructor(message) {
        super(message);
        // Set a more specific name. This will show up in e.g. console.log.
        this.name = this.constructor.toString();
    }
}
exports.CustomError = CustomError;
/**
 * @description Base class for a custom error which wraps another error.
 */
class WrappedError extends CustomError {
    constructor(message, originalError) {
        super(message);
        this.originalError = originalError;
    }
}
exports.WrappedError = WrappedError;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZXJyb3JzLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2NsaWVudHMvbGliL2Vycm9ycy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7QUFJQTs7R0FFRztBQUNILE1BQWEscUJBQXNCLFNBQVEsS0FBSztJQUM5QztRQUVFLEtBQUssQ0FBQyw0Q0FBNEMsQ0FBQyxDQUFDO1FBQ3BELElBQUksQ0FBQyxJQUFJLEdBQUcsdUJBQXVCLENBQUM7SUFDdEMsQ0FBQztDQUNGO0FBTkQsc0RBTUM7QUFFRDs7R0FFRztBQUNILE1BQWEsb0JBQXFCLFNBQVEsS0FBSztJQUs3QyxZQUNFLE9BQWUsRUFDZixNQUErQjtRQUUvQixLQUFLLENBQUMsT0FBTyxDQUFDLENBQUM7UUFDZixJQUFJLENBQUMsSUFBSSxHQUFHLGdCQUFnQixDQUFDO1FBQzdCLElBQUksQ0FBQyxNQUFNLEdBQUcsTUFBTSxDQUFDO1FBQ3JCLElBQUksQ0FBQyxJQUFJLEdBQUcsTUFBTSxDQUFDLElBQUksQ0FBQztRQUN4QixJQUFJLENBQUMsU0FBUyxHQUFHLE1BQU0sQ0FBQyxTQUFTLENBQUM7SUFDcEMsQ0FBQztDQUNGO0FBZkQsb0RBZUM7QUFFRDs7R0FFRztBQUNILE1BQWEsU0FBVSxTQUFRLEtBQUs7SUFDbEMsWUFBWSxPQUFlO1FBQ3pCLEtBQUssQ0FBQyxPQUFPLENBQUMsQ0FBQztRQUNmLElBQUksQ0FBQyxJQUFJLEdBQUcsV0FBVyxDQUFDO0lBQzFCLENBQUM7Q0FDRjtBQUxELDhCQUtDO0FBRUQ7O0dBRUc7QUFDSCxNQUFhLFdBQVksU0FBUSxLQUFLO0lBQ3BDLFlBQVksT0FBZTtRQUN6QixLQUFLLENBQUMsT0FBTyxDQUFDLENBQUM7UUFDZixtRUFBbUU7UUFDbkUsSUFBSSxDQUFDLElBQUksR0FBRyxJQUFJLENBQUMsV0FBVyxDQUFDLFFBQVEsRUFBRSxDQUFDO0lBQzFDLENBQUM7Q0FDRjtBQU5ELGtDQU1DO0FBRUQ7O0dBRUc7QUFDSCxNQUFhLFlBQWEsU0FBUSxXQUFXO0lBR3pDLFlBQ0UsT0FBZSxFQUNmLGFBQW9CO1FBRXBCLEtBQUssQ0FBQyxPQUFPLENBQUMsQ0FBQztRQUNmLElBQUksQ0FBQyxhQUFhLEdBQUcsYUFBYSxDQUFDO0lBQ3JDLENBQUM7Q0FDSjtBQVZELG9DQVVDIn0=