package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type buildWorker struct {
	wg    *sync.WaitGroup
	queue chan *params
}

func (w buildWorker) run() {
	for params := range w.queue {
		newCargo(params).build()
	}
	w.wg.Done()
}

func startWorker(queue chan *params, limit int) {
	var wg sync.WaitGroup
	wg.Add(limit)

	for i := 0; i < limit; i++ {
		go buildWorker{
			wg:    &wg,
			queue: queue,
		}.run()
	}

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGQUIT)
	go func() {
		for _ = range signalChan {
			close(queue)
		}
	}()

	wg.Wait()
}
