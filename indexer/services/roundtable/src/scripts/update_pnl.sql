INSERT INTO pnl (
    "subaccountId",
    "createdAt", 
    "createdAtHeight",
    "equity",
    "netTransfers",
    "totalPnl"
)
WITH previous_pnl AS (
    SELECT 
        "subaccountId",
        "totalPnl" as prev_total_pnl,
        "netTransfers" as prev_net_transfers
    FROM pnl
    WHERE "createdAtHeight" = :start
),
transfer_aggregated AS (
    SELECT 
        "subaccountId",
        SUM(transfer_amount) as transfer_delta
    FROM (
        SELECT "senderSubaccountId" as "subaccountId", -"size" as transfer_amount
        FROM transfers
        WHERE "createdAtHeight" > :start AND "createdAtHeight" <= :end
        UNION ALL
        SELECT "recipientSubaccountId" as "subaccountId", "size" as transfer_amount
        FROM transfers  
        WHERE "createdAtHeight" > :start AND "createdAtHeight" <= :end
    ) transfer_data
    GROUP BY "subaccountId"
),
all_relevant_subaccounts AS (
    SELECT "subaccountId" as "id" FROM previous_pnl
    UNION
    SELECT "subaccountId" as "id" FROM transfer_aggregated
),
end_time AS (
   SELECT MAX("createdAt") AS timestamp_at_end
   FROM funding_payments
   WHERE "createdAtHeight" = :end
),
funding_data AS (
    SELECT 
        "subaccountId",
        "createdAtHeight",
        "payment",
        "size" * "oraclePrice" as position_value
    FROM funding_payments
    WHERE "createdAtHeight" IN (:start, :end)
       OR ("createdAtHeight" > :start AND "createdAtHeight" <= :end)
),
funding_aggregated AS (
    SELECT 
        "subaccountId",
        SUM(CASE 
            WHEN "createdAtHeight" > :start AND "createdAtHeight" <= :end 
            THEN "payment" ELSE 0 
        END) as total_funding_payments,
        SUM(CASE WHEN "createdAtHeight" = :start THEN position_value ELSE 0 END) as position_value_start,
        SUM(CASE WHEN "createdAtHeight" = :end THEN position_value ELSE 0 END) as position_value_end
    FROM funding_data
    GROUP BY "subaccountId"
),
trade_cash_flows AS (
    SELECT 
        "subaccountId",
        SUM(CASE 
            WHEN "side" = 'SELL' THEN "quoteAmount"
            WHEN "side" = 'BUY' THEN -"quoteAmount"
        END) as net_cash_flow
    FROM fills
    WHERE "createdAtHeight" > :start 
      AND "createdAtHeight" <= :end
    GROUP BY "subaccountId"
)
SELECT 
    s."id" as "subaccountId",
    et.timestamp_at_end as "createdAt",
    :end as "createdAtHeight",
    -- Calculate equity = totalPnl + netTransfers
    COALESCE(pp.prev_total_pnl, 0) + 
    COALESCE(fa.total_funding_payments, 0) + 
    (COALESCE(fa.position_value_end, 0) - COALESCE(fa.position_value_start, 0)) +
    COALESCE(tcf.net_cash_flow, 0) +
    COALESCE(pp.prev_net_transfers, 0) + 
    COALESCE(ta.transfer_delta, 0) as "equity",
    -- Calculate netTransfers = previous + new transfers
    COALESCE(pp.prev_net_transfers, 0) + COALESCE(ta.transfer_delta, 0) as "netTransfers",
    -- Calculate totalPnl = previous + funding + position effects  
    COALESCE(pp.prev_total_pnl, 0) + 
    COALESCE(fa.total_funding_payments, 0) + 
    (COALESCE(fa.position_value_end, 0) - COALESCE(fa.position_value_start, 0)) +
    COALESCE(tcf.net_cash_flow, 0) as "totalPnl"
FROM all_relevant_subaccounts s
CROSS JOIN end_time et
LEFT JOIN previous_pnl pp ON s."id" = pp."subaccountId"
LEFT JOIN funding_aggregated fa ON s."id" = fa."subaccountId"
LEFT JOIN trade_cash_flows tcf ON s."id" = tcf."subaccountId"
LEFT JOIN transfer_aggregated ta ON s."id" = ta."subaccountId"
ORDER BY s."id";