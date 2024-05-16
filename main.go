package main

import (
	"flag"
	"fmt"
	"math"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var tickerLeakCount int = 0
var tickerLeakThread int = 0
var busyLoopThread int = 0

func main() {
	flag.IntVar(&tickerLeakCount, "ticker-leak-count", 30000, "单线程泄露的Ticker数量")
	flag.IntVar(&tickerLeakThread, "ticker-leak-thread", 2, "ticker泄漏的线程数")
	flag.IntVar(&busyLoopThread, "busy-loop-thread", 1, "空忙线程数")
	flag.Parse()
	// 启动pprof接口
	go func() {
		fmt.Println("PPROF server started on :6060")
		http.ListenAndServe(":6060", nil)
	}()

	for i := 0; i < busyLoopThread; i++ {
		go busyLoop(i)
	}
	for i := 0; i < tickerLeakThread; i++ {
		go tickerLeak(i)
	}

	select {}
}

func busyLoop(id int) {
	fmt.Printf("Starting busy loop thread id: %d\n", id)

	for {
		// 模拟计算密集型任务
		for i := 0; i < 2000000; i++ {
			_ = math.Sqrt(float64(i))
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func tickerLeak(id int) {
	fmt.Printf("Starting ticker leak thread id: %d\n", id)
	i := 0
	var ticker *time.Ticker
	for {
		i++
		if i > tickerLeakCount {
			time.Sleep(time.Second)
			continue
		}

		if ticker != nil {
			ticker = time.NewTicker(3 * time.Millisecond)
		}
		for range ticker.C {
			// do something
			break
		}
	}
}
