package utils

import "time"

const TIME_LAYOUT = "2006-01-02 15:04:05"

// GetCurrentTime get current time
func GetCurrentTime() string {
	return time.Now().Format(TIME_LAYOUT)
}

// GetTimeObject switch time string to time object
func GetTimeObject(format string) (*time.Time, error) {
	currentTime := time.Now()
	timeObject, err := time.Parse(format, currentTime.Format(format))
	if err != nil {
		return nil, err
	}
	return &timeObject, nil
}
