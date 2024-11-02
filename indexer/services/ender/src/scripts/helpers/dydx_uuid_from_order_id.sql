CREATE OR REPLACE FUNCTION klyra_uuid_from_order_id(order_id jsonb) RETURNS uuid AS $$
/**
  Returns a UUID using the JSON.stringify format of an IndexerOrderId.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    return klyra_uuid_from_order_id_parts(
        klyra_uuid_from_subaccount_id(order_id->'subaccountId'),
        order_id->>'clientId',
        order_id->>'clobPairId',
        order_id->>'orderFlags');
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
