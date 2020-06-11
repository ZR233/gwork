package gwork

import (
	"context"
	"time"
)

type WorkPool struct {
	works map[string]Work
	WorkPoolOptions
	ctx    context.Context
	cancel context.CancelFunc
}

func (w *WorkPool) AddIntervalWork(name string, interval time.Duration, loopFunc LoopFunc) (work Work) {
	work = newWorkInterval(w, name, interval, loopFunc)
	w.addWork(work)
	return
}
func (w *WorkPool) AddScheduleWork(name string, schedule *Schedule, loopFunc LoopFunc) (work Work) {
	work = newWorkSchedule(w, name, schedule, loopFunc)
	w.addWork(work)
	return
}

func (w *WorkPool) addWork(work Work) {
	w.works[work.Id()] = work
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
