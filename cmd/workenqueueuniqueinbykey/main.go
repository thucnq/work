package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/thucnq/work"
	"github.com/gomodule/redigo/redis"
)

var redisHostPort = flag.String("redis", ":6379", "redis hostport")
var redisNamespace = flag.String("ns", "work", "redis namespace")
var jobName = flag.String("job", "olsl", "job name")
var jobArgs = flag.String("args", "{}", "job arguments")

func main() {
	flag.Parse()
	key := "unique"

	if *jobName == "" {
		fmt.Println("no job specified")
		os.Exit(1)
	}

	pool := newPool(*redisHostPort)

	var args map[string]interface{}
	err := json.Unmarshal([]byte(*jobArgs), &args)
	if err != nil {
		fmt.Println("invalid args:", err)
		os.Exit(1)
	}

	en := work.NewEnqueuer(*redisNamespace, pool)
	go en.EnqueueUniqueInByKey(*jobName, 10,  args, map[string]interface{}{"key": key})

	wp := work.NewWorkerPool(context{}, 5, *redisNamespace, pool)
	wp.Job("foobar", epsilonHandler)
	wp.Start()
}
func epsilonHandler(job *work.Job) error {
	fmt.Println("epsilon")
	time.Sleep(time.Second)

	if rand.Intn(2) == 0 {
		return fmt.Errorf("random error")
	}
	return nil
}
type context struct{}
func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxActive:   20,
		MaxIdle:     20,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			return c, nil
			//return redis.NewLoggingConn(c, log.New(os.Stdout, "", 0), "redis"), err
		},
		Wait: true,
		//TestOnBorrow: func(c redis.Conn, t time.Time) error {
		//	_, err := c.Do("PING")
		//	return err
		//},
	}
}