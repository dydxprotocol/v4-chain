CREATE OR REPLACE FUNCTION dydx_open_interest_update_handler(event_data jsonb) RETURNS jsonb AS $$
/**
  Parameters:
    - event_data: The 'data' field of the IndexerTendermintEvent (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/proto/dydxprotocol/indexer/indexer_manager/event.proto#L25)
        converted to JSON format. Conversion to JSON is expected to be done by JSON.stringify.
  Returns: JSON object containing fields:
    - open_interest_update: The updated perpetual market open interest

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    perpetual_market_record perpetual_markets%ROWTYPE;
    updates_array jsonb[];
    open_interest_update jsonb;
BEGIN
    FOR open_interest_update IN SELECT * FROM jsonb_array_elements(event_data->'openInterestUpdates') LOOP
      perpetual_market_record."id" = (open_interest_update->'perpetualId')::bigint;
      perpetual_market_record."openInterest" = dydx_from_serializable_int(open_interest_update->'openInterest');
      


          UPDATE perpetual_markets
          SET
              "openInterest" = perpetual_market_record."openInterest"
          WHERE
              "id" = perpetual_market_record."id"
          RETURNING * INTO perpetual_market_record;

      updates_array = array_append(updates_array, dydx_to_jsonb(perpetual_market_record));

    END LOOP;


    RETURN jsonb_build_object(
        'open_interest_updates',
        updates_array
    );
END;
$$ LANGUAGE plpgsql;