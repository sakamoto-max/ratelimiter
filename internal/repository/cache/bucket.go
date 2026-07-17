package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/sakamoto-max/ratelimiter/internal/domain"
	myErr "github.com/sakamoto-max/ratelimiter/internal/pkg/myerrors"

	"github.com/redis/go-redis/v9"
)

type BucketIface interface {
	CheckLimit(ctx context.Context, userIp string, policy domain.Policy) (*domain.UserBucket, error)
}

type Bucket struct {
	client *redis.Client
}

func (c *Bucket) CheckLimit(ctx context.Context, userIp string, policy domain.Policy) (*domain.UserBucket, error) {

	timeStart := time.Now()

	script := redis.NewScript(`
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
	`)

	mainKey := fmt.Sprintf("user_ip:%v:resource_name:%v:bucket", userIp, policy.ResourceName)

	keys := []string{mainKey}

	cmd, err := script.Eval(ctx, c.client, keys, policy.BucketSize, policy.RefillPerSecond).Result()
	if err != nil {
		return nil, myErr.WrapErr(fmt.Errorf("failed to implement the redis script : %w", err), myErr.InternalServerErr)
	}

	totalTime := time.Since(timeStart).Microseconds()

	data, ok := cmd.([]any)

	if !ok {
		return nil, myErr.WrapErr(fmt.Errorf("failed to parse the data %v from redis", data), myErr.InternalServerErr)
	}

	var allowed int64
	var tokensLeft float64

	allowed, ok = data[0].(int64)
	if !ok {
		return nil, myErr.WrapErr(fmt.Errorf("failed to parse the data to get allowed from redis"), myErr.InternalServerErr)
	}

	tokensLeftStr, ok := data[1].(string)
	if !ok {
		return nil, myErr.WrapErr(fmt.Errorf("failed to parse the data to get tokens left from redis"), myErr.InternalServerErr)
	}

	tokensLeft, err = strconv.ParseFloat(tokensLeftStr, 64)
	if err != nil {
		return nil, myErr.WrapErr(fmt.Errorf("failed to parse tokens left from string to floaat : %w", err), myErr.InternalServerErr)
	}

	userBucket := domain.UserBucket{
		Allowed:                    allowed,
		TokensLeft:                 tokensLeft,
		LuaScriptExecutionTimeInMS: totalTime,
	}

	return &userBucket, nil
}
