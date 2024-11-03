CREATE OR REPLACE FUNCTION klyra_uuid_from_subaccount_id_parts(address text, subaccount_number text) RETURNS uuid AS $$
/**
  Returns a UUID using the parts of an IndexerSubaccountId.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    RETURN klyra_uuid(concat(address, '-', subaccount_number));
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
