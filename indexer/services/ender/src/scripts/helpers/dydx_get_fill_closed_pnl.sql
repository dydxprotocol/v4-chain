/**
  Returns the realized (closed) PnL of a single fill: the price-PnL of the
  portion of the fill that CLOSES existing position, measured against the
  position's average entry price at fill time. Excludes fees (the fills table
  carries the fee separately) and funding.

  Returns NULL for fills that open or increase a position — they realize
  nothing — so consumers can distinguish "no realization" from "realized 0".

  This is the per-fill decomposition of the position-level aggregate
  maintained by dydx_apply_fill_realized_effects (which also nets fees in);
  the two must stay consistent, so that function delegates to this one.
*/
CREATE OR REPLACE FUNCTION dydx_get_fill_closed_pnl(
    fill_side text,
    fill_size numeric,
    fill_price numeric,
    pos_side_before text,
    pos_size_before numeric,
    entry_price_before numeric
) RETURNS numeric AS $$
DECLARE
    closing_amount numeric;
BEGIN
    IF pos_size_before IS NULL OR pos_size_before = 0
        OR entry_price_before IS NULL OR pos_side_before IS NULL THEN
        RETURN NULL; -- no position to close against
    END IF;

    IF NOT ((pos_side_before = 'LONG'  AND fill_side = 'SELL')
         OR (pos_side_before = 'SHORT' AND fill_side = 'BUY')) THEN
        RETURN NULL; -- increasing fill: nothing realizes
    END IF;

    /* A fill larger than the position flips it; only the closing portion
       realizes PnL (the remainder opens the new position at fill_price). */
    closing_amount := LEAST(fill_size, pos_size_before);

    IF pos_side_before = 'LONG' THEN
        RETURN dydx_trim_scale((fill_price - entry_price_before) * closing_amount);
    ELSE
        RETURN dydx_trim_scale((entry_price_before - fill_price) * closing_amount);
    END IF;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
