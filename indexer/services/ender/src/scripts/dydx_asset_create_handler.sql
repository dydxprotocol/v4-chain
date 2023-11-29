CREATE OR REPLACE FUNCTION dydx_asset_create_handler(event_data jsonb) RETURNS jsonb AS $$
/**
  Parameters:
    - event_data: The 'data' field of the IndexerTendermintEvent (https://github.com/dydxprotocol/v4-proto/blob/8d35c86/dydxprotocol/indexer/indexer_manager/event.proto#L25)
        converted to JSON format. Conversion to JSON is expected to be done by JSON.stringify.
  Returns: JSON object containing fields:
    - asset: The created asset in asset-model format (https://github.com/dydxprotocol/indexer/blob/cc70982/packages/postgres/src/models/asset-model.ts).

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    market_record_id integer;
    asset_record assets%ROWTYPE;
BEGIN
    asset_record."id" = event_data->>'id';
    asset_record."atomicResolution" = (event_data->'atomicResolution')::integer;
    asset_record."symbol" = event_data->>'symbol';

    asset_record."hasMarket" = (event_data->'hasMarket')::bool;
    if asset_record."hasMarket" THEN
        market_record_id = (event_data->'marketId')::integer;
        SELECT "id" INTO asset_record."marketId" FROM markets WHERE "id" = market_record_id;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'Unable to find market with id: %', market_record_id;
        END IF;
    END IF;

    INSERT INTO assets VALUES (asset_record.*);

    RETURN jsonb_build_object(
            'asset',
            dydx_to_jsonb(asset_record)
        );
END;
$$ LANGUAGE plpgsql;