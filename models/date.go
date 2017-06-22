package models

import (
	"fmt"
	"time"
	"regexp"
	"strconv"
)

type Date struct {
	_time time.Time

	Year  int
	Month int
	Day   int
}

type InvalidDate struct {}

func (id InvalidDate) Error() string {
	return "the date is invalid"
}

const dateRegexp = "^(\\d\\d\\d\\d)-(\\d\\d)-(\\d\\d)$"

func NewDate(t time.Time) Date {
	date := Date{}

	date._time = t
	date.Year = t.Year()
	date.Month = int(t.Month())
	date.Day = t.Day()

	return date
}

func DateFromString(s string) (*Date, error) {
	matched, err := regexp.MatchString(dateRegexp, s)

	if err != nil {
		return nil, err
	} else if !matched {
		return nil, InvalidDate{}
	}

	r := regexp.MustCompile(dateRegexp)

	matchData := r.FindStringSubmatch(s)

	year, err := strconv.ParseInt(matchData[1], 10, 64)
	month, err := strconv.ParseInt(matchData[2], 10, 64)
	day, err := strconv.ParseInt(matchData[3], 10, 64)

	if err != nil {
		return nil, err
	}

	location, err := time.LoadLocation(time.Now().Location().String())

	if err != nil {
		return nil, err
	}

	out := NewDate(time.Date(int(year), time.Month(month), int(day), 12, 0, 0, 0, location))

	return &out, nil
}

func (d Date) StartOfWeek() Date {
	unix := d._time.Unix()
	weekday := int64(d._time.Weekday())

	if weekday == 0 {
		return NewDate(time.Unix(unix-24*60*60*6, 0))
	} else {
		return NewDate(time.Unix(unix-24*60*60*(weekday-1), 0))
	}
}

func (d Date) EndOfWeek() Date {
	unix := d._time.Unix()
	weekday := int64(d._time.Weekday())

	if weekday == 0 {
		return d
	} else {
		return NewDate(time.Unix(unix+24*60*60*(7-weekday), 0))
	}
}

func (d Date) Format(s string) string {
	return d._time.Format(s)
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

	t := time.Unix(d._time.Unix()-int64(days)*24*60*60, 0)

	out._time = t
	out.Year = t.Year()
	out.Month = int(t.Month())
	out.Day = t.Day()

	return out
}

func (d Date) Plus(days int) Date {
	out := Date{}

	t := time.Unix(d._time.Unix()+int64(days)*24*60*60, 0)

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

func BeginningOfMonth() Date {
	t := time.Now()

	return NewDate(time.Date(t.Year(), t.Month(), 1, 12, 0, 0, 0, t.Location()))
}

func EndOfMonth() Date {
	t := time.Now()

	return NewDate(time.Date(t.Year(), t.Month(), endOfMonth(int(t.Month()), t.Year()), 12, 0, 0, 0, t.Location()))
}

func BeginningOfLastMonth() Date {
	t := time.Now()

	return NewDate(time.Date(t.Year(), t.Month() - 1, 1, 12, 0, 0, 0, t.Location()))
}

func EndOfLastMonth() Date {
	t := time.Now()

	return NewDate(time.Date(t.Year(), t.Month() - 1, endOfMonth(int(t.Month() - 1), t.Year()), 12, 0, 0, 0, t.Location()))
}


func endOfMonth(month, year int) int {
	m := map[int]int {
		1: 31,
		3: 31,
		4: 30,
		5: 31,
		6: 30,
		7: 31,
		8: 31,
		9: 30,
		10: 31,
		11: 30,
		12: 31,
	}

	if month == 2 && year % 4 == 0 {
		return 29
	} else if month == 2 {
		return 28
	} else {
		return m[month]
	}
}
