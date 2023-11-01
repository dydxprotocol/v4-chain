/**
  Returns the order status given the total filled amount, the order size and whether the order was cancelled.
*/
CREATE OR REPLACE FUNCTION dydx_get_order_status(total_filled numeric, size numeric, is_cancelled boolean, order_flags bigint, time_in_force text)
RETURNS text AS $$
DECLARE
    order_status text;
BEGIN
    IF is_cancelled = true THEN
        order_status = 'BEST_EFFORT_CANCELED';
    ELSIF total_filled >= size THEN
        order_status = 'FILLED';
    ELSE
        order_status = 'OPEN';
    END IF;

    RETURN order_status;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
