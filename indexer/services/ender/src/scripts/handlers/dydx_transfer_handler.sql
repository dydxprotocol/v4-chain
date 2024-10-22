CREATE OR REPLACE FUNCTION dydx_transfer_handler(
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
    - asset: The existing asset in asset-model format (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/indexer/packages/postgres/src/models/asset-model.ts).
    - transfer: The new transfer in transfer-model format (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/indexer/packages/postgres/src/models/transfer-model.ts).

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    asset_record assets%ROWTYPE;
    recipient_subaccount_record subaccounts%ROWTYPE;
    recipient_wallet_record wallets%ROWTYPE;
    sender_wallet_record wallets%ROWTYPE;
    transfer_record transfers%ROWTYPE;
BEGIN
    asset_record."id" = event_data->>'assetId';
    SELECT * INTO asset_record FROM assets WHERE "id" = asset_record."id";

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Unable to find asset with assetId: %', asset_record."id";
    END IF;

    IF event_data->'recipient'->'subaccountId' IS NOT NULL THEN
        transfer_record."recipientSubaccountId" = dydx_uuid_from_subaccount_id(event_data->'recipient'->'subaccountId');

        recipient_subaccount_record."id" = transfer_record."recipientSubaccountId";
        recipient_subaccount_record."address" = event_data->'recipient'->'subaccountId'->>'owner';
        recipient_subaccount_record."subaccountNumber" = (event_data->'recipient'->'subaccountId'->'number')::int;
        recipient_subaccount_record."updatedAtHeight" = block_height;
        recipient_subaccount_record."updatedAt" = block_time;

        INSERT INTO subaccounts VALUES (recipient_subaccount_record.*)
        ON CONFLICT ("id") DO
            UPDATE
            SET
                "updatedAtHeight" = recipient_subaccount_record."updatedAtHeight",
                "updatedAt" = recipient_subaccount_record."updatedAt";

        recipient_wallet_record."address" = event_data->'recipient'->'subaccountId'->>'owner';
        recipient_wallet_record."totalTradingRewards" = '0';
        recipient_wallet_record."totalVolume" = '0';
        INSERT INTO wallets VALUES (recipient_wallet_record.*) ON CONFLICT DO NOTHING;
    END IF;

    IF event_data->'sender'->'subaccountId' IS NOT NULL THEN
        transfer_record."senderSubaccountId" = dydx_uuid_from_subaccount_id(event_data->'sender'->'subaccountId');
    END IF;

    IF event_data->'recipient'->'address' IS NOT NULL THEN
        transfer_record."recipientWalletAddress" = event_data->'recipient'->>'address';

        recipient_wallet_record."address" = transfer_record."recipientWalletAddress";
        recipient_wallet_record."totalTradingRewards" = '0';
        recipient_wallet_record."totalVolume" = '0';
        INSERT INTO wallets VALUES (recipient_wallet_record.*) ON CONFLICT DO NOTHING;
    END IF;

    IF event_data->'sender'->'address' IS NOT NULL THEN
        transfer_record."senderWalletAddress" = event_data->'sender'->>'address';

        sender_wallet_record."address" = transfer_record."senderWalletAddress";
        sender_wallet_record."totalTradingRewards" = '0';
        sender_wallet_record."totalVolume" = '0';
        INSERT INTO wallets VALUES (sender_wallet_record.*) ON CONFLICT DO NOTHING;
    END IF;

    transfer_record."assetId" = event_data->>'assetId';
    transfer_record."size" = dydx_trim_scale(dydx_from_jsonlib_long(event_data->'amount') * power(10, asset_record."atomicResolution")::numeric);
    transfer_record."eventId" = dydx_event_id_from_parts(block_height, transaction_index, event_index);
    transfer_record."transactionHash" = transaction_hash;
    transfer_record."createdAt" = block_time;
    transfer_record."createdAtHeight" = block_height;
    transfer_record."id" = dydx_uuid_from_transfer_parts(
        transfer_record."eventId",
        transfer_record."assetId",
        transfer_record."senderSubaccountId",
        transfer_record."recipientSubaccountId",
        transfer_record."senderWalletAddress",
        transfer_record."recipientWalletAddress");

    BEGIN
        INSERT INTO transfers VALUES (transfer_record.*);
    EXCEPTION
        WHEN check_violation THEN
            RAISE EXCEPTION 'Record: %, event: %', transfer_record, event_data;
    END;

    RETURN jsonb_build_object(
        'asset',
        dydx_to_jsonb(asset_record),
        'transfer',
        dydx_to_jsonb(transfer_record)
    );
END;
$$ LANGUAGE plpgsql;
