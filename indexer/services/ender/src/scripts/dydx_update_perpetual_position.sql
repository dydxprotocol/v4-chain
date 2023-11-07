CREATE OR REPLACE FUNCTION dydx_update_perpetual_position(
    subaccount_uuid uuid,
    perpetual_id bigint,
    side text,
    size numeric,
    price numeric
) RETURNS perpetual_positions AS $$
DECLARE
    perpetual_position_record RECORD;
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

    -- Update the perpetual position record based on the side
    IF dydx_perpetual_position_and_order_side_matching(perpetual_position_record."side", side) THEN
        perpetual_position_record."sumOpen" := dydx_trim_scale(perpetual_position_record."sumOpen" + size);
        perpetual_position_record."entryPrice" := dydx_get_weighted_average(
            perpetual_position_record."entryPrice", perpetual_position_record."sumOpen", price, size
        );
    ELSE
        perpetual_position_record."sumClose" := dydx_trim_scale(perpetual_position_record."sumClose" + size);
        perpetual_position_record."exitPrice" := dydx_get_weighted_average(
            perpetual_position_record."exitPrice", perpetual_position_record."sumClose", price, size
        );
    END IF;

    -- Perform the actual update in the database
    UPDATE perpetual_positions
    SET
        "sumOpen" = perpetual_position_record."sumOpen",
        "entryPrice" = perpetual_position_record."entryPrice",
        "sumClose" = perpetual_position_record."sumClose",
        "exitPrice" = perpetual_position_record."exitPrice"
    WHERE "id" = perpetual_position_record.id;

    -- Return the updated perpetual position record as jsonb
    RETURN perpetual_position_record;
END;
$$ LANGUAGE plpgsql;
