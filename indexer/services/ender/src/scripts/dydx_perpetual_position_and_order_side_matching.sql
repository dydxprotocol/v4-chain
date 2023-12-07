CREATE OR REPLACE FUNCTION dydx_perpetual_position_and_order_side_matching(
    perpetual_position_side text, order_side text) RETURNS boolean AS $$
/**
  Returns true iff perpetual_position_side is LONG and order_side is BUY or if perpetual_position_side is SHORT and
  order_side is SELL.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    RETURN (perpetual_position_side = 'LONG' AND order_side = 'BUY') OR
           (perpetual_position_side = 'SHORT' AND order_side = 'SELL');
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
