CREATE OR REPLACE FUNCTION klyra_clob_pair_status_to_market_status(status jsonb)
    RETURNS text AS $$
/**
  Returns the market status
  from the clob pair status.

  Parameters:
    - status: the ClobPairStatus 

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    CASE status
        WHEN '1'::jsonb THEN RETURN 'ACTIVE'; /** CLOB_PAIR_STATUS_ACTIVE */
        WHEN '2'::jsonb THEN RETURN 'PAUSED'; /** CLOB_PAIR_STATUS_PAUSED */
        WHEN '3'::jsonb THEN RETURN 'CANCEL_ONLY'; /** CLOB_PAIR_STATUS_CANCEL_ONLY */
        WHEN '4'::jsonb THEN RETURN 'POST_ONLY'; /** CLOB_PAIR_STATUS_POST_ONLY */
        WHEN '5'::jsonb THEN RETURN 'INITIALIZING'; /** CLOB_PAIR_STATUS_INITIALIZING */
        WHEN '6'::jsonb THEN RETURN 'FINAL_SETTLEMENT'; /** CLOB_PAIR_STATUS_FINAL_SETTLEMENT */
        ELSE RAISE EXCEPTION 'Invalid clob pair status: %', status;
    END CASE;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;