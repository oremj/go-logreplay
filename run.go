package main

import "time"

func sendUrls(urls []string, requests int, req_chan chan string) {
    for i := 0; i < requests; i++ {
        req_chan <- urls[i%len(urls)]
    }
    close(req_chan)
}

func startWorkers(workers int) (stats chan []int, req_chan chan string) {
    stats = make(chan []int, workers)
    req_chan = make(chan string)

    for i := 0; i < workers; i++ {
        go func() {
            stats <- RunWorker(req_chan)
        }()
    }

    return
}

func recvStats(stats chan []int, workers int, requests int) (durations []int) {
    durations = make([]int, 0, requests)
    for i := 0; i < workers; i++ {
        tmp := <-stats
        durations = append(durations, tmp...)
    }

    return
}

func RunTest(requests int, workers int, urls []string) (float64, []int) {

    start := time.Now().UnixNano()

    stats, req_chan := startWorkers(workers)
    sendUrls(urls, requests, req_chan)
    durations := recvStats(stats, workers, requests)

    test_time := (float64(time.Now().UnixNano()) - float64(start)) / 1e9

    return test_time, durations
}

