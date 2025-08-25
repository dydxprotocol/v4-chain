CREATE OR REPLACE FUNCTION dydx_perpetual_market_v2_handler(event_data jsonb) RETURNS jsonb AS $$
/**
  Parameters:
    - event_data: The 'data' field of the IndexerTendermintEvent (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/proto/dydxprotocol/indexer/indexer_manager/event.proto#L25)
        converted to JSON format. Conversion to JSON is expected to be done by JSON.stringify.
  Returns: JSON object containing fields:
    - perpetual_market: The updated perpetual market in perpetual-market-model format (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/indexer/packages/postgres/src/models/perpetual-market-model.ts).

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    perpetual_market_record perpetual_markets%ROWTYPE;
BEGIN
    perpetual_market_record."id" = (event_data->'id')::bigint;
    perpetual_market_record."clobPairId" = (event_data->'clobPairId')::bigint;
    perpetual_market_record."ticker" = event_data->>'ticker';
    perpetual_market_record."marketId" = (event_data->'marketId')::integer;
    perpetual_market_record."status" = dydx_clob_pair_status_to_market_status(event_data->'status');
    perpetual_market_record."priceChange24H" = 0;
    perpetual_market_record."trades24H" = 0;
    perpetual_market_record."volume24H" = 0;
    perpetual_market_record."nextFundingRate" = 0;
    perpetual_market_record."openInterest"= 0;
    perpetual_market_record."quantumConversionExponent" = (event_data->'quantumConversionExponent')::integer;
    perpetual_market_record."atomicResolution" = (event_data->'atomicResolution')::integer;
    perpetual_market_record."subticksPerTick" = (event_data->'subticksPerTick')::integer;
    perpetual_market_record."stepBaseQuantums" = dydx_from_jsonlib_long(event_data->'stepBaseQuantums');
    perpetual_market_record."liquidityTierId" = (event_data->'liquidityTier')::integer;
    perpetual_market_record."marketType" = dydx_protocol_market_type_to_perpetual_market_type(event_data->'marketType');
    perpetual_market_record."baseOpenInterest" = 0;

    INSERT INTO perpetual_markets (
        "id", "clobPairId", "ticker", "marketId", "status", "priceChange24H",
        "trades24H", "volume24H", "nextFundingRate", "openInterest",
        "quantumConversionExponent", "atomicResolution", "subticksPerTick",
        "stepBaseQuantums", "liquidityTierId", "marketType", "baseOpenInterest"
    ) VALUES (
        perpetual_market_record."id", perpetual_market_record."clobPairId", perpetual_market_record."ticker",
        perpetual_market_record."marketId", perpetual_market_record."status", perpetual_market_record."priceChange24H",
        perpetual_market_record."trades24H", perpetual_market_record."volume24H", perpetual_market_record."nextFundingRate",
        perpetual_market_record."openInterest", perpetual_market_record."quantumConversionExponent",
        perpetual_market_record."atomicResolution", perpetual_market_record."subticksPerTick",
        perpetual_market_record."stepBaseQuantums", perpetual_market_record."liquidityTierId",
        perpetual_market_record."marketType", perpetual_market_record."baseOpenInterest"
    ) RETURNING "id", "clobPairId", "ticker", "marketId", "status", "priceChange24H",
        "trades24H", "volume24H", "nextFundingRate", "openInterest",
        "quantumConversionExponent", "atomicResolution", "subticksPerTick",
        "stepBaseQuantums", "liquidityTierId", "marketType", "baseOpenInterest"
    INTO perpetual_market_record;

    RETURN jsonb_build_object(
            'perpetual_market',
            dydx_to_jsonb(perpetual_market_record)
        );
END;
$$ LANGUAGE plpgsql;
