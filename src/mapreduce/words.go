package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

const DataFile = "loremipsum.txt"

// Return the word frequencies of the text argument.
//
// Split load optimally across processor cores.
func WordCount(text string) map[string]int {
	freqs := make(map[string]int)
	results := make(chan map[string]int)
	// wg := new(sync.WaitGroup)

	numChunks := 10
	chunks := make([]string, numChunks)
	chunkLength := len(text) / numChunks
	next := 0
	for i := 0; i < numChunks; i += 1 {
		start := next
		end := (i + 1) * chunkLength

		// if this is the first chunk we'll include the whole thing
		if i == 0 {
			start = 0
		} else {
			if text[start] != ' ' {
				// find the next space
				spaceIndex := strings.Index(text[start:], " ")

				if spaceIndex != -1 {
					start = start + spaceIndex
				}
			}
		}

		// if this is the last chunk we'll include the whole thing
		if i == numChunks-1 {
			end = len(text)
		} else {
			if text[end] != ' ' {
				// find the next space
				spaceIndex := strings.Index(text[end:], " ")

				if spaceIndex != -1 {
					end = end + spaceIndex
				}
			}
		}

		chunks[i] = text[start:end]

		// next should start at the end of this chunk
		next = end
	}

	// now, for each chunk, count the word frequencies
	// merge the results
	for _, chunk := range chunks {
		go countChunk(chunk, results)
	}

	// receive the results from the channel
	// merge the results into freqs
	for i := 0; i < numChunks; i += 1 {
		data := <-results
		for k, v := range data {
			freqs[k] = freqs[k] + v
		}

	}

	// close the channel
	close(results)

	// return the results
	return freqs
}

// Count the word frequencies in a chunk of text.
func countChunk(text string, send chan map[string]int) {
	freqs := make(map[string]int)

	// iterate over the words in text
	// for each word, increment the count in freqs
	for _, word := range strings.Fields(text) {
		// remove any punctuation from the word, trailing or leading
		word = strings.Trim(word, ".,;:!?\"'")

		// convert the word to lowercase, and increment the count in freqs
		freqs[strings.ToLower(word)]++
	}

	// send the results back to the main thread
	send <- freqs

}

// Benchmark how long it takes to count word frequencies in text numRuns times.
//
// Return the total time elapsed.
func benchmark(text string, numRuns int) int64 {
	start := time.Now()
	for i := 0; i < numRuns; i++ {
		WordCount(text)
	}
	runtimeMillis := time.Since(start).Nanoseconds() / 1e6

	return runtimeMillis
}

// Print the results of a benchmark
func printResults(runtimeMillis int64, numRuns int) {
	fmt.Printf("amount of runs: %d\n", numRuns)
	fmt.Printf("total time: %d ms\n", runtimeMillis)
	average := float64(runtimeMillis) / float64(numRuns)
	fmt.Printf("average time/run: %.2f ms\n", average)
}

func main() {
	// read in DataFile as a string called data
	data, err := ioutil.ReadFile(DataFile)
	if err != nil {
		panic(err)
	}

	numRuns := 100
	runtimeMillis := benchmark(string(data), numRuns)
	printResults(runtimeMillis, numRuns)
}
