local canceledOrdersCacheKey = KEYS[1]
local canceledOrderWindowSize = tonumber(KEYS[2])
local orderId = ARGV[1]
local timestamp = tonumber(ARGV[2])

-- Remove canceled orders that expired
redis.call("ZREMRANGEBYSCORE", canceledOrdersCacheKey, "-inf", "(" .. timestamp)

-- Add the new canceled order with its expiration time
local expirationTime = timestamp + canceledOrderWindowSize
local numCancelledOrders = redis.call("ZADD", canceledOrdersCacheKey, expirationTime, orderId)
return numCancelledOrders
