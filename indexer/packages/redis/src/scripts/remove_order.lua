-- TODO(CORE-512): add info/resources around caches.
-- Doc: https://www.notion.so/dydx/Indexer-Technical-Spec-a6b15644502048f994c98dee35b96e96#61d5f8ca5117476caab78b3f0691b1d0

-- Key for the value in the orders cache, see `src/caches/orders-cache.ts` for more details
local orderKey = KEYS[1];
-- Key for the value in the orders data cache, see `src/caches/orders-data-cache.ts` for more
-- details
local orderDataKey = KEYS[2];
-- Key for the value in the subaccount order ids cache, see
-- `src/caches/subaccount-order-ids-cache.ts` for more details
local subaccountKey = KEYS[3];
-- Key for the value in the orders data cache, see `src/caches/orders-expiry-cache.ts` for more
-- details
local orderExpiryKey = KEYS[4];

-- UUID of the order being removed
local orderId = ARGV[1];

-- This script returns the following values in an array
-- 1. Was an order removed - 1 if an order is removed, 0 if not
-- 2. Total filled size of the removed order - "0" if an order was not removed
-- 3. Was the removed order resting on the book - "true/false", "false" if the order was not resting
--    on the book, or if no order was removed, "true" otherwise
-- 4. Encoded removed order - "" if an order was not removed

-- order data has the format:
-- [good-til-block or sequence number of order]_[total filled]_[true/false, if order is on the book]

local removedOrder = redis.call("get", orderKey);
if not removedOrder then
  return {0, "0", "false", ""};
else
  local removedOrderData = redis.call("get", orderDataKey);
  -- refer to above comment on order data format
  -- ignore the expiry (good-til-block or sequence number of the order), the protocol ignores
  -- invalid order removals including those with a expiry less than an existing orders so there's
  -- no need to check the expiry in the indexer
  local j = string.find(removedOrderData, "_");
  local i = string.find(removedOrderData, "_", j + 1);
  local removedFilled = string.sub(removedOrderData, j + 1, i - 1);
  local removedRestingOnBook = string.sub(removedOrderData, i + 1);

  -- remove the
  -- * order
  -- * order data
  -- * order id from the subaccount order id cache
  -- * order id from the order expiry cache
  redis.call("del", orderKey);
  redis.call("del", orderDataKey);
  redis.call("hdel", subaccountKey, orderId);
  redis.call("zrem", orderExpiryKey, orderId);

  return {1, removedFilled, removedRestingOnBook, removedOrder};
end
