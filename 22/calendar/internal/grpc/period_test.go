package grpc

import (
	"testing"
	"time"
)

func TestGetDayPeriod(t *testing.T) {
	now := time.Date(2019, 11, 17, 14, 33, 12, 0, time.UTC)
	period, err := NewDayPeriod(now)
	if err != nil {
		t.Fatalf("must not error happened on constuction period %s\n", err)
	}
	expectedStart := ts(2019, 11, 17, 0, 0)
	expectedEnd := ts(2019, 11, 17, 23, 59)
	if !isTimestampEquals(period.start, expectedStart) {
		t.Errorf("start must be %s insteadof %s", expectedStart, period.start)
	}
	if !isTimestampEquals(period.end, expectedEnd) {
		t.Errorf("end must be %s insteadof %s", expectedEnd, period.end)
	}
}

func TestGetWeekPeriod1(t *testing.T) {
	now := time.Date(2019, 11, 17, 14, 33, 12, 0, time.UTC)
	period, err := NewWeekPeriod(now)
	if err != nil {
		t.Fatalf("must not error happened on constuction period %s\n", err)
	}
	expectedStart := ts(2019, 11, 11, 0, 0)
	expectedEnd := ts(2019, 11, 17, 23, 59)
	if !isTimestampEquals(period.start, expectedStart) {
		t.Errorf("start must be %s insteadof %s", expectedStart, period.start)
	}
	if !isTimestampEquals(period.end, expectedEnd) {
		t.Errorf("end must be %s insteadof %s", expectedEnd, period.end)
	}
}

func TestGetWeekPeriod2(t *testing.T) {
	now := time.Date(2019, 11, 13, 14, 33, 12, 0, time.UTC)
	period, err := NewWeekPeriod(now)
	if err != nil {
		t.Fatalf("must not error happened on constuction period %s\n", err)
	}
	expectedStart := ts(2019, 11, 11, 0, 0)
	expectedEnd := ts(2019, 11, 17, 23, 59)
	if !isTimestampEquals(period.start, expectedStart) {
		t.Errorf("start must be %s insteadof %s", expectedStart, period.start)
	}
	if !isTimestampEquals(period.end, expectedEnd) {
		t.Errorf("end must be %s insteadof %s", expectedEnd, period.end)
	}
}

func TestGetWeekPeriod3(t *testing.T) {
	now := time.Date(2019, 11, 11, 14, 33, 12, 0, time.UTC)
	period, err := NewWeekPeriod(now)
	if err != nil {
		t.Fatalf("must not error happened on constuction period %s\n", err)
	}
	expectedStart := ts(2019, 11, 11, 0, 0)
	expectedEnd := ts(2019, 11, 17, 23, 59)
	if !isTimestampEquals(period.start, expectedStart) {
		t.Errorf("start must be %s insteadof %s", expectedStart, period.start)
	}
	if !isTimestampEquals(period.end, expectedEnd) {
		t.Errorf("end must be %s insteadof %s", expectedEnd, period.end)
	}
}

func TestGetMonthPeriod1(t *testing.T) {
	now := time.Date(2019, 11, 11, 14, 33, 12, 0, time.UTC)
	period, err := NewMonthPeriod(now)
	if err != nil {
		t.Fatalf("must not error happened on constuction period %s\n", err)
	}
	expectedStart := ts(2019, 11, 01, 0, 0)
	expectedEnd := ts(2019, 11, 30, 23, 59)
	if !isTimestampEquals(period.start, expectedStart) {
		t.Errorf("start must be %s insteadof %s", expectedStart, period.start)
	}
	if !isTimestampEquals(period.end, expectedEnd) {
		t.Errorf("end must be %s insteadof %s", expectedEnd, period.end)
	}
}

func TestGetMonthPeriod2(t *testing.T) {
	now := time.Date(2019, 12, 11, 14, 33, 12, 0, time.UTC)
	period, err := NewMonthPeriod(now)
	if err != nil {
		t.Fatalf("must not error happened on constuction period %s\n", err)
	}
	expectedStart := ts(2019, 12, 01, 0, 0)
	expectedEnd := ts(2019, 12, 31, 23, 59)
	if !isTimestampEquals(period.start, expectedStart) {
		t.Errorf("start must be %s insteadof %s", expectedStart, period.start)
	}
	if !isTimestampEquals(period.end, expectedEnd) {
		t.Errorf("end must be %s insteadof %s", expectedEnd, period.end)
	}
}

func TestGetMonthPeriod3(t *testing.T) {
	now := time.Date(2020, 2, 11, 14, 33, 12, 0, time.UTC)
	period, err := NewMonthPeriod(now)
	if err != nil {
		t.Fatalf("must not error happened on constuction period %s\n", err)
	}
	expectedStart := ts(2020, 02, 01, 0, 0)
	expectedEnd := ts(2020, 02, 29, 23, 59)
	if !isTimestampEquals(period.start, expectedStart) {
		t.Errorf("start must be %s insteadof %s", expectedStart, period.start)
	}
	if !isTimestampEquals(period.end, expectedEnd) {
		t.Errorf("end must be %s insteadof %s", expectedEnd, period.end)
	}
}

func TestGetMonthPeriod4(t *testing.T) {
	now := time.Date(2019, 2, 11, 14, 33, 12, 0, time.UTC)
	period, err := NewMonthPeriod(now)
	if err != nil {
		t.Fatalf("must not error happened on constuction period %s\n", err)
	}
	expectedStart := ts(2019, 02, 01, 0, 0)
	expectedEnd := ts(2019, 02, 28, 23, 59)
	if !isTimestampEquals(period.start, expectedStart) {
		t.Errorf("start must be %s insteadof %s", expectedStart, period.start)
	}
	if !isTimestampEquals(period.end, expectedEnd) {
		t.Errorf("end must be %s insteadof %s", expectedEnd, period.end)
	}
}
