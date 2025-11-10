CREATE OR REPLACE FUNCTION dydx_subaccount_update_handler(
    block_height int, block_time timestamp, event_data jsonb, event_index int, transaction_index int)
    RETURNS jsonb AS $$
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
    - subaccount: The upserted subaccount in subaccount-model format (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/indexer/packages/postgres/src/models/subaccount-model.ts).
    - perpetual_positions: A JSON array of upserted perpetual positions in perpetual-position-model format (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/indexer/packages/postgres/src/models/perpetual-position-model.ts).
    - asset_positions: A JSON array of upserted asset positions in asset-position-model format (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/indexer/packages/postgres/src/models/asset-position-model.ts).
    - markets: A JSON object mapping market ids to market records.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    QUOTE_CURRENCY_ATOMIC_RESOLUTION constant numeric = -6;
    event_id bytea;
    subaccount_record subaccounts;
    perpetual_position_update jsonb;
    perpetual_position_record_updates jsonb[];
    asset_position_update jsonb;
    asset_position_record_updates jsonb[];
    market_map jsonb;
BEGIN
    event_id = dydx_event_id_from_parts(
            block_height, transaction_index, event_index);
    perpetual_position_record_updates = array[]::jsonb[];
    asset_position_record_updates = array[]::jsonb[];

    subaccount_record."id" = dydx_uuid_from_subaccount_id(event_data->'subaccountId');
    subaccount_record."address" = jsonb_extract_path_text(event_data, 'subaccountId', 'owner');
    subaccount_record."subaccountNumber" = jsonb_extract_path(event_data, 'subaccountId', 'number')::int;
    subaccount_record."updatedAtHeight" = block_height;
    subaccount_record."updatedAt" = block_time;
    INSERT INTO subaccounts (
        "id", "address", "subaccountNumber", "updatedAtHeight", "updatedAt"
    ) VALUES (
        subaccount_record."id", subaccount_record."address", 
        subaccount_record."subaccountNumber", subaccount_record."updatedAtHeight", 
        subaccount_record."updatedAt"
    )
    ON CONFLICT ("id") DO UPDATE
        SET "updatedAtHeight" = subaccount_record."updatedAtHeight", "updatedAt" = subaccount_record."updatedAt";

    -- Process all the perpetual position updates that are part of the update.
    FOR perpetual_position_update IN SELECT * FROM jsonb_array_elements(event_data->'updatedPerpetualPositions') LOOP
        DECLARE
            perpetual_id bigint;
            perpetual_market perpetual_markets%ROWTYPE;
            perpetual_position_record perpetual_positions%ROWTYPE;
            _size numeric;
            max_size numeric;
            side text;
            settled_funding numeric;
            perpetual_position_found boolean;
            existing_funding numeric;
            new_settled_funding numeric;
        BEGIN
            perpetual_id = (perpetual_position_update->'perpetualId')::bigint;
            SELECT * INTO perpetual_market FROM perpetual_markets WHERE id = perpetual_id;
            SELECT * INTO perpetual_position_record FROM perpetual_positions
                                                    WHERE "subaccountId" = subaccount_record."id"
                                                      AND "perpetualId" = perpetual_id
                                                      AND "status" = 'OPEN';
            perpetual_position_found = FOUND;
            _size = dydx_trim_scale(dydx_from_serializable_int(perpetual_position_update->'quantums') *
                   power(10, perpetual_market."atomicResolution")::numeric);
            side = CASE WHEN _size > 0 THEN 'LONG' ELSE 'SHORT' END;
            existing_funding = CASE WHEN perpetual_position_found THEN perpetual_position_record."settledFunding" ELSE 0 END;
            settled_funding = dydx_trim_scale(-dydx_from_serializable_int(perpetual_position_update->'fundingPayment')
                                                  * power(10, QUOTE_CURRENCY_ATOMIC_RESOLUTION)::numeric + existing_funding);
            new_settled_funding = CASE WHEN (perpetual_position_found AND perpetual_position_record.side != side AND _size != 0)
                                        THEN 0
                                        ELSE settled_funding END;

            -- Handle updating the existing perpetual record.
            IF perpetual_position_found AND (_size = 0 OR perpetual_position_record.side != side) THEN
                -- Close the existing position since the new size is 0 or the sides changed.
                IF perpetual_position_record.status = 'CLOSED' THEN
                    RAISE EXCEPTION 'Unable to close % because position is closed', perpetual_position_record."id";
                END IF;
                UPDATE perpetual_positions SET "closedAt" = block_time,
                                              "closedAtHeight" = block_height,
                                              "closeEventId" = event_id,
                                              "lastEventId" = event_id,
                                              "settledFunding" = settled_funding,
                                              "status" = 'CLOSED',
                                              "size" = 0
                                          WHERE "id" = perpetual_position_record."id"
                                          RETURNING * INTO perpetual_position_record;
                perpetual_position_record_updates = array_append(
                    perpetual_position_record_updates, dydx_to_jsonb(perpetual_position_record));
            ELSEIF perpetual_position_found AND perpetual_position_record.side = side THEN
                max_size = CASE
                    WHEN perpetual_position_record."maxSize" IS NOT NULL
                             AND perpetual_position_record."maxSize" >= _size
                    THEN perpetual_position_record."maxSize"
                    ELSE _size END;
                -- Since the sides match update the existing position
                UPDATE perpetual_positions SET "size" = _size,
                                               "lastEventId" = event_id,
                                               "settledFunding" = settled_funding,
                                               "maxSize" = max_size
                                            WHERE "id" = perpetual_position_record."id" RETURNING * INTO perpetual_position_record;
                perpetual_position_record_updates = array_append(
                    perpetual_position_record_updates, dydx_to_jsonb(perpetual_position_record));
            END IF;

            -- Insert a new perpetual record if necessary.
            IF NOT perpetual_position_found
                   OR (perpetual_position_found AND perpetual_position_record.side != side AND _size != 0) THEN
                -- Since no perpetual position was found or we closed an existing perpetual position because it changed
                -- sides we must create a new perpetual position.
                perpetual_position_record."id" = dydx_uuid_from_perpetual_position_parts(subaccount_record.id, event_id);
                perpetual_position_record."subaccountId" = subaccount_record.id;
                perpetual_position_record."perpetualId" = perpetual_id;
                perpetual_position_record."side" = side;
                perpetual_position_record."status" = 'OPEN';
                perpetual_position_record."size" = _size;
                perpetual_position_record."maxSize" = _size;
                perpetual_position_record."entryPrice" = 0;
                perpetual_position_record."exitPrice" = NULL;
                perpetual_position_record."sumOpen" = 0;
                perpetual_position_record."sumClose" = 0;
                perpetual_position_record."createdAt" = block_time;
                perpetual_position_record."closedAt" = NULL;
                perpetual_position_record."createdAtHeight" = block_height;
                perpetual_position_record."closedAtHeight" = NULL;
                perpetual_position_record."openEventId" = event_id;
                perpetual_position_record."closeEventId" = NULL;
                perpetual_position_record."lastEventId" = event_id;
                perpetual_position_record."settledFunding" = new_settled_funding;
                INSERT INTO perpetual_positions VALUES (perpetual_position_record.*);
                perpetual_position_record_updates = array_append(
                    perpetual_position_record_updates, dydx_to_jsonb(perpetual_position_record));
            END IF;
        END;
    END LOOP;

    -- Process all the asset position updates that are part of the update.
    FOR asset_position_update IN SELECT * FROM jsonb_array_elements(event_data->'updatedAssetPositions') LOOP
        DECLARE
            asset_id text;
            asset_record assets%ROWTYPE;
            asset_position_record asset_positions%ROWTYPE;
            size numeric;
        BEGIN
            asset_id = asset_position_update->>'assetId';
            SELECT * INTO asset_record FROM assets WHERE "id" = asset_id;
            IF NOT FOUND THEN
                RAISE EXCEPTION 'Unable to find asset with id %', asset_id;
            END IF;

            size = dydx_trim_scale(dydx_from_serializable_int(asset_position_update->'quantums') *
                   power(10, asset_record."atomicResolution")::numeric);
            asset_position_record.id = dydx_uuid_from_asset_position_parts(subaccount_record.id, asset_id);
            asset_position_record."subaccountId" = subaccount_record."id";
            asset_position_record."assetId" = asset_id;
            asset_position_record."size" = abs(size);
            asset_position_record."isLong" = size > 0;

            INSERT INTO asset_positions VALUES (asset_position_record.*)
                                        ON CONFLICT ("id") DO UPDATE
                                            SET size = asset_position_record.size,
                                                "isLong" = asset_position_record."isLong";
            asset_position_record_updates = array_append(
                asset_position_record_updates, dydx_to_jsonb(asset_position_record));
        END;
    END LOOP;

    -- Fetch all markets
    market_map := (SELECT jsonb_object_agg(id, dydx_to_jsonb(markets)) FROM markets);
    RETURN jsonb_build_object(
        'subaccount', dydx_to_jsonb(subaccount_record),
        'perpetual_positions', to_jsonb(perpetual_position_record_updates),
        'asset_positions', to_jsonb(asset_position_record_updates),
        'markets', market_map
        );
END;
$$ LANGUAGE plpgsql;
