package utils

import (
	"fmt"
	"time"
)

func DateTimeStamp() string {
	// dts := fmt.Sprint("Date: ", time.Now())

	// d := time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC)

	d := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 12, 30, 0, 0, time.UTC)
	year, month, day := d.Date()

	return fmt.Sprintf("%v/%v/%v", month, day, year)
}

func Print(msg string) {
	fmt.Println(msg)
}
