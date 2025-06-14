CREATE OR REPLACE FUNCTION dydx_get_builder_address(fill_liquidity text, event_data jsonb) RETURNS text AS $$
/**
  Returns the builder address given the liquidity side.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    IF fill_liquidity = 'TAKER' THEN
        RETURN event_data->>'takerBuilderAddress';
    ELSE
        RETURN event_data->>'makerBuilderAddress';
    END IF;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;

