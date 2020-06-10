package gwork

import (
	"testing"
	"time"
)

func Test_parseSchedule(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2020-12-31T23:59:59+08:00")
	now, _ = time.Parse(time.RFC3339, "2020-12-29T09:00:59+08:00")

	nowStr := now.Format(time.RFC3339)
	println(nowStr)
	month := time.March
	weekday := time.Wednesday
	hour := 10
	minute := 23

	type args struct {
		schedule *Schedule
	}
	tests := []struct {
		name         string
		args         args
		wantNextTime string
	}{
		{"nextMonthTime", args{&Schedule{Month: &month}}, "2021-03-01T00:00:00+08:00"},
		{"nextWeekday", args{&Schedule{Weekday: &weekday}}, "2020-12-30T00:00:00+08:00"},
		{"nextDay", args{&Schedule{Hour: &hour, Minute: &minute}}, "2020-12-29T10:23:00+08:00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNextTime := parseSchedule(tt.args.schedule, now).Format(time.RFC3339); gotNextTime != tt.wantNextTime {
				t.Errorf("parseSchedule() = %v, want %v", gotNextTime, tt.wantNextTime)
			}
		})
	}
}
