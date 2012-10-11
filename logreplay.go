package main

import "fmt"
import "net/http"
import "runtime"
import "sync"

func fetch(url string) string {
    res, err := http.Get(url)
    if err != nil {
        fmt.Println(err)
        return ""
    }

/*
    buf := make([]byte, 0, 1024 * 1024)

    m := 0
    for {
        n, _ := res.Body.Read(buf[m:m + 20])
        if n == 0 {
            break
        }
        m += n
        buf = buf[0:m]
    }
*/
    res.Body.Close()
    return ""
//    return string(buf)
}

func main() {
    runtime.GOMAXPROCS(8)
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            for j :=0; j < 200; j++ {
                fetch("https://addons-dev.allizom.org/en-US/firefox/")
            }
            wg.Done()
        }()
    }
    wg.Wait()
}
