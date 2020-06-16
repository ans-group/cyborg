package main

import "time"

type Stats struct {
	RequestTime chan time.Duration
	Success     chan struct{}
	Fail        chan struct{}

	MaxRequestTime     time.Duration
	MinRequestTime     time.Duration
	SuccessfulRequests int
	FailedRequests     int
}

func NewStats() *Stats {
	return &Stats{
		RequestTime: make(chan time.Duration),
		Success:     make(chan struct{}),
		Fail:        make(chan struct{}),
	}
}

func (s *Stats) Start() {
	go func() {
		for requestTime := range s.RequestTime {
			if s.MaxRequestTime == 0 || requestTime > s.MaxRequestTime {
				s.MaxRequestTime = requestTime
			}
			if s.MinRequestTime == 0 || requestTime < s.MinRequestTime {
				s.MinRequestTime = requestTime
			}
		}
	}()

	go func() {
		for range s.Success {
			s.SuccessfulRequests++
		}
	}()
	go func() {
		for range s.Fail {
			s.FailedRequests++
		}
	}()
}
