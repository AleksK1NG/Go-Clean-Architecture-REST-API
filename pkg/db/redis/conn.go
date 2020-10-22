package redis

// Create new redis Pool
//func NewRedisPool(cfg *config.Config) (*redis.Pool, error) {
//	redisHost := cfg.Redis.RedisAddr
//
//	if redisHost == "" {
//		redisHost = ":6379"
//	}
//	Pool := newPool(redisHost)
//
//	cleanupHook(Pool)
//
//	if err := pingRedis(Pool); err != nil {
//		return nil, err
//	}
//
//	return Pool, nil
//}

//func newPool(server string) *redis.Pool {
//
//	return &redis.Pool{
//		MaxIdle:     50,
//		MaxActive:   300,
//		Wait:        true,
//		IdleTimeout: 240 * time.Second,
//
//		Dial: func() (redis.Conn, error) {
//			c, err := redis.Dial("tcp", server)
//			if err != nil {
//				return nil, err
//			}
//			return c, err
//		},
//
//		TestOnBorrow: func(c redis.Conn, t time.Time) error {
//			_, err := c.Do("PING")
//			return err
//		},
//	}
//}

//func cleanupHook(Pool *redis.Pool) {
//
//	c := make(chan os.Signal, 1)
//	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
//
//	go func() {
//		<-c
//		if err := Pool.Close(); err != nil {
//			log.Printf("POOL CLOSE ERROR: %s", err.Error())
//			return
//		}
//		os.Exit(0)
//	}()
//}
//
//// Redis ping method
//func pingRedis(Pool *redis.Pool) error {
//	conn := Pool.Get()
//	defer conn.Close()
//
//	ping, err := redis.String(conn.Do("PING"))
//	if err != nil {
//		return fmt.Errorf("cannot 'PING' db: %v", err)
//	}
//	log.Printf("REDIS PING: %v", ping)
//	return nil
//}
