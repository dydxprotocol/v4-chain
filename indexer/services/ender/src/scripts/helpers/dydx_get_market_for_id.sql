CREATE OR REPLACE FUNCTION dydx_get_market_for_id(
    market_id bigint
) RETURNS markets AS $$
/**
  Returns the market record for the provided market ID.

  Parameters:
    - market_id: The market id.
  Returns: the only market for the given ID. Throws an exception if not exactly one row is found.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    market_record markets%ROWTYPE;
BEGIN
    SELECT * INTO STRICT market_record FROM markets WHERE id = market_id;
    RETURN market_record;
EXCEPTION
    WHEN NO_DATA_FOUND THEN
        RAISE EXCEPTION 'Unable to find market with id: %', market_id;
    WHEN TOO_MANY_ROWS THEN
        /** This should never happen and if it ever were to would indicate that the table has malformed data. */
        RAISE EXCEPTION 'Found multiple markets with id: %', market_id;
END;
$$ LANGUAGE plpgsql;
