CREATE OR REPLACE FUNCTION dydx_uuid_from_yield_params_parts(block_height int) RETURNS uuid AS $$
/**
  Returns a UUID using the parts of the yield params.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    return dydx_uuid(block_height::text);
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
