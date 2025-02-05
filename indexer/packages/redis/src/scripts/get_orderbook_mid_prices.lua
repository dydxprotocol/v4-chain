-- KEYS is an array of cache keys for a market

local results = {}
for i, key in ipairs(KEYS) do
  local price = redis.call("GET", key)
  results[i] = price or false
end

return results
