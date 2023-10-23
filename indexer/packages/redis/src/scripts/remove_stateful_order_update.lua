-- Hash of the ZSET tracking when a stateful order update was added to the cache
local statefulOrderUpdateIdHash = KEYS[1]
-- Hash of the HSET tracking the stateful order updates
local statefulOrderUpdateHash = KEYS[2]

-- Order id of the stateful order update to add to the cache
local statefulOrderId = ARGV[1]
-- Maximum timestamp of the stateful order update to remove
-- Needed for removing old stateful order updates
local maxTimestamp = ARGV[2]

-- This script attempts to remove a stateful order update associated with an order id
-- from the ZSET / HSET tracking stateful order updates
-- This script returns the either an empty string if no stateful order update was removed
-- or the encoded stateful order update protobuf that was removed

local oldTimestamp = redis.call("ZSCORE", statefulOrderUpdateIdHash, statefulOrderId)
-- if there is no timestamp, return an empty string
if not oldTimestamp then
  return ""
end

-- If the timestamp of the stateful order update is less than the maximum tiemstamp
-- to remove, return an empty string and don't remove the stateful order update
if tonumber(oldTimestamp) > tonumber(maxTimestamp) then
  return ""
end

-- The timestamp exists and is less than the maximum timestamp, and so delete the order id
-- from the ZSET
redis.call("ZREM", statefulOrderUpdateIdHash, statefulOrderId)

local oldStatefulOrderUpdate = redis.call("HGET", statefulOrderUpdateHash, statefulOrderId)
-- If there's no order update, return empty string
if not oldStatefulOrderUpdate then
  return ""
end

redis.call("HDEL", statefulOrderUpdateHash, statefulOrderId)

return oldStatefulOrderUpdate
