package main

import (
	"asynq-task/constants"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

type ExampleTaskPayload struct {
	UserID string
	Msg    string
	// 业务需要的其他字段
}

// 生产 异步队列数据
func NewExampleTask(userID string, msg string) (*asynq.Task, error) {
	payload, err := json.Marshal(ExampleTaskPayload{UserID: userID, Msg: msg})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(constants.TypeExampleTask, payload,
		asynq.MaxRetry(10),
		asynq.Timeout(3*time.Minute),
		asynq.ProcessIn(3*time.Second)), nil
}

var client *asynq.Client

func main() {

	client = asynq.NewClient(asynq.RedisClientOpt{Addr: constants.RedisAddr, Password: constants.RedisPasswd, DB: 0})
	defer client.Close()

	//go startExampleTask()
	startExampleTask()

	//startGithubUpdate() // 定时触发
}

func startExampleTask() {

	fmt.Println("开始执行一次性的任务")
	// 立刻执行
	task1, err := NewExampleTask("10001", "mashangzhixing!")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}

	info, err := client.Enqueue(task1)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("task1 -> enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// 10秒后执行(定时执行)
	task2, err := NewExampleTask("10002", "10s houzhixing")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}

	info, err = client.Enqueue(task2, asynq.ProcessIn(10*time.Second))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("task2 -> enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// 30s后执行(定时执行)
	task3, err := NewExampleTask("10003", "30s houzhixing")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}

	theTime := time.Now().Add(30 * time.Second)
	info, err = client.Enqueue(task3, asynq.ProcessAt(theTime))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("task3 -> enqueued task: id=%s queue=%s", info.ID, info.Queue)
}
