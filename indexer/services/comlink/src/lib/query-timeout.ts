import express from 'express';

// SQLSTATE for a statement cancelled by statement_timeout (or pg_cancel_backend).
const QUERY_CANCELED_SQLSTATE: string = '57014';

/**
 * Objection wraps driver errors in a DBError that keeps the original under `nativeError`, while
 * raw knex queries reject with the driver error itself, so check both shapes.
 */
export function isQueryCanceledError(error: unknown): boolean {
  if (typeof error !== 'object' || error === null) {
    return false;
  }
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const err: any = error;
  return err.code === QUERY_CANCELED_SQLSTATE ||
    err.nativeError?.code === QUERY_CANCELED_SQLSTATE;
}

/**
 * A cancelled read means the request outlived the database's time budget, not a server bug:
 * answer 504 like the load balancer would have, so it stays out of the 500 rate.
 */
export function createQueryTimeoutResponse(res: express.Response): express.Response {
  return res.status(504).json({
    errors: [{
      msg: 'Request timed out',
    }],
  });
}
