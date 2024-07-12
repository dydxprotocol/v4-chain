CREATE OR REPLACE FUNCTION dydx_stateful_order_handler(
    block_height int, block_time timestamp, event_data jsonb) RETURNS jsonb AS $$
/**
  Parameters:
    - block_height: the height of the block being processing.
    - block_time: the time of the block being processed.
    - event_data: The 'data' field of the IndexerTendermintEvent (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/proto/dydxprotocol/indexer/indexer_manager/event.proto#L25)
        converted to JSON format. Conversion to JSON is expected to be done by JSON.stringify.
  Returns: JSON object containing fields:
    - order: The upserted order in order-model format (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/indexer/packages/postgres/src/models/order-model.ts).

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    QUOTE_CURRENCY_ATOMIC_RESOLUTION constant numeric = -6;

    order_ jsonb;
    order_id jsonb;
    clob_pair_id bigint;
    subaccount_id uuid;
    perpetual_market_record perpetual_markets%ROWTYPE;
    order_record orders%ROWTYPE;
    subaccount_record subaccounts%ROWTYPE;
BEGIN
    /** TODO(IND-334): Remove after deprecating StatefulOrderPlacementEvent. */
    IF event_data->'orderPlace' IS NOT NULL OR event_data->'longTermOrderPlacement' IS NOT NULL OR event_data->'conditionalOrderPlacement' IS NOT NULL THEN
        order_ = coalesce(event_data->'orderPlace'->'order', event_data->'longTermOrderPlacement'->'order', event_data->'conditionalOrderPlacement'->'order');
        clob_pair_id = (order_->'orderId'->'clobPairId')::bigint;

        perpetual_market_record = dydx_get_perpetual_market_for_clob_pair(clob_pair_id);

        /**
          Calculate sizes, prices, and fill amounts.

          TODO(IND-238): Extract out calculation of quantums and subticks to their own SQL functions.
        */
        order_record."id" = dydx_uuid_from_order_id(order_->'orderId');
        order_record."subaccountId" = dydx_uuid_from_subaccount_id(order_->'orderId'->'subaccountId');
        order_record."clientId" = jsonb_extract_path_text(order_, 'orderId', 'clientId')::bigint;
        order_record."clobPairId" = clob_pair_id;
        order_record."side" = dydx_from_protocol_order_side(order_->'side');
        order_record."size" = dydx_trim_scale(dydx_from_jsonlib_long(order_->'quantums') *
                                              power(10, perpetual_market_record."atomicResolution")::numeric);
        order_record."totalFilled" = 0;
        order_record."price" = dydx_trim_scale(dydx_from_jsonlib_long(order_->'subticks') *
                                               power(10, perpetual_market_record."quantumConversionExponent" +
                                                         QUOTE_CURRENCY_ATOMIC_RESOLUTION -
                                                         perpetual_market_record."atomicResolution")::numeric);
        order_record."timeInForce" = dydx_from_protocol_time_in_force(order_->'timeInForce');
        order_record."reduceOnly" = (order_->>'reduceOnly')::boolean;
        order_record."orderFlags" = (order_->'orderId'->'orderFlags')::bigint;
        order_record."goodTilBlockTime" = to_timestamp((order_->'goodTilBlockTime')::double precision);
        order_record."clientMetadata" = (order_->'clientMetadata')::bigint;
        order_record."createdAtHeight" = block_height;
        order_record."updatedAt" = block_time;
        order_record."updatedAtHeight" = block_height;

        CASE
            WHEN event_data->'conditionalOrderPlacement' IS NOT NULL THEN
                order_record."type" = dydx_protocol_condition_type_to_order_type(order_->'conditionType');
                order_record."status" = 'UNTRIGGERED';
                order_record."triggerPrice" = dydx_trim_scale(dydx_from_jsonlib_long(order_->'conditionalOrderTriggerSubticks') *
                                                              power(10, perpetual_market_record."quantumConversionExponent" +
                                                                        QUOTE_CURRENCY_ATOMIC_RESOLUTION -
                                                                        perpetual_market_record."atomicResolution")::numeric);
            ELSE
                order_record."type" = 'LIMIT';
                order_record."status" = 'OPEN';
        END CASE;

        INSERT INTO orders VALUES (order_record.*) ON CONFLICT ("id") DO
            UPDATE SET
                       "subaccountId" = order_record."subaccountId",
                       "clientId" = order_record."clientId",
                       "clobPairId" = order_record."clobPairId",
                       "side" = order_record."side",
                       "size" = order_record."size",
                       "totalFilled" = order_record."totalFilled",
                       "price" = order_record."price",
                       "timeInForce" = order_record."timeInForce",
                       "reduceOnly" = order_record."reduceOnly",
                       "orderFlags" = order_record."orderFlags",
                       "goodTilBlockTime" = order_record."goodTilBlockTime",
                       "clientMetadata" = order_record."clientMetadata",
                       "createdAtHeight" = order_record."createdAtHeight",
                       "updatedAt" = order_record."updatedAt",
                       "updatedAtHeight" = order_record."updatedAtHeight",
                       "type" = order_record."type",
                       "status" = order_record."status",
                       "triggerPrice" = order_record."triggerPrice"
        RETURNING * INTO order_record;

        RETURN jsonb_build_object(
                'order',
                dydx_to_jsonb(order_record),
                'perpetual_market',
                dydx_to_jsonb(perpetual_market_record)
            );
    ELSIF event_data->'conditionalOrderTriggered' IS NOT NULL OR event_data->'orderRemoval' IS NOT NULL THEN
        CASE
            WHEN event_data->'conditionalOrderTriggered' IS NOT NULL THEN
                order_id = event_data->'conditionalOrderTriggered'->'triggeredOrderId';
                order_record."status" = 'OPEN';
            ELSE
                order_id = event_data->'orderRemoval'->'removedOrderId';
                order_record."status" = 'CANCELED';
        END CASE;

        clob_pair_id = (order_id->'clobPairId')::bigint;
        perpetual_market_record = dydx_get_perpetual_market_for_clob_pair(clob_pair_id);

        subaccount_id = dydx_uuid_from_subaccount_id(order_id->'subaccountId');
        SELECT * INTO subaccount_record FROM subaccounts WHERE "id" = subaccount_id;
        IF NOT FOUND THEN
            RAISE EXCEPTION 'Subaccount for order not found: %', order_;
        END IF;

        order_record."id" = dydx_uuid_from_order_id(order_id);
        order_record."updatedAt" = block_time;
        order_record."updatedAtHeight" = block_height;
        UPDATE orders
        SET
            "status" = order_record."status",
            "updatedAt" = order_record."updatedAt",
            "updatedAtHeight" = order_record."updatedAtHeight"
        WHERE "id" = order_record."id"
        RETURNING * INTO order_record;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'Unable to update order status with orderId: %', dydx_uuid_from_order_id(order_id);
        END IF;

        RETURN jsonb_build_object(
                'order',
                dydx_to_jsonb(order_record),
                'perpetual_market',
                dydx_to_jsonb(perpetual_market_record),
                'subaccount',
                dydx_to_jsonb(subaccount_record)
            );
    ELSE
        RAISE EXCEPTION 'Unkonwn sub-event type %', event_data;
    END IF;
END;
$$ LANGUAGE plpgsql;
