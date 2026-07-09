local main_key = KEYS[1]

local tokens_left_key = "tokens_left"

local last_updated_key = "last_updated"

local bucket_size = tonumber(ARGV[1])

local refill_rate_per_second = tonumber(ARGV[2])

local current_time_arr = redis.call("TIME")

local current_time = tonumber(current_time_arr[1]) + tonumber(current_time_arr[2]) / 1000000

local result = redis.call("HMGET", main_key, tokens_left_key, last_updated_key)

local tokens_left

local last_updated

local allowed = 0
-- 0 -> dont allow
-- 1 -> allow

if result[1] == false and result[2] == false then
    tokens_left = bucket_size - 1
    last_updated = current_time
    redis.call("HSET", main_key, tokens_left_key, tostring(tokens_left), last_updated_key, tostring(last_updated))
    allowed = 1

    tokens_left = tostring(tokens_left)

    return {allowed, tokens_left}
else
    tokens_left = tonumber(result[1])
    last_updated = tonumber(result[2])
end

local elapsed_time = current_time - last_updated

local tokens_to_refill = elapsed_time * refill_rate_per_second

tokens_left = math.min(bucket_size, tokens_left + tokens_to_refill)

if tokens_left >= 1 then
    allowed = 1
    tokens_left = tokens_left - 1
end

tokens_left = tostring(tokens_left)

redis.call("HSET", main_key, tokens_left_key, tokens_left, last_updated_key, current_time)

return { allowed, tokens_left }







