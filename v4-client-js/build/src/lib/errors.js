"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UserError = exports.BroadcastErrorObject = exports.UnexpectedClientError = void 0;
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
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZXJyb3JzLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL2xpYi9lcnJvcnMudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7O0FBSUE7O0dBRUc7QUFDSCxNQUFhLHFCQUFzQixTQUFRLEtBQUs7SUFDOUM7UUFFRSxLQUFLLENBQUMsNENBQTRDLENBQUMsQ0FBQztRQUNwRCxJQUFJLENBQUMsSUFBSSxHQUFHLHVCQUF1QixDQUFDO0lBQ3RDLENBQUM7Q0FDRjtBQU5ELHNEQU1DO0FBRUQ7O0dBRUc7QUFDSCxNQUFhLG9CQUFxQixTQUFRLEtBQUs7SUFLN0MsWUFDRSxPQUFlLEVBQ2YsTUFBK0I7UUFFL0IsS0FBSyxDQUFDLE9BQU8sQ0FBQyxDQUFDO1FBQ2YsSUFBSSxDQUFDLElBQUksR0FBRyxnQkFBZ0IsQ0FBQztRQUM3QixJQUFJLENBQUMsTUFBTSxHQUFHLE1BQU0sQ0FBQztRQUNyQixJQUFJLENBQUMsSUFBSSxHQUFHLE1BQU0sQ0FBQyxJQUFJLENBQUM7UUFDeEIsSUFBSSxDQUFDLFNBQVMsR0FBRyxNQUFNLENBQUMsU0FBUyxDQUFDO0lBQ3BDLENBQUM7Q0FDRjtBQWZELG9EQWVDO0FBRUQ7O0dBRUc7QUFDSCxNQUFhLFNBQVUsU0FBUSxLQUFLO0lBQ2xDLFlBQVksT0FBZTtRQUN6QixLQUFLLENBQUMsT0FBTyxDQUFDLENBQUM7UUFDZixJQUFJLENBQUMsSUFBSSxHQUFHLFdBQVcsQ0FBQztJQUMxQixDQUFDO0NBQ0Y7QUFMRCw4QkFLQyJ9