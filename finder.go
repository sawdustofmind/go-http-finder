package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

// N determines the number of workers
const N = 5

func main() {
	reader := bufio.NewReader(os.Stdin)
	var total int32
	wg := sync.WaitGroup{}
	doneCh := make(chan struct{}, N-1)
	for {
		url, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		url = strings.TrimRight(url, "\n")
		if url == "" {
			continue
		}

		wg.Add(1)
		go func() {
			defer func() {
				<-doneCh
				wg.Done()
			}()
			count, err := countMathesInPage(url)
			if err != nil {
				fmt.Printf("Error while process %s: %s\n", url, err)
				return
			}
			fmt.Printf("Count for %s: %d\n", url, count)
			atomic.AddInt32(&total, int32(count))
		}()
		doneCh <- struct{}{}
	}
	wg.Wait()
	fmt.Println("Total:", total)
}

func countMathesInPage(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	return strings.Count(string(payload), "Go"), nil
}
