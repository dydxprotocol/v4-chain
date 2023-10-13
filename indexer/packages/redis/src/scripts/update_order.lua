-- TODO(CORE-512): add info/resources around caches.
-- Doc: https://www.notion.so/dydx/Indexer-Technical-Spec-a6b15644502048f994c98dee35b96e96#61d5f8ca5117476caab78b3f0691b1d0

-- Key for the value in the orders cache, see `src/caches/orders-cache.ts` for more details
local orderKey = KEYS[1];
-- Key for the value in the orders data cache, see `src/caches/orders-data-cache.ts` for more
-- details
local orderDataKey = KEYS[2];

-- Updated total filled in quantums for the order
local newTotalFilled = ARGV[1];

-- This script returns the following values in an array
-- 1. Was an order updated - 1 if an order is updated, 0 if not
-- 2. Old total filled amount in quantums of the updated order - "0" if no order was updated
-- 3. If the old order was resting on the book - "true/false", "false" if no order was updated
-- 4. Encoded order for the updated order - "" if an order was not updated


-- order data has the format:
-- [good-til-block or sequence number of order]_[total filled]_[true/false, if order is on the book]

local order = redis.call("get", orderKey);
if not order then
  return {0, "0", "false", ""}
else
  local oldOrderData = redis.call("get", orderDataKey);
  -- refer to above comment on order data format
  local j = string.find(oldOrderData, "_");
  local orderExpiry = string.sub(oldOrderData, 0, j - 1);
  local i = string.find(oldOrderData, "_", j + 1);
  local oldFilled = string.sub(oldOrderData, j + 1, i - 1);
  local oldRestingOnBook = string.sub(oldOrderData, i + 1);

  -- if an order has an update to it's total filled quantums, it has sucessfully matched and placed
  -- on the book, so update the order data to indicate the order is on the book
  redis.call("set", orderDataKey, orderExpiry .. "_" .. newTotalFilled .. "_true");

  return {1, oldFilled,  oldRestingOnBook, order};
end
