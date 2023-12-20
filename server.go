package main

import (
	"asynq-task/constants"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/hibiken/asynq"
)

var AsynqServer *asynq.Server // 异步任务server

func initTaskServer() error {
	// 初始化异步任务服务端
	AsynqServer = asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     constants.RedisAddr,
			Password: constants.RedisPasswd, //与client对应
			DB:       0,
		},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 100,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6, //关键队列中的任务将被处理 60% 的时间
				"default":  3, //默认队列中的任务将被处理 30% 的时间
				"low":      1, //低队列中的任务将被处理 10% 的时间
			},
			// See the godoc for other configuration options
		},
	)
	return nil
}

func main() {
	initTaskServer()
	mux := asynq.NewServeMux()
	// 消费 异步队列数据
	mux.HandleFunc(constants.TypeExampleTask, HandleExampleTask)
	// ...register other handlers...

	if err := AsynqServer.Run(mux); err != nil {
		fmt.Printf("could not run asynq server: %v", err)
	}
}

func HandleExampleTask(ctx context.Context, t *asynq.Task) error {

	res := make(map[string]string)

	spew.Dump("t.Payload() is:", t.Payload())
	err := json.Unmarshal(t.Payload(), &res)
	if err != nil {
		fmt.Printf("rum session, can not parse payload: %s,  err: %v", t.Payload(), err)
		return nil
	}
	//-----------具体处理逻辑------------
	spew.Println("拿到的入参为:", res, "接下来将进行具体处理")
	fmt.Println()
	// 模拟具体的处理
	time.Sleep(5 * time.Second)
	fmt.Println("--------------处理了5s，处理完成-----------------")

	return nil

}
