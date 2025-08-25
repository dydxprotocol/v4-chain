CREATE OR REPLACE FUNCTION dydx_get_order_router_fee(fill_liquidity text, event_data jsonb) RETURNS numeric AS $$
/**
  Returns the order router fee given the liquidity side.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    IF fill_liquidity = 'TAKER' THEN
        RETURN dydx_from_jsonlib_long(event_data->'takerOrderRouterFee');
    ELSE
        RETURN dydx_from_jsonlib_long(event_data->'makerOrderRouterFee');
    END IF;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;

