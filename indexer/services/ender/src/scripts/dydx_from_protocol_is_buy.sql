/**
 Converts a boolean 'isBuy' field to 'BUY' or 'SELL'.
 */
CREATE OR REPLACE FUNCTION dydx_from_protocol_is_buy(is_buy boolean) RETURNS text AS $$
BEGIN
    IF is_buy THEN
        RETURN 'BUY';
    ELSE
        RETURN 'SELL';
    END IF;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
