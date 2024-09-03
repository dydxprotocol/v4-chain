CREATE OR REPLACE FUNCTION dydx_register_affiliate_handler(block_height int, event_data jsonb) RETURNS jsonb AS $$
/**
  Parameters:
    - event_data: The 'data' field of the IndexerTendermintEvent (https://github.com/dydxprotocol/v4-proto/blob/8d35c86/dydxprotocol/indexer/indexer_manager/event.proto#L25)
        converted to JSON format. Conversion to JSON is expected to be done by JSON.stringify.
  Returns: JSON object containing fields:
    - asset: The created asset in asset-model format (https://github.com/dydxprotocol/indexer/blob/cc70982/packages/postgres/src/models/asset-model.ts).

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    referral_record affiliate_referred_users%ROWTYPE;
BEGIN
    referral_record."affiliateAddress" = event_data->>'affiliate';
    referral_record."refereeAddress" = event_data->>'referee';
    referral_record."referredAtBlock" = block_height;

    INSERT INTO affiliate_referred_users VALUES (referral_record.*);

    RETURN jsonb_build_object(
            'referral',
            dydx_to_jsonb(referral_record)
        );
END;
$$ LANGUAGE plpgsql;
