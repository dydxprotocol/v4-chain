import { Rpc } from "../../../helpers";
import { MsgAuthorizeCircuitBreaker, MsgAuthorizeCircuitBreakerResponse, MsgTripCircuitBreaker, MsgTripCircuitBreakerResponse, MsgResetCircuitBreaker, MsgResetCircuitBreakerResponse } from "./tx";
/** Msg defines the circuit Msg service. */
export interface Msg {
    /**
     * AuthorizeCircuitBreaker allows a super-admin to grant (or revoke) another
     * account's circuit breaker permissions.
     */
    authorizeCircuitBreaker(request: MsgAuthorizeCircuitBreaker): Promise<MsgAuthorizeCircuitBreakerResponse>;
    /** TripCircuitBreaker pauses processing of Msg's in the state machine. */
    tripCircuitBreaker(request: MsgTripCircuitBreaker): Promise<MsgTripCircuitBreakerResponse>;
    /**
     * ResetCircuitBreaker resumes processing of Msg's in the state machine that
     * have been been paused using TripCircuitBreaker.
     */
    resetCircuitBreaker(request: MsgResetCircuitBreaker): Promise<MsgResetCircuitBreakerResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    authorizeCircuitBreaker(request: MsgAuthorizeCircuitBreaker): Promise<MsgAuthorizeCircuitBreakerResponse>;
    tripCircuitBreaker(request: MsgTripCircuitBreaker): Promise<MsgTripCircuitBreakerResponse>;
    resetCircuitBreaker(request: MsgResetCircuitBreaker): Promise<MsgResetCircuitBreakerResponse>;
}
