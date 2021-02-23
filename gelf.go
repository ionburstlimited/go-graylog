package graylog

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

// GELF wrapper for Graylog instance
type GELF struct {
	Graylog *Graylog
	Source  string
	Logger  string
}

// GELFConfig definitions for Graylog host details
type GELFConfig struct {
	GraylogServer string
	GraylogPort   uint
	GraylogSource string
}

// Logger instance for GELF
var Logger *GELF

// CreateLogger - create instance of GELF
func CreateLogger(settings GELFConfig) *GELF {
	Logger = &GELF{}
	// Initialize a new graylog client with TLS
	g, err := NewGraylogTLS(Endpoint{
		Transport: TCP,
		Address:   settings.GraylogServer,
		Port:      settings.GraylogPort,
	}, 3*time.Second, nil)
	if err != nil {
		fmt.Println("Failed to initialise Graylog client:", err.Error())
		panic(err)
	} else {
		Logger.Graylog = g
		Logger.Source = settings.GraylogSource
		Logger.Logger, err = os.Executable()
		if err != nil {
			fmt.Println("Error getting executable name:", err.Error())
		} else {
			return Logger
		}
	}

	return nil
}

// LogFatal to provide .NET ILogger style function
func (l *GELF) LogFatal(shortMessage string, message string) {
	l.Log(shortMessage, message, 1)
}

// LogCritical to provide .NET ILogger style function
func (l *GELF) LogCritical(shortMessage string, message string) {
	l.Log(shortMessage, message, 2)
}

// LogError to provide .NET ILogger style function
func (l *GELF) LogError(shortMessage string, message string) {
	l.Log(shortMessage, message, 3)
}

// LogWarning to provide .NET ILogger style function
func (l *GELF) LogWarning(shortMessage string, message string) {
	l.Log(shortMessage, message, 4)
}

// LogInformation to provide .NET ILogger style function
func (l *GELF) LogInformation(shortMessage string, message string) {
	l.Log(shortMessage, message, 5)
}

// LogDebug to provide .NET ILogger style function
func (l *GELF) LogDebug(shortMessage string, message string) {
	l.Log(shortMessage, message, 6)
}

// LogTrace to provide .NET ILogger style function
func (l *GELF) LogTrace(shortMessage string, message string) {
	l.Log(shortMessage, message, 7)
}

// Log generic log to Graylog function
func (l *GELF) Log(shortMessage string, message string, severity uint) {
	if shortMessage == "" {
		if len(message) <= 128 {
			shortMessage = message
		} else {
			shortMessage = message[0:127]
		}
	}
	err := l.Graylog.Send(Message{
		Version:      "1.1",
		Host:         l.Source,
		ShortMessage: shortMessage,
		FullMessage:  message,
		Timestamp:    MakeTimestamp(),
		Level:        severity,
		Extra: map[string]string{
			"logger":     path.Base(l.Logger),
			"executable": l.Logger}})
	if err != nil {
		fmt.Println("Graylog send failed:", err.Error())
		panic(err)
	}
}

func MakeTimestamp() float64 {
	var millisecInt int64 = time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
	var secF float64 = float64(millisecInt) / float64(1000)
	var secString = fmt.Sprintf("%.3f", secF)
	result, err := strconv.ParseFloat(secString, 64)
	if err != nil {
		fmt.Println(fmt.Sprintf("ParseFloat error: %s", err.Error()))
	}

	return result
}
