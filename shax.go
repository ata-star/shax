package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/minio/sha256-simd"
)

func main() {
	words, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer words.Close()

	hashes, err := os.Open(os.Args[2])
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer hashes.Close()

	hashList := make(map[string]bool)
	scanner := bufio.NewScanner(hashes)
	for scanner.Scan() {
		hashList[scanner.Text()] = true
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var wg sync.WaitGroup

	scanner = bufio.NewScanner(words)
	var found bool
	var lock sync.Mutex
	for scanner.Scan() && !found {
		wg.Add(1)
		go func(word string) {
			defer wg.Done()
			hash := sha256.Sum256([]byte(word))
			hashString := fmt.Sprintf("%x", hash)

			lock.Lock()
			if hashList[hashString] {
				fmt.Println("Match found:", word, "Hash:", hashString)
				delete(hashList, hashString)
				if len(hashList) == 0 {
					found = true
				}
			}
			lock.Unlock()
		}(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	wg.Wait()
}
