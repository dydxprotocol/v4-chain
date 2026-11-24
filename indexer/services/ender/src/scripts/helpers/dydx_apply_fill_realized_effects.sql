/**
  Folds one fill's realized effects into the position-level aggregate
  perpetual_positions."totalRealizedPnl": the closing price-PnL (per
  dydx_get_fill_closed_pnl, the same quantity persisted per fill as
  fills."closedPnl") plus fees, which are charged on every fill — including
  opening/increasing fills, where no price-PnL realizes.
*/
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
    pnl numeric;
BEGIN
    pnl := COALESCE(dydx_get_fill_closed_pnl(
        fill_side, fill_size, fill_price,
        pos_side_before, pos_size_before, entry_price_before), 0);

    UPDATE perpetual_positions
        SET "totalRealizedPnl" =
            COALESCE("totalRealizedPnl", 0) + pnl - COALESCE(fill_fee, 0)
        WHERE "id" = position_id;
END;
$$ LANGUAGE plpgsql;
