CREATE OR REPLACE FUNCTION dydx_get_fee(fill_liquidity text, event_data jsonb) RETURNS numeric AS $$
/**
  Returns the fee given the liquidity side.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    IF fill_liquidity = 'TAKER' THEN
        RETURN dydx_from_jsonlib_long(event_data->'takerFee');
    ELSE
        RETURN dydx_from_jsonlib_long(event_data->'makerFee');
    END IF;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;

