CREATE OR REPLACE FUNCTION klyra_uuid_from_subaccount_id(subaccount_id jsonb) RETURNS uuid AS $$
/**
  Returns a UUID using the JSON.stringify format of an IndexerSubaccountId.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    RETURN klyra_uuid_from_subaccount_id_parts(subaccount_id->>'owner', subaccount_id->>'number');
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
