CREATE OR REPLACE FUNCTION dydx_get_order_router_address(fill_liquidity text, event_data jsonb) RETURNS text AS $$
/**
  Returns the order router address given the liquidity side.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    IF fill_liquidity = 'TAKER' THEN
        RETURN event_data->>'takerOrderRouterAddress';
    ELSE
        RETURN event_data->>'makerOrderRouterAddress';
    END IF;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;

