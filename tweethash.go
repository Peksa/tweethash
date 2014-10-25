package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type Attempt struct {
	id   int64
	hash []byte
}

func main() {
	workers := 128
	runtime.GOMAXPROCS(workers)

	rand.Seed(time.Now().UnixNano())

	bestAttempt := Attempt{0, make([]byte, 32)}
	bestAttempt.hash[0] = 0xff

	results := make(chan Attempt, 5000000)
	counter := make(chan int, 5000000)

	speed := 0.0

	workerStarts := make([]int64, workers)

	for w := 0; w < workers; w++ {
		workerStarts[w] = rand.Int63n(131621703842267136)
		go startWorker(results, counter, workerStarts[w])
	}
	start := time.Now()
	total := 0

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "Speed: %.0f\nLowest: %s %x\nStarts: %d", speed, strconv.FormatInt(bestAttempt.id, 36), bestAttempt.hash, workerStarts)
	})

	go http.ListenAndServe(":8080", nil)

	for {
		select {
		case attempt := <-results:
			if bytes.Compare(attempt.hash, bestAttempt.hash) == -1 {
				fmt.Printf("%s: %x\n", strconv.FormatInt(attempt.id, 36), attempt.hash)
				bestAttempt = attempt
			}
		default:
		}

		select {
		case count := <-counter:
			total += count
			if total > 10000000 {
				millis := time.Now().Sub(start) / time.Millisecond
				speed = float64(total) / float64(millis) * 1000
				fmt.Printf("Speed: %.2fM/s\n", speed/1000000)
				start = time.Now()
				total = 0
			}
		default:
		}
	}
}

func startWorker(results chan<- Attempt, counter chan<- int, currentValue int64) {
	bestAttempt := make([]byte, 32)
	bestAttempt[0] = 0xff
	count := 0

	prefix := []byte{'h', 't', 't', 'p', 's', ':', '/', '/', 't', 'w', 'i', 't', 't', 'e', 'r', '.', 'c', 'o', 'm', '/', 'p', 'e', 'k', '_'}
	suffix := []byte{'/', 's', 't', 'a', 't', 'u', 's', '/', '5', '2', '5', '6', '4', '4', '1', '4', '0', '8', '6', '5', '1', '4', '2', '7', '8', '4'}

	for {
		h := sha256.New()
		h.Write(prefix)
		h.Write([]byte(strconv.FormatInt(currentValue, 36)))
		h.Write(suffix)
		hash := h.Sum(nil)

		if bytes.Compare(hash, bestAttempt) == -1 {
			results <- Attempt{currentValue, hash}
			bestAttempt = hash
		}

		currentValue++
		count++
		if count > 10000 {
			counter <- count
			count = 0
		}
	}
}
