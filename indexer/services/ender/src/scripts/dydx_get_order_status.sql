/**
  * The obvious case is if totalFilled >= size, then the order status should always be `FILLED`.
  * The difficult case is if totalFilled < size after a fill, then we need to keep the following
  * cases in mind:
  * 1. Stateful Orders - All cancelations are on-chain events, so the will be `OPEN`. The
  *     CanceledOrdersCache does not store any stateful orders and we never send
  *     BEST_EFFORT_CANCELED notifications for stateful orders.
  * 2. Short-term FOK - FOK orders can never be `OPEN`, since they don't rest on the orderbook, so
  *    totalFilled cannot be < size. By the end of the block, the order will be filled, so we mark
  *    it as `FILLED`.
  * 3. Short-term IOC - Protocol guarantees that an IOC order will only ever be filled in a single
  *    block, so status should be `CANCELED`.
  * 4. Short-term Limit & Post-only - If the order is in the CanceledOrdersCache, then it should be
  *    set to the corresponding CanceledOrderStatus, otherwise `OPEN`.
  * @param isCanceled - if the order is in the CanceledOrderCache, always false for liquidiation
  * orders
  */
CREATE OR REPLACE FUNCTION dydx_get_order_status(total_filled numeric, size numeric, order_canceled_status text, order_flags bigint, time_in_force text)
RETURNS text AS $$
DECLARE
    order_status text;
BEGIN
    IF total_filled >= size THEN
        RETURN 'FILLED';
    /** Order flag of 64 is a stateful term order */
    ELSIF order_flags = 64 THEN /** 1. Stateful Order */
        RETURN 'OPEN';
    ELSIF time_in_force = 'FOK' THEN /** 2. Short-term FOK */
        RETURN 'FILLED';
    ELSIF time_in_force = 'IOC' THEN /** 3. Short-term IOC */
        RETURN 'CANCELED';
    ELSIF order_canceled_status = 'BEST_EFFORT_CANCELED' THEN /** 4. Short-term Limit & Postonly */
        RETURN 'BEST_EFFORT_CANCELED';
    ELSIF order_canceled_status = 'CANCELED' THEN
        RETURN 'CANCELED';
    ELSE
        order_status = 'OPEN';
    END IF;

    RETURN order_status;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
