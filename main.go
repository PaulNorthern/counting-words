package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

const SUMMARY string = "SUMMARY"

var (
	words []string
)

type KeyValue struct {
	dict  map[string]int
	mutex sync.RWMutex
}

func main() {
	logger := slog.Default()
	t0 := time.Now()

	words = os.Args[1:]
	if len(words) == 0 {
		fmt.Println("Enter the words via arguments:")
		os.Exit(1)
	}
	//if len(words) != 5 {
	//	fmt.Println("5 слов в командной строке не было найдено, введите их через пробел:")
	//	reader, err := bufio.NewReader(os.Stdin).ReadString('\n')
	//	words = strings.Split(reader, " ")
	//	if err != nil {
	//		os.Exit(1)
	//	}
	//	if len(words) != 5 {
	//		os.Exit(1)
	//	}
	//}

	if err := run(); err != nil {
		logger.Error("%v", err)
		os.Exit(1)
	}

	logger.Info("Exit")
	t1 := time.Now()
	fmt.Printf("Elapsed time: %v", t1.Sub(t0))

}

func run() error {
	file, err := os.Open("input")
	if err != nil {
		return err
	}
	defer file.Close()

	//w, err := os.OpenFile("output", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	//if err != nil {
	//	return err
	//}
	//defer w.Close()

	countThreads := 2
	out := make(chan KeyValue, countThreads)

	wg := sync.WaitGroup{}
	scanChan := make(chan string)
	scanner := bufio.NewScanner(file)
	go func() {
		for scanner.Scan() {
			scanChan <- scanner.Text()
		}
		defer close(scanChan)
	}()

	for i := 0; i < countThreads; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			read(scanChan, out)
		}()
	}

	wg.Wait()
	close(out)
	printCountWords(out)

	return nil
}

func read(in <-chan string, out chan<- KeyValue) {
	keyValue := initStorage()
	for line := range in {
		parts := strings.Split(line, " ")
		for i := 0; i < len(parts); i++ {
			keyValue.mutex.Lock()
			for key := range keyValue.dict {
				if strings.Contains(strings.ToLower(parts[i]), key) {
					keyValue.dict[key]++
					keyValue.dict[SUMMARY]++
				}
			}
			keyValue.mutex.Unlock()
		}
	}
	out <- keyValue
}

func initStorage() KeyValue {
	keyValue := KeyValue{
		dict:  make(map[string]int),
		mutex: sync.RWMutex{},
	}
	for _, str := range words {
		keyValue.dict[strings.ToLower(str)] = 0
	}
	keyValue.dict[SUMMARY] = 0
	return keyValue
}

func printCountWords(outChannels ...<-chan KeyValue) {
	merged := make(chan KeyValue)
	done := make(chan struct{})

	for _, out := range outChannels {
		go func(ch <-chan KeyValue) {
			defer close(merged)
			for kv := range ch {
				merged <- kv
			}
			done <- struct{}{}
		}(out)
	}

	go func() {
		for i := 0; i < len(outChannels); i++ {
			<-done
		}
		close(done)
	}()

	totalResult := KeyValue{
		dict:  make(map[string]int),
		mutex: sync.RWMutex{},
	}

	for kv := range merged {
		for key, count := range kv.dict {
			totalResult.mutex.Lock()
			totalResult.dict[key] += count
			totalResult.mutex.Unlock()
		}
	}

	for _, v := range words {
		count, ok := totalResult.dict[v]
		if ok {
			fmt.Printf("'%s': %d \n", v, count)
		}
	}
	fmt.Printf("'%s': %d \n", "SUMMARY", totalResult.dict[SUMMARY])
}
