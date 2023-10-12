/**
  Returns the order status given the total filled amount, the order size, whether the order was
  cancelled, order flags, and time in force.
  * The obvious case is if totalFilled >= size, then the order status should always be `FILLED`.
  * The difficult case is if totalFilled < size after a fill, then we need to keep the following
  * cases in mind:
  * 1. Stateful Orders - All cancelations are on-chain events, so the order can be `OPEN` or 
  *    `BEST_EFFORT_CANCELED` if the order is in the CanceledOrdersCache.
  * 2. Short-term FOK - FOK orders can never be `OPEN`, since they don't rest on the orderbook, so
  *    totalFilled cannot be < size. By the end of the block, the order will be filled, so we mark
  *    it as `FILLED`.
  * 3. Short-term IOC - Protocol guarantees that an IOC order will only ever be filled in a single
  *    block, so status should be `CANCELED`.
  * 4. Short-term Limit & Post-only - If the order is in the CanceledOrdersCache, then it should be
  *    set to `BEST_EFFORT_CANCELED`, otherwise `OPEN`.
*/
CREATE OR REPLACE FUNCTION get_order_status(total_filled numeric, size numeric, is_cancelled boolean, order_flags bigint, time_in_force text)
RETURNS text AS $$
BEGIN
    IF total_filled >= size THEN
        RETURN 'FILLED';
    /** Order flag of 64 is a stateful term order */
    ELSIF order_flags = 64 THEN /** 1. Stateful Order */
        IF is_cancelled THEN
            RETURN 'BEST_EFFORT_CANCELED';
        ELSE
            RETURN 'OPEN';
        END IF;
    ELSIF time_in_force = 'FOK' THEN /** 2. Short-term FOK */
        RETURN 'FILLED';
    ELSIF time_in_force = 'IOC' THEN /** 3. Short-term IOC */
        RETURN 'CANCELED';
    ELSIF is_cancelled THEN /** 4. Short-term Limit & Postonly */
        RETURN 'BEST_EFFORT_CANCELED';
    ELSE
        RETURN 'OPEN';
    END IF;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
