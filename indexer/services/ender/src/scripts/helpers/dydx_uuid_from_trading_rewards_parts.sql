CREATE OR REPLACE FUNCTION dydx_uuid_from_trading_rewards_parts(address text, block_height int) RETURNS uuid AS $$
/**
  Returns a UUID using the parts of a trading rewards.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    return dydx_uuid(concat(address, '-', block_height));
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;

