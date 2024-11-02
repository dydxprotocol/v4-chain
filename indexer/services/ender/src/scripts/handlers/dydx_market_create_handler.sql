CREATE OR REPLACE FUNCTION klyra_market_create_handler(event_data jsonb) RETURNS jsonb AS $$
/**
  Parameters:
    - event_data: The 'data' field of the IndexerTendermintEvent
        converted to JSON format. Conversion to JSON is expected to be done by JSON.stringify.
  Returns: JSON object containing fields:
    - market: The created market in market-model format.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    market_record_id integer;
    market_record markets%ROWTYPE;
BEGIN
    market_record_id = (event_data->'marketId')::integer;
    SELECT * INTO market_record FROM markets WHERE "id" = market_record_id;

    IF FOUND THEN
        RAISE EXCEPTION 'Market in MarketCreate already exists. Record: %', market_record;
    END IF;

    market_record."id" = market_record_id;
    market_record."pair" = event_data->'marketCreate'->'base'->>'pair';
    market_record."exponent" = (event_data->'marketCreate'->'exponent')::integer;
    market_record."minPriceChangePpm" = (event_data->'marketCreate'->'base'->'minPriceChangePpm')::integer;

    INSERT INTO markets VALUES (market_record.*);

    RETURN jsonb_build_object(
        'market',
        klyra_to_jsonb(market_record)
    );
END;
$$ LANGUAGE plpgsql;