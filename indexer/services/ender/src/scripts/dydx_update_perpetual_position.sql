CREATE OR REPLACE FUNCTION dydx_update_perpetual_position(
    subaccount_uuid uuid,
    perpetual_id bigint,
    side text,
    size numeric,
    price numeric
) RETURNS perpetual_positions AS $$
DECLARE
    perpetual_position_record RECORD;
    sum_open numeric;
    entry_price numeric;
    sum_close numeric;
    exit_price numeric;
BEGIN
    -- Retrieve the latest perpetual position record
    SELECT * INTO perpetual_position_record
    FROM perpetual_positions
    WHERE "subaccountId" = subaccount_uuid
      AND "perpetualId" = perpetual_id
    ORDER BY "createdAtHeight" DESC
    LIMIT 1;

    -- Check if a perpetual position record was found
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Unable to find existing perpetual position, subaccountId: %, perpetualId: %', subaccount_uuid, perpetual_id;
    END IF;

    sum_open = perpetual_position_record."sumOpen";
    entry_price = perpetual_position_record."entryPrice";
    sum_close = perpetual_position_record."sumClose";
    exit_price = perpetual_position_record."exitPrice";

    -- Update the perpetual position record based on the side
    IF dydx_perpetual_position_and_order_side_matching(perpetual_position_record."side", side) THEN
        sum_open := dydx_trim_scale(perpetual_position_record."sumOpen" + size);
        entry_price := dydx_get_weighted_average(
            perpetual_position_record."entryPrice", perpetual_position_record."sumOpen", price, size
        );
        perpetual_position_record."sumOpen" = sum_open;
        perpetual_position_record."entryPrice" = entry_price;
    ELSE
        sum_close := dydx_trim_scale(perpetual_position_record."sumClose" + size);
        exit_price := dydx_get_weighted_average(
            perpetual_position_record."exitPrice", perpetual_position_record."sumClose", price, size
        );
        perpetual_position_record."sumClose" = sum_close;
        perpetual_position_record."exitPrice" = exit_price;
    END IF;

    -- Perform the actual update in the database
    UPDATE perpetual_positions
    SET
        "sumOpen" = sum_open,
        "entryPrice" = entry_price,
        "sumClose" = sum_close,
        "exitPrice" = exit_price
    WHERE "id" = perpetual_position_record.id;

    -- Return the updated perpetual position record as jsonb
    RETURN perpetual_position_record;
END;
$$ LANGUAGE plpgsql;
