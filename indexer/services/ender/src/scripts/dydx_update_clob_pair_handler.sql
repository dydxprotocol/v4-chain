CREATE OR REPLACE FUNCTION dydx_update_clob_pair_handler(event_data jsonb) RETURNS jsonb AS $$
/**
  Parameters:
    - event_data: The 'data' field of the IndexerTendermintEvent (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/proto/dydxprotocol/indexer/indexer_manager/event.proto#L25)
        converted to JSON format. Conversion to JSON is expected to be done by JSON.stringify.
  Returns: JSON object containing fields:
    - perpetual_market: The updated perpetual market in perpetual-market-model format (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/indexer/packages/postgres/src/models/perpetual-market-model.ts).

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    clob_pair_id bigint;
    perpetual_market_record perpetual_markets%ROWTYPE;
BEGIN
    clob_pair_id = (event_data->'clobPairId')::bigint;
    perpetual_market_record."status" = dydx_clob_pair_status_to_market_status(event_data->'status');
    perpetual_market_record."quantumConversionExponent" = (event_data->'quantumConversionExponent')::integer;
    perpetual_market_record."subticksPerTick" = (event_data->'subticksPerTick')::integer;
    perpetual_market_record."stepBaseQuantums" = dydx_from_jsonlib_long(event_data->'stepBaseQuantums');

    UPDATE perpetual_markets
    SET
        "status" = perpetual_market_record."status",
        "quantumConversionExponent" = perpetual_market_record."quantumConversionExponent",
        "subticksPerTick" = perpetual_market_record."subticksPerTick",
        "stepBaseQuantums" = perpetual_market_record."stepBaseQuantums"
    WHERE "clobPairId" = clob_pair_id
    RETURNING * INTO perpetual_market_record;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Could not find perpetual market with corresponding clobPairId %', event_data;
    END IF;

    RETURN jsonb_build_object(
            'perpetual_market',
            dydx_to_jsonb(perpetual_market_record)
        );
END;
$$ LANGUAGE plpgsql;