package redis

import (
	"context"
	"errors"
	redis "github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"net/url"
	"os"
)

type RedisService struct {
	Red *redis.Client
	log *log.Logger
}

func verifyParams(urlParams *url.Values) error {
	if ok := urlParams.Has("int1"); !ok {
		return errors.New("int1 missing")
	}
	if ok := urlParams.Has("int2"); !ok {
		return errors.New("int2 missing")
	}
	if ok := urlParams.Has("limit"); !ok {
		return errors.New("limit missing")
	}
	if ok := urlParams.Has("str1"); !ok {
		return errors.New("str1 missing")
	}
	if ok := urlParams.Has("str2"); !ok {
		return errors.New("str2 missing")
	}

	return nil
}

func StartRedis() *RedisService {
	ctx := context.Background()

	rds := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var redisService = &RedisService{
		Red: rds,
		log: log.New(os.Stdout, "[REDIS]:	", 3),
	}

	_, err := rds.Ping(ctx).Result()
	if err != nil {
		redisService.log.Fatalf("	error initializing redis client : %s", err.Error())
	} else {
		redisService.log.Printf("	Redis client successfully initialized")
	}

	return redisService
}

func (rds *RedisService) GetConn(ctx context.Context) *redis.Conn {
	return rds.Red.Conn(ctx)
}

func (rds *RedisService) SetKey(key string, value interface{}, ctx context.Context) error {
	err := rds.Red.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rds *RedisService) GetKey(key string, ctx context.Context) (string, error) {
	stringCmd := rds.Red.Get(ctx, key)
	err := stringCmd.Err()
	if err != nil {
		return "", err
	}
	return stringCmd.Val(), nil
}

func (rds *RedisService) Simple(res http.ResponseWriter, req *http.Request) {
	urlParams := req.URL.Query()

	err := verifyParams(&urlParams)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(err.Error()))
	}
}
