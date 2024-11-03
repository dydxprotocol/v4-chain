CREATE OR REPLACE FUNCTION klyra_uuid_from_oracle_price_parts(market_id int, block_height int) RETURNS uuid AS $$
/**
  Returns a UUID using the parts of an OraclePrice.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    return klyra_uuid(concat(market_id, '-', block_height));
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;