-- The `KEYS` table contains the market cache keys
-- The `ARGV` table contains the prices for each market and a single timestamp value at the end.
local numKeys = #KEYS
local numArgs = #ARGV
local timestamp = tonumber(ARGV[numArgs])
local fiveSeconds = 5

if not timestamp then
  return redis.error_reply("Invalid timestamp")
end

local cutoffTime = timestamp - fiveSeconds

for i = 1, numKeys do
  local priceCacheKey = KEYS[i]
  local price = tonumber(ARGV[i])

  if not price then
    return redis.error_reply("Invalid price for key " .. priceCacheKey)
  end

  redis.call("zadd", priceCacheKey, timestamp, price)
  redis.call("zremrangebyscore", priceCacheKey, "-inf", cutoffTime)
end

return true