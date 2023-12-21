CREATE OR REPLACE FUNCTION dydx_process_trading_reward_event(
    trading_reward jsonb, block_height int, block_time timestamp, transaction_index int, transaction_hash text,
    event_index int) RETURNS jsonb AS $$
/**
  Parameters:
    - trading_reward: the trading reward to process, which should match AddressTradingReward (proto/dydxprotocol/indexer/events/events.proto#AddressTradingReward).
    - block_height: the height of the block being processing.
    - block_time: the time of the block being processed.
    - transaction_index: The transaction_index of the IndexerTendermintEvent after the conversion that takes into
        account the block_event (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/indexer/services/ender/src/lib/helper.ts#L41)
    - transaction_hash: The transaction hash corresponding to this event from the IndexerTendermintBlock 'tx_hashes'.
    - event_index: The 'event_index' of the IndexerTendermintEvent.
  Returns: trading_rewards row that was created as a result of the trading_reward jsonb event

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    trading_reward_record trading_rewards%ROWTYPE;
    wallet_record wallets%ROWTYPE;
    amount_in_human_readable numeric;
BEGIN
    amount_in_human_readable = dydx_from_serializable_int(trading_reward->>'denom_amount') * power(10, 18)::numeric;

    SELECT * INTO STRICT wallet_record FROM wallets WHERE "address" = trading_reward->>'owner';

    IF NO_DATA_FOUND THEN
        INSERT INTO wallets ("address", "totalTradingRewards") VALUES (
            trading_reward->>'owner',
            amount_in_human_readable);
    ELSE
        UPDATE wallets
        SET
            "totalTradingRewards" = wallet_record."totalTradingRewards" + amount_in_human_readable
        WHERE "address" = trading_reward->>'owner';
    END IF;

    INSERT INTO trading_rewards 
        ("id", "address", "blockTime", "blockHeight", "amount")
    VALUES (dydx_uuid_from_trading_rewards_parts(trading_reward->>'owner', block_height),
            trading_reward->>'owner',
            block_time,
            block_height,
            amount_in_human_readable)
    RETURNING * INTO trading_reward_record;

    RETURN trading_reward_record;
END;
$$ LANGUAGE plpgsql;
