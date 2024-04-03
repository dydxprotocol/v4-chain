CREATE OR REPLACE FUNCTION dydx_funding_handler(
    block_height int, block_time timestamp, event_data jsonb, event_index int, transaction_index int) RETURNS jsonb AS $$
/**
  Parameters:
    - block_height: the height of the block being processing.
    - block_time: the time of the block being processed.
    - event_data: The 'data' field of the IndexerTendermintEvent (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/proto/dydxprotocol/indexer/indexer_manager/event.proto#L25)
        converted to JSON format. Conversion to JSON is expected to be done by JSON.stringify.
    - event_index: The 'event_index' of the IndexerTendermintEvent.
    - transaction_index: The transaction_index of the IndexerTendermintEvent after the conversion that takes into
        account the block_event (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/indexer/services/ender/src/lib/helper.ts#L41)
  Returns: JSON object containing fields:
    - perpetual_markets: A mapping from perpetual market id to the associated perpetual market in perpetual-market-model format (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/indexer/packages/postgres/src/models/perpetual-market-model.ts).
    - errors: An array containing an error string (or NULL if no error occurred) for each FundingEventUpdate.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    PPM_EXPONENT constant numeric = -6;
    FUNDING_RATE_FROM_PROTOCOL_IN_HOURS constant numeric = 8;
    QUOTE_CURRENCY_ATOMIC_RESOLUTION constant numeric = -6;

    TYPE_PREMIUM_SAMPLE constant jsonb = '1';
    TYPE_FUNDING_RATE_AND_INDEX constant jsonb = '2';

    perpetual_market_id bigint;
    perpetual_market_record perpetual_markets%ROWTYPE;
    funding_index_updates_record funding_index_updates%ROWTYPE;
    oracle_prices_record oracle_prices%ROWTYPE;

    funding_update jsonb;
    perpetual_markets_response jsonb = jsonb_build_object();
    funding_update_response jsonb = jsonb_build_object();
    errors_response jsonb[];
    event_id bytea;
BEGIN
    FOR funding_update IN SELECT * FROM jsonb_array_elements(event_data->'updates') LOOP
        perpetual_market_id = (funding_update->'perpetualId')::bigint;
        SELECT * INTO perpetual_market_record FROM perpetual_markets WHERE "id" = perpetual_market_id;
        IF NOT FOUND THEN
            errors_response = array_append(errors_response, '"Received FundingUpdate with unknown perpetualId."'::jsonb);
            CONTINUE;
        END IF;

        perpetual_markets_response = jsonb_set(perpetual_markets_response, ARRAY[(perpetual_market_record."id")::text], dydx_to_jsonb(perpetual_market_record));

        CASE event_data->'type'
            WHEN TYPE_PREMIUM_SAMPLE THEN
                /** Here we just need to return the associated perpetual market. */
            WHEN TYPE_FUNDING_RATE_AND_INDEX THEN
                /** Returns the latest oracle price <= current block_height. */
                SELECT * INTO oracle_prices_record
                         FROM oracle_prices
                         WHERE "marketId" = perpetual_market_record."marketId" AND "effectiveAtHeight" <= block_height
                         ORDER BY "effectiveAtHeight"
                         DESC LIMIT 1;
                IF NOT FOUND THEN
                    errors_response = array_append(errors_response, '"oracle_price not found for marketId."'::jsonb);
                    CONTINUE;
                END IF;

                event_id = dydx_event_id_from_parts(block_height, transaction_index, event_index);

                funding_index_updates_record."id" = dydx_uuid_from_funding_index_update_parts(
                    block_height,
                    event_id,
                    perpetual_market_record."id");
                funding_index_updates_record."perpetualId" = perpetual_market_id;
                funding_index_updates_record."eventId" = event_id;
                funding_index_updates_record."effectiveAt" = block_time;
                funding_index_updates_record."rate" = dydx_trim_scale(
                    power(10, PPM_EXPONENT) /
                    FUNDING_RATE_FROM_PROTOCOL_IN_HOURS *
                    (funding_update->'fundingValuePpm')::numeric);
                funding_index_updates_record."oraclePrice" = oracle_prices_record."price";
                funding_index_updates_record."fundingIndex" = dydx_trim_scale(
                    dydx_from_serializable_int(funding_update->'fundingIndex') *
                    power(10,
                        PPM_EXPONENT + QUOTE_CURRENCY_ATOMIC_RESOLUTION - perpetual_market_record."atomicResolution"));
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
