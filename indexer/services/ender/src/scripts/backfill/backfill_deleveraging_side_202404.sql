CREATE OR REPLACE FUNCTION backfill_deleveraging_side() RETURNS void AS $$
/**
  Backfills deleveraging fills and perpetual updates due to incorrect deleveraging side.
*/
DECLARE
    fill_record fills%ROWTYPE;
BEGIN
    FOR fill_record IN (SELECT * FROM fills WHERE type IN ('DELEVERAGED', 'OFFSETTING')) 
    LOOP
        DECLARE
            subaccount_uuid uuid;
            clob_pair_id bigint;
            perpetual_market_record perpetual_markets%ROWTYPE;
            perpetual_id bigint;
            incorrect_side text;
            correct_side text;
            size numeric;
            neg_size numeric;
            price numeric;
            reverted_perpetual_position_record perpetual_positions%ROWTYPE;
            updated_perpetual_position_record perpetual_positions%ROWTYPE;
        BEGIN
            subaccount_uuid = fill_record."subaccountId";
            price = fill_record."price";

            -- get perpetual id
            clob_pair_id = fill_record."clobPairId";
            perpetual_market_record = dydx_get_perpetual_market_for_clob_pair(clob_pair_id);
            perpetual_id = perpetual_market_record."id";

            -- get the size of the fill.
            size = fill_record."size";
            neg_size = size * -1;

            -- side was incorrectly flipped for both deleveraged and offsetting 
            -- subaccounts.
            incorrect_side = fill_record."side";
            correct_side = CASE WHEN incorrect_side = 'BUY' THEN 'SELL' ELSE 'BUY' END;

            -- revert the perpetual position update for incorrect side.
            reverted_perpetual_position_record = dydx_update_perpetual_position_aggregate_fields(
                subaccount_uuid,
                perpetual_id,
                -- use incorrect side to revert the perpetual position update
                incorrect_side, 
                -- use negative size to revert the perpetual position update
                neg_size, 
                price);

            -- apply the perpetual position update for correct side
            updated_perpetual_position_record = dydx_update_perpetual_position_aggregate_fields(
                subaccount_uuid,
                perpetual_id,
                -- use correct side to apply the perpetual position update
                correct_side, 
                -- use size to apply the perpetual position update
                size, 
                price);

            -- update the fill record with the correct side
            UPDATE fills
            SET
                "side" = correct_side
            WHERE "id" = fill_record."id";
        END;
    END LOOP;
END;
$$ LANGUAGE plpgsql;
