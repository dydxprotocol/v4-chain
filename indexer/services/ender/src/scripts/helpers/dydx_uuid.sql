CREATE OR REPLACE FUNCTION dydx_uuid(name text) RETURNS uuid AS $$
/**
  Computes a UUID using a well known namespace.

  The namespace must match the well known constant defined in
  https://github.com/dydxprotocol/indexer/blob/6aafb97/packages/postgres/src/helpers/uuid.ts#L4.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    RETURN uuid_generate_v5('0f9da948-a6fb-4c45-9edc-4685c3f3317d', name);
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
