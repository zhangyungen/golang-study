package goroutine

type Worker struct {
	stop chan struct{}
	done chan struct{}
}

func NewWorker(args ...func()) *Worker {
	w := &Worker{
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	go w.doWork()
	return w
}
func (w *Worker) doWork() {
	defer close(w.done)
	for {
		<-w.stop
	}

}

// Shutdown 告诉 worker 停止
// 并等待它完成。
func (w *Worker) Shutdown() {
	close(w.stop)
	<-w.done
}
