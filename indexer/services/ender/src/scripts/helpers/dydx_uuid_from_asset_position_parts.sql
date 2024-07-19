CREATE OR REPLACE FUNCTION dydx_uuid_from_asset_position_parts(subaccount_uuid uuid, asset_id text) RETURNS uuid AS $$
/**
  Returns a UUID using the parts of an asset position.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    return dydx_uuid(concat(subaccount_uuid, '-', asset_id));
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
