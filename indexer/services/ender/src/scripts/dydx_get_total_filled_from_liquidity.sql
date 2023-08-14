/**
  Returns the order total filled amount given the liquidity side.
*/
CREATE OR REPLACE FUNCTION get_total_filled(fill_liquidity text, event_data jsonb) RETURNS numeric AS $$
BEGIN
    IF fill_liquidity = 'TAKER' THEN
        RETURN dydx_from_jsonlib_long(event_data->'totalFilledTaker');
    ELSE
        RETURN dydx_from_jsonlib_long(event_data->'totalFilledMaker');
    END IF;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
