package gwork

import (
	"context"
	"time"
)

type WorkPool struct {
	works  map[string]Work
	ctx    context.Context
	cancel context.CancelFunc
}

func (w *WorkPool) AddIntervalWork(interval time.Duration, loopFunc LoopFunc, onError OnError) {
	w.addWork(newWorkInterval(w, interval, loopFunc, onError))
}
func (w *WorkPool) AddScheduleWork(schedule *Schedule, loopFunc LoopFunc, onError OnError) {
	w.addWork(newWorkSchedule(w, schedule, loopFunc, onError))
}

func (w *WorkPool) addWork(work Work) {
	w.works[work.Id()] = work
	go runWork(work)
}

func runWork(work Work) {
	defer work.finish()
	if work.isRunImmediately() {
		work.loop()
	}
	for {
		select {
		case <-work.done():
			return
		case <-work.getLoopTimer().C:
			work.loop()
		}
	}
}
func (w *WorkPool) Join() {
	for _, work := range w.works {
		work.Join()
	}
}

func (w *WorkPool) Close() error {
	w.cancel()
	return nil
}
