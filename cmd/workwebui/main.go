package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/thucnq/work/webui"
)
//redis://:123456789@localhost:6379/0
var (
	redisHostPort  = flag.String("redis", "redis://:123456789@localhost:6379/0", "redis hostport")
	redisNamespace = flag.String("ns", "test_service", "redis namespace")
	webHostPort    = flag.String("listen", ":5040", "hostport to listen for HTTP JSON API")
)

func main() {
	flag.Parse()

	fmt.Println("Starting workwebui:")
	fmt.Println("redis = ", *redisHostPort)
	fmt.Println("namespace = ", *redisNamespace)
	fmt.Println("listen = ", *webHostPort)

	pool := newPool(*redisHostPort)

	server := webui.NewServer(*redisNamespace, pool, *webHostPort)
	server.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c

	server.Stop()

	fmt.Println("\nQuitting...")
}

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxActive:   3,
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(addr)
		},
		Wait: true,
	}
}
