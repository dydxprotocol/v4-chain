import { TIME_IN_FORCE_TO_API_TIME_IN_FORCE } from '../constants';
import { APITimeInForce, TimeInForce } from '../types';

/**
 * Gets post-only boolean from Indexer TimeInForce value.
 * Only returns true if the passed in TimeInForce value is POST_ONLY
 * @param timeInForce
 * @returns
 */
export function isOrderTIFPostOnly(timeInForce: TimeInForce): boolean {
  return timeInForce === TimeInForce.POST_ONLY;
}

/**
 * Converts Indexer TimeInForce to APITimeInForce value.
 * Special cases POST_ONLY as GTT.
 * @param timeInForce
 * @returns
 */
export function orderTIFToAPITIF(timeInForce: TimeInForce): APITimeInForce {
  return TIME_IN_FORCE_TO_API_TIME_IN_FORCE[timeInForce];
}
