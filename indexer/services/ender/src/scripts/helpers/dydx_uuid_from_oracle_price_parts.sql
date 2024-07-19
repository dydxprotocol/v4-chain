CREATE OR REPLACE FUNCTION dydx_uuid_from_oracle_price_parts(market_id int, block_height int) RETURNS uuid AS $$
/**
  Returns a UUID using the parts of an OraclePrice (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/indexer/packages/postgres/src/stores/oracle-price-table.ts#L24).

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    return dydx_uuid(concat(market_id, '-', block_height));
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;