package main


import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "runtime"
    "strings"
)

var concurrent = flag.Int("c", 1, "Concurrent users")
var requests = flag.Int("n", 1, "Number of requests")
var url_file = flag.String("f", "", "File containing URLs")

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
    runtime.GOMAXPROCS(runtime.NumCPU() - 2)
    flag.Parse()

    urls := getUrlsFromFile(*url_file)

    workers := *concurrent
    requests := *requests

    test_time, durations := RunTest(requests, workers, urls)

    fmt.Printf("%0.2f req/s\n", float64(requests)/test_time)
    fmt.Printf("Avg: %dms\n", sum(durations)/requests)
}
