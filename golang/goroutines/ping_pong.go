package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/trace"
	"time"
)

type Ball struct{ hits int }

func main() {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	//runtime.GOMAXPROCS(1)
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
	}()

	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()
	table := make(chan *Ball)
	go player("ping", table)
	go player("pong", table)

	table <- new(Ball) // game on; toss the ball

	time.Sleep(50 * time.Second)
	<-table // game over; grab the ball
}

func player(name string, table chan *Ball) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	for {
		fmt.Println("Number of Goroutines:", runtime.NumGoroutine())
		ball := <-table
		ball.hits++
		fmt.Println(name, ball.hits)
		a := []string{}
		for i := 0; i < 100; i++ {
			a = append(a, "ahdkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkks")
		}
		time.Sleep(1000 * time.Millisecond)
		//panic("show me the stacks")
		table <- ball
	}
}
