-- It calculates funding payments for subaccounts between heights $1 and $2
-- where $1 is the start height and $2 is the end height.
INSERT INTO funding_payments (
    "subaccountId",
    "createdAt",
    "createdAtHeight",
    "perpetualId",
    ticker,
    "oraclePrice",
    size,
    side,
    rate,
    payment
)
WITH
    net AS (
        SELECT
            "subaccountId",
            "clobPairId", -- align the names
            SUM(
                CASE
                    WHEN side = 'BUY' THEN size
                    WHEN side = 'SELL' THEN - size
                END
            ) AS net_size
        FROM
            fills
        -- contains inclusively fills from [$1 + 1, $2]
        WHERE "createdAtHeight" > $1 AND "createdAtHeight" <= $2
        GROUP BY
            "subaccountId",
            "clobPairId"
    ),
    -- figure out what the last funding payment was.
    last_funding_payment AS (
        SELECT DISTINCT ON ("subaccountId", "perpetualId")
            "subaccountId",
            "perpetualId",
            size as last_snapshot_size,
            "createdAtHeight"
        FROM funding_payments
        -- snapshot at height $1.
        WHERE "createdAtHeight" = $1
        ORDER BY "subaccountId", "perpetualId", "createdAtHeight" DESC
    ),
    paired AS (
        SELECT
            n."subaccountId",
            pm.id AS "perpetualId",
            pm."marketId",
            n."clobPairId",
            pm.ticker,
            COALESCE(n.net_size, 0) + COALESCE(lfp.last_snapshot_size, 0) AS net_size
        FROM
            net n
            JOIN perpetual_markets pm ON pm."clobPairId" = n."clobPairId"
            -- okay, but what if the clob_pair_id is not in the perpetual_markets table?
            -- how do we handle a clob_pair_id that we can't find a perpetual_id for?
            LEFT JOIN last_funding_payment lfp ON lfp."subaccountId" = n."subaccountId" 
                AND lfp."perpetualId" = pm.id
    ),
    funding AS (
        /* Grab the latest funding index update per perpetual_id */
        SELECT DISTINCT
            ON (f."perpetualId") f."perpetualId" AS "perpetualId",
            f.rate,
            f."oraclePrice" AS "oraclePrice",
            f."effectiveAt" AS "effectiveAt"
        FROM
            funding_index_updates f
        ORDER BY
            f."perpetualId",
            f."effectiveAtHeight" DESC
    )
SELECT
    p."subaccountId",
    CURRENT_TIMESTAMP as "createdAt",
    $2 as "createdAtHeight",
    p."perpetualId",
    p.ticker,
    f."oraclePrice",
    p.net_size AS size,
    CASE
        WHEN p.net_size > 0 THEN 'LONG'
        ELSE 'SHORT'
    END AS side,
    f.rate,
    - p.net_size * f."oraclePrice" * f.rate AS payment
FROM
    paired p
    LEFT JOIN funding f ON f."perpetualId" = p."perpetualId"
WHERE
    p.net_size != 0
ORDER BY
    p."subaccountId",
    p."perpetualId";
