INSERT INTO pnl (
    "subaccountId", 
    "createdAt", 
    "createdAtHeight", 
    "deltaFundingPayments", 
    "deltaPositionEffects", 
    "totalPnl"
)
WITH subaccounts_with_transfers AS (
    SELECT "id" FROM subaccounts
    WHERE "id" IN (
        SELECT "senderSubaccountId" FROM transfers
        WHERE "createdAtHeight" <= :end
        UNION
        SELECT "recipientSubaccountId" FROM transfers
        WHERE "createdAtHeight" <= :end
    )
),
end_timestamp AS (
    -- Get timestamp from oracle prices at the latest height <= end
    SELECT "effectiveAt"
    FROM oracle_prices
    WHERE "effectiveAtHeight" = (
        SELECT MAX("effectiveAtHeight")
        FROM oracle_prices
        WHERE "effectiveAtHeight" <= :end
    )
    LIMIT 1
),
latest_oracle_prices AS (
    SELECT 
        "marketId",
        'start' as price_type,
        "price",
        "effectiveAtHeight"
    FROM oracle_prices op1
    WHERE "effectiveAtHeight" = (
        SELECT MAX("effectiveAtHeight")
        FROM oracle_prices op2
        WHERE op2."marketId" = op1."marketId"
          AND "effectiveAtHeight" <= :start
    )
    UNION ALL
    SELECT 
        "marketId",
        'end' as price_type,
        "price",
        "effectiveAtHeight"
    FROM oracle_prices op1
    WHERE "effectiveAtHeight" = (
        SELECT MAX("effectiveAtHeight")
        FROM oracle_prices op2
        WHERE op2."marketId" = op1."marketId"
          AND "effectiveAtHeight" <= :end
    )
),
open_position_pnl AS (
    SELECT 
        pp."subaccountId",
        SUM((op_end."price" - 
             CASE 
                 WHEN pp."createdAtHeight" <= :start THEN op_start."price"
                 ELSE pp."entryPrice"
             END
            ) * pp."size"
        ) as open_pnl
    FROM perpetual_positions pp
    JOIN perpetual_markets pm ON pp."perpetualId" = pm."id"
    JOIN latest_oracle_prices op_end ON pm."marketId" = op_end."marketId" 
        AND op_end.price_type = 'end'
    LEFT JOIN latest_oracle_prices op_start ON pm."marketId" = op_start."marketId" 
        AND op_start.price_type = 'start'
        AND pp."createdAtHeight" <= :start
    WHERE pp."status" = 'OPEN'
      AND pp."createdAtHeight" <= :end
    GROUP BY pp."subaccountId"
),
closed_position_pnl AS (
    SELECT 
        pp."subaccountId",
        SUM((pp."exitPrice" - 
             CASE 
                 WHEN pp."createdAtHeight" <= :start THEN op_start."price"
                 ELSE pp."entryPrice"
             END
            ) * pp."size"
        ) as closed_pnl
    FROM perpetual_positions pp
    JOIN perpetual_markets pm ON pp."perpetualId" = pm."id"
    LEFT JOIN latest_oracle_prices op_start ON pm."marketId" = op_start."marketId" 
        AND op_start.price_type = 'start'
        AND pp."createdAtHeight" <= :start
    WHERE pp."status" IN ('CLOSED', 'LIQUIDATED')
      AND pp."closedAtHeight" > :start 
      AND pp."closedAtHeight" <= :end
    GROUP BY pp."subaccountId"
),
funding_payments_sum AS (
    SELECT 
        "subaccountId",
        SUM("payment") as total_funding_payments
    FROM funding_payments
    WHERE "createdAtHeight" > :start 
      AND "createdAtHeight" <= :end
    GROUP BY "subaccountId"
)
SELECT 
    s."id" as "subaccountId",
    et."effectiveAt" as "createdAt",
    :end as "createdAtHeight",
    COALESCE(fp.total_funding_payments, 0) as "deltaFundingPayments",
    COALESCE(open_pnl.open_pnl, 0) + COALESCE(closed_pnl.closed_pnl, 0) as "deltaPositionEffects",
    COALESCE(p."totalPnl", 0) + 
    COALESCE(fp.total_funding_payments, 0) + 
    COALESCE(open_pnl.open_pnl, 0) + 
    COALESCE(closed_pnl.closed_pnl, 0) as "totalPnl"
FROM subaccounts_with_transfers s
CROSS JOIN end_timestamp et
LEFT JOIN pnl p ON s."id" = p."subaccountId" 
    AND p."createdAtHeight" = :start
LEFT JOIN funding_payments_sum fp ON s."id" = fp."subaccountId"
LEFT JOIN open_position_pnl open_pnl ON s."id" = open_pnl."subaccountId"
LEFT JOIN closed_position_pnl closed_pnl ON s."id" = closed_pnl."subaccountId"
ORDER BY s."id";