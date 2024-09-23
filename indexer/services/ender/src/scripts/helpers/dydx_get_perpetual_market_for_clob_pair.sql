DROP TYPE IF EXISTS perpetual_market_filtered CASCADE;

CREATE TYPE perpetual_market_filtered AS (
    id bigint,
    "clobPairId" bigint,
    ticker character varying(255),
    "marketId" integer,
    status text,
    "quantumConversionExponent" integer,
    "atomicResolution" integer,
    "subticksPerTick" integer,
    "stepBaseQuantums" integer,
    "liquidityTierId" integer,
    "marketType" text
);

CREATE OR REPLACE FUNCTION dydx_to_jsonb(row_t perpetual_market_filtered) RETURNS jsonb AS $$
BEGIN
    RETURN jsonb_build_object(
        'id', row_t.id::text,
        'clobPairId', row_t."clobPairId"::text,
        'ticker', row_t.ticker,
        'marketId', row_t."marketId",
        'status', row_t.status,
        'quantumConversionExponent', row_t."quantumConversionExponent",
        'atomicResolution', row_t."atomicResolution",
        'subticksPerTick', row_t."subticksPerTick",
        'stepBaseQuantums', row_t."stepBaseQuantums",
        'liquidityTierId', row_t."liquidityTierId",
        'marketType', row_t."marketType"
    );
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;

DROP FUNCTION IF EXISTS dydx_get_perpetual_market_for_clob_pair(bigint);

CREATE OR REPLACE FUNCTION dydx_get_perpetual_market_for_clob_pair(
    clob_pair_id bigint
) RETURNS perpetual_market_filtered AS $$
/**
  Returns the perpetual market record with selected fields for the provided clob pair.

  Parameters:
    - clob_pair_id: The clob pair id.
  Returns: the filtered perpetual market fields for the clob pair.
  Throws an exception if not exactly one row is found.
*/
DECLARE
    perpetual_market_record perpetual_market_filtered;
BEGIN
    SELECT
        id,
        "clobPairId",
        ticker,
        "marketId",
        status,
        "quantumConversionExponent",
        "atomicResolution",
        "subticksPerTick",
        "stepBaseQuantums",
        "liquidityTierId",
        "marketType"
    INTO STRICT perpetual_market_record
    FROM perpetual_markets
    WHERE "clobPairId" = clob_pair_id;

    RETURN perpetual_market_record;
EXCEPTION
    WHEN NO_DATA_FOUND THEN
        RAISE EXCEPTION 'Unable to find perpetual market with clobPairId: %', clob_pair_id;
    WHEN TOO_MANY_ROWS THEN
        RAISE EXCEPTION 'Found multiple perpetual markets with clobPairId: %', clob_pair_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION dydx_get_perpetual_market_for_id(
    perpetual_id bigint
) RETURNS perpetual_market_filtered AS $$
/**
  Returns the perpetual market record with selected fields for the provided perpetual id.

  Parameters:
    - perpetual_id: The perpetual market id.
  Returns: the filtered perpetual market fields for the provided perpetual id.
  Throws an exception if not exactly one row is found.
*/
DECLARE
    perpetual_market_record perpetual_market_filtered;
BEGIN
    SELECT
        id,
        "clobPairId",
        ticker,
        "marketId",
        status,
        "quantumConversionExponent",
        "atomicResolution",
        "subticksPerTick",
        "stepBaseQuantums",
        "liquidityTierId",
        "marketType"
    INTO STRICT perpetual_market_record
    FROM perpetual_markets
    WHERE "id" = perpetual_id;

    RETURN perpetual_market_record;
EXCEPTION
    WHEN NO_DATA_FOUND THEN
        RAISE EXCEPTION 'Unable to find perpetual market with id: %', perpetual_id;
    WHEN TOO_MANY_ROWS THEN
        RAISE EXCEPTION 'Found multiple perpetual markets with id: %', perpetual_id;
END;
$$ LANGUAGE plpgsql;
