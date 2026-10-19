package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sakamoto-max/ratelimiter/internal/domain"
	"github.com/sakamoto-max/ratelimiter/internal/opentel"

	"github.com/redis/go-redis/v9"
)

var script = redis.NewScript(`
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

type BucketIface interface {
	CheckLimit(ctx context.Context, userIp string, policy domain.Policy) (*domain.UserBucket, error)
}

func newBucket(client *redis.Client, reg *prometheus.Registry) *Bucket {
	return &Bucket{
		Client: client,
		LuaScriptExecutionTime: promauto.With(reg).NewHistogram(prometheus.HistogramOpts{
			Name: "lua_script_execution_time_in_seconds",
			Help: "Time taken to execute lua script",
			Buckets: []float64{
				0.0005, // 0.5 ms
				0.001,  // 1 ms
				0.0025, // 2.5 ms
				0.005,  // 5 ms
				0.01,   // 10 ms
				0.025,  // 25 ms
				0.05,   // 50 ms
				0.1,    // 100 ms
				0.5,    // 500 ms
				1,      // 1 s
				2,      // 2 s
			},
		},
		),
	}
}

type Bucket struct {
	Client                 *redis.Client
	LuaScriptExecutionTime prometheus.Histogram
}

func (c *Bucket) CheckLimit(ctx context.Context, userIp string, policy domain.Policy) (*domain.UserBucket, error) {

	ctx, span := opentel.CreateNewSpan(ctx, "redis_rate_limit_check")
	defer span.End()

	mainKey := fmt.Sprintf("user_ip:%v:resource_name:%v:bucket", userIp, policy.ResourceName)

	keys := []string{mainKey}

	timeStart := time.Now()

	cmd, err := script.Eval(ctx, c.Client, keys, policy.BucketSize, policy.RefillPerSecond).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to implement the redis script : %w", err)
	}

	c.LuaScriptExecutionTime.Observe(float64(time.Since(timeStart).Seconds()))

	data, ok := cmd.([]any)

	if !ok {
		return nil, fmt.Errorf("failed to parse the data %v from redis", data)
	}

	var allowed int64
	var tokensLeft float64

	allowed, ok = data[0].(int64)
	if !ok {
		return nil, fmt.Errorf("failed to parse the data[0] from redis")
	}

	tokensLeftStr, ok := data[1].(string)
	if !ok {
		return nil, fmt.Errorf("failed to parse the data[1] from redis")
	}

	tokensLeft, err = strconv.ParseFloat(tokensLeftStr, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the data[1] from redis")
	}

	userBucket := domain.UserBucket{
		Allowed:                    allowed,
		TokensLeft:                 tokensLeft,
		LuaScriptExecutionTimeInMS: 0,
	}

	return &userBucket, nil
}
