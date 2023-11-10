/**
  Parameters:
    - block_height: the height of the block being processing.
    - block_time: the time of the block being processed.
    - event_data: The 'data' field of the IndexerTendermintEvent (https://github.com/dydxprotocol/v4-proto/blob/8d35c86/dydxprotocol/indexer/indexer_manager/event.proto#L25)
        converted to JSON format. Conversion to JSON is expected to be done by JSON.stringify.
  Returns: JSON object containing fields:
    - market: The updated market in market-model format (https://github.com/dydxprotocol/indexer/blob/cc70982/packages/postgres/src/models/market-model.ts).
    - oracle_price: The created oracle price in oracle-price-model format (https://github.com/dydxprotocol/indexer/blob/cc70982/packages/postgres/src/models/oracle-price-model.ts).
*/
CREATE OR REPLACE FUNCTION dydx_market_price_update_handler(block_height int, block_time timestamp, event_data jsonb) RETURNS jsonb AS $$
DECLARE
    market_record_id integer;
    market_record markets%ROWTYPE;
    oracle_price numeric;
    oracle_price_record oracle_prices%ROWTYPE;
BEGIN
    market_record_id = (event_data->'marketId')::integer;
    SELECT * INTO market_record FROM markets WHERE "id" = market_record_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'MarketPriceUpdateEvent contains a non-existent market id. Id: %', market_record_id;
    END IF;

    oracle_price = dydx_trim_scale(
        dydx_from_jsonlib_long(event_data->'priceUpdate'->'priceWithExponent') *
        power(10, market_record.exponent::numeric));

    market_record."oraclePrice" = oracle_price;

    UPDATE markets
    SET
        "oraclePrice" = market_record."oraclePrice"
    WHERE id = market_record."id";

    oracle_price_record."id" = dydx_uuid_from_oracle_price_parts(market_record_id, block_height);
    oracle_price_record."effectiveAt" = block_time;
    oracle_price_record."effectiveAtHeight" = block_height;
    oracle_price_record."marketId" = market_record_id;
    oracle_price_record."price" = oracle_price;

    INSERT INTO oracle_prices VALUES (oracle_price_record.*);

    RETURN jsonb_build_object(
        'market',
        dydx_to_jsonb(market_record),
        'oracle_price',
        dydx_to_jsonb(oracle_price_record)
    );
END;
$$ LANGUAGE plpgsql;