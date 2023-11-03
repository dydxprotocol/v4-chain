CREATE OR REPLACE FUNCTION dydx_clob_pair_status_to_market_status(status jsonb)
    RETURNS text AS $$
BEGIN
    CASE status
        WHEN '1'::jsonb THEN RETURN 'ACTIVE'; /** CLOB_PAIR_STATUS_ACTIVE */
        WHEN '2'::jsonb THEN RETURN 'PAUSED'; /** CLOB_PAIR_STATUS_PAUSED */
        WHEN '3'::jsonb THEN RETURN 'CANCEL_ONLY'; /** CLOB_PAIR_STATUS_CANCEL_ONLY */
        WHEN '4'::jsonb THEN RETURN 'POST_ONLY'; /** CLOB_PAIR_STATUS_POST_ONLY */
        WHEN '5'::jsonb THEN RETURN 'INITIALIZING'; /** CLOB_PAIR_STATUS_INITIALIZING */
        ELSE RAISE EXCEPTION 'Invalid clob pair status: %', status;
    END CASE;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;