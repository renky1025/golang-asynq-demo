package test

import (
	"fmt"
	"os"
	"sync"
	"testing"

	ants "github.com/panjf2000/ants/v2"
)

func wrapper(i int, wg *sync.WaitGroup) func() {
	return func() {
		fmt.Printf("hello from task:%d\n", i)
		if i%2 == 0 {
			panic(fmt.Sprintf("panic from task:%d", i))
		}
		wg.Done()
	}
}
func panicHandler(err interface{}) {
	fmt.Printf("hello from error => %v", err)
	fmt.Fprintln(os.Stderr, err)
}
func Test_Pool(b *testing.T) {
	// ants.WithMaxBlockingTasks(2) 最大等待长度2，可以设置等待队列的最大长度。超过这个长度，提交任务直接返回错误
	p, _ := ants.NewPool(5, ants.WithNonblocking(true), ants.WithPanicHandler(panicHandler)) // 非阻塞线程池, 拦截非法crash
	defer p.Release()

	var wg sync.WaitGroup
	wg.Add(5)
	for i := 1; i <= 6; i++ {
		err := p.Submit(wrapper(i, &wg))
		if err != nil {
			fmt.Printf("task:%d err:%v\n", i, err)
			wg.Done()
		}
	}

	wg.Wait()

}
