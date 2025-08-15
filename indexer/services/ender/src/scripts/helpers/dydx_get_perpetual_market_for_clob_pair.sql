CREATE OR REPLACE FUNCTION dydx_get_perpetual_market_for_clob_pair(
    clob_pair_id bigint
) RETURNS perpetual_markets AS $$
/**
  Returns the perpetual market record for the provided clob pair.

  Parameters:
    - clob_pair_id: The clob pair id.
  Returns: the only perpetual market for the clob pair. Throws an exception if not exactly one row is found.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    perpetual_market_record perpetual_markets%ROWTYPE;
BEGIN
    SELECT * INTO STRICT perpetual_market_record FROM perpetual_markets WHERE "clobPairId" = clob_pair_id;
    RETURN perpetual_market_record;
EXCEPTION
    WHEN NO_DATA_FOUND THEN
        RAISE EXCEPTION 'Unable to find perpetual market with clobPairId: %', clob_pair_id;
    WHEN TOO_MANY_ROWS THEN
        /** This should never happen and if it ever were to would indicate that the table has malformed data. */
        RAISE EXCEPTION 'Found multiple perpetual markets with clobPairId: %', clob_pair_id;
END;
$$ LANGUAGE plpgsql;