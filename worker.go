package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Worker struct {
	Number               int
	Delay                *time.Duration
	Request              func() (*http.Request, error)
	RequestLoggerManager *RequestLoggerManager
	Stats                *Stats
	Client               *http.Client
}

func (w *Worker) Start(ctx context.Context) error {
	w.log("Starting worker")

	for {
		select {
		case <-ctx.Done():
			w.log("Worker stopped")
			return ctx.Err()
		default:
			w.work()
			if w.Delay != nil {
				time.Sleep(*w.Delay)
			}
		}
	}
}

func (w *Worker) work() {
	reqStart := time.Now()

	req, err := w.Request()
	if err != nil {
		w.failure(fmt.Errorf("Failed to create request: %s", err))
		return
	}

	r, err := w.Client.Do(req)
	if err != nil {
		w.failure(fmt.Errorf("Failed to invoke request: %s", err))
		return
	}
	if r == nil {
		w.failure(errors.New("No response"))
		return
	}

	if r.Body != nil {
		r.Body.Close()
	}

	requestTime := time.Since(reqStart)
	w.Stats.RequestTime <- requestTime
	w.success("Request successful", r.StatusCode, requestTime)
}

func (w *Worker) success(msg string, statusCode int, requestTime time.Duration) {
	w.Stats.Success <- struct{}{}
	w.RequestLoggerManager.Log(RequestLogEntry{
		StatusCode:  statusCode,
		RequestTime: requestTime,
		Worker:      w.Number,
	})
}

func (w *Worker) failure(err error) {
	w.Stats.Fail <- struct{}{}
	w.RequestLoggerManager.Log(RequestLogEntry{
		Error:  err,
		Worker: w.Number,
	})
}
func (w *Worker) log(msg string) {
	log.Printf("[Worker %d] %s", w.Number, msg)
}
