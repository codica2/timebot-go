package models

import (
	"time"
	"fmt"
)

type Date struct {
	_time time.Time

	Year int
	Month int
	Day int
}

func NewDate(t time.Time) Date {
	date := Date{}

	date._time = t
	date.Year = t.Year()
	date.Month = int(t.Month())
	date.Day = t.Day()

	return date
}

func (d Date) StartOfWeek() Date {
	unix := d._time.Unix()
	weekday := int64(d._time.Weekday())

	if weekday == 0 {
		return NewDate(time.Unix(unix - 24 * 60 * 60 * 6, 0))
	} else {
		return NewDate(time.Unix(unix - 24 * 60 * 60 * (weekday - 1), 0))
	}
}

func (d Date) EndOfWeek() Date {
	unix := d._time.Unix()
	weekday := int64(d._time.Weekday())

	if weekday == 0 {
		return d
	} else {
		return NewDate(time.Unix(unix + 24 * 60 * 60 * (7 - weekday), 0))
	}
}

func Today() Date {
	return NewDate(time.Now())
}

func (d Date) String() string {
	return fmt.Sprintf("%d-%s-%s", d.Year, padNumber(d.Month), padNumber(d.Day))
}

func padNumber(number int) string {
	s := fmt.Sprint(number)

	if len(s) < 2 {
		s = "0" + s
	}

	return s
}

func (d Date) Minus(days int) Date {
	out := Date{}

	t := time.Unix(d._time.Unix() - int64(days) * 24 * 60 * 60, 0)

	out._time = t
	out.Year = t.Year()
	out.Month = int(t.Month())
	out.Day = t.Day()

	return out
}

func (d Date) Plus(days int) Date {
	out := Date{}

	t := time.Unix(d._time.Unix() + int64(days) * 24 * 60 * 60, 0)

	out._time = t
	out.Year = t.Year()
	out.Month = int(t.Month())
	out.Day = t.Day()

	return out
}

func (d Date) Equal(dt *Date) bool {
	return d.Year == dt.Year && d.Month == dt.Month && d.Day == dt.Day
}

func (d Date) SQL() string {
	return fmt.Sprintf("'%s'", d.String())
}

func (d Date) CompareTo(other *Date) int {
	if d.Year > other.Year {
		return 1
	} else if d.Year < other.Year {
		return -1
	} else if d.Month > other.Month {
		return 1
	} else if d.Month < other.Month {
		return -1
	} else if d.Day > other.Day {
		return 1
	} else if d.Day < other.Day {
		return -1
	} else {
		return 0
	}
}
