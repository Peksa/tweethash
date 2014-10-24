package main

import (
	"crypto/sha256"
	"fmt"
	"bytes"
	"runtime"
	"time"
	"math/rand"
	"net/http"
)

type Attempt struct {
	id int64
	hash []byte
}

func main() {
	cpus := 8

	runtime.GOMAXPROCS(cpus)

	bestAttempt := Attempt{0, make([]byte, 32)}
	bestAttempt.hash[0] = 0xff

	results := make(chan Attempt, 5000000)
	counter := make(chan int, 5000000)

	speed := 0.0

	starts := make([]int64, cpus)

	for w := 0; w < 8; w++ {
		starts[w] = rand.Int63()
		go getHash(results, counter, starts[w])
	}

	start := time.Now()
	total := 0

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Speed: %.0f\nLowest: %x %x\n", speed, bestAttempt.id, bestAttempt.hash)
	})

	go http.ListenAndServe(":8080", nil)


	for {
		select {
		case attempt := <- results:
			if bytes.Compare(attempt.hash, bestAttempt.hash) == -1 {
				fmt.Printf("%x: %x\n", attempt.id, attempt.hash)
				bestAttempt = attempt
			}
		default:
		}

		select {
		case count := <- counter:
			total += count
			if total > 10000000 {
				millis := int64(time.Now().Sub(start)/time.Millisecond)
				speed = float64(total)/float64(millis)*1000
				fmt.Printf("Speed: %.2fM/s\n", speed/1000000)
				start = time.Now()
				total = 0
			}
		default:
		}
	}
}

func getHash(results chan<- Attempt, counter chan<- int, start int64) {
	bestAttempt := make([]byte, 32)
	bestAttempt[0] = 0xff
	count := 0
	for {
		h := sha256.New()
		h.Write([]byte("https://twitter.com/p"))
		h.Write([]byte(fmt.Sprintf("%x", start)))
		h.Write([]byte("/status/525644140865142784"))
		hash := h.Sum(nil)

		if bytes.Compare(hash, bestAttempt) == -1 {
			results <- Attempt{start, hash}
			bestAttempt = hash
		}

		start++
		count++
		if count > 500000 {
			counter <- count
			count = 0
		}
	}
}
