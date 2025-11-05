CREATE OR REPLACE FUNCTION dydx_liquidation_fill_handler_per_order(
    field text, block_height int, block_time timestamp, event_data jsonb, event_index int, transaction_index int,
    transaction_hash text, fill_liquidity text, fill_type text, usdc_asset_id text) RETURNS jsonb AS $$
DECLARE
    order_ jsonb;
    maker_order jsonb;
    clob_pair_id bigint;
    subaccount_uuid uuid;
    perpetual_market_record perpetual_markets%ROWTYPE;
    market_record markets%ROWTYPE;
    order_record orders%ROWTYPE;
    fill_record fills%ROWTYPE;
    perpetual_position_record perpetual_positions%ROWTYPE;
    read_perpetual_position_record perpetual_positions%ROWTYPE;
    asset_record assets%ROWTYPE;
    order_uuid uuid;
    order_side text;
    order_size numeric;
    order_price numeric;
    order_client_metadata bigint;
    fee numeric;
    builder_fee numeric;
    builder_address text;
    affiliate_rev_share numeric;
    fill_amount numeric;
    total_filled numeric;
    maker_price numeric;
    event_id bytea;
    order_router_address text;
    order_router_fee numeric;
    snap_size_before numeric;
    snap_entry_before numeric;
    snap_side_before text;
