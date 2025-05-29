DROP TABLE IF EXISTS subaccount_perp_positions;

-- This aggregate.sql file performs the second suite of tests.
-- It calculates funding payments for subaccounts between heights $1 and $2
-- where $1 is the start height and $2 is the end height
CREATE TABLE
    subaccount_funding_payments AS
WITH
    net AS (
        SELECT
            "subaccount_id",
            "clob_pair_id", -- align the names
            SUM(
                CASE
                    WHEN side = 'BUY' THEN size
                    WHEN side = 'SELL' THEN - size
                END
            ) AS net_size
        FROM
            fills
        -- contains inclusively fills from [$1 + 1, $2]
        WHERE created_at_height > $1 AND created_at_height <= $2
        GROUP BY
            "subaccount_id",
            "clob_pair_id"
    ),
    -- figure out what the last funding payment was.
    last_funding_payment AS (
        SELECT DISTINCT ON (subaccount_id, perpetual_id)
            subaccount_id,
            perpetual_id,
            size as last_snapshot_size,
            created_at_height
        FROM funding_payments
        -- snapshot at height $1.
        WHERE created_at_height = $1
        ORDER BY subaccount_id, perpetual_id, created_at_height DESC
    ),
    paired AS (
        SELECT
            n."subaccount_id",
            pm.id AS perpetual_id,
            pm.market_id,
            n.clob_pair_id,
            pm.ticker,
            COALESCE(n.net_size, 0) + COALESCE(lfp.last_snapshot_size, 0) AS net_size
        FROM
            net n
            JOIN perpetual_markets pm ON pm.clob_pair_id = n.clob_pair_id
            -- okay, but what if the clob_pair_id is not in the perpetual_markets table?
            -- how do we handle a clob_pair_id that we can't find a perpetual_id for?
            LEFT JOIN last_funding_payment lfp ON lfp.subaccount_id = n.subaccount_id 
                AND lfp.perpetual_id = pm.id
    ),
    funding AS (
        /* Grab the latest funding index update per perpetual_id */
        SELECT DISTINCT
            ON (f."perpetualId") f."perpetualId" AS perpetual_id,
            f.rate,
            f."oraclePrice" AS oracle_price,
            f."effectiveAt" AS effective_at,
            f."effectiveAtHeight" AS effective_at_height
        FROM
            funding_index_updates f
        ORDER BY
            f."perpetualId",
            f."effectiveAtHeight" DESC
    )
SELECT
    p."subaccount_id",
    p.perpetual_id,
    p.ticker,
    p.net_size AS size,
    CASE
        WHEN p.net_size > 0 THEN 'BUY'
        WHEN p.net_size < 0 THEN 'SELL'
        ELSE 'FLAT'
    END AS side,
    f.oracle_price,
    f.rate,
    - p.net_size * f.oracle_price * f.rate AS payment,
    f.effective_at,
    f.effective_at_height
FROM
    paired p
    LEFT JOIN funding f ON f.perpetual_id = p.perpetual_id
ORDER BY
    p."subaccount_id",
    p.perpetual_id;
