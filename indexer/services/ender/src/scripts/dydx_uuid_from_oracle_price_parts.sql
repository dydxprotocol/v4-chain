/**
  Returns a UUID using the parts of an OraclePrice (https://github.com/dydxprotocol/v4-chain/blob/755b0b928be793072d19eb3a1608e7a2503f396a/indexer/packages/postgres/src/stores/oracle-price-table.ts#L24).
*/
CREATE OR REPLACE FUNCTION dydx_uuid_from_oracle_price_parts(market_id int, block_height int) RETURNS uuid AS $$
BEGIN
    return dydx_uuid(concat(market_id, '-', block_height));
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;