package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func runGoroutinesRace(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	wg := &sync.WaitGroup{}
	switch r.Method {
	case "GET":
		channel := make(chan string)
		ctx, cancel := context.WithCancel(context.Background())
		for i := 1; i <= 10; i++ {
			wg.Add(1)
			go getMeString(i, wg, channel, ctx)
		}

		value := <-channel
		fmt.Printf(value)
		cancel()
		wg.Wait()
		close(channel)
		_, err := w.Write([]byte(value))
		if err != nil {
			return
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET method is supported.")
	}
	fmt.Println("This is the end of main func!")
}

func main() {
	http.HandleFunc("/", runGoroutinesRace)

	fmt.Printf("Starting server for testing HTTP...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func getMeString(i int, wg *sync.WaitGroup, channel chan string, ctx context.Context) {
	ticker := time.Tick(time.Duration(i) * time.Second)
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case <-ticker:
			channel <- fmt.Sprintf("I'am %d goroutine\n", i)
		}
	}
}
