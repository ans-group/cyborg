# cyborg

A simple tool for testing availability of a given URI, utilising configurable concurrency

## Help

```
Usage: cyborg [options...] <url>
      --body string          Body of request
      --config string        Path to config (default "$HOME/.cyborg")
      --delay string         Delay per request (default "1s")
  -H, --header stringArray   Header(s) to use e.g. Accept: application/json
      --host string          Optional host header
  -k, --httpsskipverify      Specifies HTTPS insecure validation should be skipped
      --method string        HTTP method to use (default "GET")
      --timeout string       Specifies timeout for HTTP requests
      --workers int          Amount of workers (default 1)
```

## Environment

* `CYBORG_NOCOLOUR`: Output log messages without colour

## Usage

Example usage:

> cyborg --delay 3s --workers 3 https://www.ans.co.uk


## Output

Example output:

```
2020/06/23 14:43:47 [Worker 2] Starting worker
2020/06/23 14:43:47 [Worker 1] Starting worker
2020/06/23 14:43:47 [Worker 3] Starting worker
2020/06/23 14:43:47 [Worker 1] Got result with status code [200] in [111.9991ms]
2020/06/23 14:43:47 [Worker 2] Got result with status code [200] in [132.9937ms]
2020/06/23 14:43:47 [Worker 3] Got result with status code [200] in [109.0029ms]
2020/06/23 14:43:50 [Worker 1] Got result with status code [200] in [70.9984ms]
2020/06/23 14:43:50 [Worker 3] Got result with status code [200] in [76.9929ms]
2020/06/23 14:43:50 [Worker 2] Got result with status code [200] in [80.9952ms]
2020/06/23 14:43:51 Caught interupt signal, instructing workers to stop
2020/06/23 14:43:53 [Worker 1] Worker stopped
2020/06/23 14:43:53 [Worker 3] Worker stopped
2020/06/23 14:43:53 [Worker 2] Worker stopped
2020/06/23 14:43:53 Executed [6] requests. Successful requests: 6, Failed requests: 0, Total execution time: 6.2312894s, Minimum request time: 70.9984ms, Maximum request time: 132.9937ms
```