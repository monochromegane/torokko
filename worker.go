package cargo

import "sync"

type buildWorker struct {
	wg    *sync.WaitGroup
	queue chan *params
}

func (w buildWorker) run() {
	for params := range w.queue {
		go func() {
			cargo{params}.build()
		}()
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
	wg.Wait()
}
