CREATE OR REPLACE FUNCTION dydx_uuid_from_order_id(order_id jsonb) RETURNS uuid AS $$
/**
  Returns a UUID using the JSON.stringify format of an IndexerOrderId (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/proto/dydxprotocol/indexer/protocol/v1/clob.proto#L15).

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    IF (order_id->>'orderFlags')::bigint = constants.order_flag_twap_suborder() THEN
        -- Twap suborders should be mapped to their parent order.
        return dydx_uuid_from_order_id_parts(
            dydx_uuid_from_subaccount_id(order_id->'subaccountId'),
            order_id->>'clientId',
            order_id->>'clobPairId',
            constants.order_flag_twap()::text);
    END IF;

    return dydx_uuid_from_order_id_parts(
        dydx_uuid_from_subaccount_id(order_id->'subaccountId'),
        order_id->>'clientId',
        order_id->>'clobPairId',
        order_id->>'orderFlags');
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
