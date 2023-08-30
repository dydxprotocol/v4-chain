/**
  Returns a UUID using the parts of an IndexerSubaccountId (https://github.com/dydxprotocol/v4-proto/blob/437f6d8/dydxprotocol/indexer/protocol/v1/subaccount.proto#L15).
*/
CREATE OR REPLACE FUNCTION dydx_uuid_from_subaccount_id_parts(address text, subaccount_number text) RETURNS uuid AS $$
BEGIN
    RETURN dydx_uuid(concat(address, '-', subaccount_number));
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
