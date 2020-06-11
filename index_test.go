package gwork

import (
	"context"
	"fmt"
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
		err = fmt.Errorf("出现错误")
		panic(err)
		return
	}
}

func TestNew(t *testing.T) {
	pool := NewPool(nil)
	pool.AddIntervalWork("testInterval", time.Second*2, testLoopFunc()).Run()
	time.Sleep(time.Second * 5)
	t.Log("close")
	pool.Close()
	pool.Join()
}
func TestNewSchedule(t *testing.T) {
	pool := NewPool(nil)
	h := 10
	m := 9
	se := 10

	s := &Schedule{
		Hour:   &h,
		Minute: &m,
		Second: &se,
	}

	pool.AddScheduleWork("testInterval", s, testLoopFunc()).Run()

	pool.Join()
}
