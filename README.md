# cyborg

A simple tool for testing availability of a given URI, utilising configurable concurrency

## Help

```
      --delay string         Optional delay per request
  -H, --header stringArray   Header(s) to use e.g. Accept: application/json
      --host string          Optional host header
  -k, --httpsskipverify      Specifies HTTPS insecure validation should be skipped
      --jsonbody string      JSON body of request
      --method string        HTTP method to use (default "GET")
      --uri string           URI to hit
      --workers int          Amount of workers (default 1)
```

## Environment

* `CYBORG_COLOUR`: Output log messages with colour

## Usage

Example usage:

> cyborg --uri https://lee.io --delay 3s --workers 3


## Output

Example output:

```
2019/07/03 10:39:48 Starting worker [1]
2019/07/03 10:39:48 Starting worker [2]
2019/07/03 10:39:48 Starting worker [3]
2019/07/03 10:39:48 [Worker 1] Got result with status code [200] in [133.0969ms]
2019/07/03 10:39:48 [Worker 2] Got result with status code [200] in [132.08ms]
2019/07/03 10:39:48 [Worker 3] Got result with status code [200] in [132.0698ms]
2019/07/03 10:39:51 [Worker 2] Got result with status code [200] in [60.2247ms]
2019/07/03 10:39:51 [Worker 3] Got result with status code [200] in [58.2248ms]
2019/07/03 10:39:51 [Worker 1] Got result with status code [200] in [62.2146ms]
2019/07/03 10:39:52 Caught interupt signal, instructing workers to stop
2019/07/03 10:39:52 [Worker 1] Instructing worker to stop..
2019/07/03 10:39:52 [Worker 2] Instructing worker to stop..
2019/07/03 10:39:52 [Worker 3] Instructing worker to stop..
2019/07/03 10:39:54 [Worker 3] Got result with status code [200] in [59.4172ms]
2019/07/03 10:39:54 [Worker 3] Worker stopped
2019/07/03 10:39:54 [Worker 1] Got result with status code [200] in [57.4239ms]
2019/07/03 10:39:54 [Worker 2] Got result with status code [200] in [61.4343ms]
2019/07/03 10:39:54 [Worker 1] Worker stopped
2019/07/03 10:39:54 [Worker 2] Worker stopped
2019/07/03 10:39:54 Executed [9] requests. Successful requests: 9, Failed requests: 0, Total execution time: 6.258393s, Minimum request time: 57.4239ms, Maximum request time: 133.0969ms
```