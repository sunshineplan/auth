package main

import (
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func run() {
	if logPath != "" {
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			log.Fatalln("Failed to open log file:", err)
		}
		gin.DefaultWriter = f
		gin.DefaultErrorWriter = f
		log.SetOutput(f)
	}

	if err := initMongo(); err != nil {
		log.Fatalln("Failed to initialize mongodb:", err)
	}

	var redisStore struct{ Endpoint, Password string }
	if err := meta.Get("account_redis", &redisStore); err != nil {
		log.Fatal(err)
	}
	store, err := redis.NewStore(10, "tcp", redisStore.Endpoint, redisStore.Password, []byte(secret))
	if err != nil {
		log.Fatal(err)
	}
	err, realStore := redis.GetRedisStore(store)
	if err != nil {
		log.Fatal(err)
	}
	realStore.DefaultMaxAge = 60 * 60 * 24
	realStore.SetKeyPrefix("account")

	router := gin.Default()
	server.Handler = router
	router.Use(sessions.Sessions("session", store))

	router.POST("/login", login)
	router.POST("/chgpwd", chgpwd)
	router.POST("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.String(200, "bye")
	})

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
