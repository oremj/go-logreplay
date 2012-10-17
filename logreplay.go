package main

import "bufio"
import "flag"
import "fmt"
import "io"
import "io/ioutil"
import "net/http"
import "os"
import "runtime"
import "strings"
import "time"

var concurrent = flag.Int("c", 1, "Concurrent users")
var requests = flag.Int("n", 1, "Number of requests")
var url_file = flag.String("f", "", "File containing URLs")

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

func RunWorker(requests chan string) []int {
    durations := make([]int, 0)

    for url := range requests {
        d := fetch(url)
        durations = append(durations, d)
    }
    return durations
}

func RunTest(requests int, workers int, urls []string) (float64, []int) {

    stats := make(chan []int, workers)
    req_chan := make(chan string)

    start := time.Now().UnixNano()
    for i := 0; i < workers; i++ {
        go func() {
            stats <- RunWorker(req_chan)
        }()
    }

    for i:= 0; i < requests; i++ {
        req_chan <- urls[i % len(urls)]
    }

    close(req_chan)

    durations := make([]int, 0, requests)
    for i := 0; i < workers; i++ {
        tmp := <-stats
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

func getUrlsFromFile(file string) []string {

    f, err := os.Open(file)
    if err != nil {
        fmt.Println("Could not open file:", err)
        os.Exit(2)
    }

    urls := make([]string, 0)
    in := bufio.NewReader(f)
    for {
        l, err := in.ReadString('\n')
        if err != nil {
            break
        }
        urls = append(urls, strings.TrimSpace(l))
    }

    return urls
}

func main() {
    runtime.GOMAXPROCS(18)
    flag.Parse()

    urls := getUrlsFromFile(*url_file)

    workers := *concurrent
    requests := *requests

    test_time, durations := RunTest(requests, workers, urls)

    fmt.Printf("%0.2f req/s\n", float64(requests)/test_time)
    fmt.Printf("Avg: %dms\n", sum(durations)/requests)
}
