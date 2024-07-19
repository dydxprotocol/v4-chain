CREATE OR REPLACE FUNCTION dydx_uuid_from_fill_event_parts(event_id bytea, liquidity text) RETURNS uuid AS $$
/**
  Returns a UUID using the parts of a fill event.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    return dydx_uuid(concat(encode(event_id, 'hex'), '-', liquidity));
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
