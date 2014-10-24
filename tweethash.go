package main

import (
	"crypto/sha256"
	"fmt"
	"bytes"
	"runtime"
	"time"
	"math/rand"
)

type Attempt struct {
	id int64
	hash []byte
}

func main() {
	cpus := 8

	runtime.GOMAXPROCS(cpus)

	bestAttempt := make([]byte, 32)
	bestAttempt[0] = 0xff

	results := make(chan Attempt, 5000000)
	counter := make(chan int, 5000000)

	starts := make([]int, cpus)

	for w := 0; w < 8; w++ {
		start[w] = rand.Int63()
		go getHash(results, counter, start[w])
	}

	start := time.Now()
	total := 0

	for {
		select {
		case attempt := <- results:
			if bytes.Compare(attempt.hash, bestAttempt) == -1 {
				fmt.Printf("%x: %x\n", attempt.id, attempt.hash)
				bestAttempt = attempt.hash
			}
		default:
		}

		select {
		case count := <- counter:
			total += count
			if total > 10000000 {
				millis := int64(time.Now().Sub(start)/time.Millisecond)
				fmt.Printf("Speed: %.2fM/s\n", float64(total)/float64(millis)/1000)
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
