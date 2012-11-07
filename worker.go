package main

import (
    "fmt"
    "net/http"
    "time"
)

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

func RunWorker(requests chan string) []int {
    durations := make([]int, 0)

    for url := range requests {
        d := fetch(url)
        durations = append(durations, d)
    }
    return durations
}
