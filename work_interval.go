package gwork

import "time"

type WorkInterval struct {
	workBase
	intervalTime time.Duration
}

func newWorkInterval(pool *WorkPool, intervalTime time.Duration, loopFunc LoopFunc, onError OnError) *WorkInterval {
	w := &WorkInterval{}
	w.loopFunc = loopFunc
	w.onError = onError
	w.init(pool)
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
