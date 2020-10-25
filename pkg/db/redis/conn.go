package redis

import (
	"encoding/json"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Redis db interface
type RedisPool interface {
	//GetPool() *redis.Pool
	PingContext(ctx context.Context) error
	GetBytesContext(ctx context.Context, key string) ([]byte, error)
	SetBytes(key string, value []byte) error
	SetexBytes(key string, durationSec int, value []byte) error
	Exists(key string) (bool, error)
	Delete(key string) error
	DeleteContext(ctx context.Context, key string) error
	GetKeys(pattern string) ([]string, error)
	Incr(counterKey string) (int, error)
	SetEXJSON(key string, seconds int, value interface{}) error
	SetexJSONContext(ctx context.Context, key string, seconds int, value interface{}) error
	GetJSON(key string, model interface{}) error
	GetJSONContext(ctx context.Context, key string, model interface{}) error
}

// Redis client
type RedisClient struct {
	config *config.Config
	client *redis.Client
}

//// Returns new redis client
//func NewRedisClient(config *config.Config) *RedisClient {
//	redisHost := config.Redis.RedisAddr
//
//	if redisHost == "" {
//		redisHost = ":6379"
//	}
//	pool := newPool(redisHost)
//
//	cleanupHook(pool)
//	return &RedisClient{config: config, pool: pool}
//}

// Returns new redis client
func NewRedisClient(config *config.Config) *RedisClient {
	redisHost := config.Redis.RedisAddr

	if redisHost == "" {
		redisHost = ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		MinIdleConns: 200,
		PoolSize:     12000,
		PoolTimeout:  240 * time.Second,
		Password:     "", // no password set
		DB:           0,  // use default DB
	})

	cleanupHook(client)
	return &RedisClient{config: config, client: client}
}

func cleanupHook(client *redis.Client) {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		client.Close()
		os.Exit(0)
	}()
}

//func newPool(server string) *redis.Pool {
//
//	return &redis.Pool{
//
//		MaxIdle:     100,
//		IdleTimeout: 10 * time.Second,
//		MaxActive:   12000,
//		Wait:        true,
//
//		Dial: func() (redis.Conn, error) {
//			c, err := redis.Dial("tcp", server)
//			if err != nil {
//				return nil, err
//			}
//			return c, err
//		},
//
//		//TestOnBorrow: func(c redis.Conn, t time.Time) error {
//		//	_, err := c.Do("PING")
//		//	return err
//		//},
//	}
//}

// Get pool
func (r *RedisClient) GetPool() *redis.Client {
	return r.client
}

// Redis ping method
func (r *RedisClient) PingContext(ctx context.Context) error {
	//conn, err := r.pool.GetContext(ctx)
	//if err != nil {
	//	return err
	//}
	//defer conn.Close()
	//
	//ping, err := redis.String(conn.Do("PING"))
	//if err != nil {
	//	return fmt.Errorf("cannot 'PING' db: %v", err)
	//}
	//log.Printf("PING: %v", ping)

	s := r.client.Ping(ctx).String()
	log.Printf("PING: %v", s)
	return nil
}

// Redis ping method
func (r *RedisClient) Ping() error {
	//conn := r.pool.Get()
	//defer conn.Close()
	//
	//ping, err := redis.String(conn.Do("PING"))
	//if err != nil {
	//	return fmt.Errorf("cannot 'PING' db: %v", err)
	//}
	//log.Printf("PING: %v", ping)
	return nil
}

// Get by key string, return []byte
func (r *RedisClient) GetBytesContext(ctx context.Context, key string) ([]byte, error) {
	//conn, err := r.pool.GetContext(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//defer conn.Close()
	//
	//var data []byte
	//data, err = redis.Bytes(conn.Do("GET", key))
	//if err != nil {
	//	return data, fmt.Errorf("error getting key %s: %v", key, err)
	//}
	//return data, err
	return nil, nil
}

// Set by key string, return []byte
func (r *RedisClient) SetBytes(key string, value []byte) error {
	//conn := r.pool.Get()
	//defer conn.Close()
	//
	//_, err := conn.Do("SET", key, value)
	//if err != nil {
	//	v := string(value)
	//	if len(v) > 15 {
	//		v = v[0:12] + "..."
	//	}
	//	return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	//}
	//return err
	return nil
}

