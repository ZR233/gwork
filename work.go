package gwork

import (
	"context"
	"fmt"
	"github.com/ZR233/gwork/errors"
	"github.com/sirupsen/logrus"
	"path"
	"sync"
	"time"
)

type Work interface {
	Run()
	WithOptions(options *WorkOptions) Work
	Id() string
	Cancel()
	Join()
	loop()
	isRunImmediately() bool
	//ctx done
	done() <-chan struct{}
	getLoopTimer() *time.Timer
	//work all loop stopped
	finish()
}

type OnError func(work Work, err error)
type LoopFunc func(ctx context.Context) (err error)

type workBase struct {
	id string
	WorkOptions
	workPool  *WorkPool
	loopFunc  LoopFunc
	ctx       context.Context
	cancel    context.CancelFunc
	stopped   chan bool
	loopTimer *time.Timer
	sync.Mutex
}

func (w *workBase) isRunImmediately() bool {
	return w.RunImmediately
}
func (w *workBase) Id() string {
	return w.id
}
func (w *workBase) done() <-chan struct{} {
	return w.ctx.Done()
}
func (w *workBase) getLoopTimer() *time.Timer {
	return w.loopTimer
}

func (w *workBase) SetId(id string) {
	w.id = id
}
func (w *workBase) checkOptions() {
	if w.OnError == nil {
		logrus.Warn(fmt.Sprintf("work[%s] using default OnError", w.Id()))
		w.OnError = defaultOnError()
	}
}

func (w *workBase) init(workPool *WorkPool, name string, loopFunc LoopFunc) {
	w.loopFunc = loopFunc
	w.workPool = workPool
	w.Name = name
	if w.Name == "" {
		panic(fmt.Errorf("work name not define"))
	}
	w.id = path.Join(workPool.Prefix, w.Name)
	if w.loopFunc == nil {
		panic(fmt.Errorf("work LoopFunc not define"))
	}

	w.WorkOptions = *NewWorkOptions()

	w.stopped = make(chan bool)
	w.ctx, w.cancel = context.WithCancel(workPool.ctx)
}

func (w *workBase) excLoop() (err error) {
	defer func() {
		p := errors.HandleRecover(recover())
		if p != nil {
			err = p
		}
	}()
	err = w.loopFunc(w.ctx)
	return
}

func (w *workBase) excLoopFunc(work Work) {
	var (
		err            error
		isOnErrorPanic = false
	)
	defer func() {
		//TODO 远程调用发送报告
	}()

	defer func() {
		if isOnErrorPanic {
			logrus.Error(err)
		}
	}()

	defer func() {
		p := errors.HandleRecover(recover())
		if p != nil {
			err = p
			isOnErrorPanic = true
		}
	}()

	if w.ctx.Err() != nil {
		return
	}

	err = w.loopFunc(w.ctx)
	if err != nil {
		w.OnError(work, err)
	}
}

func (w *workBase) Cancel() {
	w.cancel()
}
func (w *workBase) finish() {
	w.stopped <- true
}
func (w *workBase) Join() {
	<-w.stopped
}
func defaultOnError() OnError {
	return func(work Work, err error) {
		logrus.Error(fmt.Sprintf("work[%s] error: %s", work.Id(), err))
	}
}
