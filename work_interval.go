package gwork

import "time"

type WorkInterval struct {
	workBase
	intervalTime time.Duration
}

func (w *WorkInterval) WithOptions(options *WorkOptions) Work {
	w.setOptions(options)
	return w
}

func newWorkInterval(pool *WorkPool,
	name string,
	intervalTime time.Duration,
	loopFunc LoopFunc) *WorkInterval {
	w := &WorkInterval{}
	w.init(pool, name, loopFunc)
	w.intervalTime = intervalTime
	w.loopTimer = time.NewTimer(w.intervalTime)
	return w
}

func (w *WorkInterval) loop() {
	if w.loopTimer != nil {
		w.loopTimer.Stop()
	}

	w.loopTimer = time.NewTimer(w.intervalTime)
	w.excLoopFunc(w)
}

func (w *WorkInterval) Run() {
	w.checkOptions()
	go runWork(w)
}