// Setex by key string, return []byte
func (r *RedisClient) SetexBytes(key string, durationSec int, value []byte) error {

	//conn := r.pool.Get()
	//defer conn.Close()
	//
	//_, err := conn.Do("SETEX", key, durationSec, value)
	//if err != nil {
	//	v := string(value)
	//	if len(v) > 15 {
	//		v = v[0:12] + "..."
	//	}
	//	return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	//}
	return nil
}

// Exists by key string, return bool
func (r *RedisClient) Exists(key string) (bool, error) {
	//conn := r.pool.Get()
	//defer conn.Close()
	//
	//ok, err := redis.Bool(conn.Do("EXISTS", key))
	//if err != nil {
	//	return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	//}
	//return ok, err
	return false, nil
}

// Delete by key string
func (r *RedisClient) Delete(key string) error {
	//conn := r.pool.Get()
	//defer conn.Close()
	//
	//_, err := conn.Do("DEL", key)
	//return err
	return nil
}

// Delete by key string
func (r *RedisClient) DeleteContext(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
	//conn, err := r.pool.GetContext(ctx)
	//if err != nil {
	//	return err
	//}
	//defer conn.Close()
	//
	//_, err = conn.Do("DEL", key)
	//return err
}

// Get by keys string
func (r *RedisClient) GetKeys(pattern string) ([]string, error) {
	//conn := r.pool.Get()
	//defer conn.Close()
	//
	//iter := 0
	//var keys []string
	//for {
	//	arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", pattern))
	//	if err != nil {
	//		return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
	//	}
	//
	//	iter, _ = redis.Int(arr[0], nil)
	//	k, _ := redis.Strings(arr[1], nil)
	//	keys = append(keys, k...)
	//
	//	if iter == 0 {
	//		break
	//	}
	//}
	//
	//return keys, nil
	return nil, nil
}

// Incr by key string
func (r *RedisClient) Incr(counterKey string) (int, error) {
	//conn := r.pool.Get()
	//defer conn.Close()
	//
	//return redis.Int(conn.Do("INCR", counterKey))
	return 0, nil
}

// Set JSON value
func (r *RedisClient) SetEXJSON(key string, seconds int, value interface{}) error {
	//conn := r.pool.Get()
	//defer conn.Close()
	//
	//bytes, err := json.Marshal(&value)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = redis.String(conn.Do("SETEX", key, seconds, bytes))
	//if err != nil {
	//	return err
	//}

	return nil
}

// Set JSON value
func (r *RedisClient) SetexJSONContext(ctx context.Context, key string, seconds int, value interface{}) error {
	bytes, err := json.Marshal(&value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, bytes, time.Second*60).Err()

	//conn, err := r.pool.GetContext(ctx)
	//if err != nil {
	//	return err
	//}
	//defer conn.Close()
	//
	//bytes, err := json.Marshal(&value)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = redis.String(conn.Do("SETEX", key, seconds, bytes))
	//if err != nil {
	//	return err
	//}
	//
	//return nil
}

// Get JSON value
func (r *RedisClient) GetJSONContext(ctx context.Context, key string, model interface{}) error {
	bytes, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, &model); err != nil {
		return err
	}
	return nil

	//ctx, _ = context.WithTimeout(context.Background(), time.Second*5)
	//conn, err := r.pool.GetContext(ctx)
	//if err != nil {
	//	return err
	//}
	//defer conn.Close()
	//
	//bytes, err := redis.Bytes(conn.Do("GET", key))
	//if err != nil {
	//	return err
	//}
	//
	//if err := json.Unmarshal(bytes, &model); err != nil {
	//	return err
	//}
	//
	//return nil
}

// Get JSON value
func (r *RedisClient) GetJSON(key string, model interface{}) error {
	//conn := r.pool.Get()
	//defer conn.Close()
	//
	//bytes, err := redis.Bytes(conn.Do("GET", key))
	//if err != nil {
	//	return err
	//}
	//
	//if err := json.Unmarshal(bytes, &model); err != nil {
	//	return err
	//}

	return nil
}

// Get JSON value
func (r *RedisClient) GetIfExistsJSON(key string, model interface{}) error {
	//conn := r.pool.Get()
	//defer conn.Close()
	//
	//ok, err := redis.Bool(conn.Do("EXISTS", key))
	//if err != nil {
	//	return err
	//}
	//if !ok {
	//	return httpErrors.NotExists
	//}
	//
	//bytes, err := redis.Bytes(conn.Do("GET", key))
	//if err != nil {
	//	return err
	//}
	//
	//if err := json.Unmarshal(bytes, &model); err != nil {
	//	return err
	//}

	return nil
}
