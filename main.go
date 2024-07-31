package main

import (
	"fmt"
	"multi-threaded--GREP/queue"
	"multi-threaded--GREP/worker"
	"os"
	"path/filepath"
	"sync"

	"github.com/alexflint/go-arg"
)

func discoverFolders(jobQueue *queue.Queue, path string) {
	entries, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("error reading folders:", err, entries)

		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			nextPath := filepath.Join(path, entry.Name())
			discoverFolders(jobQueue, nextPath)
		} else {
			jobQueue.Add(queue.NewJob(filepath.Join(path, entry.Name())))
		}
	}
}

var arguments struct {
	SearchQuery  string `arg:"positional,required"`
	SearchFolder string `arg:"positional"`
}

func main() {
	arg.MustParse(&arguments)

	var workersWaitGroup sync.WaitGroup
	const NumberOfWorkers = 10
	jobQueue := queue.NewQueue(70)
	results := make(chan worker.Result, 70)

	workersWaitGroup.Add(1)
	go func() {
		defer workersWaitGroup.Done()
		discoverFolders(&jobQueue, arguments.SearchFolder)
		jobQueue.Finalise(NumberOfWorkers)
	}()

	for i := 0; i < NumberOfWorkers; i++ {
		workersWaitGroup.Add(1)
		go func() {
			defer workersWaitGroup.Done()

			for {
				job := jobQueue.Next()

				if job.Path != "" {
					workerResults := worker.FindInFile(job.Path, arguments.SearchQuery)

					if workerResults != nil {
						for _, workerResult := range workerResults.Results {
							results <- workerResult
						}
					}
				} else {
					return
				}
			}
		}()
	}

	workersBlocker := make(chan int)
	go func() {
		workersWaitGroup.Wait()
		close(workersBlocker)
	}()

	var displayWaitGroup sync.WaitGroup

	displayWaitGroup.Add(1)
	go func() {
		for {
			select {
			case result := <-results:
				fmt.Printf("%v[%v]:%v\n", result.Path, result.LineNumber, result.Line)
			case <-workersBlocker:
				// ensure channel is empty before aborting display-goRoutine
				if len(results) == 0 {
					displayWaitGroup.Done()

					return
				}
			}
		}
	}()
	displayWaitGroup.Wait()
}