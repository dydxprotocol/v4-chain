-- KEYS is an array of cache keys for a market

local results = {}
for i, key in ipairs(KEYS) do
  -- Get the prices for each key, but limit to a maximum of 10
  local prices = redis.call("ZRANGE", key, 0, 9)
  results[i] = prices
end

return results
