package main

import "flag"
import "fmt"
import "io"
import "io/ioutil"
import "net/http"
import "runtime"
import "sync"
import "time"

var concurrent = flag.Int("c", 1, "Concurrent users")
var requests = flag.Int("n", 1, "Number of requests")
var url = flag.String("u", "", "The url to test")

func readBody(body io.Reader) string {
    b, _ := ioutil.ReadAll(body)

    return string(b)
}

func fetch(url string) int64 {
    start := time.Now().UnixNano()
    res, err := http.Get(url)
    if err != nil {
        fmt.Println(err)
        return 0
    }
    res.Body.Close()
    duration := (time.Now().UnixNano() - start) / 1e6
    return duration
}

func main() {
    runtime.GOMAXPROCS(18)
    flag.Parse()


    workers := *concurrent
    requests := *requests
    url := *url
    requests_per_worker := requests / workers
    stats := make(chan []int64, workers)

    start := time.Now().UnixNano()
    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            durations := make([]int64, 0, requests_per_worker)
            for j :=0; j < requests_per_worker; j++ {
                d := fetch(url)
                durations = append(durations, d)
            }
            wg.Done()
            stats <- durations
        }()
    }
    wg.Wait()
    test_time := (float64(time.Now().UnixNano()) - float64(start)) / 1e9

    durations := make([]int64, 0, requests)
    for i := 0; i < workers; i++ {
        tmp := <- stats
        durations = append(durations, tmp...)
    }

    sum := int64(0)
    for _, v := range durations {
        sum += v
    }

    fmt.Printf("%0.2f req/s\n", float64(requests) / test_time)
    fmt.Printf("Avg: %dms\n", sum / int64(requests))
}
