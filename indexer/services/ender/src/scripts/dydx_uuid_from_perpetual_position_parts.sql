/**
  Returns a UUID using the parts of a perpetual position.
*/
CREATE OR REPLACE FUNCTION dydx_uuid_from_perpetual_position_parts(subaccount_uuid uuid, open_event_id bytea) RETURNS uuid AS $$
BEGIN
    return dydx_uuid(concat(subaccount_uuid, '-', encode(open_event_id, 'hex')));
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
