-- Hash of the ZSET tracking when a stateful order update was added to the cache
local statefulOrderUpdateIdHash = KEYS[1]
-- Hash of the HSET tracking the stateful order updates
local statefulOrderUpdateHash = KEYS[2]

-- Order id of the stateful order update to add to the cache
local statefulOrderId = ARGV[1]
-- Encoded stateful order update protobuf
local statefulOrderUpdate = ARGV[2]
-- Timestamp of when the order update was added
local timestamp = ARGV[3]

redis.call("ZADD", statefulOrderUpdateIdHash, timestamp, statefulOrderId)
redis.call("HSET", statefulOrderUpdateHash, statefulOrderId, statefulOrderUpdate)

return 1
