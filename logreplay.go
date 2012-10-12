package main

import "flag"
import "fmt"
import "io"
import "io/ioutil"
import "net/http"
import "runtime"
import "time"

var concurrent = flag.Int("c", 1, "Concurrent users")
var requests = flag.Int("n", 1, "Number of requests")
var url = flag.String("u", "", "The url to test")

func readBody(body io.Reader) string {
    b, _ := ioutil.ReadAll(body)

    return string(b)
}

func fetch(url string) int {
    start := time.Now().UnixNano()
    res, err := http.Get(url)
    if err != nil {
        fmt.Println(err)
        return 0
    }
    res.Body.Close()
    duration := (time.Now().UnixNano() - start) / 1e6
    return int(duration)
}

func RunWorker(requests int, url string) []int {
    durations := make([]int, 0, requests)

    for i :=0; i < requests; i++ {
        d := fetch(url)
        durations = append(durations, d)
    }
    return durations
}

func RunTest(requests int, workers int, url string) (float64, []int) {
    requests_per_worker := requests / workers

    stats := make(chan []int, workers)

    start := time.Now().UnixNano()
    for i := 0; i < workers; i++ {
        go func() {
            stats <- RunWorker(requests_per_worker, url)
        }()
    }

    durations := make([]int, 0, requests)
    for i := 0; i < workers; i++ {
        tmp := <- stats
        durations = append(durations, tmp...)
    }

    test_time := (float64(time.Now().UnixNano()) - float64(start)) / 1e9

    return test_time, durations
}

func sum(a []int) int {
    s := 0
    for _, v := range a {
        s += v
    }

    return s
}

func main() {
    runtime.GOMAXPROCS(18)
    flag.Parse()

    workers := *concurrent
    requests := *requests
    url := *url

    test_time, durations := RunTest(requests, workers, url)

    fmt.Printf("%0.2f req/s\n", float64(requests) / test_time)
    fmt.Printf("Avg: %dms\n", sum(durations) / requests)
}
