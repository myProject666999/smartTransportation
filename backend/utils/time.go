package utils

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type Time time.Time

const (
	timeFormat1 = "2006-01-02 15:04:05"
	timeFormat2 = "2006-01-02T15:04:05Z07:00"
	timeFormat3 = "2006-01-02T15:04:05"
	timeFormat4 = "2006-01-02"
)

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		*t = Time(time.Time{})
		return nil
	}

	str := strings.Trim(string(data), `"`)

	parseFormats := []string{
		timeFormat1,
		timeFormat2,
		timeFormat3,
		timeFormat4,
		"2006/01/02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05-07:00",
	}

	var parsedTime time.Time
	var err error

	for _, format := range parseFormats {
		parsedTime, err = time.ParseInLocation(format, str, time.Local)
		if err == nil {
			*t = Time(parsedTime)
			return nil
		}
	}

	parsedTime, err = time.Parse(time.RFC3339, str)
	if err == nil {
		*t = Time(parsedTime)
		return nil
	}

	parsedTime, err = time.Parse(time.RFC3339Nano, str)
	if err == nil {
		*t = Time(parsedTime)
		return nil
	}

	return fmt.Errorf("无法解析时间: %s, 错误: %v", str, err)
}

func (t Time) MarshalJSON() ([]byte, error) {
	tm := time.Time(t)
	if tm.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, tm.Format(timeFormat1))), nil
}

func (t Time) Value() (driver.Value, error) {
	tm := time.Time(t)
	if tm.IsZero() {
		return nil, nil
	}
	return tm, nil
}

func (t *Time) Scan(value interface{}) error {
	if value == nil {
		*t = Time(time.Time{})
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*t = Time(v)
		return nil
	case string:
		parseFormats := []string{
			timeFormat1,
			timeFormat2,
			timeFormat3,
			timeFormat4,
		}
		for _, format := range parseFormats {
			parsedTime, err := time.ParseInLocation(format, v, time.Local)
			if err == nil {
				*t = Time(parsedTime)
				return nil
			}
		}
		return fmt.Errorf("无法扫描时间字符串: %s", v)
	default:
		return fmt.Errorf("不支持的时间类型: %T", v)
	}
}

func (t Time) String() string {
	tm := time.Time(t)
	if tm.IsZero() {
		return ""
	}
	return tm.Format(timeFormat1)
}

func (t Time) Time() time.Time {
	return time.Time(t)
}
