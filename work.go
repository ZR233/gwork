package gwork

import (
	"context"
	"github.com/ZR233/gwork/errors"
	"github.com/sirupsen/logrus"
	"time"
)

type Work interface {
	SetId(string)
	Id() string
	Cancel()
	Join()
	SetDescription(string)
	GetDescription() string
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
	id             string
	runImmediately bool
	description    string
	workPool       *WorkPool
	loopFunc       LoopFunc
	onError        OnError
	ctx            context.Context
	cancel         context.CancelFunc
	stopped        chan bool
	loopTimer      *time.Timer
}

func (w *workBase) isRunImmediately() bool {
	return w.runImmediately
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

func (w *workBase) SetDescription(str string) {
	w.description = str
}
func (w *workBase) GetDescription() string {
	return w.description
}

func (w *workBase) init(workPool *WorkPool) {
	w.workPool = workPool
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
		w.onError(work, err)
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
