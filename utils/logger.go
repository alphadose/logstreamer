package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// tag definitions
const (
	ErrorTAG = iota
	InfoTAG
	DebugTAG
)

var (
	tagToString = map[int]string{
		ErrorTAG: "[ERROR]",
		InfoTAG:  "[INFO]",
		DebugTAG: "[DEBUG]",
	}
	logfile       *os.File
	outFileStream = make(chan string, 1<<8)
)

func getTimeStamp() string {
	return fmt.Sprintf("%v", time.Now().Unix())
}

func out(context, message string, TAG int) {
	currentTime := time.Now()
	hours := fmt.Sprintf("%d", currentTime.Hour())
	if currentTime.Hour() < 10 {
		hours = "0" + hours
	}
	minutes := fmt.Sprintf("%d", currentTime.Minute())
	if currentTime.Minute() < 10 {
		minutes = "0" + minutes
	}
	seconds := fmt.Sprintf("%d", currentTime.Second())
	if currentTime.Second() < 10 {
		seconds = "0" + seconds
	}
	timeLog := fmt.Sprintf(
		"%d-%d-%d %s:%s:%s",
		currentTime.Day(),
		currentTime.Month(),
		currentTime.Year(),
		hours,
		minutes,
		seconds,
	)
	logOut := fmt.Sprintf("%s %s -> %s\n", tagToString[TAG], timeLog, message)
	print(logOut)
	outFileStream <- logOut
}

// Log logs to the console with your custom TAG
func Log(context, message string, TAG int) {
	out(context, message, TAG)
}

// LogInfo logs information to the console
func LogInfo(context, template string, args ...interface{}) {
	out(context, fmt.Sprintf(template, args...), InfoTAG)
}

// LogDebug logs debug messages to console
func LogDebug(context, template string, args ...interface{}) {
	out(context, fmt.Sprintf(template, args...), DebugTAG)
}

// LogError logs type error to console
func LogError(context string, err error) {
	if err == nil {
		return
	}
	out(context, err.Error(), ErrorTAG)
}

// GracefulExit exits gracefully without panicking
func GracefulExit(context string, err error) {
	LogError(context, err)
	os.Exit(1)
}

func init() {
	_ = os.MkdirAll("logs", 0755)
	logfile, _ = os.Create(filepath.Join("logs", filepath.Base(fmt.Sprintf("app-%d-%s.log", os.Getpid(), getTimeStamp()))))
	go func() {
		for {
			if _, err := logfile.WriteString(<-outFileStream); err != nil {
				println(err.Error())
			}
		}
	}()
}
