package main

import "fmt"

type Logger interface {
	Log() string
}

type LogLevel string

const (
	Info  LogLevel = "Info"
	Error LogLevel = "Error"
)

type log struct {
	LogLevel string
}

func (l log) log(messge string) {
	if l.LogLevel = "Info" {
		fmt.Println("INFO: " + messge)
	} else {
		fmt.Println("ERROR: " + messge)
	}
}
