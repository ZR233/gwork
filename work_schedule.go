package gwork

import "time"

type Schedule struct {
	Month   *time.Month
	Weekday *time.Weekday
	Day     *int
	Hour    *int
	Minute  *int
	Second  *int
}

func (s *Schedule) SetMonth(v time.Month) *Schedule {
	s.Month = &v
	return s
}
func (s *Schedule) SetWeekday(v time.Weekday) *Schedule {
	s.Weekday = &v
	return s
}
func (s *Schedule) SetDay(v int) *Schedule {
	s.Day = &v
	return s
}
func (s *Schedule) SetHour(v int) *Schedule {
	s.Hour = &v
	return s
}
func (s *Schedule) SetMinute(v int) *Schedule {
	s.Minute = &v
	return s
}
func (s *Schedule) SetSecond(v int) *Schedule {
	s.Second = &v
	return s
}

type WorkSchedule struct {
	workBase
	schedule *Schedule
}

func newWorkSchedule(pool *WorkPool, name string, schedule *Schedule, loopFunc LoopFunc) *WorkSchedule {
	w := &WorkSchedule{}

	w.init(pool, name, loopFunc)
	w.schedule = schedule
	w.newTimer()
	return w
}
func initNils(list ...**int) {
	for _, v := range list {
		initNil(v)
	}
}

func initNil(i **int) {
	if *i == nil {
		n := 0
		*i = &n
	}
}
func nextMonthTime(schedule *Schedule, now time.Time) (next time.Time) {
	initNils(&schedule.Day, &schedule.Hour, &schedule.Minute, &schedule.Second)
	next = time.Date(
		now.Year(),
		*schedule.Month,
		*schedule.Day,
		*schedule.Hour,
		*schedule.Minute,
		*schedule.Second, 0, now.Location())
	if next.Sub(now) <= 0 {
		next = next.AddDate(1, 0, 0)
	}

	return
}
func nextMonthDay(schedule *Schedule, now time.Time) (next time.Time) {
	initNils(&schedule.Hour, &schedule.Minute, &schedule.Second)
	next = time.Date(
		now.Year(),
		now.Month(),
		*schedule.Day,
		*schedule.Hour,
		*schedule.Minute,
		*schedule.Second, 0, now.Location())
	if next.Sub(now) <= 0 {
		next = next.AddDate(0, 1, 0)
	}
	return
}
func nextWeekDay(schedule *Schedule, now time.Time) (next time.Time) {
	initNils(&schedule.Hour, &schedule.Minute, &schedule.Second)
	days := *schedule.Weekday - now.Weekday()
	if days < 0 {
		days += 7
	}
	next = time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		*schedule.Hour,
		*schedule.Minute,
		*schedule.Second, 0, now.Location())
	if next.Sub(now) <= 0 {
		next = next.AddDate(0, 0, int(days))
	}
	return
}
func nextDay(schedule *Schedule, now time.Time) (next time.Time) {
	initNils(&schedule.Hour, &schedule.Minute, &schedule.Second)
	next = time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		*schedule.Hour,
		*schedule.Minute,
		*schedule.Second, 0, now.Location())
	if next.Sub(now) <= 0 {
		next = next.AddDate(0, 0, 1)
	}
	return
}
func nextHour(schedule *Schedule, now time.Time) (next time.Time) {
	initNils(&schedule.Minute, &schedule.Second)
	next = time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		*schedule.Minute,
		*schedule.Second, 0, now.Location())
	if next.Sub(now) <= 0 {
		next = next.Add(time.Hour)
	}
	return
}
func nextMinute(schedule *Schedule, now time.Time) (next time.Time) {
	initNils(&schedule.Second)
	next = time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		*schedule.Second, 0, now.Location())
	if next.Sub(now) <= 0 {
		next = next.Add(time.Minute)
	}
	return
}
func parseSchedule(schedule *Schedule, now time.Time) (nextTime time.Time) {

	if schedule.Month != nil {
		nextTime = nextMonthTime(schedule, now)
	} else {
		if schedule.Day != nil {
			nextTime = nextMonthDay(schedule, now)
		} else {
			if schedule.Weekday != nil {
				nextTime = nextWeekDay(schedule, now)
			} else {
				if schedule.Hour != nil {
					nextTime = nextDay(schedule, now)
				} else {
					if schedule.Minute != nil {
						nextTime = nextHour(schedule, now)
					} else {
						if schedule.Second != nil {
							nextTime = nextMinute(schedule, now)
						}
					}
				}
			}
		}
	}

	return
}
func (w *WorkSchedule) newTimer() {

	nextTime := parseSchedule(w.schedule, time.Now())

	w.loopTimer = time.NewTimer(nextTime.Sub(time.Now()))
}

func (w *WorkSchedule) loop() {
	if w.loopTimer != nil {
		w.loopTimer.Stop()
	}

	w.newTimer()
	w.excLoopFunc(w)
}
func (w *WorkSchedule) Run() {
	w.checkOptions()
	go runWork(w)
}

func (w *WorkSchedule) WithOptions(options *WorkOptions) Work {
	w.setOptions(options)
	return w
}
