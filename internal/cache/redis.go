package cache

// // RedisClient holds the Redis connection.
// var RedisClient *_redis.Client
//
// // InitRedis connects to Redis.
// func InitRedis(selectDB ...int) {
// 	var redisHost = os.Getenv("REDIS_HOST")
// 	var redisPassword = os.Getenv("REDIS_PASSWORD")
//
// 	redisDB := 0
// 	if len(selectDB) > 0 {
// 		redisDB = selectDB[0]
// 	}
//
// 	RedisClient = _redis.NewClient(&_redis.Options{
// 		Addr:     redisHost,
// 		Password: redisPassword,
// 		DB:       redisDB,
// 	})
//
// 	if err := RedisClient.Ping().Err(); err != nil {
// 		log.Fatalf("Critical Error: Could not connect to Redis: %v", err)
// 	}
//
// }
//
