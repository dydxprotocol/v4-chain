/**
  Returns a UUID using the parts of an IndexerOrderId (https://github.com/dydxprotocol/v4-proto/blob/437f6d8/dydxprotocol/indexer/protocol/v1/clob.proto#L15).
*/
CREATE OR REPLACE FUNCTION dydx_uuid_from_order_id_parts(subaccount_id uuid, client_id text, clob_pair_id text, order_flags text) RETURNS uuid AS $$
BEGIN
    return dydx_uuid(concat(subaccount_id, '-', client_id, '-', clob_pair_id, '-', order_flags));
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
