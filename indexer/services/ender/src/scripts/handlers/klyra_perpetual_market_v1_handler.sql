CREATE OR REPLACE FUNCTION klyra_perpetual_market_v1_handler(event_data jsonb) RETURNS jsonb AS $$
/**
  Note: This is a deprecated handler, see `klyra_perpetual_market_v2_handler` for the latest handler.
  Parameters:
    - event_data: The 'data' field of the IndexerTendermintEvent
        converted to JSON format. Conversion to JSON is expected to be done by JSON.stringify.
  Returns: JSON object containing fields:
    - perpetual_market: The updated perpetual market in perpetual-market-model format.
  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    perpetual_market_record perpetual_markets%ROWTYPE;
BEGIN
    perpetual_market_record."id" = (event_data->'id')::bigint;
    perpetual_market_record."clobPairId" = (event_data->'clobPairId')::bigint;
    perpetual_market_record."ticker" = event_data->>'ticker';
    perpetual_market_record."marketId" = (event_data->'marketId')::integer;
    perpetual_market_record."status" = klyra_clob_pair_status_to_market_status(event_data->'status');
    perpetual_market_record."priceChange24H" = 0;
    perpetual_market_record."trades24H" = 0;
    perpetual_market_record."volume24H" = 0;
    perpetual_market_record."nextFundingRate" = 0;
    perpetual_market_record."openInterest"= 0;
    perpetual_market_record."quantumConversionExponent" = (event_data->'quantumConversionExponent')::integer;
    perpetual_market_record."atomicResolution" = (event_data->'atomicResolution')::integer;
    perpetual_market_record."dangerIndexPpm" = 1000000;
    perpetual_market_record."isolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock" = 0;
    perpetual_market_record."subticksPerTick" = (event_data->'subticksPerTick')::integer;
    perpetual_market_record."stepBaseQuantums" = klyra_from_jsonlib_long(event_data->'stepBaseQuantums');
    perpetual_market_record."liquidityTierId" = (event_data->'liquidityTier')::integer;
    perpetual_market_record."marketType" = 'CROSS';
    perpetual_market_record."baseOpenInterest" = 0;
    perpetual_market_record."perpYieldIndex" = '0/1';


    INSERT INTO perpetual_markets VALUES (perpetual_market_record.*) RETURNING * INTO perpetual_market_record;

    RETURN jsonb_build_object(
            'perpetual_market',
            klyra_to_jsonb(perpetual_market_record)
        );
END;
$$ LANGUAGE plpgsql;