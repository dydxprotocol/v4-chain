/**
  Returns a UUID using the JSON.stringify format of an IndexerSubaccountId (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/proto/dydxprotocol/indexer/protocol/v1/subaccount.proto#L15).
*/
CREATE OR REPLACE FUNCTION dydx_uuid_from_subaccount_id(subaccount_id jsonb) RETURNS uuid AS $$
BEGIN
    RETURN dydx_uuid_from_subaccount_id_parts(subaccount_id->>'owner', subaccount_id->>'number');
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