/**
  Parameters:
    - field: the field storing the order to process.
    - block_height: the height of the block being processing.
    - block_time: the time of the block being processed.
    - event_data: The 'data' field of the IndexerTendermintEvent (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/proto/dydxprotocol/indexer/indexer_manager/event.proto#L25)
        converted to JSON format. Conversion to JSON is expected to be done by JSON.stringify.
    - event_index: The 'event_index' of the IndexerTendermintEvent.
    - transaction_index: The transaction_index of the IndexerTendermintEvent after the conversion that takes into
        account the block_event (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/indexer/services/ender/src/lib/helper.ts#L41)
    - transaction_hash: The transaction hash corresponding to this event from the IndexerTendermintBlock 'tx_hashes'.
    - fill_liquidity: The liquidity for the fill record.
    - fill_type: The type for the fill record.
    - usdc_asset_id: The USDC asset id.
  Returns: JSON object containing fields:
    - order: The updated order in order-model format (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/indexer/packages/postgres/src/models/order-model.ts).
        Only returned if field == 'makerOrder'.
    - fill: The updated fill in fill-model format (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/indexer/packages/postgres/src/models/fill-model.ts).
    - perpetual_market: The perpetual market for the order in perpetual-market-model format (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/indexer/packages/postgres/src/models/perpetual-market-model.ts).
    - perpetual_position: The updated perpetual position in perpetual-position-model format (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/indexer/packages/postgres/src/models/perpetual-position-model.ts).

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    order_ = event_data->field;
    maker_order = event_data->'makerOrder';

    IF field = 'makerOrder' THEN
        clob_pair_id = jsonb_extract_path(order_, 'orderId', 'clobPairId')::bigint;
    ELSE
        clob_pair_id = jsonb_extract_path(order_, 'clobPairId')::bigint;
    END IF;

    perpetual_market_record = dydx_get_perpetual_market_for_clob_pair(clob_pair_id);
    market_record = dydx_get_market_for_id(perpetual_market_record."marketId");

    BEGIN
        SELECT * INTO STRICT asset_record FROM assets WHERE "id" = usdc_asset_id;
    EXCEPTION
        WHEN NO_DATA_FOUND THEN
            RAISE EXCEPTION 'Unable to find asset with id %', usdc_asset_id;
    END;

    /**
      Calculate sizes, prices, and fill amounts.

      TODO(IND-238): Extract out calculation of quantums and subticks to their own SQL functions.
    */
    fill_amount = dydx_trim_scale(dydx_from_jsonlib_long(event_data->'fillAmount') *
                                  power(10, perpetual_market_record."atomicResolution")::numeric);
    maker_price = dydx_trim_scale(dydx_from_jsonlib_long(maker_order->'subticks') *
                                  power(10, perpetual_market_record."quantumConversionExponent" +
                                            asset_record."atomicResolution" -
                                            perpetual_market_record."atomicResolution")::numeric);
    total_filled = dydx_trim_scale(dydx_get_total_filled(fill_liquidity, event_data) *
                                   power(10, perpetual_market_record."atomicResolution")::numeric);
    fee = dydx_trim_scale(dydx_get_fee(fill_liquidity, event_data) *
                          power(10, asset_record."atomicResolution")::numeric);
    builder_fee = dydx_trim_scale(dydx_get_builder_fee(fill_liquidity, event_data) *
                                  power(10, asset_record."atomicResolution")::numeric);
    builder_address = dydx_get_builder_address(fill_liquidity, event_data);
    affiliate_rev_share = dydx_trim_scale(dydx_from_jsonlib_long(event_data->'affiliateRevShare') *
                                    power(10, asset_record."atomicResolution")::numeric);
    order_price = dydx_trim_scale(dydx_from_jsonlib_long(order_->'subticks') *
                                  power(10, perpetual_market_record."quantumConversionExponent" +
                                            asset_record."atomicResolution" -
                                            perpetual_market_record."atomicResolution")::numeric);
    order_router_fee = dydx_trim_scale(dydx_get_order_router_fee(fill_liquidity, event_data) *
                                    power(10, asset_record."atomicResolution")::numeric);
    order_router_address = dydx_get_order_router_address(fill_liquidity, event_data);

    IF field = 'makerOrder' THEN
        order_uuid = dydx_uuid_from_order_id(order_->'orderId');
        subaccount_uuid = dydx_uuid_from_subaccount_id(jsonb_extract_path(order_, 'orderId', 'subaccountId'));
        order_client_metadata = (order_->'clientMetadata')::bigint;
        order_side = dydx_from_protocol_order_side(order_->'side');
    ELSE
        order_uuid = NULL;
        subaccount_uuid = dydx_uuid_from_subaccount_id(jsonb_extract_path(order_, 'liquidated'));
        order_client_metadata = NULL;
        /** Liquidation order proto has an isBuy property rather than a side property **/
        order_side = CASE WHEN (order_->'isBuy')::bool THEN 'BUY' ELSE 'SELL' END;
    END IF;

    IF field = 'makerOrder' THEN
        order_size = dydx_trim_scale(dydx_from_jsonlib_long(order_->'quantums') *
                                     power(10, perpetual_market_record."atomicResolution")::numeric);

        /** Upsert the order, populating the order_record fields with what will be in the database. */
        SELECT * INTO order_record FROM orders WHERE "id" = order_uuid;
        order_record."side" = order_side;
        order_record."size" = order_size;
        order_record."price" = order_price;
        order_record."timeInForce" = dydx_from_protocol_time_in_force(order_->'timeInForce');
        order_record."reduceOnly" = (order_->>'reduceOnly')::boolean;
        order_record."orderFlags" = jsonb_extract_path(order_, 'orderId', 'orderFlags')::bigint;
        order_record."goodTilBlock" = (order_->'goodTilBlock')::bigint;
        order_record."goodTilBlockTime" = to_timestamp((order_->'goodTilBlockTime')::double precision);
        order_record."clientMetadata" = order_client_metadata;
        order_record."updatedAt" = block_time;
        order_record."updatedAtHeight" = block_height;
        order_record."orderRouterAddress" = order_->'orderRouterAddress'::text;

        IF FOUND THEN
            IF jsonb_extract_path(order_, 'orderId', 'orderFlags')::bigint = constants.order_flag_twap_suborder() THEN
                order_record."price" = dydx_get_weighted_average(order_record."price", order_record."totalFilled", maker_price, fill_amount);
                order_record."totalFilled" = order_record."totalFilled" + fill_amount;

                order_record."status" = dydx_get_order_status(order_record."totalFilled", order_record."size", 'NOT_CANCELED', jsonb_extract_path(order_, 'orderId', 'orderFlags')::bigint, order_record."timeInForce");
                UPDATE orders
                SET
                    "status" = order_record."status",
                    "price" = order_record."price",
                    "updatedAt" = order_record."updatedAt",
                    "updatedAtHeight" = order_record."updatedAtHeight",
                    "totalFilled" = order_record."totalFilled" -- keep track of fill amount for the parent order
                WHERE "id" = order_uuid;

            ELSE
                order_record."totalFilled" = total_filled;
                order_record."status" = dydx_get_order_status(total_filled, order_record.size, 'NOT_CANCELED', order_record."orderFlags", order_record."timeInForce");

                UPDATE orders
                SET
                    "side" = order_record."side",
                    "size" = order_record."size",
                    "totalFilled" = order_record."totalFilled",
                    "price" = order_record."price",
                    "status" = order_record."status",
                    "orderFlags" = order_record."orderFlags",
                    "goodTilBlock" = order_record."goodTilBlock",
                    "goodTilBlockTime" = order_record."goodTilBlockTime",
                    "timeInForce" = order_record."timeInForce",
                    "reduceOnly" = order_record."reduceOnly",
                    "clientMetadata" = order_record."clientMetadata",
                    "updatedAt" = order_record."updatedAt",
                    "updatedAtHeight" = order_record."updatedAtHeight",
                    "builderAddress" = order_record."builderAddress",
                    "feePpm" = order_record."feePpm",
                    "orderRouterAddress" = order_record."orderRouterAddress"
                WHERE id = order_uuid;
            END IF;
        ELSE
            order_record."id" = order_uuid;
            order_record."subaccountId" = subaccount_uuid;
            order_record."clientId" = jsonb_extract_path_text(order_, 'orderId', 'clientId')::bigint;
            order_record."clobPairId" = clob_pair_id;
            order_record."side" = order_side;
            order_record."type" = 'LIMIT';

            order_record."totalFilled" = fill_amount;
            order_record."status" = dydx_get_order_status(fill_amount, order_size, 'NOT_CANCELED', order_record."orderFlags", order_record."timeInForce");
            order_record."createdAtHeight" = block_height;

            IF jsonb_extract_path(order_, 'orderId', 'orderFlags')::bigint = constants.order_flag_twap_suborder() THEN
                -- This is a handled case but is not expected for twap. Parent orders should always exist
                RAISE WARNING 'Twap suborders should always have a parent order. Order UUID: %', order_uuid;
                order_record."orderFlags" = constants.order_flag_twap(); -- Twap suborders should be mapped to their parent order.
                order_record."type" = 'TWAP';
                order_record."price" = maker_price;
            END IF;

            order_record."duration" = NULL;
            order_record."interval" = NULL;
            order_record."priceTolerance" = NULL;

            INSERT INTO orders
            ("id", "subaccountId", "clientId", "clobPairId", "side", "size", "totalFilled", "price", "type",
             "status", "timeInForce", "reduceOnly", "orderFlags", "goodTilBlock", "goodTilBlockTime", "createdAtHeight",
             "clientMetadata", "triggerPrice", "updatedAt", "updatedAtHeight", "builderAddress", "feePpm",
             "orderRouterAddress", "duration", "interval", "priceTolerance")
            VALUES (
                order_record."id", order_record."subaccountId", order_record."clientId", order_record."clobPairId",
                order_record."side", order_record."size", order_record."totalFilled", order_record."price", order_record."type",
                order_record."status", order_record."timeInForce", order_record."reduceOnly", order_record."orderFlags",
                order_record."goodTilBlock", order_record."goodTilBlockTime", order_record."createdAtHeight",
                order_record."clientMetadata", order_record."triggerPrice", order_record."updatedAt", order_record."updatedAtHeight",
                order_record."builderAddress", order_record."feePpm", order_record."orderRouterAddress", order_record."duration",
                order_record."interval", order_record."priceTolerance"
            );
        END IF;
    END IF;

    -- Retrieve the latest perpetual position record.
    SELECT * INTO read_perpetual_position_record
    FROM perpetual_positions
    WHERE "subaccountId" = subaccount_uuid
      AND "perpetualId" = perpetual_market_record."id"
    ORDER BY "openEventId" DESC
    LIMIT 1;

    snap_size_before = COALESCE(ABS(read_perpetual_position_record."sumOpen"), 0) - COALESCE(ABS(read_perpetual_position_record."sumClose"), 0);
    snap_entry_before = NULLIF(read_perpetual_position_record."entryPrice", 0);
    snap_side_before = read_perpetual_position_record."side";

    /* Insert the associated fill record for this order_fill event. */
    event_id = dydx_event_id_from_parts(
            block_height, transaction_index, event_index);
    INSERT INTO fills
        ("id", "subaccountId", "side", "liquidity", "type", "clobPairId", "orderId", "size", "price", "quoteAmount",
         "eventId", "transactionHash", "createdAt", "createdAtHeight", "clientMetadata", "fee", "affiliateRevShare",
         "builderFee", "builderAddress", "orderRouterFee", "orderRouterAddress", "positionSizeBefore", "entryPriceBefore", "positionSideBefore")
    VALUES (dydx_uuid_from_fill_event_parts(event_id, fill_liquidity),
            subaccount_uuid,
            order_side,
            fill_liquidity,
            fill_type,
            clob_pair_id,
            order_uuid,
            fill_amount,
            maker_price,
            dydx_trim_scale(fill_amount * maker_price),
            event_id,
            transaction_hash,
            block_time,
            block_height,
            order_client_metadata,
            fee,
            affiliate_rev_share,
            NULLIF(builder_fee, 0),
            NULLIF(builder_address, ''),
            NULLIF(order_router_fee, 0),
            NULLIF(order_router_address, ''),
            snap_size_before,
            snap_entry_before,
            snap_side_before)
    RETURNING * INTO fill_record;

    /* Upsert the perpetual_position record for this order_fill event. */
    SELECT * INTO perpetual_position_record FROM perpetual_positions WHERE "subaccountId" = subaccount_uuid
                                                                       AND "perpetualId" = perpetual_market_record."id"
    ORDER BY "openEventId" DESC;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Unable to find existing perpetual position, subaccountId: %, perpetualId: %', subaccount_uuid, perpetual_market_record."id";
    END IF;
    DECLARE
        sum_open numeric = perpetual_position_record."sumOpen";
        entry_price numeric = perpetual_position_record."entryPrice";
        sum_close numeric = perpetual_position_record."sumClose";
        exit_price numeric = perpetual_position_record."exitPrice";
    BEGIN
        IF dydx_perpetual_position_and_order_side_matching(
                perpetual_position_record."side", order_side) THEN
            sum_open = dydx_trim_scale(perpetual_position_record."sumOpen" + fill_amount);
            entry_price = dydx_get_weighted_average(
                    perpetual_position_record."entryPrice", perpetual_position_record."sumOpen",
                    maker_price, fill_amount);
            perpetual_position_record."sumOpen" = sum_open;
            perpetual_position_record."entryPrice" = entry_price;
        ELSE
            sum_close = dydx_trim_scale(perpetual_position_record."sumClose" + fill_amount);
            exit_price = dydx_get_weighted_average(
                    perpetual_position_record."exitPrice", perpetual_position_record."sumClose",
                    maker_price, fill_amount);
            perpetual_position_record."sumClose" = sum_close;
            perpetual_position_record."exitPrice" = exit_price;
        END IF;
        UPDATE perpetual_positions
        SET
            "sumOpen" = sum_open,
            "entryPrice" = entry_price,
            "sumClose" = sum_close,
            "exitPrice" = exit_price
        WHERE "id" = perpetual_position_record.id;
    END;

    IF field = 'makerOrder' THEN
        RETURN jsonb_build_object(
                'order',
                dydx_to_jsonb(order_record),
                'fill',
                dydx_to_jsonb(fill_record),
                'perpetual_market',
                dydx_to_jsonb(perpetual_market_record),
                'market',
                dydx_to_jsonb(market_record),
                'perpetual_position',
                dydx_to_jsonb(perpetual_position_record)
            );
    ELSE
        RETURN jsonb_build_object(
                'fill',
                dydx_to_jsonb(fill_record),
                'perpetual_market',
                dydx_to_jsonb(perpetual_market_record),
                'market',
                dydx_to_jsonb(market_record),
                'perpetual_position',
                dydx_to_jsonb(perpetual_position_record)
            );
    END IF;
END;
$$ LANGUAGE plpgsql;
