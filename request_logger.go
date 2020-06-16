package main

import "time"

type RequestLogEntry struct {
	Worker      int           `json:"worker"`
	StatusCode  int           `json:"status_code"`
	RequestTime time.Duration `json:"request_time"`
	Error       error         `json:"error"`
}

type RequestLogger interface {
	Log(entry RequestLogEntry)
}
