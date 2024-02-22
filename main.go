package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options...] <url>\n", os.Args[0])
		flag.PrintDefaults()
	}

	var flagMethod = flag.String("method", "GET", "HTTP method to use")
	var flagHeaders = flag.StringArrayP("header", "H", []string{}, "Header(s) to use e.g. Accept: application/json")
	var flagHost = flag.String("host", "", "Optional host header")
	var flagWorkers = flag.Int("workers", 1, "Amount of workers")
	var flagDelay = flag.String("delay", "1s", "Delay per request")
	var flagBody = flag.String("body", "", "Body of request")
	var flagHTTPSSkipVerify = flag.BoolP("httpsskipverify", "k", false, "Specifies HTTPS insecure validation should be skipped")
	var flagTimeout = flag.String("timeout", "", "Specifies timeout for HTTP requests")
	var flagNoColour = flag.Bool("no-colour", false, "Disables coloured output")
	var flagConfig = flag.String("config", "$HOME/.cyborg", "Path to config")

	flag.Parse()

	err := initConfig(*flagConfig)
	if err != nil {
		log.Fatalf("failed to init config: %s", err)
	}

	if *flagHTTPSSkipVerify {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &http.Client{}

	url := flag.Arg(0)
	if url == "" {
		log.Fatal("must provide URL")
	}

	parsedDelayDuration, err := parseDurationString(flagDelay)
	if err != nil {
		log.Fatal(err.Error())
	}

	parsedTimeoutDuration, err := parseDurationString(flagTimeout)
	if err != nil {
		log.Fatal(err.Error())
	}

	if parsedTimeoutDuration != nil {
		client.Timeout = *parsedTimeoutDuration
	}

	parsedHeaders := parseHeadersFlag(flagHeaders)

	stats := NewStats()
	stats.Start()

	logManager := NewRequestLoggerManager()
	logManager.AddLogger(&RequestLoggerStdout{NoColour: *flagNoColour})

	if viper.GetBool("logger.elk.enabled") {
		elkLogger, err := NewRequestLoggerELK()
		if err != nil {
			log.Fatalf("failed to initialise ELK logger: %s", err)
		}
		logManager.AddLogger(elkLogger)
	}

	workers := make([]*Worker, *flagWorkers)
	for i := 0; i < *flagWorkers; i++ {
		workerNum := i + 1
		worker := Worker{
			Number:               workerNum,
			Delay:                parsedDelayDuration,
			RequestLoggerManager: logManager,
			Stats:                stats,
			Client:               client,
			Request: func() (*http.Request, error) {
				bodyBuf := new(bytes.Buffer)
				if *flagBody != "" {
					bodyBuf.WriteString(*flagBody)
				}
				req, err := http.NewRequest(*flagMethod, url, bodyBuf)
				if err != nil {
					return nil, fmt.Errorf("failed to create request: %s", err)
				}

				if *flagHost != "" {
					req.Host = *flagHost
				}

				req.Header = parsedHeaders

				return req, nil
			},
		}

		workers[i] = &worker
	}

	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	start := time.Now()
	for _, worker := range workers {
		w := worker
		g.Go(func() error {
			return w.Start(ctx)
		})
	}

	flushResults := func() {
		elapsed := time.Since(start)
		log.Printf("Executed [%d] requests. Successful requests: %d, Failed requests: %d, Total execution time: %s, Minimum request time: %s, Maximum request time: %s", stats.SuccessfulRequests+stats.FailedRequests, stats.SuccessfulRequests, stats.FailedRequests, elapsed, stats.MinRequestTime, stats.MaxRequestTime)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Printf("Caught interupt signal, instructing workers to stop")
			cancel()
		}
	}()

	g.Wait()
	flushResults()
}

func parseDurationString(delay *string) (*time.Duration, error) {
	if *delay == "" {
		return nil, nil
	}

	d, err := time.ParseDuration(*delay)
	if err != nil {
		return nil, errors.New("invalid delay")
	}

	return &d, nil
}

func parseHeadersFlag(flagHeaders *[]string) http.Header {
	parsedHeaders := make(map[string][]string)

	for _, header := range *flagHeaders {
		headerSplit := strings.Split(header, ":")
		parsedHeaders[strings.TrimSpace(headerSplit[0])] = []string{strings.TrimSpace(headerSplit[1])}
	}

	return parsedHeaders
}
