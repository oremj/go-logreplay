package main

import "fmt"
import "io"
import "io/ioutil"
import "net/http"
import "runtime"
import "sync"
import "time"


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
    readBody(res.Body)
    res.Body.Close()
    duration := (time.Now().UnixNano() - start) / 1e6
    return duration
}

func main() {
    runtime.GOMAXPROCS(8)

    workers := 20
    requests := 200
    requests_per_worker := requests / workers
    stats := make(chan []int64, workers)

    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            durations := make([]int64, 0, requests_per_worker)
            for j :=0; j < requests_per_worker; j++ {
                d := fetch("https://addons-dev.allizom.org/en-US/firefox/")
                durations = append(durations, d)
            }
            wg.Done()
            stats <- durations
        }()
    }
    wg.Wait()

    durations := make([]int64, 0, requests)
    for i := 0; i < workers; i++ {
        tmp := <- stats
        durations = append(durations, tmp...)
    }

    sum := int64(0)
    for _, v := range durations {
        sum += v
    }

    fmt.Printf("Avg: %dms\n", sum / int64(requests))
}
