CREATE OR REPLACE FUNCTION dydx_trading_rewards_handler(
    block_height int, block_time timestamp, event_data jsonb, event_index int, transaction_index int,
    transaction_hash text) RETURNS jsonb AS $$
/**
  Parameters:
    - block_height: the height of the block being processing.
    - block_time: the time of the block being processed.
    - event_data: The 'data' field of the IndexerTendermintEvent (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/proto/dydxprotocol/indexer/indexer_manager/event.proto#L25)
        converted to JSON format. Conversion to JSON is expected to be done by JSON.stringify.
    - event_index: The 'event_index' of the IndexerTendermintEvent.
    - transaction_index: The transaction_index of the IndexerTendermintEvent after the conversion that takes into
        account the block_event (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/indexer/services/ender/src/lib/helper.ts#L41)
    - transaction_hash: The transaction hash corresponding to this event from the IndexerTendermintBlock 'tx_hashes'.
  Returns: JSON object containing fields:
    - tradingRewards: A list of the trading rewards in the trading-reward-model format (indexer/packages/postgres/src/models/trading-reward-model.ts)

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    trading_rewards_array jsonb[];
    trading_rewards jsonb;
BEGIN
    trading_rewards_array = array_fill(NULL::jsonb, ARRAY[coalesce(jsonb_array_length(event_data->'tradingRewards'), 0)]::integer[]);

    /** Note that arrays are 1-indexed in PostgreSQL and empty arrays return NULL for array_length. */
    FOR i IN 1..coalesce(array_length(trading_rewards_array, 1), 0) LOOP
        trading_rewards_array[i] = dydx_process_trading_reward_event(
            jsonb_array_element(event_data->'tradingRewards', i-1),
            block_height,
            block_time,
            transaction_index,
            transaction_hash,
            event_index
        );
    END LOOP;

    RETURN jsonb_build_object(
        'trading_rewards',
        to_jsonb(trading_rewards_array)
    );
END;
$$ LANGUAGE plpgsql;
