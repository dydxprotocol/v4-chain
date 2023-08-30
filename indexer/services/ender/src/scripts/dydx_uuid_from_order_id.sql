/**
  Returns a UUID using the JSON.stringify format of an IndexerOrderId (https://github.com/dydxprotocol/v4-proto/blob/437f6d8/dydxprotocol/indexer/protocol/v1/clob.proto#L15).
*/
CREATE OR REPLACE FUNCTION dydx_uuid_from_order_id(order_id jsonb) RETURNS uuid AS $$
BEGIN
    return dydx_uuid_from_order_id_parts(
        dydx_uuid_from_subaccount_id(order_id->'subaccountId'),
        order_id->>'clientId',
        order_id->>'clobPairId',
        order_id->>'orderFlags');
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
