package gwork

import (
	"context"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func testOnError() OnError {
	return func(work Work, err error) {
		logrus.Errorf("[%s]:%s", work.Id(), err)
	}
}
func testLoopFunc() LoopFunc {
	return func(ctx context.Context) (err error) {
		logrus.Info("run once")
		return
	}
}

func TestNew(t *testing.T) {
	pool := New()
	pool.AddIntervalWork(time.Second*2, testLoopFunc(), testOnError())
	time.Sleep(time.Second * 5)
	t.Log("close")
	pool.Close()
	pool.Join()
}
func TestNewSchedule(t *testing.T) {
	pool := New()
	se := 10

	s := &Schedule{
		Second: &se,
	}

	pool.AddScheduleWork(s, testLoopFunc(), testOnError())

	pool.Join()
}
