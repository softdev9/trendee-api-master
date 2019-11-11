package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Describe the number of max day in a month
var DayMaxInMonth = map[int]int{
	1:  31,
	2:  28,
	3:  31,
	4:  30,
	5:  31,
	6:  30,
	7:  31,
	8:  31,
	9:  30,
	10: 31,
	11: 30,
	12: 31,
}

/*
Chech the date for the given format DD-MM-YYYY
*/
func CheckDate(date string) error {
	components := strings.Split(date, "-")
	if len(components) != 3 {
		return errors.New("Date format DD-MM-YYYY")
	}
	// Check the year
	if err := checkYear(components[2]); err != nil {
		return err
	}
	if err := checkMonth(components[1]); err != nil {
		return err
	}
	y, _ := strconv.Atoi(components[2])
	m, _ := strconv.Atoi(components[1])
	if err := checkDay(components[0], m, y); err != nil {
		return err
	}
	return nil
}

func checkYear(year string) error {
	if len(year) != 4 {
		return errors.New("Date format DD-MM-YYYY")
	}
	y, err := strconv.Atoi(year)
	if err != nil {
		return errors.New("Invalid Year")
	}
	if y < 1900 && y > time.Now().UTC().Year() {
		return errors.New("Invalid Year")
	}
	return nil
}

func checkMonth(month string) error {
	m, err := strconv.Atoi(month)
	if err != nil {
		return err
	}
	if m < 1 && m > 12 {
		return errors.New("Invalid Month")
	}
	return nil
}

func checkDay(day string, month int, year int) error {
	d, err := strconv.Atoi(day)
	if err != nil {
		return errors.New("Invalid day given L02")
	}
	if d < 1 {
		return errors.New("Invalid day given L105")
	}
	if month == 2 && year%4 == 0 {
		if d > 29 {
			return errors.New("Invalid day given L109")
		}
		return nil
	}
	if d > DayMaxInMonth[month] {
		return errors.New("Invalid day given L113")
	}
	return nil
}

func ParseDate(date string, separator string) time.Time {
	components := strings.Split(date, separator)
	if len(components) == 3 {
		day, _ := strconv.Atoi(components[0])
		month, _ := strconv.Atoi(components[1])
		year, _ := strconv.Atoi(components[2])
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	}
	if len(components) == 1 && len(components[0]) == 4 {
		year, _ := strconv.Atoi(components[0])
		return time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	}
	if len(components) == 2 {
		day, _ := strconv.Atoi(components[1])
		month, _ := strconv.Atoi(components[0])
		return time.Date(1, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	}
	return time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
}

func DateToString(time time.Time) string {
	return fmt.Sprintf("%d-%d-%d", time.Day(), time.Month(), time.Year())
}
