CREATE OR REPLACE FUNCTION dydx_apply_fill_realized_effects(
    position_id uuid,
    fill_side text,
    fill_size numeric,
    fill_price numeric,
    fill_fee numeric,
    pos_side_before text,
    pos_size_before numeric,
    entry_price_before numeric
) RETURNS void AS $$
DECLARE
    is_reducing boolean;
    closing_amount numeric;
    pnl numeric := 0;
BEGIN
    IF pos_size_before IS NULL OR pos_size_before = 0 THEN
        -- opening trade: only fees realize
        UPDATE perpetual_positions
            SET "totalRealizedPnl" =
                COALESCE("totalRealizedPnl", 0) - COALESCE(fill_fee, 0)
            WHERE "id" = position_id;
        RETURN;
    END IF;

    is_reducing := (pos_side_before = 'LONG'  AND fill_side = 'SELL')
                OR (pos_side_before = 'SHORT' AND fill_side = 'BUY');

    IF NOT is_reducing THEN
        -- increasing: fees only
        UPDATE perpetual_positions
            SET "totalRealizedPnl" =
                COALESCE("totalRealizedPnl", 0) - COALESCE(fill_fee, 0)
            WHERE "id" = position_id;
        RETURN;
    END IF;

    closing_amount := LEAST(fill_size, pos_size_before); -- cap to existing

    IF pos_side_before = 'LONG' THEN
        pnl := (fill_price - entry_price_before) * closing_amount;
    ELSE
        pnl := (entry_price_before - fill_price) * closing_amount;
    END IF;

    UPDATE perpetual_positions
        SET "totalRealizedPnl" =
            COALESCE("totalRealizedPnl", 0) + COALESCE(pnl, 0) - COALESCE(fill_fee, 0)
        WHERE "id" = position_id;
END;
$$ LANGUAGE plpgsql;
