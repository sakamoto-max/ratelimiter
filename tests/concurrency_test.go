package tests

import (
	"context"
	"fmt"
	"log"
	"github.com/sakamoto-max/ratelimiter/internal/config"
	"github.com/sakamoto-max/ratelimiter/internal/database"
	"github.com/sakamoto-max/ratelimiter/internal/domain"
	"github.com/sakamoto-max/ratelimiter/internal/interceptors"
	"github.com/sakamoto-max/ratelimiter/internal/repository/cache"
	"github.com/sakamoto-max/ratelimiter/internal/service"
	"sync"
	"testing"

	rateLimiterPb "github.com/sakamoto-max/rate_limiter_proto/shared/rate_limiter"
)

func Test_Concurrency(t *testing.T) {
	config := config.Config{
		Redis: config.RedisConfig{
			Host:     "localhost",
			Port:     "6379",
			Db:       "0",
			UserName: "default",
			Password: "",
		},
	}

	redisClient := database.NewRedisClient(&config)
	cache := cache.NewCache(redisClient)

	resultChan := make(chan *rateLimiterPb.CheckResponse, 100)

	policy := domain.Policy{
		OwnerName:         "max",
		ResourceName:      "resource",
		BucketSize:        10,
		IntervalInSeconds: 60,
		RefillPerSecond:   float64(10) / float64(60),
	}

	cache.Policy.SetPolicy(context.Background(), policy)

	grpcService := service.Grpc{
		Cache: cache,
	}

	var wg sync.WaitGroup

	for i := range 100 {
		wg.Add(1)
		go Go(i, &wg, &grpcService, resultChan)
	}

	wg.Wait()

	totalAllowed, totalBlocked, totalRequests := GetData(resultChan)

	fmt.Println("TotalRequests", totalRequests)
	fmt.Println("TotalAllowed", totalAllowed)
	fmt.Println("TotalBlocked", totalBlocked)

	close(resultChan)
	redisClient.Close()
}

func Go(id int, wg *sync.WaitGroup, grpcService *service.Grpc, reslutChan chan<- *rateLimiterPb.CheckResponse) {
	defer wg.Done()

	ctx := context.WithValue(context.Background(), interceptors.OwnerNameKey, "max")

	resp, err := grpcService.Check(ctx, &rateLimiterPb.CheckRequest{
		ClientIp:     "127.0.0.1",
		ResourceName: "resource",
	})
	if err != nil {
		log.Println("err occured", err)
		return
	}

	reslutChan <- resp
}

func GetData(resultChan chan *rateLimiterPb.CheckResponse) (int, int, int) {

	var TotalAllowed int
	var TotalBlocked int
	var TotalRequests int

	for range 100 {
		data, ok := <-resultChan
		if !ok {
			return TotalAllowed, TotalBlocked, TotalRequests
		}

		switch data.Allowed {
		case true:
			TotalAllowed++
		case false:
			TotalBlocked++
		}

		TotalRequests++
	}

	return TotalAllowed, TotalBlocked, TotalRequests
}
