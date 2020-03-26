package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strconv"

	"golang.org/x/sync/semaphore"
)

func main() {
	if len(os.Args) < 3 {
		panic("Usage: parallel script_file [script_arg script_arg...]")
	}

	script := os.Args[1]
	args := os.Args[2:]

	maxProcStr := os.Getenv("MAX_PROCS")
	maxProcs, err := strconv.Atoi(maxProcStr)
	if err != nil {
		maxProcs = 8
	}

	ctx := context.TODO()

	sem := semaphore.NewWeighted(int64(maxProcs))
	i := 0
	for _, a := range args {
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Printf("Failed to acquire semaphore: %v", err)
			break
		}

		go func(a string) {
			defer sem.Release(1)

			cmd := exec.Command(script, a)
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			i = i + 1
			if err != nil {
				log.Printf("Finished Unsuccessfully %d/%d", i, len(args))
			} else {
				log.Printf("Finished Successfully %d/%d", i, len(args))
			}
		}(a)
	}

	if err := sem.Acquire(ctx, int64(maxProcs)); err != nil {
		log.Printf("Failed to acquire semaphore: %v", err)
	}

}
