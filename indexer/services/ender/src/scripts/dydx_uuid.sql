/**
  Computes a UUID using a well known namespace.

  The namespace must match the well known constant defined in
  https://github.com/dydxprotocol/indexer/blob/6aafb97/packages/postgres/src/helpers/uuid.ts#L4.
*/
CREATE OR REPLACE FUNCTION dydx_uuid(name text) RETURNS uuid AS $$
BEGIN
    RETURN uuid_generate_v5('0f9da948-a6fb-4c45-9edc-4685c3f3317d', name);
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
