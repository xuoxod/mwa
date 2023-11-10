package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func DateTimeStamp() string {
	// dts := fmt.Sprint("Date: ", time.Now())

	// d := time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC)

	d := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 12, 30, 0, 0, time.UTC)
	year, month, day := d.Date()

	return fmt.Sprintf("%v/%v/%v", month, day, year)
}
func DateStamp() string {
	// dts := fmt.Sprint("Date: ", time.Now())

	// d := time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC)

	d := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 12, 30, 0, 0, time.UTC)
	year, month, day := d.Date()

	return fmt.Sprintf("%v %v %v", month, day, year)
}

func DTS() string {
	// dts := fmt.Sprint("Date: ", time.Now())
	// d := time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC)

	d := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 12, 30, 0, 0, time.UTC)
	year, month, day := d.Date()
	hour, minute, second := time.Now().Local().Clock()

	var suffix string
	var strDay string = fmt.Sprintf("%d", day)

	if strings.HasSuffix(strDay, "1") {
		suffix = "st"
	} else if strings.HasSuffix(strDay, "2") {
		suffix = "nd"
	} else if strings.HasSuffix(strDay, "3") {
		suffix = "rd"
	} else {
		suffix = "th"
	}

	return fmt.Sprintf("%v %v%s %v %v:%v:%v", month, day, suffix, year, hour, minute, second)
}

func DS() string {
	d := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 12, 30, 0, 0, time.UTC)
	year, month, day := d.Date()
	var suffix string
	var strDay string = fmt.Sprintf("%d", day)

	if strings.HasSuffix(strDay, "1") {
		suffix = "st"
	} else if strings.HasSuffix(strDay, "2") {
		suffix = "nd"
	} else if strings.HasSuffix(strDay, "3") {
		suffix = "rd"
	} else {
		suffix = "th"
	}

	return fmt.Sprintf("%v %v%s %v", month, day, suffix, year)
}

func TS() string {
	hour, minute, second := time.Now().Local().Clock()
	return fmt.Sprintf("%v:%v:%v", hour, minute, second)
}

func Print(msg string) {
	fmt.Println(msg)
}

func GenerateRandomNumber() (int, error) {
	min := 111111
	max := 999999
	return min + rand.Intn(max-min), nil
}

func CopyrightDate() string {
	d := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 12, 30, 0, 0, time.UTC)
	year, _, _ := d.Date()

	return fmt.Sprintf("RmediaTech %v", year)
}
