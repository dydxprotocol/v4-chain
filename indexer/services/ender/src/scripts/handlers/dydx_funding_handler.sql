CREATE OR REPLACE FUNCTION dydx_funding_handler(
    block_height int, block_time timestamp, event_data jsonb, event_index int, transaction_index int) RETURNS jsonb AS $$
DECLARE
    PPM_EXPONENT constant numeric = -6;
    FUNDING_RATE_FROM_PROTOCOL_IN_HOURS constant numeric = 8;
    QUOTE_CURRENCY_ATOMIC_RESOLUTION constant numeric = -6;

    TYPE_PREMIUM_SAMPLE constant jsonb = '1';
    TYPE_FUNDING_RATE_AND_INDEX constant jsonb = '2';

    perpetual_market_id bigint;
    funding_index_updates_record funding_index_updates%ROWTYPE;

    funding_update jsonb;
    perpetual_markets_response jsonb = jsonb_build_object();
    funding_update_response jsonb = jsonb_build_object();
    errors_response jsonb[] = '{}';
    event_id bytea;

    -- Declare JSONB objects for the maps
    perpetual_market_map jsonb;
    oracle_price_map jsonb;
BEGIN
    -- Build perpetual_market_map using jsonb_object_agg
    SELECT jsonb_object_agg(id::text, dydx_to_jsonb(perpetual_markets::perpetual_markets))
    INTO perpetual_market_map
    FROM perpetual_markets;

    -- Build oracle_price_map using jsonb_object_agg with latest prices
    SELECT jsonb_object_agg("marketId"::text, dydx_to_jsonb(op::oracle_prices))
    INTO oracle_price_map
    FROM (
        SELECT DISTINCT ON ("marketId") *
        FROM oracle_prices
        WHERE "effectiveAtHeight" <= block_height
        ORDER BY "marketId", "effectiveAtHeight" DESC
    ) op;

    -- Process each funding update
    FOR funding_update IN SELECT * FROM jsonb_array_elements(event_data->'updates') LOOP
        perpetual_market_id = (funding_update->'perpetualId')::bigint;

        -- Retrieve perpetual market from map
        PERFORM jsonb_populate_record(null::perpetual_markets, perpetual_market_map->(perpetual_market_id::text));

        IF perpetual_market_map->(perpetual_market_id::text) IS NULL THEN
            errors_response = array_append(errors_response, '"Received FundingUpdate with unknown perpetualId."'::jsonb);
            CONTINUE;
        END IF;

        perpetual_markets_response = jsonb_set(perpetual_markets_response, ARRAY[(perpetual_market_id::text)], perpetual_market_map->(perpetual_market_id::text));

        CASE event_data->'type'
            WHEN TYPE_PREMIUM_SAMPLE THEN
                -- Do nothing for premium samples
            WHEN TYPE_FUNDING_RATE_AND_INDEX THEN
                -- Retrieve the latest oracle price for the marketId
                IF oracle_price_map->(perpetual_market_map->(perpetual_market_id::text)->>'marketId') IS NULL THEN
                    errors_response = array_append(errors_response, '"oracle_price not found for marketId."'::jsonb);
                    CONTINUE;
                END IF;

                event_id = dydx_event_id_from_parts(block_height, transaction_index, event_index);

                funding_index_updates_record."id" = dydx_uuid_from_funding_index_update_parts(
                    block_height,
                    event_id,
                    perpetual_market_id);
                funding_index_updates_record."perpetualId" = perpetual_market_id;
                funding_index_updates_record."eventId" = event_id;
                funding_index_updates_record."effectiveAt" = block_time;
                funding_index_updates_record."rate" = dydx_trim_scale(
                    power(10, PPM_EXPONENT) /
                    FUNDING_RATE_FROM_PROTOCOL_IN_HOURS *
                    (funding_update->'fundingValuePpm')::numeric);
                funding_index_updates_record."oraclePrice" = (oracle_price_map->(perpetual_market_map->(perpetual_market_id::text)->>'marketId'))->>'price';
                funding_index_updates_record."fundingIndex" = dydx_trim_scale(
                    dydx_from_serializable_int(funding_update->'fundingIndex') *
                    power(10, PPM_EXPONENT + QUOTE_CURRENCY_ATOMIC_RESOLUTION - (perpetual_market_map->(perpetual_market_id::text)->>'atomicResolution')::numeric));
                funding_index_updates_record."effectiveAtHeight" = block_height;

                INSERT INTO funding_index_updates VALUES (funding_index_updates_record.*);
                funding_update_response = jsonb_set(funding_update_response, ARRAY[(funding_index_updates_record."perpetualId")::text], dydx_to_jsonb(funding_index_updates_record));

            ELSE
                errors_response = array_append(errors_response, '"Received unknown FundingEvent type."'::jsonb);
                CONTINUE;
        END CASE;

        errors_response = array_append(errors_response, NULL);
    END LOOP;

    RETURN jsonb_build_object(
        'perpetual_markets',
        perpetual_markets_response,
        'funding_index_updates',
        funding_update_response,
        'errors',
        to_jsonb(errors_response)
    );
END;
$$ LANGUAGE plpgsql;
